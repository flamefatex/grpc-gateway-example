package myclient

import (
	"sync"

	"github.com/flamefatex/grpc-gateway-example/pkg/lib/config"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/log"
	"github.com/go-resty/resty/v2"
)

var (
	instance *MyClient // 单例实体
	once     sync.Once
)

type MyClient struct {
	httpClient *resty.Client
}

func GetInstance() *MyClient {
	once.Do(func() {
		httpClient := resty.New()
		httpClient.SetLogger(log.GetLogger())
		httpClient.SetDebug(config.Config().GetBool("http.debug"))
		httpClient.SetTimeout(config.Config().GetDuration("http.timeout"))

		instance = &MyClient{
			httpClient: httpClient,
		}
	})
	return instance
}
