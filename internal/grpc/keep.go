package grpcsrv

// import (
// 	"context"
// 	"crypto/sha1"
// 	"encoding/hex"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"game/ent"
// 	"game/ent/dohiem"
// 	"game/ent/he"
// 	"game/ent/loaivukhi"
// 	"game/ent/vukhi"
// 	"game/internal/events"
// 	"game/pkg/mapping"
// 	redisx "game/pkg/redis"
// 	v1 "game/v1"

// 	entsql "entgo.io/ent/dialect/sql"
// 	"github.com/go-sql-driver/mysql"
// 	"github.com/redis/go-redis/v9"
// 	"github.com/segmentio/kafka-go"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/metadata"
// 	"google.golang.org/grpc/status"
// 	"google.golang.org/protobuf/encoding/protojson"
// 	"google.golang.org/protobuf/proto"
// 	"google.golang.org/protobuf/types/known/emptypb"
// )

// // VuKhiGRPCServer cần có s.rdb *redis.Client
// type VuKhiGRPCServer struct {
// 	v1.UnimplementedVuKhiServiceServer
// 	ent *ent.Client
// 	rdb *redis.Client

// 	ttl                time.Duration // TTL cache cho trang search
// 	kafkaWriterCreated *kafka.Writer
// 	kafkaWriterUpdated *kafka.Writer
// }

// // Mặc định TTL = 10 phút
// func NewVuKhiGRPCServer(entClient *ent.Client, rdb *redis.Client) *VuKhiGRPCServer {
// 	return &VuKhiGRPCServer{ent: entClient, rdb: rdb, ttl: 10 * time.Minute}
// }

// // Khởi tạo kèm TTL tuỳ biến
// func NewVuKhiGRPCServerWithTTL(entClient *ent.Client, rdb *redis.Client, ttl time.Duration) *VuKhiGRPCServer {
// 	if ttl <= 0 {
// 		ttl = 10 * time.Minute
// 	}
// 	return &VuKhiGRPCServer{ent: entClient, rdb: rdb, ttl: ttl}
// }

// func NewVuKhiGRPCServerKafka(entClient *ent.Client, rdb *redis.Client, kwCreated, kwUpdated *kafka.Writer) *VuKhiGRPCServer {
// 	return &VuKhiGRPCServer{
// 		ent:                entClient,
// 		rdb:                rdb,
// 		ttl:                10 * time.Minute,
// 		kafkaWriterCreated: kwCreated,
// 		kafkaWriterUpdated: kwUpdated,
// 	}
// }

// func toPB(item *ent.VuKhi) *v1.VuKhi {
// 	out := &v1.VuKhi{
// 		MaVuKhi:        int32(item.ID),
// 		TenVuKhi:       item.TenVuKhi,
// 		SatThuongCoBan: int32(item.SatThuongCoBan),
// 		TocDoDanh:      item.TocDoDanh,
// 		TamDanh:        int32(item.TamDanh),
// 		MoTa:           item.MoTa,
// 		MaLoai:         int32(item.MaLoai),
// 		MaDoHiem:       int32(item.MaDoHiem),
// 		MaHe:           int32(item.MaHe),
// 	}
// 	if item.Edges.Loai != nil {
// 		out.TenLoai = item.Edges.Loai.TenLoai
// 	}
// 	if item.Edges.He != nil {
// 		out.TenHe = item.Edges.He.TenHe
// 	}
// 	if item.Edges.DoHiem != nil {
// 		out.TenDoHiem = item.Edges.DoHiem.TenDoHiem
// 		out.MauSac = item.Edges.DoHiem.MauSac
// 		out.SoLuong = int32(item.Edges.DoHiem.SoLuong)
// 		out.SatThuongBonus = item.Edges.DoHiem.SatThuongBonus
// 		out.TocDoDanhBonus = item.Edges.DoHiem.TocDoDanhBonus
// 		out.CapBac = int32(item.Edges.DoHiem.CapBac)
// 	}
// 	return out
// }

// func listToPB(rows []*ent.VuKhi) *v1.DanhSachVuKhi {
// 	resp := &v1.DanhSachVuKhi{}
// 	for _, it := range rows {
// 		resp.Items = append(resp.Items, toPB(it))
// 	}
// 	return resp
// }

// const (
// 	kVuKhiSearchPrefix = "vukhi:search:v2"
// )

// // vukhi:search:v1:<sha1(json_request)>
// func (s *VuKhiGRPCServer) khoaTimKiem(req *v1.TimKiemRequest, limit, offset int) string {
// 	cp := proto.Clone(req).(*v1.TimKiemRequest)
// 	cp.Limit = int32(limit)
// 	cp.Cursor = encodeCursor(offset)

// 	b, _ := (protojson.MarshalOptions{
// 		UseProtoNames:   true,
// 		EmitUnpopulated: true,
// 	}).Marshal(cp)

// 	sum := sha1.Sum(b)
// 	return kVuKhiSearchPrefix + ":" + hex.EncodeToString(sum[:])
// }

// // CRUD
// func (s *VuKhiGRPCServer) CreateVuKhi(
// 	ctx context.Context,
// 	req *v1.TaoVuKhiRequest,
// ) (*v1.VuKhi, error) {

// 	// 1️⃣ Validate
// 	name := strings.TrimSpace(req.TenVuKhi)
// 	if name == "" {
// 		return nil, status.Error(codes.InvalidArgument, "ten_vu_khi là bắt buộc")
// 	}

