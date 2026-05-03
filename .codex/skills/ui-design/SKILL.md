---
name: ui-design
description: Codex 側の UI 設計作業プロトコル。UI 要件契約、task-local UIプロトタイプ、agent-browser 確認結果を固定する基準を提供する。
---
# UI Design

## 目的

`ui-design` は作業プロトコルである。
`designer` agent が UI を言葉だけで固定せず、UI 要件契約、task-local UIプロトタイプ、agent-browser 確認結果として扱うための、表示項目、操作、状態差分、導線、主要操作後の画面変化、UX 確認観点の見方を提供する。

実行境界、正本、引き継ぎ、停止 / 戻し は [design-bundle](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md) を参照する。

## 対応ロール

- `designer` が使う。
- 呼び出し元は `implement_lane` または人間とする。
- 返却先は 人間レビュー または `implement_lane` とする。
- 担当成果物は `ui-design` の出力規約で固定する。

## 入力規約

- task 内成果物: UI 要件契約、task-local UIプロトタイプ、agent-browser 確認結果の根拠にする設計成果物。
- 根拠参照: UI 判断の根拠にする要件、シナリオ、既存画面。
- 承認状態: 呼び出し元が渡す承認済みまたは未承認の状態。

## 外部参照規約

- エージェント実行定義と実行境界は [designer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.toml) に従う。
- 要件正本: [spec.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/spec.md) とする。
- architecture 正本: [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) とする。
- ER 正本: [er.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/er.md) と [diagrams/er](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/diagrams/er/) とする。
- 画面正本: [screen-design](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/README.md) とする。
- 上位シナリオ詳細仕様正本: [detail-specs](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/README.md) とする。
- scenario 正本: [scenario-tests](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/scenario-tests/README.md) とする。
- UI 設計雛形: [ui-design.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/ui-design/assets/ui-design.md)
- UIプロトタイプ確認サーバー: [serve-prototype.mjs](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/scripts/dev/serve-prototype.mjs)
- `agent-browser` 利用規約: [agent-browser.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/references/agent-browser.md)
- 実行定義 skill: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md)
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

### 拘束観点

- 画面表示文言、表示項目、主要操作、ボタン有効条件
- 画面区画、状態差分、配置制約、アクセシビリティ
- UIプロトタイプの種類、土台、task-local 正本配置
- UIプロトタイプの主要導線、情報密度、入力と結果の見え方
- UIプロトタイプ確認サーバーの URL、起動 command、人間確認中の起動状態
- `agent-browser` による task-local UIプロトタイプの UX 確認結果
- task-local UIプロトタイプ内の `mock-data/` とモックデータ移植禁止範囲
- 読み込み中、空、エラー、無効、進行中、再試行、成功
- デスクトップ / モバイルで破綻してはいけない条件と実装後確認観点

### UX 確認観点

| 観点 | 確認する内容 |
| --- | --- |
| 目的達成 | 主要利用者が開始操作から期待結果まで迷わず進めるかを確認する |
| 情報優先度 | 最初に見るべき情報、次に見るべき情報、補助情報の順序が画面構造で分かるかを確認する |
| 操作順 | 主要操作、戻る、取消、再試行、破壊的操作の位置と文言が作業順に沿うかを確認する |
| 状態理解 | 読み込み中、空、エラー、無効、進行中、成功の差分が見た目と文言で区別できるかを確認する |
| 回復可能性 | エラーや無効状態から、利用者が次に取れる操作を判断できるかを確認する |
| 表示文言 | 専門知識がなくても画面表示文言から次に何をするか分かるかを確認する |
| 入力負荷 | 入力量、選択肢、初期値、確認操作が目的に対して過剰でないかを確認する |
| 視線移動 | 見出し、ラベル、操作、結果表示が上下左右に散らばりすぎていないかを確認する |
| 表示幅追従 | desktop と mobile で文字、操作、状態表示が重ならず、主要導線が保たれるかを確認する |

### 表示文言変換例

| 固定名または内部状態名 | 画面表示文言 |
| --- | --- |
| `credential missing` | APIキーが未設定です |
| `dirty-validation` | 設定を変更したため、もう一度確認が必要です |
| `getModels failure` | モデル一覧を取得できませんでした |
| `Create job` は無効 | ジョブを作成できません |
| `Ready job` の read-only summary | 作成後の設定内容 |

## 判断規約

