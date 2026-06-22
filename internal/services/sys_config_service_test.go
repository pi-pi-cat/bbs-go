package services

import "testing"

func TestParseModulesConfig_BackfillsQaFromTopicForLegacyConfig(t *testing.T) {
	cfg := parseModulesConfig(`{"tweet":true,"topic":true,"article":false}`)

	if !cfg.QA {
		t.Fatalf("expected legacy config without qa to keep QA enabled when topic is enabled")
	}
}

func TestParseModulesConfig_RespectsExplicitQaSwitch(t *testing.T) {
	cfg := parseModulesConfig(`{"tweet":true,"topic":true,"qa":false,"article":true}`)

	if cfg.QA {
		t.Fatalf("expected explicit qa=false to disable QA independently from topic")
	}
}

func TestDefaultNotificationTypeKeysDropRemovedCommunityGrowthEvents(t *testing.T) {
	keys := defaultNotificationTypeKeys()
	for _, removed := range []string{"topicFavorite", "userLevelUp", "userBadgeGrant"} {
		if containsString(keys, removed) {
			t.Fatalf("expected notification type %q to be removed, got %#v", removed, keys)
		}
	}
	for _, expected := range []string{"topicComment", "commentReply", "topicLike", "topicRecommend", "topicDelete", "articleComment", "qaAnswerAccepted"} {
		if !containsString(keys, expected) {
			t.Fatalf("expected notification type %q to remain, got %#v", expected, keys)
		}
	}
}
