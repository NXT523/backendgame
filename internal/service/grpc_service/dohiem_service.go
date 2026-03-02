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
	"google.golang.org/protobuf/types/known/emptypb"
)

type DoHiemService struct {
	grpc_repository *grpc_repository.DoHiemRepo
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

func NewDoHiemService(
	ent *ent.Client,
	rdb *redis.Client,
	lock *utils.RedisLock,
) *DoHiemService {
	return &DoHiemService{
		grpc_repository: grpc_repository.NewDoHiemRepo(ent),
		ent:             ent,
		rdb:             rdb,
		cache:           cache.NewRedisCacheService(rdb),
		lock:            lock,
		batchQueue:      make(chan *updateJob, 5000),
	}
}

func (s *DoHiemService) CreateDoHiemService(ctx context.Context, req *v1.TaoDoHiemRequest) (*v1.DoHiem, error) {

	name := utils.NormalizeKeepCase(req.TenDoHiem)
	normalized := utils.NormalizeString(name)

	lockKey := "lock:dohiem:create:" + normalized
	var out *v1.DoHiem

	err := s.lock.WithDistributedLock(ctx, lockKey, 10*time.Second, func(_ string) error {
		tx, err := s.ent.Tx(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		row, err := s.grpc_repository.CreateDoHiemRepo(ctx, tx, req, name)
		if err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				return status.Error(codes.AlreadyExists, "Tên độ hiếm đã tồn tại")
			}
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		full, _ := s.grpc_repository.GetByIDDoHiemRepo(
			context.Background(), row.ID,
		)
		out = mappers.DoHiemEntToPB(full)
		return nil
	})

	if err != nil {
		return nil, err
	}
	if s.rdb != nil {
		_ = s.rdb.Incr(ctx, "dohiem:search:version").Err()
		go s.invalidateDoHiemCache(context.Background())
	}

	return out, nil
}

func (s *DoHiemService) UpdateByNameDoHiemService(
	ctx context.Context,
	req *v1.CapNhatTheoTenDoHiemRequest,
) (*v1.DoHiem, error) {

	name := strings.TrimSpace(req.TenDoHiem)
	lockKey := "lock:dohiem:name:" + utils.NormalizeString(name)
	var out *v1.DoHiem
	err := s.lock.WithDistributedLock(ctx, lockKey, 10*time.Second, func(_ string) error {

		tx, err := s.ent.Tx(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		cur, err := s.grpc_repository.LockByNameDoHiemRepo(ctx, tx, name)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Error(codes.NotFound, "Không tìm thấy độ hiếm")
			}
			return err
		}

		updated, err := s.grpc_repository.UpdateDoHiemRepo(
			ctx, tx, cur.ID, req,
		)

		if err := tx.Commit(); err != nil {
			return err
		}

		full, _ := s.grpc_repository.GetByIDDoHiemRepo(ctx, updated.ID)
		out = mappers.DoHiemEntToPB(full)
		return nil
	})

	if err != nil {
		return nil, err
	}

	if s.rdb != nil {
		_ = s.rdb.Incr(ctx, "dohiem:search:version").Err()
		go s.invalidateDoHiemCache(context.Background())
	}

	return out, nil

}

func (s *DoHiemService) DeleteByNameDoHiemService(
	ctx context.Context,
	req *v1.XoaTheoTenDoHiemRequest,
) (*v1.XoaTheoTenDoHiemResponse, error) {

	name := strings.TrimSpace(req.TenDoHiem)
	lockKey := "lock:dohiem:name:" + utils.NormalizeString(name)
	var cur *ent.DoHiem

	err := s.lock.WithDistributedLock(ctx, lockKey, 10*time.Second, func(_ string) error {

		tx, err := s.ent.Tx(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		cur, err = s.grpc_repository.LockByNameDoHiemRepo(ctx, tx, name)
		if err != nil {
			return err
		}
		if cur == nil {
			return status.Error(codes.NotFound, "Không tìm thấy hệ")
		}

		if err := s.grpc_repository.DeleteByIDDoHiemRepo(ctx, tx, cur.ID); err != nil {
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

	_ = s.rdb.Incr(ctx, "dohiem:search:version").Err()
	go s.invalidateDoHiemCache(context.Background())

	return &v1.XoaTheoTenDoHiemResponse{
		DeletedName: name,
		MaDoHiem:    int32(cur.ID),
	}, nil

}

func (s *DoHiemService) GetAllDoHiemService(ctx context.Context) (*v1.DanhSachDoHiem, error) {
	cacheKey := "dohiem:get_all_dohiem"

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachDoHiem{}
	}); ok {
		return msg.(*v1.DanhSachDoHiem), nil
	}

	rows, err := s.grpc_repository.GetAllDoHiemRepo(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]*v1.DoHiem, 0, len(rows))
	for _, r := range rows {
		items = append(items, mappers.DoHiemEntToPB(r))
	}

	resp := &v1.DanhSachDoHiem{
		Items: items,
		Total: int32(len(items)),
	}

	_ = s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute)
	return resp, nil
}

func (s *DoHiemService) GetAllTenDoHiemService(
	ctx context.Context,
) (*v1.DanhSachStringDoHiem, error) {
	cacheKey := "dohiem:get_all_ten_dohiem"

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachStringDoHiem{}
	}); ok {
		return msg.(*v1.DanhSachStringDoHiem), nil
	}

	names, err := s.grpc_repository.GetAllTenDoHiemRepo(ctx)
	if err != nil {
		return nil, err
	}

	resp := &v1.DanhSachStringDoHiem{
		Total: int32(len(names)),
		Items: names,
	}

	_ = s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute)

	return resp, nil
}

