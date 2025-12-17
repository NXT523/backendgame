1. Khi nào tin nhắn được coi là "Hoàn thành" (Commit)?
    Tin nhắn được coi là hoàn thành khi Consumer báo cho Kafka Broker biết: 
    "Tôi đã xử lý xong tin nhắn ở Offset X, lần sau hãy đưa tôi tin nhắn ở Offset X+1".
    
Trong segmentio/kafka-go, có 2 cách để làm việc này:
Cách 1: Commit Tự Động (ReadMessage) - Ít an toàn hơn

    + Khi bạn gọi reader.ReadMessage(ctx), thư viện sẽ lấy tin nhắn về và tự động commit ngay sau đó (hoặc theo chu kỳ CommitInterval).
    + Rủi ro: Nếu code của bạn bị lỗi (panic/crash) trong quá trình xử lý logic (sau khi đã nhận tin), tin nhắn đó đã bị đánh dấu là "xong" rồi. Bạn sẽ bị mất tin nhắn đó (không được xử lý lại).

Cách 2: Commit Thủ Công (FetchMessage + CommitMessages) - Khuyên dùng

    B1: Bạn gọi reader.FetchMessage(ctx). Tin nhắn được tải về ứng dụng, nhưng Offset trên Kafka chưa thay đổi.
    B2: Bạn chạy logic xử lý (lưu DB, tính toán...).
    B3: Nếu xử lý thành công, bạn gọi reader.CommitMessages(ctx, msg).

    Kết quả: Lúc này Kafka mới ghi nhận là tin nhắn đã hoàn thành.

2. Khi nào tin nhắn bị "Đẩy lại" (Retry)?
Trong Kafka, "đẩy lại" thực chất là việc đọc lại một Offset cũ chưa được Commit. Điều này xảy ra trong các trường hợp sau:
Trường hợp A: Consumer bị Crash / Restart (Chưa Commit)
    Nếu bạn dùng FetchMessage để lấy tin nhắn, đang xử lý thì app bị sập (hoặc mất điện, kill process) trước khi gọi CommitMessages:
        1.Offset trên Kafka vẫn nằm ở vị trí cũ.
        2.Khi App khởi động lại, nó hỏi Kafka: "Offset gần nhất của tôi ở đâu?".
        3.Kafka trả về Offset cũ.
        4.App nhận lại tin nhắn đó và xử lý lại từ đầu.

Trường hợp B: Xử lý lỗi (Logic Retry trong code)
    Nếu logic của bạn bị lỗi (ví dụ: lỗi kết nối DB), bạn không được gọi CommitMessages.

    + Tuy nhiên, kafka-go sẽ không tự động nhảy lại tin nhắn đó ngay lập tức trong vòng lặp hiện tại (vì nó là một dòng chảy stream).
    + Để retry ngay lập tức, bạn phải viết vòng lặp for bên trong code xử lý của mình.

Trường hợp C: Rebalance (Consumer Group thay đổi)
    Nếu Consumer của bạn bị treo quá lâu (quá thời gian Heartbeat hoặc SessionTimeout) mà chưa commit:

    1.Kafka coi Consumer này đã chết.
    2.Kafka chuyển phân vùng (Partition) đó cho một Consumer khác trong nhóm.
    3.Consumer mới sẽ đọc lại từ Offset được commit cuối cùng -> Tin nhắn được xử lý lại bởi consumer khác.