// 	// 0️⃣ Idempotency key
// 	if req.IdempotencyKey != "" {
// 		idKey := "idemp:vukhi:create:" + req.IdempotencyKey

// 		// đã xử lý rồi → trả lại response cũ
// 		cached, err := s.rdb.Get(ctx, idKey).Bytes()
// 		if err == nil {
// 			var out v1.VuKhi
// 			_ = json.Unmarshal(cached, &out)
// 			return &out, nil
// 		}

// 		// đánh dấu đang xử lý
// 		ok, err := s.rdb.SetNX(ctx, idKey, "processing", time.Hour).Result()
// 		if err != nil {
// 			return nil, status.Errorf(codes.Internal, "Idempotency redis lỗi")
// 		}
// 		if !ok {
// 			return nil, status.Error(codes.Aborted, "Request đang được xử lý")
// 		}

// 		// khi thành công → cache response
// 		var result *v1.VuKhi

// 		defer func() {
// 			if result != nil {
// 				b, _ := json.Marshal(result)
// 				_ = s.rdb.Set(context.Background(), idKey, b, time.Hour).Err()
// 			}
// 		}()
// 	}

// 	// 2️⃣ Distributed lock (tránh create trùng khi multi-server)
// 	normalized := strings.ToLower(strings.TrimSpace(name))
// 	normalized = strings.Join(strings.Fields(normalized), " ")

// 	lockKey := "lock:vukhi:create:" + normalized
// 	lockTTL := 2 * time.Second

// 	ok, err := s.rdb.SetNX(ctx, lockKey, "1", lockTTL).Result()
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "Redis lock lỗi: %v", err)
// 	}
// 	if !ok {
// 		return nil, status.Error(codes.ResourceExhausted, "Vũ khí đang được tạo, thử lại sau")
// 	}
// 	defer s.rdb.Del(ctx, lockKey)

// 	// 3️⃣ Transaction
// 	tx, err := s.ent.Tx(ctx)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "Khởi tạo transaction lỗi: %v", err)
// 	}
// 	defer tx.Rollback()

// 	// 5️⃣ Create
// 	v, err := tx.VuKhi.Create().
// 		SetTenVuKhi(name).
// 		SetSatThuongCoBan(int(req.SatThuongCoBan)).
// 		SetTocDoDanh(req.TocDoDanh).
// 		SetTamDanh(int(req.TamDanh)).
// 		SetMoTa(req.MoTa).
// 		SetMaLoai(int(req.MaLoai)).
// 		SetMaDoHiem(int(req.MaDoHiem)).
// 		SetMaHe(int(req.MaHe)).
// 		Save(ctx)
// 	if err != nil {
// 		var mysqlErr *mysql.MySQLError
// 		if errors.As(err, &mysqlErr) {
// 			switch mysqlErr.Number {
// 			case 1062: // UNIQUE
// 				return nil, status.Error(codes.AlreadyExists, "Tên vũ khí đã tồn tại")
// 			case 1452: // FOREIGN KEY
// 				return nil, status.Error(codes.InvalidArgument, "Loại / Độ hiếm / Hệ không tồn tại")
// 			default:
// 				return nil, status.Errorf(codes.Internal, "Lỗi DB (%d)", mysqlErr.Number)
// 			}
// 		}

// 		return nil, status.Errorf(codes.Internal, "Tạo vũ khí lỗi: %v", err)
// 	}

// 	if err := tx.Commit(); err != nil {
// 		return nil, status.Errorf(codes.Internal, "Commit transaction lỗi: %v", err)
// 	}

// 	// 8️⃣ Eager load sau commit (read-only)
// 	v, _ = s.ent.VuKhi.Query().
// 		Where(vukhi.IDEQ(v.ID)).
// 		WithLoai().
// 		WithDoHiem().
// 		WithHe().
// 		Only(ctx)

// 	if s.kafkaWriterCreated != nil {
// 		ev := events.SuKienVuKhiCreated{
// 			Type:           "VU_KHI_CREATED",
// 			ID:             v.ID,
// 			TenVuKhi:       v.TenVuKhi,
// 			SatThuongCoBan: v.SatThuongCoBan,
// 			TocDoDanh:      v.TocDoDanh,
// 			TamDanh:        v.TamDanh,
// 			MoTa:           v.MoTa,
// 			MaLoai:         v.MaLoai,
// 			MaHe:           v.MaHe,
// 			MaDoHiem:       v.MaDoHiem,
// 			CreatedAt:      time.Now().Unix(),
// 			CreatedBy:      "system",
// 		}
// 		data, _ := json.Marshal(ev)

// 		if err := s.rdb.RPush(ctx, "outbox:vukhicreate", data).Err(); err != nil {
// 			return nil, status.Errorf(codes.Internal, "Push outbox lỗi: %v", err)
// 		}
// 	}

// 	result := toPB(v)
// 	_ = s.xoaCacheSearch(ctx)
// 	return result, nil
// }

// func (s *VuKhiGRPCServer) UpdateByNameVuKhi(ctx context.Context, req *v1.CapNhatTheoTenVuKhiRequest) (*v1.VuKhi, error) {
// 	// 1️⃣ Validate input
// 	name := strings.TrimSpace(req.Name)
// 	if name == "" {
// 		return nil, status.Error(codes.InvalidArgument, "name không hợp lệ")
// 	}

// 	normalized := strings.ToLower(strings.TrimSpace(name))
// 	normalized = strings.Join(strings.Fields(normalized), " ")
// 	// 2️⃣ Idempotency key
// 	idemKey := getIdempotencyKey(ctx)
// 	if idemKey != "" {
// 		idemProcessingKey := "idemp:vukhi:update:processing:" + normalized + ":" + idemKey

