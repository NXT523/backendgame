// pkg/elastic/elastic.go
package elastic

import (
	"context"

	"game/internal/config"
	"game/pkg/logx"
)

// Tạo logger Elasticsearch đọc từ config.
func TaoElasticLogger(_ context.Context, cfg config.CauHinh) (*logx.LoggerElastic, error) {
	return logx.TaoLoggerElastic(logx.Opts{
		ESURL:   cfg.ESURL,
		Service: cfg.ESService,
		Env:     cfg.ESEnv,
		Index:   cfg.ESIndex,
	})
}
