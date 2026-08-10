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
