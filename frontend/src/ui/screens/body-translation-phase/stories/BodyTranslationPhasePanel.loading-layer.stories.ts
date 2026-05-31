import type { Meta, StoryObj } from "@storybook/svelte-vite"

import BodyTranslationPhasePanel from "../BodyTranslationPhasePanel.svelte"
import { bodyTranslationPhasePanelFixture } from "../__fixtures__/body-phase-card-fixture"

const meta = {
  title: "Screen Components/Body Translation Phase/BodyTranslationPhasePanel",
  component: BodyTranslationPhasePanel,
  parameters: {
    layout: "fullscreen"
  },
  args: {
    viewModel: bodyTranslationPhasePanelFixture,
    onAction: () => {}
  }
} satisfies Meta<typeof BodyTranslationPhasePanel>

export default meta

type Story = StoryObj<typeof meta>

/**
 * 初回取得中ローディングレイヤー表示（initialFetchDone=false）。
 * selector: body-translation-phase-processing-target-loading
 */
export const ProcessingTargetLoading: Story = {
  args: {
    initialFetchDone: false
  }
}

/**
 * 初回取得完了後（initialFetchDone=true）、件数あり状態。
 */
export const ProcessingTargetWithItems: Story = {
  args: {
    initialFetchDone: true,
    processingTargetPageState: {
      items: [
        {
          id: "field:1",
          name: "FULL",
          detail: "原文: The ancient city of Markarth",
          titleParts: [{ text: "FULL" }],
          metadata: [{ label: "訳語", value: "古代都市マルカルス" }]
        }
      ],
      metadata: [],
      page: 1,
      pageSize: 50,
      totalCount: 1,
      searchQuery: "",
      busy: false
    }
  }
}

/**
 * 初回取得完了後（initialFetchDone=true）、空状態（母数0件）。
 * selector: body-translation-phase-processing-target-empty
 */
export const ProcessingTargetEmpty: Story = {
  args: {
    initialFetchDone: true,
    processingTargetPageState: {
      items: [],
      metadata: [],
      page: 1,
      pageSize: 50,
      totalCount: 0,
      searchQuery: "",
      busy: false
    }
  }
}
