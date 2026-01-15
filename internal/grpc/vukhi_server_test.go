package grpcsrv

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"game/ent/vukhi"
	"game/internal/config"
	"game/internal/db"
	v1 "game/v1"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	logOnce sync.Once
	logFile *os.File
)

// initTestLogger: Cấu hình log ra file để debug concurrency dễ hơn
func initTestLogger() {
	logOnce.Do(func() {
		// Mở file log (tạo mới hoặc ghi đè)
		f, err := os.OpenFile(
			"test_vukhi_concurrent.log",
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC, // Truncate để mỗi lần test là file mới
			0666,
		)
		if err != nil {
			panic(fmt.Sprintf("Failed to open log file: %v", err))
		}

		log.SetOutput(f)
		// Quan trọng: Log micro-seconds để soi thứ tự race condition
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
		logFile = f
	})
}

// closeLogger: Đóng file log khi test xong
func closeLogger() {
	if logFile != nil {
		_ = logFile.Sync()
		_ = logFile.Close()
	}
}

var logCounter int64

const maxLogs = 50000 // Giới hạn số dòng log để tránh file quá to

func logAction(gid int64, action, name string, maVuKhi int64, msg string) {
	if atomic.AddInt64(&logCounter, 1) > maxLogs {
		return
	}
	// Format: [GID] ACTION Name (ID): Message
	log.Printf("[%d] %-6s %-15s (id=%d): %s\n", gid, action, name, maVuKhi, msg)
}

func goid() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	var id int64
	fmt.Sscanf(string(buf[:n]), "goroutine %d ", &id)
	return id
}

func setupTestServer(t *testing.T) *VuKhiGRPCServer {
	t.Helper()
	ctx := context.Background()
	cfg := config.DocCauHinh()
	testDB := cfg.TenDBEnt + "_test_concurrent" // DB riêng cho test này

	// Reset DB
	_ = db.DropDBEnt(cfg.DSNChuaDB, testDB) // Ignore err
	require.NoError(t, db.CreateDBEnt(cfg.DSNChuaDB, testDB))

	entClient, err := db.OpenEnt(cfg.DSNChuaDB, testDB)
	require.NoError(t, err)
	require.NoError(t, db.RunMigrationEnt(ctx, entClient))

	t.Cleanup(func() { _ = entClient.Close() })

	// Mock Redis
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { rdb.Close() })

	return &VuKhiGRPCServer{
		ent: entClient,
		rdb: rdb,
		ttl: 1 * time.Second, // TTL ngắn để test cache invalidation
	}
}

type RefDataIDs struct {
	Loai   []int32
	DoHiem []int32
	He     []int32
}

func seedRefData(ctx context.Context, s *VuKhiGRPCServer) RefDataIDs {
	var ids RefDataIDs
	for i := 1; i <= 3; i++ {
		entObj, _ := s.ent.LoaiVuKhi.Create().SetTenLoai(fmt.Sprintf("Loai-%d", i)).Save(ctx)
		ids.Loai = append(ids.Loai, int32(entObj.ID))
	}
	for i := 1; i <= 3; i++ {
		entObj, _ := s.ent.DoHiem.Create().
			SetTenDoHiem(fmt.Sprintf("DoHiem-%d", i)).
			SetSoLuong(10).SetCapBac(i).
			SetSatThuongBonus(1.0).SetTocDoDanhBonus(0.1).
			Save(ctx)
		ids.DoHiem = append(ids.DoHiem, int32(entObj.ID))
	}
	for i := 1; i <= 3; i++ {
		entObj, _ := s.ent.He.Create().SetTenHe(fmt.Sprintf("He-%d", i)).Save(ctx)
		ids.He = append(ids.He, int32(entObj.ID))
	}
	return ids
}

var (
	vkPool   []string
	poolLock sync.Mutex
)

func poolAdd(name string) {
	poolLock.Lock()
	defer poolLock.Unlock()
	vkPool = append(vkPool, name)
}

