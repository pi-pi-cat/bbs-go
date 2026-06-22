import { serverApiFetch as apiFetch } from "./server"

export interface AboutConfig {
  content?: string
}

export interface InstallStatus {
  installed?: boolean
}

export function getAboutConfig() {
  return apiFetch<AboutConfig>("/api/config/about")
}

export function getInstallStatus() {
  return apiFetch<InstallStatus>("/api/install/status")
}
