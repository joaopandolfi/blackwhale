package request

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"

	"fmt"

	"github.com/pkg/errors"
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

	payloadData := bytes.NewBuffer(data)

	if head["Content-Encoding"] == "gzip" {
		var compressedData bytes.Buffer
		gzipBuff := gzip.NewWriter(&compressedData)
		if _, err := gzipBuff.Write(data); err != nil {
			return nil, http.StatusExpectationFailed, fmt.Errorf("gzipping body: %w", err)
		}
		gzipBuff.Close()
		payloadData = &compressedData
	}

	req, err := http.NewRequest(method, url, payloadData)
	if err != nil {
		return nil, 0, fmt.Errorf("making requester: %w", err)
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
		return nil, statusCode, fmt.Errorf("[RequestWithHeader] - Error on make %s request, URL: %s, DATA: %s , ERROR: %w", method, url, string(data), err)
	}

	defer resp.Body.Close()

	var b []byte

	if resp.Header.Get("Content-Encoding") != "gzip" {
		b, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return b, resp.StatusCode, fmt.Errorf("[RequestWithHeader] - Error on Read Body result, URL: %s, DATA: %s , ERROR: %w", url, string(data), err)
		}

	} else {
		r, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, http.StatusExpectationFailed, fmt.Errorf("reading gzip body: %w", err)
		}
		var resB bytes.Buffer
		_, err = resB.ReadFrom(r)
		if err != nil {
			return nil, http.StatusExpectationFailed, fmt.Errorf("reading gzip bytes: %w", err)
		}
		r.Close()
		b = resB.Bytes()
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
		return nil, fmt.Errorf("making requester: %w", err)
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
		return nil, fmt.Errorf("making requester: %w", err)
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
