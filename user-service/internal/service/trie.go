package service

import (
	"strings"
	"unicode"

	"github.com/namnv2496/user-service/internal/entity"
	"golang.org/x/text/unicode/norm"
)

type ITrie interface {
	SearchSuggestion(word string, limit int64) []*entity.LocationInfo
}

type TrieNode struct {
	children map[string]*TrieNode
	isWord   bool
	detail   *entity.LocationInfo
}

type Trie struct {
	root *TrieNode
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[string]*TrieNode),
		isWord:   false,
		detail:   &entity.LocationInfo{},
	}
}

func NewTrie(
	location ILocation,
) *Trie {
	root := NewTrieNode()
	for _, locationInfo := range location.GetLocation() {
		curr := root
		address := sannitizeVietnameseString(locationInfo.BuilderName)
		for _, char := range address {
			charStr := string(char)
			if nextNode, exist := curr.children[charStr]; exist {
				curr = nextNode
			} else {
				newNode := NewTrieNode()
				curr.children[charStr] = newNode
				curr = newNode
			}
		}
		curr.isWord = true
		curr.detail = locationInfo
	}
	return &Trie{
		root: root,
	}
}

func (_self *Trie) SearchSuggestion(word string, limit int64) []*entity.LocationInfo {
	word = sannitizeVietnameseString(word)
	root := _self.root
	curr := root
	for _, char := range word {
		if curr == nil {
			return nil
		}
		if node, exist := curr.children[string(char)]; exist {
			curr = node
		} else {
			return nil
		}
	}
	return searchWord(curr, limit)
}

func searchWord(trieNode *TrieNode, limit int64) []*entity.LocationInfo {
	result := make([]*entity.LocationInfo, 0)
	if trieNode.isWord {
		result = append(result, trieNode.detail)
	}
	if len(result) >= int(limit) {
		return result
	}
	for _, node := range trieNode.children {
		result = append(result, searchWord(node, limit)...)
		if len(result) >= int(limit) {
			return result
		}
	}
	return result
}

func sannitizeVietnameseString(searchWord string) string {
	searchWord = strings.ToLower(strings.TrimSpace(searchWord))
	searchWord = removePrefixes(searchWord)
	searchWord = removeSpecialChar(searchWord)
	searchWord = removeVietnameseAccents(searchWord)

	return strings.Join(strings.Fields(searchWord), " ")
}

func removeVietnameseAccents(searchWord string) string {
	replacements := map[string]string{
		"đ": "d",
		"Đ": "d",
	}
	result := searchWord
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}
	t := norm.NFD.String(result)
	b := make([]rune, 0, len(t))
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		b = append(b, r)
	}
	return string(b)
}

func removePrefixes(searchWord string) string {
	removeWords := []string{
		"phường", "xã", "thị trấn", "quận", "huyện", "tỉnh", "thành phố", "tp", "thị xã",
	}
	result := searchWord
	for _, word := range removeWords {
		result = strings.ReplaceAll(result, word, "")
	}
	return result
}

func removeSpecialChar(searchWord string) string {
	var b strings.Builder
	for _, r := range searchWord {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}
