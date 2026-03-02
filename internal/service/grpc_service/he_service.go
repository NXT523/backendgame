package grpc_service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"game/ent"
	mappers "game/internal/mapping/ent"
	grpc_repository "game/internal/repository/grpc_repository"
	"game/internal/utils"
	"game/pkg/cache"
	"game/pkg/kafka"
	v1 "game/v1/proto"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type HeService struct {
	grpc_repository *grpc_repository.HeRepo
	ent             *ent.Client
	rdb             *redis.Client
	mc              *memcache.Client
	cache           cache.RedisCacheService
	kw              *kafka.KafkaWriters
	lock            *utils.RedisLock
	ttl             time.Duration
	es              *elasticsearch.Client
	batchQueue      chan *updateJob
}

func NewHeService(
	ent *ent.Client,
	rdb *redis.Client,
	lock *utils.RedisLock,
) *HeService {
	return &HeService{
		grpc_repository: grpc_repository.NewHeRepo(ent),
		ent:             ent,
		rdb:             rdb,
		cache:           cache.NewRedisCacheService(rdb),
		lock:            lock,
		batchQueue:      make(chan *updateJob, 5000),
	}
}

func (s *HeService) CreateHeService(ctx context.Context, req *v1.TaoHeRequest) (*v1.He, error) {

	name := utils.NormalizeKeepCase(req.TenHe)
	normalized := utils.NormalizeString(name)

	lockKey := "lock:he:create:" + normalized
	var out *v1.He

	err := s.lock.WithDistributedLock(ctx, lockKey, 10*time.Second, func(_ string) error {
		tx, err := s.ent.Tx(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		row, err := s.grpc_repository.CreateHeRepo(ctx, tx, req, name)
		if err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				return status.Error(codes.AlreadyExists, "Tên hệ đã tồn tại")
			}
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		full, _ := s.grpc_repository.GetByIDHeRepo(
			context.Background(), row.ID,
		)
		out = mappers.HeEntToPB(full)
		return nil
	})

	if err != nil {
		return nil, err
	}
	if s.rdb != nil {
		_ = s.rdb.Incr(ctx, "he:search:version").Err()
		go s.invalidateHeCache(context.Background())
	}

	return out, nil
}

func (s *HeService) UpdateByNameHeService(
	ctx context.Context,
	req *v1.CapNhatTheoTenHeRequest,
) (*v1.He, error) {

	name := strings.TrimSpace(req.TenHe)
	lockKey := "lock:he:name:" + utils.NormalizeString(name)
	var out *v1.He
	err := s.lock.WithDistributedLock(ctx, lockKey, 10*time.Second, func(_ string) error {

		tx, err := s.ent.Tx(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		cur, err := s.grpc_repository.LockByNameHeRepo(ctx, tx, name)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Error(codes.NotFound, "Không tìm thấy hệ")
			}
			return err
		}

		updated, err := s.grpc_repository.UpdateHeRepo(
			ctx, tx, cur.ID, req,
		)

		if err := tx.Commit(); err != nil {
			return err
		}

		full, _ := s.grpc_repository.GetByIDHeRepo(ctx, updated.ID)
		out = mappers.HeEntToPB(full)
		return nil
	})

	if err != nil {
		return nil, err
	}

	if s.rdb != nil {
		_ = s.rdb.Incr(ctx, "he:search:version").Err()
		go s.invalidateHeCache(context.Background())
	}

	return out, nil

}

