export type ShellRouteId =
  | "dashboard"
  | "provider-settings"
  | "master-dictionary"
  | "master-persona"
  | "translation-management"
  | "output-management"

export type TranslationManagementViewId =
  | "input-review"
  | "job-management"
  | "job-run"
  | "term-translation"
  | "persona-generation"
  | "body-translation"
  | "translation-complete"
  | "output-management"

export interface ShellRouteContract {
  id: ShellRouteId
  label: string
  state: string
  lead: string
  description: string
}

const SHELL_ROUTE_CONTRACT: ReadonlyArray<ShellRouteContract> = [
  {
    id: "dashboard",
    label: "ダッシュボード",
    state: "既定表示",
    lead: "最初に移動したい作業を選び、共通ナビゲーションからいつでも別の主要ページへ切り替えられます。",
    description: "主要ページへの入口をまとめて確認します。"
  },
  {
    id: "provider-settings",
    label: "AIサービス設定",
    state: "設定状態を確認",
    lead: "AIサービスごとの接続設定をまとめて確認し、保存と接続確認を行うページです。",
    description: "エンドポイントと APIキー状態を AIサービスごとに確認します。"
  },
  {
    id: "master-dictionary",
    label: "マスター辞書",
    state: "確認可能",
    lead: "用語と訳語の基盤データを確認するページです。",
    description: "用語と訳語の基盤データを確認します。"
  },
  {
    id: "master-persona",
    label: "マスターペルソナ",
    state: "確認可能",
    lead: "ベースゲームや大型 Mod の NPC を対象に、翻訳前の準備としてペルソナをまとめて作成する。作成後は一覧と詳細で同じ画面から確認できる。",
    description: "翻訳に使うペルソナ設定を確認します。"
  },
  {
    id: "translation-management",
    label: "翻訳管理",
    state: "未完了ジョブ入口",
    lead: "未完了ジョブ一覧から新規翻訳の開始と途中再開を選び、対象ジョブを固定して翻訳を進めるページです。",
    description:
      "未完了ジョブ一覧、新規翻訳の開始、対象ジョブの現在の翻訳段階への再開を扱います。"
  },
  {
    id: "output-management",
    label: "出力管理",
    state: "確認可能",
    lead: "生成された成果物を確認するページです。",
    description: "生成物と書き出し結果を確認します。"
  }
]

interface ShellState {
  defaultRouteId: ShellRouteId
  routes: ShellRouteContract[]
  defaultTranslationManagementViewId: TranslationManagementViewId
}

export function createShellState(): ShellState {
  return {
    defaultRouteId: "dashboard",
    routes: SHELL_ROUTE_CONTRACT.map((route) => ({ ...route })),
    defaultTranslationManagementViewId: "job-management"
  }
}
