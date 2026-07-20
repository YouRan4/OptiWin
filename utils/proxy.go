//go:build windows

package utils

import (
	"net/http"
	"net/url"
	"strings"
	"time"
	"golang.org/x/sys/windows/registry"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		Proxy: getProxy,
	},
}

func GetHttpClient() *http.Client {
	return httpClient
}

func GetProxyInfo() string {
	u := registryProxy()
	if u != nil {
		return u.Host
	}
	return "未使用代理"
}

func getProxy(req *http.Request) (*url.URL, error) {
	u := registryProxy()
	if u != nil {
		return u, nil
	}
	return nil, nil
}

func registryProxy() *url.URL {
	k, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Internet Settings`,
		registry.QUERY_VALUE)
	if err != nil {
		return nil
	}
	defer k.Close()

	enabled, _, _ := k.GetIntegerValue("ProxyEnable")
	if enabled == 0 {
		return nil
	}

	proxyServer, _, _ := k.GetStringValue("ProxyServer")
	if proxyServer == "" {
		return nil
	}

	if strings.Contains(proxyServer, "=") {
		for _, part := range strings.Split(proxyServer, ";") {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "http=") {
				proxyServer = strings.TrimPrefix(part, "http=")
				break
			} else if strings.HasPrefix(part, "https=") {
				proxyServer = strings.TrimPrefix(part, "https=")
				break
			}
		}
	}
	if !strings.HasPrefix(proxyServer, "http://") && !strings.HasPrefix(proxyServer, "https://") {
		return &url.URL{Host: proxyServer, Scheme: "http"}
	}
	u, err := url.Parse(proxyServer)
	if err != nil {
		return nil
	}
	return u
}
