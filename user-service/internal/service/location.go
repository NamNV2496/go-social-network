package service

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/namnv2496/user-service/internal/configs"
	"github.com/namnv2496/user-service/internal/entity"
	"github.com/namnv2496/user-service/internal/service/location"
)

type ILocation interface {
	GetLocation() []*entity.LocationInfo
	SearchLocation(word string) []*entity.LocationInfo
}

type Location struct {
	locations         []*entity.LocationInfo
	dictionaryMatcher location.IDictionaryMatcher
	fuzzyMatcher      location.IFuzzyMatcher
}

func NewLocation(
	conf *configs.Config,
) ILocation {
	records, err := readCsv(conf.Location.FilePath)
	if err != nil {
		return nil
	}
	locations, locationIndex := mapLocationInfo(records)
	return &Location{
		locations:         locations,
		dictionaryMatcher: location.NewDictionaryMatcher(locationIndex),
		fuzzyMatcher:      location.NewFuzzyMatcher(locationIndex, 0.8),
	}
}

func readCsv(filePath string) ([][]string, error) {
	f, err := os.Open("./internal/service/location/data/" + filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func mapLocationInfo(records [][]string) ([]*entity.LocationInfo, *entity.LocationIndex) {
	//      0             1           2           3.         4            5            6.       7               8          9                      10
	// Tên tỉnh cũ	| Mã TP cũ	| Tên QH cũ	| Tên PX cũ	| Mã QH cũ	| Mã PX cũ	| Tên PX mới	| Mã PX mới | Độ khớp | QH_variants_normalized | PX_variants_normalized
	index := &entity.LocationIndex{
		NormalizedCityName:     make(map[string][]*entity.LocationInfo),
		NormalizedDistrictName: make(map[string][]*entity.LocationInfo),
		NormalizedWardName:     make(map[string][]*entity.LocationInfo),
		NormalizedName:         make(map[string][]*entity.LocationInfo),
	}

	locations := make([]*entity.LocationInfo, 0)
	for i := 1; i < len(records); i++ {
		record := records[i]
		buildName := fmt.Sprintf("%s %s %s", record[3], record[2], record[0])
		buildName = strings.ReplaceAll(buildName, "xã", "")
		buildName = strings.ReplaceAll(buildName, "phường", "")
		buildName = strings.ReplaceAll(buildName, "quận", "")
		buildName = strings.ReplaceAll(buildName, "huyện", "")
		buildName = strings.ReplaceAll(buildName, "thành phố", "")
		variants := append(strings.Split(strings.TrimSpace(record[9]), "|"), strings.Split(strings.TrimSpace(record[10]), "|")...)
		location := &entity.LocationInfo{
			CityV1:     strings.TrimSpace(record[1]),
			CityV1Name: strings.TrimSpace(record[0]),
			// CityV2:       "new_id",
			// CityV2Name:   "new_name",
			District:     strings.TrimSpace(record[4]),
			DistrictName: strings.TrimSpace(record[2]),
			WardV1:       strings.TrimSpace(record[5]),
			WardV1Name:   strings.TrimSpace(record[3]),
			WardV2:       strings.TrimSpace(record[7]),
			WardV2Name:   strings.TrimSpace(record[6]),
			BuilderName:  buildName,
			Variants:     variants,
		}
		locations = append(locations, location)
		// index the location - store all matches, not just first
		if location.WardV1Name != "" {
			normalizedWardName := normalizeWard(location.WardV1Name)
			index.NormalizedWardName[normalizeKey(normalizedWardName)] = append(index.NormalizedWardName[normalizeKey(normalizedWardName)], location)
		}
		if location.WardV2Name != "" {
			normalizedWardName := normalizeWard(location.WardV2Name)
			index.NormalizedWardName[normalizeKey(normalizedWardName)] = append(index.NormalizedWardName[normalizeKey(normalizedWardName)], location)
		}
		if location.DistrictName != "" {
			normalizedDistrictName := normalizeDistrict(location.DistrictName)
			index.NormalizedDistrictName[normalizeKey(normalizedDistrictName)] = append(index.NormalizedDistrictName[normalizeKey(normalizedDistrictName)], location)
		}
		if location.CityV1Name != "" {
			normalizedCityV1Name := normalizeCity(location.CityV1Name)
			index.NormalizedCityName[normalizeKey(normalizedCityV1Name)] = append(index.NormalizedCityName[normalizeKey(normalizedCityV1Name)], location)
			index.NormalizedName[normalizeKey(normalizedCityV1Name)] = append(index.NormalizedName[normalizeKey(normalizedCityV1Name)], location)
		}
		for _, variant := range variants {
			normalizedVariant := normalizeKey(variant)
			if normalizedVariant != "" {
				index.NormalizedName[normalizedVariant] = append(index.NormalizedName[normalizedVariant], location)
				// Extract abbreviations from variants (e.g., "ben tre" -> "ben_tr", "can tho" -> "can_th")
				// abbrev := extractAbbreviation(normalizedVariant)
				// if abbrev != "" && abbrev != normalizedVariant {
				// 	index.NormalizedName[abbrev] = append(index.NormalizedName[abbrev], location)
				// }
			}
		}
	}
	return locations, index
}

// extractAbbreviation generates abbreviation from a normalized location name
// e.g., "ben_tre" -> "ben_tr", "can_tho" -> "can_th"
func extractAbbreviation(normalizedName string) string {
	parts := strings.Split(normalizedName, "_")
	if len(parts) < 2 {
		return ""
	}

	// Build abbreviation by taking first 2-3 chars of each word
	var abbrev []string
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		// For short words (1-2 chars), use whole word; for longer, use first 2 chars
		if len(part) <= 2 {
			abbrev = append(abbrev, part)
		} else {
			abbrev = append(abbrev, part[:2])
		}
	}

	if len(abbrev) == 0 {
		return ""
	}

	return strings.Join(abbrev, "_")
}

func (_self *Location) GetLocation() []*entity.LocationInfo {
	return _self.locations
}

func (_self *Location) SearchLocation(word string) []*entity.LocationInfo {
	result := _self.extract(word)
	if result == nil {
		return []*entity.LocationInfo{}
	}
	districtUnique := make(map[string]bool, 0)
	resp := make([]*entity.LocationInfo, 0)
	for _, loc := range result {
		key := fmt.Sprintf("%s_%s", loc.Location.District, loc.Location.CityV1)
		if _, exist := districtUnique[key]; !exist {
			resp = append(resp, loc.Location)
			districtUnique[key] = true
		}
	}
	return resp
}

func normalizeQuery(in string) string {
	return location.NormalizeWhitespace(location.RemoveVietnameseDiacritics(strings.TrimSpace(strings.ToLower(in))))
}

func normalizeDistrict(in string) string {
	normalizedDistrictName := normalizeQuery(in)
	normalizedDistrictName = strings.ReplaceAll(normalizedDistrictName, "quan", "")
	normalizedDistrictName = strings.ReplaceAll(normalizedDistrictName, "huyen", "")
	return normalizedDistrictName
}
func normalizeWard(in string) string {
	normalizedWardName := normalizeQuery(in)
	normalizedWardName = strings.ReplaceAll(normalizedWardName, "xa", "")
	normalizedWardName = strings.ReplaceAll(normalizedWardName, "phuong", "")
	normalizedWardName = strings.ReplaceAll(normalizedWardName, "thi tran", "")
	return normalizedWardName
}
func normalizeCity(in string) string {
	normalizedCityName := normalizeQuery(in)
	normalizedCityName = strings.ReplaceAll(normalizedCityName, "thanh pho", "")
	normalizedCityName = strings.ReplaceAll(normalizedCityName, "tinh", "")
	return normalizedCityName
}

func normalizeKey(in string) string {
	return strings.ReplaceAll(location.NormalizeWhitespace(location.RemoveVietnameseDiacritics(strings.TrimSpace(strings.ToLower(in)))), " ", "_")
}

func (_self *Location) extract(query string) []*entity.MatchResult {
	tokens := strings.Fields(normalizeQuery(query))

	if len(tokens) == 0 {
		return nil
	}
	if len(tokens) > 20 {
		tokens = tokens[:20]
	}

	results := make([]*entity.MatchResult, 0)
	var bestScore float64
	var matchedString string

	// Try dictionary match first (exact hierarchical lookup), then fuzzy
	for n := min(4, len(tokens)); n >= 1; n-- {
		for i := 0; i <= len(tokens)-n; i++ {
			ngram := normalizeKey(strings.Join(tokens[i:i+n], " "))

			// 1. Try dictionary match first (ward, district, city hierarchical lookup)
			if locs, found := _self.dictionaryMatcher.Match(ngram, "default"); found && len(locs) > 0 {
				score := 1.0
				if score >= bestScore {
					for _, loc := range locs {
						results = append(results, &entity.MatchResult{Location: loc, MatchedText: ngram, Score: score, MatchType: "dictionary"})
					}
				}
			}

			// 1b. Try matching abbreviation of the ngram
			abbrev := extractAbbreviation(ngram)
			if abbrev != "" && abbrev != ngram {
				if locs, found := _self.dictionaryMatcher.Match(abbrev, "default"); found && len(locs) > 0 {
					score := 1.0
					if score >= bestScore {
						for _, loc := range locs {
							results = append(results, &entity.MatchResult{Location: loc, MatchedText: ngram, Score: score, MatchType: "dictionary"})
						}
					}
				}
			}

			// 2. Try fuzzy match if dictionary didn't find anything
			if locs, score := _self.fuzzyMatcher.Match(ngram); locs != nil && score >= 0.8 && score > bestScore {
				for _, loc := range locs {
					results = append(results, &entity.MatchResult{Location: loc, MatchedText: ngram, Score: score, MatchType: "fuzzy"})
				}
				bestScore = score
				matchedString = ngram
				break
			}
		}
	}
	if len(results) > 1 {
		// Try to use hierarchical info to filter results if we have multiple tokens
		return _self.filterResultsByHierarchy(results, tokens, matchedString)
	}
	return results
}

// filterResultsByHierarchy tries to filter multiple results using hierarchical information
// It looks for city/district names in the remaining tokens to disambiguate duplicate results
func (_self *Location) filterResultsByHierarchy(results []*entity.MatchResult, tokens []string, matchedString string) []*entity.MatchResult {
	if len(results) <= 1 {
		return results
	}

	cityMap := make(map[string]*entity.MatchResult, 0)
	districtMap := make(map[string]*entity.MatchResult, 0)
	for _, result := range results {
		cityMap[normalizeKey(normalizeCity(result.Location.CityV1Name))] = result
		districtMap[normalizeKey(normalizeDistrict(result.Location.DistrictName))] = result
	}
	// Try to find city or district name in the remaining tokens
	// Search from longest to shortest n-grams
	matchResults := make([]*entity.MatchResult, 0)
	for n := min(3, len(tokens)); n >= 1; n-- {
		for i := 0; i <= len(tokens)-n; i++ {
			ngram := normalizeKey(strings.Join(tokens[i:i+n], " "))
			if ngram == normalizeKey(matchedString) {
				continue
			}
			// Try to match city name
			if cityLocs, found := _self.dictionaryMatcher.Match(ngram, "city"); found {
				if len(cityLocs) > 0 {
					for _, loc := range cityLocs {
						if normalizeKey(normalizeDistrict(loc.DistrictName)) == matchedString {
							matchResults = append(matchResults, &entity.MatchResult{Location: loc, MatchedText: ngram, Score: 1.0, MatchType: "city"})
						}
					}
					return matchResults
				}
			}

			// Try to match district name
			if districtLocs, found := _self.dictionaryMatcher.Match(ngram, "district"); found {
				if len(districtLocs) > 0 {
					for _, loc := range districtLocs {
						matchResults = append(matchResults, &entity.MatchResult{Location: loc, MatchedText: ngram, Score: 1.0, MatchType: "city"})
					}
				}
				return matchResults
			}
		}
	}
	if len(matchResults) > 0. {
		return matchResults
	}

	// If no hierarchical match found, return all results
	return results
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
