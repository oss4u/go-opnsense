package opnsense

import (
	"crypto/tls"
	"fmt"
	"os"
	"sync"

	"github.com/go-resty/resty/v2"
)

var createLock = &sync.Mutex{}
var requestLock = &sync.Mutex{}
var opnSenseApi *OpnSenseApi

func NewOpnSenseClient(host string, key string, secret string) *OpnSenseApi {
	if host == "" || key == "" || secret == "" {
		host = os.Getenv("OPNSENSE_ADDRESS")
		key = os.Getenv("OPNSENSE_KEY")
		secret = os.Getenv("OPNSENSE_SECRET")
	}

	return &OpnSenseApi{
		address: host,
		key:     key,
		secret:  secret,
	}
}

func GetOpnSenseClient(host string, key string, secret string) *OpnSenseApi {
	if opnSenseApi == nil {
		createLock.Lock()
		defer createLock.Unlock()
		opnSenseApi = NewOpnSenseClient(host, key, secret)
	}
	return opnSenseApi
}

type OpnSenseApi struct {
	client  *resty.Client
	address string
	key     string
	secret  string
}

func (c *OpnSenseApi) get_client() *resty.Client {
	if c.client != nil {
		return c.client
	}
	c.client = resty.New()
	c.client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	return c.client
}

func (c *OpnSenseApi) ModifyingRequest(module string, controller string, command string, data string, params []string) (string, error) {
	client := c.get_client()
	requestLock.Lock()
	defer requestLock.Unlock()
	request := client.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(c.key, c.secret)
	url := ""

	if len(params) > 0 {
		p := ""
		for _, v := range params {
			p = fmt.Sprintf("%s/%s", p, v)
			url = fmt.Sprintf("%s/api/%s/%s/%s%s", c.address, module, controller, command, p)
		}
	} else {
		url = fmt.Sprintf("%s/api/%s/%s/%s", c.address, module, controller, command)
	}

	if len(data) > 0 {
		request = request.SetBody(data)
	} else {
		request = request.SetBody(`{}`)
	}
	res, err := request.Post(url)
	if err != nil {
		return res.String(), err
	}

	return res.String(), nil
}

func (c *OpnSenseApi) NonModifyingRequest(module string, controller string, command string, params []string) (string, int, error) {
	client := c.get_client()
	requestLock.Lock()
	defer requestLock.Unlock()
	request := client.R().
		SetBasicAuth(c.key, c.secret)
	url := ""
	if len(params) > 0 {
		p := ""
		for _, v := range params {
			p = fmt.Sprintf("%s/%s", p, v)
			url = fmt.Sprintf("%s/api/%s/%s/%s%s", c.address, module, controller, command, p)
		}
	} else {
		url = fmt.Sprintf("%s/api/%s/%s/%s", c.address, module, controller, command)
	}
	res, err := request.Get(url)
	if err != nil {
		return res.String(), res.StatusCode(), err
	}
	return res.String(), res.StatusCode(), nil
}
