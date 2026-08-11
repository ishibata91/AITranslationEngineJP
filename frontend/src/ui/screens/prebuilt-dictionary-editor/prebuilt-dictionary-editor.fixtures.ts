import type {
  PrebuiltDictionaryFilters,
  PrebuiltDictionaryForm,
  PrebuiltDictionaryRow
} from "./prebuilt-dictionary-editor-view"

export const PREBUILT_DICTIONARY_PAGE_SIZE = 50

export const emptyPrebuiltDictionaryFilters: PrebuiltDictionaryFilters = {
  source: "",
  destination: "",
  partOfSpeech: "",
  category: ""
}

export const emptyPrebuiltDictionaryForm: PrebuiltDictionaryForm = {
  source: "",
  destination: "",
  partOfSpeech: "",
  meaning: ""
}

export const prebuiltDictionaryRows: PrebuiltDictionaryRow[] = [
  {
    id: "1",
    source: "Companion",
    destination: "同胞団",
    partOfSpeech: "noun",
    categories: ["NPC_"]
  },
  {
    id: "2",
    source: "Companion",
    destination: "仲間",
    partOfSpeech: "noun",
    categories: ["DIAL", "INFO"]
  },
  {
    id: "3",
    source: "Dragonborn",
    destination: "ドラゴンボーン",
    partOfSpeech: "noun",
    categories: ["QUST"]
  }
]

export const prebuiltDictionaryPageRows: PrebuiltDictionaryRow[] = Array.from(
  { length: PREBUILT_DICTIONARY_PAGE_SIZE },
  (_, index) => {
    const sample = prebuiltDictionaryRows[index % prebuiltDictionaryRows.length]
    return {
      ...sample,
      id: `page-${index + 1}`,
      source: `${sample.source} ${index + 1}`,
      destination: `${sample.destination} ${index + 1}`,
      categories: [...sample.categories]
    }
  }
)