func poolRemove(name string) {
	poolLock.Lock()
	defer poolLock.Unlock()
	for i, n := range vkPool {
		if n == name {
			vkPool[i] = vkPool[len(vkPool)-1]
			vkPool = vkPool[:len(vkPool)-1]
			return
		}
	}
}

func poolReplace(oldName, newName string) {
	poolLock.Lock()
	defer poolLock.Unlock()
	for i, n := range vkPool {
		if n == oldName {
			vkPool[i] = newName
			return
		}
	}
}

func poolRandom() (string, bool) {
	poolLock.Lock()
	defer poolLock.Unlock()
	if len(vkPool) == 0 {
		return "", false
	}
	return vkPool[rand.Intn(len(vkPool))], true
}

func randomCreateReq(name string, ref RefDataIDs) *v1.TaoVuKhiRequest {
	return &v1.TaoVuKhiRequest{
		TenVuKhi:       name,
		SatThuongCoBan: int32(rand.Intn(200)),
		TocDoDanh:      rand.Float64()*2 + 0.5,
		TamDanh:        int32(rand.Intn(5)),
		MaLoai:         ref.Loai[rand.Intn(len(ref.Loai))],
		MaDoHiem:       ref.DoHiem[rand.Intn(len(ref.DoHiem))],
		MaHe:           ref.He[rand.Intn(len(ref.He))],
		IdempotencyKey: "",
	}
}

func randomUpdateReq(oldName string, ref RefDataIDs, version int32) *v1.CapNhatTheoTenVuKhiRequest {
	newName := oldName
	if rand.Intn(2) == 0 {
		newName = fmt.Sprintf("vk-%04d", rand.Intn(10000))
	}

	return &v1.CapNhatTheoTenVuKhiRequest{
		Name:           oldName,
		TenVuKhi:       newName,
		SatThuongCoBan: int32(rand.Intn(300)),
		TocDoDanh:      rand.Float64()*2 + 0.5,
		TamDanh:        int32(rand.Intn(5)),
		MaLoai:         ref.Loai[rand.Intn(len(ref.Loai))],
		MaDoHiem:       ref.DoHiem[rand.Intn(len(ref.DoHiem))],
		MaHe:           ref.He[rand.Intn(len(ref.He))],
		Version:        version, // ✅ đúng
	}
}

func (s *VuKhiGRPCServer) ServeHTTP() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/vukhi/create", func(w http.ResponseWriter, r *http.Request) {
		var req v1.TaoVuKhiRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		vk, err := s.CreateVuKhi(r.Context(), &req)
		if err != nil {
			code, msg := mapGRPCErrorToHTTP(err)
			http.Error(w, msg, code)
			logAction(goid(), "CREATE", req.TenVuKhi, 0, fmt.Sprintf("Fail: %s", msg))
			return
		}
		json.NewEncoder(w).Encode(vk)
	})

	mux.HandleFunc("/vukhi/update", func(w http.ResponseWriter, r *http.Request) {
		var req v1.CapNhatTheoTenVuKhiRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		vk, err := s.UpdateByNameVuKhi(r.Context(), &req)
		if err != nil {
			code, msg := mapGRPCErrorToHTTP(err)
			http.Error(w, msg, code)
			logAction(goid(), "UPDATE", req.Name, 0, fmt.Sprintf("Fail: %s", msg))
			return
		}
		json.NewEncoder(w).Encode(vk)
	})

	mux.HandleFunc("/vukhi/delete", func(w http.ResponseWriter, r *http.Request) {
		var req v1.XoaTheoTenVuKhiRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := s.DeleteByNameVuKhi(r.Context(), &req)
		if err != nil {
			code, msg := mapGRPCErrorToHTTP(err)
			http.Error(w, msg, code)
			logAction(goid(), "DELETE", req.Name, 0, fmt.Sprintf("Fail: %s", msg))
			return
		}
		json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/vukhi/search", func(w http.ResponseWriter, r *http.Request) {
		var req v1.TimKiemRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := s.Search(r.Context(), &req)
		if err != nil {
			code, msg := mapGRPCErrorToHTTP(err)
			http.Error(w, msg, code)
			logAction(goid(), "SEARCH", req.TenVuKhi, 0, fmt.Sprintf("Fail: %s", msg))
			return
		}

		// Optimized Log
		if len(resp.Items) == 0 {
			logAction(goid(), "SEARCH", req.TenVuKhi, 0, "không tìm thấy")
		} else {
			msg := fmt.Sprintf("tìm thấy %d kết quả (First ID: %d)", len(resp.Items), resp.Items[0].MaVuKhi)
			logAction(goid(), "SEARCH", req.TenVuKhi, int64(resp.Items[0].MaVuKhi), msg)
		}
		json.NewEncoder(w).Encode(resp)
	})

	return mux
}

