import { cache } from "react"

import { getSiteConfig } from "@/lib/api/site"
import type { SiteConfig, UserSummary } from "@/lib/api/types"
import { getSessionUser } from "@/lib/auth/session"
import { createT, type Locale, type TFunction } from "@/lib/i18n"
import { getServerLocale } from "@/lib/i18n/server"
import { normalizeSiteConfig } from "@/lib/site-config"

export type AppState = {
  config: SiteConfig | null
  currentUser: UserSummary | null
  locale: Locale
  unreadMessageCount: number
  t: TFunction
  isLogin: boolean
}

export const getAppState = cache(async (): Promise<AppState> => {
  const [config, currentUser] = await Promise.all([
    getSiteConfig().catch(() => null),
    getSessionUser().catch(() => null),
  ])
  const locale = await getServerLocale(config?.language)

  return {
    config: normalizeSiteConfig(config),
    currentUser,
    locale,
    unreadMessageCount: 0,
    t: createT(locale),
    isLogin: Boolean(currentUser),
  }
})
