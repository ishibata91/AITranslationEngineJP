import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TranslationRunScreen from "./TranslationRunScreen.svelte"
import {
  EMPTY_STATE,
  MODELS_LOADING_STATE,
  READY_STATE,
  RUNNING_EXTRACT_STATE,
  RUNNING_TRANSLATE_STATE,
  DONE_STATE,
  DONE_PERSONA_STATE,
  ERROR_STATE,
  XAI_EMPTY_STATE,
  XAI_READY_STATE,
  XAI_SUBMITTED_STATE,
  XAI_CHECKING_STATE,
  XAI_PROPER_PROCESSING_STATE,
  XAI_PROPER_READY_STATE,
  XAI_BODY_PROCESSING_STATE,
  XAI_BODY_READY_STATE,
  XAI_DONE_STATE
} from "./translation-run.fixtures"

// プラグイン選択を外し、翻訳対象プラグイン画面へ移した実行・結果画面。翻訳対象は読み取り専用表示。
// Storybook 人間レビュー承認済み（通常分類 Screens）。
const meta = {
  title: "Screens/翻訳実行",
  component: TranslationRunScreen,
  parameters: {
    layout: "fullscreen"
  },
  args: {
    onFieldInput: () => {},
    onLoadModels: () => {},
    onRun: () => {},
    onPagePrev: () => {},
    onPageNext: () => {},
    onProviderChange: () => {},
    onSubmit: () => {},
    onRefresh: () => {},
    onCheckStatus: () => {},
    onApply: () => {}
  }
} satisfies Meta<typeof TranslationRunScreen>

export default meta
type Story = StoryObj<typeof meta>

// 初期表示。未入力・モデル未取得で実行できない。
export const Empty: Story = {
  name: "空状態",
  args: { ...EMPTY_STATE }
}

// getModels でモデル一覧を取得中。
export const ModelsLoading: Story = {
  name: "モデル取得中",
  args: { ...MODELS_LOADING_STATE }
}

// 入力が揃い、実行できる。
export const Ready: Story = {
  name: "入力済み",
  args: { ...READY_STATE }
}

// 実行中（台詞抽出）。不定の進捗バーが出る。
export const RunningExtract: Story = {
  name: "実行中（抽出）",
  args: { ...RUNNING_EXTRACT_STATE }
}

// 実行中（本文翻訳）。done/total の確定進捗バーが出る。
export const RunningTranslate: Story = {
  name: "実行中（本文翻訳）",
  args: { ...RUNNING_TRANSLATE_STATE }
}

// 完了（叙述文）。原文と訳文が並ぶ。
export const Done: Story = {
  name: "完了",
  args: { ...DONE_STATE }
}

// 完了（台詞）。訳文と口調指示が並び、話者ごとの口調差が観測できる。
export const DonePersona: Story = {
  name: "完了（口調差）",
  args: { ...DONE_PERSONA_STATE }
}

// 実行に失敗した。
export const Errored: Story = {
  name: "エラー",
  args: { ...ERROR_STATE }
}

// xAI（batch）を選び、未入力で送信できない。
export const XaiEmpty: Story = {
  name: "xAI・未入力",
  args: { ...XAI_EMPTY_STATE }
}

// xAI（batch）で接続情報とモデルが揃い、状態未確認。主アクションは「送信して開始」で活性。
export const XaiReady: Story = {
  name: "xAI・送信可",
  args: { ...XAI_READY_STATE }
}

// xAI（batch）を送信した直後。案内が出る。
export const XaiSubmitted: Story = {
  name: "xAI・送信後",
  args: { ...XAI_SUBMITTED_STATE }
}

// xAI（batch）で状態確認中。状態確認ボタンにスピナーが出る。
export const XaiChecking: Story = {
  name: "xAI・状態確認中",
  args: { ...XAI_CHECKING_STATE }
}

// xAI（batch）で固有名段が処理中。ステッパー現在地=固有名、主アクションはグレーアウト。
export const XaiProperProcessing: Story = {
  name: "xAI・固有名処理中",
  args: { ...XAI_PROPER_PROCESSING_STATE }
}

// xAI（batch）で固有名段が完了。主アクションは「取り込んで本文を送信」。
export const XaiProperReady: Story = {
  name: "xAI・固有名完了",
  args: { ...XAI_PROPER_READY_STATE }
}

// xAI（batch）で本文段が処理中。ステッパー現在地=本文、主アクションはグレーアウト。
export const XaiBodyProcessing: Story = {
  name: "xAI・本文処理中",
  args: { ...XAI_BODY_PROCESSING_STATE }
}

// xAI（batch）で本文段が完了。主アクションは「取り込んで完了」。
export const XaiBodyReady: Story = {
  name: "xAI・本文完了",
  args: { ...XAI_BODY_READY_STATE }
}

// xAI（batch）で取り込みが終わり、結果が入った。段=完了、主アクションは再送信の「送信して開始」。
export const XaiDone: Story = {
  name: "xAI・取り込み済み",
  args: { ...XAI_DONE_STATE }
}
