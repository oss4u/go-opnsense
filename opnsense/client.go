package opnsense

import (
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
	"sync"
)

var (
	createLock  = &sync.Mutex{}
	modifyLock  = &sync.Mutex{}
	opnSenseApi *OpnSenseApi
)

func GetOpnSenseClient(host, key, secret string) *OpnSenseApi {
	if opnSenseApi == nil {
		createLock.Lock()
		defer createLock.Unlock()
		if host == "" || key == "" || secret == "" {
			host = os.Getenv("OPNSENSE_ADDRESS")
			key = os.Getenv("OPNSENSE_KEY")
			secret = os.Getenv("OPNSENSE_SECRET")
		}
		opnSenseApi = &OpnSenseApi{
			address: host,
			key:     key,
			secret:  secret,
		}
	}
	return opnSenseApi
}

type OpnSenseApi struct {
	client  *resty.Client
	address string
	key     string
	secret  string
}

func (c *OpnSenseApi) getClient() *resty.Client {
	if c.client == nil {
		c.client = resty.New()
		c.client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	return c.client
}

func (c *OpnSenseApi) ModifyingRequest(module, controller, command, data string, params []string) (string, error) {
	client := c.getClient()
	modifyLock.Lock()
	defer modifyLock.Unlock()

	url := fmt.Sprintf("%s/api/%s/%s/%s", c.address, module, controller, command)
	if len(params) > 0 {
		url += "/" + fmt.Sprintf("%s", params)
	}

	request := client.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(c.key, c.secret).
		SetBody(data)

	res, err := request.Post(url)
	if err != nil {
		return "", err
	}
	return res.String(), nil
}

func (c *OpnSenseApi) NonModifyingRequest(module, controller, command string, params []string) (string, int, error) {
	client := c.getClient()

	url := fmt.Sprintf("%s/api/%s/%s/%s", c.address, module, controller, command)
	if len(params) > 0 {
		url += "/" + fmt.Sprintf("%s", params)
	}

	request := client.R().
		SetBasicAuth(c.key, c.secret)

	res, err := request.Get(url)
	if err != nil {
		return "", res.StatusCode(), err
	}
	return res.String(), res.StatusCode(), nil
}
