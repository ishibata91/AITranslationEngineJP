export type PhaseStateToken =
  | "ready"
  | "running"
  | "paused"
  | "recoverable_failed"
  | "completed"
  | "failed"
  | "canceled"
  | "loading"
  | "idle_ready"
  | "empty_completed"
  | "blocked"
  | "not_started"
  | "snapshot_missing"

export interface PhaseMetricCounter {
  label: string
  value: string
}

export interface PhaseDetailItem {
  label: string
  value: string
  note?: string
}

export interface PhaseActionItem<ActionId extends string = string> {
  id: ActionId
  label: string
  disabled: boolean
  blockedReason: string
  tone: "default" | "primary" | "warning"
}
