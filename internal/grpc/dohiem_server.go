// grpcsrv/dohiem_server.go
package grpcsrv

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"time"

	"game/ent"
	"game/ent/dohiem"
	"game/pkg/memcached"
	v1 "game/v1"

	"github.com/bradfitz/gomemcache/memcache"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DoHiemGRPCServer implements v1.DoHiemServiceServer
type DoHiemGRPCServer struct {
	v1.UnimplementedDoHiemServiceServer
	ent *ent.Client

	mc  *memcache.Client
	ttl time.Duration
}

func NewDoHiemGRPCServer(entClient *ent.Client) *DoHiemGRPCServer {
	return &DoHiemGRPCServer{ent: entClient}
}

func NewDoHiemGRPCServerMemcached(entClient *ent.Client, mc *memcache.Client, ttl time.Duration) *DoHiemGRPCServer {
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return &DoHiemGRPCServer{ent: entClient, mc: mc, ttl: ttl}
}

func toPBDoHiem(m *ent.DoHiem) *v1.DoHiem {
	if m == nil {
		return nil
	}
	return &v1.DoHiem{
		MaDoHiem:       int32(m.ID),
		TenDoHiem:      m.TenDoHiem,
		SoLuong:        int32(m.SoLuong),
		MauSac:         m.MauSac,
		SatThuongBonus: m.SatThuongBonus,
		TocDoDanhBonus: m.TocDoDanhBonus,
		CapBac:         int32(m.CapBac),
	}
}

const (
	kDoHiemSearchPrefix = "dohiem:search:v1" // cache kết quả tìm kiếm
)

// khóa search theo request: dohiem:search:v1:<sha1(json_request)>
func (s *DoHiemGRPCServer) khoaTimKiem(req *v1.TimKiemDoHiemRequest) string {
	b, _ := (protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}).Marshal(req)
	sum := sha1.Sum(b)
	return kDoHiemSearchPrefix + ":" + hex.EncodeToString(sum[:])
}

//  CRUD & GET LIST

func (s *DoHiemGRPCServer) CreateDoHiem(ctx context.Context, in *v1.TaoDoHiemRequest) (*v1.DoHiem, error) {
	create := s.ent.DoHiem.Create().
		SetTenDoHiem(strings.TrimSpace(in.TenDoHiem)).
		SetSoLuong(int(in.SoLuong)).
		SetMauSac(strings.TrimSpace(in.MauSac)).
		SetSatThuongBonus(in.SatThuongBonus).
		SetTocDoDanhBonus(in.TocDoDanhBonus).
		SetCapBac(int(in.CapBac))

	rec, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}
	return toPBDoHiem(rec), nil
}

func (s *DoHiemGRPCServer) UpdateByNameDoHiem(ctx context.Context, req *v1.CapNhatTheoTenDoHiemRequest) (*v1.DoHiem, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name không hợp lệ")
	}
	cur, err := s.ent.DoHiem.Query().Where(dohiem.TenDoHiemEQ(name)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "Không tìm thấy độ hiếm")
		}
		return nil, status.Errorf(codes.Internal, "Lỗi truy vấn: %v", err)
	}
	up := s.ent.DoHiem.UpdateOneID(cur.ID).
		SetTenDoHiem(req.TenDoHiem).
		SetSoLuong(int(req.SoLuong)).
		SetMauSac(req.MauSac).
		SetSatThuongBonus(req.SatThuongBonus).
		SetTocDoDanhBonus(req.TocDoDanhBonus).
		SetCapBac(int(req.CapBac))
	after, err := up.Save(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cập nhật lỗi: %v", err)
	}
	after, _ = s.ent.DoHiem.Query().Where(dohiem.IDEQ(after.ID)).Only(ctx)
	return toPBDoHiem(after), nil
}

