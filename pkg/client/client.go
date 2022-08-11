package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Axway/agents-layer7/pkg/models/policy"
	"github.com/Axway/agents-layer7/pkg/models/policyversion"
	"github.com/Axway/agents-layer7/pkg/models/service"
	xj "github.com/basgys/goxml2json"
	"github.com/sirupsen/logrus"
)

// GatewayClient -
type GatewayClient struct {
	client HTTPClient
	host   string
	user   string
	pass   string
}

// HTTPClient an interface for the go HTTP Client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewGatewayClient creates a client to connect to the layer7 gateway
func NewGatewayClient(host, user, pass string) *GatewayClient {
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 30 * time.Second,
	}

	return &GatewayClient{
		client: c,
		host:   host,
		user:   user,
		pass:   pass,
	}
}

// GetServices gets a list of services from layer7
func (c *GatewayClient) GetServices() (*service.ListServices, error) {
	url := fmt.Sprintf("%s/services", c.host)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.user, c.pass)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		logFailed(res)
		return nil, fmt.Errorf("expected a 200 response but received %d", res.StatusCode)
	}

	buf, err := xj.Convert(res.Body)
	if err != nil {
		return nil, err
	}

	list := &service.ListServices{}
	err = json.Unmarshal(buf.Bytes(), list)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return list, nil
}

// GetActivePolicy gets the active policy for a service
func (c *GatewayClient) GetActivePolicy(apiID string) (*policy.PolicyItem, error) {
	url := fmt.Sprintf("%s/services/%s/versions/active", c.host, apiID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.user, c.pass)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		logFailed(res)
		return nil, fmt.Errorf("expected a 200 response but received %d", res.StatusCode)
	}

	buf, err := xj.Convert(res.Body)
	if err != nil {
		return nil, err
	}

	pv := &policyversion.PolicyVersionRes{}
	err = json.Unmarshal(buf.Bytes(), pv)
	if err != nil {
		return nil, err
	}

	r := strings.NewReader(pv.Item.Resource.PolicyVersion.XML)
	buf, err = xj.Convert(r)
	if err != nil {
		return nil, err
	}

	p := &policy.PolicyItem{}
	err = json.Unmarshal(buf.Bytes(), p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// GetSpec gets a spec from the provided host
func (c *GatewayClient) GetSpec(host string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, host, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		logFailed(res)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		logFailed(res)
		return nil, fmt.Errorf("expected a 200 response but received %d", res.StatusCode)
	}

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bts, nil
}

func logFailed(res *http.Response) {
	if res == nil || res.Body == nil {
		return
	}
	b, _ := ioutil.ReadAll(res.Body)
	logrus.Errorf("request failed: %s", string(b))
}
