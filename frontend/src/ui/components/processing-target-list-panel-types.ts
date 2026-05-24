export interface ProcessingTargetListItem {
  id: string
  name: string
  titleParts?: ProcessingTargetListItemTitlePart[]
  detail: string
  metadata: ProcessingTargetListItemMetadata[]
  actions?: ProcessingTargetListItemAction[]
}

export interface ProcessingTargetListItemTitlePart {
  text: string
}

export interface ProcessingTargetListItemMetadata {
  label: string
  value: string
}

export interface ProcessingTargetListItemAction {
  label: string
  variant?: "secondary" | "danger"
  disabled?: boolean
  onAction: () => void
}

export interface ProcessingTargetListPageState {
  items: ProcessingTargetListItem[]
  page: number
  pageSize: number
  totalCount: number
  searchQuery?: string
  busy?: boolean
}
