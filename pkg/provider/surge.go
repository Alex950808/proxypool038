package provider

import (
	"strings"

	"github.com/Sansui233/proxypool/pkg/tool"

	"github.com/Sansui233/proxypool/pkg/proxy"
)

type Surge struct {
	Base
}

func (s Surge) Provide() string {
	s.preFilter()

	var resultBuilder strings.Builder
	for _, p := range *s.Proxies {
		if checkSurgeSupport(p) {
			resultBuilder.WriteString(p.ToClash() + "\n")
		}
	}
	return resultBuilder.String()
}

func checkSurgeSupport(p proxy.Proxy) bool {
	switch p.(type) {
	case *proxy.ShadowsocksR:
		return false
	case *proxy.Vmess:
		return true
	case *proxy.Shadowsocks:
		ss := p.(*proxy.Shadowsocks)
		if tool.CheckInList(proxy.SSCipherList, ss.Cipher) {
			return true
		}
	default:
		return false
	}
	return false
}
