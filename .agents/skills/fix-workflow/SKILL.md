---
name: fix-workflow
description: "修正フローの入口オーケストレーター。人間が確認した不具合、レビュー非通過、検証失敗の観測記録から、作業 branch と plan.md を固定し、観測ログ駆動の再現確認と原因究明を investigation.md に、どう直すかの修正方針を design.md に、修正後の確定仕様（要求ごとの仕様）を spec.md にまとめ、design-review（design_reviewer agent によるソース照合検証と記述検証）と人間修正レビューを通してから implementation-workflow へ渡す。TRIGGER when: 修正系 task（バグ修正・refactor）の入口。SKIP when: 仕様変更や機能追加が必要と判明した場合は本モジュールを停止し、design-workflow へ迂回する。"
---
# Fix Workflow

## 目的

`fix-workflow` は、修正フローの入口オーケストレーターである。人間が確認した不具合、レビュー非通過、検証失敗の観測記録から、作業 branch と `plan.md` を固定する。再現確認と原因究明は `investigation.md`、as-is と to-be の方針は `design.md`、修正後の確定仕様は `spec.md`、未決定事項とブロッカーは `pending.md`、参照は `references.md` に分けて Codex 本体が固定する。
`pending.md` が空の場合だけ、`design-review` と人間修正レビューを通して `implementation-workflow` へ渡す。`design.md` と `spec.md` は `plan.md` の要求ごとの節に分ける。両者が食い違う場合は `spec.md` を優先する。
再現確認・原因究明（investigation）と、どう直すかの設計（design）は責務を分ける。investigation は「何が起きてなぜか」を確定し、design は「どう直すか」だけを扱う。
修正は fail-test ベースで進める前提で、先に不具合を検出できるテスト観点も引き継ぎに含める。

## 呼び出し関係

- 呼び出し元: 人間、または上位 agent。
- 返却先: 呼び出し元。
- モジュールが呼ぶ下位 skill: `fix-decision`（原因究明の判断基準）。`design-protocol`（設計作成）。`specification-protocol`（仕様作成）。`design-review`（設計と仕様の検証）。`presentation`（人間修正レビュー材料と、現況とあるべき形の 2 図の作成）。
- モジュールが呼ぶ下位 agent: `design_reviewer`（`fresh`）。人間修正レビューの前に `design.md` と `spec.md` を実ソースと照合して検証する。調査、設計、仕様の本文は Codex 本体が task 文脈を持ったまま書く。
- 下流の段: `implementation-workflow`、`finalization-module`。

## 入口条件

- 呼び出し元から依頼要約と `task-id` が渡されている。
- 統合先 branch（既定 `master`）が local に存在する。
- 人間観測記録（人間が見た画面、操作、ログ、失敗、期待との差分）が固定されている。
- 修正対象の Wails 接続対象（プロセスまたは接続先）が単一化されている。
- ブラウザ操作（`chrome-devtools` MCP ツール群、`mcp__plugin_chrome-devtools-mcp_chrome-devtools__*`）が利用できる。

## 出口条件

- 作業 branch が `codex/<task-id>` として local に存在する。
- 作業計画フォルダ `docs/exec-plans/active/<task-id>/` が存在する。
- `plan.md` に対象 task の合意済み要求が記録されている。
- `investigation.md` に観測済み問題、画面再現確認、原因仮説、観測ログ検証、確定原因が記録されている。
- `design.md` の各要求の節に as-is と to-be の方針が記録されている。
- `spec.md` の各要求の節に仕様が記録され、各仕様に前提条件と確かめ方が書かれている。
- `pending.md` が空である。
- `references.md` に source と外部資料の所在が記録されている。
- `design-review` が通過済み。
- `人間修正レビュー` 承認済み。
- 仕様変更または仕様追加が必要と判断された場合は停止して呼び出し元へ戻す（`design-workflow` への迂回が必要か、人間判断を仰ぐ）。

## 担当成果物

| 成果物ID | 担当 | 依存対象 |
| --- | --- | --- |
| `branch 準備` | Codex 本体 | 入口条件 |
| `plan.md` | Codex 本体 | `branch 準備` |
| `人間観測記録` | Codex 本体 | `plan.md` |
| `調査` | Codex 本体 | `人間観測記録` |
| `設計と仕様` | Codex 本体 | `調査` |
| `design-review` | `design_reviewer`（`fresh`） | `設計と仕様` |
| `人間修正レビュー` | 人間 | `design-review` |
| `実装への引き継ぎ` | Codex 本体 | `人間修正レビュー` |

## branch 準備

- 作業 branch を既定名 `codex/<task-id>` として local に作成または切り替える。
- 統合先 branch（既定 `master`）からの分岐元 commit を控える。
- remote repository を変更しない（push、tag push、remote branch delete は行わない）。

