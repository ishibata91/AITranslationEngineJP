import type { TranslationJobManagementJobRunTarget } from "@application/contract/translation-job-management/translation-job-management-screen-types"

export interface JobRunTargetSummaryProps {
  target: TranslationJobManagementJobRunTarget
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
