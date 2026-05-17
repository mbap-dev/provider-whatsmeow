package provider

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	ilog "your.org/provider-whatsmeow/internal/log"
)

// ProxyConfig mirrors Evolution Go's proxy model: protocol/host/port plus
// optional credentials. Protocol may be inferred from the port when omitted.
type ProxyConfig struct {
	Protocol string `json:"protocol,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Enabled  bool   `json:"enabled"`
}

var defaultProxyConfig struct {
	sync.RWMutex
	cfg ProxyConfig
}

func SetDefaultProxyConfig(cfg ProxyConfig) {
	defaultProxyConfig.Lock()
	defaultProxyConfig.cfg = normalizeProxyConfig(cfg)
	if defaultProxyConfig.cfg.Host != "" && defaultProxyConfig.cfg.Port != "" {
		defaultProxyConfig.cfg.Enabled = true
	}
	defaultProxyConfig.Unlock()
}

func defaultProxy() ProxyConfig {
	defaultProxyConfig.RLock()
	cfg := defaultProxyConfig.cfg
	defaultProxyConfig.RUnlock()
	return cfg
}

func normalizeProxyProtocol(protocol, port string) string {
	normalized := strings.ToLower(strings.TrimSpace(protocol))
	switch normalized {
	case "socks":
		return "socks5"
	case "http", "https", "socks5":
		return normalized
	}

	switch strings.TrimSpace(port) {
	case "1080", "2080":
		return "socks5"
	}
	if portNum, err := strconv.Atoi(strings.TrimSpace(port)); err == nil && portNum >= 42000 && portNum <= 43000 {
		return "socks5"
	}
	return "http"
}

func buildProxyAddress(protocol, host, port, username, password string) (string, error) {
	host = strings.TrimSpace(host)
	port = strings.TrimSpace(port)
	if host == "" {
		return "", fmt.Errorf("proxy host is required")
	}
	if port == "" {
		return "", fmt.Errorf("proxy port is required")
	}

	protocol = normalizeProxyProtocol(protocol, port)
	if protocol != "http" && protocol != "https" && protocol != "socks5" {
		return "", fmt.Errorf("unsupported proxy protocol %q", protocol)
	}

	u := &url.URL{
		Scheme: protocol,
		Host:   net.JoinHostPort(host, port),
	}
	username = strings.TrimSpace(username)
	if username != "" {
		if password != "" {
			u.User = url.UserPassword(username, password)
		} else {
			u.User = url.User(username)
		}
	}
	return u.String(), nil
}

func normalizeProxyConfig(cfg ProxyConfig) ProxyConfig {
	cfg.Protocol = normalizeProxyProtocol(cfg.Protocol, cfg.Port)
	cfg.Host = strings.TrimSpace(cfg.Host)
	cfg.Port = strings.TrimSpace(cfg.Port)
	cfg.Username = strings.TrimSpace(cfg.Username)
	return cfg
}

func redactedProxyAddress(cfg ProxyConfig) string {
	cfg = normalizeProxyConfig(cfg)
	host := strings.TrimSpace(cfg.Host)
	port := strings.TrimSpace(cfg.Port)
	if host == "" || port == "" {
		return ""
	}
	return normalizeProxyProtocol(cfg.Protocol, cfg.Port) + "://" + net.JoinHostPort(host, port)
}

func httpClientForProxyConfig(sessionID string, cfg ProxyConfig) *http.Client {
	cfg = normalizeProxyConfig(cfg)
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment}
	if cfg.Enabled && cfg.Host != "" && cfg.Port != "" {
		transport.Proxy = func(req *http.Request) (*url.URL, error) {
			addr, err := buildProxyAddress(cfg.Protocol, cfg.Host, cfg.Port, cfg.Username, cfg.Password)
			if err != nil {
				ilog.WithSession(sessionID).Debug(
					"proxy_http_route error target=%s protocol=%s host=%s port=%s err=%v",
					req.URL.Host,
					normalizeProxyProtocol(cfg.Protocol, cfg.Port),
					cfg.Host,
					cfg.Port,
					err,
				)
				return nil, err
			}
			u, err := url.Parse(addr)
			if err != nil {
				ilog.WithSession(sessionID).Debug(
					"proxy_http_route parse_error target=%s proxy=%s err=%v",
					req.URL.Host,
					redactedProxyAddress(cfg),
					err,
				)
				return nil, err
			}
			ilog.WithSession(sessionID).Debug(
				"proxy_http_route enabled target=%s proxy=%s username_set=%t password_set=%t",
				req.URL.Host,
				redactedProxyAddress(cfg),
				strings.TrimSpace(cfg.Username) != "",
				cfg.Password != "",
			)
			return u, nil
		}
	} else {
		ilog.WithSession(sessionID).Debug(
			"proxy_http_route disabled using_environment_proxy enabled=%t host_set=%t port_set=%t",
			cfg.Enabled,
			cfg.Host != "",
			cfg.Port != "",
		)
	}
	return &http.Client{Transport: transport}
}

func proxyFuncFromConfigOrEnvironment(req *http.Request) (*url.URL, error) {
	cfg := defaultProxy()
	if cfg.Enabled && strings.TrimSpace(cfg.Host) != "" && strings.TrimSpace(cfg.Port) != "" {
		addr, err := buildProxyAddress(cfg.Protocol, cfg.Host, cfg.Port, cfg.Username, cfg.Password)
		if err != nil {
			return nil, err
		}
		u, err := url.Parse(addr)
		if err != nil {
			return nil, err
		}
		ilog.Debugf(
			"proxy_http_route enabled target=%s proxy=%s username_set=%t password_set=%t source=default",
			req.URL.Host,
			redactedProxyAddress(cfg),
			strings.TrimSpace(cfg.Username) != "",
			cfg.Password != "",
		)
		return u, nil
	}
	return http.ProxyFromEnvironment(req)
}
