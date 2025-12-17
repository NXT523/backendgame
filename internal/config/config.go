package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type CauHinh struct {
	DSNChuaDB string // ví dụ: root:123456@tcp(localhost:3306)/
	TenDBEnt  string // ví dụ: game
	TenDBGorm string // ví dụ: game1
	TenDBRaw  string // ví dụ: game2
	// HTTPPortEnt  string // ví dụ: :8080
	HTTPPortEnt1 string // ví dụ: :8060
	HTTPPortEnt2 string // ví dụ: :8070
	HTTPPortEnt3 string // ví dụ: :8080
	HTTPPortGorm string // ví dụ: :8081
	HTTPPortRaw  string // ví dụ: :8082

	// Proxy
	ProxyPort string // ví dụ: ":9000"

	// --- REDIS ---
	RedisAddr          string        // ví dụ: localhost:6379
	RedisUsername      string        // nếu có
	RedisPassword      string        // nếu có
	RedisDB            int           // ví dụ: 0
	RedisProtocol      int           // 2 cho RESP2, 3 cho RESP3 (go-redis v9 dùng 2|3)
	RedisTLS           bool          // bật TLS nếu true
	RedisTLSSkipVerify bool          // bỏ qua verify cert (chỉ dùng khi cần)
	RedisDialTimeout   time.Duration // ví dụ: 5s
	RedisReadTimeout   time.Duration // ví dụ: 3s
	RedisWriteTimeout  time.Duration // ví dụ: 3s

	// Memcached
	MemcachedAddr         string
	MemcachedTimeout      time.Duration
	MemcachedMaxIdleConns int
	MemcachedHealthKey    string
	MemcachedHealthTTL    time.Duration

	// --- ELASTIC ---
	ESURL     string // ví dụ: http://localhost:9200
	ESService string // ví dụ: game-app
	ESEnv     string // ví dụ: dev
	ESIndex   string // ví dụ: app-logs

	// --- KAFKA ---
	KafkaBrokers         string // ví dụ: "localhost:9093,localhost:9094"
	KafkaAllowAutoCreate bool   // cho phép auto-create topic nếu cluster cho phép
	KafkaTopic           string // ví dụ: vu_khi_created
	KafkaPartitions      int    // ví dụ: 3
	KafkaReplication     int    // ví dụ: 2
}

func DocCauHinh() CauHinh {
	return CauHinh{
		DSNChuaDB:    env("MYSQL_DSN_NO_DB", "root:123456@tcp(localhost:3306)/"),
		TenDBEnt:     env("MYSQL_DB_ENT", "game"),
		TenDBGorm:    env("MYSQL_DB_GORM", "game1"),
		TenDBRaw:     env("MYSQL_DB_RAW", "game2"),
		HTTPPortEnt1: env("HTTP_PORT_ENT1", ":8060"),
		HTTPPortEnt2: env("HTTP_PORT_ENT2", ":8070"),
		HTTPPortEnt3: env("HTTP_PORT_ENT3", ":8080"),
		HTTPPortGorm: env("HTTP_PORT_GORM", ":8081"),
		HTTPPortRaw:  env("HTTP_PORT_RAW", ":8082"),

		// Proxy
		ProxyPort: env("PROXY_PORT", ":9000"),

		// --- REDIS (ENV override được hết) ---
		RedisAddr:          env("REDIS_ADDR", "localhost:6379"), //"redis-16269.crce214.us-east-1-3.ec2.redns.redis-cloud.com:16269"
		RedisUsername:      env("REDIS_USERNAME", "default"),
		RedisPassword:      env("REDIS_PASSWORD", "abc"),
		RedisDB:            envInt("REDIS_DB", 0),
		RedisProtocol:      envInt("REDIS_PROTOCOL", 2), // RESP2
		RedisTLS:           envBool("REDIS_TLS", false), // URL dùng redis:// (không TLS). Nếu cần TLS, set true
		RedisTLSSkipVerify: envBool("REDIS_TLS_SKIP_VERIFY", false),
		RedisDialTimeout:   envDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
		RedisReadTimeout:   envDuration("REDIS_READ_TIMEOUT", 3*time.Second),
		RedisWriteTimeout:  envDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),

		// --- MEMCACHED --- 👈 THÊM MỚI
		MemcachedAddr:         env("MEMCACHED_ADDR", "127.0.0.1:11211"),
		MemcachedTimeout:      envDuration("MEMCACHED_TIMEOUT", 200*time.Millisecond),
		MemcachedMaxIdleConns: envInt("MEMCACHED_MAX_IDLE_CONNS", 100),

		// --- ELASTIC ---
		ESURL:     env("ES_URL", "http://localhost:9200"),
		ESService: env("ES_SERVICE", "game-app"),
		ESEnv:     env("ES_ENV", "dev"),
		ESIndex:   env("ES_INDEX", "app-logs"),

		// --- KAFKA ---
		KafkaBrokers:         env("KAFKA_BROKERS", "localhost:9093"),
		KafkaAllowAutoCreate: envBool("KAFKA_ALLOW_AUTO_CREATE", false),
		KafkaTopic:           env("KAFKA_TOPIC", "vu_khi_created"),
		KafkaPartitions:      envInt("KAFKA_PARTITIONS", 3),
		KafkaReplication:     envInt("KAFKA_REPLICATION", 1),
	}
}

// Tiện ích ghép DSN có DB cho MySQL (dùng chung cho Ent/GORM/Raw)
func DSNCoDB(dsnNoDB, tenDB string) string {
	// Kết quả ví dụ: root:123456@tcp(localhost:3306)/game1?charset=utf8mb4&parseTime=True&loc=Local
	return fmt.Sprintf("%s%s?charset=utf8mb4&parseTime=True&loc=Local", dsnNoDB, tenDB)
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func envInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func envBool(k string, def bool) bool {
	if v := os.Getenv(k); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

func envDuration(k string, def time.Duration) time.Duration {
	if v := os.Getenv(k); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
