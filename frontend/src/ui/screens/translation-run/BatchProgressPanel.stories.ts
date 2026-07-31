import type { Meta, StoryObj } from "@storybook/svelte-vite"
import BatchProgressPanel from "./BatchProgressPanel.svelte"

// OpenAI と xAI に共通する batch の進行状況パネル。固有名 → 本文 → 完了 のステッパーで現在地を見せ、現段 batch の件数を出す。
const meta = {
  title: "UI Components/進行状況パネル",
  component: BatchProgressPanel,
  parameters: {
    layout: "padded"
  }
} satisfies Meta<typeof BatchProgressPanel>

export default meta
type Story = StoryObj<typeof meta>

// 状態未確認（状態確認前）。ステッパーは中立、件数は出さず状態確認を促す。
export const Unchecked: Story = {
  name: "未確認",
  args: { progress: undefined }
}

// 固有名段が処理中。ステッパー現在地=固有名、処理待ちが残る。
export const ProperProcessing: Story = {
  name: "固有名処理中",
  args: {
    progress: {
      stage: "proper",
      total: 12,
      pending: 5,
      succeeded: 7,
      failed: 0,
      canApply: false
    }
  }
}

// 固有名段が完了（処理待ち 0）。取り込んで本文へ進める。
export const ProperReady: Story = {
  name: "固有名完了",
  args: {
    progress: {
      stage: "proper",
      total: 12,
      pending: 0,
      succeeded: 12,
      failed: 0,
      canApply: true
    }
  }
}

// 本文段が処理中。ステッパー現在地=本文、固有名段は完了済み。
export const BodyProcessing: Story = {
  name: "本文処理中",
  args: {
    progress: {
      stage: "body",
      total: 113,
      pending: 2,
      succeeded: 111,
      failed: 0,
      canApply: false
    }
  }
}

// 本文段が完了（処理待ち 0）。取り込んで完了へ進める。
export const BodyReady: Story = {
  name: "本文完了",
  args: {
    progress: {
      stage: "body",
      total: 113,
      pending: 0,
      succeeded: 113,
      failed: 0,
      canApply: true
    }
  }
}

// 本文段が完了し、成功と失敗が混在する。成功分を取り込み、失敗分を未訳として残せる。
export const BodyReadyWithFailures: Story = {
  name: "本文完了（失敗あり）",
  args: {
    progress: {
      stage: "body",
      total: 113,
      pending: 0,
      succeeded: 110,
      failed: 3,
      canApply: true
    }
  }
}

// すべて完了（段=完了）。全段に完了マークが付く。
export const AllDone: Story = {
  name: "全完了",
  args: {
    progress: {
      stage: "done",
      total: 113,
      pending: 0,
      succeeded: 113,
      failed: 0,
      canApply: false
    }
  }
}