func (s *DoHiemGRPCServer) DeleteByNameDoHiem(ctx context.Context, req *v1.XoaTheoTenDoHiemRequest) (*v1.XoaTheoTenDoHiemResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name không hợp lệ")
	}

	rec, err := s.ent.DoHiem.Query().Where(dohiem.TenDoHiemEQ(name)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "Không tìm thấy độ hiếm")
		}
		return nil, status.Errorf(codes.Internal, "Lỗi truy vấn: %v", err)
	}

	if err := s.ent.DoHiem.DeleteOneID(rec.ID).Exec(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "Xóa lỗi: %v", err)
	}

	return &v1.XoaTheoTenDoHiemResponse{DeletedName: name}, nil
}

func (s *DoHiemGRPCServer) GetAllDoHiem(ctx context.Context, in *v1.LayTatCaDoHiemRequest) (*v1.DanhSachDoHiem, error) {
	qb := s.ent.DoHiem.Query()

	if q := strings.TrimSpace(in.Q); q != "" {
		qb = qb.Where(dohiem.TenDoHiemContainsFold(q))
	}

	total, err := qb.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}

	items, err := qb.Order(ent.Asc(dohiem.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}

	out := &v1.DanhSachDoHiem{
		Total: int32(total),
		Items: make([]*v1.DoHiem, 0, len(items)),
	}
	for _, m := range items {
		out.Items = append(out.Items, toPBDoHiem(m))
	}
	return out, nil
}

func (s *DoHiemGRPCServer) GetAllTenDoHiem(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachStringDoHiem, error) {
	items, err := s.ent.DoHiem.Query().Order(ent.Asc(dohiem.FieldID)).All(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Lỗi truy vấn: %v", err)
	}
	out := &v1.DanhSachStringDoHiem{
		Total: int32(len(items)),
		Items: make([]string, 0, len(items)),
	}
	for _, m := range items {
		out.Items = append(out.Items, m.TenDoHiem)
	}
	return out, nil
}

func (s *DoHiemGRPCServer) GetAllSoLuongDoHiem(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachIntDoHiem, error) {
	items, err := s.ent.DoHiem.Query().Order(ent.Asc(dohiem.FieldID)).All(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Lỗi truy vấn: %v", err)
	}
	out := &v1.DanhSachIntDoHiem{
		Total: int32(len(items)),
		Items: make([]int32, 0, len(items)),
	}
	for _, m := range items {
		out.Items = append(out.Items, int32(m.SoLuong))
	}
	return out, nil
}

func (s *DoHiemGRPCServer) GetAllMauSacDoHiem(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachStringDoHiem, error) {
	items, err := s.ent.DoHiem.Query().Order(ent.Asc(dohiem.FieldID)).All(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Lỗi truy vấn: %v", err)
	}
	out := &v1.DanhSachStringDoHiem{
		Total: int32(len(items)),
		Items: make([]string, 0, len(items)),
	}
	for _, m := range items {
		out.Items = append(out.Items, m.MauSac)
	}
	return out, nil
}

func (s *DoHiemGRPCServer) GetAllSatThuongDoHiem(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachDoubleDoHiem, error) {
	items, err := s.ent.DoHiem.Query().Order(ent.Asc(dohiem.FieldID)).All(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Lỗi truy vấn: %v", err)
	}
	out := &v1.DanhSachDoubleDoHiem{
		Total: int32(len(items)),
		Items: make([]float64, 0, len(items)),
	}
	for _, m := range items {
		out.Items = append(out.Items, m.SatThuongBonus)
	}
	return out, nil
}

func (s *DoHiemGRPCServer) GetAllTocDoDoHiem(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachDoubleDoHiem, error) {
	items, err := s.ent.DoHiem.Query().Order(ent.Asc(dohiem.FieldID)).All(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Lỗi truy vấn: %v", err)
	}
	out := &v1.DanhSachDoubleDoHiem{
		Total: int32(len(items)),
		Items: make([]float64, 0, len(items)),
	}
	for _, m := range items {
		out.Items = append(out.Items, m.TocDoDanhBonus)
	}
	return out, nil
}

//  SEARCH + CACHE (Memcached JSON)

