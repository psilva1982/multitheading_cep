package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type AddressViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Unidade     string `json:"unidade"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Source      string // Source information
}

type AddressBrasilApi struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"street"`
	Bairro      string `json:"neighborhood"`
	Localidade  string `json:"city"`
	Uf          string `json:"state"`
	Source      string // Source information
}

func main() {
	// Parse command line arguments
	var cep string
	flag.StringVar(&cep, "cep", "", "CEP to lookup")
	flag.Parse()

	if cep == "" {
		fmt.Println("You must provide a CEP to lookup using the -cep flag")
		os.Exit(1)
	}

	// Channel to receive results from goroutines
	resultCh := make(chan interface{}, 2)

	// Perform concurrent requests to both APIs
	go fetchFromAPI1(cep, resultCh)
	go fetchFromAPI2(cep, resultCh)

	// Wait for results from both goroutines
	timeout := time.After(1 * time.Second)
	var result interface{}

	for i := 0; i < 2; i++ {
		select {
		case res := <-resultCh:
			// If we haven't set a result yet or this one came in faster
			if result == nil || res != nil {
				result = res
			}
		case <-timeout:
			fmt.Println("Timeout exceeded while waiting for responses")
			os.Exit(1)
		}
	}

	// Display the result
	switch addr := result.(type) {
	case *AddressViaCep:
		displayAddress("ViaCep", addr)
	case *AddressBrasilApi:
		displayAddress("BrasilApi", addr)
	default:
		fmt.Println("No address found for CEP:", cep)
	}
}

func fetchFromAPI1(cep string, ch chan<- interface{}) {
	url := "https://brasilapi.com.br/api/cep/v1/" + cep
	source := "https://brasilapi.com.br"

	address, err := fetchAddressBrasilApi(url)
	if err != nil {
		ch <- nil
		return
	}
	address.Source = source
	ch <- address
}

func fetchFromAPI2(cep string, ch chan<- interface{}) {
	url := "http://viacep.com.br/ws/" + cep + "/json/"
	source := "http://viacep.com.br"

	address, err := fetchAddressViaCep(url)
	if err != nil {
		ch <- nil
		return
	}
	address.Source = source
	ch <- address
}

func fetchAddressBrasilApi(url string) (*AddressBrasilApi, error) {
	client := http.Client{
		Timeout: 1 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	address := &AddressBrasilApi{}
	err = json.Unmarshal(body, address)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func fetchAddressViaCep(url string) (*AddressViaCep, error) {
	client := http.Client{
		Timeout: 1 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	address := &AddressViaCep{}
	err = json.Unmarshal(body, address)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func displayAddress(source string, addr interface{}) {
	switch a := addr.(type) {
	case *AddressViaCep:
		fmt.Println("Address found using:", source)
		fmt.Println("CEP:", a.Cep)
		fmt.Println("Logradouro:", a.Logradouro)
		fmt.Println("Complemento:", a.Complemento)
		fmt.Println("Bairro:", a.Bairro)
		fmt.Println("Localidade:", a.Localidade)
		fmt.Println("UF:", a.Uf)
		fmt.Println("Unidade:", a.Unidade)
		fmt.Println("IBGE:", a.Ibge)
		fmt.Println("GIA:", a.Gia)
	case *AddressBrasilApi:
		fmt.Println("Address found using:", source)
		fmt.Println("CEP:", a.Cep)
		fmt.Println("Logradouro:", a.Logradouro)
		fmt.Println("Bairro:", a.Bairro)
		fmt.Println("Localidade:", a.Localidade)
		fmt.Println("UF:", a.Uf)
	default:
		fmt.Println("No address found")
	}
}