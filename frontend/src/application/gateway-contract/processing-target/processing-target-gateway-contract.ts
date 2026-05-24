export interface ProcessingTargetListRequest {
  jobId: number
  phase: string
  page: number
  pageSize: number
  searchQuery: string
}

export interface ProcessingTargetListItemTitlePart {
  text: string
}

export interface ProcessingTargetListItemMetadata {
  label: string
  value: string
}

export interface ProcessingTargetListItem {
  id: string
  name: string
  detail: string
  titleParts: ProcessingTargetListItemTitlePart[]
  metadata: ProcessingTargetListItemMetadata[]
}

export interface ProcessingTargetListResponse {
  items: ProcessingTargetListItem[]
  metadata: ProcessingTargetListItemMetadata[]
  page: number
  pageSize: number
  totalCount: number
  searchQuery: string
}

export interface ProcessingTargetListPageState extends ProcessingTargetListResponse {
  busy?: boolean
}

export type ProcessingTargetListPageStatesByPhase = Partial<
  Record<string, ProcessingTargetListPageState>
>
