# Storybook Usage Draft

- `status`: draft
- `scope`: Storybook の機能の使い方を決めるためのたたき台
- `canonicalization`: 後続 task で削ってから docs 正本へ反映する
- `references`:
  - Storybook docs: `writing-stories`
  - Storybook docs: `args`
  - Storybook docs: `actions`
  - Storybook docs: `decorators`

## 目的

Storybook は UI component を isolated state で確認する場所とする。
Storybook は backend、Wails runtime、Gateway、AI provider、DB、翻訳実行フローの代替にしない。

この文書は、story、fixture、action、decorator、controls の書き方を固定するための draft である。
プロダクト仕様や画面設計の正本ではない。

## Story の単位

1 story は、1 component の 1 表示状態を表す。
複数の状態を 1 story の中で切り替えない。

良い単位:

- `FixedProps`: 固定 props だけで標準表示を確認する。
- `Disabled`: 操作不可状態を確認する。
- `Warning`: 警告表示を確認する。
- `LongText`: 長い文言や折り返しを確認する。
- `Empty`: 選択肢なし、結果なし、未設定を確認する。

避ける単位:

- backend の成功から失敗までを story 内で再現する。
- Wails の generated binding を story から呼ぶ。
- 1 story 内で複数 screen の遷移を再現する。
- fixture を更新して業務フローの状態遷移を表す。

## Args

Storybook の `args` は component props の入力値として扱う。
`args` は UI の見た目と操作可能状態を変えるために使う。

この repo では、標準状態の `args` は component 横の fixture から渡す。

```ts
const meta = {
  title: "UI Components/AIModelSelectionCard",
  component: AIModelSelectionCard,
  args: aiModelSelectionCardFixture
} satisfies Meta<typeof AIModelSelectionCard>
```

個別 story は、標準 fixture との差分だけを `args` に書く。

```ts
export const Disabled: Story = {
  args: {
    providerDisabled: true,
    modelDisabled: true,
    actionButtonDisabled: true
  }
}
```

`args` に入れないもの:

- secret、API key、token
- 実ユーザーデータ
- backend DTO
- provider 応答原文
- file system path
- Wails runtime object

## Fixtures

fixture は component が受け取る props の固定値を置く。
fixture は component 横の `__fixtures__` に置く。

配置例:

```text
frontend/src/ui/components/AIModelSelectionCard.svelte
frontend/src/ui/components/AIModelSelectionCard.stories.ts
frontend/src/ui/components/__fixtures__/ai-model-selection-card-fixture.ts
```

fixture は view model か props の shape に合わせる。
fixture は backend DTO mock にしない。

fixture に入れてよいもの:

- 表示ラベル
- 選択肢
- disabled や selected などの表示状態
- 空配列、長文、警告文などの表示確認値
- no-op callback

fixture に入れないもの:

- gateway mock
- repository mock
- backend DTO mock
- generated `wailsjs`
- 外部 provider response
- 実行中 job の実データ

## Actions

Storybook の action は callback prop が呼ばれたことを確認するために使う。
action は業務処理の代替にしない。

標準では、見た目確認だけの story は no-op callback を fixture に置いてよい。
操作確認が必要な story では `@storybook/test` の `fn()` を使う。

```ts
import { fn } from "@storybook/test"

export const WithActions: Story = {
  args: {
    ...aiModelSelectionCardFixture,
    onProviderChange: fn(),
    onModelChange: fn(),
    onAction: fn()
  }
}
```

action で確認するもの:

- button click callback が呼ばれる。
- select change callback が呼ばれる。
- input callback が呼ばれる。

action で確認しないもの:

- backend に保存されたか。
- Wails command が呼ばれたか。
- provider API が成功したか。
- DB が更新されたか。

## Controls

Controls は props の表示調整と edge case 探索に使う。
Controls は仕様変更の根拠にしない。

Controls に出してよいもの:

- text label
- status tone
- disabled flag
- selected value
- warning 表示 flag

Controls に出さないもの:

- secret 本体
- file path
- provider credential
- backend object
- 実データの全文

`argTypes` は必要になった時だけ追加する。
追加する場合は、表示確認に必要な option set に限定する。

```ts
const meta = {
  component: AIModelSelectionCard,
  argTypes: {
    statusTone: {
      control: "select",
      options: ["neutral", "warning", "success"]
    }
  }
} satisfies Meta<typeof AIModelSelectionCard>
```

## Decorators

decorator は Storybook の表示環境を component の実利用環境に近づけるために使う。
decorator は production wiring を作らない。

global decorator に向いているもの:

- dark background
- shell padding
- CSS variables
- layout width

story local decorator に向いているもの:

- narrow width の確認
- modal や overlay の表示 container
- specific surface の確認

decorator に入れないもの:

- Gateway provider
- Wails runtime provider
- production controller factory
- backend mock server
- 実行フロー再現

## Story Naming

story 名は表示状態を短く表す。
画面や業務フローの説明を story 名に詰め込まない。

推奨:

- `FixedProps`
- `Disabled`
- `Warning`
- `LongText`
- `Empty`
- `Loading`

非推奨:

- `UserCanSelectProviderAndSaveSettings`
- `BackendReturnsError`
- `FullTranslationFlow`
- `WailsGatewayMocked`

## Review URL

人間レビュー用 URL は Storybook の story URL だけにする。
fakeAPI URL、Wails runtime URL、backend API URL は Storybook review URL にしない。

記録する項目:

- Storybook 起動 command
- review URL
- iframe URL
- story ID
- 確認状態
- 未確認理由
- build-storybook 結果

記録しない項目:

- command output 全文
- local absolute path
- secret、API key、token
- 実ユーザーデータ

## Test Boundary

Storybook は visual state の確認入口である。
unit test、scenario test、system test の代替ではない。

Storybook で確認する:

- component の表示状態
- props 差分による見た目
- callback が UI 操作で呼ばれること
- dark shell の中で読めること

Storybook で確認しない:

- service / usecase の correctness
- repository の保存結果
- Wails bridge の接続結果
- provider API の通信結果
- translation job の状態遷移

## 禁止事項

- story から `generated wailsjs` を import しない。
- story から Gateway を import しない。
- story から backend DTO mock を import しない。
- fixture に secret や token を入れない。
- story を業務フロー再現に使わない。
- Storybook 専用 utility を production code から import しない。

## 初期テンプレート

```ts
import type { Meta, StoryObj } from "@storybook/svelte-vite"
import { fn } from "@storybook/test"
import TargetComponent from "./TargetComponent.svelte"
import { targetComponentFixture } from "./__fixtures__/target-component-fixture"

const meta = {
  title: "UI Components/TargetComponent",
  component: TargetComponent,
  args: targetComponentFixture
} satisfies Meta<typeof TargetComponent>

export default meta

type Story = StoryObj<typeof meta>

export const FixedProps: Story = {}

export const WithActions: Story = {
  args: {
    ...targetComponentFixture,
    onChange: fn()
  }
}
```

## 削る候補

- `argTypes` は、Controls を積極利用しないなら削る。
- `story local decorator` は、実例が出るまで削る。
- `Review URL` は別 docs に分離してもよい。
