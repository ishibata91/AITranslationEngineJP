import type { Meta, StoryObj } from "@storybook/svelte-vite"
import type {
  TranslationInputScreenViewModel,
  TranslationInputStagedFile
} from "@application/gateway-contract/translation-input"

import InputReviewPage from "../InputReviewPage.svelte"
import {
  baseReviewItem,
  createInputReviewPageControllerFixture,
  inputReviewPageSelectedViewModel
} from "../__fixtures__/translation-input-panel-fixtures"

const stagedFile: TranslationInputStagedFile = {
  fileName: "sample-translation-input.json",
  filePath: "選択済み JSON",
  fileHash: "sample-hash-20260518"
}

const waitingState = {
  items: [],
  selectedItemId: null,
  stagedFile: null,
  operationState: "idle",
  selectedItem: null,
  hasStagedFile: false,
  canImport: false,
  isImporting: false,
  stagedFileName: "未選択",
  stagedFilePath: "-",
  stagedFileHash: "-",
  operationStatusLabel: "待機中",
  operationStatusText: "JSON を選ぶと読み込み準備を開始できます。",
  latestOutcomeTitle: "未選択",
  latestOutcomeText: "読み込み済みデータを選んでください。",
  selectionStatusText: "読み込み済みデータを選んでください。",
  totalItemCountLabel: "0 件",
  emptyStateText: "読み込み済みデータはありません。",
  canRebuildSelected: false
} satisfies Partial<TranslationInputScreenViewModel>

const readyState = {
  items: [],
  selectedItemId: null,
  stagedFile,
  operationState: "ready",
  selectedItem: null,
  hasStagedFile: true,
  canImport: true,
  isImporting: false,
  stagedFileName: stagedFile.fileName,
  stagedFilePath: stagedFile.filePath,
  stagedFileHash: stagedFile.fileHash,
  operationStatusLabel: "準備完了",
  operationStatusText: "選択した JSON を登録できます。",
  latestOutcomeTitle: "登録前",
  latestOutcomeText: "登録後に内容を確認できます。",
  selectionStatusText: "読み込み済みデータはまだありません。",
  totalItemCountLabel: "0 件",
  emptyStateText: "読み込み済みデータはありません。",
  canRebuildSelected: false
} satisfies Partial<TranslationInputScreenViewModel>

const runningState = {
  ...readyState,
  operationState: "importing",
  canImport: false,
  isImporting: true,
  operationStatusLabel: "実行中",
  operationStatusText: "JSON を登録しています。"
} satisfies Partial<TranslationInputScreenViewModel>

const completedState = {
  ...inputReviewPageSelectedViewModel,
  operationState: "idle",
  operationStatusLabel: "完了",
  operationStatusText: "入力データの登録が完了しました。",
  latestOutcomeTitle: "登録完了",
  latestOutcomeText: "次に単語翻訳へ進めます。",
  selectedItem: {
    ...baseReviewItem,
    status: "registered"
  },
  selectedItemId: baseReviewItem.localId
} satisfies Partial<TranslationInputScreenViewModel>

const meta = {
  title: "Screens/Translation Input/InputReviewPage",
  component: InputReviewPage,
  parameters: {
    layout: "fullscreen"
  },
  args: {
    createController: createInputReviewPageControllerFixture(),
    onOpenJobRun: () => {}
  }
} satisfies Meta<typeof InputReviewPage>

export default meta

type Story = StoryObj<typeof meta>

export const SelectedInputReady: Story = {}

export const Waiting: Story = {
  name: "待機中",
  args: {
    createController: createInputReviewPageControllerFixture(waitingState)
  }
}

export const Ready: Story = {
  name: "準備完了",
  args: {
    createController: createInputReviewPageControllerFixture(readyState)
  }
}

export const Running: Story = {
  name: "実行中",
  args: {
    createController: createInputReviewPageControllerFixture(runningState)
  }
}

export const Completed: Story = {
  name: "完了",
  args: {
    createController: createInputReviewPageControllerFixture(completedState)
  }
}
