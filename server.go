package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"ga-framework/internal/engine"
)

func runServe(args []string) {
	addr := ":8080"
	for i := 0; i < len(args); i++ {
		if args[i] == "--addr" && i+1 < len(args) {
			addr = args[i+1]
			i++
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/run", handleRun)

	fmt.Printf("ga-framework serving on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "serve: %v\n", err)
		os.Exit(1)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Size        int     `json:"size"`
		Genes       int     `json:"genes"`
		Generations int     `json:"generations"`
		MutateRate  float64 `json:"mutate_rate"`
		Elite       int     `json:"elite"`
		Seed        int64   `json:"seed"`
	}

	if ct := r.Header.Get("Content-Type"); ct == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		q := r.URL.Query()
		req.Size, _ = strconv.Atoi(q.Get("size"))
		req.Genes, _ = strconv.Atoi(q.Get("genes"))
		req.Generations, _ = strconv.Atoi(q.Get("generations"))
		req.MutateRate, _ = strconv.ParseFloat(q.Get("mutate_rate"), 64)
		req.Elite, _ = strconv.Atoi(q.Get("elite"))
		req.Seed, _ = strconv.ParseInt(q.Get("seed"), 10, 64)
	}

	if req.Genes == 0 {
		req.Genes = 16
	}
	if req.Size == 0 {
		req.Size = 50
	}
	if req.Generations == 0 {
		req.Generations = 100
	}
	if req.MutateRate == 0 {
		req.MutateRate = 0.1
	}

	cfg := engine.Config{
		Size:        req.Size,
		Genes:       req.Genes,
		Generations: req.Generations,
		TournamentK: 2,
		MutateRate:  req.MutateRate,
		Elite:       req.Elite,
		Seed:        req.Seed,
	}

	res, err := engine.Run(cfg)
	if err != nil {
		http.Error(w, "run: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"best_fitness": res.BestFit,
		"generations":  res.Generations,
		"gene_length":  len(res.Best.Genes),
	})
}
