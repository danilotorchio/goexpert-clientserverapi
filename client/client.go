package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Cotacao string `json:"cotacao"`
}

const (
	RequestTimeout = 300 * time.Millisecond

	ServerUri = "http://localhost:8080/cotacao"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ServerUri, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		log.Fatalf("Server response error: %s", string(data))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var cotacao Cotacao
	if err := json.Unmarshal(data, &cotacao); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("cotacao.txt", []byte(fmt.Sprintf("Dólar: %s", cotacao.Cotacao)), os.ModePerm); err != nil {
		log.Fatal(err)
	}

	log.Default().Println("Request success! Check the cotacao.txt file to see the result.")
}