func (s *HeService) DeleteByNameHeService(
	ctx context.Context,
	req *v1.XoaTheoTenHeRequest,
) (*v1.XoaTheoTenHeResponse, error) {

	name := strings.TrimSpace(req.TenHe)
	lockKey := "lock:he:name:" + utils.NormalizeString(name)
	var cur *ent.He

	err := s.lock.WithDistributedLock(ctx, lockKey, 10*time.Second, func(_ string) error {

		tx, err := s.ent.Tx(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		cur, err = s.grpc_repository.LockByNameHeRepo(ctx, tx, name)
		if err != nil {
			return err
		}
		if cur == nil {
			return status.Error(codes.NotFound, "Không tìm thấy hệ")
		}

		if err := s.grpc_repository.DeleteByIDHeRepo(ctx, tx, cur.ID); err != nil {
			return err
		}

		return tx.Commit()
	})

	if err != nil {
		return nil, err
	}

	if cur == nil {
		return nil, status.Error(codes.Internal, "Xóa thất bại, dữ liệu không hợp lệ")
	}

	_ = s.rdb.Incr(ctx, "he:search:version").Err()
	go s.invalidateHeCache(context.Background())

	return &v1.XoaTheoTenHeResponse{
		DeletedName: name,
		MaHe:        int32(cur.ID),
	}, nil

}

func (s *HeService) GetAllHeService(ctx context.Context) (*v1.DanhSachHe, error) {
	cacheKey := "he:get_all_he"

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachHe{}
	}); ok {
		return msg.(*v1.DanhSachHe), nil
	}

	rows, err := s.grpc_repository.GetAllHeRepo(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]*v1.He, 0, len(rows))
	for _, r := range rows {
		items = append(items, mappers.HeEntToPB(r))
	}

	resp := &v1.DanhSachHe{
		Items: items,
		Total: int32(len(items)),
	}

	_ = s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute)
	return resp, nil
}

func (s *HeService) GetAllTenHeService(
	ctx context.Context,
) (*v1.DanhSachStringHe, error) {
	cacheKey := "he:get_all_ten_he"

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachStringHe{}
	}); ok {
		return msg.(*v1.DanhSachStringHe), nil
	}

	names, err := s.grpc_repository.GetAllTenHeRepo(ctx)
	if err != nil {
		return nil, err
	}

	resp := &v1.DanhSachStringHe{
		Total: int32(len(names)),
		Items: names,
	}

	_ = s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute)

	return resp, nil
}

func (s *HeService) GetAllMoTaHeService(
	ctx context.Context,
) (*v1.DanhSachStringHe, error) {
	cacheKey := "he:get_all_mo_ta_he"

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachStringHe{}
	}); ok {
		return msg.(*v1.DanhSachStringHe), nil
	}

	names, err := s.grpc_repository.GetAllMoTaHeRepo(ctx)
	if err != nil {
		return nil, err
	}

	resp := &v1.DanhSachStringHe{
		Total: int32(len(names)),
		Items: names,
	}

	_ = s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute)

	return resp, nil
}

func (s *HeService) SearchHeService(ctx context.Context, req *v1.TimKiemHeRequest) ([]*v1.He, int64, error) {
	limit := normalizeLimit(req.Limit)
	offset, err := parseCursor(req.Cursor)
	if err != nil {
		return nil, 0, status.Error(codes.InvalidArgument, err.Error())
	}

	ver, err := s.rdb.Get(ctx, "he:search:version").Int()
	if err != nil {
		ver = 1
	}

	cp := proto.Clone(req).(*v1.TimKiemHeRequest)
	cp.Limit = int32(limit)
	cp.Cursor = strconv.Itoa(offset)

	raw, _ := protojson.Marshal(cp)
	sum := sha1.Sum(raw)

	cacheKey := fmt.Sprintf(
		"he:search:v%d:%s",
		ver,
		hex.EncodeToString(sum[:]),
	)

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachHe{}
	}); ok {
		resp := msg.(*v1.DanhSachHe)
		return resp.Items, int64(resp.Total), nil
	}

	ents, total, err := s.grpc_repository.SearchHeRepo(ctx, req, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return nil, 0, status.Error(codes.NotFound, "Không tìm thấy hệ")
	}

	items := make([]*v1.He, 0, len(ents))
	for _, e := range ents {
		items = append(items, mappers.HeEntToPB(e))
	}

	resp := &v1.DanhSachHe{
		Total:     int32(total),
		ItemCount: int32(len(items)),
		Items:     items,
	}

	if err := s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute); err != nil {
		log.Println("Redis lỗi:", err)
	}

	return resp.Items, int64(resp.Total), nil
}
