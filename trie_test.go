package routermatcher

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	topics = []string {
		"$fogcloud/+/+/thing/event/+/post_reply",

		"$fogcloud/+/+/thing/event/property/post_reply",
		"$fogcloud/+/+/thing/event/property/set",
		"$fogcloud/+/+/thing/service/+",

		"$fogcloud/+/+/thing/event/property/post",
		"$fogcloud/+/+/thing/event/+/post",
		"$fogcloud/+/+/thing/service/+/reply",

		"$fogcloud/+/+/event/#/hello",

		"$fogcloud/+/+/all/#",
	}

	routerPaths = []string{
		"/projects/:project_id/users",
		"/projects/:project_id/users/:user_id",
		"/projects/:project_id/*Prod",

		"/groups/:group_id/users/:user_id",
		"/groups/:group_id/owner",
		"/groups/:group_id/users",
	}
)

func TestMqttTopicMatch(t *testing.T)  {
	m := NewMqttTopicMatcher()
	var err error
	for i := range topics {
		err = m.AddPathWithPriority(topics[i], i)
		assert.NoError(t, err)
	}

	m1 := NewMqttTopicMatcher()
	m1.AddPathWithPriority(topics[0], 2)
	m1.AddPathWithPriority(topics[1], 1)
	m1.AddPathWithPriority(topics[8], 3)

	dstTopics := []string{
		"$fogcloud/pk1/device1/thing/event/hello/post",
		"$fogcloud/pk1/device1/thing/event/property/post",
		"$fogcloud/pk2/device2/thing/service/lightMode",
		"$fogcloud/pk3/device3/thing/event/property/post_reply",

		"$fogcloud/pk3/device3/thing/event/alarm/post_reply/err",
		"$fogcloud/pk3/thing/",

		"$fogcloud/pk3/dn3/event/sdfdsfsdfsdfsd",
		"$fogcloud/pk1/dn1/all/#",
	}

	t.Run("match with priority", func(t *testing.T) {
		matchedTopic, _, _ := m.MatchWithAnonymousParams(dstTopics[3])
		assert.Equal(t, topics[0], matchedTopic)

		matchedTopic, _, _ = m1.MatchWithAnonymousParams(dstTopics[3])
		assert.Equal(t, topics[1], matchedTopic)

		matchedTopic, params, ok := m1.MatchWithAnonymousParams(dstTopics[7])
		assert.Equal(t, topics[8], matchedTopic)
		assert.Equal(t, true, ok)
		fmt.Println(params)
	})

	t.Run("srcPath split error", func(t *testing.T) {
		err = m.AddPath("")
		assert.Error(t, err)
		_, _, ok := m.Match("")
		assert.Equal(t, false, ok)
	})
}

func TestRouterPathMatch(t *testing.T) {
	m := NewRouterPathMatcher()
	var err error
	for i := range routerPaths {
		err = m.AddPath(routerPaths[i])
		assert.NoError(t, err)
	}

	testPaths := []string {
		"/groups/1/users",
		"/groups/windy/users/what",
		"/groups/3/users/3",
		"/projects/8/products/ting/1",
	}

	src, params, ok := m.Match(testPaths[0])
	assert.Equal(t, true, ok)
	assert.Equal(t, routerPaths[5], src)
	log.Printf("\npath: %s\nparams: %v\nmatched: %v", src, params, ok)

	src, params, ok = m.Match(testPaths[1])
	assert.Equal(t, true, ok)
	assert.Equal(t, routerPaths[3], src)
	log.Printf("\npath: %s\nparams: %v\nmatched: %v", src, params, ok)

	src, params, ok = m.Match(testPaths[2])
	assert.Equal(t, true, ok)
	assert.Equal(t, routerPaths[3], src)
	log.Printf("\npath: %s\nparams: %v\nmatched: %v", src, params, ok)

	src, params, ok = m.Match(testPaths[3])
	assert.Equal(t, true, ok)
	assert.Equal(t, routerPaths[2], src)
	log.Printf("\npath: %s\nparams: %v\nmatched: %v", src, params, ok)
}

func TestMatcher_AddPathWithPriority(t *testing.T) {
	m := NewRouterPathMatcher()
	var err error
	for i := range routerPaths {
		err = m.AddPathWithPriority(routerPaths[i], i)
		assert.NoError(t, err)
	}
	log.Print(m)
}

func BenchmarkTopicMatcher_Match(b *testing.B) {
	m := NewMqttTopicMatcher()
	for i := range topics {
		_ = m.AddPath(topics[i])
	}

	dstTopics := []string{
		"$fogcloud/pk1/device1/thing/event/property/post_reply",
		"$fogcloud/pk2/device2/thing/service/lightMode",
		"$fogcloud/pk3/device3/thing/event/alarm/post_reply",
		"$fogcloud/pk3/device3/thing/event/alarm/post_reply/err",
	}

	b.Run("Match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.Match(dstTopics[2])
		}
	})

	b.Run("Match with anonymous params", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.MatchWithAnonymousParams(dstTopics[2])
		}
	})
}



