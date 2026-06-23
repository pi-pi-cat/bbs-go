import assert from "node:assert/strict"
import { existsSync, readFileSync } from "node:fs"
import { resolve } from "node:path"

const webRoot = resolve(import.meta.dirname, "..")
const repoRoot = resolve(webRoot, "..")

function source(path) {
  return readFileSync(resolve(webRoot, path), "utf8")
}

function repoSource(path) {
  return readFileSync(resolve(repoRoot, path), "utf8")
}

const siteHeader = source("components/layout/site-header.tsx")

for (const forbidden of [
  "/user/messages",
  "/user/favorites",
  "useUnreadMessageCount",
  "MsgNotice",
  "common.header.tasks",
  "common.header.favorites",
]) {
  assert.equal(
    siteHeader.includes(forbidden),
    false,
    `site header should not expose community utility entry: ${forbidden}`
  )
}

const siteConfigNormalizer = source("lib/site-config.ts")
for (const expected of ["normalizeSiteConfig", "removedSiteNavUrls", "removedFooterLinkUrls"]) {
  assert.equal(
    siteConfigNormalizer.includes(expected),
    true,
    `site config should normalize removed legacy configuration: ${expected}`
  )
}

for (const forbidden of ['href="/tasks"', 'href: "/tasks"', '"/tasks"']) {
  assert.equal(
    siteHeader.includes(forbidden),
    false,
    `site header should not render removed tasks navigation: ${forbidden}`
  )
}

const siteFooter = source("components/layout/site-footer.tsx")
for (const forbidden of ["Powered by", "BBS-GO", "https://bbs-go.com"]) {
  assert.equal(
    siteFooter.includes(forbidden),
    false,
    `site footer should not expose removed product attribution: ${forbidden}`
  )
}

for (const expected of ['"/about"', '"/links"']) {
  assert.equal(
    siteConfigNormalizer.includes(expected),
    true,
    `site config should hide removed footer link from existing config: ${expected}`
  )
}

for (const file of [
  "app/root.tsx",
  "lib/app-state/client.ts",
  "lib/app-state/server.ts",
]) {
  assert.equal(
    source(file).includes("normalizeSiteConfig"),
    true,
    `${file} should normalize legacy config before rendering`
  )
}

const dashboardUserMenu = source("components/dashboard/nav-user.tsx")
assert.equal(
  dashboardUserMenu.includes("/user/messages"),
  false,
  "dashboard user menu should not expose messages"
)

const profileCard = source("components/user/user-profile-card.tsx")
for (const forbidden of ["FollowButton", "level-badge", "badge-icon", "Medal"]) {
  assert.equal(
    profileCard.includes(forbidden),
    false,
    `user profile card should not expose social growth UI: ${forbidden}`
  )
}

const userSidebar = source("components/user/user-sidebar.tsx")
for (const forbidden of [
  "UserCountsCard",
  "UserBadgesWidget",
  "FollowWidget",
  "UserFollowList",
  "component.myCounts",
  "component.userBadges",
  "component.fansWidget",
  "component.followWidget",
]) {
  assert.equal(
    userSidebar.includes(forbidden),
    false,
    `user center sidebar should not expose community growth UI: ${forbidden}`
  )
}

const userInfo = source("components/user/user-info.tsx")
for (const forbidden of [
  "component.userInfo.level",
  "component.userInfo.score",
  "user.level",
  "user.score",
]) {
  assert.equal(
    userInfo.includes(forbidden),
    false,
    `topic author card should not expose growth metrics: ${forbidden}`
  )
}

const searchRoute = source("app/routes/search.tsx")
for (const forbidden of [
  "SearchUserList",
  "SearchUser",
  '"user"',
  "/api/search/user",
  "pages.search.tabs.user",
]) {
  assert.equal(
    searchRoute.includes(forbidden),
    false,
    `search route should focus on issue and knowledge content, not users: ${forbidden}`
  )
}

const dashboardTopicsRoute = source("app/routes/dashboard.topics.tsx")
assert.match(
  dashboardTopicsRoute,
  /createAdminInitialFilters\(\{\s*status:\s*0,\s*type:\s*2\s*\},\s*20\)/,
  "dashboard issue list should default to normal QA issues"
)

