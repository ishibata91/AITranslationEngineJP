import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TranslationRunScreen from "./TranslationRunScreen.svelte"
import {
  OPENAI_READY_STATE,
  OPENAI_RUNNING_STATE,
  OPENAI_PAUSED_STATE,
  OPENAI_NO_UNTRANSLATED_STATE,
  OPENAI_BATCH_UNTRANSLATED_STATE,
  OPENAI_FAILED_STATE
} from "./translation-run.fixtures"

// 状態確認と手動取り込みを削除し、一つの主操作で自動処理する batch 画面の代表状態。
// Storybook 人間レビュー中は作業中分類に置く。承認後は Screens/翻訳実行へ戻す。
const meta = {
  title: "Screens/翻訳実行",
  component: TranslationRunScreen,
  tags: ["autodocs"],
  parameters: {
    layout: "fullscreen"
  },
  args: {
    onFieldInput: () => {},
    onLoadModels: () => {},
    onRun: () => {},
    onPagePrev: () => {},
    onPageNext: () => {},
    onUntranslatedOnlyChange: () => {},
    onProviderChange: () => {},
    onSubmit: () => {}
  }
} satisfies Meta<typeof TranslationRunScreen>

export default meta
type Story = StoryObj<typeof meta>

// 保存済みの進行がない初回表示。「バッチ実行」だけを表示する。
export const NotStarted: Story = {
  name: "未開始",
  args: { ...OPENAI_READY_STATE },
  parameters: {
    docs: {
      description: {
        story: `### 前提条件

- OpenAI の接続情報が揃い、保存済みの batch 進行がない。

### 期待値

- 主操作は「バッチ実行」で有効になる。
- 状態確認ボタンと手動取り込みボタンは表示しない。
- 進行状況パネルは「未開始」と、「バッチ実行」を押すと処理を開始する案内を表示する。`
      }
    }
  }
}

// 開始処理または自動状態確認中。spinner 付き「実行中…」を無効表示する。
export const Running: Story = {
  name: "実行中",
  args: { ...OPENAI_RUNNING_STATE },
  parameters: {
    docs: {
      description: {
        story: `### 前提条件

- 固有名段の batch を処理している。

### 期待値

- 主操作は spinner 付き「実行中…」で無効になる。
- 状態表示は「実行中」になる。
- 状態確認ボタンと手動取り込みボタンは表示しない。
- 進行状況パネルは「処理中」と、完了すると次の処理へ自動で進む案内を表示する。`
      }
    }
  }
}

// 画面を閉じた後などで処理が止まった状態。「バッチ実行を再開」を表示する。
export const Paused: Story = {
  name: "途中停止",
  args: { ...OPENAI_PAUSED_STATE },
  parameters: {
    docs: {
      description: {
        story: `### 前提条件

- 本文段の保存済み進行があり、処理が止まっている。

### 期待値

- 主操作は「バッチ実行を再開」で有効になる。
- 状態確認ボタンと手動取り込みボタンは表示しない。
- 進行状況パネルは「再開待ち」と、主操作で続ける案内を表示する。`
      }
    }
  }
}

// 全件の翻訳が完了した状態。「完了」を無効表示する。
export const Done: Story = {
  name: "完了（未訳なし）",
  args: { ...OPENAI_NO_UNTRANSLATED_STATE },
  parameters: {
    docs: {
      description: {
        story: `### 前提条件

- 固有名段と本文段が完了し、未訳が残っていない。

### 期待値

- 主操作は「完了」で無効になる。
- 状態確認ボタンと手動取り込みボタンは表示しない。
- 進行状況パネルは「完了」と、すべての翻訳が完了した案内を表示する。`
      }
    }
  }
}

// 未訳が残った完了状態。「未訳だけを再送信」を表示する。
export const DoneWithUntranslated: Story = {
  name: "完了（未訳あり）",
  args: { ...OPENAI_BATCH_UNTRANSLATED_STATE },
  parameters: {
    docs: {
      description: {
        story: `### 前提条件

- 固有名段と本文段が完了し、未訳が3件残っている。

### 期待値

- 主操作は「未訳だけを再送信」で有効になる。
- 状態確認ボタンと手動取り込みボタンは表示しない。
- 進行状況パネルは「完了」を表示し、画面は未訳件数を案内する。`
      }
    }
  }
}

// 外部 batch の失敗理由を表示し、保存済みの進行を再開できる状態。
export const Failed: Story = {
  name: "失敗表示",
  args: { ...OPENAI_FAILED_STATE },
  parameters: {
    docs: {
      description: {
        story: `### 前提条件

- 保存済みの本文段で外部 batch が失敗し、処理が止まっている。

### 期待値

- 主操作は「バッチ実行を再開」で有効になる。
- 状態確認ボタンと手動取り込みボタンは表示しない。
- 外部 batch ID と失敗理由を表示し、進行状況パネルは「再開待ち」を表示する。`
      }
    }
  }
}
