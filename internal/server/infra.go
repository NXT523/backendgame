package server

import (
	"context"
	"game/ent"
	"game/internal/db"
	"game/pkg/cache"
	"game/pkg/elastic"
	kconsumer "game/pkg/kafka"
	"game/pkg/logx"
	"game/pkg/memcached"
	"log"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/redis/go-redis/v9"
)

type Infra struct {
	Elastic *logx.LoggerElastic
	Redis   *redis.Client
	Memc    *memcache.Client
	Kafka   *kconsumer.KafkaWriters
	Ent     *ent.Client
}

func InitInfra(ctx context.Context) (*Infra, error) {
	elog, err := elastic.TaoElasticLogger(ctx)
	if err != nil {
		return nil, err
	}

	rdb, err := cache.NewRedisClient(ctx)
	if err != nil {
		return nil, err
	}

	mc, err := memcached.NewMemcacheClient(ctx)
	if err != nil {
		return nil, err
	}

	kafkaWriter, err := kconsumer.TaoKafkaWriters()
	if err != nil {
		return nil, err
	}

	if err := db.CreateDBEnt(); err != nil {
		log.Fatalf("[ENT] Tạo DB lỗi: %v", err)
	}

	entClient, err := db.OpenEnt()
	if err != nil {
		return nil, err
	}

	if err := db.RunMigrationEnt(context.Background(), entClient); err != nil {
		log.Fatalf("[ENT] Migration lỗi: %v", err)
	}

	return &Infra{
		Elastic: elog,
		Redis:   rdb,
		Memc:    mc,
		Kafka:   kafkaWriter,
		Ent:     entClient,
	}, nil
}

func (i *Infra) Close() {
	if i == nil {
		return
	}

	if i.Redis != nil {
		_ = i.Redis.Close()
	}

	if i.Ent != nil {
		_ = i.Ent.Close()
	}

	if i.Memc != nil {
		_ = i.Memc.Close()
	}

	if i.Kafka != nil {
		i.Kafka.CloseAll()
	}
}