## plan.md

`plan.md` は合意済み要求だけを持つ。設計判断、未決定事項、参照、判断履歴、検証結果、修正結果は書かない。

- やることの要点: この修正で何をするかを、人間の依頼内容をそのまま要約して書く。後から `plan.md` 単体で何の task か分かる粒度にする。原因仮説は `investigation.md`、修正方針は `design.md` が持つ。

判断履歴は `plan.md` に残さず、恒久的に残す判断は `docs/changelog.md` に書く。

## 各成果物の詳細

### 人間観測記録

- Codex 本体が固定する。
- 人間が渡した不具合の観測記録から、確認済みの不具合、期待との差分、観測された操作または条件を task 内成果物として固定する。
- 観測事実だけを書く。仮説と推測は `調査` の原因仮説が持つ。

### 調査

- 再現確認と原因究明だけを扱う。どう直すかは `設計と仕様` で扱う。
- Codex 本体が `fix-decision` skill を Skill ツールで読んで適用する。
- 入力: 人間観測記録、作業計画フォルダ、Wails 接続対象、画面の正本（Storybook の story と svelte コンポーネント）、観測ログ仕様。
- テンプレートは `docs/exec-plans/templates/task-folder/investigation.md` を `investigation.md` として使う。
- 結果を `investigation.md` に固定する。次を含める:
    - 観測済み問題（人間観測記録から確認できる問題）。
    - 画面再現確認（`chrome-devtools` MCP ツールで実装済み画面の `data-testid` またはセレクタに従って再現した結果と再現手順）。
    - 原因仮説（複数仮説と検証順序）。
    - 観測ログ検証（一時ログ追加、観測結果、削除確認）。
    - 確定原因（観測で確定したもののみ）。

### 設計と仕様

- `design.md` は `design-protocol` skill に従って作る。
- `spec.md` は `specification-protocol` skill に従って作る。
- Codex 本体が両 skill を読み、`plan.md` の要求、`references.md`、対象 repository、`investigation.md` の確定原因、`docs/vocabulary.md`、各出力先を渡して適用する。
- `design.md` は `docs/exec-plans/templates/task-folder/design.md` を雛形にする。
- `spec.md` は `docs/exec-plans/templates/task-folder/spec.md` を雛形にする。
- `design.md`: as-is と to-be の方針。`investigation.md` の確定原因を根拠にする。
- `spec.md`: 修正後に成立させる仕様。
- `pending.md`: 未決定事項とブロッカー。解決した項目は結論を正本へ反映してから削除する。
- `references.md`: source と外部資料の所在。設計と仕様からは参照 ID だけを使う。
- 仕様は、`実装への引き継ぎ` の追加する fail-test の観点と同じ文にする。修正前に fail し、修正後に pass する対象がこの文で一意に決まるようにする。
- 本モジュールは書く順序と承認順序を固定する。
- 設計の書き方は `design-protocol` skill が持つ。
- 仕様の書き方は `specification-protocol` skill が持つ。

### design-review

- `design-review` は、人間修正レビューの前に、実現可能でない案と誤読の余地がある記述を否決する AI 検証である。`design_reviewer` agent（`fresh`、読み取り専用）が担う。
- 検証内容、判断範囲、出力は `design-review` skill に従う。
- 起動入力: `plan.md` の要求、`design.md`、`spec.md`、`pending.md`、`references.md`、対象 repository、`docs/vocabulary.md`、`investigation.md` の確定原因。`log.jsonl` は渡さない。
- 戻し先: 本モジュール。
- 完了: `design.md` の全ての as-is と to-be、および `spec.md` の全ての仕様に判定が付いている。
- 回数と3回目の扱いは `design-review` に従う。各回の判定と更新した正本を指す event を `log.jsonl` へ追記する。

### 人間修正レビュー

- `調査`、`設計と仕様`、`design-review` が固まった時点で人間へ返す。
- レビュー材料は `presentation` skill に従い、人間がわかりやすい構成で書く。人間が仕様を短時間で確認できるよう `spec.md` の要求ごとの節を先に示す。
- 差し戻しまたは追加質問の場合は、Codex 本体が同じ文脈で書き直す。

### 実装への引き継ぎ

- Codex 本体が固定する。`implementation-workflow` へ渡すものを返す成果物としてまとめる。
    - 確定原因（`investigation.md` から）。
    - 承認済み修正方針（`design.md` の直し方）。
    - 影響ファイル候補（観測事実に基づく候補）。
    - 確定仕様（`spec.md` の要求ごとの仕様）。
    - 再現手順と修正後に満たすべき期待状態（`investigation.md` の画面再現確認から）。
    - 追加する fail-test の観点（`spec.md` の仕様と同じ文。修正前に追加して fail を確認、修正後に pass を確認）。

