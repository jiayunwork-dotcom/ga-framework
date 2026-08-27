package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"ga-framework/internal/engine"
)

func usage(w io.Writer) {
	fmt.Fprintln(w, "usage:")
	fmt.Fprintln(w, "  ga-framework run [--size N] [--genes N] [--generations N] [--mutate RATE] [--elite N] [--seed N]")
}

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "serve" {
		runServe(os.Args[2:])
		return
	}
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) < 1 || args[0] != "run" {
		usage(stderr)
		return fmt.Errorf("missing command")
	}
	var size, genes, generations, elite int
	var mutate float64
	var seed int64
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--size":
			if i+1 < len(args) {
				size, _ = strconv.Atoi(args[i+1])
				i++
			}
		case "--genes":
			if i+1 < len(args) {
				genes, _ = strconv.Atoi(args[i+1])
				i++
			}
		case "--generations":
			if i+1 < len(args) {
				generations, _ = strconv.Atoi(args[i+1])
				i++
			}
		case "--mutate":
			if i+1 < len(args) {
				mutate, _ = strconv.ParseFloat(args[i+1], 64)
				i++
			}
		case "--elite":
			if i+1 < len(args) {
				elite, _ = strconv.Atoi(args[i+1])
				i++
			}
		case "--seed":
			if i+1 < len(args) {
				seed, _ = strconv.ParseInt(args[i+1], 10, 64)
				i++
			}
		}
	}
	if genes == 0 {
		genes = 16
	}
	if size == 0 {
		size = 50
	}
	if generations == 0 {
		generations = 100
	}
	if mutate == 0 {
		mutate = 0.1
	}
	cfg := engine.Config{
		Size: size, Genes: genes, Generations: generations,
		TournamentK: 2, MutateRate: mutate, Elite: elite, Seed: seed,
	}
	res, err := engine.Run(cfg)
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "best_fitness=%.0f generations=%d genes=%d\n",
		res.BestFit, res.Generations, len(res.Best.Genes))
	return nil
}
