import type { BodyTranslationFieldResultItem } from "@application/contract/body-translation-phase"
import type { PhaseDetailItem } from "../../components/phase-panel-types"

export type { PhaseDetailItem }

export interface BodyInputSummaryCardProps {
  readinessReason: string
  details: PhaseDetailItem[]
}

export interface BodyExecutionSummaryCardProps {
  providerStateLabel: string
  details: PhaseDetailItem[]
}

export interface BodyResultSummaryCardProps {
  outputReadinessLabel: string
  details: PhaseDetailItem[]
}

export interface FieldResultListPanelProps {
  availabilityLabel: string
  items: BodyTranslationFieldResultItem[]
}

export interface OutputReadinessCardProps {
  outputReadinessLabel: string
  details: PhaseDetailItem[]
}
