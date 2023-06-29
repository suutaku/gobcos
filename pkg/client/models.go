package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mitchellh/mapstructure"
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

type CommonResponse struct {
	Status   int32       `json:"status"`
	Message  string      `json:"message"`
	Payload  interface{} `json:"resp"`
	BizSeqNo string      `json:"bizSeqNo"`
}

func NewCommonResponse(r io.Reader) (*CommonResponse, error) {
	ret := &CommonResponse{}
	if err := json.NewDecoder(r).Decode(ret); err != nil {
		return nil, fmt.Errorf("common response: %v", err)
	}

	return ret, nil
}

func (c *CommonResponse) GetPayload(payload interface{}) error {
	if c.Status != http.StatusOK {
		return fmt.Errorf("common response: response with error [%v]: %v", c.Status, c.Message)
	}
	if err := mapstructure.Decode(c.Payload, payload); err != nil {
		return fmt.Errorf("common response: payload with error: %v", err)
	}

	return nil
}

type PublicRequestHeader struct {
	Token    string `json:"access-token"`
	BizSeqNo string `json:"bizSeqNo"`
}

type PublicRequest struct {
	ApiName string      `json:"apiName"`
	Payload interface{} `json:"reqData"`
}

type UserRequest struct {
	SignUserId  string `json:"signUserId"`
	Description string `json:"description"`
}

type UserResponse struct {
	SignUserId  string `json:"signUserId"`
	AppId       string `json:"appId"`
	PrivateKey  string `json:"privateKey"`
	Address     string `json:"address"`
	PublicKey   string `json:"publicKey"`
	Description string `json:"description"`
	EncryptType int    `json:"encryptType"`
}

type QueryTransactionRequest struct {
	Name    string `json:"funcName"`
	Encoded string `json:"encodeStr"`
	Abi     []byte `json:"functionAbi"`
	Address string `json:"contractAddress"`
}

type GetTransactionResponse struct {
	BlockHash        string    `json:"blockHash"`
	BlockNumber      int64     `json:"blockNumber"`
	From             string    `json:"from"`
	Gas              string    `json:"gas"`
	Hash             string    `json:"hash"`
	Input            string    `json:"input"`
	Nonce            string    `json:"nonce"`
	To               string    `json:"to"`
	TransactionIndex string    `json:"transactionIndex"`
	Value            string    `json:"value"`
	GasPrice         string    `json:"gasPrice"`
	BlockLimit       string    `json:"blockLimit"`
	ChainId          string    `json:"chainId"`
	Extra            string    `json:"extraData"`
	Signature        Signature `json:"signature"`
}

type Signature struct {
	R         string `json:"r"`
	S         int64  `json:"s"`
	V         string `json:"v"`
	Signature string `json:"signature"`
}

type TransactionRecipetResponse struct {
	Constant        bool   `json:"constant"`
	QueryInfo       string `json:"queryInfo"`
	TransactionHash string `json:"transactionHash"`
	BlockHash       string `json:"blockHash"`
	BlockNumber     int64  `json:"blockNumber"`
	GasUsed         int64  `json:"gasUsed"`
	Status          string `json:"status"`
	From            string `json:"from"`
	To              string `json:"to"`
	Input           string `json:"input"`
	Output          string `json:"output"`
	ContractAddress string `json:"contractAddress"`
	LogsBloom       string `json:"logsBloom"`
	Logs            []Log  `json:"logs"`
}

type Log struct {
	Address     string `json:"address"`
	Topics      []byte `json:"topics"`
	Data        string `json:"data"`
	BlockNumber int64  `json:"blockNumber"`
}

type SignedTransactionRequest struct {
	Signed string `json:"signedStr"`
	Sync   bool   `json:"bool"`
}

type SignedTransactionResponse struct {
}

type SignRequest struct {
	SignUserId  string `json:"signUserId"`
	MessageHash string `json:"messageHash"`
}

type SignResponse struct {
	Signature string `json:"signDataStr"`
}

type DeployRequest struct {
	Params      [][]byte `json:"funcParam"`
	Abi         []byte   `json:"contractAbi"`
	SignUserId  string   `json:"signUserId"`
	BytecodeBin string   `json:"bytecodeBin"`
}

type DeployResponse struct {
	Constant        bool   `json:"constant"`
	QueryInfo       string `json:"queryInfo"`
	TransactionHash string `json:"transactionHash"`
	BlockHash       string `json:"blockHash"`
	BlockNumber     int64  `json:"blockNumber"`
	GasUsed         int64  `json:"gasUsed"`
	Status          string `json:"status"`
	From            string `json:"from"`
	To              string `json:"to"`
	Input           string `json:"input"`
	Output          string `json:"output"`
	ContractAddress string `json:"contractAddress"`
	Message         string `json:"message"`
}

type SendTransactionRequest struct {
	Name            string   `json:"funcName"`
	Params          [][]byte `json:"funcParam"`
	Abi             []byte   `json:"functionAbi"`
	SignUserId      string   `json:"signUserId"`
	ContractAddress string   `json:"contractAddress"`
}

type SendTransactionResponse DeployResponse
