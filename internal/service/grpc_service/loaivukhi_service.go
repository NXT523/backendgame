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
	v1 "game/v1/proto"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type LoaiVuKhiService struct {
	grpc_repository *grpc_repository.LoaiVuKhiRepo
	ent             *ent.Client
	lock            *utils.RedisLock
	rdb             *redis.Client
	cache           cache.RedisCacheService
}

func NewLoaiVuKhiService(
	ent *ent.Client,
	rdb *redis.Client,
	lock *utils.RedisLock,
) *LoaiVuKhiService {
	return &LoaiVuKhiService{
		grpc_repository: grpc_repository.NewLoaiVuKhiRepo(ent),
		ent:             ent,
		lock:            lock,
		rdb:             rdb,
		cache:           cache.NewRedisCacheService(rdb),
	}
}

func (s *LoaiVuKhiService) CreateLoaiVuKhiService(
	ctx context.Context,
	req *v1.TaoLoaiVuKhiRequest,
) (*v1.LoaiVuKhi, error) {

	name := utils.NormalizeKeepCase(req.TenLoaiVuKhi)
	normalized := utils.NormalizeString(name)
	lockKey := "lock:loaivukhi:create:" + normalized

	var out *v1.LoaiVuKhi

	err := s.lock.WithDistributedLock(ctx, lockKey, 10*time.Second, func(_ string) error {

		tx, err := s.ent.Tx(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		row, err := s.grpc_repository.CreateLoaiVuKhiRepo(ctx, tx, req, name)
		if err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				return status.Error(codes.AlreadyExists, "Tên loại vũ khí đã tồn tại")
			}
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		full, _ := s.grpc_repository.GetByIDLoaiVuKhiRepo(context.Background(), row.ID)
		out = mappers.LoaiVuKhiEntToPB(full)

		return nil
	})

	if err != nil {
		return nil, err
	}

	if s.rdb != nil {
		_ = s.rdb.Incr(ctx, "loaivukhi:search:version").Err()
		go s.invalidateLoaiVuKhiCache(context.Background())
	}

	return out, nil
}

func (s *LoaiVuKhiService) UpdateByNameLoaiVuKhiService(
	ctx context.Context,
	req *v1.CapNhatTheoTenLoaiVuKhiRequest,
) (*v1.LoaiVuKhi, error) {

	name := strings.TrimSpace(req.TenLoaiVuKhi)
	lockKey := "lock:loaivukhi:name:" + utils.NormalizeString(name)

	var out *v1.LoaiVuKhi

	err := s.lock.WithDistributedLock(ctx, lockKey, 10*time.Second, func(_ string) error {

		tx, err := s.ent.Tx(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		cur, err := s.grpc_repository.LockByNameLoaiVuKhiRepo(ctx, tx, name)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Error(codes.NotFound, "Không tìm thấy loại vũ khí")
			}
			return err
		}

		updated, err := s.grpc_repository.UpdateLoaiVuKhiRepo(ctx, tx, cur.ID, req)
		if err != nil {
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		full, _ := s.grpc_repository.GetByIDLoaiVuKhiRepo(ctx, updated.ID)
		out = mappers.LoaiVuKhiEntToPB(full)

		return nil
	})

	if err != nil {
		return nil, err
	}

	if s.rdb != nil {
		_ = s.rdb.Incr(ctx, "loaivukhi:search:version").Err()
		go s.invalidateLoaiVuKhiCache(context.Background())
	}

	return out, nil
}

