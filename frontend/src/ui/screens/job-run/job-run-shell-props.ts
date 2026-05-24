import type { TranslationJobManagementJobRunTarget } from "@application/contract/translation-job-management/translation-job-management-screen-types"
import type {
  ProcessingTargetListItem,
  ProcessingTargetListItemMetadata,
  ProcessingTargetListItemTitlePart
} from "@ui/components/processing-target-list-panel-types"

export interface JobRunTargetSummaryProps {
  target: TranslationJobManagementJobRunTarget
  currentPhasePage?: JobRunPhaseStepId
}

export type JobRunPhaseStepId = "term" | "persona" | "body" | "complete"

export type {
  ProcessingTargetListItem,
  ProcessingTargetListItemMetadata,
  ProcessingTargetListItemTitlePart
}

export interface JobUnselectedGuidanceProps {
  onOpenJobManagement?: () => void
}

export interface PhaseNavigationFooterProps {
  title: string
  titleId: string
  description: string
  reasons: string[]
  primaryLabel?: string
  primaryDisabled?: boolean
  showPrimary?: boolean
  showBack?: boolean
  showOutput?: boolean
  onPrimary?: () => void
  onBack?: () => void
  onOutput?: () => void
  dataTestId?: string
}
