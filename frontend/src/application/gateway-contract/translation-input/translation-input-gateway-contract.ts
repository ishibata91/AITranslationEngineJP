export type TranslationInputErrorKind =
  | "invalid_json"
  | "duplicate_input_hash"
  | "unsupported_extract_shape"
  | "missing_required_field"
  | "source_file_missing"
  | "cache_missing"

export type TranslationInputWarningKind = "unknown_field_definition"

export interface ImportTranslationInputRequest {
  filePath: string
}

export interface RebuildTranslationInputCacheRequest {
  inputId: number
}

export interface TranslationInputImportedInput {
  id: number
  sourceFilePath: string
  sourceTool: string
  targetPluginName: string
  targetPluginType: string
  recordCount: number
  importedAt: string
}

export interface TranslationInputCategoryCount {
  category: string
  recordCount: number
  fieldCount: number
}

export interface TranslationInputSampleField {
  recordType: string
  subrecordType: string
  formId: string
  editorId: string
  sourceText: string
  translatable: boolean
}

export interface TranslationInputWarning {
  kind: TranslationInputWarningKind
  recordType: string
  subrecordType: string
  message: string
}

export interface TranslationInputImportSummary {
  input: TranslationInputImportedInput
  translationRecordCount: number
  translationFieldCount: number
  categories: TranslationInputCategoryCount[]
  sampleFields: TranslationInputSampleField[]
  warnings: TranslationInputWarning[]
}

export type TranslationInputOperationState =
  | "idle"
  | "ready"
  | "importing"
  | "rebuilding"

export type TranslationInputReviewStatus =
  | "registered"
  | "warning"
  | "failed"
  | "rebuild-required"

export interface TranslationInputStagedFile {
  fileName: string
  filePath: string
  fileHash: string
}

export interface TranslationInputReviewItem {
  localId: string
  inputId: number | null
  fileName: string
  filePath: string
  fileHash: string
  importTimestamp: string
  status: TranslationInputReviewStatus
  accepted: boolean
  canRebuild: boolean
  lastAction: "import" | "rebuild"
  errorKind: TranslationInputErrorKind | null
  warnings: TranslationInputWarning[]
  summary: TranslationInputImportSummary | null
}

export interface TranslationInputScreenState {
  items: TranslationInputReviewItem[]
  selectedItemId: string | null
  stagedFile: TranslationInputStagedFile | null
  operationState: TranslationInputOperationState
  errorMessage: string
  latestResponse: TranslationInputCommandResponse | null
}

export interface TranslationInputScreenViewModel
  extends TranslationInputScreenState {
  selectedItem: TranslationInputReviewItem | null
  gatewayStatus: string
  hasStagedFile: boolean
  canImport: boolean
  canRebuildSelected: boolean
  isImporting: boolean
  isRebuilding: boolean
  stagedFileName: string
  stagedFilePath: string
  stagedFileHash: string
  operationStatusLabel: string
  operationStatusText: string
  latestOutcomeTitle: string
  latestOutcomeText: string
  selectionStatusText: string
  totalItemCountLabel: string
  emptyStateText: string
}

export interface TranslationInputCommandResponse {
  accepted: boolean
  summary?: TranslationInputImportSummary
  errorKind?: TranslationInputErrorKind
  warnings: TranslationInputWarning[]
}

export interface TranslationInputGatewayContract {
  importTranslationInput(
    request: ImportTranslationInputRequest
  ): Promise<TranslationInputCommandResponse>
  rebuildTranslationInputCache(
    request: RebuildTranslationInputCacheRequest
  ): Promise<TranslationInputCommandResponse>
}