for (const removedRoute of [
  "app/routes/tasks.tsx",
  "app/routes/links.tsx",
  "app/routes/user.favorites.tsx",
  "app/routes/user.messages.tsx",
  "app/routes/user.scores.tsx",
  "app/routes/user_.$userId.badges.tsx",
  "app/routes/user_.$userId.fans.tsx",
  "app/routes/user_.$userId.followed.tsx",
]) {
  assert.equal(
    existsSync(resolve(webRoot, removedRoute)),
    false,
    `${removedRoute} should be removed from the simplified product routes`
  )
}

for (const removedFile of [
  "app/route-helpers/private-center.tsx",
  "lib/api/tasks.ts",
  "components/tasks/task-widgets.tsx",
  "components/tasks/tasks-page.tsx",
  "components/user/private-user-center-page.tsx",
  "components/search/search-user-list.tsx",
  "components/user/user-follow-list.tsx",
  "components/user/user-lists.tsx",
  "components/user/follow-button.tsx",
  "components/user/user-badges.tsx",
]) {
  assert.equal(
    existsSync(resolve(webRoot, removedFile)),
    false,
    `${removedFile} should be removed after dropping community growth routes`
  )
}

for (const removedDashboardRoute of [
  "app/routes/dashboard.badges.tsx",
  "app/routes/dashboard.levels.tsx",
  "app/routes/dashboard.links.tsx",
  "app/routes/dashboard.tasks.tsx",
  "app/routes/dashboard.user-badges.tsx",
  "app/routes/dashboard.user-exp-logs.tsx",
  "app/routes/dashboard.user-task-logs.tsx",
]) {
  assert.equal(
    existsSync(resolve(webRoot, removedDashboardRoute)),
    false,
    `${removedDashboardRoute} should be removed from dashboard routes`
  )
}

const dashboardLayout = source("app/routes/dashboard.tsx")
for (const forbidden of [
  "/dashboard/links",
  "/dashboard/badges",
  "/dashboard/levels",
  "/dashboard/tasks",
  "/dashboard/user-badges",
  "/dashboard/user-exp-logs",
  "/dashboard/user-task-logs",
  "dashboard.nav.growth",
  "dashboard.nav.userBadges",
  "dashboard.nav.userExpLogs",
  "dashboard.nav.userTaskLogs",
]) {
  assert.equal(
    dashboardLayout.includes(forbidden),
    false,
    `dashboard layout should not expose removed dashboard route: ${forbidden}`
  )
}

const dashboardSettings = source("app/routes/dashboard.settings.tsx")
for (const forbidden of [
  '"topicFavorite"',
  '"userLevelUp"',
  '"userBadgeGrant"',
  "enableQaBounty",
  "qaBountyMin",
  "qaBountyMax",
  "qaBountyRequired",
  "sectionQaBounty",
]) {
  assert.equal(
    dashboardSettings.includes(forbidden),
    false,
    `settings page should not expose removed notification or bounty feature: ${forbidden}`
  )
}

const dashboardUsers = source("app/routes/dashboard.users.tsx")
assert.equal(
  dashboardUsers.includes('"score"'),
  false,
  "dashboard users should not expose score column"
)

const userProfilePage = source("components/user/user-profile-client-page.tsx")
for (const forbidden of [
  "UserBadgesClientPage",
  "UserFansClientPage",
  "UserFollowedClientPage",
  "UserFollowClientPage",
  "UserFollowList",
  "/api/badge/badges",
  "/api/fans/",
]) {
  assert.equal(
    userProfilePage.includes(forbidden),
    false,
    `user profile pages should not expose community growth views: ${forbidden}`
  )
}

const miscApi = source("lib/api/misc.ts")
for (const forbidden of ["/api/link/list", "/api/link/top_links", "FriendLink"]) {
  assert.equal(
    miscApi.includes(forbidden),
    false,
    `misc api should not expose removed link API: ${forbidden}`
  )
}

