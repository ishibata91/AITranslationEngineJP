import type { Meta, StoryObj } from "@storybook/svelte-vite"
import type {
  MasterPersonaPreviewStateEntry,
  MasterPersonaScreenViewModel
} from "@application/gateway-contract/master-persona"

import { createMasterPersonaPageControllerFixture } from "../../__fixtures__/screen-page-controller-fixtures"
import MasterPersonaPage from "../MasterPersonaPage.svelte"

const previewFileName = "sample-persona-source.json"
const previewTargetPlugin = "SamplePersonaPlugin.esp"
const previewExistingCount = 32

const preview: MasterPersonaPreviewStateEntry = {
  fileName: previewFileName,
  targetPlugin: previewTargetPlugin,
  candidateCount: 128,
  newlyAddableCount: 96,
  existingCount: previewExistingCount,
  status: "previewed"
}

const waitingState = {
  canStartGeneration: false,
  hasPreview: false,
  isRunActive: false,
  preview: null,
  progressPercent: 0,
  runStatus: {
    runState: "待機中",
    targetPlugin: "",
    processedCount: 0,
    successCount: 0,
    existingSkipCount: 0,
    zeroDialogueSkipCount: 0,
    genericNpcCount: 0,
    currentActorLabel: "",
    message: "JSON を選ぶと生成状態を確認できます。"
  },
  selectedFileName: "未選択",
  selectedFileReference: null
} satisfies Partial<MasterPersonaScreenViewModel>

const readyState = {
  canStartGeneration: true,
  hasPreview: true,
  isRunActive: false,
  preview,
  progressPercent: 0,
  runStatus: {
    runState: "準備完了",
    targetPlugin: previewTargetPlugin,
    processedCount: 0,
    successCount: 0,
    existingSkipCount: previewExistingCount,
    zeroDialogueSkipCount: 0,
    genericNpcCount: 0,
    currentActorLabel: "",
    message: "ペルソナ作成を開始できます。"
  },
  selectedFileName: previewFileName,
  selectedFileReference: "storybook-master-persona-json"
} satisfies Partial<MasterPersonaScreenViewModel>

const runningState = {
  canStartGeneration: false,
  hasPreview: true,
  isRunActive: true,
  preview,
  progressPercent: 54,
  runStatus: {
    runState: "実行中",
    targetPlugin: previewTargetPlugin,
    processedCount: 52,
    successCount: 48,
    existingSkipCount: previewExistingCount,
    zeroDialogueSkipCount: 2,
    genericNpcCount: 2,
    currentActorLabel: "SAMPLE_NPC_052",
    message: "ペルソナを作成しています。"
  },
  selectedFileName: previewFileName,
  selectedFileReference: "storybook-master-persona-json"
} satisfies Partial<MasterPersonaScreenViewModel>

const completedState = {
  canStartGeneration: false,
  hasPreview: true,
  isRunActive: false,
  preview,
  progressPercent: 100,
  runStatus: {
    runState: "完了",
    targetPlugin: previewTargetPlugin,
    processedCount: 96,
    successCount: 96,
    existingSkipCount: previewExistingCount,
    zeroDialogueSkipCount: 0,
    genericNpcCount: 0,
    currentActorLabel: "",
    message: "ペルソナ作成が完了しました。"
  },
  selectedFileName: previewFileName,
  selectedFileReference: "storybook-master-persona-json"
} satisfies Partial<MasterPersonaScreenViewModel>

const meta = {
  title: "Screens/Master Persona/MasterPersonaPage",
  component: MasterPersonaPage,
  parameters: {
    layout: "fullscreen"
  },
  args: {
    createController: createMasterPersonaPageControllerFixture()
  }
} satisfies Meta<typeof MasterPersonaPage>

export default meta

type Story = StoryObj<typeof meta>

export const Disconnected: Story = {}

export const Waiting: Story = {
  name: "待機中",
  args: {
    createController: createMasterPersonaPageControllerFixture(waitingState)
  }
}

export const Ready: Story = {
  name: "準備完了",
  args: {
    createController: createMasterPersonaPageControllerFixture(readyState)
  }
}

export const Running: Story = {
  name: "実行中",
  args: {
    createController: createMasterPersonaPageControllerFixture(runningState)
  }
}

export const Completed: Story = {
  name: "完了",
  args: {
    createController: createMasterPersonaPageControllerFixture(completedState)
  }
}
