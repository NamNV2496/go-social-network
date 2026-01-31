package location

import (
	"regexp"

	"github.com/namnv2496/user-service/internal/entity"
)

type IFuzzyMatcher interface {
	Match(normalizedText string) ([]*entity.LocationInfo, float64)
}

type searchPattern struct {
	pattern *regexp.Regexp
	expand  string
}

type fuzzyMatcher struct {
	threshold float64
	index     *entity.LocationIndex
	patterns  []*searchPattern
}

func NewFuzzyMatcher(indexLocation *entity.LocationIndex, threshold float64) *fuzzyMatcher {
	return &fuzzyMatcher{
		threshold: threshold,
		index:     indexLocation,
		patterns: []*searchPattern{
			// District: q7, q.7, q 7 -> quan 7
			{regexp.MustCompile(`(?i)^q\.?\s*(\d+)$`), "quan"},
			{regexp.MustCompile(`(?i)^quan\s*(\d+)$`), "quan"},
			{regexp.MustCompile(`(?i)^district\s*(\d+)$`), "quan"},
			{regexp.MustCompile(`(?i)^d\.?\s*(\d+)$`), "quan"},
			// Subdistrict: huyen, thi tran
			{regexp.MustCompile(`(?i)^huyen\s+(.+)$`), "huyen"},
			{regexp.MustCompile(`(?i)^thi\s+tran\s+(.+)$`), "thi tran"},
			// Ward: p7, p.7 -> phuong 7
			{regexp.MustCompile(`(?i)^p\.?\s*(\d+)$`), "phuong"},
			{regexp.MustCompile(`(?i)^phuong\s*(\d+)$`), "phuong"},
			{regexp.MustCompile(`(?i)^xa\s*(\d+)$`), "xa"},
			// City: thanh pho, tp, tinh
			{regexp.MustCompile(`(?i)^thanh\s+pho\s+(.+)$`), "thanh pho"},
			{regexp.MustCompile(`(?i)^tp\.?\s+(.+)$`), "thanh pho"},
			{regexp.MustCompile(`(?i)^tinh\s+(.+)$`), "tinh"},
			// Underscore-separated formats
			{regexp.MustCompile(`(?i)^quan_(.+)$`), "quan"},
			{regexp.MustCompile(`(?i)^thanh_pho_(.+)$`), "thanh pho"},
			{regexp.MustCompile(`(?i)^tinh_(.+)$`), "tinh"},
			{regexp.MustCompile(`(?i)^huyen_(.+)$`), "huyen"},
			{regexp.MustCompile(`(?i)^phuong_(.+)$`), "phuong"},
		},
	}
}

// Match tries abbreviation expansion first, then Levenshtein similarity, returns all matches.
func (_self *fuzzyMatcher) Match(normalizedText string) ([]*entity.LocationInfo, float64) {
	if _self.index == nil || len(normalizedText) == 0 {
		return nil, 0
	}

	// 1. Try abbreviation expansion
	for _, p := range _self.patterns {
		if matches := p.pattern.FindStringSubmatch(normalizedText); len(matches) == 2 {
			if locs, found := _self.index.NormalizedName[p.expand+matches[1]]; found {
				return locs, 1.0
			}
		}
	}

	// 2. Levenshtein similarity matching
	var bestMatches []*entity.LocationInfo
	var bestScore float64

	for indexedName, locs := range _self.index.NormalizedName {
		score := levenshteinSimilarity(normalizedText, indexedName)
		if score >= _self.threshold {
			if score > bestScore {
				bestScore = score
				bestMatches = locs
			}
		}
	}

	return bestMatches, bestScore
}

func levenshteinSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}
	r1, r2 := []rune(s1), []rune(s2)
	maxLen := max(len(r1), len(r2))
	if maxLen == 0 {
		return 1.0
	}
	return 1.0 - float64(levenshteinDistance(r1, r2))/float64(maxLen)
}

func levenshteinDistance(r1, r2 []rune) int {
	len1, len2 := len(r1), len(r2)
	if len1 == 0 {
		return len2
	}
	if len2 == 0 {
		return len1
	}

	prev := make([]int, len2+1)
	curr := make([]int, len2+1)
	for j := range prev {
		prev[j] = j
	}

	for i := 1; i <= len1; i++ {
		curr[0] = i
		for j := 1; j <= len2; j++ {
			cost := 0
			if r1[i-1] != r2[j-1] {
				cost = 1
			}
			curr[j] = min(prev[j]+1, curr[j-1]+1, prev[j-1]+cost)
		}
		prev, curr = curr, prev
	}
	return prev[len2]
}