func Test_HTTP_Concurrent_Safe(t *testing.T) {
	initTestLogger()
	defer closeLogger()

	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()

	// Init Server
	s := setupTestServer(t)
	refIDs := seedRefData(ctx, s)
	ts := httptest.NewServer(s.ServeHTTP())
	defer ts.Close()

	// Reset pool
	vkPool = vkPool[:0]

	// Cấu hình tải
	const (
		workers         = 100 // Số luồng đồng thời
		roundsPerWorker = 50  // Số vòng lặp mỗi luồng
	)

	var wg sync.WaitGroup
	wg.Add(workers)

	client := &http.Client{Timeout: 5 * time.Second}

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < roundsPerWorker; j++ {
				action := rand.Intn(4) // Random hành động

				switch action {
				case 0: // CREATE
					name := fmt.Sprintf("vk-%04d", rand.Intn(10000))
					req := randomCreateReq(name, refIDs)
					b, _ := json.Marshal(req)
					resp, err := client.Post(ts.URL+"/vukhi/create", "application/json", bytes.NewReader(b))

					if err != nil {
						logAction(goid(), "CREATE", name, 0, "Err: "+err.Error())
						continue
					}

					body, _ := io.ReadAll(resp.Body)
					resp.Body.Close()

					if resp.StatusCode >= 400 {
						logAction(goid(), "CREATE", name, 0, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)))
						continue
					}

					var vk v1.VuKhi
					if err := json.Unmarshal(body, &vk); err == nil {
						poolAdd(name)
						logAction(goid(), "CREATE", name, int64(vk.MaVuKhi), "thành công")
					}

				case 1: // UPDATE
					oldName, ok := poolRandom()
					if !ok {
						continue
					}

					verResp, err := s.LayVersion(ctx, &v1.LayVersionRequest{
						TenVuKhi: oldName,
					})
					if err != nil {
						logAction(goid(), "UPDATE", oldName, 0, "Fail get version: "+err.Error())
						continue
					}

					oldVersion := verResp.Version // ⭐ lưu version cũ

					req := randomUpdateReq(oldName, refIDs, oldVersion)

					// Logic đảm bảo tên thay đổi để test
					if req.TenVuKhi == oldName {
						req.TenVuKhi = fmt.Sprintf("%s-%d", oldName, rand.Int63())
					}

					b, _ := json.Marshal(req)
					request, _ := http.NewRequest(http.MethodPut, ts.URL+"/vukhi/update", bytes.NewReader(b))
					request.Header.Set("Content-Type", "application/json")

					resp, err := client.Do(request)
					if err != nil {
						logAction(goid(), "UPDATE", oldName, 0, "Err: "+err.Error())
						continue
					}

					body, _ := io.ReadAll(resp.Body)
					resp.Body.Close()

					if resp.StatusCode >= 400 {
						logAction(goid(), "UPDATE", oldName, 0,
							fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)))
						continue
					}

					var vk v1.VuKhi
					if err := json.Unmarshal(body, &vk); err == nil {
						poolReplace(oldName, req.TenVuKhi)

						// ⭐ LOG CÓ VERSION
						msg := fmt.Sprintf(
							"thành công -> đổi thành: %s | version %d → %d",
							req.TenVuKhi,
							oldVersion,
							vk.Version,
						)

						logAction(goid(), "UPDATE", oldName, int64(vk.MaVuKhi), msg)
					}

				case 2: // DELETE
					name, ok := poolRandom()
					if !ok {
						continue
					}

					req := v1.XoaTheoTenVuKhiRequest{Name: name}
					b, _ := json.Marshal(&req)
					request, _ := http.NewRequest(http.MethodPost, ts.URL+"/vukhi/delete", bytes.NewReader(b))
					request.Header.Set("Content-Type", "application/json")
					resp, err := client.Do(request)

					if err != nil {
						logAction(goid(), "DELETE", name, 0, "Err: "+err.Error())
						continue
					}
					body, _ := io.ReadAll(resp.Body)
					resp.Body.Close()

					if resp.StatusCode >= 400 {
						logAction(goid(), "DELETE", name, 0, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)))
						continue
					}

					var delResp v1.XoaTheoTenVuKhiResponse
					if err := json.Unmarshal(body, &delResp); err == nil {
						poolRemove(name)
						logAction(goid(), "DELETE", name, int64(delResp.MaVuKhi), "thành công")
					}

				case 3: // SEARCH
					name, ok := poolRandom()
					if !ok {
						continue
					}

					req := v1.TimKiemRequest{TenVuKhi: name}
					b, _ := json.Marshal(&req)
					resp, err := client.Post(ts.URL+"/vukhi/search", "application/json", bytes.NewReader(b))

					if err != nil {
						logAction(goid(), "SEARCH", name, 0, "Err: "+err.Error())
						continue
					}
					body, _ := io.ReadAll(resp.Body)
					resp.Body.Close()

					if resp.StatusCode >= 400 {
						logAction(goid(), "SEARCH", name, 0, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)))
						continue
					}
					// Log đã được xử lý trong ServeHTTP
				}
				time.Sleep(2 * time.Millisecond)
			}
		}()
	}

	wg.Wait()
	verifyDBInvariants(t, s)
}