// 		ok, err := s.rdb.SetNX(ctx, idemProcessingKey, "1", 30*time.Second).Result()
// 		if err != nil {
// 			return nil, status.Errorf(codes.Internal, "Idempotency redis lỗi")
// 		}
// 		if !ok {
// 			return nil, status.Error(codes.Aborted, "Request đang được xử lý")
// 		}

// 		defer s.rdb.Del(ctx, idemProcessingKey)
// 	}

// 	// 2️⃣ Thiết lập distributed lock (Redis)

// 	lockKey := "lock:vukhi:update:" + normalized
// 	lockTTL := 5 * time.Second // lock thời gian ngắn, đủ cho transaction

// 	// Thử lấy lock
// 	ok, err := s.rdb.SetNX(ctx, lockKey, "1", lockTTL).Result()
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "Redis lock lỗi: %v", err)
// 	}
// 	if !ok {
// 		return nil, status.Error(codes.ResourceExhausted, "Vũ khí đang được cập nhật, thử lại sau")
// 	}
// 	// đảm bảo release lock khi xong
// 	defer s.rdb.Del(ctx, lockKey)

// 	// 3️⃣ Bắt đầu transaction trên DB (Ent)
// 	tx, err := s.ent.Tx(ctx)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "Khởi tạo transaction lỗi: %v", err)
// 	}

// 	defer func() {
// 		if p := recover(); p != nil {
// 			tx.Rollback()
// 			panic(p)
// 		}
// 	}()

// 	cur, err := tx.VuKhi.Query().Where(vukhi.TenVuKhiEQ(name)).Only(ctx)
// 	if err != nil {
// 		tx.Rollback()
// 		if ent.IsNotFound(err) {
// 			return nil, status.Error(codes.NotFound, "Không tìm thấy vũ khí")
// 		}
// 		return nil, status.Errorf(codes.Internal, "Truy vấn lỗi: %v", err)
// 	}

// 	after, err := tx.VuKhi.UpdateOneID(cur.ID).
// 		SetTenVuKhi(req.TenVuKhi).
// 		SetSatThuongCoBan(int(req.SatThuongCoBan)).
// 		SetTocDoDanh(req.TocDoDanh).
// 		SetTamDanh(int(req.TamDanh)).
// 		SetMoTa(req.MoTa).
// 		SetMaLoai(int(req.MaLoai)).
// 		SetMaDoHiem(int(req.MaDoHiem)).
// 		SetMaHe(int(req.MaHe)).
// 		Save(ctx)

// 	if err != nil {
// 		if ent.IsConstraintError(err) {
// 			return nil, status.Error(codes.AlreadyExists, "Tên vũ khí đã tồn tại")
// 		}
// 		tx.Rollback()
// 		return nil, status.Errorf(codes.Internal, "Cập nhật lỗi: %v", err)
// 	}

// 	if err := tx.Commit(); err != nil {
// 		return nil, status.Errorf(codes.Internal, "Commit transaction lỗi: %v", err)
// 	}

// 	after, _ = s.ent.VuKhi.Query().Where(vukhi.IDEQ(after.ID)).WithLoai().WithDoHiem().WithHe().Only(ctx)
// 	resp := toPB(after)

// 	// 5️⃣ Cache idempotency result
// 	if idemKey != "" {
// 		b, _ := protojson.Marshal(resp)
// 		cacheKey := "idemp:vukhi:update:" + normalized + ":" + idemKey
// 		_ = s.rdb.Set(ctx, cacheKey, b, 2*time.Minute).Err()
// 	}
// 	_ = s.xoaCacheSearch(ctx)

// 	return resp, nil
// }

// func (s *VuKhiGRPCServer) DeleteByNameVuKhi(
// 	ctx context.Context,
// 	req *v1.XoaTheoTenVuKhiRequest,
// ) (*v1.XoaTheoTenVuKhiResponse, error) {

// 	// 1️⃣ Validate
// 	name := strings.TrimSpace(req.Name)
// 	if name == "" {
// 		return nil, status.Error(codes.InvalidArgument, "name không hợp lệ")
// 	}

// 	// 2️⃣ Distributed lock (multi-server safe)
// 	normalized := strings.ToLower(strings.TrimSpace(name))
// 	normalized = strings.Join(strings.Fields(normalized), " ")

// 	lockKey := "lock:vukhi:delete:" + normalized
// 	lockTTL := 5 * time.Second

// 	ok, err := s.rdb.SetNX(ctx, lockKey, "1", lockTTL).Result()
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "Redis lock lỗi: %v", err)
// 	}
// 	if !ok {
// 		return nil, status.Error(codes.ResourceExhausted, "Vũ khí đang được xử lý, thử lại sau")
// 	}
// 	defer s.rdb.Del(ctx, lockKey)

// 	// 3️⃣ Transaction
// 	tx, err := s.ent.Tx(ctx)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "Khởi tạo transaction lỗi: %v", err)
// 	}
// 	defer tx.Rollback()

// 	// 4️⃣ Query trong transaction
// 	cur, err := tx.VuKhi.
// 		Query().
// 		Where(vukhi.TenVuKhiEQ(name)).
// 		Only(ctx)
// 	if err != nil {
// 		if ent.IsNotFound(err) {
// 			return nil, status.Error(codes.NotFound, "Không tìm thấy vũ khí")
// 		}
// 		return nil, status.Errorf(codes.Internal, "Lỗi truy vấn: %v", err)
// 	}

