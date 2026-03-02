package grpc_service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"game/ent"
	mappers "game/internal/mapping/ent"
	grpc_repository "game/internal/repository/grpc_repository"
	"game/internal/utils"
	"game/pkg/cache"
	v1 "game/v1/proto"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type jobResult struct {
	Index int
	TenVuKhiNew  string
	Err   error
}

type updateJob struct {
	Index   int
	req     *v1.CapNhatTheoTenVuKhiRequest
	resChan chan jobResult
}

type VuKhiService struct {
	grpc_repository *grpc_repository.VuKhiRepo
	ent             *ent.Client
	rdb             *redis.Client
	cache           cache.RedisCacheService
	lock            *utils.RedisLock
	rsync           *redsync.Redsync
	batchQueue      chan *updateJob
}

func NewVuKhiService(
	ent *ent.Client,
	rdb *redis.Client,
	lock *utils.RedisLock,
) *VuKhiService {
	s := &VuKhiService{
		grpc_repository: grpc_repository.NewVuKhiRepo(ent),
		ent:             ent,
		rdb:             rdb,
		cache:           cache.NewRedisCacheService(rdb),
		lock:            lock,
		batchQueue:      make(chan *updateJob, 5000),
	}

	for i := 0; i < 20; i++ {
		go s.startBatchWorker()
	}

	return s
}
func (s *VuKhiService) CreateVuKhiService(
	ctx context.Context,
	req *v1.TaoVuKhiRequest,
) (*v1.VuKhi, error) {

	name := utils.NormalizeKeepCase(req.TenVuKhi)
	normalized := utils.NormalizeString(name)

	lockKey := "lock:vukhi:create:" + normalized
	var out *v1.VuKhi

	err := s.lock.WithDistributedLock(ctx, lockKey, 10*time.Second, func(_ string) error {
		tx, err := s.ent.Tx(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		row, err := s.grpc_repository.CreateVuKhiRepo(ctx, tx, req, name)
		if err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				return status.Error(codes.AlreadyExists, "Tên vũ khí đã tồn tại")
			}
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		full, _ := s.grpc_repository.GetByIDVuKhiRepo(
			context.Background(), row.ID,
		)
		out = mappers.VuKhiEntToPB(full)
		return nil
	})

	if err != nil {
		return nil, err
	}
	go s.invalidateVuKhiCache(context.Background())

	return out, nil
}

func (s *VuKhiService) UpdateByNameVuKhiService(
	ctx context.Context,
	req *v1.CapNhatTheoTenVuKhiRequest,
) (*v1.VuKhi, error) {

	name := strings.TrimSpace(req.TenVuKhi)
	lockKey := "lock:vukhi:name:" + utils.NormalizeString(name)
	var out *v1.VuKhi
	err := s.lock.WithDistributedLock(ctx, lockKey, 10*time.Second, func(_ string) error {

		tx, err := s.ent.Tx(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		cur, err := s.grpc_repository.LockByNameVuKhiRepo(ctx, tx, name)
		if err != nil {
			if ent.IsNotFound(err) {
				return status.Error(codes.NotFound, "Không tìm thấy vũ khí")
			}
			return err
		}

		if cur == nil {
			return status.Error(codes.NotFound, "Không tìm thấy vũ khí")
		}

		if int32(cur.Version) != req.Version {
			return status.Error(codes.FailedPrecondition, "version không khớp")
		}

		updated, err := s.grpc_repository.UpdateVuKhiRepo(
			ctx, tx, cur.ID, req, cur.Version+1,
		)
		if err != nil {
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		full, _ := s.grpc_repository.GetByIDVuKhiRepo(ctx, updated.ID)
		out = mappers.VuKhiEntToPB(full)
		return nil
	})

	if err != nil {
		return nil, err
	}

	go s.invalidateVuKhiCache(context.Background())

	return out, nil

}

func (s *VuKhiService) DeleteByNameVuKhiService(
	ctx context.Context,
	req *v1.XoaTheoTenVuKhiRequest,
) (*v1.XoaTheoTenVuKhiResponse, error) {

	name := strings.TrimSpace(req.TenVuKhi)
	lockKey := "lock:vukhi:name:" + utils.NormalizeString(name)
	var cur *ent.VuKhi

	err := s.lock.WithDistributedLock(ctx, lockKey, 10*time.Second, func(_ string) error {

		tx, err := s.ent.Tx(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		cur, err = s.grpc_repository.LockByNameVuKhiRepo(ctx, tx, name)
		if err != nil {
			return err
		}
		if cur == nil {
			return status.Error(codes.NotFound, "Không tìm thấy vũ khí")
		}

		if err := s.grpc_repository.DeleteByIDVuKhiRepo(ctx, tx, cur.ID); err != nil {
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

	go s.invalidateVuKhiCache(context.Background())

	return &v1.XoaTheoTenVuKhiResponse{
		DeletedName: name,
		MaVuKhi:     int32(cur.ID),
	}, nil

}

func (s *VuKhiService) GetAllVuKhiService(ctx context.Context) (*v1.DanhSachVuKhi, error) {
	cacheKey := "vukhi:get_all_vukhi"

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachVuKhi{}
	}); ok {
		return msg.(*v1.DanhSachVuKhi), nil
	}

	rows, err := s.grpc_repository.GetAllVuKhiRepo(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]*v1.VuKhi, 0, len(rows))
	for _, r := range rows {
		items = append(items, mappers.VuKhiEntToPB(r))
	}

	resp := &v1.DanhSachVuKhi{
		Items: items,
		Total: int32(len(items)),
	}

	_ = s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute)
	return resp, nil
}

func (s *VuKhiService) GetAllTenVuKhiService(
	ctx context.Context,
) (*v1.DanhSachStringVuKhi, error) {
	cacheKey := "vukhi:get_all_ten_vukhi"

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachStringVuKhi{}
	}); ok {
		return msg.(*v1.DanhSachStringVuKhi), nil
	}

	names, err := s.grpc_repository.GetAllTenVuKhiRepo(ctx)
	if err != nil {
		return nil, err
	}

	resp := &v1.DanhSachStringVuKhi{
		Total: int32(len(names)),
		Items: names,
	}

	_ = s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute)

	return resp, nil
}

