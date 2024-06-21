package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ExchangeRate struct {
	Bid string `json:"bid"`
}

const ServerUrl = "http://localhost:8080/cotacao"
const OutputFile = "cotacao.txt"
const HttpRequestTimeout = 300 * time.Millisecond

func main() {
	exchRate, err := GetCurrencyExchangeRate()
	if err != nil {
		log.Println("[error]", err)
		return
	}

	err = PersistCurrencyExchangeInfo(exchRate)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	log.Println("[info] USD to BRL: " + exchRate.Bid)
}

func GetCurrencyExchangeRate() (*ExchangeRate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), HttpRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", ServerUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Get \"" + ServerUrl + "\": " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var exchRate ExchangeRate

	err = json.Unmarshal(body, &exchRate)
	if err != nil {
		return nil, err
	}

	return &exchRate, nil
}

func PersistCurrencyExchangeInfo(exchRate *ExchangeRate) error {
	file, err := os.Create(OutputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	size, err := fmt.Fprintf(file, "Dólar:{%s}", string(exchRate.Bid))
	if err != nil {
		return err
	}

	log.Println("[info] File " + OutputFile + " saved with success! Size : " + strconv.Itoa(size) + " bytes")

	return nil
}