// 	// 5️⃣ Delete
// 	if err := tx.VuKhi.DeleteOneID(cur.ID).Exec(ctx); err != nil {
// 		return nil, status.Errorf(codes.Internal, "Xóa lỗi: %v", err)
// 	}

// 	// 6️⃣ Commit
// 	if err := tx.Commit(); err != nil {
// 		return nil, status.Errorf(codes.Internal, "Commit lỗi: %v", err)
// 	}
// 	_ = s.xoaCacheSearch(ctx)

// 	return &v1.XoaTheoTenVuKhiResponse{
// 		DeletedName: name,
// 	}, nil
// }

// //  GET LIST

// func (s *VuKhiGRPCServer) GetAllVuKhi(ctx context.Context, req *v1.LayTatCaRequest) (*v1.DanhSachVuKhi, error) {
// 	qb := s.ent.VuKhi.Query().WithLoai().WithDoHiem().WithHe()
// 	if strings.TrimSpace(req.Q) != "" {
// 		qb = qb.Where(vukhi.TenVuKhiContainsFold(strings.TrimSpace(req.Q)))
// 	}
// 	items, err := qb.All(ctx)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "Lấy danh sách lỗi: %v", err)
// 	}
// 	page := listToPB(items)
// 	page.Total = int32(len(items))
// 	page.ItemCount = int32(len(items))
// 	page.NextCursor = ""
// 	page.HasMore = false

// 	return page, nil
// }

// func (s *VuKhiGRPCServer) GetAllTenVuKhi(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachStringVuKhi, error) {
// 	names, err := s.ent.VuKhi.Query().Select(vukhi.FieldTenVuKhi).Strings(ctx)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "Lấy tên lỗi: %v", err)
// 	}
// 	return &v1.DanhSachStringVuKhi{Total: int32(len(names)), Items: names}, nil
// }

// func (s *VuKhiGRPCServer) GetAllSatThuongCoBan(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachIntVuKhi, error) {
// 	vals, err := s.ent.VuKhi.Query().Select(vukhi.FieldSatThuongCoBan).Ints(ctx)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "Lấy sát thương cơ bản lỗi: %v", err)
// 	}
// 	out := make([]int32, 0, len(vals))
// 	for _, v := range vals {
// 		out = append(out, int32(v))
// 	}
// 	return &v1.DanhSachIntVuKhi{Total: int32(len(out)), Items: out}, nil
// }

// func (s *VuKhiGRPCServer) GetAllTocDoDanh(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachDoubleVuKhi, error) {
// 	vals, err := s.ent.VuKhi.Query().Select(vukhi.FieldTocDoDanh).Float64s(ctx)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "Lấy tốc độ đánh lỗi: %v", err)
// 	}
// 	return &v1.DanhSachDoubleVuKhi{Total: int32(len(vals)), Items: vals}, nil
// }

// func (s *VuKhiGRPCServer) GetAllTamDanh(ctx context.Context, _ *emptypb.Empty) (*v1.DanhSachIntVuKhi, error) {
// 	vals, err := s.ent.VuKhi.Query().Select(vukhi.FieldTamDanh).Ints(ctx)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "Lấy tầm đánh lỗi: %v", err)
// 	}
// 	out := make([]int32, 0, len(vals))
// 	for _, v := range vals {
// 		out = append(out, int32(v))
// 	}
// 	return &v1.DanhSachIntVuKhi{Total: int32(len(out)), Items: out}, nil
// }

// //  SEARCH (cursor-based + cache JSON)

// func (s *VuKhiGRPCServer) Search(ctx context.Context, req *v1.TimKiemRequest) (*v1.DanhSachVuKhi, error) {
// 	// 1) Chuẩn hoá limit/cursor
// 	limit := normalizeLimit(req.Limit)
// 	offset, err := parseCursor(req.Cursor)
// 	if err != nil {
// 		return nil, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	// 2) Cache key (giống DoHiem: dùng hàm khoaTimKiem + prefix)
// 	cacheKey := s.khoaTimKiem(req, limit, offset)

// 	// 3) Đọc cache (Redis JSON)
// 	if page, ok := redisx.LayTrangCacheJSON(ctx, s.rdb, cacheKey, func() *v1.DanhSachVuKhi {
// 		return &v1.DanhSachVuKhi{}
// 	}); ok {
// 		// fallback cho cache cũ chưa có item_count
// 		if page.ItemCount == 0 && len(page.Items) > 0 {
// 			page.ItemCount = int32(len(page.Items))
// 		}
// 		return page, nil
// 	}

// 	// 4) Xây filters chung
// 	applyFilters := func(q *ent.VuKhiQuery) *ent.VuKhiQuery {
// 		// ---- Vũ khí / tên chứa ----
// 		if req.MaVuKhi != nil {
// 			q = q.Where(vukhi.IDEQ(int(req.MaVuKhi.Value)))
// 		}
// 		if t := strings.TrimSpace(req.TenVuKhi); t != "" {
// 			q = q.Where(vukhi.TenVuKhiContainsFold(t))
// 		}
// 		if req.MaLoai != nil {
// 			q = q.Where(vukhi.MaLoaiEQ(int(req.MaLoai.Value)))
// 		}
// 		if t := strings.TrimSpace(req.TenLoai); t != "" {
// 			q = q.Where(vukhi.HasLoaiWith(loaivukhi.TenLoaiContainsFold(t)))
// 		}
// 		if req.MaHe != nil {
// 			q = q.Where(vukhi.MaHeEQ(int(req.MaHe.Value)))
// 		}
// 		if t := strings.TrimSpace(req.TenHe); t != "" {
// 			q = q.Where(vukhi.HasHeWith(he.TenHeContainsFold(t)))
// 		}
// 		if req.MaDoHiem != nil {
// 			q = q.Where(vukhi.MaDoHiemEQ(int(req.MaDoHiem.Value)))
// 		}
// 		if t := strings.TrimSpace(req.TenDoHiem); t != "" {
// 			q = q.Where(vukhi.HasDoHiemWith(dohiem.TenDoHiemContainsFold(t)))
// 		}
// 		if t := strings.TrimSpace(req.MauSac); t != "" {
// 			q = q.Where(vukhi.HasDoHiemWith(dohiem.MauSacEQ(t)))
// 		}