func (s *VuKhiService) GetAllSatThuongCoBanVuKhiService(
	ctx context.Context,
	_ *emptypb.Empty,
) (*v1.DanhSachIntVuKhi, error) {

	cacheKey := "vukhi:get_all_sat_thuong_co_ban_vukhi"

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachIntVuKhi{}
	}); ok {
		return msg.(*v1.DanhSachIntVuKhi), nil
	}

	items, err := s.grpc_repository.GetAllSatThuongCoBanVuKhiRepo(ctx)
	if err != nil {
		return nil, err
	}

	resp := &v1.DanhSachIntVuKhi{
		Total: int32(len(items)),
		Items: items,
	}

	_ = s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute)
	return resp, nil
}

func (s *VuKhiService) GetAllTamDanhVuKhiService(
	ctx context.Context,
	_ *emptypb.Empty,
) (*v1.DanhSachIntVuKhi, error) {

	cacheKey := "vukhi:get_all_tam_danh_vukhi"

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachIntVuKhi{}
	}); ok {
		return msg.(*v1.DanhSachIntVuKhi), nil
	}

	items, err := s.grpc_repository.GetAllTamDanhVuKhiRepo(ctx)
	if err != nil {
		return nil, err
	}

	resp := &v1.DanhSachIntVuKhi{
		Total: int32(len(items)),
		Items: items,
	}

	_ = s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute)
	return resp, nil
}

func (s *VuKhiService) GetAllTocDoDanhVuKhiService(
	ctx context.Context,
	_ *emptypb.Empty,
) (*v1.DanhSachDoubleVuKhi, error) {

	cacheKey := "vukhi:get_all_toc_do_vukhi"

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachDoubleVuKhi{}
	}); ok {
		return msg.(*v1.DanhSachDoubleVuKhi), nil
	}

	items, err := s.grpc_repository.GetAllTocDoDanhVuKhiRepo(ctx)
	if err != nil {
		return nil, err
	}

	resp := &v1.DanhSachDoubleVuKhi{
		Total: int32(len(items)),
		Items: items,
	}

	_ = s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute)
	return resp, nil
}

func (s *VuKhiService) SearchVuKhiService(
	ctx context.Context,
	req *v1.TimKiemVuKhiRequest,
) ([]*v1.VuKhi, int64, error) {

	limit := normalizeLimit(req.Limit)
	offset, err := parseCursor(req.Cursor)
	if err != nil {
		return nil, 0, status.Error(codes.InvalidArgument, err.Error())
	}

	ver, err := s.rdb.Get(ctx, "vukhi:search:version").Int()
	if err != nil {
		ver = 1
	}

	cp := proto.Clone(req).(*v1.TimKiemVuKhiRequest)
	cp.Limit = int32(limit)
	cp.Cursor = strconv.Itoa(offset)

	raw, _ := protojson.Marshal(cp)
	sum := sha1.Sum(raw)

	cacheKey := fmt.Sprintf(
		"vukhi:search:v%d:%s",
		ver,
		hex.EncodeToString(sum[:]),
	)

	if msg, ok := s.cache.GetProto(ctx, cacheKey, func() proto.Message {
		return &v1.DanhSachVuKhi{}
	}); ok {
		resp := msg.(*v1.DanhSachVuKhi)
		return resp.Items, int64(resp.Total), nil
	}

	ents, total, err := s.grpc_repository.SearchVuKhiRepo(ctx, req, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return nil, 0, status.Error(codes.NotFound, "Không tìm thấy vũ khí")
	}

	items := make([]*v1.VuKhi, 0, len(ents))
	for _, e := range ents {
		items = append(items, mappers.VuKhiEntToPB(e))
	}

	resp := &v1.DanhSachVuKhi{
		Total:      int32(total),
		ItemCount:  int32(len(items)),
		Items:      items,
		HasMore:    len(items) == limit,
		NextCursor: strconv.Itoa((offset + limit)),
	}

	if err := s.cache.SetProto(ctx, cacheKey, resp, 10*time.Minute); err != nil {
		log.Println("Redis lỗi:", err)
	}

	return resp.Items, int64(resp.Total), nil
}

func (s *VuKhiService) GetVersionVuKhiService(
	ctx context.Context,
	req *v1.LayVersionRequest,
) (*v1.LayVersionResponse, error) {

	name := strings.TrimSpace(req.TenVuKhi)
	vk, err := s.grpc_repository.GetVersionByNameVuKhiRepo(ctx, name)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "Không tìm thấy vũ khí")
		}
		return nil, status.Errorf(codes.Internal, "Lỗi DB: %v", err)
	}

	return &v1.LayVersionResponse{
		Version: int32(vk.Version),
	}, nil
}
