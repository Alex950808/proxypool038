package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ivpusic/grpool"

	"github.com/Dreamacro/clash/adapters/outbound"
)

const defaultURLTestTimeout = time.Second * 5

func testDelay(p Proxy) (delay uint16, err error) {
	pmap := make(map[string]interface{})
	err = json.Unmarshal([]byte(p.String()), &pmap)
	if err != nil {
		return
	}

	pmap["port"] = int(pmap["port"].(float64))
	if p.TypeName() == "vmess" {
		pmap["alterId"] = int(pmap["alterId"].(float64))
	}

	clashProxy, err := outbound.ParseProxy(pmap)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultURLTestTimeout)
	delay, err = clashProxy.URLTest(ctx, "http://www.gstatic.com/generate_204")
	cancel()
	return delay, err
}

func CleanBadProxiesWithGrpool(proxies []Proxy) (cproxies []Proxy) {
	// Note: Grpool实现对go并发管理的封装，主要是在数据量大时减少内存占用，不会提高效率。
	pool := grpool.NewPool(500, 200)

	c := make(chan checkResult)
	defer close(c)

	pool.WaitCount(len(proxies))
	// 线程：延迟测试，测试过程通过grpool的job并发
	go func() {
		for _, p := range proxies {
			pp := p // 复制一份，否则job执行时是按当前的p测试的
			pool.JobQueue <- func() {
				defer pool.JobDone()
				delay, err := testDelay(pp)
				if err == nil {
					c <- checkResult{
						name:  pp.Identifier(),
						delay: delay,
					}
				}
			}
		}
	}()
	done := make(chan struct{}) // 用于多线程的运行结束标识
	defer close(done)

	go func() {
		pool.WaitAll()
		pool.Release()
		done <- struct{}{}
	}()

	okMap := make(map[string]struct{})
	for { // Note: 无限循环，直到能读取到done。处理并发也算是挺有创意的写法
		select {
		case r := <-c:
			if r.delay > 0 {
				okMap[r.name] = struct{}{}
			}
		case <-done:
			cproxies = make(ProxyList, 0, 500) // 定义返回的proxylist
			for _, p := range proxies {
				if _, ok := okMap[p.Identifier()]; ok {
					cproxies = append(cproxies, p.Clone())
				}
			}
			return
		}
	}
}

func CleanBadProxies(proxies []Proxy) (cproxies []Proxy) {
	c := make(chan checkResult, 40)
	wg := &sync.WaitGroup{}
	wg.Add(len(proxies))
	for _, p := range proxies {
		go testProxyDelayToChan(p, c, wg)
	}
	go func() {
		wg.Wait()
		close(c)
	}()

	okMap := make(map[string]struct{})
	for r := range c {
		if r.delay > 0 {
			okMap[r.name] = struct{}{}
		}
	}
	cproxies = make(ProxyList, 0, 500)
	for _, p := range proxies {
		if _, ok := okMap[p.Identifier()]; ok {
			p.SetUseable(true)
			cproxies = append(cproxies, p.Clone())
		} else {
			p.SetUseable(false)
		}
	}
	return
}

type checkResult struct {
	name  string
	delay uint16
}

func testProxyDelayToChan(p Proxy, c chan checkResult, wg *sync.WaitGroup) {
	defer wg.Done()
	delay, err := testDelay(p)
	if err == nil {
		c <- checkResult{
			name:  p.Identifier(),
			delay: delay,
		}
	}
}
