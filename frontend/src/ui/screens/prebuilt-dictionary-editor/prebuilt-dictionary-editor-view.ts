export interface PrebuiltDictionaryRow {
  id: string
  source: string
  destination: string
  partOfSpeech: string
  categories: string[]
  pending?: "edited" | "deleted"
}

export interface PrebuiltDictionaryFilters {
  source: string
  destination: string
  partOfSpeech: string
  category: string
}

export interface PrebuiltDictionaryForm {
  source: string
  destination: string
  partOfSpeech: string
  meaning: string
}
