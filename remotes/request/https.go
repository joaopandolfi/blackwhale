package request

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func Get(url string) (body []byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error on", r)
		}
	}()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	return
}

func PostWithHeader(url string, head map[string]string, data []byte) (body []byte, err error) {
	/*defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error on", r)
		}
	}()*/

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))

	//Setting Headers
	for k, v := range head {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("[PostWithHeader] - Error on make POST request, URL: %s, DATA: %s , ERROR: %s", url, string(data), err.Error()))
		return
	}

	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		err = errors.New(fmt.Sprintf("[PostWithHeader] - Error on Read Body result, URL: %s, DATA: %s , ERROR: %s", url, string(data), err.Error()))
	}

	if resp.StatusCode == 400 {
		err = errors.New("[PostWithHeader] - Got Message error 400")
	}

	return
}

func Post(url string, data []byte) (body []byte, err error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error on HTTP POST", r)
		}
	}()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("[Post] - Error on make POST request, URL: %s, DATA: %s , ERROR: %s", url, string(data), err.Error()))
		return
	}

	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		err = errors.New(fmt.Sprintf("[Post] - Error on Read Body result, URL: %s, DATA: %s , ERROR: %s", url, string(data), err.Error()))
	}

	if resp.StatusCode == 400 {
		err = errors.New("[Post] - Got Message error 400")
	}

	return
}

func GetWithHeader(url string, head map[string]string) (body []byte, err error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, _ := http.NewRequest("GET", url, nil)

	//Setting Headers
	for k, v := range head {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("[GetWithHeader] - Error on make GET request, URL: %s , ERROR: %s", url, err.Error()))
		return
	}

	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		err = errors.New(fmt.Sprintf("[GetWithHeader] - Error on Read Body result, URL: %s, ERROR: %s", url, err.Error()))
	}

	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("[GetWithHeader] - Got Message error %d", resp.StatusCode))
	}

	return
}