func (s *DoHiemGRPCServer) SearchDoHiem(ctx context.Context, req *v1.TimKiemDoHiemRequest) (*v1.DanhSachDoHiem, error) {
	// ---- Cache theo request hash ----
	cacheKey := s.khoaTimKiem(req)
	if page, ok := memcached.LayTrangCacheJSONMem(s.mc, cacheKey, func() *v1.DanhSachDoHiem {
		return &v1.DanhSachDoHiem{}
	}); ok {
		return page, nil
	}

	// ---- Xây query ----
	qb := s.ent.DoHiem.Query()

	// ID
	if req.MaDoHiem != nil {
		qb = qb.Where(dohiem.IDEQ(int(req.MaDoHiem.Value)))
	}

	// Màu sắc
	if c := strings.TrimSpace(req.MauSac); c != "" {
		qb = qb.Where(dohiem.MauSacEQ(c))
	}

	// Từ khoá tên (name chỉ để search theo tên, KHÔNG còn parse 1-5/alias)
	if q := strings.TrimSpace(req.Q); q != "" {
		qb = qb.Where(dohiem.TenDoHiemContainsFold(q))
	}
	if t := strings.TrimSpace(req.TenDoHiem); t != "" {
		qb = qb.Where(dohiem.TenDoHiemContainsFold(t))
	}

	// Số lượng (eq / gte / lte)
	if req.SoLuong != nil {
		qb = qb.Where(dohiem.SoLuongEQ(int(req.SoLuong.Value)))
	}
	if req.MinSoLuong != nil {
		qb = qb.Where(dohiem.SoLuongGTE(int(req.MinSoLuong.Value)))
	}
	if req.MaxSoLuong != nil {
		qb = qb.Where(dohiem.SoLuongLTE(int(req.MaxSoLuong.Value)))
	}

	// Sát thương bonus (eq / gte / lte)
	if req.SatThuongBonus != nil {
		qb = qb.Where(dohiem.SatThuongBonusEQ(req.SatThuongBonus.Value))
	}
	if req.MinSatThuongBonus != nil {
		qb = qb.Where(dohiem.SatThuongBonusGTE(req.MinSatThuongBonus.Value))
	}
	if req.MaxSatThuongBonus != nil {
		qb = qb.Where(dohiem.SatThuongBonusLTE(req.MaxSatThuongBonus.Value))
	}

	// Tốc độ đánh bonus (eq / gte / lte)
	if req.TocDoDanhBonus != nil {
		qb = qb.Where(dohiem.TocDoDanhBonusEQ(req.TocDoDanhBonus.Value))
	}
	if req.MinTocDoDanhBonus != nil {
		qb = qb.Where(dohiem.TocDoDanhBonusGTE(req.MinTocDoDanhBonus.Value))
	}
	if req.MaxTocDoDanhBonus != nil {
		qb = qb.Where(dohiem.TocDoDanhBonusLTE(req.MaxTocDoDanhBonus.Value))
	}

	// Cấp bậc: CHỈ dùng min_cap_bac / max_cap_bac
	if req.MinCapBac != nil && req.MaxCapBac != nil && req.MinCapBac.Value > req.MaxCapBac.Value {
		// Hoán đổi nếu người dùng nhập ngược
		minV := int(req.MaxCapBac.Value)
		maxV := int(req.MinCapBac.Value)
		qb = qb.Where(dohiem.CapBacGTE(minV)).Where(dohiem.CapBacLTE(maxV))
	} else {
		if req.MinCapBac != nil {
			qb = qb.Where(dohiem.CapBacGTE(int(req.MinCapBac.Value)))
		}
		if req.MaxCapBac != nil {
			qb = qb.Where(dohiem.CapBacLTE(int(req.MaxCapBac.Value)))
		}
	}

	// ---- Thực thi ----
	items, err := qb.Order(ent.Asc(dohiem.FieldID)).All(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Tìm kiếm độ hiếm lỗi: %v", err)
	}

	out := &v1.DanhSachDoHiem{
		Total: int32(len(items)),
		Items: make([]*v1.DoHiem, 0, len(items)),
	}
	for _, m := range items {
		out.Items = append(out.Items, toPBDoHiem(m))
	}

	// ---- Lưu cache ----
	memcached.LuuTrangCacheJSONMem(s.mc, cacheKey, out, s.ttl)
	return out, nil
}
