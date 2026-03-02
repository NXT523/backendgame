package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	v1 "game/v1/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TestData struct {
	TenVuKhi           string  `json:"ten_vu_khi"`
	SatThuongCoBan int32   `json:"sat_thuong_co_ban"`
	Version        int32   `json:"version"`
	TocDoDanh      float64 `json:"toc_do_danh"`
	TamDanh        int32   `json:"tam_danh"`
	MoTa           string  `json:"mo_ta"`
	MaLoaiVuKhi    int32   `json:"ma_loai_vu_khi"`
	MaDoHiem       int32   `json:"ma_do_hiem"`
	MaHe           int32   `json:"ma_he"`
}

const (
	Address      = "127.0.0.1:9000" // Cổng Proxy
	Concurrency  = 100              // Số luồng song song
	TestDuration = 30 * time.Second // Thời gian chạy test
)

// --- THỐNG KÊ CHI TIẾT ---
var (
	// Mode 0: Unary
	cntSuccessUnary int32
	cntFailUnary    int32

	// Mode 1: Bidi Stream
	cntSuccessBidi int32
	cntFailBidi    int32

	// Mode 2: Client Stream
	cntSuccessClient int32
	cntFailClient    int32

	// Mode 3: Server Stream
	cntSuccessServer int32
	cntFailServer    int32

	// Latency chung
	latencies []time.Duration
	latMu     sync.Mutex

	// Biến đếm tổng request đã gửi
	totalSent int64
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// 1. Setup Data & Connection
	fmt.Println("--> Đang đọc file data_stress.json...")
	file, err := os.ReadFile("data_stress.json")
	if err != nil {
		log.Fatalf("Lỗi đọc file: %v", err)
	}

	var dataList []TestData
	if err := json.Unmarshal(file, &dataList); err != nil {
		log.Fatalf("Lỗi parse JSON: %v", err)
	}

	if len(dataList) == 0 {
		log.Fatal("File data rỗng!")
	}

	fmt.Printf("--> Đã load %d items.\n", len(dataList))
	fmt.Printf("--> Bắt đầu Stress Test trong %v với %d Concurrency...\n", TestDuration, Concurrency)

	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Không thể kết nối Server: %v", err)
	}
	defer conn.Close()
	client := v1.NewVuKhiServiceClient(conn)

	// 2. Chạy Test Time-based
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, Concurrency)

	timer := time.NewTimer(TestDuration)
	startTotal := time.Now()

Loop:
	for {
		select {
		case <-timer.C:
			fmt.Println("\n🛑 Hết thời gian! Đang đợi các request còn lại hoàn tất...")
			break Loop

		case semaphore <- struct{}{}:
			wg.Add(1)
			currentIdx := atomic.AddInt64(&totalSent, 1)
			item := dataList[currentIdx%int64(len(dataList))]
			mode := int(currentIdx % 4)

			go func(d TestData, m int) {
				defer wg.Done()
				defer func() { <-semaphore }()

				switch m {
				case 0: // Unary Update
					runUnaryWorker(client, d)
				case 1: // Bidirectional Stream
					runBidiStreamWorker(client, d)
				case 2: // Client Side Stream
					runClientStreamWorker(client, d)
				case 3: // Server Side Stream
					runServerStreamWorker(client, d)
				}
			}(item, mode)
		}
	}

	wg.Wait()
	totalDuration := time.Since(startTotal)
	printStats(totalDuration, int(totalSent))
}

// 1. WORKER: UNARY UPDATE

func runUnaryWorker(client v1.VuKhiServiceClient, d TestData) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	start := time.Now()
	req := convertToReq(d)

	for retry := 0; retry < 5; retry++ {
		ver, err := client.LayVersion(ctx, &v1.LayVersionRequest{TenVuKhi: d.TenVuKhi})
		if err == nil {
			req.Version = ver.Version
		}

		_, err = client.UpdateByNameVuKhi(ctx, req)
		if err == nil {
			atomic.AddInt32(&cntSuccessUnary, 1)
			recordLatency(time.Since(start))
			return
		}
		time.Sleep(time.Duration(retry*2) * time.Millisecond)
	}

	// Hết retry mới đếm lỗi
	atomic.AddInt32(&cntFailUnary, 1)
	trackError(context.DeadlineExceeded) // Hoặc biến err nếu bạn muốn chính xác
}


// 2. WORKER: BIDIRECTIONAL STREAM

