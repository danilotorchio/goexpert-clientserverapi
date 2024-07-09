package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/danilotorchio/goexpert-clientserverapi/database"
	"github.com/danilotorchio/goexpert-clientserverapi/models"
)

const (
	RequestTimeout  = 200 * time.Millisecond
	DatabaseTimeout = 10 * time.Millisecond

	ExchangeAPI = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
)

func ExchangeHandler(w http.ResponseWriter, r *http.Request) {
	reqCtx, reqCancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer reqCancel()

	resp, err := fetchUsdExchange(reqCtx)
	if err != nil {
		handleError(w, err, "Error fetching exchange rate")
		return
	}

	dbCtx, dbCancel := context.WithTimeout(context.Background(), DatabaseTimeout)
	defer dbCancel()

	if err := database.InsertNewExchangeRate(dbCtx, resp.Bid); err != nil {
		handleError(w, err, "Error inserting exchange rate")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResp := map[string]string{"cotacao": resp.Bid}
	if err := json.NewEncoder(w).Encode(jsonResp); err != nil {
		handleError(w, err, "Error encoding response")
		return
	}
}

func HistoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), DatabaseTimeout)
	defer cancel()

	limit := 10

	if param := r.URL.Query().Get("limit"); param != "" {
		l, err := strconv.Atoi(param)
		if err == nil {
			limit = l
		}
	}

	resp, err := database.GetExchangeHistory(ctx, limit)
	if err != nil {
		handleError(w, err, "Error fetching exchange history")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		handleError(w, err, "Error encoding response")
		return
	}
}

func fetchUsdExchange(ctx context.Context) (*models.ExchangeApiResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ExchangeAPI, nil)
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

	var data models.ExchangeApiResponse
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func handleError(w http.ResponseWriter, err error, msg string) {
	log.Default().Println(msg, "-", err)

	if errors.Is(err, context.DeadlineExceeded) {
		http.Error(w, "Timeout excedeed", http.StatusRequestTimeout)
		return
	}

	http.Error(w, "Internal server error", http.StatusInternalServerError)
}