func verifyDBInvariants(t *testing.T, s *VuKhiGRPCServer) {
	verifyUniqueTenVuKhi(t, s)
	verifyNoGhostRecords(t, s)
	verifySearchNoGhost(t, s)
}

func verifyUniqueTenVuKhi(t *testing.T, s *VuKhiGRPCServer) {
	ctx := context.Background()
	names, err := s.ent.VuKhi.Query().Select(vukhi.FieldTenVuKhi).Strings(ctx)
	require.NoError(t, err)
	seen := map[string]bool{}
	for _, n := range names {
		if seen[n] {
			t.Fatalf("DATA INTEGRITY FAIL: Duplicate name found: %s", n)
		}
		seen[n] = true
	}
}

func verifyNoGhostRecords(t *testing.T, s *VuKhiGRPCServer) {
	ctx := context.Background()
	rows, err := s.ent.VuKhi.Query().All(ctx)
	require.NoError(t, err)
	for _, r := range rows {
		exists, _ := s.ent.VuKhi.Query().Where(vukhi.IDEQ(r.ID)).Exist(ctx)
		if !exists {
			t.Fatalf("DATA INTEGRITY FAIL: Ghost record in memory id=%d", r.ID)
		}
	}
}

func verifySearchNoGhost(t *testing.T, s *VuKhiGRPCServer) {
	ctx := context.Background()
	resp, err := s.Search(ctx, &v1.TimKiemRequest{})
	if err != nil {
		return
	}
	for _, it := range resp.Items {
		exists, _ := s.ent.VuKhi.Query().Where(vukhi.IDEQ(int(it.MaVuKhi))).Exist(ctx)
		if !exists {
			t.Fatalf("DATA INTEGRITY FAIL: Search returned deleted ID=%d", it.MaVuKhi)
		}
	}
}

func mapGRPCErrorToHTTP(err error) (int, string) {
	st, ok := status.FromError(err)
	if !ok {
		return 500, "Internal server error"
	}
	switch st.Code() {
	case codes.NotFound:
		return 404, st.Message()
	case codes.InvalidArgument:
		return 400, st.Message()
	case codes.ResourceExhausted:
		return 409, st.Message()
	case codes.Aborted, codes.FailedPrecondition:
		return 409, st.Message()
	case codes.Internal:
		return 500, st.Message()
	default:
		return 500, st.Message()
	}
}