func runBidiStreamWorker(client v1.VuKhiServiceClient, d TestData) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	start := time.Now()

	stream, err := client.StreamBidirectionalUpdateVuKhi(ctx)
	if err != nil {
		atomic.AddInt32(&cntFailBidi, 1)
		trackError(err)
		return
	}
	defer stream.CloseSend()

	// [ĐÃ XÓA ĐOẠN CODE SAI Ở ĐÂY]

	req := convertToReq(d)

	// Lấy version mới nhất
	ver, err := client.LayVersion(ctx, &v1.LayVersionRequest{TenVuKhi: d.TenVuKhi})
	if err == nil {
		req.Version = ver.Version
	}

	for retry := 0; retry < 1; retry++ {
		if err := stream.Send(req); err != nil {
			break
		}
		resp, err := stream.Recv()
		if err != nil {
			break
		}

		if resp.Success {
			atomic.AddInt32(&cntSuccessBidi, 1)
			recordLatency(time.Since(start))
			return
		}
		// Nếu Server trả về Success=false (do sai version), ta coi như thất bại logic
		// trackError(fmt.Errorf("Bidi Logic Error: %s", resp.ErrorMessage))
	}

	// Nếu chạy xuống đây tức là thất bại
	atomic.AddInt32(&cntFailBidi, 1)
	trackError(fmt.Errorf("Bidi Thất bại or Version Không khớp"))
}

// 3. WORKER: CLIENT SIDE STREAM

func runClientStreamWorker(client v1.VuKhiServiceClient, d TestData) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	start := time.Now()

	stream, err := client.StreamClientSideUpdateVuKhi(ctx)
	if err != nil {
		atomic.AddInt32(&cntFailClient, 1)
		return
	}

	req := convertToReq(d)
	ver, err := client.LayVersion(ctx, &v1.LayVersionRequest{TenVuKhi: d.TenVuKhi})
	if err == nil {
		req.Version = ver.Version
	}

	if err := stream.Send(req); err != nil {
		atomic.AddInt32(&cntFailClient, 1) // <--- SỬA
		return
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		atomic.AddInt32(&cntFailClient, 1)
		trackError(err) // [THÊM] Lỗi mạng/gRPC
		return
	}

	// Check kết quả bên trong
	if resp.Success > 0 {
		atomic.AddInt32(&cntSuccessClient, 1)
		recordLatency(time.Since(start))
	} else {
		atomic.AddInt32(&cntFailClient, 1)
	}

	// [THÊM] Quét lỗi chi tiết trong Batch
	if len(resp.Details) > 0 {
		for _, det := range resp.Details {
			if !det.Success {
				// Gom nhóm lỗi từ message trả về
				trackError(fmt.Errorf("Batch Error: %s", det.Message))
			}
		}
	}
}

// 4. WORKER: SERVER SIDE STREAM (READ)

func runServerStreamWorker(client v1.VuKhiServiceClient, d TestData) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	start := time.Now()

	req := &v1.LayTatCaRequest{Q: d.TenVuKhi}

	stream, err := client.StreamServerSideGetVuKhi(ctx, req)
	if err != nil {
		atomic.AddInt32(&cntFailServer, 1)
		trackError(err)
		return
	}

	var totalItems int32 = 0

	// CHỈ DÙNG 1 VÒNG LẶP DUY NHẤT
	for {
		resp, err := stream.Recv()

		// 1. Nếu đọc hết Stream (Thành công)
		if err == io.EOF {
			atomic.AddInt32(&cntSuccessServer, 1)
			recordLatency(time.Since(start))
			// fmt.Printf("Đã đọc xong stream %s, tổng items: %d\n", d.Name, totalItems) // Uncomment để debug
			return
		}

		// 2. Nếu gặp lỗi
		if err != nil {
			atomic.AddInt32(&cntFailServer, 1)
			trackError(err)
			return
		}

		// 3. Nếu đọc được data -> Cộng dồn
		totalItems += int32(len(resp.Items))
	}
}

// HELPERS & BÁO CÁO

func convertToReq(d TestData) *v1.CapNhatTheoTenVuKhiRequest {
	return &v1.CapNhatTheoTenVuKhiRequest{
		TenVuKhiNew:    d.TenVuKhi,
		TenVuKhi:       d.TenVuKhi,
		SatThuongCoBan: d.SatThuongCoBan,
		TocDoDanh:      d.TocDoDanh,
		TamDanh:        d.TamDanh,
		MoTa:           d.MoTa,
		MaLoaiVuKhi:    d.MaLoaiVuKhi,
		MaDoHiem:       d.MaDoHiem,
		MaHe:           d.MaHe,
		Version:        d.Version,
	}
}

