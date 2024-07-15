package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type brasilApiDTO struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func GetAddressBrasilApi(cep string, data chan Message) {
	var msg Message

	req, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if err != nil {
		log.Println("[error] GET BrasilAPI failed: ", err)
		return
	}

	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("[error] Failed to read BrasilAPI response:", err)
		return
	}

	var brasilApi brasilApiDTO
	err = json.Unmarshal(res, &brasilApi)
	if err != nil {
		log.Println("[error] Failed to parse BrasilAPI response:", err)
		return
	}

	msg.Address = brasilApi.Street + ", " + brasilApi.Neighborhood + ", " + brasilApi.City +
		" - " + brasilApi.State + ", " + brasilApi.Cep

	msg.API = "BrasilAPI"

	data <- msg
}
