package client

import "net/url"

const (
	GETE_WAY = "wbbc-dappgateway"
	MESSAGE  = "message"
)

type Method func() string

func (c *Client) GETSALT() string {
	ret, err := url.JoinPath(c.host, GETE_WAY, MESSAGE, "getsalt")
	if err != nil {
		panic(err)
	}
	return ret
}

func (c *Client) LOGIN() string {
	ret, err := url.JoinPath(c.host, GETE_WAY, MESSAGE, "login")
	if err != nil {
		panic(err)
	}
	return ret
}

func (c *Client) PUBLICREQUEST() string {
	ret, err := url.JoinPath(c.host, GETE_WAY, MESSAGE, "gateway")
	if err != nil {
		panic(err)
	}
	return ret
}