for (const clientSource of [
  "lib/api/users.ts",
  "lib/actions/user.ts",
  "lib/app-state/client.ts",
  "components/topic/topic-action-context.tsx",
  "components/topic/topic-detail-actions.tsx",
  "components/topic/topic-side-action-bar.tsx",
]) {
  const text = source(clientSource)
  for (const forbidden of [
    "/api/search/user",
    "/api/user/favorites",
    "/api/user/messages",
    "/api/user/msg_recent",
    "/api/user/score_logs",
    "/api/user/score/rank",
    "/api/fans/",
    "/api/badge/badges",
    "/api/favorite/",
    "/api/task/",
    "/api/checkin/",
    "/api/link/",
  ]) {
    assert.equal(
      text.includes(forbidden),
      false,
      `${clientSource} should not call removed community growth API: ${forbidden}`
    )
  }
}

const robots = source("public/robots.txt")
for (const forbidden of [
  "/tasks/create",
  "/tasks/edit",
  "/user/favorites",
  "/user/messages",
  "/user/scores",
]) {
  assert.equal(
    robots.includes(forbidden),
    false,
    `robots.txt should not carry removed page rule: ${forbidden}`
  )
}

for (const [file, forbiddenList] of [
  [
    "migrations/000001_migration_script_init.go",
    ['"url": "/tasks"', '"url":"/tasks"'],
  ],
  [
    "migrations/000011_migration_script_site_content_config.go",
    ['Url:             "/links"', 'Url:"/links"'],
  ],
  [
    "migrations/000008_migration_script_notification_types.go",
    ['"topicFavorite"', '"userLevelUp"', '"userBadgeGrant"'],
  ],
  [
    "migrations/migration.go",
    ["qa bounty config defaults", "migrate_qa_bounty_config"],
  ],
  ["internal/services/seo_sitemap_service.go", ['AbsUrl("/links")']],
]) {
  const text = repoSource(file)
  for (const forbidden of forbiddenList) {
    assert.equal(
      text.includes(forbidden),
      false,
      `${file} should not expose removed product feature: ${forbidden}`
    )
  }
}

for (const i18nFile of [
  "lib/i18n/messages/en-US.ts",
  "lib/i18n/messages/zh-CN.ts",
]) {
  const text = source(i18nFile)
  for (const forbidden of [
    "favorites:",
    "userBadges:",
    "userExpLogs:",
    "userTaskLogs:",
    "fansWidget:",
    "followWidget:",
    "followBtn:",
    "userFans:",
    "badgesTitle:",
    "badgesSubtitle:",
    "noBadges:",
    "Community Links",
    "Growth",
    "成长体系",
    "友情链接",
    "Search discussions, articles, or members",
    "搜索帖子、文章、用户",
    "tasks:",
    "enableQaBounty",
    "qaBounty",
    "bountyLabel",
    "attachmentScoreRequired",
    "scorePlaceholder",
    "taskList:",
  ]) {
    assert.equal(
      text.includes(forbidden),
      false,
      `${i18nFile} should not retain removed i18n fragment: ${forbidden}`
    )
  }
}

for (const [file, forbiddenList] of [
  [
    "components/editor/rich-text-editor.tsx",
    [
      "extension-task-list",
      "extension-task-item",
      "taskList",
      "toggleTaskList",
      "任务列表",
      "待办",
    ],
  ],
  [
    "components/topic/topic-create-form.tsx",
    ["bountyScore", "downloadScore", "update_download_score", "scorePlaceholder"],
  ],
  [
    "components/topic/topic-edit-form.tsx",
    ["downloadScore", "update_download_score", "scorePlaceholder"],
  ],
  [
    "components/topic/topic-meta.tsx",
    ["bountyScore", "bountyLabel"],
  ],
  [
    "components/topic/topic-list-item.tsx",
    ["bountyScore", "bountyLabel"],
  ],
  [
    "components/topic/topic-attachments.tsx",
    ["downloadScore", "attachmentScoreRequired"],
  ],
  [
    "lib/api/types.ts",
    [
      "score?:",
      "downloadScore?:",
      "bountyScore?:",
      "enableQaBounty?:",
      "qaBountyMin",
      "qaBountyMax",
      "qaBountyRequired",
    ],
  ],
]) {
  const text = source(file)
  for (const forbidden of forbiddenList) {
    assert.equal(
      text.includes(forbidden),
      false,
      `${file} should not retain removed score/bounty/task-list feature: ${forbidden}`
    )
  }
}

console.log("product simplification tests passed")
