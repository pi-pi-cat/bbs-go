import Link from "@/components/common/link"
import { ChevronRight } from "lucide-react"

import { WidgetCard } from "@/components/common/widget-card"
import { UserCenterOperations } from "@/components/user/user-center-operations"
import type { UserSummary } from "@/lib/api/types"
import type { TFunction } from "@/lib/i18n"

export function MyProfileCard({
  user,
  currentUser,
  t,
}: {
  user: UserSummary
  currentUser?: UserSummary | null
  t: TFunction
}) {
  const canEdit = currentUser?.id === user.id
  return (
    <WidgetCard
      title={t("component.myProfile.title")}
      actions={
        canEdit ? (
          <Link href="/user/profile" className="inline-flex items-center gap-1">
            {t("component.myProfile.editProfile")}
            <ChevronRight className="h-4 w-4" />
          </Link>
        ) : null
      }
    >
      <div className="stable">
        <div className="str">
          <div className="slabel">{t("component.myProfile.nickname")}</div>
          <div className="svalue">{user.nickname}</div>
        </div>
        <div className="str">
          <div className="slabel">{t("component.myProfile.description")}</div>
          <div className="svalue">{user.description}</div>
        </div>
        {user.homePage ? (
          <div className="str">
            <div className="slabel">{t("component.myProfile.homePage")}</div>
            <div className="svalue">
              <a href={user.homePage} target="_blank" rel="nofollow noreferrer">
                {user.homePage}
              </a>
            </div>
          </div>
        ) : null}
      </div>
    </WidgetCard>
  )
}

export function UserCenterSidebar({
  user,
  currentUser,
  t,
}: {
  user: UserSummary
  currentUser?: UserSummary | null
  t: TFunction
}) {
  return (
    <div className="left-container space-y-4">
      <MyProfileCard user={user} currentUser={currentUser} t={t} />
      <UserCenterOperations user={user} currentUser={currentUser} />
    </div>
  )
}
