Các thành phần chính trong Kafka
1. Broker
    + Là server Kafka thực sự, nhận message từ producer và cung cấp cho consumer.
    + Một cluster Kafka có thể có nhiều broker (ví dụ: localhost:9092, localhost:9093…).
    + ví dụ:
        Brokers: []string{"localhost:9093"}
    + Broker là nơi producer gửi message và consumer lấy message.

2. Topic
    + Là kênh logic để gửi/nhận dữ liệu.
    + Producer gửi message vào topic, consumer đọc từ topic.
    + Một topic có thể chia thành nhiều partition để tăng throughput.
    + Ví dụ: topic "vu_khi_created".
        Topic: "vu_khi_created"

3. Partition
    + Một topic có thể có nhiều partition.
    + Partition đảm bảo thứ tự message trong partition đó.
    + Mỗi message trong partition có offset (số thứ tự, auto tăng).
    + Ví dụ: topic có 3 partition → messages được phân chia theo key hoặc round-robin.

4. Replication Factor
    + Số bản sao của partition trên các broker khác nhau.
    + Mục đích: tăng tính sẵn sàng (high availability).
    + Ví dụ: nếu replication factor = 3, mỗi partition có 3 bản sao trên 3 broker khác nhau.

5. Producer
    + Là client gửi message đến Kafka.
    + Producer có thể:
        * Gửi message với key hoặc không có key.
        * Key sẽ quyết định partition nào nhận message.
    + Ví dụ:
        writer := kafka.NewWriter(kafka.WriterConfig{
            Brokers: []string{"localhost:9093"},
            Topic:   "vu_khi_created",
        })

6. Consumer
    + Là client nhận message từ Kafka.
    + Có 2 cách đọc:
        * Individual Consumer: đọc trực tiếp từng partition.
        * Consumer Group: nhiều consumer cùng group, Kafka tự động phân phối message giữa các consumer trong group.
    + Ví dụ Go với (segmentio/kafka-go):
        reader := kafka.NewReader(kafka.ReaderConfig{
            Brokers: []string{"localhost:9093"},
            Topic:   "vu_khi_created",
            GroupID: "game_group",
        })
    + GroupID quyết định cách phân phối message: mỗi message chỉ được 1 consumer trong group xử lý.

7. Message / Record
    Message gồm
    + Key (tùy chọn)
    + Value (dữ liệu thực)
    + Timestamp (thời gian gửi)
    + Partition & Offset (vị trí trong partition)
    + Ví dụ:
    msg := kafka.Message{
        Key:       []byte("id123"),                  // Tùy chọn, dùng để phân vùng
        Value:     []byte("Hello Kafka"),           // Dữ liệu chính
        Time:      time.Now(),                       // Timestamp, nếu không cung cấp Kafka sẽ tự sinh
        Partition: 0,                                // Thường Kafka tự chọn nếu không ghi
        Offset:    0,                                // Kafka gán khi message được ghi
    }

8. Offset
    + Là số thứ tự message trong partition.
    + Consumer dùng offset để xác định message nào đã đọc.
    + Kafka cho phép:
        * Auto commit offset (tự động lưu)
        * Manual commit offset (tùy chỉnh khi đã xử lý xong)

9. Retention Policy
    + Kafka lưu message trong khoảng thời gian hoặc theo dung lượng.
    + Ví dụ:
        retention.ms = 604800000 → giữ 7 ngày.
        retention.bytes = 1GB → giữ tối đa 1GB dữ liệu.

10. MinBytes / MaxBytes / MaxWait (Go kafka-go config)
    + MinBytes: Số byte tối thiểu đọc từ broker trước khi return message.
    + MaxBytes: Số byte tối đa đọc 1 lần.
    + MaxWait: Thời gian chờ tối đa nếu chưa đủ MinBytes.
    + Ví dụ:
        reader := kafka.NewReader(kafka.ReaderConfig{
            MinBytes: 1,          // ít nhất 1 byte
            MaxBytes: 10e6,       // tối đa 10MB
            MaxWait:  3 * time.Second,
        })

11. Key phân phối message
    + Nếu message có key, Kafka dùng hash(key) % num_partitions để xác định partition.
    + Nếu không có key, Kafka dùng round-robin.

12. Consumer Group Scaling
    + Một topic có N partition → tối đa N consumer đồng thời trong 1 group.
    + Nếu consumer > N, một số consumer sẽ không nhận message nào.