// 		// ---- RANGE: màu sắc (min_mau_sac / max_mau_sac) ----
// 		{
// 			var (
// 				haveMinColor, haveMaxColor bool
// 				minColor, maxColor         string
// 			)

// 			if c := strings.TrimSpace(req.MinMauSac); c != "" {
// 				haveMinColor = true
// 				minColor = c
// 			}
// 			if c := strings.TrimSpace(req.MaxMauSac); c != "" {
// 				haveMaxColor = true
// 				maxColor = c
// 			}

// 			if haveMinColor && haveMaxColor && strings.Compare(minColor, maxColor) > 0 {
// 				minColor, maxColor = maxColor, minColor
// 			}

// 			if haveMinColor && haveMaxColor {
// 				q = q.Where(
// 					vukhi.HasDoHiemWith(
// 						dohiem.MauSacGTE(minColor),
// 						dohiem.MauSacLTE(maxColor),
// 					),
// 				)
// 			} else if haveMinColor {
// 				q = q.Where(vukhi.HasDoHiemWith(dohiem.MauSacGTE(minColor)))
// 			} else if haveMaxColor {
// 				q = q.Where(vukhi.HasDoHiemWith(dohiem.MauSacLTE(maxColor)))
// 			}
// 		}

// 		// ---- RANGE: cap_bac (giao min/max số và min/max theo tên) ----
// 		var (
// 			haveMinCap, haveMaxCap bool
// 			minCap, maxCap         int
// 		)
// 		if req.MinCapBac != nil {
// 			haveMinCap = true
// 			minCap = int(req.MinCapBac.Value)
// 		}
// 		if req.MaxCapBac != nil {
// 			haveMaxCap = true
// 			maxCap = int(req.MaxCapBac.Value)
// 		}
// 		if sMin := strings.TrimSpace(req.MinTenDoHiem); sMin != "" {
// 			if capV, ok := s.resolveCapFromTenDoHiem(ctx, sMin); ok {
// 				if !haveMinCap || capV > minCap {
// 					haveMinCap, minCap = true, capV
// 				}
// 			}
// 		}
// 		if sMax := strings.TrimSpace(req.MaxTenDoHiem); sMax != "" {
// 			if capV, ok := s.resolveCapFromTenDoHiem(ctx, sMax); ok {
// 				if !haveMaxCap || capV < maxCap {
// 					haveMaxCap, maxCap = true, capV
// 				}
// 			}
// 		}
// 		if haveMinCap && haveMaxCap && minCap > maxCap {
// 			minCap, maxCap = maxCap, minCap
// 		}
// 		if haveMinCap && haveMaxCap {
// 			q = q.Where(vukhi.HasDoHiemWith(dohiem.CapBacGTE(minCap), dohiem.CapBacLTE(maxCap)))
// 		} else if haveMinCap {
// 			q = q.Where(vukhi.HasDoHiemWith(dohiem.CapBacGTE(minCap)))
// 		} else if haveMaxCap {
// 			q = q.Where(vukhi.HasDoHiemWith(dohiem.CapBacLTE(maxCap)))
// 		}

// 		// ---- RANGE: so_luong ----
// 		if req.MinSoLuong != nil && req.MaxSoLuong != nil && req.MinSoLuong.Value > req.MaxSoLuong.Value {
// 			minV := int(req.MaxSoLuong.Value)
// 			maxV := int(req.MinSoLuong.Value)
// 			q = q.Where(vukhi.HasDoHiemWith(dohiem.SoLuongGTE(minV), dohiem.SoLuongLTE(maxV)))
// 		} else {
// 			if req.MinSoLuong != nil {
// 				q = q.Where(vukhi.HasDoHiemWith(dohiem.SoLuongGTE(int(req.MinSoLuong.Value))))
// 			}
// 			if req.MaxSoLuong != nil {
// 				q = q.Where(vukhi.HasDoHiemWith(dohiem.SoLuongLTE(int(req.MaxSoLuong.Value))))
// 			}
// 		}

// 		// ---- RANGE: sat_thuong_bonus ----
// 		if req.MinSatThuongBonus != nil && req.MaxSatThuongBonus != nil && req.MinSatThuongBonus.Value > req.MaxSatThuongBonus.Value {
// 			minV := req.MaxSatThuongBonus.Value
// 			maxV := req.MinSatThuongBonus.Value
// 			q = q.Where(vukhi.HasDoHiemWith(dohiem.SatThuongBonusGTE(minV), dohiem.SatThuongBonusLTE(maxV)))
// 		} else {
// 			if req.MinSatThuongBonus != nil {
// 				q = q.Where(vukhi.HasDoHiemWith(dohiem.SatThuongBonusGTE(req.MinSatThuongBonus.Value)))
// 			}
// 			if req.MaxSatThuongBonus != nil {
// 				q = q.Where(vukhi.HasDoHiemWith(dohiem.SatThuongBonusLTE(req.MaxSatThuongBonus.Value)))
// 			}
// 		}

