package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type CurrencyExchangeInfo struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type ExchangeRate struct {
	Bid string `json:"bid"`
}

const HttpRequestTimeout = 200 * time.Millisecond
const DataPersistenceTimeout = 10 * time.Millisecond

func main() {
	InitDB()

	mux := http.NewServeMux()

	mux.HandleFunc("/cotacao", DollarExchangeRateHandler)

	log.Println("Listening on :3000...")
	http.ListenAndServe(":8080", mux)
}

func DollarExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	exchInfo, err := GetCurrencyExchangeInfo()
	if err != nil {
		log.Println("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = PersistCurrencyExchangeInfo(exchInfo)
	if err != nil {
		log.Println("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var exchRate ExchangeRate
	exchRate.Bid = exchInfo.Usdbrl.Bid

	json.NewEncoder(w).Encode(exchRate)
}

func GetCurrencyExchangeInfo() (*CurrencyExchangeInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), HttpRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var exchInfo CurrencyExchangeInfo

	err = json.Unmarshal(body, &exchInfo)
	if err != nil {
		return nil, err
	}

	log.Println(exchInfo)

	return &exchInfo, nil
}

func InitDB() {
	// Open an SQLite database (creates a new file if not exists).
	db, err := sql.Open("sqlite3", "./exchange.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create a table (if not exists).
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS rates (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            code TEXT,
			code_in TEXT,
			name TEXT,
			high TEXT,
			low TEXT,
			var_bid TEXT,
			pct_change TEXT,
			bid REAL TEXT,
			ask REAL TEXT,
			timestamp TEXT,
			create_date TEXT
        )
    `)
	if err != nil {
		panic(err)
	}
}

func PersistCurrencyExchangeInfo(exchInfo *CurrencyExchangeInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), DataPersistenceTimeout)
	defer cancel()

	db, err := sql.Open("sqlite3", "./exchange.db")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO rates(code, code_in, name, high, low, var_bid, pct_change, " +
		"bid, ask, timestamp, create_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		exchInfo.Usdbrl.Code,
		exchInfo.Usdbrl.Codein,
		exchInfo.Usdbrl.Name,
		exchInfo.Usdbrl.High,
		exchInfo.Usdbrl.Low,
		exchInfo.Usdbrl.VarBid,
		exchInfo.Usdbrl.PctChange,
		exchInfo.Usdbrl.Bid,
		exchInfo.Usdbrl.Ask,
		exchInfo.Usdbrl.Timestamp,
		exchInfo.Usdbrl.CreateDate)
	if err != nil {
		return err
	}

	return nil
}
