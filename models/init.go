package models

import (
	"sort"
	"strings"
)

var models map[string]map[string]int = make(map[string]map[string]int)

// Struct used to sort trigrams
type valSorter struct {
	keys   []string
	values []int
}

// Returns list of all models
func GetModels() map[string]map[string]int {
	return models
}

// Create a list of trigrams in content sorted by frequency.
func GetOrderedModel(content string) []string {
	var (
		trigrams map[string]int = make(map[string]int)
		trigram  string
		runes    []rune = []rune(strings.ToLower(content))
	)

	for i := 0; i < len(runes)-2; i++ {
		trigram = string(runes[i : i+3])

		if _, ok := trigrams[trigram]; ok {
			trigrams[trigram]++
		} else {
			trigrams[trigram] = 1
		}
	}

	vs := getValSorter(trigrams)
	vs.Sort()

	return vs.keys
}

func getValSorter(m map[string]int) *valSorter {
	vs := &valSorter{
		keys:   make([]string, 0, len(m)),
		values: make([]int, 0, len(m)),
	}

	for k, v := range m {
		vs.keys = append(vs.keys, k)
		vs.values = append(vs.values, v)
	}

	return vs
}

func (vs *valSorter) Sort() {
	sort.Sort(sort.Reverse(vs))
}

func (vs *valSorter) Len() int {
	return len(vs.values)
}

func (vs *valSorter) Less(i, j int) bool {
	return vs.values[i] < vs.values[j]
}

func (vs *valSorter) Swap(i, j int) {
	vs.values[i], vs.values[j] = vs.values[j], vs.values[i]
	vs.keys[i], vs.keys[j] = vs.keys[j], vs.keys[i]
}
