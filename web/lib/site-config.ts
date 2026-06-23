import type { SiteConfig, SiteNav } from "@/lib/api/types"
import { normalizeLocale } from "@/lib/i18n"

const removedSiteNavUrls = new Set(["/tasks"])
const removedFooterLinkUrls = new Set(["/about", "/links"])

function normalizeSiteNavTitle(nav: SiteNav, language?: string) {
  if (nav.url === "/topics") {
    return normalizeLocale(language) === "zh-CN" ? "问答" : "Q&A"
  }
  if (nav.url === "/articles") {
    return normalizeLocale(language) === "zh-CN" ? "知识库" : "Knowledge Base"
  }
  return nav.title
}

function normalizeSiteNav(nav: SiteNav, language?: string): SiteNav {
  return {
    ...nav,
    title: normalizeSiteNavTitle(nav, language),
    children: nav.children
      ?.filter((child) => !removedSiteNavUrls.has(child.url))
      .map((child) => normalizeSiteNav(child, language)),
  }
}

export function normalizeSiteConfig(config: SiteConfig | null | undefined) {
  if (!config) return config ?? null

  return {
    ...config,
    siteNavs: config.siteNavs
      ?.filter((nav) => !removedSiteNavUrls.has(nav.url))
      .map((nav) => normalizeSiteNav(nav, config.language)),
    footerLinks: config.footerLinks?.filter(
      (link) => !removedFooterLinkUrls.has(link.url || "")
    ),
  }
}
