---
name: feature-workflow
description: "新規実装フローの入口オーケストレーター。プロダクト変更を伴う実装系 task の入口で、作業 branch と plan.md を固定し、実装方針（現況の理解とあるべき形と変更点をソース根拠つきで示し、流れ・関係・責務が変わる箇所は図を添える）を design.md に、確定仕様（要求ごとの仕様）を spec.md にまとめ、design-review（design_reviewer agent によるソース照合検証と記述検証）と人間設計レビューを通してから下流の段（storybook-module / implementation-module / finalization-module）へ渡す。実装範囲・テスト設計は扱わない。TRIGGER when: プロダクト変更を伴う実装系 task の入口。SKIP when: 修正系 task（バグ修正・refactor）は fix-workflow へ。"
---
# Feature Workflow

## 目的

`feature-workflow` は、新規実装フローの入口オーケストレーターである。プロダクト変更を伴う実装系 task の入口で、作業 branch と `plan.md` を固定し、実装方針を `design.md` に、確定仕様を `spec.md` にまとめ、`design-review` と人間設計レビューを通してから下流の段へ渡す。
`design.md` は人間が読んで判断する説明を持ち、`spec.md` は下流が実装根拠にする仕様を持つ。`design.md` と `spec.md` は `plan.md` の要求ごとの節に分ける。両者が食い違う場合は `spec.md` を優先する。

実装系 task は重さによらず本モジュールを通す。
修正系 task（バグ修正、refactor）は `fix-workflow` が扱う。

## 呼び出し関係

- 呼び出し元: 人間、または上位 agent。
- 返却先: 呼び出し元。
- モジュールが呼ぶ下位 skill: `design-protocol`（Codex 本体が読んで適用。`design.md` と `spec.md` の書き方、`design-review` の検証規約）。`presentation`（参照のみ。Codex 本体が読んで適用し、現況とあるべき形の 2 図と設計レビュー材料を作る）。
- モジュールが呼ぶ下位 agent: `design_reviewer`（`fresh`）。人間設計レビューの前に `design.md` の現況の理解と変更点を実ソースと照合して検証し、`spec.md` の各仕様が要求の節に置かれ、確かめ方を持つかを検証する。設計判断と `design.md`、`spec.md` は Codex 本体が task 文脈を持ったまま書く。
- 下流の段: `storybook-module`（画面表示の変更がある場合）、`implementation-module`、`finalization-module`。

## 入口条件

- 呼び出し元から依頼要約と `task-id` が渡されている。
- 統合先 branch（既定 `master`）が local に存在する。

## 出口条件

- 作業 branch が `codex/<task-id>` として local に存在する。
- 作業計画フォルダ `docs/exec-plans/active/<task-id>/` が存在する。
- `plan.md` に branch 情報と対象 task でやること・扱わないことの要点が記録されている。
- `design.md` の各要求の節に現況の理解、あるべき形、変更点が記録され、末尾に検討が必要なことが記録されている。
- `spec.md` の各要求の節に仕様が記録され、各仕様に前提条件と確かめ方が書かれている。
- `design.md` の検討が必要なことに未解決の論点が残っていない（人間回答で解消済み）。
- `design-review` が通過済み。
- `人間設計レビュー` が承認済み。

## 担当成果物

| 成果物ID | 担当 | 依存対象 | 起動先 |
| --- | --- | --- | --- |
| `branch 準備` | Codex 本体 | 入口条件 | なし |
| `plan.md` | Codex 本体 | `branch 準備` | なし |
| `design.md` | Codex 本体 | `plan.md` | なし |
| `spec.md` | Codex 本体 | `design.md` | なし |
| `design-review` | `design_reviewer`（`fresh`） | `design.md`, `spec.md` | `design_reviewer` agent |
| `人間設計レビュー` | 人間 | `design-review` | 人間 |

## branch 準備

- 作業 branch を既定名 `codex/<task-id>` として local に作成または切り替える。
- 統合先 branch（既定 `master`）からの分岐元 commit を控える。
- remote repository を変更しない（push、tag push、remote branch delete は行わない）。

## plan.md

`plan.md` は次の 3 つだけを持つ。設計判断、判断履歴、検証結果、実装結果は書かない。

- branch 情報: 作業 branch 名、統合先 branch 名、分岐元 commit。
- やることの要点: 対象 task で何をするかを、人間の依頼内容をそのまま要約して書く。後から `plan.md` 単体で何の task か分かる粒度にする。手段選定や設計判断は `design.md` が持つ。
- やらないことの要点: 対象 task で扱わない範囲を大まかに書く。goal に要る手段が対象として残っているかを確かめる。

