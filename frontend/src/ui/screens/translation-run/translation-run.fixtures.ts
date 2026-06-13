// 翻訳実行画面の表示用 fixture。
// 画面単位の代表状態（空・入力済み・実行中・完了・エラー）を固定する。
import type {
  TranslationRunForm,
  NarrationResultRow,
  RunPhase
} from "./translation-run-view"

// 画面全体の表示状態。story の固定状態を組むための型。
interface TranslationRunViewModel {
  form: TranslationRunForm
  phase: RunPhase
  canRun: boolean
  models: string[]
  modelsLoading: boolean
  results: NarrationResultRow[]
  errorMessage: string
}

const EMPTY_FORM: TranslationRunForm = {
  pluginPath: "",
  endpoint: "",
  apiKey: "",
  model: ""
}

const FILLED_FORM: TranslationRunForm = {
  pluginPath: "/Users/me/Skyrim/Data/Dawnguard.esm",
  endpoint: "https://api.openai.com/v1",
  apiKey: "sk-demo-key",
  model: "gpt-4o-mini"
}

// getModels で取得した一覧の表示用サンプル。
const MODELS = ["gpt-4o-mini", "gpt-4o", "gpt-4.1-mini", "o4-mini"]

const RESULT_ROWS: NarrationResultRow[] = [
  {
    edid: "DLC1BookSerana",
    source:
      "I have walked these halls for centuries, and still the cold of Castle Volkihar finds me.",
    dest: "私は何世紀もこの広間を歩いてきたが、それでもヴォルキハル城の冷気は私を捉えて離さない。",
    statusLabel: "仮訳"
  },
  {
    edid: "DLC1NoteFromArvak",
    source:
      "Whoever finds this: the Soul Cairn takes more than it gives. Do not bargain with the Ideal Masters.",
    dest: "これを見つけた者へ。ソウル・ケルンは与える以上に奪う。理想の支配者と取引してはならない。",
    statusLabel: "仮訳"
  },
  {
    edid: "DLC1BookAetherium",
    source:
      "The Aetherium Forge lies beneath Blackreach, its fires long since dimmed by the dwarves who fled.",
    dest: "",
    statusLabel: "未訳"
  }
]

// 初期表示。何も入力しておらず、モデルも未取得で実行できない。
export const EMPTY_STATE: TranslationRunViewModel = {
  form: EMPTY_FORM,
  phase: "idle",
  canRun: false,
  models: [],
  modelsLoading: false,
  results: [],
  errorMessage: ""
}

// エンドポイントと API キーを入れて、モデルを取得中。
export const MODELS_LOADING_STATE: TranslationRunViewModel = {
  form: { ...FILLED_FORM, model: "" },
  phase: "idle",
  canRun: false,
  models: [],
  modelsLoading: true,
  results: [],
  errorMessage: ""
}

// plugin・接続情報・モデルが揃い、実行できる状態。
export const READY_STATE: TranslationRunViewModel = {
  form: FILLED_FORM,
  phase: "idle",
  canRun: true,
  models: MODELS,
  modelsLoading: false,
  results: [],
  errorMessage: ""
}

// 抽出と翻訳を実行中。ボタンは押せない。
export const RUNNING_STATE: TranslationRunViewModel = {
  form: FILLED_FORM,
  phase: "running",
  canRun: false,
  models: MODELS,
  modelsLoading: false,
  results: [],
  errorMessage: ""
}

// 完了し、原文と訳文が並ぶ。一部は未訳のまま残る場合もある。
export const DONE_STATE: TranslationRunViewModel = {
  form: FILLED_FORM,
  phase: "done",
  canRun: true,
  models: MODELS,
  modelsLoading: false,
  results: RESULT_ROWS,
  errorMessage: ""
}

// 実行が失敗し、原因と対応を伝える。
export const ERROR_STATE: TranslationRunViewModel = {
  form: FILLED_FORM,
  phase: "error",
  canRun: true,
  models: MODELS,
  modelsLoading: false,
  results: [],
  errorMessage:
    "翻訳 API への接続に失敗しました。エンドポイントと API キーを確認して、もう一度実行してください。"
}
