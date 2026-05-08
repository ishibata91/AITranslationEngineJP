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

export interface TranslationManagementViewContract {
  id: TranslationManagementViewId
  label: string
  description: string
  stepNumber: number
  directNavigation: boolean
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

const TRANSLATION_MANAGEMENT_VIEW_CONTRACT: ReadonlyArray<TranslationManagementViewContract> =
  [
    {
      id: "job-management",
      label: "未完了のジョブ",
      description: "新しい翻訳を始めるか、途中のジョブを再開します。",
      stepNumber: 1,
      directNavigation: true
    },
    {
      id: "input-review",
      label: "入力データの確認",
      description: "翻訳に使う入力データを選び、登録結果を確認します。",
      stepNumber: 2,
      directNavigation: false
    },
    {
      id: "job-setup",
      label: "翻訳設定",
      description: "入力データと AI 設定を確認し、ジョブを作成します。",
      stepNumber: 3,
      directNavigation: false
    },
    {
      id: "term-translation",
      label: "単語翻訳",
      description: "選択したジョブで、単語翻訳を実行します。",
      stepNumber: 4,
      directNavigation: false
    },
    {
      id: "persona-generation",
      label: "NPC ペルソナ生成",
      description: "単語翻訳の完了後に、NPC の話し方や役割を整理します。",
      stepNumber: 5,
      directNavigation: false
    },
    {
      id: "body-translation",
      label: "本文翻訳",
      description: "NPC ペルソナを参照できる状態で、本文の翻訳を実行します。",
      stepNumber: 6,
      directNavigation: false
    },
    {
      id: "translation-complete",
      label: "翻訳結果の確認",
      description: "本文翻訳が完了した後に、原文と訳文を確認します。",
      stepNumber: 7,
      directNavigation: false
    },
    {
      id: "output-management",
      label: "出力管理",
      description: "翻訳結果を確認した後に、出力するジョブを選びます。",
      stepNumber: 8,
      directNavigation: false
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
    defaultTranslationManagementViewId: "job-management",
    translationManagementViews: TRANSLATION_MANAGEMENT_VIEW_CONTRACT.map(
      (view) => ({
        ...view
      })
    )
  }
}
