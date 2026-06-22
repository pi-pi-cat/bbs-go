"use client"

import * as React from "react"

import { Input } from "@/components/ui/input"
import { useI18n } from "@/lib/i18n/provider"

export type IssueFieldValues = {
  platformArea: string
  businessScene: string
  issueSource: string
  issuePriority: string
  issueSeverity: string
  issueOwner: string
}

const issueFieldNames: Array<keyof IssueFieldValues> = [
  "platformArea",
  "businessScene",
  "issueSource",
  "issuePriority",
  "issueSeverity",
  "issueOwner",
]

export function IssueFields({
  value,
  onChange,
}: {
  value: IssueFieldValues
  onChange: (next: Partial<IssueFieldValues>) => void
}) {
  const { t } = useI18n()

  return (
    <div className="field grid gap-3 rounded-md border bg-muted/20 p-3 sm:grid-cols-2">
      {issueFieldNames.map((name) => (
        <Input
          key={name}
          value={value[name]}
          placeholder={t(`pages.topic.create.issueFields.${name}`)}
          onChange={(event) => onChange({ [name]: event.currentTarget.value })}
        />
      ))}
    </div>
  )
}
