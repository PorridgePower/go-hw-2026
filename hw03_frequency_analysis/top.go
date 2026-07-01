package hw03frequencyanalysis

import (
	"cmp"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"
)

var expr = regexp.MustCompile(`(?:[\pL\p{So}]|-)+`)

func Top10(input string) []string {
	result := make(map[string]int)

	for _, w := range expr.FindAllString(input, -1) {
		result[strings.ToLower(w)]++
	}

	delete(result, "-")

	keys := slices.SortedFunc(maps.Keys(result), func(a, b string) int {
		if result[a] != result[b] {
			return cmp.Compare(result[b], result[a])
		}
		return cmp.Compare(a, b)
	})

	resLen := min(len(keys), 10)
	res := make([]string, 0, resLen)
	for _, k := range keys[:resLen] {
		res = append(res, k)
		fmt.Println(k, result[k])
	}
	return res
}
