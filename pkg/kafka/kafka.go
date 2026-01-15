package kafka

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"game/internal/config"

	"github.com/segmentio/kafka-go"
)

// tách "a,b,c" -> []string{"a","b","c"}, fallback mặc định
func tachBrokers(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}
	if len(out) == 0 {
		out = []string{"localhost:9093"}
	}
	return out
}

func TaoTopic(ctx context.Context, brokers []string, topic string, partitions, replication int) (bool, error) {
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return false, fmt.Errorf("kết nối Kafka lỗi: %v", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return false, fmt.Errorf("lấy controller lỗi: %v", err)
	}

	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return false, fmt.Errorf("kết nối controller lỗi: %v", err)
	}
	defer controllerConn.Close()

	// ---- Kiểm tra topic đã tồn tại chưa ----
	partitionsInfo, err := controllerConn.ReadPartitions()
	if err != nil {
		return false, fmt.Errorf("đọc partitions lỗi: %v", err)
	}
	for _, p := range partitionsInfo {
		if p.Topic == topic {
			fmt.Printf("ℹ Topic '%s' đã tồn tại, bỏ qua tạo mới\n", topic)
			return false, nil // trả về false = không tạo mới
		}
	}

	// ---- Tạo topic mới ----
	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replication,
	})
	if err != nil {
		return false, fmt.Errorf("tạo topic lỗi: %v", err)
	}

	fmt.Printf("✅ Topic '%s' đã được tạo với %d partition(s)\n", topic, partitions)
	return true, nil // trả về true = đã tạo
}

// ====== WRITER ======
type KafkaWriters struct {
	Created *kafka.Writer
	Updated *kafka.Writer
}

// TaoKafkaWriters sẽ tạo sẵn 2 writer cho 2 topic cố định
func TaoKafkaWriters(ctx context.Context, cfg config.CauHinh) (*KafkaWriters, error) {
	createdWriter, err := taoWriter(ctx, cfg, cfg.KafkaTopicCreate, 3, cfg.KafkaReplication)
	if err != nil {
		return nil, err
	}

	updatedWriter, err := taoWriter(ctx, cfg, cfg.KafkaTopicUpdate, 1, cfg.KafkaReplication)
	if err != nil {
		createdWriter.Close()
		return nil, err
	}

	return &KafkaWriters{
		Created: createdWriter,
		Updated: updatedWriter,
	}, nil
}

// Hàm nội bộ dùng để tạo 1 writer cho topic cụ thể
func taoWriter(ctx context.Context, cfg config.CauHinh, topic string, partitions, replication int) (*kafka.Writer, error) {
	brokers := tachBrokers(cfg.KafkaBrokers)

	created, err := TaoTopic(ctx, brokers, topic, partitions, replication)
	if err != nil {
		return nil, err
	}
	if created {
		log.Printf("[kafka] Topic '%s' được tạo (partitions=%d, replication=%d)\n", topic, partitions, replication)
	} else {
		log.Printf("[kafka] Topic '%s' đã tồn tại\n", topic)
	}

	w := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.Hash{},
		RequiredAcks:           kafka.RequireAll,
		AllowAutoTopicCreation: false,
		WriteTimeout:           10 * time.Second,
		ReadTimeout:            10 * time.Second,
	}

	return w, nil
}

func XoaTopic(brokers []string, topic string) error {
	// Kết nối tới broker đầu tiên
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return fmt.Errorf("không kết nối Kafka: %v", err)
	}
	defer conn.Close()

	// Lấy controller
	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("lấy controller lỗi: %v", err)
	}

	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return fmt.Errorf("kết nối controller lỗi: %v", err)
	}
	defer controllerConn.Close()

	// Xóa topic
	err = controllerConn.DeleteTopics(topic)
	if err != nil {
		return fmt.Errorf("xóa topic '%s' lỗi: %v", topic, err)
	}

	fmt.Printf("✅ Topic '%s' đã được xóa thành công\n", topic)
	return nil
}
