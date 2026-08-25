package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type GAConfig struct {
	PopSize     int     `json:"pop_size"`
	Genes       int     `json:"genes"`
	MaxGen      int     `json:"max_gen"`
	MutateRate  float64 `json:"mutate_rate"`
	CrossRate   float64 `json:"cross_rate"`
	Elite       int     `json:"elite"`
	Seed        int64   `json:"seed"`
	Selection   string  `json:"selection"`
	Crossover   string  `json:"crossover"`
	Mutation    string  `json:"mutation"`
	Fitness     string  `json:"fitness"`
	Islands     int     `json:"islands"`
	MigrateRate float64 `json:"migrate_rate"`
	MigrateFreq int     `json:"migrate_freq"`
	Encoding    string  `json:"encoding"`
	BoundsLo    float64 `json:"bounds_lo"`
	BoundsHi    float64 `json:"bounds_hi"`
}

func DefaultConfig() GAConfig {
	return GAConfig{
		PopSize:    50,
		Genes:      20,
		MaxGen:     200,
		MutateRate: 0.05,
		CrossRate:  0.8,
		Elite:      2,
		Selection:  "tournament",
		Crossover:  "single_point",
		Mutation:   "flip",
		Fitness:    "onemax",
		Encoding:   "binary",
		BoundsLo:   -5.0,
		BoundsHi:   5.0,
	}
}

func (c GAConfig) Validate() error {
	var errs []string
	if c.PopSize < 2 {
		errs = append(errs, "pop_size must be >= 2")
	}
	if c.Genes < 1 {
		errs = append(errs, "genes must be >= 1")
	}
	if c.MaxGen < 1 {
		errs = append(errs, "max_gen must be >= 1")
	}
	if c.MutateRate < 0 || c.MutateRate > 1 {
		errs = append(errs, "mutate_rate must be in [0, 1]")
	}
	if c.CrossRate < 0 || c.CrossRate > 1 {
		errs = append(errs, "cross_rate must be in [0, 1]")
	}
	if c.Elite < 0 {
		errs = append(errs, "elite must be >= 0")
	}
	if c.Elite >= c.PopSize {
		errs = append(errs, "elite must be < pop_size")
	}
	if c.Islands < 0 {
		errs = append(errs, "islands must be >= 0")
	}
	if c.BoundsLo >= c.BoundsHi {
		errs = append(errs, "bounds_lo must be < bounds_hi")
	}
	validSel := map[string]bool{"tournament": true, "roulette": true, "rank": true}
	if !validSel[c.Selection] {
		errs = append(errs, fmt.Sprintf("unknown selection: %q", c.Selection))
	}
	validCx := map[string]bool{"single_point": true, "two_point": true, "uniform": true}
	if !validCx[c.Crossover] {
		errs = append(errs, fmt.Sprintf("unknown crossover: %q", c.Crossover))
	}
	validMut := map[string]bool{"flip": true, "gaussian": true, "swap": true}
	if !validMut[c.Mutation] {
		errs = append(errs, fmt.Sprintf("unknown mutation: %q", c.Mutation))
	}
	validEnc := map[string]bool{"binary": true, "real": true, "permutation": true, "integer": true}
	if !validEnc[c.Encoding] {
		errs = append(errs, fmt.Sprintf("unknown encoding: %q", c.Encoding))
	}
	if len(errs) > 0 {
		return bindEmptyErr(errors.New(strings.Join(errs, "; ")))
	}
	return nil
}

func (c GAConfig) MarshalJSON() ([]byte, error) {
	type Alias GAConfig
	return json.Marshal((Alias)(c))
}

func (c GAConfig) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(c)
}

func ReadJSON(r io.Reader) (GAConfig, error) {
	var c GAConfig
	dec := json.NewDecoder(r)
	if err := dec.Decode(&c); err != nil {
		return GAConfig{}, fmt.Errorf("decode config: %w", err)
	}
	return c, nil
}

func Merge(base, overlay GAConfig) GAConfig {
	result := base
	if overlay.PopSize > 0 {
		result.PopSize = overlay.PopSize
	}
	if overlay.Genes > 0 {
		result.Genes = overlay.Genes
	}
	if overlay.MaxGen > 0 {
		result.MaxGen = overlay.MaxGen
	}
	if overlay.MutateRate > 0 {
		result.MutateRate = overlay.MutateRate
	}
	if overlay.CrossRate > 0 {
		result.CrossRate = overlay.CrossRate
	}
	if overlay.Elite > 0 {
		result.Elite = overlay.Elite
	}
	if overlay.Seed != 0 {
		result.Seed = overlay.Seed
	}
	if overlay.Selection != "" {
		result.Selection = overlay.Selection
	}
	if overlay.Crossover != "" {
		result.Crossover = overlay.Crossover
	}
	if overlay.Mutation != "" {
		result.Mutation = overlay.Mutation
	}
	if overlay.Fitness != "" {
		result.Fitness = overlay.Fitness
	}
	if overlay.Islands > 0 {
		result.Islands = overlay.Islands
	}
	if overlay.MigrateRate > 0 {
		result.MigrateRate = overlay.MigrateRate
	}
	if overlay.MigrateFreq > 0 {
		result.MigrateFreq = overlay.MigrateFreq
	}
	if overlay.Encoding != "" {
		result.Encoding = overlay.Encoding
	}
	if overlay.BoundsLo != 0 || overlay.BoundsHi != 0 {
		result.BoundsLo = overlay.BoundsLo
		result.BoundsHi = overlay.BoundsHi
	}
	return result
}

func (c GAConfig) Summary() string {
	return fmt.Sprintf("pop=%d genes=%d gen=%d sel=%s cx=%s mut=%s enc=%s",
		c.PopSize, c.Genes, c.MaxGen, c.Selection, c.Crossover, c.Mutation, c.Encoding)
}
