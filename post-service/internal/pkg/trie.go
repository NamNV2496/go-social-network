package pkg

type ITrie interface {
	Insert(word string)
	Search(word string) bool
}

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
		isEnd:    false,
	}
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: NewTrieNode(),
	}
}

func (_self *Trie) GetRoot() *TrieNode {
	return _self.root
}
func (_self *Trie) Insert(word string) {
	node := _self.root
	for _, char := range word {
		if _, ok := node.children[char]; !ok {
			node.children[char] = NewTrieNode()
		}
		node = node.children[char]
	}
	node.isEnd = true
}

func (_self *Trie) Search(word string) bool {
	node := _self.root
	for _, char := range word {
		if _, ok := node.children[char]; !ok {
			return false
		}
		node = node.children[char]
	}
	return node.isEnd
}

func (_self *Trie) SearchSubstring(text string) (bool, string) {
	matchWords := ""
	for i := 0; i < len(text); i++ {
		node := _self.root
		for j := i; j < len(text); j++ {
			matchWords += string(text[j])
			char := rune(text[j])
			if _, ok := node.children[char]; !ok {
				matchWords = ""
				break
			}
			node = node.children[char]
			if node.isEnd {
				return true, matchWords
			}
		}
	}
	return false, ""
}
