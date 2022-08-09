package client

import (
	"fmt"
	"strings"
	"testing"

	xj "github.com/basgys/goxml2json"
	"github.com/stretchr/testify/assert"
)

func TestClientGetServices(t *testing.T) {
	host := "https://ec2-44-205-180-158.compute-1.amazonaws.com:9443/restman1/1.0"
	client := NewGatewayClient(host, "admin", "L7Secure$0@")
	svcs, err := client.GetServices()
	// svcs, err := client.GetActivePolicy("55b5c175f0b7924dcb53865072a44014")
	assert.Nil(t, err)
	assert.NotNil(t, svcs)
}

func TestXMLConversion(t *testing.T) {
	t.Skip()
	reader := strings.NewReader(s)
	buf, _ := xj.Convert(reader)

	fmt.Println(buf)
}

const s = ""
