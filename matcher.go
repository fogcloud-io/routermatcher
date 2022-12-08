package routermatcher

import (
	"errors"
	"strings"
)

type KeyMatcher = func(sub string) (string, bool)
type PathSpliter = func(string) ([]string, error)

type Matcher struct {
	matchParam    KeyMatcher
	matchWildcard KeyMatcher
	split         PathSpliter
	root          *pathTrieNode
}

var (
	ErrInvalidPath = errors.New("invalid path")
)

var (
	RouterParamMatcher KeyMatcher = func(s string) (string, bool) {
		if len(s) <= 1 {
			return s, false
		}
		if s[0:1] == ":" {
			return s[1:], true
		} else {
			return s, false
		}
	}
	RouterWildcardMatcher KeyMatcher = func(sub string) (string, bool) {
		if len(sub) == 0 {
			return sub, false
		}
		if sub[0:1] == "*" {
			return "*", true
		} else {
			return sub, false
		}
	}
)

var (
	MqttTopicParamMatcher KeyMatcher = func(sub string) (string, bool) {
		if sub == "+" {
			return "+", true
		}
		return sub, false
	}
	MqttTopicWildMatcher KeyMatcher = func(sub string) (string, bool) {
		return sub, sub == "#"
	}

	MqttTopicPathSpliter PathSpliter = func(s string) ([]string, error) {
		if len(s) <= 0 {
			return nil, ErrInvalidPath
		}
		return strings.Split(s, "/"), nil
	}
)

func NewMatcher(paramMatcher, wildcardMatcher KeyMatcher, spliter PathSpliter) Matcher {
	return Matcher{
		matchParam:    paramMatcher,
		matchWildcard: wildcardMatcher,
		root: &pathTrieNode{
			value: "/",
			path:  "",
			child: make(map[string]*pathTrieNode),
		},
		split: spliter,
	}
}

func NewMqttTopicMatcher() Matcher {
	return NewMatcher(MqttTopicParamMatcher, MqttTopicWildMatcher, MqttTopicPathSpliter)
}

func NewRouterPathMatcher() Matcher {
	return NewMatcher(RouterParamMatcher, RouterWildcardMatcher, MqttTopicPathSpliter)
}

func (t *Matcher) AddPath(path string) error {
	return t.AddPathWithPriority(path, 0)
}

func (t *Matcher) AddPathWithPriority(path string, prior int) error {
	subs, err := t.split(path)
	if err != nil {
		return err
	}
	node := t.root
	key, isParam, isWildcard := "", false, false
	for _, sub := range subs {
		if node.child[sub] == nil {
			key, isParam = t.matchParam(sub)
			if !isParam {
				key, isWildcard = t.matchWildcard(sub)
			}
			node.child[sub] = &pathTrieNode{
				value:     key,
				path:      "",
				child:     make(map[string]*pathTrieNode),
				wordFlag:  false,
				paramFlag: isParam,
				wildFlag:  isWildcard,
				priority:  prior,
			}
			if node.child[sub].wildFlag {
				node.child[sub].path = path
				return nil
			}
		}
		node = node.child[sub]
	}
	node.wordFlag = true
	node.path = path
	return nil
}

func (t *Matcher) Match(dstTopic string) (matchedPath string, params map[string]string, ok bool) {
	subs, err := t.split(dstTopic)
	if err != nil {
		return "", nil, false
	}
	return t.root.match(subs)
}

func (t *Matcher) MatchWithAnonymousParams(dstTopic string) (matchedPath string, params []string, ok bool) {
	subs, err := t.split(dstTopic)
	if err != nil {
		return "", nil, false
	}
	return t.root.matchWithAnonymousParams(subs)
}
