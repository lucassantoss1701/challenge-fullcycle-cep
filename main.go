package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Response struct {
	ApiName    string
	StatusCode int
	Body       string
}

func main() {
	cep := "13167-616"
	addressChan := make(chan *Response, 2)

	go func() {
		address, err := getApiCep(cep)
		if err != nil {
			fmt.Println(err)
		} else {
			addressChan <- address
		}
	}()

	go func() {
		address, err := getViaCep(cep)
		if err != nil {
			fmt.Println(err)
		} else {
			addressChan <- address
		}
	}()

	select {
	case response := <-addressChan:
		fmt.Printf("api-name: %s \nstatus code: %d \nbody: %s \n ", response.ApiName, response.StatusCode, response.Body)
	case <-time.After(time.Second):
		fmt.Printf("TIMEOUT \n")
	}

}

func getApiCep(cep string) (*Response, error) {
	apiName := "api-cep"
	url := fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", cep)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return &Response{ApiName: apiName, StatusCode: res.StatusCode, Body: string(body)}, nil
}

func getViaCep(cep string) (*Response, error) {
	apiName := "via-cep"
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return &Response{ApiName: apiName, StatusCode: res.StatusCode, Body: string(body)}, nil
}
