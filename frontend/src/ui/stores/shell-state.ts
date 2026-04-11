export type ShellRouteId =
  | "dashboard"
  | "master-dictionary"
  | "master-persona"
  | "translation-management"
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
    id: "master-dictionary",
    label: "マスター辞書",
    state: "準備中",
    lead: "用語と訳語の基盤データを確認するページです。",
    description: "用語と訳語の基盤データを確認します。"
  },
  {
    id: "master-persona",
    label: "マスターペルソナ",
    state: "準備中",
    lead: "翻訳に使うペルソナ設定を確認するページです。",
    description: "翻訳に使うペルソナ設定を確認します。"
  },
  {
    id: "translation-management",
    label: "翻訳管理",
    state: "準備中",
    lead: "翻訳準備と実行状況をまとめて確認するページです。",
    description: "準備状況と翻訳ジョブの進行をまとめて確認します。"
  },
  {
    id: "output-management",
    label: "出力管理",
    state: "準備中",
    lead: "生成された成果物を確認するページです。",
    description: "生成物と書き出し結果を確認します。"
  }
]

interface ShellState {
  defaultRouteId: ShellRouteId
  routes: ShellRouteContract[]
}

export function createShellState(): ShellState {
  return {
    defaultRouteId: "dashboard",
    routes: SHELL_ROUTE_CONTRACT.map((route) => ({ ...route }))
  }
}
