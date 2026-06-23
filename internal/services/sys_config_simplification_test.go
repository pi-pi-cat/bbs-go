package services

import (
	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/config"
	"testing"
)

func TestNormalizeSiteNavsRemovesTasksAndMapsContentLabels(t *testing.T) {
	original := config.Instance
	t.Cleanup(func() { config.Instance = original })
	config.Instance = &config.Config{Language: config.LanguageZhCN}

	navs := []dto.ActionLink{
		{Title: "任务", Url: "/tasks"},
		{Title: "话题", Url: "/topics"},
		{Title: "文章", Url: "/articles"},
		{
			Title: "更多",
			Url:   "/more",
			Children: []dto.ActionLink{
				{Title: "任务", Url: "/tasks"},
				{Title: "文章", Url: "/articles"},
			},
		},
	}

	got := normalizeSiteNavs(navs)

	if len(got) != 3 {
		t.Fatalf("expected tasks nav to be removed, got %#v", got)
	}
	if got[0].Url != "/topics" || got[0].Title != "问答" {
		t.Fatalf("expected topics nav to be mapped to Q&A, got %#v", got[0])
	}
	if got[1].Url != "/articles" || got[1].Title != "知识库" {
		t.Fatalf("expected articles nav to be mapped to knowledge base, got %#v", got[1])
	}
	if len(got[2].Children) != 1 || got[2].Children[0].Url != "/articles" || got[2].Children[0].Title != "知识库" {
		t.Fatalf("expected child navs to be filtered and mapped, got %#v", got[2].Children)
	}
}

func TestNormalizeSiteNavsUsesEnglishLabelsForEnglishSites(t *testing.T) {
	original := config.Instance
	t.Cleanup(func() { config.Instance = original })
	config.Instance = &config.Config{Language: config.LanguageEnUS}

	got := normalizeSiteNavs([]dto.ActionLink{
		{Title: "Topics", Url: "/topics"},
		{Title: "Articles", Url: "/articles"},
	})

	if len(got) != 2 {
		t.Fatalf("expected content navs to remain, got %#v", got)
	}
	if got[0].Title != "Q&A" || got[1].Title != "Knowledge Base" {
		t.Fatalf("expected English content labels, got %#v", got)
	}
}

func TestFilterFooterLinksRemovesAboutAndLinks(t *testing.T) {
	links := []dto.FooterLink{
		{Url: "/about", Text: dto.LocalizedText{"zh-CN": "关于"}, Visible: true},
		{Url: "/links", Text: dto.LocalizedText{"zh-CN": "友情链接"}, Visible: true},
		{Url: "", Text: dto.LocalizedText{"zh-CN": "ICP备案号"}, Visible: true},
	}

	got := filterFooterLinks(links)

	if len(got) != 1 {
		t.Fatalf("expected only filing link to remain, got %#v", got)
	}
	if got[0].Url != "" || got[0].Text["zh-CN"] != "ICP备案号" {
		t.Fatalf("unexpected remaining footer link: %#v", got[0])
	}
}
