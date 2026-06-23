package migrations

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/config"
	"strings"
	"testing"
)

func TestSeedForLanguageDoesNotExposeRemovedCommunityGrowthNavs(t *testing.T) {
	original := config.Instance
	t.Cleanup(func() { config.Instance = original })

	expectedTitles := map[config.Language]map[string]string{
		config.LanguageEnUS: {"/topics": "Q&A", "/articles": "Knowledge Base"},
		config.LanguageZhCN: {"/topics": "问答", "/articles": "知识库"},
	}

	for _, lang := range []config.Language{config.LanguageEnUS, config.LanguageZhCN} {
		config.Instance = &config.Config{Language: lang}
		seed := seedForLanguage()
		navs := findSysConfigValue[[]map[string]string](t, seed.SysConfigs, constants.SysConfigSiteNavs)

		for _, nav := range navs {
			if nav["url"] == "/tasks" {
				t.Fatalf("default %s nav should not expose removed tasks page: %#v", lang, navs)
			}
			if expected, ok := expectedTitles[lang][nav["url"]]; ok && nav["title"] != expected {
				t.Fatalf("default %s nav %s should be titled %q, got %#v", lang, nav["url"], expected, navs)
			}
		}
	}
}

func TestDefaultFooterLinksDoesNotExposeRemovedLinksPage(t *testing.T) {
	links := defaultFooterLinks()
	for _, link := range links {
		if link.Url == "/about" || link.Url == "/links" {
			t.Fatalf("default footer links should not expose removed page %s: %#v", link.Url, links)
		}
	}
}

func TestDefaultAboutContentUsesIssueKnowledgeBasePositioning(t *testing.T) {
	contents := defaultAboutPageContent()
	for lang, content := range contents {
		for _, removed := range []string{"BBS-GO", "open source community", "开源社区"} {
			if strings.Contains(content, removed) {
				t.Fatalf("default about content for %s should not contain legacy community copy %q: %s", lang, removed, content)
			}
		}
	}
	if !strings.Contains(contents["zh-CN"], "问题库") || !strings.Contains(contents["zh-CN"], "知识库") {
		t.Fatalf("zh-CN about content should describe issue and knowledge base positioning: %s", contents["zh-CN"])
	}
	if !strings.Contains(contents["en-US"], "Issue") || !strings.Contains(contents["en-US"], "Knowledge") {
		t.Fatalf("en-US about content should describe issue and knowledge base positioning: %s", contents["en-US"])
	}
}

func TestNotificationTypeDefaultsDropRemovedCommunityGrowthEvents(t *testing.T) {
	defaults := defaultNotificationTypeConfigs()
	for _, removed := range []string{"topicFavorite", "userLevelUp", "userBadgeGrant"} {
		if _, ok := defaults[removed]; ok {
			t.Fatalf("default notification type %q should be removed, got %#v", removed, defaults)
		}
	}
	for _, expected := range []string{"topicComment", "commentReply", "topicLike", "topicRecommend", "topicDelete", "articleComment", "qaAnswerAccepted"} {
		if _, ok := defaults[expected]; !ok {
			t.Fatalf("default notification type %q should remain, got %#v", expected, defaults)
		}
	}
}

func findSysConfigValue[T any](t *testing.T, items []sysConfigSeed, key string) T {
	t.Helper()

	for _, item := range items {
		if item.Key != key {
			continue
		}
		value, ok := item.Value.(T)
		if !ok {
			t.Fatalf("sys config %q has unexpected type %T", key, item.Value)
		}
		return value
	}

	var zero T
	t.Fatalf("sys config %q not found", key)
	return zero
}
