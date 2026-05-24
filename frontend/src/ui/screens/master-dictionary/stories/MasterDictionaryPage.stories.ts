import type { Meta, StoryObj } from "@storybook/svelte-vite"
import type { MasterDictionaryScreenViewModel } from "@application/contract/master-dictionary/master-dictionary-screen-types"

import { createMasterDictionaryPageControllerFixture } from "../../__fixtures__/screen-page-controller-fixtures"
import MasterDictionaryPage from "../MasterDictionaryPage.svelte"

const waitingState = {
  hasStagedFile: false,
  importProgress: 0,
  importStatusText: "XML を選ぶと取込状態を確認できます。",
  importStatusValue: "待機中",
  importSummary: null,
  isImportRunning: false,
  selectedFileName: "未選択",
  selectedFileReference: null
} satisfies Partial<MasterDictionaryScreenViewModel>

const readyState = {
  hasStagedFile: true,
  importProgress: 0,
  importStatusText: "XML を取り込めます。",
  importStatusValue: "準備完了",
  importSummary: null,
  isImportRunning: false,
  selectedFileName: "synthetic-master-dictionary.xml",
  selectedFileReference: "storybook-master-dictionary-xml"
} satisfies Partial<MasterDictionaryScreenViewModel>

const runningState = {
  hasStagedFile: true,
  importProgress: 48,
  importStatusText: "XML を取り込み中です。",
  importStatusValue: "48%",
  importSummary: null,
  isImportRunning: true,
  selectedFileName: "synthetic-master-dictionary.xml",
  selectedFileReference: "storybook-master-dictionary-xml"
} satisfies Partial<MasterDictionaryScreenViewModel>

const completedState = {
  hasStagedFile: true,
  importProgress: 100,
  importStatusText: "XML 取り込みが完了しました。",
  importStatusValue: "完了",
  importSummary: {
    fileName: "synthetic-master-dictionary.xml",
    importedCount: 12,
    updatedCount: 4,
    totalCount: 48,
    selectedSource: "Dragon Priest"
  },
  isImportRunning: false,
  selectedFileName: "synthetic-master-dictionary.xml",
  selectedFileReference: "storybook-master-dictionary-xml"
} satisfies Partial<MasterDictionaryScreenViewModel>

const meta = {
  title: "Screens/Master Dictionary/MasterDictionaryPage",
  component: MasterDictionaryPage,
  parameters: {
    layout: "fullscreen"
  },
  args: {
    createController: createMasterDictionaryPageControllerFixture()
  }
} satisfies Meta<typeof MasterDictionaryPage>

export default meta

type Story = StoryObj<typeof meta>

export const Disconnected: Story = {}

export const Waiting: Story = {
  name: "待機中",
  args: {
    createController: createMasterDictionaryPageControllerFixture(waitingState)
  }
}

export const Ready: Story = {
  name: "準備完了",
  args: {
    createController: createMasterDictionaryPageControllerFixture(readyState)
  }
}

export const Running: Story = {
  name: "実行中",
  args: {
    createController: createMasterDictionaryPageControllerFixture(runningState)
  }
}

export const Completed: Story = {
  name: "完了",
  args: {
    createController: createMasterDictionaryPageControllerFixture(
      completedState
    )
  }
}
