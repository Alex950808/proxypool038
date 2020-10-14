package getter

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	"github.com/Alex950808/proxypool038/pkg/proxy"
	"github.com/Alex950808/proxypool038/pkg/tool"
)

func init() {
	Register("web-fanqiangdang", NewWebFanqiangdangGetter)
}

type WebFanqiangdang struct {
	c       *colly.Collector
	Url     string
	results proxy.ProxyList
}

func NewWebFanqiangdangGetter(options tool.Options) (getter Getter, err error) {
	urlInterface, found := options["url"]
	if found {
		url, err := AssertTypeStringNotNull(urlInterface)
		if err != nil {
			return nil, err
		}
		return &WebFanqiangdang{
			c:   colly.NewCollector(),
			Url: url,
		}, nil
	}
	return nil, ErrorUrlNotFound
}

func (w *WebFanqiangdang) Get() proxy.ProxyList {
	w.results = make(proxy.ProxyList, 0)
	w.c.OnHTML("td.t_f", func(e *colly.HTMLElement) {
		w.results = append(w.results, FuzzParseProxyFromString(e.Text)...)
		subUrls := urlRe.FindAllString(e.Text, -1)
		for _, url := range subUrls {
			w.results = append(w.results, (&Subscribe{Url: url}).Get()...)
		}
	})

	w.c.OnHTML("th.new>a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		if strings.HasPrefix(url, "https://fanqiangdang.com/thread") {
			_ = e.Request.Visit(url)
		}
	})

	w.results = make(proxy.ProxyList, 0)
	err := w.c.Visit(w.Url)
	if err != nil {
		_ = fmt.Errorf("%s", err.Error())
	}

	return w.results
}

func (w *WebFanqiangdang) Get2Chan(pc chan proxy.Proxy, wg *sync.WaitGroup) {
	defer wg.Done()
	nodes := w.Get()
	log.Printf("STATISTIC: Fanqiangdang\tcount=%d\turl=%s\n", len(nodes), w.Url)
	for _, node := range nodes {
		pc <- node
	}
}