// 		// ---- RANGE: toc_do_danh_bonus ----
// 		if req.MinTocDoDanhBonus != nil && req.MaxTocDoDanhBonus != nil && req.MinTocDoDanhBonus.Value > req.MaxTocDoDanhBonus.Value {
// 			minV := req.MaxTocDoDanhBonus.Value
// 			maxV := req.MinTocDoDanhBonus.Value
// 			q = q.Where(vukhi.HasDoHiemWith(dohiem.TocDoDanhBonusGTE(minV), dohiem.TocDoDanhBonusLTE(maxV)))
// 		} else {
// 			if req.MinTocDoDanhBonus != nil {
// 				q = q.Where(vukhi.HasDoHiemWith(dohiem.TocDoDanhBonusGTE(req.MinTocDoDanhBonus.Value)))
// 			}
// 			if req.MaxTocDoDanhBonus != nil {
// 				q = q.Where(vukhi.HasDoHiemWith(dohiem.TocDoDanhBonusLTE(req.MaxTocDoDanhBonus.Value)))
// 			}
// 		}

// 		// ---- Vũ khí: so khớp chính xác (EQ) ----
// 		if req.SatThuongCoBan != nil {
// 			q = q.Where(vukhi.SatThuongCoBanEQ(int(req.SatThuongCoBan.Value)))
// 		}
// 		if req.TocDoDanh != nil {
// 			q = q.Where(vukhi.TocDoDanhEQ(req.TocDoDanh.Value))
// 		}
// 		if req.TamDanh != nil {
// 			q = q.Where(vukhi.TamDanhEQ(int(req.TamDanh.Value)))
// 		}
// 		if req.SoLuong != nil { // EQ so_luong theo proto field 15
// 			q = q.Where(vukhi.HasDoHiemWith(dohiem.SoLuongEQ(int(req.SoLuong.Value))))
// 		}

// 		// ---- RANGE: sat_thuong_co_ban (min_sat_thuong_co_ban / max_sat_thuong_co_ban) ----
// 		if req.MinSatThuongCoBan != nil && req.MaxSatThuongCoBan != nil &&
// 			req.MinSatThuongCoBan.Value > req.MaxSatThuongCoBan.Value {
// 			minV := int(req.MaxSatThuongCoBan.Value)
// 			maxV := int(req.MinSatThuongCoBan.Value)
// 			q = q.Where(
// 				vukhi.SatThuongCoBanGTE(minV),
// 				vukhi.SatThuongCoBanLTE(maxV),
// 			)
// 		} else {
// 			if req.MinSatThuongCoBan != nil {
// 				q = q.Where(vukhi.SatThuongCoBanGTE(int(req.MinSatThuongCoBan.Value)))
// 			}
// 			if req.MaxSatThuongCoBan != nil {
// 				q = q.Where(vukhi.SatThuongCoBanLTE(int(req.MaxSatThuongCoBan.Value)))
// 			}
// 		}

// 		// ---- RANGE: toc_do_danh (min_toc_do_danh / max_toc_do_danh) ---- // NEW
// 		if req.MinTocDoDanh != nil && req.MaxTocDoDanh != nil &&
// 			req.MinTocDoDanh.Value > req.MaxTocDoDanh.Value {
// 			minV := req.MaxTocDoDanh.Value
// 			maxV := req.MinTocDoDanh.Value
// 			q = q.Where(
// 				vukhi.TocDoDanhGTE(minV),
// 				vukhi.TocDoDanhLTE(maxV),
// 			)
// 		} else {
// 			if req.MinTocDoDanh != nil {
// 				q = q.Where(vukhi.TocDoDanhGTE(req.MinTocDoDanh.Value))
// 			}
// 			if req.MaxTocDoDanh != nil {
// 				q = q.Where(vukhi.TocDoDanhLTE(req.MaxTocDoDanh.Value))
// 			}
// 		}

// 		// ---- RANGE: tam_danh (min_tam_danh / max_tam_danh) ---- // NEW
// 		if req.MinTamDanh != nil && req.MaxTamDanh != nil &&
// 			req.MinTamDanh.Value > req.MaxTamDanh.Value {
// 			minV := int(req.MaxTamDanh.Value)
// 			maxV := int(req.MinTamDanh.Value)
// 			q = q.Where(
// 				vukhi.TamDanhGTE(minV),
// 				vukhi.TamDanhLTE(maxV),
// 			)
// 		} else {
// 			if req.MinTamDanh != nil {
// 				q = q.Where(vukhi.TamDanhGTE(int(req.MinTamDanh.Value)))
// 			}
// 			if req.MaxTamDanh != nil {
// 				q = q.Where(vukhi.TamDanhLTE(int(req.MaxTamDanh.Value)))
// 			}
// 		}

// 		return q
// 	}

// 	// 5) Đếm tổng
// 	total, err := applyFilters(s.ent.VuKhi.Query()).Count(ctx)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "Đếm tổng lỗi: %v", err)
// 	}

