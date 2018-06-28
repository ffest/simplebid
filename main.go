package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

type row struct {
	Price  int    `json:"price"`
	Source string `json:"source"`
}

func main() {
	listenAddr := flag.String("addr", ":8080", "http listen address")
	flag.Parse()

	http.HandleFunc("/winner", handler())

	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		fmt.Errorf("%+v", err)
	}
}

func handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		urls := r.URL.Query()["s"]

		if len(urls) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Не совсем понятно из условия, 100ms на общий запрос или на внешние апи
		// Если на внешние, то ставим 100ms в контекст или делаем клиент с таким таймаутом
		ctx, cancel := context.WithTimeout(context.Background(), 90*time.Millisecond)
		defer cancel()
		dataChan := make(chan []row, len(urls))

		bids := make([]row, 0)

		go func() {
			select {
			case re := <-dataChan:
				bids = append(bids, re...)
			case <-ctx.Done():
				fmt.Println(ctx.Err())
			}
		}()

		wg := sync.WaitGroup{}
		for _, url := range urls {
			wg.Add(1)
			go worker(ctx, url, dataChan, &wg)
		}
		wg.Wait()

		if len(bids) < 2 {
			// DO SOMETHING/ 204 maybe
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(sortAndGetResult(bids))
	}
}

func worker(ctx context.Context, url string, dc chan []row, wg *sync.WaitGroup) {
	defer wg.Done()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp == nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	result := make([]row, 0)
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return
	}

	for i := range result {
		result[i].Source = url
	}

	dc <- result
}

func sortAndGetResult(bids []row) row {
	sort.Slice(bids, func(i, j int) bool {
		if bids[i].Price < bids[j].Price {
			return true
		}
		return false
	})

	return bids[len(bids)-2]
}
