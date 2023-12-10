// main.go
package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"
)

type RequestPayload struct {
	ToSort [][]int `json:"to_sort"`
}

type ResponsePayload struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNS       int64   `json:"time_ns"`
}

func sortSequential(toSort [][]int) [][]int {
	for i := range toSort {
		sort.Ints(toSort[i])
	}
	return toSort
}

func sortConcurrent(toSort [][]int) [][]int {
	var wg sync.WaitGroup
	for i := range toSort {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sort.Ints(toSort[i])
		}(i)
	}
	wg.Wait()
	return toSort
}

func processSingleHandler(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	sortedArrays := sortSequential(payload.ToSort)
	timeTaken := time.Since(startTime)

	response := ResponsePayload{
		SortedArrays: sortedArrays,
		TimeNS:       timeTaken.Nanoseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func processConcurrentHandler(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	sortedArrays := sortConcurrent(payload.ToSort)
	timeTaken := time.Since(startTime)

	response := ResponsePayload{
		SortedArrays: sortedArrays,
		TimeNS:       timeTaken.Nanoseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/process-single", processSingleHandler)
	http.HandleFunc("/process-concurrent", processConcurrentHandler)

	http.ListenAndServe(":8000", nil)
}
