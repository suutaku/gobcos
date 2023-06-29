package client

import (
	"fmt"
	"net/http"
	"time"

	"crypto/tls"

	"github.com/sirupsen/logrus"
	"github.com/suutaku/gobcos/pkg/option"
)

type Client struct {
	c        *http.Client
	dapp     *Dapp
	host     string
	bizSeqNo string
	signUser string
}

func NewClient(opts ...option.ClientOption) *Client {
	cOpt := option.PrepareOpts(opts)
	ret := &Client{
		c:        cOpt.Client,
		host:     cOpt.Host,
		signUser: cOpt.SignUser,
	}
	if ret.c == nil {
		ret.c = http.DefaultClient
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	randNum := int32(time.Now().Unix())
	salt, err := ret.GetSalt(randNum)
	if err != nil {
		panic(err)
	}
	if _, err := ret.Login(fmt.Sprintf("%v", randNum), salt, cOpt.AppId, cOpt.Appkey); err != nil {
		panic(err)
	}
	go func() {
		for {
			<-time.After(25 * time.Minute)
			if err := ret.UpdateToken(); err != nil {
				logrus.Error(err)
			}
		}
	}()
	return ret
}

func (c *Client) UpdateToken() error {
	payload := map[string]interface{}{
		"access-token": c.dapp.Token,
		"exStatus":     "true",
		"bizSeqNo":     c.bizSeqNo,
	}
	commonResp, err := c.invoke(c.PUBLICREQUEST, payload, nil)
	if err != nil {
		return err
	}
	c.bizSeqNo = commonResp.BizSeqNo
	return nil
}

func (c *Client) GetSalt(randNum int32) (string, error) {
	salt := ""
	commonResp, err := c.invoke(c.GETSALT, nil, map[string]string{"rand-num": fmt.Sprintf("%v", randNum)})
	if err != nil {
		return salt, err
	}
	err = commonResp.GetPayload(&salt)
	return salt, err
}

func (c *Client) Login(randNum, salt, appId, appKey string) (*Dapp, error) {

	encodeKey, err := encryptAES(appKey, salt)
	if err != nil {
		panic(err)
	}
	p2 := map[string]string{
		"rand-num": randNum,
		"app-id":   appId,
		"app-key":  encodeKey,
	}
	commResp, err := c.invoke(c.LOGIN, nil, p2)
	if err != nil {
		return nil, err
	}

	c.dapp = &Dapp{}
	commResp.GetPayload(c.dapp)
	c.bizSeqNo = commResp.BizSeqNo
	return c.dapp, nil
}

func (c *Client) NewUser(desc string) (*UserResponse, error) {
	req := UserRequest{SignUserId: c.signUser, Description: desc}
	pubReq := PublicRequest{ApiName: "newUserApi", Payload: req}
	cResp, err := c.invoke(c.PUBLICREQUEST, pubReq, c.publicHeader())
	if err != nil {
		return nil, err
	}
	resp := &UserResponse{}
	err = cResp.GetPayload(resp)
	return resp, err
}

func (c *Client) UserInfo() (*UserResponse, error) {
	pubReq := PublicRequest{ApiName: "userInfoApi", Payload: c.signUser}
	cResp, err := c.invoke(c.PUBLICREQUEST, pubReq, c.publicHeader())
	if err != nil {
		return nil, err
	}
	resp := &UserResponse{}
	err = cResp.GetPayload(resp)
	return resp, err
}

func (c *Client) QueryTransaction(req QueryTransactionRequest) ([]byte, error) {
	pubReq := PublicRequest{ApiName: "sendQueryTransactionApi", Payload: req}
	cResp, err := c.invoke(c.PUBLICREQUEST, pubReq, c.publicHeader())
	if err != nil {
		return nil, err
	}
	resp := make([]byte, 0)
	err = cResp.GetPayload(resp)
	return resp, err
}

func (c *Client) GetTransactionByHash(hash string) (*GetTransactionResponse, error) {
	pubReq := PublicRequest{ApiName: "getTransactionByHashApi", Payload: hash}
	cResp, err := c.invoke(c.PUBLICREQUEST, pubReq, c.publicHeader())
	if err != nil {
		return nil, err
	}
	resp := &GetTransactionResponse{}
	err = cResp.GetPayload(resp)
	return resp, err
}

// func (c *Client) GetTransactionReceipt(hash string) (*GetTransactionByHashResponse, error) {
// 	pubReq := PublicRequest{ApiName: "getTransactionReceiptApi", Payload: hash}
// 	cResp, err := c.invoke(c.PUBLICREQUEST, pubReq, c.publicHeader())
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp := &GetTransactionByHashResponse{}
// 	err = cResp.GetPayload(resp)
// 	return resp, err
// }

func (c *Client) GetBlockNumber() (int64, error) {
	pubReq := PublicRequest{ApiName: "getBlockNumber"}
	cResp, err := c.invoke(c.PUBLICREQUEST, pubReq, c.publicHeader())
	if err != nil {
		return 0, err
	}
	resp := int64(0)
	err = cResp.GetPayload(&resp)
	return resp, err
}

func (c *Client) Sign(req SignRequest) (*SignResponse, error) {
	pubReq := PublicRequest{ApiName: "signApi", Payload: req}
	cResp, err := c.invoke(c.PUBLICREQUEST, pubReq, c.publicHeader())
	if err != nil {
		return nil, err
	}
	resp := &SignResponse{}
	err = cResp.GetPayload(&resp)
	return resp, err
}

func (c *Client) Deploy(req DeployRequest) (*DeployResponse, error) {
	pubReq := PublicRequest{ApiName: "deployApi", Payload: req}
	cResp, err := c.invoke(c.PUBLICREQUEST, pubReq, c.publicHeader())
	if err != nil {
		return nil, err
	}
	resp := &DeployResponse{}
	err = cResp.GetPayload(&resp)
	return resp, err
}

func (c *Client) SendTransaction(req SendTransactionRequest) (*SendTransactionResponse, error) {
	pubReq := PublicRequest{ApiName: "sendNewApi", Payload: req}
	cResp, err := c.invoke(c.PUBLICREQUEST, pubReq, c.publicHeader())
	if err != nil {
		return nil, err
	}
	resp := &SendTransactionResponse{}
	err = cResp.GetPayload(&resp)
	return resp, err
}
