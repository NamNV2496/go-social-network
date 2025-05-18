package pkg

import (
	"github.com/google/wire"
)

var TrieNodeSet = wire.NewSet(
	NewTrie,
	wire.Bind(new(ITrie), new(*Trie)),
)
