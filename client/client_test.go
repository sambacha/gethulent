package client

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeClient struct {
}

func setResultValue(ptr interface{}, v string) error {
	val := reflect.ValueOf(ptr)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer type")
	}
	m := map[string]string{"result": v}
	bytes, _ := json.Marshal(m)
	return json.Unmarshal(bytes, ptr)
}

func getResultValue(result interface{}) string {
	m := result.(map[string]interface{})
	return m["result"].(string)
}

func (fc *fakeClient) Call(result interface{}, method string, args ...interface{}) error {
	retStr := method
	for _, arg := range args {
		retStr += fmt.Sprintf("%v", arg)
	}

	return setResultValue(&result, retStr)
}

func (fc *fakeClient) Close() {
	return
}

func newFakeClientAgent() *GethAgent {
	return &GethAgent{client: &fakeClient{}}
}

func TestCallMethod(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		method string
		args   []string
	}{
		{
			"eth_getBlockByNumber",
			[]string{"latest", "true"},
		},
		{
			"eth_getBlockByNumber",
			[]string{"0x123456", "true"},
		},
		{
			"eth_getBalance",
			[]string{"0x12345678", "latest"},
		},
		{
			"eth_blockNumber",
			[]string{},
		},
		{
			"eth_getTransactionCount",
			[]string{"0x12345678", "latest"},
		},
	}
	agent := newFakeClientAgent()

	for _, tt := range tests {
		var result interface{}
		err := agent.CallMethod(&result, tt.method, tt.args)
		assert.NoError(err)
		assert.Equal(getResultValue(result), func(method string, args []string) string {
			strs := []string{method}
			strs = append(strs, args...)
			return strings.Join(strs, "")
		}(tt.method, tt.args))
	}

}
