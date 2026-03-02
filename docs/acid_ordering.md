# ACID & Ordering khi nhiều server gửi request đồng thời

## Kết luận nhanh
- **ACID**: hiện tại hệ thống đảm bảo ở **mức từng transaction của từng item**.
- **Ordering (thứ tự toàn cục)**: hiện tại **không đảm bảo thứ tự tuyệt đối** giữa nhiều server.
- **Exactly-once**: hiện tại **chưa có** (cần idempotency key + dedupe store).

## Vì sao chưa có thứ tự toàn cục?
- Mỗi instance có queue nội bộ `batchQueue`, không phải queue dùng chung toàn cụm.
- Có 20 worker xử lý song song nên thứ tự xử lý thực tế có thể khác thứ tự request đến.

## Đoạn `for i := 0; i < 20; i++ { go s.startBatchWorker() }` có ổn không?
- **Dùng được** nếu mục tiêu là tăng throughput và giảm độ trễ khi tải cao.
- Tuy nhiên con số `20` là **hard-code**, không tối ưu cho mọi môi trường (CPU/DB pool khác nhau).
- Quan trọng: tăng worker **không tạo global ordering**; ngược lại còn làm thứ tự xử lý phân tán hơn.

### Nên đổi như thế nào cho thực tế production
- Đưa số worker thành cấu hình (ví dụ `STREAM_WORKERS`), default theo `min(2*CPU, DB_MAX_OPEN_CONNS)`.
- Theo dõi metric trước khi tăng/giảm:
  - queue depth,
  - p95/p99 latency,
  - tỉ lệ lỗi `Aborted`/timeout,
  - DB saturation (open conns, lock wait).
- Nếu cần **thứ tự theo thực thể**, thay vì tăng worker nội bộ:
  - partition queue theo key thực thể (`id` / `ten_vu_khi`),
  - mỗi partition xử lý tuần tự.

## Cái đang có để đảm bảo đúng dữ liệu
1. Distributed lock theo key cho một số luồng ghi.
2. Optimistic locking theo `version` khi update stream.
3. Transaction DB cho từng update item.
4. Unique index DB để chống trùng dữ liệu.

## Nếu cần vừa ACID tốt vừa có thứ tự khi multi-server
### Mức 1: Đảm bảo đúng dữ liệu (không cần strict order)
- Giữ optimistic locking + retry có backoff + jitter.
- Thêm idempotency key cho lệnh ghi.
- Chuẩn hóa error mapping (`Aborted`, `FailedPrecondition`, `AlreadyExists`) để client retry đúng.

### Mức 2: Đảm bảo thứ tự theo từng thực thể (per-entity ordering)
- Dùng message broker có partition key theo thực thể (ví dụ `ten_vu_khi` hoặc `id`).
- Một key chỉ đi vào 1 partition + 1 consumer group semantics => giữ order theo key.

### Mức 3: Strict global ordering
- Dùng sequencer toàn cục (global monotonic sequence) rồi áp dụng tuần tự.
- Đổi lại: latency cao hơn, throughput thấp hơn, tăng độ phức tạp.

## Gợi ý thực tế cho hệ thống hiện tại
- Chọn **per-entity ordering** thay vì global ordering.
- Kết hợp:
  - Redis/Kafka queue phân tán,
  - partition theo `ten_vu_khi` (hoặc ID),
  - optimistic locking version,
  - idempotency key.

Cách này giữ được tính đúng dữ liệu khi tải cao từ nhiều server và vẫn mở rộng tốt.
