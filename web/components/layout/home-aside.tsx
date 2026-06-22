"use client"

import * as React from "react"

import { useAppState } from "@/components/app/app-provider"
import { useI18n } from "@/lib/i18n/provider"

function WidgetCard({
  title,
  actions,
  children,
}: {
  title?: string
  actions?: React.ReactNode
  children: React.ReactNode
}) {
  return (
    <section className="rounded-md bg-background px-3 py-1">
      {title || actions ? (
        <div className="flex items-center justify-between border-b py-2 text-base font-medium">
          <span>{title}</span>
          {actions ? (
            <div className="shrink-0 text-sm font-normal">{actions}</div>
          ) : null}
        </div>
      ) : null}
      <div className="py-2 break-all">{children}</div>
    </section>
  )
}

function SiteNotice({ title, content }: { title: string; content?: string }) {
  if (!content) {
    return null
  }

  return (
    <WidgetCard title={title}>
      <div
        className="prose prose-sm max-w-none text-sm text-muted-foreground"
        dangerouslySetInnerHTML={{ __html: content }}
      />
    </WidgetCard>
  )
}

export function HomeAside() {
  const { config } = useAppState()
  const { t } = useI18n()

  return (
    <React.Fragment>
      <SiteNotice
        title={t("component.siteNotice.title")}
        content={config?.siteNotification}
      />
    </React.Fragment>
  )
}
