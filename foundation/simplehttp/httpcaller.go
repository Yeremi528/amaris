// A Resty wrapper
package simplehttp

import (
	"net"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	*resty.Client
}

type Config struct {
	BaseURL          string
	RetryCount       int
	RetryMaxWaitTime time.Duration
	Timeout          time.Duration
}

func New(cfg Config) *Client {
	c := resty.New()

	c.SetBaseURL(cfg.BaseURL)
	c.SetTimeout(cfg.Timeout)
	c.SetRetryCount(cfg.RetryCount)
	c.SetRetryMaxWaitTime(cfg.RetryMaxWaitTime)
	c.SetLogger(RestyLogger{})

	return &Client{c}
}

type PrettyTraceInfo struct {
	DNSLookup      string
	ConnTime       string
	TCPConnTime    string
	TLSHandshake   string
	ServerTime     string
	ResponseTime   string
	TotalTime      string
	IsConnReused   bool
	IsConnWasIdle  bool
	ConnIdleTime   string
	RequestAttempt int
	RemoteAddr     net.Addr
}

func FormatTraceInfo(tc resty.TraceInfo) PrettyTraceInfo {
	return PrettyTraceInfo{
		DNSLookup:      tc.DNSLookup.String(),
		ConnTime:       tc.ConnTime.String(),
		TCPConnTime:    tc.TCPConnTime.String(),
		TLSHandshake:   tc.TLSHandshake.String(),
		ServerTime:     tc.ServerTime.String(),
		ResponseTime:   tc.ResponseTime.String(),
		TotalTime:      tc.TotalTime.String(),
		IsConnReused:   tc.IsConnReused,
		IsConnWasIdle:  tc.IsConnWasIdle,
		ConnIdleTime:   tc.ConnIdleTime.String(),
		RequestAttempt: tc.RequestAttempt,
		RemoteAddr:     tc.RemoteAddr,
	}
}
