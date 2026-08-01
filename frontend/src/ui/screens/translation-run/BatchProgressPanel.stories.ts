import type { Meta, StoryObj } from "@storybook/svelte-vite"
import BatchProgressPanel from "./BatchProgressPanel.svelte"

// 処理中、再開待ち、完了を区別する batch の進行状況パネル。
// Storybook 人間レビュー中は作業中分類に置く。承認後は UI Components/進行状況パネルへ戻す。
const meta = {
  title: "UI Components/進行状況パネル",
  component: BatchProgressPanel,
  tags: ["autodocs"],
  parameters: {
    layout: "padded"
  }
} satisfies Meta<typeof BatchProgressPanel>

export default meta
type Story = StoryObj<typeof meta>

// 保存済みの進行がない。主操作から処理を開始することを案内する。
export const NotStarted: Story = {
  name: "未開始",
  args: { progress: undefined, running: false },
  parameters: {
    docs: {
      description: {
        story: `### 前提条件

- 保存済みの batch 進行がなく、処理も始まっていない。

### 期待値

- 状態は「未開始」になる。
- ステッパーは全段を中立表示にする。
- 「バッチ実行」を押すと処理を開始する案内を表示する。`
      }
    }
  }
}

// batch の開始応答を待っている。進行段を得る前から自動処理中と分かる。
export const Starting: Story = {
  name: "開始中",
  args: { progress: undefined, running: true },
  parameters: {
    docs: {
      description: {
        story: `### 前提条件

- batch の開始応答を待っており、進行段はまだ取得していない。

### 期待値

- 状態は「開始中」になる。
- ステッパーは全段を中立表示にする。
- 開始すると固有名の処理へ進む案内を表示する。`
      }
    }
  }
}

// 固有名段を処理している。
export const ProperRunning: Story = {
  name: "固有名・処理中",
  args: {
    running: true,
    progress: {
      stage: "proper",
      total: 1000,
      pending: 742,
      succeeded: 258,
      failed: 0,
      canApply: false,
      untranslatedCount: 0
    }
  },
  parameters: {
    docs: {
      description: {
        story: `### 前提条件

- 固有名段に処理待ちが残り、処理を継続している。

### 期待値

- 状態は「処理中」になり、ステッパーは固有名を現在段として表示する。
- 総数、処理待ち、成功、失敗の件数を表示する。
- 完了すると次の処理へ自動で進む案内を表示する。`
      }
    }
  }
}

// 固有名段が完了した。人の取り込みを促さず、次の処理へ自動で進むことを示す。
export const ProperReady: Story = {
  name: "固有名・次へ進行中",
  args: {
    running: true,
    progress: {
      stage: "proper",
      total: 1000,
      pending: 0,
      succeeded: 998,
      failed: 2,
      canApply: true,
      untranslatedCount: 0
    }
  },
  parameters: {
    docs: {
      description: {
        story: `### 前提条件

- 固有名段が完了し、自動処理が結果の取り込みと本文段への移行を続けている。

### 期待値

- 状態は「処理中」になり、固有名段の完了件数を表示する。
- 人の取り込み操作を促さない。
- 結果を取り込み、次の処理へ自動で進む案内を表示する。`
      }
    }
  }
}

// 本文段の途中で処理が止まり、人の再開操作を待っている。
export const Paused: Story = {
  name: "本文・再開待ち",
  args: {
    running: false,
    progress: {
      stage: "body",
      total: 1000,
      pending: 412,
      succeeded: 586,
      failed: 2,
      canApply: false,
      untranslatedCount: 0
    }
  },
  parameters: {
    docs: {
      description: {
        story: `### 前提条件

- 本文段に処理待ちが残り、処理が止まっている。

### 期待値

- 状態は「再開待ち」になり、ステッパーは本文を現在段として表示する。
- 現在の件数を保持して表示する。
- 「バッチ実行を再開」で続ける案内を表示する。`
      }
    }
  }
}

// 固有名段と本文段が完了した。
export const Done: Story = {
  name: "完了",
  args: {
    running: false,
    progress: {
      stage: "done",
      total: 113,
      pending: 0,
      succeeded: 113,
      failed: 0,
      canApply: false,
      untranslatedCount: 0
    }
  },
  parameters: {
    docs: {
      description: {
        story: `### 前提条件

- 固有名段と本文段が完了している。

### 期待値

- 状態は「完了」になる。
- ステッパーは全段に完了マークを表示する。
- すべての翻訳が完了した案内を表示する。`
      }
    }
  }
}
