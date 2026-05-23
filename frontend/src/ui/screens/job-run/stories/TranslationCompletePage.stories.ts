import type { Meta, StoryObj } from "@storybook/svelte-vite"
import type { BodyTranslationFieldResultItem } from "@application/contract/body-translation-phase"

import TranslationCompletePage from "../TranslationCompletePage.svelte"

const rows: BodyTranslationFieldResultItem[] = [
  {
    fieldId: "field-name-001",
    fieldLabel: "Name",
    recordTypeLabel: "NPC_",
    fieldTypeLabel: "FULL",
    formIdLabel: "0x00001234",
    editorIdLabel: "SampleActorOne",
    sourceExcerpt: "Sample source line for completion review.",
    translatedText: "完了画面確認用のサンプル訳文です。",
    outputStatus: "ready",
    protectionValidation: "passed",
    retryCountLabel: "0",
    rawItem: null
  },
  {
    fieldId: "field-dialogue-002",
    fieldLabel: "Dialogue",
    recordTypeLabel: "INFO",
    fieldTypeLabel: "NAM1",
    formIdLabel: "0x00005678",
    editorIdLabel: "SampleDialogueTopic",
    sourceExcerpt:
      "A longer synthetic source line for pagination and wrapping review.",
    translatedText:
      "ページングと折り返しを確認するための長めの合成訳文です。",
    outputStatus: "ready",
    protectionValidation: "passed",
    retryCountLabel: "1",
    rawItem: null
  }
]

const waitingRows: BodyTranslationFieldResultItem[] = []

const runningRows: BodyTranslationFieldResultItem[] = rows.map((row) => ({
  ...row,
  outputStatus: "実行中",
  protectionValidation: "確認中",
  translatedText: row.translatedText
    ? `${row.translatedText}（生成中）`
    : "生成中"
}))

const completedRows: BodyTranslationFieldResultItem[] = rows.map((row) => ({
  ...row,
  outputStatus: "完了",
  protectionValidation: "passed"
}))

const meta = {
  title: "Screens/Job Run/TranslationCompletePage",
  component: TranslationCompletePage,
  parameters: {
    layout: "fullscreen"
  },
  args: {
    jobId: 1001,
    rows
  }
} satisfies Meta<typeof TranslationCompletePage>

export default meta

type Story = StoryObj<typeof meta>

export const Waiting: Story = {
  name: "待機中",
  args: {
    rows: waitingRows
  }
}

export const Ready: Story = {
  name: "準備完了"
}

export const Running: Story = {
  name: "実行中",
  args: {
    rows: runningRows
  }
}

export const Completed: Story = {
  name: "完了",
  args: {
    rows: completedRows
  }
}
