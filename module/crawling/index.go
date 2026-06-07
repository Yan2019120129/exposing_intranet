package crawling

import (
	"net/http"

	"github.com/gocolly/colly"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

const (
	// LogMsg 日志信息标识
	LogMsg = "crawling"
)

type Crawling struct {
	*colly.Collector
	rdsConn redis.Conn
	err     error
}

// NewCrawling 新建实例
func NewCrawling() (*Crawling, error) {
	amazon := &Crawling{
		Collector: colly.NewCollector(),
	}
	amazon.Async = true
	amazon.AllowURLRevisit = true

	amazon.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	if err := amazon.SetStorage(NewStorage("crawling")); err != nil {
		zap.L().Error(LogMsg, zap.Error(err))
		return nil, err
	}

	amazon.OnError(func(rsp *colly.Response, err error) {
		zap.L().Error(LogMsg, zap.Error(err), zap.String("body", string(rsp.Body)))
		amazon.err = err
	})

	return amazon, amazon.err
}

func (_Crawling *Crawling) SetRequestHeaders(header map[string]string) *Crawling {
	_Crawling.OnRequest(func(r *colly.Request) {
		for k, v := range header {
			r.Headers.Set(k, v)
		}
	})
	return _Crawling
}

func (_Crawling *Crawling) SetOnHtml(selector string, fn func(e *colly.HTMLElement)) *Crawling {
	_Crawling.OnHTML(selector, fn)
	return _Crawling
}

func (_Crawling *Crawling) SetOnResponse(fn func(response *colly.Response)) *Crawling {
	_Crawling.OnResponse(fn)
	return _Crawling
}

func (_Crawling *Crawling) SetOnScraped(fn func(response *colly.Response)) *Crawling {
	_Crawling.OnScraped(fn)
	return _Crawling
}

// Run 执行程序
func (_Crawling *Crawling) Run(rwaUrl string) error {
	if err := _Crawling.Visit(rwaUrl); err != nil {
		return err
	}
	return _Crawling.err
}
