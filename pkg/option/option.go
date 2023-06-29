package option

import "net/http"

type clientOption struct {
	AppId    string
	Appkey   string
	Client   *http.Client
	Host     string
	SignUser string
}

type ClientOption func(opts *clientOption)

func WithHost(host string) ClientOption {
	return func(opts *clientOption) {
		opts.Host = host
	}
}

func WithAppId(appId string) ClientOption {
	return func(opts *clientOption) {
		opts.AppId = appId
	}
}

func WithAppKey(appKey string) ClientOption {
	return func(opts *clientOption) {
		opts.Appkey = appKey
	}
}

func WithHTTPClient(c *http.Client) ClientOption {
	return func(opts *clientOption) {
		opts.Client = c
	}
}

func WithSignUser(u string) ClientOption {
	return func(opts *clientOption) {
		opts.SignUser = u
	}
}

func PrepareOpts(opts []ClientOption) *clientOption {
	ret := &clientOption{}
	for _, v := range opts {
		v(ret)
	}
	return ret
}