## 不変条件

### 調査（investigation.md）の必須観点

| 観点 | 拘束する判断 |
| --- | --- |
| `観測済み問題` | 人間観測記録から確認できる問題だけを固定する。 |
| `画面再現確認` | 実装済み画面の `data-testid` またはセレクタに従い、人間観測記録のユーザー操作を `chrome-devtools` MCP ツールで再現する。 |
| `原因仮説` | 観測事実から複数の原因候補を立て、検証する順序と根拠を固定する。 |
| `観測ログ検証` | 仮説を否定または支持するために追加した一時ログ、観測結果、削除確認を固定する。 |
| `確定原因` | 観測で確定した原因だけを固定する。 |

### 設計と仕様の必須観点

- `design.md` の必須観点は `design-protocol` skill が持つ。
- `spec.md` の必須観点は `specification-protocol` skill が持つ。
- 修正フロー固有の拘束は次の 2 つとする。
    - 修正方針: `investigation.md` の確定原因に対応する直し方だけを固定する。仕様が不足していない場合だけ恒久修正を固定する。
    - fail-test との一致: 仕様の文と、追加する fail-test の観点を同じ文にする。

### レビュー順序ゲート

- `design-review` 通過なしで `人間修正レビュー` を依頼しない。
- `pending.md` が空でない場合は `design-review` と `人間修正レビュー` を依頼しない。

### 責務境界

- 調査（再現確認・原因究明）と設計（どう直すか）と仕様（何が成立すべきか）は別成果物に分ける。`design.md` に再現確認・原因究明を混ぜない。`spec.md` に直し方と設計理由を混ぜない。
- 調査、設計、仕様の本文は Codex 本体が task 文脈を持ったまま書く。`fresh` には渡さない。
- 影響ファイルは候補として扱う。実装の変更ファイルを本モジュールで確定しない。

### 安全境界

- 一時観測ログ以外のプロダクトコード変更は本モジュールでは行わない。
- プロダクトテスト変更、docs 正本本文変更は本モジュールでは行わない。
- 追加した一時観測ログは `調査` の `investigation.md` を固定する前に削除する。

### 仕様不足の停止

- 修正方針が仕様変更・機能追加・受け入れ条件の新規判断を必要とすると Codex 本体が判断した場合は、本モジュールを停止する。停止時は呼び出し元へ「`design-workflow` 経由の仕様変更が必要」と戻す。
- 修正方針が画面の見た目を変える場合は、本モジュールを停止する。停止時は呼び出し元へ、`design-workflow` で Storybook 人間レビューを完了し、承認された story、fixture、表示コンポーネントを `implementation-workflow` へ渡す必要があると返す。

## 返す成果物

- 作業 branch 名、統合先 branch、分岐元 commit。
- 作業計画フォルダのパス、`plan.md` のパス、`investigation.md` のパス、`design.md` のパス、`spec.md` のパス、`pending.md` のパス、`references.md` のパス。
- 確定仕様の要約: `spec.md` の要求ごとの仕様。
- 観測済み問題: 根拠から確認できる問題。
- 画面再現確認: 再現手順、操作結果、画面状態、証跡参照。
- 確定原因: 観測で確定した原因。
- 採用する修正方針: 恒久修正として採用する直し方。
- 影響ファイル候補: 観測事実に基づく候補。
- `design-review` の判定結果、否決理由、漏れ候補。
- 実装への引き継ぎ: `implementation-workflow` へ渡すもの（確定原因、承認済み修正方針、確定仕様、影響ファイル候補、再現手順と期待状態、追加する fail-test の観点）。
- 停止判定: 停止理由、不足項目、戻し先。

## 作業を止める条件

- 依頼要約または `task-id` が渡されていない。
- 統合先 branch が存在しない、または local に取得できない。
- 人間観測記録が不足する。
- Wails 接続対象が単一化されていない。
- 実装済み画面の `data-testid` またはセレクタに従ったユーザー操作を `chrome-devtools` MCP ツールで再現確認できない。
- 観測ログを追加または確認できず、仮説を検証できない。
- 追加した一時観測ログを削除できない。
- 原因が仮説に留まり、採用する修正方針を固定できない。
- 仕様が `plan.md` のどの要求の節にも置けない。
- `pending.md` が空でない。
- `design-review` の否決理由を `design.md` または `spec.md` の書き直しで解消できない。
- `design-review` が定める回数の上限後も `design_reviewer` が `目的未固定` または否決を返す。
- 修正方針が画面の見た目を変える。
- 修正方針が仕様変更・機能追加・受け入れ条件の新規判断を必要とすると判断した。
- `人間修正レビュー` 承認が得られない、または差し戻しを解消できない。
- 停止時は不足項目、固定できない判断、戻し先を返す。
