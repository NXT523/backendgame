// pkg/logx/logx.go
package logx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	esv8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type LoggerElastic struct {
	es      *esv8.Client
	bulk    esutil.BulkIndexer
	service string // -> service.name
	env     string // -> service.environment
	index   string
}

type Opts struct {
	ESURL   string
	Service string
	Env     string
	Index   string
}

func TaoLoggerElastic(o Opts) (*LoggerElastic, error) {
	if o.Index == "" {
		o.Index = "app-logs"
	}
	es, err := esv8.NewClient(esv8.Config{Addresses: []string{o.ESURL}})
	if err != nil {
		return nil, err
	}
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:        es,
		Index:         o.Index,
		FlushBytes:    5 << 20,
		FlushInterval: 2 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &LoggerElastic{
		es:      es,
		bulk:    bi,
		service: o.Service,
		env:     o.Env,
		index:   o.Index,
	}, nil
}

func (l *LoggerElastic) DongBo(ctx context.Context) error {
	return l.bulk.Close(ctx)
}

func (l *LoggerElastic) Ghi(ctx context.Context, level, message string, kv map[string]any) {
	doc := map[string]any{
		"@timestamp":          time.Now().UTC().Format(time.RFC3339Nano),
		"service.name":        l.service,
		"service.environment": l.env,
		"level":               level,   // "info" | "warn" | "error"
		"message":             message, // ví dụ: "HTTP request handled"
	}
	for k, v := range kv {
		if k == "error" {
			if _, ok := v.(map[string]any); !ok {
				// bỏ qua hoặc chuyển sang object rỗng
				v = map[string]any{}
			}
		}
		doc[k] = v
	}
	body, _ := json.Marshal(doc)
	_ = l.bulk.Add(ctx, esutil.BulkIndexerItem{
		Action: "index",
		Body:   bytes.NewReader(body),
		OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, resp esutil.BulkIndexerResponseItem, err error) {
			fmt.Println("❌ Lỗi bulk:", err, resp.Error)
		},
	})
}

func (l *LoggerElastic) Info(ctx context.Context, msg string, kv map[string]any) {
	l.Ghi(ctx, "info", msg, kv)
}
func (l *LoggerElastic) Warn(ctx context.Context, msg string, kv map[string]any) {
	l.Ghi(ctx, "warn", msg, kv)
}
func (l *LoggerElastic) Error(ctx context.Context, msg string, kv map[string]any) {
	l.Ghi(ctx, "error", msg, kv)
}
