# Storybook 規約

この文書は Storybook の作成、起動、分類、確認資源、`fixture` の共通規約である。
この文書は Storybook をどう作るか、どの分類と確認資源で使うかだけを扱う。

## 起動

- 固定 URL: Storybook は `http://localhost:6008/` を標準 URL にする。
- 起動 command: Storybook は `npm --prefix frontend run storybook` で起動する。
- port 固定: Storybook は `6008` だけを使い、別 port で追加起動しない。
- 変更反映: Storybook 確認中に frontend または story を変更した場合は、既存 Storybook を停止し、同じ command で再起動する。
- port 使用中: `6008` が使用中の場合は、別 port に逃がさず、既存 Storybook を停止してから再起動する。
- build 確認: Storybook 確認資源を追加または更新した場合は、`npm --prefix frontend run build-storybook` の結果または未実行理由を残す。

## 分類

| 対象 | 作業中分類 | 通常分類 | 使う条件 |
| --- | --- | --- | --- |
| 画面 | `Review/Changed Screens/<画面名>/<状態名>` | `Screens/<画面名>/<状態名>` | 画面全体、主要導線、複数コンポーネントの配置、画面単位の状態を確認する場合。 |
| コンポーネント | `Review/Changed Components/<コンポーネント名>/<状態名>` | `UI Components/<コンポーネント名>/<状態名>` | 再利用コンポーネント、コンポーネント単体の表示状態、コンポーネント内部の状態差分を確認する場合。 |

画面上でしか意味が確定しない変更は、画面の分類を使う。
コンポーネント単体で意味が閉じる変更は、コンポーネントの分類を使う。
画面とコンポーネントの両方を確認する必要がある変更は、画面の story とコンポーネントの story を分ける。
作業中または差し戻し対応中に変更した story は、作業中分類へ置く。
完了後に確定した story は、通常分類へ戻す。
返却材料には、作業中分類、通常分類、現在分類を分けて書く。

## 確認資源

| 確認資源 | 意味 |
| --- | --- |
| 変更画面 | 今回の frontend 変更で表示、導線、配置、画面単位の状態が変わった画面。 |
| 追加画面 | 今回の frontend 変更で追加した画面。 |
| 変更コンポーネント | 今回の frontend 変更で表示または振る舞いが変わったコンポーネント。 |
| 追加コンポーネント | 今回の frontend 変更で追加したコンポーネント。 |
| 変更表示状態 | 今回の frontend 変更で表示、遷移、操作可否、エラー、空状態、読み込み、実行中、完了が変わった状態。 |
| 追加表示状態 | 今回の frontend 変更で追加した表示状態。 |
| story | 変更画面、追加画面、変更コンポーネント、追加コンポーネント、変更表示状態、追加表示状態を Storybook で確認する入口。 |
| `fixture` | story が固定表示に使う入力値、前提状態、表示用データ。 |
| 関連資源 | story の表示に必要な props 変換、定数、表示専用の補助資源。 |

Storybook の story は固定 props または固定 `fixture` で表示できる状態にする。
Storybook 確認資源は frontend 表示と操作確認用であり、backend 実装、統合境界実装、永続化仕様の代替にしない。

## fixture 種類基準

| `fixture` 種類 | 使う条件 | 含める内容 |
| --- | --- | --- |
| props `fixture` | コンポーネントの表示差分を props だけで再現できる場合。 | props、表示用ラベル、選択肢、callback の placeholder。 |
| 画面入力 `fixture` | 画面全体、主要導線、複数コンポーネントの配置、画面単位の操作可否を確認する場合。 | 画面 view model、選択対象、画面状態、導線上の前提状態。 |
| 一覧 `fixture` | 一覧、ページング、選択、空状態、大量件数の表示範囲が変わる場合。 | 行データ、選択ID、ページ、表示件数、画面に表示する総件数。 |
| フォーム `fixture` | 入力値、未保存、検証エラー、保存中、保存済み、操作可否が変わる場合。 | 入力値、field error、dirty 状態、submitting 状態、保存結果、disabled 理由。 |
| 実行状態 `fixture` | job、phase、import、generation など非同期処理の状態で表示、進捗、操作可否が変わる場合。 | 実行状態、進捗、running、cancel、retry、error、complete、表示 message。 |
| 境界結果 `fixture` | gateway または usecase の結果を画面表示用データとして固定する場合。 | 変換後 view model、表示用 DTO 相当データ、成功結果、失敗結果。 |
| 表示環境 `fixture` | viewport、theme、locale など表示環境で見た目や文言が実質的に変わる場合。 | viewport、theme、locale、表示環境に依存する前提値。 |

画面 fixture は、画面単位の判断に必要な代表状態だけを作る。
状態ラベル、バッジ文言、単一フィールド表示だけが変わり、配置、操作可否、導線、表示量が変わらない場合は、画面 fixture を増やさず、該当コンポーネントの story で扱う。
コンポーネント story で変化を確認できる状態は、同じ状態を画面 story へ重複追加しない。
全状態の組み合わせ表を作らず、画面判断またはコンポーネント判断に必要な状態だけを置く。
