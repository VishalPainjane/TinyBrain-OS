package v2

import "sync"

// PrefixNode represents a cached block of tokens.
type PrefixNode struct {
	BlockID  int
	Tokens   []int32 // Will always be BlockSize length
	Children map[int32]*PrefixNode
	Parent   *PrefixNode
}

// PrefixCache manages the Radix Trie of cached blocks.
type PrefixCache struct {
	mu       sync.RWMutex
	Root     *PrefixNode
	BlockMap map[int]*PrefixNode // Maps BlockID -> Node for O(1) eviction
}

func NewPrefixCache() *PrefixCache {
	return &PrefixCache{
		Root:     &PrefixNode{BlockID: -1, Children: make(map[int32]*PrefixNode)},
		BlockMap: make(map[int]*PrefixNode),
	}
}

// Match finds the longest prefix of fully cached blocks.
// Returns the physical block IDs, and the number of tokens matched.
func (p *PrefixCache) Match(tokens []int32, blockSize int) ([]int, int) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var matchedBlocks []int
	matchedTokens := 0
	
	curr := p.Root
	for matchedTokens+blockSize <= len(tokens) {
		chunk := tokens[matchedTokens : matchedTokens+blockSize]
		
		// Find matching child based on first token of chunk
		child, ok := curr.Children[chunk[0]]
		if !ok {
			break
		}
		
		// Verify full chunk matches
		match := true
		for i := 0; i < blockSize; i++ {
			if child.Tokens[i] != chunk[i] {
				match = false
				break
			}
		}
		if !match {
			break
		}
		
		matchedBlocks = append(matchedBlocks, child.BlockID)
		matchedTokens += blockSize
		curr = child
	}
	
	return matchedBlocks, matchedTokens
}

// Insert adds a new fully-formed block to the cache.
// `parentBlocks` is used to traverse to the insertion point.
func (p *PrefixCache) Insert(parentBlocks []int, tokens []int32, newBlockID int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Traverse to the parent node
	curr := p.Root
	for _, blockID := range parentBlocks {
		if node, ok := p.BlockMap[blockID]; ok {
			curr = node
		} else {
			// Parent block was evicted or doesn't exist; cannot attach child
			return
		}
	}
	
	if len(tokens) == 0 {
		return
	}
	
	// Check if already exists
	if _, ok := curr.Children[tokens[0]]; ok {
		return
	}
	
	newNode := &PrefixNode{
		BlockID:  newBlockID,
		Tokens:   append([]int32(nil), tokens...),
		Children: make(map[int32]*PrefixNode),
		Parent:   curr,
	}
	
	curr.Children[tokens[0]] = newNode
	p.BlockMap[newBlockID] = newNode
}

// Evict removes a block from the cache when LRU reclaims it.
func (p *PrefixCache) Evict(blockID int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	node, ok := p.BlockMap[blockID]
	if !ok {
		return
	}
	
	// Remove from parent
	if node.Parent != nil {
		delete(node.Parent.Children, node.Tokens[0])
	}
	
	// Prune all children recursively
	var prune func(n *PrefixNode)
	prune = func(n *PrefixNode) {
		delete(p.BlockMap, n.BlockID)
		for _, child := range n.Children {
			prune(child)
		}
	}
	
	for _, child := range node.Children {
		prune(child)
	}
	
	delete(p.BlockMap, blockID)
}
