import { UserProfileCard } from "@/components/user/user-profile-card"
import { UserCenterSidebar } from "@/components/user/user-sidebar"
import type { UserSummary } from "@/lib/api/types"
import type { TFunction } from "@/lib/i18n"

export function UserCenterShell({
  user,
  currentUser,
  t,
  children,
}: {
  user: UserSummary
  currentUser?: UserSummary | null
  t: TFunction
  children: React.ReactNode
}) {
  return (
    <section className="main">
      <div className="container">
        <UserProfileCard user={user} currentUser={currentUser} />
      </div>
      <div className="container main-container right-main side-size-360">
        <UserCenterSidebar user={user} currentUser={currentUser} t={t} />
        <div className="right-container">{children}</div>
      </div>
    </section>
  )
}
