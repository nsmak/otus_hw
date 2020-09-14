package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"sort"
	"strings"
)

type RatedWord struct {
	Text  string
	Count int
}

type RatedWords []RatedWord

func (rWords RatedWords) Words() []string {
	words := make([]string, 0, len(rWords))
	for _, word := range rWords {
		words = append(words, word.Text)
	}
	return words
}

func Top10(s string) []string {
	if s == "" {
		return nil
	}

	uniqWords := make(map[string]int)
	words := strings.Fields(s) // Вместо strings.Split, так как показалось удобнее, чтобы убрать пробелы-мусор.
	for _, word := range words {
		uniqWords[word]++
	}

	ratedWords := make([]RatedWord, 0, len(uniqWords))
	for k, v := range uniqWords {
		ratedWords = append(ratedWords, RatedWord{Text: k, Count: v})
	}
	sort.Slice(ratedWords, func(i, j int) bool {
		return ratedWords[i].Count > ratedWords[j].Count
	})

	if len(ratedWords) > 10 {
		return RatedWords(ratedWords[0:10]).Words()
	}
	return RatedWords(ratedWords).Words()
}
