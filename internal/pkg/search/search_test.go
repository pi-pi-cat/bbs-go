package search

import "testing"

func setupTestIndex(t *testing.T) {
	t.Helper()
	idx := newIndex(t.TempDir())
	if idx == nil {
		t.Fatal("expected test index")
	}
	index = idx
	t.Cleanup(func() {
		_ = idx.Close()
	})
}

func TestSearchTopicScopesResultsToTopics(t *testing.T) {
	setupTestIndex(t)

	mustIndex(t, searchDocID(EntityTypeTopic, 1), &TopicDocument{
		Type:       EntityTypeTopic,
		Id:         1,
		CategoryId: 10,
		UserId:     100,
		Nickname:   "Ada",
		Title:      "Golang search topic",
		Content:    "Topic content about bleve search.",
		Status:     0,
		CreateTime: 1000,
	})
	mustIndex(t, searchDocID(EntityTypeArticle, 2), &ArticleDocument{
		Type:       EntityTypeArticle,
		Id:         2,
		UserId:     100,
		Nickname:   "Ada",
		Title:      "Golang search article",
		Summary:    "Article summary about bleve search.",
		Content:    "Article body.",
		Status:     0,
		CreateTime: 1000,
	})

	docs, _, err := SearchTopic("Golang", 0, nil, 0, 1, 20)
	if err != nil {
		t.Fatalf("SearchTopic returned error: %v", err)
	}
	if len(docs) != 1 {
		t.Fatalf("expected one topic result, got %d", len(docs))
	}
	if docs[0].Id != 1 {
		t.Fatalf("expected topic id 1, got %d", docs[0].Id)
	}
}

func TestSearchArticleFindsArticleFields(t *testing.T) {
	setupTestIndex(t)

	mustIndex(t, searchDocID(EntityTypeArticle, 11), &ArticleDocument{
		Type:       EntityTypeArticle,
		Id:         11,
		UserId:     101,
		Nickname:   "Grace",
		Title:      "React router article",
		Summary:    "A compact guide for framework mode.",
		Content:    "Article content about loaders and actions.",
		Status:     0,
		CreateTime: 1000,
	})

	docs, _, err := SearchArticle("framework", 0, 1, 20)
	if err != nil {
		t.Fatalf("SearchArticle returned error: %v", err)
	}
	if len(docs) != 1 {
		t.Fatalf("expected one article result, got %d", len(docs))
	}
	if docs[0].Id != 11 {
		t.Fatalf("expected article id 11, got %d", docs[0].Id)
	}
	if docs[0].Summary == "" {
		t.Fatal("expected article summary to be returned")
	}
}

func TestSearchAllReturnsOnlyTopicAndArticlePreview(t *testing.T) {
	setupTestIndex(t)

	mustIndex(t, searchDocID(EntityTypeTopic, 1), &TopicDocument{
		Type:       EntityTypeTopic,
		Id:         1,
		Title:      "Unified search topic",
		Content:    "Topic content.",
		CreateTime: 1000,
	})
	mustIndex(t, searchDocID(EntityTypeArticle, 2), &ArticleDocument{
		Type:       EntityTypeArticle,
		Id:         2,
		Title:      "Unified search article",
		Summary:    "Article summary.",
		Content:    "Article content.",
		CreateTime: 1000,
	})

	result, err := SearchAll("Unified", 5)
	if err != nil {
		t.Fatalf("SearchAll returned error: %v", err)
	}
	if len(result.Topics) != 1 || len(result.Articles) != 1 {
		t.Fatalf("expected grouped topic/article preview, got topics=%d articles=%d", len(result.Topics), len(result.Articles))
	}
}

func mustIndex(t *testing.T, id string, doc any) {
	t.Helper()
	if err := index.Index(id, doc); err != nil {
		t.Fatalf("index %s: %v", id, err)
	}
}
