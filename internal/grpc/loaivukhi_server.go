package grpcsrv

import (
	"context"
	"strings"

	"game/ent"
	"game/ent/loaivukhi"
	v1 "game/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// LoaiVuKhiGRPCServer implements v1.LoaiVuKhiServiceServer
type LoaiVuKhiGRPCServer struct {
	v1.UnimplementedLoaiVuKhiServiceServer
	ent *ent.Client
}

func NewLoaiVuKhiGRPCServer(entClient *ent.Client) *LoaiVuKhiGRPCServer {
	return &LoaiVuKhiGRPCServer{ent: entClient}
}

func toPBLoaiVuKhi(m *ent.LoaiVuKhi) *v1.LoaiVuKhi {
	if m == nil {
		return nil
	}
	return &v1.LoaiVuKhi{
		MaLoai:  int32(m.ID),
		TenLoai: m.TenLoai,
		MoTa:    m.MoTa,
	}
}

// POST /loaivukhi/ent/post
func (s *LoaiVuKhiGRPCServer) CreateLoaiVuKhi(ctx context.Context, in *v1.TaoLoaiVuKhiRequest) (*v1.LoaiVuKhi, error) {
	create := s.ent.LoaiVuKhi.Create().
		SetTenLoai(strings.TrimSpace(in.TenLoai)).
		SetMoTa(strings.TrimSpace(in.MoTa))

	rec, err := create.Save(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Tạo loại vũ khí lỗi: %v", err)
	}
	return toPBLoaiVuKhi(rec), nil
}

// PUT /loaivukhi/ent/put-by-name/{name}
func (s *LoaiVuKhiGRPCServer) UpdateByNameLoaiVuKhi(ctx context.Context, req *v1.CapNhatTheoTenLoaiVuKhiRequest) (*v1.LoaiVuKhi, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name không hợp lệ")
	}

	cur, err := s.ent.LoaiVuKhi.Query().Where(loaivukhi.TenLoaiEQ(name)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "Không tìm thấy loại vũ khí")
		}
		return nil, status.Errorf(codes.Internal, "Lỗi truy vấn: %v", err)
	}

	after, err := s.ent.LoaiVuKhi.UpdateOneID(cur.ID).
		SetTenLoai(strings.TrimSpace(req.TenLoai)).
		SetMoTa(strings.TrimSpace(req.MoTa)).
		Save(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cập nhật lỗi: %v", err)
	}

	after, _ = s.ent.LoaiVuKhi.Query().Where(loaivukhi.IDEQ(after.ID)).Only(ctx)
	return toPBLoaiVuKhi(after), nil
}

// DELETE /loaivukhi/ent/delete-by-name/{name}
func (s *LoaiVuKhiGRPCServer) DeleteByNameLoaiVuKhi(ctx context.Context, req *v1.XoaTheoTenLoaiVuKhiRequest) (*v1.XoaTheoTenLoaiVuKhiResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name không hợp lệ")
	}

	rec, err := s.ent.LoaiVuKhi.Query().Where(loaivukhi.TenLoaiEQ(name)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "Không tìm thấy loại vũ khí")
		}
		return nil, status.Errorf(codes.Internal, "Lỗi truy vấn: %v", err)
	}

	if err := s.ent.LoaiVuKhi.DeleteOneID(rec.ID).Exec(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "Xóa lỗi: %v", err)
	}

	return &v1.XoaTheoTenLoaiVuKhiResponse{DeletedName: name}, nil
}

// GET /loaivukhi/ent/getall
func (s *LoaiVuKhiGRPCServer) GetAllLoaiVuKhi(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachLoaiVuKhi, error) {
	items, err := s.ent.LoaiVuKhi.Query().Order(ent.Asc(loaivukhi.FieldID)).All(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Lỗi truy vấn: %v", err)
	}

	out := &v1.DanhSachLoaiVuKhi{
		Total: int32(len(items)),
		Items: make([]*v1.LoaiVuKhi, 0, len(items)),
	}
	for _, m := range items {
		out.Items = append(out.Items, toPBLoaiVuKhi(m))
	}
	return out, nil
}

// GET /loaivukhi/ent/ten
func (s *LoaiVuKhiGRPCServer) GetAllTenLoaiVuKhi(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachStringLoaiVuKhi, error) {
	items, err := s.ent.LoaiVuKhi.Query().Order(ent.Asc(loaivukhi.FieldID)).All(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Lỗi truy vấn: %v", err)
	}

	out := &v1.DanhSachStringLoaiVuKhi{
		Total: int32(len(items)),
		Items: make([]string, 0, len(items)),
	}
	for _, m := range items {
		out.Items = append(out.Items, m.TenLoai)
	}
	return out, nil
}

// GET /loaivukhi/ent/mota
func (s *LoaiVuKhiGRPCServer) GetAllMoTaLoaiVuKhi(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachStringLoaiVuKhi, error) {
	items, err := s.ent.LoaiVuKhi.Query().Order(ent.Asc(loaivukhi.FieldID)).All(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Lỗi truy vấn: %v", err)
	}

	out := &v1.DanhSachStringLoaiVuKhi{
		Total: int32(len(items)),
		Items: make([]string, 0, len(items)),
	}
	for _, m := range items {
		out.Items = append(out.Items, m.MoTa)
	}
	return out, nil
}

// POST /loaivukhi/ent/search
func (s *LoaiVuKhiGRPCServer) SearchLoaiVuKhi(ctx context.Context, req *v1.TimKiemLoaiVuKhiRequest) (*v1.DanhSachLoaiVuKhi, error) {
	qb := s.ent.LoaiVuKhi.Query()

	// Lọc theo id
	if req.MaLoai != nil {
		qb = qb.Where(loaivukhi.IDEQ(int(req.MaLoai.Value)))
	}

	// Lọc theo tên & mô tả (không phân biệt hoa thường)
	if t := strings.TrimSpace(req.TenLoai); t != "" {
		qb = qb.Where(loaivukhi.TenLoaiContainsFold(t))
	}
	if d := strings.TrimSpace(req.MoTa); d != "" {
		qb = qb.Where(loaivukhi.MoTaContainsFold(d))
	}

	// q: search nhanh trên tên
	if q := strings.TrimSpace(req.Q); q != "" {
		qb = qb.Where(loaivukhi.TenLoaiContainsFold(q))
	}

	items, err := qb.Order(ent.Asc(loaivukhi.FieldID)).All(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Tìm kiếm loại vũ khí lỗi: %v", err)
	}

	out := &v1.DanhSachLoaiVuKhi{
		Total: int32(len(items)),
		Items: make([]*v1.LoaiVuKhi, 0, len(items)),
	}
	for _, m := range items {
		out.Items = append(out.Items, toPBLoaiVuKhi(m))
	}
	return out, nil
}
