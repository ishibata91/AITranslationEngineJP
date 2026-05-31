import type { Meta, StoryObj } from "@storybook/svelte-vite"

import TermTranslationPhasePanel from "../TermTranslationPhasePanel.svelte"
import { termTranslationPhasePanelFixture } from "../__fixtures__/term-phase-card-fixture"

const meta = {
  title: "Screen Components/Term Translation Phase/TermTranslationPhasePanel",
  component: TermTranslationPhasePanel,
  parameters: {
    layout: "fullscreen"
  },
  args: {
    viewModel: termTranslationPhasePanelFixture,
    onAction: () => {}
  }
} satisfies Meta<typeof TermTranslationPhasePanel>

export default meta

type Story = StoryObj<typeof meta>

/**
 * 初回取得中ローディングレイヤー表示（initialFetchDone=false）。
 * 処理対象一覧領域の最前面にオーバーレイが表示され、検索・ページ操作が排他される状態。
 * selector: term-translation-phase-processing-target-loading
 */
export const ProcessingTargetLoading: Story = {
  args: {
    initialFetchDone: false
  }
}

/**
 * 初回取得完了後（initialFetchDone=true）、件数あり状態。
 * ローディングレイヤーが消え、処理対象一覧が表示される状態。
 */
export const ProcessingTargetWithItems: Story = {
  args: {
    initialFetchDone: true,
    processingTargetPageState: {
      items: [
        {
          id: "term:1",
          name: "FULL",
          detail: "原文: Hello",
          titleParts: [{ text: "原文: Hello" }],
          metadata: [{ label: "AI 訳語候補", value: "こんにちは" }]
        },
        {
          id: "term:2",
          name: "NAME",
          detail: "原文: World",
          titleParts: [{ text: "原文: World" }],
          metadata: [{ label: "AI 訳語候補", value: "世界" }]
        }
      ],
      metadata: [],
      page: 1,
      pageSize: 50,
      totalCount: 2,
      searchQuery: "",
      busy: false
    }
  }
}

/**
 * 初回取得完了後（initialFetchDone=true）、空状態（母数0件）。
 * ローディングレイヤーが消え、空状態が表示される状態。
 * selector: term-translation-phase-processing-target-empty
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