判断履歴は `plan.md` に残さず、恒久的に残す判断は `docs/changelog.md` に書く。

## design.md と spec.md

書き方は `design-protocol` skill に従う。Codex 本体が Skill ツールで読み、フロー種別として新規実装フローを渡して適用する。

- `design.md`: 人間が読んで判断する説明（実装方針、検討が必要なこと）。
- `spec.md`: 下流の段が実装根拠にする確定仕様（要求ごとの仕様）。
- 本モジュールが固定するのは書く順序と承認順序とし、両 file の書き方は `design-protocol` skill が持つ。

## design-review

`design-review` は、人間設計レビューの前に、実現可能でない案と誤読の余地がある記述を否決して人間との無駄な往復を減らすための AI 検証である。`design_reviewer` agent（`fresh`、読み取り専用）が担う。

- 検証内容、判断範囲、出力、否決時の扱いは `design-protocol` skill の `design-review` 節に従う。
- 起動入力: `design.md` のパス、`spec.md` のパス、`task-id`、対象 repo のルート、フロー種別（新規実装フロー）。
- 戻し先: 本モジュール。
- 完了: `design.md` の全ての現況の理解と全ての変更点、および `spec.md` の全ての仕様に判定が付いている。

## 人間設計レビュー

- `design-review` 通過後に人間へレビューを依頼する。人間が仕様を短時間で確認できるよう `spec.md` の要求ごとの節を先に示し、`design.md` は理由と変え方の説明として添える。`design-review` の判定結果（漏れ候補を含む）もレビュー材料に添える。画面表示の視覚レビューは `storybook-module` の Storybook 人間レビューループで行う。
- レビュー材料は `presentation` skill に従い、人間がわかりやすい構成で書く。材料は対象 task に当たる論点だけ残す。
- 差し戻しまたは追加質問の場合は、Codex 本体が同じ文脈で `design.md` と `spec.md` を書き直す。
- `検討が必要なこと` に未解決の論点が残る場合は、人間の回答を得るまで承認へ進めない。

## 下流への引き継ぎ

- 画面表示の変更（layout、文言、style、表示構造、svelte 表示コンポーネント、props、story、fixture のいずれか）がある場合は `storybook-module` を経由する。
- 画面表示の変更がない場合は `storybook-module` を bypass して `implementation-module` へ進む。
- 実装系 task の経路: `feature-workflow` →（画面表示の変更があれば `storybook-module`）→ `implementation-module` → `finalization-module`。

## 不変条件

- `design-review` 通過なしで `人間設計レビュー` を依頼しない。
- `人間設計レビュー` 承認なしで下流の段へ進めない。
- `検討が必要なこと` に未解決の論点が残るまま下流の段へ進めない。
- `spec.md` の要求ごとの仕様を固定せずに下流の段へ進めない。
- 差し戻し時は Codex 本体が同じ文脈で書き直す。`fresh` に分割しない。
- プロダクトコード、プロダクトテスト、docs 正本本文を本モジュールでは変更しない。
- remote repository を変更しない。

## 返す成果物

- 作業 branch 名、統合先 branch、分岐元 commit。
- 作業計画フォルダのパス、`plan.md` のパス、`design.md` のパス、`spec.md` のパス。
- 確定仕様の要約: `spec.md` の要求ごとの仕様。
- `design-review` の判定結果、否決理由、漏れ候補。
- 人間設計レビューの承認状態、差し戻し記録。
- 下流への引き継ぎ: `task-id`、`design.md` の実装方針、`spec.md` の要求ごとの仕様、画面表示の変更の有無、`storybook-module` へ進むかどうか。

## 作業を止める条件

- 依頼要約または `task-id` が渡されていない。
- 既存の作業 branch と `task-id` が衝突する。
- 統合先 branch が存在しない、または local に取得できない。
- goal とやらないことが矛盾し、実装方針を矛盾なく固定できない。
- 仕様が `plan.md` のどの要求の節にも置けない。
- `検討が必要なこと` の未解決の論点に人間回答が得られない。
- `design-review` の否決理由を `design.md` または `spec.md` の書き直しで解消できない。
- `人間設計レビュー` 承認が得られない、または差し戻しを解消できない。
- 設計判断が AI 単独で確定できる範囲を越え、人間判断が要るのに得られない。
- 停止時は不足項目、衝突箇所、固定できない判断、戻し先を返す。
