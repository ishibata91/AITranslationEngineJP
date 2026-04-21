export interface MasterDictionaryListFilters {
  query: string
  category: string
  page: number
  pageSize: number
}

export interface MasterDictionaryFrontendRefresh {
  query: string
  category: string
  page: number
  pageSize: number
}

export interface ListMasterDictionaryEntriesRequest {
  filters: MasterDictionaryListFilters
}

export interface MasterDictionaryEntrySummary {
  id: string
  source: string
  translation: string
  category: string
  origin: string
  updatedAt: string
}

export interface ListMasterDictionaryEntriesResponse {
  entries: MasterDictionaryEntrySummary[]
  totalCount: number
  page: number
  pageSize: number
}

export interface GetMasterDictionaryEntryRequest {
  id: string
}

export interface MasterDictionaryEntryDetail extends MasterDictionaryEntrySummary {
  note: string
}

export interface GetMasterDictionaryEntryResponse {
  entry: MasterDictionaryEntryDetail | null
}

export interface MasterDictionaryPageEntry {
  id: number
  source: string
  translation: string
  category: string
  origin: string
  updatedAt: string
}

export interface MasterDictionaryPageState {
  items: MasterDictionaryPageEntry[]
  totalCount: number
  page: number
  pageSize: number
  selectedId?: number
}

export interface MasterDictionaryUpsertPayload {
  source: string
  translation: string
  category: string
  origin: string
}

export interface CreateMasterDictionaryEntryRequest {
  payload: MasterDictionaryUpsertPayload
  refresh?: MasterDictionaryFrontendRefresh
}

export interface CreateMasterDictionaryEntryResponse {
  entry: MasterDictionaryEntryDetail
  refreshTargetId: string
  page?: MasterDictionaryPageState
}

export interface UpdateMasterDictionaryEntryRequest {
  id: string
  payload: MasterDictionaryUpsertPayload
  refresh?: MasterDictionaryFrontendRefresh
}

export interface UpdateMasterDictionaryEntryResponse {
  entry: MasterDictionaryEntryDetail
  refreshTargetId: string
  page?: MasterDictionaryPageState
}

export interface DeleteMasterDictionaryEntryRequest {
  id: string
  refresh?: MasterDictionaryFrontendRefresh
}

export interface DeleteMasterDictionaryEntryResponse {
  deletedId: string
  nextSelectedId: string | null
  page?: MasterDictionaryPageState
}

export interface ImportMasterDictionaryXmlRequest {
  filePath: string
  fileReference?: string
  refresh?: MasterDictionaryFrontendRefresh
}

export interface MasterDictionaryImportSummary {
  filePath: string
  fileName: string
  importedCount: number
  updatedCount: number
  skippedCount: number
  lastEntryId: number
}

export interface ImportMasterDictionaryXmlResponse {
  accepted: boolean
  page?: MasterDictionaryPageState
  summary?: MasterDictionaryImportSummary
}

export interface MasterDictionaryGatewayContract {
  listMasterDictionaryEntries(
    request: ListMasterDictionaryEntriesRequest
  ): Promise<ListMasterDictionaryEntriesResponse>
  getMasterDictionaryEntry(
    request: GetMasterDictionaryEntryRequest
  ): Promise<GetMasterDictionaryEntryResponse>
  createMasterDictionaryEntry(
    request: CreateMasterDictionaryEntryRequest
  ): Promise<CreateMasterDictionaryEntryResponse>
  updateMasterDictionaryEntry(
    request: UpdateMasterDictionaryEntryRequest
  ): Promise<UpdateMasterDictionaryEntryResponse>
  deleteMasterDictionaryEntry(
    request: DeleteMasterDictionaryEntryRequest
  ): Promise<DeleteMasterDictionaryEntryResponse>
  importMasterDictionaryXml(
    request: ImportMasterDictionaryXmlRequest
  ): Promise<ImportMasterDictionaryXmlResponse>
}