// 	hadFilter :=
// 		req.MaVuKhi != nil ||
// 			strings.TrimSpace(req.TenVuKhi) != "" ||
// 			req.MaLoai != nil ||
// 			strings.TrimSpace(req.TenLoai) != "" ||
// 			req.MaHe != nil ||
// 			strings.TrimSpace(req.TenHe) != "" ||
// 			req.MaDoHiem != nil ||
// 			strings.TrimSpace(req.TenDoHiem) != "" ||
// 			strings.TrimSpace(req.MauSac) != "" ||
// 			req.MinCapBac != nil || req.MaxCapBac != nil ||
// 			strings.TrimSpace(req.MinTenDoHiem) != "" || strings.TrimSpace(req.MaxTenDoHiem) != "" ||
// 			req.MinSoLuong != nil || req.MaxSoLuong != nil ||
// 			req.MinSatThuongBonus != nil || req.MaxSatThuongBonus != nil ||
// 			req.MinTocDoDanhBonus != nil || req.MaxTocDoDanhBonus != nil ||
// 			req.SatThuongCoBan != nil || req.TocDoDanh != nil || req.TamDanh != nil || req.SoLuong != nil ||
// 			req.MinSatThuongCoBan != nil || req.MaxSatThuongCoBan != nil ||
// 			req.MinTocDoDanh != nil || req.MaxTocDoDanh != nil ||
// 			req.MinTamDanh != nil || req.MaxTamDanh != nil

// 	if hadFilter && total == 0 {
// 		return nil, status.Error(codes.NotFound, "Không tìm thấy vũ khí theo bộ lọc")
// 	}
// 	if !hasAnyFilter(req) {
// 		if strings.TrimSpace(req.Cursor) == "" && req.Limit == 0 {
// 			return nil, status.Error(
// 				codes.InvalidArgument,
// 				"Phải nhập ít nhất một điều kiện tìm kiếm hoặc limit",
// 			)
// 		}
// 	}

// 	qb := applyFilters(s.ent.VuKhi.Query().WithLoai().WithDoHiem().WithHe())
// 	// ----- ORDER: chỉ cho 1 field & 1 chiều -----
// 	ascRaw := strings.TrimSpace(req.ArrangeAsc)
// 	descRaw := strings.TrimSpace(req.ArrangeDesc)

// 	if ascRaw != "" && descRaw != "" {
// 		return nil, status.Error(codes.InvalidArgument, "Không được truyền đồng thời arrange_asc và arrange_desc")
// 	}
// 	fieldAsc, errAsc := parseSingleOrderField(ascRaw)
// 	fieldDesc, errDesc := parseSingleOrderField(descRaw)
// 	if errAsc != nil {
// 		return nil, status.Errorf(codes.InvalidArgument, "arrange_asc không hợp lệ: %v", errAsc)
// 	}
// 	if errDesc != nil {
// 		return nil, status.Errorf(codes.InvalidArgument, "arrange_desc không hợp lệ: %v", errDesc)
// 	}

// 	orderField := ""
// 	orderAsc := true
// 	if fieldAsc != "" {
// 		orderField, orderAsc = fieldAsc, true
// 	} else if fieldDesc != "" {
// 		orderField, orderAsc = fieldDesc, false
// 	}

// 	hasCustomOrder := orderField != ""
// 	if hasCustomOrder && !validOrderField(orderField) {
// 		return nil, status.Errorf(codes.InvalidArgument, "Trường sắp xếp không hợp lệ: %s", orderField)
// 	}

// 	if hasCustomOrder {
// 		qb = qb.Order(func(sel *entsql.Selector) {
// 			if col, ok := orderableVK[orderField]; ok {
// 				if orderAsc {
// 					sel.OrderBy(sel.C(col))
// 				} else {
// 					sel.OrderBy(entsql.Desc(sel.C(col)))
// 				}
// 			} else if col, ok := orderableDH[orderField]; ok {
// 				doh := entsql.Table(dohiem.Table)
// 				sel.Join(doh).On(sel.C(vukhi.FieldMaDoHiem), doh.C(dohiem.FieldID))
// 				if orderAsc {
// 					sel.OrderBy(doh.C(col))
// 				} else {
// 					sel.OrderBy(entsql.Desc(doh.C(col)))
// 				}
// 			}
// 			// Tie-breaker ổn định
// 			sel.OrderBy(sel.C(vukhi.FieldID))
// 		})
// 	} else {
// 		// Mặc định theo ID ASC
// 		qb = qb.Order(vukhi.ByID())
// 	}

// 	qb = qb.Offset(offset)

// 	// Lấy limit+1 để biết còn trang sau không
// 	rows, err := qb.Limit(limit + 1).All(ctx)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "Tìm kiếm lỗi: %v", err)
// 	}

// 	// 7) Tính hasMore, cắt dư, tạo nextCursor
// 	hasMore := false
// 	if len(rows) > limit {
// 		hasMore = true
// 		rows = rows[:limit]
// 	}

// 	nextCursor := ""
// 	if hasMore {
// 		nextCursor = encodeCursor(offset + limit) // ví dụ: offset=0, limit=5 -> nextCursor="5"
// 	}

// 	page := listToPB(rows)

// 	page.Total = int32(total)
// 	page.ItemCount = int32(len(rows))
// 	page.NextCursor = nextCursor
// 	page.HasMore = hasMore

// 	// 8) Ghi cache (Redis JSON) với TTL của server
// 	redisx.LuuTrangCacheJSON(ctx, s.rdb, cacheKey, page, s.ttl)

// 	return page, nil
// }

