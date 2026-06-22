package migrations

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/sqls"
)

func migrate_site_content_config() error {
	return sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		tx := ctx.Tx
		now := dates.NowTimestamp()

		if repositories.SysConfigRepository.GetByKey(tx, constants.SysConfigAboutPageConfig) == nil {
			value, err := jsons.ToStr(defaultAboutPageConfig())
			if err != nil {
				return err
			}
			name, desc := aboutPageConfigMetaByLanguage()
			if err := repositories.SysConfigRepository.Create(tx, &models.SysConfig{
				Key:         constants.SysConfigAboutPageConfig,
				Value:       value,
				Name:        name,
				Description: desc,
				CreateTime:  now,
				UpdateTime:  now,
			}); err != nil {
				return err
			}
		}

		if repositories.SysConfigRepository.GetByKey(tx, constants.SysConfigFooterLinks) == nil {
			value, err := jsons.ToStr(defaultFooterLinks())
			if err != nil {
				return err
			}
			name, desc := footerLinksMetaByLanguage()
			if err := repositories.SysConfigRepository.Create(tx, &models.SysConfig{
				Key:         constants.SysConfigFooterLinks,
				Value:       value,
				Name:        name,
				Description: desc,
				CreateTime:  now,
				UpdateTime:  now,
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

func defaultAboutPageConfig() dto.AboutPageConfig {
	return dto.AboutPageConfig{
		Content: defaultAboutPageContent(),
	}
}

func defaultAboutPageContent() dto.LocalizedText {
	return dto.LocalizedText{
		"en-US": "# About\n\nThis platform is used to collect product issues, track handling status, and maintain reusable knowledge entries.\n\n## What It Supports\n\n- Issue intake and triage by category, platform area, business scene, source, priority, severity, and owner.\n- Answer adoption and issue status updates for a lightweight problem-closing loop.\n- Knowledge Base articles for standard answers, operating notes, and reusable troubleshooting material.\n",
		"zh-CN": "# 关于\n\n本平台用于沉淀平台问题库和知识库，支撑问题收集、处理状态跟踪、答案采纳与知识条目维护。\n\n## 当前定位\n\n- 问题库：按分类、平台模块、业务场景、来源、优先级、严重程度和责任人记录问题。\n- 问题闭环：通过问答帖、采纳答案和处理状态更新完成基础闭环。\n- 知识库：通过文章维护标准答案、处理经验和可复用排查资料。\n",
	}
}

func defaultFooterLinks() []dto.FooterLink {
	return []dto.FooterLink{
		{
			Text: dto.LocalizedText{
				"en-US": "About",
				"zh-CN": "关于",
			},
			Url:             "/about",
			OpenInNewWindow: false,
			Visible:         true,
		},
		{
			Text: dto.LocalizedText{
				"en-US": "ICP Filing",
				"zh-CN": "ICP备案号",
			},
			Url:             "",
			OpenInNewWindow: true,
			Visible:         false,
		},
		{
			Text: dto.LocalizedText{
				"en-US": "Public Security Filing",
				"zh-CN": "公安网备信息",
			},
			Url:             "",
			OpenInNewWindow: true,
			Visible:         false,
		},
	}
}

func aboutPageConfigMetaByLanguage() (name, description string) {
	if config.Instance.Language == config.LanguageEnUS {
		return "About Page Config", "Configurable about page content with i18n text"
	}
	return "关于页配置", "可配置的关于页内容，支持中英文文本"
}

func footerLinksMetaByLanguage() (name, description string) {
	if config.Instance.Language == config.LanguageEnUS {
		return "Footer Links", "Footer links and filing information configuration"
	}
	return "底部链接", "底部链接与备案信息等配置"
}
