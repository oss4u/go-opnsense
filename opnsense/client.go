package opnsense

import (
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
	"sync"
)

<<<<<<< Updated upstream
var (
	createLock  = &sync.Mutex{}
	modifyLock  = &sync.Mutex{}
	opnSenseApi *OpnSenseApi
)

func GetOpnSenseClient(host, key, secret string) *OpnSenseApi {
=======
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
>>>>>>> Stashed changes
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

func (c *OpnSenseApi) getClient() *resty.Client {
	if c.client == nil {
		c.client = resty.New()
		c.client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	return c.client
}

<<<<<<< Updated upstream
func (c *OpnSenseApi) ModifyingRequest(module, controller, command, data string, params []string) (string, error) {
	client := c.getClient()
	modifyLock.Lock()
	defer modifyLock.Unlock()

	url := fmt.Sprintf("%s/api/%s/%s/%s", c.address, module, controller, command)
	if len(params) > 0 {
		url += "/" + fmt.Sprintf("%s", params)
	}

=======
func (c *OpnSenseApi) ModifyingRequest(module string, controller string, command string, data string, params []string) (string, error) {
	client := c.get_client()
	requestLock.Lock()
	defer requestLock.Unlock()
>>>>>>> Stashed changes
	request := client.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(c.key, c.secret).
		SetBody(data)

<<<<<<< Updated upstream
	res, err := request.Post(url)
=======
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
	res, err := request.
		Post(url)
>>>>>>> Stashed changes
	if err != nil {
		return "", err
	}
	return res.String(), nil
}

<<<<<<< Updated upstream
func (c *OpnSenseApi) NonModifyingRequest(module, controller, command string, params []string) (string, int, error) {
	client := c.getClient()

	url := fmt.Sprintf("%s/api/%s/%s/%s", c.address, module, controller, command)
	if len(params) > 0 {
		url += "/" + fmt.Sprintf("%s", params)
	}

	request := client.R().
		SetBasicAuth(c.key, c.secret)

=======
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
>>>>>>> Stashed changes
	res, err := request.Get(url)
	if err != nil {
		return "", res.StatusCode(), err
	}
	return res.String(), res.StatusCode(), nil
}
