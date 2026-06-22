package services

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/idcodec"
	"bbs-go/internal/pkg/msg"
	"testing"

	"github.com/mlogclub/simple/sqls"
)

func TestBuildEmailNoticeSubjectAvoidsBlankSitePrefix(t *testing.T) {
	got := MessageService.buildEmailNoticeSubject("", "你的话题被设为推荐")
	if got != "你的话题被设为推荐" {
		t.Fatalf("expected subject without blank site prefix, got %q", got)
	}
}

func TestBuildEmailNoticeSubjectIncludesSiteTitle(t *testing.T) {
	got := MessageService.buildEmailNoticeSubject("BBS-GO", "你的话题被设为推荐")
	if got != "BBS-GO - 你的话题被设为推荐" {
		t.Fatalf("expected subject with site title, got %q", got)
	}
}

func TestBuildEmailNoticeContentFallsBackToNoticeTitle(t *testing.T) {
	got := MessageService.buildEmailNoticeContent("", "你的话题被设为推荐")
	if got != "你的话题被设为推荐" {
		t.Fatalf("expected notice title fallback, got %q", got)
	}
}

func TestBuildEmailNoticeDetailURLUsesTopicForRecommend(t *testing.T) {
	idcodec.Init(1)

	got := MessageService.buildEmailNoticeDetailURL(&models.Message{
		Type:      int(msg.TypeTopicRecommend),
		ExtraData: `{"topicId":123}`,
	})

	if got != "/topic/"+idcodec.Encode(123) {
		t.Fatalf("expected topic detail url, got %q", got)
	}
}

func TestBuildEmailNoticeDetailURLFallbackDoesNotPointToRemovedMessagesPage(t *testing.T) {
	setupMessageServiceURLTestDB(t)

	got := MessageService.buildEmailNoticeDetailURL(&models.Message{Type: 999})
	if got != "/user/profile" {
		t.Fatalf("expected unknown message fallback to profile, got %q", got)
	}
}

func setupMessageServiceURLTestDB(t *testing.T) {
	t.Helper()

	db := setupTestDB(t)
	if err := db.AutoMigrate(&models.SysConfig{}); err != nil {
		t.Fatalf("auto migrate sys config: %v", err)
	}
	cache.SysConfigCache.Invalidate(constants.SysConfigBaseURL)
	t.Cleanup(func() { cache.SysConfigCache.Invalidate(constants.SysConfigBaseURL) })
	if err := db.Create(&models.SysConfig{Key: constants.SysConfigBaseURL, Value: "/"}).Error; err != nil {
		t.Fatalf("create base URL config: %v", err)
	}
	_ = sqls.DB()
}
