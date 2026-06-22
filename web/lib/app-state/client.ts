"use client"

import { apiFetch } from "@/lib/api/client"
import type { SiteConfig, UserSummary } from "@/lib/api/types"

export type ClientAppStateHydration = {
  config?: SiteConfig | null
  currentUser: UserSummary | null
  unreadMessageCount: number
}

export async function loadClientAppState(): Promise<ClientAppStateHydration> {
  const [config, currentUser] = await Promise.all([
    apiFetch<SiteConfig>("/api/config/configs").catch(() => undefined),
    apiFetch<UserSummary | null>("/api/user/current").catch(() => null),
  ])

  return {
    config,
    currentUser,
    unreadMessageCount: 0,
  }
}
