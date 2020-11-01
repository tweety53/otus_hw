package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"regexp"
	"sort"
	"strings"
)

type wordCnt struct {
	w string
	n int
}

var spacesRegex = regexp.MustCompile(`\s+`)

func Top10(str string) []string {
	str = spacesRegex.ReplaceAllString(str, " ")
	if str == "" || str == " " {
		return []string{}
	}

	counters := collectCounters(str)
	wordsCnt := make([]wordCnt, 0, len(counters))
	for w, n := range counters {
		wordsCnt = append(wordsCnt, wordCnt{
			w: w,
			n: n,
		})
	}

	sort.Slice(wordsCnt, func(i, j int) bool {
		return wordsCnt[i].n > wordsCnt[j].n
	})

	res := make([]string, 0, 10)
	for i := range wordsCnt {
		if i == 10 {
			break
		}
		res = append(res, wordsCnt[i].w)
	}

	return res
}

func collectCounters(str string) map[string]int {
	counters := make(map[string]int)

	words := strings.Split(str, " ")
	for i := range words {
		counters[words[i]]++
	}

	return counters
}
