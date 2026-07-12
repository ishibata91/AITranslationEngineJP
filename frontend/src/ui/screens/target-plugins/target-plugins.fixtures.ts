// 翻訳対象プラグイン画面の表示用 fixture。
// 画面単位の代表状態（空・読み込み中・一覧・選択済み・削除確認中・削除実行中・エラー）を固定する。
import type { TargetPluginRow } from "./target-plugins-view"

// 画面全体の表示状態。story の固定状態を組むための型。
interface TargetPluginsViewModel {
  newPluginPath: string
  plugins: TargetPluginRow[]
  loading?: boolean
  pendingDeletePlugin?: string
  deleting?: boolean
  errorMessage?: string
}

// 一覧の表示用サンプル。完了・翻訳中・未着手・大量件数の進捗差が観測できる。
const PLUGIN_ROWS: TargetPluginRow[] = [
  {
    plugin: "Dawnguard.esm",
    createdLabel: "2026-07-12 14:32",
    total: 812,
    translated: 812
  },
  {
    plugin: "Dragonborn.esm",
    createdLabel: "2026-07-11 09:15",
    total: 1340,
    translated: 604
  },
  {
    plugin: "Inigo.esp",
    createdLabel: "2026-07-10 22:48",
    total: 356,
    translated: 0
  },
  {
    plugin: "3DNPC.esp",
    createdLabel: "2026-07-09 18:03",
    total: 20418,
    translated: 20418
  }
]

// まだ翻訳を実行しておらず、対象プラグインが無い初期表示。
export const EMPTY_STATE: TargetPluginsViewModel = { newPluginPath: "", plugins: [] }

// 一覧を読み込み中。
export const LOADING_STATE: TargetPluginsViewModel = {
  newPluginPath: "",
  plugins: [],
  loading: true
}

// 複数のプラグインが並ぶ。完了・翻訳中・未着手の進捗差が出る。
export const LIST_STATE: TargetPluginsViewModel = {
  newPluginPath: "",
  plugins: PLUGIN_ROWS
}

// 新しいプラグインを選び終え、翻訳へ進める状態。
export const SELECTED_STATE: TargetPluginsViewModel = {
  newPluginPath: "/Users/me/Skyrim/Data/Falskaar.esm",
  plugins: PLUGIN_ROWS
}

// 1 件の削除確認中。対象行だけ確認 UI を出す。
export const CONFIRM_DELETE_STATE: TargetPluginsViewModel = {
  newPluginPath: "",
  plugins: PLUGIN_ROWS,
  pendingDeletePlugin: "Dragonborn.esm"
}

// 削除実行中。確定ボタンがスピナー・無効化される。
export const DELETING_STATE: TargetPluginsViewModel = {
  newPluginPath: "",
  plugins: PLUGIN_ROWS,
  pendingDeletePlugin: "Dragonborn.esm",
  deleting: true
}

// 削除に失敗した。
export const ERROR_STATE: TargetPluginsViewModel = {
  newPluginPath: "",
  plugins: PLUGIN_ROWS,
  errorMessage: "プラグインの翻訳成果の削除に失敗しました。もう一度お試しください。"
}
