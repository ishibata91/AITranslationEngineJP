export type ShellRouteId =
  | "dashboard"
  | "provider-settings"
  | "master-dictionary"
  | "master-persona"
  | "translation-management"
  | "output-management"

export type TranslationManagementViewId =
  | "input-review"
  | "job-setup"
  | "job-management"
  | "job-run"

export interface ShellRouteContract {
  id: ShellRouteId
  label: string
  state: string
  lead: string
  description: string
}

export interface TranslationManagementViewContract {
  id: TranslationManagementViewId
  label: string
  description: string
  stepNumber: number
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
    state: "ジョブ管理を含む",
    lead: "ジョブ管理、データロード、セットアップ、実行を順番に切り替え、未完了 job 管理から翻訳実行表示までを確認するページです。",
    description:
      "ジョブ管理、データロード、validation、ready job 作成、term phase、persona phase、body phase の実行状況をまとめて確認します。"
  },
  {
    id: "output-management",
    label: "出力管理",
    state: "確認可能",
    lead: "生成された成果物を確認するページです。",
    description: "生成物と書き出し結果を確認します。"
  }
]

const TRANSLATION_MANAGEMENT_VIEW_CONTRACT: ReadonlyArray<TranslationManagementViewContract> =
  [
    {
      id: "job-management",
      label: "ジョブ管理",
      description: "未完了 job の一覧、詳細、操作可否を最初に確認します。",
      stepNumber: 1
    },
    {
      id: "input-review",
      label: "データロード",
      description: "入力ファイルの登録結果と再構築判断を確認します。",
      stepNumber: 2
    },
    {
      id: "job-setup",
      label: "セットアップ",
      description: "validation と ready job 作成を確認します。",
      stepNumber: 3
    },
    {
      id: "job-run",
      label: "実行",
      description:
        "term phase、persona phase、body phase の progress、result summary、output readiness を確認します。",
      stepNumber: 4
    }
  ]

interface ShellState {
  defaultRouteId: ShellRouteId
  routes: ShellRouteContract[]
  defaultTranslationManagementViewId: TranslationManagementViewId
  translationManagementViews: TranslationManagementViewContract[]
}

export function createShellState(): ShellState {
  return {
    defaultRouteId: "dashboard",
    routes: SHELL_ROUTE_CONTRACT.map((route) => ({ ...route })),
    defaultTranslationManagementViewId: "input-review",
    translationManagementViews: TRANSLATION_MANAGEMENT_VIEW_CONTRACT.map(
      (view) => ({
        ...view
      })
    )
  }
}