func recordLatency(d time.Duration) {
	latMu.Lock()
	latencies = append(latencies, d)
	latMu.Unlock()
}

// --- THỐNG KÊ CHI TIẾT ---
var (
	// ... các biến đếm cũ giữ nguyên ...

	// [THÊM MỚI] Map lưu thông điệp lỗi
	errorCounts = make(map[string]int)
	errMu       sync.Mutex
)

// [THÊM MỚI] Hàm helper để đếm lỗi
func trackError(err error) {
	if err == nil {
		return
	}
	msg := err.Error()
	// Rút gọn thông báo lỗi gRPC để dễ gom nhóm
	// Ví dụ: gom các lỗi "Version không khớp" lại, bỏ qua chi tiết râu ria
	if len(msg) > 100 {
		msg = msg[:100] + "..."
	}

	errMu.Lock()
	errorCounts[msg]++
	errMu.Unlock()
}

func printStats(totalDuration time.Duration, totalReq int) {
	rps := float64(totalReq) / totalDuration.Seconds()

	// Tính tổng
	totalSuccess := cntSuccessUnary + cntSuccessBidi + cntSuccessClient + cntSuccessServer
	totalFail := cntFailUnary + cntFailBidi + cntFailClient + cntFailServer

	fmt.Println("\n================ BÁO CÁO CHI TIẾT CHAOS TEST ================")
	fmt.Printf("Tổng thời gian:  %v\n", totalDuration)
	fmt.Printf("Tổng Requests:   %d\n", totalReq)
	fmt.Printf("Tốc độ (RPS):    %.2f req/s\n", rps)
	fmt.Printf("Thành công:      %d (%.2f%%)\n", totalSuccess, float64(totalSuccess)/float64(totalReq)*100)
	fmt.Printf("Thất bại:        %d (%.2f%%)\n", totalFail, float64(totalFail)/float64(totalReq)*100)
	fmt.Println("-------------------------------------------------------------")
	fmt.Printf("%-20s | %-10s | %-10s | %s\n", "LOẠI TEST", "THÀNH CÔNG", "THẤT BẠI", "TỶ LỆ OK")
	fmt.Println("-------------------------------------------------------------")

	printRow("Unary Update", cntSuccessUnary, cntFailUnary)
	printRow("Bidi Stream", cntSuccessBidi, cntFailBidi)
	printRow("Client Stream", cntSuccessClient, cntFailClient)
	printRow("Server Stream (Read)", cntSuccessServer, cntFailServer)

	fmt.Println("-------------------------------------------------------------")

	if len(latencies) > 0 {
		sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
		var totalLat time.Duration
		for _, l := range latencies {
			totalLat += l
		}
		avg := totalLat / time.Duration(len(latencies))

		fmt.Println("\n--- Độ trễ chung (Latency) ---")
		fmt.Printf("Min: %v | Avg: %v | Max: %v\n", latencies[0], avg, latencies[len(latencies)-1])
		fmt.Printf("P95: %v | P99: %v\n", getPercentile(latencies, 95), getPercentile(latencies, 99))
	}
	fmt.Println("=============================================================")
	fmt.Println("-------------------------------------------------------------")
	// [THÊM ĐOẠN NÀY]
	fmt.Println("\n================ PHÂN TÍCH NGUYÊN NHÂN LỖI ================")
	if len(errorCounts) == 0 {
		fmt.Println("✅ Không có lỗi nào được ghi nhận!")
	} else {
		// Sắp xếp lỗi để in cho đẹp
		type errStat struct {
			Msg   string
			Count int
		}
		var sorted []errStat
		for msg, count := range errorCounts {
			sorted = append(sorted, errStat{Msg: msg, Count: count})
		}
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Count > sorted[j].Count
		})

		fmt.Printf("%-10s | %s\n", "SỐ LƯỢNG", "NỘI DUNG LỖI")
		fmt.Println("-----------|-------------------------------------------------")
		for _, s := range sorted {
			fmt.Printf("%-10d | %s\n", s.Count, s.Msg)
		}
	}
	fmt.Println("=============================================================")
}

func printRow(name string, success, fail int32) {
	total := success + fail
	rate := 0.0
	if total > 0 {
		rate = float64(success) / float64(total) * 100
	}
	fmt.Printf("%-20s | %-10d | %-10d | %.2f%%\n", name, success, fail, rate)
}

func getPercentile(sorted []time.Duration, p int) time.Duration {
	idx := int(math.Ceil(float64(len(sorted))*float64(p)/100)) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}
