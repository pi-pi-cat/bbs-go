package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/req"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/search"
	"bbs-go/internal/repositories"
	"testing"
	"time"

	"github.com/mlogclub/simple/sqls"
)

func setupTopicPublishServiceTestDB(t *testing.T) {
	t.Helper()
	config.Instance = &config.Config{Language: config.DefaultLanguage}
	search.Init()
	db := setupTestDB(t)
	if err := db.AutoMigrate(&models.Category{}, &models.Topic{}, &models.Tag{}, &models.TopicTag{}); err != nil {
		t.Fatalf("auto migrate topic publish: %v", err)
	}
}

func mustCreateCategory(t *testing.T, category *models.Category) *models.Category {
	t.Helper()
	if category.Status == 0 {
		category.Status = constants.StatusOk
	}
	if category.Type == "" {
		category.Type = constants.CategoryTypeQA
	}
	if err := repositories.CategoryRepository.Create(sqls.DB(), category); err != nil {
		t.Fatalf("create category: %v", err)
	}
	return category
}

func TestTopicPublishServicePublishStoresIssueBusinessFields(t *testing.T) {
	setupTopicPublishServiceTestDB(t)
	now := time.Date(2026, 6, 22, 10, 0, 0, 0, time.UTC).UnixMilli()
	user := mustCreateUser(t, now)
	category := mustCreateCategory(t, &models.Category{
		Name:       "仿真平台问题",
		Type:       constants.CategoryTypeQA,
		CreateTime: now,
	})

	topic, err := TopicPublishService.Publish(user.Id, req.CreateTopicReq{
		Type:          constants.TopicTypeQA,
		CategoryId:    category.Id,
		Title:         "模型加载失败",
		Content:       "导入模型后平台提示加载失败。",
		ContentType:   constants.ContentTypeMarkdown,
		PlatformArea:  "仿真建模平台",
		BusinessScene: "模型导入",
		IssueSource:   "用户反馈",
		IssuePriority: "high",
		IssueSeverity: "major",
		IssueOwner:    "研发支持组",
	})
	if err != nil {
		t.Fatalf("publish topic: %v", err)
	}

	got := TopicService.Get(topic.Id)
	if got == nil {
		t.Fatalf("expected topic to be saved")
	}
	if got.PlatformArea != "仿真建模平台" {
		t.Fatalf("expected platform area saved, got %q", got.PlatformArea)
	}
	if got.BusinessScene != "模型导入" {
		t.Fatalf("expected business scene saved, got %q", got.BusinessScene)
	}
	if got.IssueSource != "用户反馈" {
		t.Fatalf("expected issue source saved, got %q", got.IssueSource)
	}
	if got.IssuePriority != "high" {
		t.Fatalf("expected issue priority saved, got %q", got.IssuePriority)
	}
	if got.IssueSeverity != "major" {
		t.Fatalf("expected issue severity saved, got %q", got.IssueSeverity)
	}
	if got.IssueOwner != "研发支持组" {
		t.Fatalf("expected issue owner saved, got %q", got.IssueOwner)
	}
	if got.IssueStatus != constants.IssueStatusOpen {
		t.Fatalf("expected new issue status open, got %q", got.IssueStatus)
	}
}

func TestTopicServiceUpdateIssueStatusUpdatesOnlyQaIssues(t *testing.T) {
	setupTopicPublishServiceTestDB(t)
	now := time.Date(2026, 6, 22, 11, 0, 0, 0, time.UTC).UnixMilli()
	user := mustCreateUser(t, now)
	category := mustCreateCategory(t, &models.Category{
		Name:       "问题流转",
		Type:       constants.CategoryTypeQA,
		CreateTime: now,
	})

	issue, err := TopicPublishService.Publish(user.Id, req.CreateTopicReq{
		Type:        constants.TopicTypeQA,
		CategoryId:  category.Id,
		Title:       "仿真任务卡住",
		Content:     "运行到 80% 后无响应。",
		ContentType: constants.ContentTypeMarkdown,
	})
	if err != nil {
		t.Fatalf("publish issue: %v", err)
	}

	if err := TopicService.UpdateIssueStatus(issue.Id, constants.IssueStatusProcessing); err != nil {
		t.Fatalf("update issue status: %v", err)
	}
	got := TopicService.Get(issue.Id)
	if got.IssueStatus != constants.IssueStatusProcessing {
		t.Fatalf("expected issue status processing, got %q", got.IssueStatus)
	}

	normalCategory := mustCreateCategory(t, &models.Category{
		Name:       "普通讨论",
		Type:       constants.CategoryTypeNormal,
		CreateTime: now,
	})
	topic, err := TopicPublishService.Publish(user.Id, req.CreateTopicReq{
		Type:        constants.TopicTypeTopic,
		CategoryId:  normalCategory.Id,
		Title:       "普通讨论",
		Content:     "不是问题单。",
		ContentType: constants.ContentTypeMarkdown,
	})
	if err != nil {
		t.Fatalf("publish topic: %v", err)
	}

	if err := TopicService.UpdateIssueStatus(topic.Id, constants.IssueStatusProcessing); err == nil {
		t.Fatalf("expected normal topic to reject issue status update")
	}
}