func (s *DoHiemService) GetAllSoLuongDoHiemService(
	ctx context.Context,
	_ *emptypb.Empty,
) (*v1.DanhSachIntDoHiem, error) {

	cacheKey := "dohiem:get_all_so_luong_dohiem"

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachIntDoHiem{}
	}); ok {
		return msg.(*v1.DanhSachIntDoHiem), nil
	}

	items, err := s.grpc_repository.GetAllSoLuongDoHiemRepo(ctx)
	if err != nil {
		return nil, err
	}

	resp := &v1.DanhSachIntDoHiem{
		Total: int32(len(items)),
		Items: items,
	}

	_ = s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute)
	return resp, nil
}

func (s *DoHiemService) GetAllMauSacDoHiemService(
	ctx context.Context,
) (*v1.DanhSachStringDoHiem, error) {
	cacheKey := "dohiem:get_all_mau_sac_dohiem"

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachStringDoHiem{}
	}); ok {
		return msg.(*v1.DanhSachStringDoHiem), nil
	}

	names, err := s.grpc_repository.GetAllMauSacDoHiemRepo(ctx)
	if err != nil {
		return nil, err
	}

	resp := &v1.DanhSachStringDoHiem{
		Total: int32(len(names)),
		Items: names,
	}

	_ = s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute)

	return resp, nil
}

func (s *DoHiemService) GetAllSatThuongBonusDoHiemService(
	ctx context.Context,
	_ *emptypb.Empty,
) (*v1.DanhSachDoubleDoHiem, error) {

	cacheKey := "dohiem:get_all_sat_thuong_bonus_dohiem"

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachDoubleDoHiem{}
	}); ok {
		return msg.(*v1.DanhSachDoubleDoHiem), nil
	}

	items, err := s.grpc_repository.GetAllSatThuongBonusDoHiemRepo(ctx)
	if err != nil {
		return nil, err
	}

	resp := &v1.DanhSachDoubleDoHiem{
		Total: int32(len(items)),
		Items: items,
	}

	_ = s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute)
	return resp, nil
}

func (s *DoHiemService) GetAllTocDoDanhBonusDoHiemService(
	ctx context.Context,
	_ *emptypb.Empty,
) (*v1.DanhSachDoubleDoHiem, error) {

	cacheKey := "dohiem:get_all_toc_do_danh_bonus_dohiem"

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachDoubleDoHiem{}
	}); ok {
		return msg.(*v1.DanhSachDoubleDoHiem), nil
	}

	items, err := s.grpc_repository.GetAllTocDoDanhBonusDoHiemRepo(ctx)
	if err != nil {
		return nil, err
	}

	resp := &v1.DanhSachDoubleDoHiem{
		Total: int32(len(items)),
		Items: items,
	}

	_ = s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute)
	return resp, nil
}

func (s *DoHiemService) GetAllCapBacDoHiemService(
	ctx context.Context,
	_ *emptypb.Empty,
) (*v1.DanhSachIntDoHiem, error) {

	cacheKey := "dohiem:get_all_cap_bac_dohiem"

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachIntDoHiem{}
	}); ok {
		return msg.(*v1.DanhSachIntDoHiem), nil
	}

	items, err := s.grpc_repository.GetAllCapBacDoHiemRepo(ctx)
	if err != nil {
		return nil, err
	}

	resp := &v1.DanhSachIntDoHiem{
		Total: int32(len(items)),
		Items: items,
	}

	_ = s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute)
	return resp, nil
}

func (s *DoHiemService) SearchDoHiemService(ctx context.Context, req *v1.TimKiemDoHiemRequest) ([]*v1.DoHiem, int64, error) {
	limit := normalizeLimit(req.Limit)
	offset, err := parseCursor(req.Cursor)
	if err != nil {
		return nil, 0, status.Error(codes.InvalidArgument, err.Error())
	}

	ver, err := s.rdb.Get(ctx, "dohiem:search:version").Int()
	if err != nil {
		ver = 1
	}

	cp := proto.Clone(req).(*v1.TimKiemDoHiemRequest)
	cp.Limit = int32(limit)
	cp.Cursor = strconv.Itoa(offset)

	raw, _ := protojson.Marshal(cp)
	sum := sha1.Sum(raw)

	cacheKey := fmt.Sprintf(
		"dohiem:search:v%d:%s",
		ver,
		hex.EncodeToString(sum[:]),
	)

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachDoHiem{}
	}); ok {
		resp := msg.(*v1.DanhSachDoHiem)
		return resp.Items, int64(resp.Total), nil
	}

	ents, total, err := s.grpc_repository.SearchDoHiemRepo(ctx, req, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return nil, 0, status.Error(codes.NotFound, "Không tìm thấy độ hiếm")
	}

	items := make([]*v1.DoHiem, 0, len(ents))
	for _, e := range ents {
		items = append(items, mappers.DoHiemEntToPB(e))
	}

	resp := &v1.DanhSachDoHiem{
		Total:     int32(total),
		ItemCount: int32(len(items)),
		Items:     items,
	}

	if err := s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute); err != nil {
		log.Println("Redis lỗi:", err)
	}

	return resp.Items, int64(resp.Total), nil
}
