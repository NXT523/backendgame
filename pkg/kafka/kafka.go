package kafka

import (
	"context"
	"fmt"
	"time"

	"game/internal/utils"

	"github.com/segmentio/kafka-go"
)

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
			return false, nil
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

type KafkaWriters struct {
	writers map[KafkaWriterKey]*kafka.Writer
}
type KafkaWriterKey string

const (
	KafkaCreate KafkaWriterKey = "create"
	KafkaUpdate KafkaWriterKey = "update"
)

func TaoKafkaWriters() (*KafkaWriters, error) {
	replication := utils.GetenvInt("KAFKA_REPLICATION", 1)

	configs := map[KafkaWriterKey]struct {
		topic     string
		partition int
	}{
		KafkaCreate: {
			topic:     utils.GetenvString("KAFKA_TOPIC_CREATE", ""),
			partition: 3,
		},
		KafkaUpdate: {
			topic:     utils.GetenvString("KAFKA_TOPIC_UPDATE", ""),
			partition: 1,
		},
	}

	writers := make(map[KafkaWriterKey]*kafka.Writer)

	for key, cfg := range configs {
		w, err := taoWriter(cfg.topic, cfg.partition, replication)
		if err != nil {
			for _, ww := range writers {
				ww.Close()
			}
			return nil, err
		}
		writers[key] = w
	}

	return &KafkaWriters{writers: writers}, nil
}

// Hàm nội bộ dùng để tạo 1 writer cho topic cụ thể
func taoWriter(topic string, partitions, replication int) (*kafka.Writer, error) {
	brokers := utils.GetenvStrings("KAFKA_BROKERS", nil)
	if len(brokers) == 0 {
		brokers = []string{"localhost:9093"}
	}

	_, err := TaoTopic(context.Background(), brokers, topic, partitions, replication)
	if err != nil {
		return nil, err
	}
	if topic == "" {
		return nil, fmt.Errorf("Kafka topic không được rỗng")
	}
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}, nil
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

func (k *KafkaWriters) CloseAll() {
	for _, w := range k.writers {
		_ = w.Close()
	}
}
