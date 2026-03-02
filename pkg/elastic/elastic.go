// pkg/elastic/elastic.go
package elastic

import (
	"context"
	"game/pkg/logx"
	"os"
)

// Tạo logger Elasticsearch đọc từ config.
func TaoElasticLogger(_ context.Context) (*logx.LoggerElastic, error) {
	return logx.TaoLoggerElastic(logx.Opts{
		ESURL:   os.Getenv("ES_URL"),
		Service: os.Getenv("ES_SERVICE"),
		Env:     os.Getenv("ES_ENV"),
		Index:   os.Getenv("ES_INDEX"),
	})
}
