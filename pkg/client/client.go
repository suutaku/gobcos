package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/suutaku/gobcos/pkg/option"
)

const (
	GETE_WAY = "wbbc-dappgetway"
	MESSAGE  = "message"
)

type Dapp struct {
	Token            string `json:"token"`
	ChainId          string `json:"chainId"`
	GroupId          string `json:"groupId"`
	AppId            string `json:"appId"`
	AppName          string `json:"appName"`
	IsV3Chain        bool   `json:"isV3Chain"`
	ChainEncryptType int    `json:"chainEncryptType"`
}

type Client struct {
	c        *http.Client
	dapp     *Dapp
	host     string
	bizSeqNo string
}

func NewClient(opts ...option.ClientOption) *Client {
	cOpt := option.PrepareOpts(opts)
	ret := &Client{

		c:    cOpt.Client,
		host: cOpt.Host,
	}
	if ret.c == nil {
		ret.c = http.DefaultClient
	}
	response := make(map[string]interface{})
	err := ret.doRequest(ret.getSaltAddr(), nil, map[string]string{"rand-num": fmt.Sprintf("%v", int32(time.Now().Unix()))}, &response)
	if err != nil {
		panic(err)
	}

	salt, err := strconv.ParseInt(response["resp"].(string), 10, 32)
	if err != nil {
		panic(err)
	}

	p2 := map[string]interface{}{
		"rand-num": salt,
		"app-id":   cOpt.AppId,
		"app-key":  cOpt.Appkey,
	}
	err = ret.doRequest(ret.loginAddr(), p2, nil, &response)
	if err != nil {
		panic(err)
	}

	tmp, err := json.Marshal(response["resp"])
	if err != nil {
		panic(err)
	}
	ret.dapp = &Dapp{}
	err = json.Unmarshal(tmp, ret.dapp)
	if err != nil {
		panic(err)
	}
	ret.bizSeqNo = response["bizSeqNo"].(string)
	go func() {
		for {
			<-time.After(25 * time.Minute)
			payload := map[string]interface{}{
				"access-token": ret.dapp.Token,
				"exStatus":     "true",
				"bizSeqNo":     ret.bizSeqNo,
			}
			response := make(map[string]interface{})
			if err := ret.doRequest(ret.publicRequestAddr(), payload, nil, &response); err != nil {
				fmt.Println(err)
				continue
			}
			ret.bizSeqNo = response["bizSeqNo"].(string)
		}
	}()
	return ret
}
func (c *Client) getSaltAddr() string {
	ret, err := url.JoinPath(c.host, GETE_WAY, MESSAGE, "getsalt")
	if err != nil {
		panic(err)
	}
	return ret
}

func (c *Client) loginAddr() string {
	ret, err := url.JoinPath(c.host, GETE_WAY, MESSAGE, "login")
	if err != nil {
		panic(err)
	}
	return ret
}

func (c *Client) publicRequestAddr() string {
	ret, err := url.JoinPath(c.host, GETE_WAY, MESSAGE, "gateway")
	if err != nil {
		panic(err)
	}
	return ret
}

func readMessage(r io.Reader) string {
	tmp := make(map[string]interface{})
	if err := readPayload(r, tmp); err != nil {
		return err.Error()
	}
	if tmp["message"] != nil {
		return tmp["message"].(string)
	}
	return "unknown message"
}

func readPayload(r io.Reader, dist interface{}) error {
	if err := json.NewDecoder(r).Decode(dist); err != nil {
		return err
	}
	return nil
}

func (c *Client) doRequest(url string, payload map[string]interface{}, headr map[string]string, dist interface{}) error {
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return err
	}
	for k, v := range headr {
		req.Header.Add(k, v)
	}
	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request faild [%v]: %v", resp.StatusCode, readMessage(resp.Body))
	}
	return readPayload(resp.Body, dist)
}
