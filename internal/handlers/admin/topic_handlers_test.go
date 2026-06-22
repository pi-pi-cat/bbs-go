package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/search"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func TestTopicUpdateIssueStatusChangesIssueStatus(t *testing.T) {
	db := setupAdminTopicTestDB(t)
	now := time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC).UnixMilli()
	topic := &models.Topic{
		Model:       models.Model{Id: 1},
		Type:        constants.TopicTypeQA,
		CategoryId:  1,
		UserId:      1,
		Title:       "问题状态流转",
		Content:     "待处理问题",
		ContentType: constants.ContentTypeMarkdown,
		QaStatus:    constants.QaStatusUnsolved,
		IssueStatus: constants.IssueStatusOpen,
		Status:      constants.StatusOk,
		CreateTime:  now,
	}
	if err := db.Create(topic).Error; err != nil {
		t.Fatalf("create topic: %v", err)
	}

	postTopicUpdateIssueStatus(t, "id=1&issueStatus=processing")

	var saved models.Topic
	if err := db.First(&saved, "id = ?", topic.Id).Error; err != nil {
		t.Fatalf("load topic: %v", err)
	}
	if saved.IssueStatus != constants.IssueStatusProcessing {
		t.Fatalf("expected issue status processing, got %q", saved.IssueStatus)
	}
}

func setupAdminTopicTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:admin_topic_test_%d?mode=memory&cache=shared&_fk=1", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	sqls.SetDB(db)
	config.Instance = &config.Config{
		Language: config.DefaultLanguage,
		Search: config.SearchConfig{
			IndexPath: filepath.Join(t.TempDir(), "index"),
		},
	}
	search.Init()
	if err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Topic{},
		&models.Tag{},
		&models.TopicTag{},
	); err != nil {
		t.Fatalf("auto migrate topic dependencies: %v", err)
	}
	return db
}

func postTopicUpdateIssueStatus(t *testing.T, body string) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(http.MethodPost, "/api/admin/topic/update_issue_status", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx.Request = req

	TopicUpdateIssueStatus(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode response %q: %v", recorder.Body.String(), err)
	}
	if !result.Success {
		t.Fatalf("expected success response, got %s", recorder.Body.String())
	}
}
