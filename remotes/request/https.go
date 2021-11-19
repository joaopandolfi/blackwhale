package request

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/xerrors"
)

var transport *http.Transport
var defaultTimeout time.Duration = time.Minute * 10

func getTransport() *http.Transport {
	if transport == nil {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		http.DefaultClient.Timeout = defaultTimeout
	}
	return transport
}

// SetDefaultTimeout to requests
func SetDefaultTimeout(timeout time.Duration) {
	defaultTimeout = timeout
}

// Get - basic call a get command
func Get(url string) (body []byte) {
	client := &http.Client{Transport: getTransport()}
	resp, err := client.Get(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		resp.Body.Close()
		if r := recover(); r != nil {
			fmt.Println("Error on", r)
		}
	}()

	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	return
}

func RequestWithHeader(method, url string, head map[string]string, data []byte) ([]byte, int, error) {
	client := &http.Client{Transport: getTransport()}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, 0, xerrors.Errorf("making requester: %w", err)
	}

	//Setting Headers
	for k, v := range head {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		statusCode := http.StatusBadRequest
		if resp != nil {
			statusCode = resp.StatusCode
		}
		return nil, statusCode, xerrors.Errorf("[RequestWithHeader] - Error on make %s request, URL: %s, DATA: %s , ERROR: %w", method, url, string(data), err)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return b, resp.StatusCode, xerrors.Errorf("[RequestWithHeader] - Error on Read Body result, URL: %s, DATA: %s , ERROR: %w", url, string(data), err)
	}

	return b, resp.StatusCode, err
}

// PostWithHeder2 - make post and aggregate statuscode response
func PostWithHeader2(url string, head map[string]string, data []byte) ([]byte, int, error) {
	return RequestWithHeader("POST", url, head, data)
}

// PostWithHeader -
func PostWithHeader(url string, head map[string]string, data []byte) (body []byte, err error) {

	body, statuscode, err := PostWithHeader2(url, head, data)

	if statuscode == 400 {
		err = errors.New("[PostWithHeader] - Got Message error 400")
	}

	return
}

// Post - simple post
func Post(url string, data []byte) (body []byte, err error) {

	client := &http.Client{Transport: getTransport()}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, xerrors.Errorf("making requester: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("[Post] - Error on make POST request, URL: %s, DATA: %s , ERROR: %s", url, string(data), err.Error()))
		return
	}

	defer func() {
		resp.Body.Close()
		if r := recover(); r != nil {
			fmt.Println("Error on HTTP POST", r)
		}
	}()

	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		err = errors.New(fmt.Sprintf("[Post] - Error on Read Body result, URL: %s, DATA: %s , ERROR: %s", url, string(data), err.Error()))
	}

	if resp.StatusCode == 400 {
		err = errors.New("[Post] - Got Message error 400")
	}

	return
}

// GetWithHeader -
func GetWithHeader(url string, head map[string]string) (body []byte, err error) {
	client := &http.Client{Transport: getTransport()}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, xerrors.Errorf("making requester: %w", err)
	}

	//Setting Headers
	for k, v := range head {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("[GetWithHeader] - Error on make GET request, URL: %s , ERROR: %s", url, err.Error()))
		return
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		err = errors.New(fmt.Sprintf("[GetWithHeader] - Error on Read Body result, URL: %s, ERROR: %s", url, err.Error()))
	}

	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("[GetWithHeader] - Got Message error %d", resp.StatusCode))
	}

	return
}
