package grpcsrv

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"time"

	"game/ent"
	"game/ent/he"
	"game/pkg/memcached"
	v1 "game/v1"

	"github.com/bradfitz/gomemcache/memcache"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
)

// HeGRPCServer implements v1.HeServiceServer
type HeGRPCServer struct {
	v1.UnimplementedHeServiceServer
	ent *ent.Client

	mc  *memcache.Client
	ttl time.Duration
}

func NewHeGRPCServer(entClient *ent.Client) *HeGRPCServer {
	return &HeGRPCServer{ent: entClient}
}

func NewHeGRPCServerMemcached(entClient *ent.Client, mc *memcache.Client, ttl time.Duration) *HeGRPCServer {
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return &HeGRPCServer{ent: entClient, mc: mc, ttl: ttl}
}

func toPBHe(m *ent.He) *v1.He {
	if m == nil {
		return nil
	}
	return &v1.He{
		MaHe:  int32(m.ID),
		TenHe: m.TenHe,
		MoTa:  m.MoTa,
	}
}

//  Khóa cache (đặt cạnh "class")

const (
	kHeSearchPrefix = "he:search:v1"
)

// khóa search theo request: he:search:v1:<sha1(json_request)>
func (s *HeGRPCServer) khoaTimKiem(req *v1.TimKiemHeRequest) string {
	b, _ := (protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}).Marshal(req)
	sum := sha1.Sum(b)
	return kHeSearchPrefix + ":" + hex.EncodeToString(sum[:])
}

//  CRUD

// POST /he/ent/post
func (s *HeGRPCServer) CreateHe(ctx context.Context, in *v1.TaoHeRequest) (*v1.He, error) {
	rec, err := s.ent.He.Create().
		SetTenHe(strings.TrimSpace(in.TenHe)).
		SetMoTa(strings.TrimSpace(in.MoTa)).
		Save(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Tạo hệ lỗi: %v", err)
	}
	return toPBHe(rec), nil
}

// PUT /he/ent/put-by-name/{name}
func (s *HeGRPCServer) UpdateByNameHe(ctx context.Context, req *v1.CapNhatTheoTenHeRequest) (*v1.He, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name không hợp lệ")
	}

	cur, err := s.ent.He.Query().Where(he.TenHeEQ(name)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "Không tìm thấy hệ")
		}
		return nil, status.Errorf(codes.Internal, "Lỗi truy vấn: %v", err)
	}

	after, err := s.ent.He.UpdateOneID(cur.ID).
		SetTenHe(strings.TrimSpace(req.TenHe)).
		SetMoTa(strings.TrimSpace(req.MoTa)).
		Save(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cập nhật lỗi: %v", err)
	}

	after, _ = s.ent.He.Query().Where(he.IDEQ(after.ID)).Only(ctx)
	return toPBHe(after), nil
}

// DELETE /he/ent/delete-by-name/{name}
func (s *HeGRPCServer) DeleteByNameHe(ctx context.Context, req *v1.XoaTheoTenHeRequest) (*v1.XoaTheoTenHeResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name không hợp lệ")
	}

	rec, err := s.ent.He.Query().Where(he.TenHeEQ(name)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "Không tìm thấy hệ")
		}
		return nil, status.Errorf(codes.Internal, "Lỗi truy vấn: %v", err)
	}

	if err := s.ent.He.DeleteOneID(rec.ID).Exec(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "Xóa lỗi: %v", err)
	}

	return &v1.XoaTheoTenHeResponse{DeletedName: name}, nil
}

//  GET LIST

// GET /he/ent/getall
func (s *HeGRPCServer) GetAllHe(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachHe, error) {
	items, err := s.ent.He.Query().Order(ent.Asc(he.FieldID)).All(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Lỗi truy vấn: %v", err)
	}

	out := &v1.DanhSachHe{
		Total: int32(len(items)),
		Items: make([]*v1.He, 0, len(items)),
	}
	for _, m := range items {
		out.Items = append(out.Items, toPBHe(m))
	}
	return out, nil
}

// GET /he/ent/ten
func (s *HeGRPCServer) GetAllTenHe(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachStringHe, error) {
	items, err := s.ent.He.Query().Order(ent.Asc(he.FieldID)).All(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Lỗi truy vấn: %v", err)
	}

	out := &v1.DanhSachStringHe{
		Total: int32(len(items)),
		Items: make([]string, 0, len(items)),
	}
	for _, m := range items {
		out.Items = append(out.Items, m.TenHe)
	}
	return out, nil
}

func (s *HeGRPCServer) GetAllMoTaHe(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachStringHe, error) {
	items, err := s.ent.He.Query().Order(ent.Asc(he.FieldID)).All(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Lỗi truy vấn: %v", err)
	}

	out := &v1.DanhSachStringHe{
		Total: int32(len(items)),
		Items: make([]string, 0, len(items)),
	}
	for _, m := range items {
		out.Items = append(out.Items, m.MoTa)
	}
	return out, nil
}

//  SEARCH + CACHE (Memcached JSON)

// POST /he/ent/search
func (s *HeGRPCServer) SearchHe(ctx context.Context, req *v1.TimKiemHeRequest) (*v1.DanhSachHe, error) {
	//  Cache theo request hash
	cacheKey := s.khoaTimKiem(req)
	if page, ok := memcached.LayTrangCacheJSONMem(s.mc, cacheKey, func() *v1.DanhSachHe {
		return &v1.DanhSachHe{}
	}); ok {
		return page, nil
	}

	//  Xây query
	qb := s.ent.He.Query()

	// Lọc theo id (ma_he)
	if req.MaHe != nil {
		qb = qb.Where(he.IDEQ(int(req.MaHe.Value)))
	}

	// Lọc theo tên, mô tả (contains-fold để không phân biệt hoa thường)
	if t := strings.TrimSpace(req.TenHe); t != "" {
		qb = qb.Where(he.TenHeContainsFold(t))
	}
	if d := strings.TrimSpace(req.MoTa); d != "" {
		qb = qb.Where(he.MoTaContainsFold(d))
	}

	// q: tìm nhanh trên tên
	if q := strings.TrimSpace(req.Q); q != "" {
		qb = qb.Where(he.TenHeContainsFold(q))
	}

	//  Thực thi
	items, err := qb.Order(ent.Asc(he.FieldID)).All(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Tìm kiếm hệ lỗi: %v", err)
	}

	out := &v1.DanhSachHe{
		Total: int32(len(items)),
		Items: make([]*v1.He, 0, len(items)),
	}
	for _, m := range items {
		out.Items = append(out.Items, toPBHe(m))
	}

	//  Lưu cache
	memcached.LuuTrangCacheJSONMem(s.mc, cacheKey, out, s.ttl)
	return out, nil
}