- UI は UI 要件契約で固定し、UIプロトタイプは task-local 確認用として扱う
- UIプロトタイプは `docs/exec-plans/active/<task-id>/` 配下を正本にする
- UIプロトタイプを作る場合は `prototype.svelte` として task folder に置く
- 既存画面変更では、既存画面または既存 UI 部品を土台にする
- 新規画面では、`docs/screen-design` の画面設計に従う
- UIプロトタイプ確認サーバーは `npm --prefix frontend run dev:prototype -- --task <task-id> --port 34116` で起動する
- UIプロトタイプは `http://127.0.0.1:34116/prototype` を `agent-browser` で開き、UX 確認観点の確認結果を `ui-design.md` に残す
- 人間確認中は UIプロトタイプ確認サーバーを起動したままにする
- 人間レビュー記録には確認 URL と起動 command を残す
- `agent-browser` 確認では `docs/references/agent-browser.md` に従い、`open`、`snapshot`、`errors`、`screenshot`、`close` を必要に応じて使う
- UIプロトタイプは task-local 確認用とし、docs 正本へ昇格しない
- UIプロトタイプは本番実装ではなく、人間承認後に implementation lane が反映する設計成果物として扱う
- 本番コードから UIプロトタイプを参照してはいけない
- UIプロトタイプのモックデータは `mock-data/` または `data-ui-prototype-sample-data-root` の範囲へ置き、frontend 実装へ移植してはいけない
- `mock-data/` 配下の値は状態表示確認用であり、product code、fixture、default state、test data へ移植してはいけない
- UIプロトタイプで対象にする動きは、入力反応、有効条件、タブ、詳細、モーダル、確認表示、読み込み中、空、エラー、成功の切り替えに限定する
- 汎用的な AI 風 UI や過剰な装飾を要求しない

- UI 契約 と シナリオ の責務を分ける
- デスクトップ と モバイル の破綻条件を実装後確認観点として残す
- 画面表示文言は日本語を優先する
- 画面表示文言は、固定名以外を日本語の業務語へ置き換える
- 内部状態名は画面に出さず、利用者の次操作を示す文へ変換する
- 英語ラベルを画面に出す場合は、利用者が設定画面で見る既存語だけに限定する
- `agent-browser` 確認後に表示文言レビューを必ず行う
- 表示文言レビューは、専門知識がなくても次に何をするか分かる表現水準かを判定する

## 非対象規約

- UI 不要 task、プロダクト frontend 実装、docs 正本反映だけの作業は扱わない。
- プロダクトコード実装と未承認 docs 正本化は扱わない。
- UIプロトタイプでは、実 API、永続化、本番 gateway 接続、業務ロジック完全再現を扱わない。
- frontend 側へ置けるものは、再利用可能な dev 専用起動基盤だけとする。
- 実装後に人間が確認すべき見た目調整を隠さない。

## 出力規約

- 判断結果: UI 要件契約の完了、未完了、停止の判定を返す。
- 根拠参照: UI 判断の根拠にした要件、シナリオ、既存画面を返す。
- UIプロトタイプ: task-local 確認用 `prototype.svelte` の作成または更新結果を返す。
- UIプロトタイプ種別: 既存画面変更または新規画面の区分を返す。
- UIプロトタイプ土台: 既存画面、既存 UI 部品、または `docs/screen-design` の参照を返す。
- UIプロトタイプ配置: task folder 内の配置パスを返す。
- UIプロトタイプ確認サーバー: 確認 URL、起動 command、人間確認中の起動要否を返す。
- 確認結果: `ui-design.md` の agent-browser 確認結果を返す。
- 操作確認結果: 主要操作後の画面変化、状態切り替え、未確認理由を返す。
- 表示文言レビュー結果: `agent-browser` 確認後に行った表示文言レビューの判定を返す。
- 不足情報: UI 要件契約を固定できない不足項目を返す。
- 次判断材料: `designer` または `implement_lane` が次を判断できる材料を返す。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコード変更の指示を含めない。

## 完了規約

- task 内成果物 が承認状態、根拠参照、未決事項を含んでいる。
- 人間レビュー が必要な判断を AI だけで完了扱いにしていない。
- 表示項目、主要操作、ボタン有効条件を確認した。
- UX 確認観点の各項目を `ui-design.md` に結果付きで残した。
- 状態、状態差分、表示幅追従、はみ出しリスク を実装後確認観点として確認した。
- `ui-design.md` は UI 要件契約と確認観点を含んでいる。
- task-local UIプロトタイプを作る場合は、主要区画、主要操作、状態差分を確認できる。
- task-local UIプロトタイプを作る場合は、主要操作後の画面変化を確認できる。
- task-local UIプロトタイプを作る場合は、`ui-design.md` に確認サーバーの URL、起動 command、人間確認中の起動要否を含めている。
- task-local UIプロトタイプを作る場合は、モックデータが `mock-data/` または `data-ui-prototype-sample-data-root` の範囲に置かれている。
- `ui-design.md` は `agent-browser` で確認した URL、起動 command、人間確認中の起動状態、画面サイズ、UX 確認観点ごとの結果、問題、未確認理由を含んでいる。
- `ui-design.md` は `agent-browser` 確認後の表示文言レビュー結果を含んでいる。
- 表示文言レビューは、固定名以外の画面表示文言が日本語の業務語になっているかを確認している。
- 表示文言レビューは、内部状態名が画面に出ず、利用者の次操作を示す文へ変換されているかを確認している。
- 表示文言レビューは、英語ラベルが利用者の設定画面で見る既存語だけに限定されているかを確認している。

## 停止規約

- UI が不要で `plan.md` の `ui_design` が `N/A` の時
- プロダクト frontend コードを実装する時
- docs 正本へ UI 仕様を反映するだけの時
- UIプロトタイプを `agent-browser` で確認できない場合は未実行理由を返して停止する。
- UIプロトタイプの土台が不明な場合は停止する。
- 本番コードから UIプロトタイプへの参照が必要になる場合は停止する。
- 停止時は不足項目、衝突箇所、戻し先を返す。
