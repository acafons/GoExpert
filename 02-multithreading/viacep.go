package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type viaCepDTO struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func GetAddressViaCEP(cep string, data chan Message) {
	var msg Message

	req, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		log.Println("[error] GET ViaCEP failed: ", err)
		return
	}

	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("[error] Failed to read ViaCEP response:", err)
		return
	}

	var viacep viaCepDTO
	err = json.Unmarshal(res, &viacep)
	if err != nil {
		log.Println("[error] Failed to parse ViaCEP response:", err)
		return
	}

	msg.Address = viacep.Logradouro + ", " + viacep.Bairro + ", " + viacep.Localidade + " - " +
		viacep.Uf + ", " + viacep.Cep

	msg.API = "ViaCep"

	data <- msg
}