// // Helpers: limit/cursor
// func hasAnyFilter(req *v1.TimKiemRequest) bool {
// 	return req.MaVuKhi != nil ||
// 		strings.TrimSpace(req.TenVuKhi) != "" ||
// 		req.MaLoai != nil ||
// 		strings.TrimSpace(req.TenLoai) != "" ||
// 		req.MaHe != nil ||
// 		strings.TrimSpace(req.TenHe) != "" ||
// 		req.MaDoHiem != nil ||
// 		strings.TrimSpace(req.TenDoHiem) != "" ||
// 		strings.TrimSpace(req.MauSac) != "" ||
// 		req.MinCapBac != nil || req.MaxCapBac != nil ||
// 		strings.TrimSpace(req.MinTenDoHiem) != "" ||
// 		strings.TrimSpace(req.MaxTenDoHiem) != "" ||
// 		req.MinSoLuong != nil || req.MaxSoLuong != nil ||
// 		req.MinSatThuongBonus != nil || req.MaxSatThuongBonus != nil ||
// 		req.MinTocDoDanhBonus != nil || req.MaxTocDoDanhBonus != nil ||
// 		req.SatThuongCoBan != nil ||
// 		req.TocDoDanh != nil ||
// 		req.TamDanh != nil ||
// 		req.SoLuong != nil ||
// 		req.MinSatThuongCoBan != nil || req.MaxSatThuongCoBan != nil ||
// 		req.MinTocDoDanh != nil || req.MaxTocDoDanh != nil ||
// 		req.MinTamDanh != nil || req.MaxTamDanh != nil
// }

// func normalizeLimit(l int32) int {
// 	n := int(l)
// 	if n <= 0 || n > 1000 { // cap cứng
// 		n = 30
// 	}
// 	return n
// }

// func parseCursor(cur string) (int, error) {
// 	cur = strings.TrimSpace(cur)

// 	if cur == "" {
// 		return 0, nil
// 	}

// 	if strings.Contains(cur, "{{") {
// 		return 0, nil
// 	}

// 	v, err := strconv.Atoi(cur)
// 	if err != nil {
// 		return 0, fmt.Errorf("cursor không phải số")
// 	}

// 	if v < 0 {
// 		return 0, fmt.Errorf("cursor không được âm")
// 	}

// 	return v, nil
// }
// func encodeCursor(offset int) string {
// 	return strconv.Itoa(offset)
// }

// func (s *VuKhiGRPCServer) resolveCapFromTenDoHiem(ctx context.Context, name string) (int, bool) {
// 	n := strings.TrimSpace(name)
// 	if n == "" {
// 		return 0, false
// 	}
// 	if v, ok := mapping.MapStringToCap(n); ok {
// 		return v, true
// 	}
// 	rec, err := s.ent.DoHiem.Query().
// 		Where(dohiem.TenDoHiemEqualFold(n)).
// 		Select(dohiem.FieldCapBac).
// 		Only(ctx)
// 	if err == nil && rec != nil {
// 		return rec.CapBac, true
// 	}
// 	return 0, false
// }

// // Trường sắp xếp thuộc bảng vukhi
// var orderableVK = map[string]string{
// 	"ma_vu_khi":         vukhi.FieldID,
// 	"ten_vu_khi":        vukhi.FieldTenVuKhi,
// 	"sat_thuong_co_ban": vukhi.FieldSatThuongCoBan,
// 	"toc_do_danh":       vukhi.FieldTocDoDanh,
// 	"tam_danh":          vukhi.FieldTamDanh,
// 	"ma_loai":           vukhi.FieldMaLoai,
// 	"ma_he":             vukhi.FieldMaHe,
// 	"ma_do_hiem":        vukhi.FieldMaDoHiem,
// }

// // Trường sắp xếp thuộc bảng dohiem (qua edge)
// var orderableDH = map[string]string{
// 	"ten_do_hiem":       dohiem.FieldTenDoHiem,
// 	"mau_sac":           dohiem.FieldMauSac,
// 	"so_luong":          dohiem.FieldSoLuong,
// 	"sat_thuong_bonus":  dohiem.FieldSatThuongBonus,
// 	"toc_do_danh_bonus": dohiem.FieldTocDoDanhBonus,
// 	"cap_bac":           dohiem.FieldCapBac,
// }

// func parseSingleOrderField(s string) (string, error) {
// 	f := strings.ToLower(strings.TrimSpace(s))
// 	if f == "" {
// 		return "", nil
// 	}
// 	if strings.Contains(f, ",") {
// 		return "", fmt.Errorf("chỉ cho phép 1 trường sắp xếp, không được có dấu phẩy")
// 	}
// 	return f, nil
// }

// func validOrderField(f string) bool {
// 	if _, ok := orderableVK[f]; ok {
// 		return true
// 	}
// 	if _, ok := orderableDH[f]; ok {
// 		return true
// 	}
// 	return false
// }

// func getIdempotencyKey(ctx context.Context) string {
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		return ""
// 	}
// 	if v := md.Get("idempotency-key"); len(v) > 0 {
// 		return v[0]
// 	}
// 	return ""
// }

// func (s *VuKhiGRPCServer) xoaCacheSearch(ctx context.Context) error {
// 	pattern := kVuKhiSearchPrefix + ":*"

// 	var cursor uint64
// 	for {
// 		keys, nextCursor, err := s.rdb.Scan(ctx, cursor, pattern, 100).Result()
// 		if err != nil {
// 			return err
// 		}
// 		if len(keys) > 0 {
// 			if err := s.rdb.Del(ctx, keys...).Err(); err != nil {
// 				return err
// 			}
// 		}
// 		cursor = nextCursor
// 		if cursor == 0 {
// 			break
// 		}
// 	}
// 	return nil
// }
