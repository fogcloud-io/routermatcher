package routermatcher

import "sort"

type priorities []*pathTrieNode

func (p priorities) Len() int {
	return len(p)
}

func (p priorities) Swap(l, r int) {
	p[l], p[r] = p[r], p[l]
}

func (p priorities) Less(l, r int) bool {
	return p[l].priority < p[r].priority
}

type pathTrieNode struct {
	value     string
	path      string
	child     map[string]*pathTrieNode
	wordFlag  bool
	paramFlag bool
	wildFlag  bool
	priority  int
}

func (t *pathTrieNode) match(subs []string) (matchedPath string, params map[string]string, ok bool) {
	params = map[string]string{}
	matchedPath, ok = t.backtrace(t, subs, params, 0)
	return
}

func (t *pathTrieNode) matchWithAnonymousParams(subs []string) (matchedPath string, params []string, ok bool) {
	matchedPath, ok = t.backtraceWithAnonymousParams(t, subs, &params, 0)

	for l, r := 0, len(params)-1; l < r; l, r = l+1, r-1 {
		params[l], params[r] = params[r], params[l]
	}

	return
}

func (t *pathTrieNode) backtrace(node *pathTrieNode, subs []string, params map[string]string, index int) (string, bool) {
	if index == len(subs) {
		return "", true
	}
	nodes := sortNodes(node.child)
	for _,subNode := range nodes {
		if subNode.value == subs[index] {
			matched, ok := t.backtrace(subNode, subs, params, index+1)
			if ok {
				if index == len(subs)-1 {
					//if subNode.wordFlag {
					//	return subNode.path, true
					//} else {
					//	continue
					//}
					return subNode.path, true
				} else {
					return matched, true
				}
			}
		} else if subNode.paramFlag {
			matched, ok := t.backtrace(subNode, subs, params, index+1)
			if ok {
				if index == len(subs)-1 {
					if subNode.wordFlag {
						params[subNode.value] = subs[index]
						return subNode.path, true
					} else {
						continue
					}
				}
				params[subNode.value] = subs[index]
				return matched, true
			}
		} else if subNode.wildFlag {
			return subNode.path, true
		} else {
			continue
		}
	}
	return "", false
}

func sortNodes(nodes map[string]*pathTrieNode) (res []*pathTrieNode) {
	n := make(priorities, len(nodes))
	i := 0
	for _, v := range nodes {
		n[i] = v
		i++
	}
	sort.Sort(n)
	res = n
	return
}

func (t *pathTrieNode) backtraceWithAnonymousParams(node *pathTrieNode, subs []string, params *[]string, index int) (string, bool) {
	if index == len(subs) {
		return "", true
	}
	nodes := sortNodes(node.child)
	for _,subNode := range nodes {
		if subNode.value == subs[index] {
			matched, ok := t.backtraceWithAnonymousParams(subNode, subs, params, index+1)
			if ok {
				if index == len(subs)-1 {
					//if subNode.wordFlag {
					//	return subNode.path, true
					//} else {
					//	continue
					//}
					return subNode.path, true
				} else {
					return matched, true
				}
			}
		} else if subNode.paramFlag {
			matched, ok := t.backtraceWithAnonymousParams(subNode, subs, params, index+1)
			if ok {
				if index == len(subs)-1 {
					if subNode.wordFlag {
						*params = append(*params, subs[index])
						return subNode.path, true
					} else {
						continue
					}
				}
				*params = append(*params, subs[index])
				return matched, true
			}
		} else if subNode.wildFlag {
			return subNode.path, true
		} else {
			continue
		}
	}
	return "", false
}
