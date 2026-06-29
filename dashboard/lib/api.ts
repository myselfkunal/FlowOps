const API_BASE = "http://localhost:8081"

export type ServiceStatus = {
  name: string
  desired_replicas: number
  ready_replicas: number
  image: string
  status: "healthy" | "degraded" | "unknown"
}

export type ReconcileEvent = {
  timestamp: string
  service_name: string
  what_changed: string
  old_value: string
  new_value: string
}

export type ConfigResponse = {
  config: string
}

export async function fetchStatus(): Promise<ServiceStatus[]> {
  const res = await fetch(`${API_BASE}/status`, { cache: "no-store" })
  return res.json()
}

export async function fetchHistory(): Promise<ReconcileEvent[]> {
  const res = await fetch(`${API_BASE}/history`, { cache: "no-store" })
  return res.json()
}

export async function fetchConfig(): Promise<ConfigResponse> {
  const res = await fetch(`${API_BASE}/config`, { cache: "no-store" })
  return res.json()
}