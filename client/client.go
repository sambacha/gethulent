package client

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/rpc"
)

type rpcClient interface {
	Call(result interface{}, method string, args ...interface{}) error
	Close()
}

//Agent defines a RPC agent interface
type Agent interface {
	CallMethod(result interface{}, method string, params []string) error
	Close()
}

//GethAgent wraps ethereum rpc client
type GethAgent struct {
	client rpcClient
}

//New returns an Ethereum rpc client
func New(url string) (Agent, error) {
	client, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	return &GethAgent{client: client}, nil
}

func getMethodArgs(args []string) []interface{} {
	ret := make([]interface{}, len(args))
	for i := range args {
		ret[i] = args[i]
	}
	return ret
}

func inDefaultBlockNum(s string) bool {
	for _, bn := range []string{"earliest", "latest", "pending"} {
		if bn == s {
			return true
		}
	}
	return false
}

//Close closes RPC client
func (ga *GethAgent) Close() {
	ga.client.Close()
}

func (ga *GethAgent) ethGetBalance(result interface{}, params []string) error {
	return ga.client.Call(result, "eth_getBalance", params[0], params[1])
}

func (ga *GethAgent) ethGetBlockByNumber(result interface{}, params []string) error {
	if len(params) < 2 {
		return fmt.Errorf("need to provide block number(0x...) and transaction flag(bool)")
	}
	blockNum := params[0]
	if !strings.HasPrefix(blockNum, "0x") && !inDefaultBlockNum(blockNum) {
		num, err := strconv.ParseInt(blockNum, 10, 64)
		if err != nil {
			return err
		}
		blockNum = "0x" + strconv.FormatInt(num, 16)
	}
	fullTx, err := strconv.ParseBool(params[1])
	if err != nil {
		return err
	}
	return ga.client.Call(result, "eth_getBlockByNumber", blockNum, fullTx)
}

//CallMethod calls RPC method with params
func (ga *GethAgent) CallMethod(result interface{}, method string, params []string) error {
	var err error
	switch method {
	case "eth_getBalance":
		err = ga.ethGetBalance(result, params)
	case "eth_getBlockByNumber":
		err = ga.ethGetBlockByNumber(result, params)
	default:
		err = ga.client.Call(result, method, getMethodArgs(params)...)
	}
	return err
}