func (s *LoaiVuKhiService) DeleteByNameLoaiVuKhiService(
	ctx context.Context,
	req *v1.XoaTheoTenLoaiVuKhiRequest,
) (*v1.XoaTheoTenLoaiVuKhiResponse, error) {

	name := strings.TrimSpace(req.TenLoaiVuKhi)
	lockKey := "lock:loaivukhi:name:" + utils.NormalizeString(name)

	var cur *ent.LoaiVuKhi

	err := s.lock.WithDistributedLock(ctx, lockKey, 10*time.Second, func(_ string) error {

		tx, err := s.ent.Tx(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		cur, err = s.grpc_repository.LockByNameLoaiVuKhiRepo(ctx, tx, name)
		if err != nil {
			return err
		}

		if cur == nil {
			return status.Error(codes.NotFound, "Không tìm thấy loại vũ khí")
		}

		if err := s.grpc_repository.DeleteByIDLoaiVuKhiRepo(ctx, tx, cur.ID); err != nil {
			return err
		}

		return tx.Commit()
	})

	if err != nil {
		return nil, err
	}

	_ = s.rdb.Incr(ctx, "loaivukhi:search:version").Err()
	go s.invalidateLoaiVuKhiCache(context.Background())

	return &v1.XoaTheoTenLoaiVuKhiResponse{
		DeletedName: name,
		MaLoaiVuKhi: int32(cur.ID),
	}, nil
}

func (s *LoaiVuKhiService) GetAllLoaiVuKhiService(ctx context.Context) (*v1.DanhSachLoaiVuKhi, error) {

	cacheKey := "loaivukhi:get_all_loaivukhi"

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachLoaiVuKhi{}
	}); ok {
		return msg.(*v1.DanhSachLoaiVuKhi), nil
	}

	rows, err := s.grpc_repository.GetAllLoaiVuKhiRepo(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]*v1.LoaiVuKhi, 0, len(rows))
	for _, r := range rows {
		items = append(items, mappers.LoaiVuKhiEntToPB(r))
	}

	resp := &v1.DanhSachLoaiVuKhi{
		Total: int32(len(items)),
		Items: items,
	}

	_ = s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute)
	return resp, nil
}

func (s *LoaiVuKhiService) GetAllTenLoaiVuKhiService(
	ctx context.Context,
) (*v1.DanhSachStringLoaiVuKhi, error) {
	cacheKey := "loaivukhi:get_all_ten_loaivukhi"

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachStringLoaiVuKhi{}
	}); ok {
		return msg.(*v1.DanhSachStringLoaiVuKhi), nil
	}

	names, err := s.grpc_repository.GetAllTenLoaiVuKhiRepo(ctx)
	if err != nil {
		return nil, err
	}

	resp := &v1.DanhSachStringLoaiVuKhi{
		Total: int32(len(names)),
		Items: names,
	}

	_ = s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute)

	return resp, nil
}

func (s *LoaiVuKhiService) GetAllMoTaLoaiVuKhiService(
	ctx context.Context,
) (*v1.DanhSachStringLoaiVuKhi, error) {
	cacheKey := "loaivukhi:get_all_mo_ta_loaivukhi"

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachStringLoaiVuKhi{}
	}); ok {
		return msg.(*v1.DanhSachStringLoaiVuKhi), nil
	}

	names, err := s.grpc_repository.GetAllMoTaLoaiVuKhiRepo(ctx)
	if err != nil {
		return nil, err
	}

	resp := &v1.DanhSachStringLoaiVuKhi{
		Total: int32(len(names)),
		Items: names,
	}

	_ = s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute)

	return resp, nil
}

func (s *LoaiVuKhiService) SearchLoaiVuKhiService(
	ctx context.Context,
	req *v1.TimKiemLoaiVuKhiRequest,
) ([]*v1.LoaiVuKhi, int64, error) {

	limit := normalizeLimit(req.Limit)
	offset, err := parseCursor(req.Cursor)
	if err != nil {
		return nil, 0, status.Error(codes.InvalidArgument, err.Error())
	}

	ver, err := s.rdb.Get(ctx, "loaivukhi:search:version").Int()
	if err != nil {
		ver = 1
	}

	cp := proto.Clone(req).(*v1.TimKiemLoaiVuKhiRequest)
	cp.Limit = int32(limit)
	cp.Cursor = strconv.Itoa(offset)

	raw, _ := protojson.Marshal(cp)
	sum := sha1.Sum(raw)

	cacheKey := fmt.Sprintf(
		"loaivukhi:search:v%d:%s",
		ver,
		hex.EncodeToString(sum[:]),
	)

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachLoaiVuKhi{}
	}); ok {
		resp := msg.(*v1.DanhSachLoaiVuKhi)
		return resp.Items, int64(resp.Total), nil
	}

	ents, total, err := s.grpc_repository.SearchLoaiVuKhiRepo(ctx, req, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return nil, 0, status.Error(codes.NotFound, "Không tìm thấy loại vũ khí")
	}

	items := make([]*v1.LoaiVuKhi, 0, len(ents))
	for _, e := range ents {
		items = append(items, mappers.LoaiVuKhiEntToPB(e))
	}

	resp := &v1.DanhSachLoaiVuKhi{
		Total:     int32(total),
		ItemCount: int32(len(items)),
		Items:     items,
	}

	if err := s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute); err != nil {
		log.Println("Redis lỗi:", err)
	}

	return resp.Items, int64(resp.Total), nil
}
