---
name: feature-workflow
description: "新規実装フローの入口オーケストレーター。プロダクト変更を伴う実装系 task の入口で、作業 branch と plan.md を固定し、実装方針（現状 AS-IS と変更後 TO-BE を文章で示し、流れ・関係・責務が変わる箇所は図を添える）を design.md にまとめ、人間設計レビューを通してから下流の段（storybook-module / implementation-module / finalization-module）へ渡す。実装範囲・テスト設計は扱わない。TRIGGER when: プロダクト変更を伴う実装系 task の入口。SKIP when: 修正系 task（バグ修正・refactor）は fix-workflow へ。"
---
# Feature Workflow

## 目的

`feature-workflow` は、新規実装フローの入口オーケストレーターである。プロダクト変更を伴う実装系 task の入口で、作業 branch と `plan.md` を固定し、実装方針を `design.md` にまとめ、人間設計レビューを通してから下流の段へ渡す。

本モジュールは実装系 task の重さで分岐しない。軽い実装系 task も本モジュールを通す。
修正系 task（バグ修正、refactor）は本モジュールで扱わず `fix-workflow` へ渡す。

## 呼び出し関係

- 呼び出し元: 人間、または上位 agent。
- 返却先: 呼び出し元。
- モジュールが呼ぶ下位 skill: `presentation`（参照のみ。Claude 本体が読んで適用し、AS-IS→TO-BE の変更図と設計レビュー材料を作る）。
- モジュールが呼ぶ下位 agent: なし。設計判断と `design.md` は Claude 本体が task 文脈を持ったまま書く。
- 下流の段: `storybook-module`（画面表示の変更がある場合）、`implementation-module`、`finalization-module`。

## 入口条件

- 呼び出し元から依頼要約と `task-id` が渡されている。
- 統合先 branch（既定 `master`）が local に存在する。

## 出口条件

- 作業 branch が `claude/<task-id>` として local に存在する。
- 作業計画フォルダ `docs/exec-plans/active/<task-id>/` が存在する。
- `plan.md` に branch 情報と対象 task で扱わないことの要点が記録されている。
- `design.md` に実装方針（AS-IS→TO-BE を文章で、必要なら図を添える）と検討が必要なことが記録されている。
- `design.md` の検討が必要なことに未解決の論点が残っていない（人間回答で解消済み）。
- `人間設計レビュー` が承認済み。

## 担当成果物

| 成果物ID | 担当 | 依存対象 | 起動先 |
| --- | --- | --- | --- |
| `branch 準備` | Claude 本体 | 入口条件 | なし |
| `plan.md` | Claude 本体 | `branch 準備` | なし |
| `design.md` | Claude 本体 | `plan.md` | なし |
| `人間設計レビュー` | 人間 | `design.md` | 人間 |

## branch 準備

- 作業 branch を既定名 `claude/<task-id>` として local に作成または切り替える。
- 統合先 branch（既定 `master`）からの分岐元 commit を控える。
- remote repository を変更しない（push、tag push、remote branch delete は行わない）。

## plan.md

`plan.md` は次の 2 つだけを持つ。設計、判断履歴、検証結果、実装結果は書かない。

- branch 情報: 作業 branch 名、統合先 branch 名、分岐元 commit。
- やらないことの要点: 対象 task で扱わない範囲を大まかに書く。goal が要る手段を除外していないかを確かめる。

判断履歴は `plan.md` に残さず、恒久的に残す判断は `docs/changelog.md` に書く。

## design.md

`design.md` は次の 2 つを持つ。実装範囲の scope 列挙とテスト設計は本モジュールで扱わない。

- 実装方針: どう実装し、どう変えるかを文章で書く。現状（AS-IS）と変更後（TO-BE）を対にして、何がどう変わるかを文章で示す。どこまで動かすか（task 後に観測できる振る舞い）と、観測点（単体テスト・実画面・実データのいずれか）を含める。最小実装・仮実装・空テーブルだけで goal を満たしたことにしない。流れ・関係・責務が変わる箇所は、AS-IS と TO-BE の 2 図を文章に添える（図作法は `presentation` skill に従う）。図だけで説明せず、文章を必ず書く。
- 検討が必要なこと: 人間の回答なしに先へ進めない未解決の論点。回答が得られるまで下流の段へ進まない。

## 人間設計レビュー

- `design.md` が揃った時点で人間へレビューを依頼する。画面表示の視覚レビューは `storybook-module` の Storybook 人間レビューループで行う。
- レビュー材料は `presentation` skill に従い、人間がわかりやすい構成で書く。材料は対象 task に当たる論点だけ残す。
- 差し戻しまたは追加質問の場合は、Claude 本体が同じ文脈で `design.md` を書き直す。
- `検討が必要なこと` に未解決の論点が残る場合は、人間の回答を得るまで承認へ進めない。

## 下流への引き継ぎ

- 画面表示の変更（layout、文言、style、表示構造、svelte 表示コンポーネント、props、story、fixture のいずれか）がある場合は `storybook-module` を経由する。
- 画面表示の変更がない場合は `storybook-module` を bypass して `implementation-module` へ進む。
- 実装系 task の経路: `feature-workflow` →（画面表示の変更があれば `storybook-module`）→ `implementation-module` → `finalization-module`。

## 不変条件

- `人間設計レビュー` 承認なしで下流の段へ進めない。
- `検討が必要なこと` に未解決の論点が残るまま下流の段へ進めない。
- 差し戻し時は Claude 本体が同じ文脈で書き直す。`fresh` に分割しない。
- プロダクトコード、プロダクトテスト、docs 正本本文を本モジュールでは変更しない。
- remote repository を変更しない。

## 返す成果物

- 作業 branch 名、統合先 branch、分岐元 commit。
- 作業計画フォルダのパス、`plan.md` のパス、`design.md` のパス。
- 人間設計レビューの承認状態、差し戻し記録。
- 下流への引き継ぎ: `task-id`、`design.md` の実装方針、画面表示の変更の有無、`storybook-module` へ進むかどうか。

## 作業を止める条件

- 依頼要約または `task-id` が渡されていない。
- 既存の作業 branch と `task-id` が衝突する。
- 統合先 branch が存在しない、または local に取得できない。
- goal とやらないことが矛盾し、実装方針を矛盾なく固定できない。
- `検討が必要なこと` の未解決の論点に人間回答が得られない。
- `人間設計レビュー` 承認が得られない、または差し戻しを解消できない。
- 設計判断が AI 単独で確定できる範囲を越え、人間判断が要るのに得られない。
- 停止時は不足項目、衝突箇所、固定できない判断、戻し先を返す。
