package indexes

import (
	"sort"
	"strings"
	"sync"
)

//Trie implements a trie data structure mapping strings to int64s
//that is safe for concurrent use.
type Trie struct {
	Name     rune
	Vals     []int64
	Children []*Trie
	Size     int
	mx       sync.RWMutex
}

//NewTrie constructs a new Trie.
func NewTrie() *Trie {
	return &Trie{}
}

//Len returns the number of entries in the trie.
func (t *Trie) Len() int {
	return t.Size
}

//Add adds a key and value to the trie.
func (t *Trie) Add(key string, value int64) {
	t.mx.Lock()
	t.Size++
	key = strings.ToLower(key)
	r := []rune(key)
	addHelper(r, value, t, 0)
	t.mx.Unlock()
}

func addHelper(r []rune, id int64, root *Trie, index int) {
	found := false
	for _, c := range root.Children {
		if c.Name == r[index] {
			root = c
			found = true
			break
		}
	}
	if !found {
		tNode := NewTrie()
		tNode.Name = r[index]
		//add new letter to trie
		root.Children = append(root.Children, tNode)
		//set root to new letter
		sort.Slice(root.Children, func(i, j int) bool {
			return root.Children[i].Name < root.Children[j].Name
		})
		root = tNode
	}
	if index != len(r)-1 {
		index++
		addHelper(r, id, root, index)
	} else {
		root.Vals = append(root.Vals, id)
	}
}

//Find finds `max` values matching `prefix`. If the trie
//is entirely empty, or the prefix is empty, or max == 0,
//or the prefix is not found, this returns a nil slice.
func (t *Trie) Find(prefix string, max int) []int64 {
	t.mx.RLock()
	defer t.mx.RUnlock()
	prefix = strings.ToLower(prefix)
	r := []rune(prefix)
	list := findHelper(t, r, 0, []int64{}, max)
	return list
}
func findHelper(root *Trie, r []rune, index int, list []int64, max int) []int64 {
	if index < len(r) {
		for _, c := range root.Children {
			if c.Name == r[index] {
				index++
				c.mx.Lock()
				list = findHelper(c, r, index, list, max)
				c.mx.Unlock()
				return list
			}
		}
		return []int64{}
	}
	if len(list) < max {
		if len(root.Vals) > 0 {
			for _, v := range root.Vals {
				if !contains(list, v) {
					list = append(list, v)
				}
			}
		}
		for _, c := range root.Children {
			if len(c.Vals) > 0 {
				for _, v := range c.Vals {
					if !contains(list, v) {
						list = append(list, v)
					}
				}
			}
			if len(c.Children) > 0 {

				list = findHelper(c, r, index, list, max)

			}
		}
	}
	return list
}

//Remove removes a key/value pair from the trie
//and trims branches with no values.
func (t *Trie) Remove(key string, value int64) {
	t.mx.Lock()
	defer t.mx.Unlock()
	t.Size--
	key = strings.ToLower(key)
	r := []rune(key)
	removeHelper(t, value, r, 0)
}

func removeHelper(root *Trie, value int64, r []rune, index int) {
	if index < len(r) {
		found := false
		var childIndex int
		for i, c := range root.Children {
			if c.Name == r[index] {
				found = true
				childIndex = i
				index++

				removeHelper(c, value, r, index)

				break
			}
		}
		if found == true {
			childRoot := root.Children[childIndex]
			if len(childRoot.Vals) < 1 && len(childRoot.Children) < 1 {
				root.Children = append(root.Children[:childIndex], root.Children[childIndex+1:]...)
			}
		}
	} else {
		if len(root.Vals) > 0 {
			var valIndex int
			for i, v := range root.Vals {
				if v == value {
					valIndex = i
				}
			}
			//remove id in slice using index
			root.Vals = append(root.Vals[:valIndex], root.Vals[valIndex+1:]...)
		}
	}
}

func contains(vals []int64, id int64) bool {
	for _, v := range vals {
		if v == id {
			return true
		}
	}
	return false
}
