---
name: feature-workflow
description: "新規実装フローの入口オーケストレーター。プロダクト変更を伴う実装系 task の入口で、作業 branch と plan.md を固定し、実装方針（現状 AS-IS と変更後 TO-BE をソース根拠つきで示し、流れ・関係・責務が変わる箇所は図を添える）を design.md に、確定仕様（対象集合と期待値表）を spec.md にまとめ、design-review（design_reviewer agent によるソース照合検証と記述検証）と人間設計レビューを通してから下流の段（storybook-module / implementation-module / finalization-module）へ渡す。実装範囲・テスト設計は扱わない。TRIGGER when: プロダクト変更を伴う実装系 task の入口。SKIP when: 修正系 task（バグ修正・refactor）は fix-workflow へ。"
---
# Feature Workflow

## 目的

`feature-workflow` は、新規実装フローの入口オーケストレーターである。プロダクト変更を伴う実装系 task の入口で、作業 branch と `plan.md` を固定し、実装方針を `design.md` に、確定仕様を `spec.md` にまとめ、`design-review` と人間設計レビューを通してから下流の段へ渡す。
`design.md` は人間が読んで判断する説明を持ち、`spec.md` は下流が実装根拠にする対象集合と期待値を持つ。両者が食い違う場合は `spec.md` を優先する。

実装系 task は重さによらず本モジュールを通す。
修正系 task（バグ修正、refactor）は `fix-workflow` が扱う。

## 呼び出し関係

- 呼び出し元: 人間、または上位 agent。
- 返却先: 呼び出し元。
- モジュールが呼ぶ下位 skill: `presentation`（参照のみ。Claude 本体が読んで適用し、AS-IS→TO-BE の変更図と設計レビュー材料を作る）。
- モジュールが呼ぶ下位 agent: `design_reviewer`（`fresh`）。人間設計レビューの前に `design.md` の AS-IS 根拠と TO-BE 実現主張を実ソースと照合して検証し、`spec.md` の記述に誤読の余地がないかを検証する。設計判断と `design.md`、`spec.md` は Claude 本体が task 文脈を持ったまま書く。
- 下流の段: `storybook-module`（画面表示の変更がある場合）、`implementation-module`、`finalization-module`。

## 入口条件

- 呼び出し元から依頼要約と `task-id` が渡されている。
- 統合先 branch（既定 `master`）が local に存在する。

## 出口条件

- 作業 branch が `claude/<task-id>` として local に存在する。
- 作業計画フォルダ `docs/exec-plans/active/<task-id>/` が存在する。
- `plan.md` に branch 情報と対象 task でやること・扱わないことの要点が記録されている。
- `design.md` に実装方針（AS-IS→TO-BE を文章で、必要なら図を添える）と検討が必要なことが記録されている。
- `spec.md` に対象集合と期待値表（AS-IS の期待値と TO-BE の期待値）が記録されている。
- `design.md` の検討が必要なことに未解決の論点が残っていない（人間回答で解消済み）。
- `design-review` が通過済み。
- `人間設計レビュー` が承認済み。

## 担当成果物

| 成果物ID | 担当 | 依存対象 | 起動先 |
| --- | --- | --- | --- |
| `branch 準備` | Claude 本体 | 入口条件 | なし |
| `plan.md` | Claude 本体 | `branch 準備` | なし |
| `design.md` | Claude 本体 | `plan.md` | なし |
| `spec.md` | Claude 本体 | `design.md` | なし |
| `design-review` | `design_reviewer`（`fresh`） | `design.md`, `spec.md` | `design_reviewer` agent |
| `人間設計レビュー` | 人間 | `design-review` | 人間 |

## branch 準備

- 作業 branch を既定名 `claude/<task-id>` として local に作成または切り替える。
- 統合先 branch（既定 `master`）からの分岐元 commit を控える。
- remote repository を変更しない（push、tag push、remote branch delete は行わない）。

## plan.md

`plan.md` は次の 3 つだけを持つ。設計判断、判断履歴、検証結果、実装結果は書かない。

- branch 情報: 作業 branch 名、統合先 branch 名、分岐元 commit。
- やることの要点: 対象 task で何をするかを、人間の依頼内容をそのまま要約して書く。後から `plan.md` 単体で何の task か分かる粒度にする。手段選定や設計判断は `design.md` が持つ。
- やらないことの要点: 対象 task で扱わない範囲を大まかに書く。goal に要る手段が対象として残っているかを確かめる。

判断履歴は `plan.md` に残さず、恒久的に残す判断は `docs/changelog.md` に書く。

## design.md

テンプレートは `docs/exec-plans/templates/task-folder/design.md` を使う。
`design.md` は人間が読んで判断する説明を持ち、次の 2 つで構成する。実装範囲の scope 列挙とテスト設計は本モジュールで扱わない。対象集合と期待値は `spec.md` が持つ。

- 実装方針: どう実装し、どう変えるかを文章で書く。現状（AS-IS）と変更後（TO-BE）を対にして、何がどう変わるかを文章で示す。AS-IS と TO-BE の対応はテンプレートの表で書き、AS-IS には根拠となるソースの場所を、TO-BE には変更予定の場所と「どこをどう変えれば TO-BE が成立するか」の実現主張を添える。場所を挙げられない TO-BE は書かない。どこまで動かすか（task 後に観測できる振る舞い）と、観測点（単体テスト・実画面・実データのいずれか）を含める。最小実装・仮実装・空テーブルだけで goal を満たしたことにしない。流れ・関係・責務が変わる箇所は、AS-IS と TO-BE の 2 図を文章に添える（図作法は `presentation` skill に従う）。説明は文章で成立させ、図は理解を助ける補助として置く。
- 検討が必要なこと: 人間の回答なしに先へ進めない未解決の論点。回答が得られるまで下流の段へ進まない。

## spec.md

テンプレートは `docs/exec-plans/templates/task-folder/spec.md` を使う。
`spec.md` は下流の段が実装根拠にする確定仕様を持つ。設計理由、変更手順、図は持たない（`design.md` が持つ）。`design.md` と食い違う場合は `spec.md` を優先する。

- 対象集合: この task の処理対象に含むものと含まないものを、それぞれ肯定形で列挙する。読み手が範囲を 1 通りに取れる粒度で書く。
- 期待値表: 1 行 1 期待値で、AS-IS の期待値（現状で成立している振る舞い。現状が無い行は「なし」）と TO-BE の期待値（変更後に成立させる振る舞い）を横に並べる。TO-BE の期待値の文は、そのまま実テストの test case 名として使う。
- 記述の拘束: 主語と目的語を書く。TO-BE の期待値の語尾は「〜こと」で終える。1 通りにしか読めない語を選ぶ。対象に含む側の境界と含まない側の境界を、それぞれ独立した 1 行として立てる。
- 「対応する実テスト」列: 本モジュールでは空にする。`implementation-module` が最終検証で埋める。

## design-review

`design-review` は、人間設計レビューの前に、実現可能でない案と誤読の余地がある記述を否決して人間との無駄な往復を減らすための AI 検証である。`design_reviewer` agent（`fresh`、読み取り専用）が担う。

- 入力: `design.md` のパス、`spec.md` のパス、`task-id`、対象 repo のルート。
- 検証内容:
  - AS-IS 検証: `design.md` の AS-IS 記述を根拠ソースの場所と照合し、現状認識が実ソースで成立しているかを確かめる。
  - TO-BE 検証: 変更予定箇所の実現主張（どこをどう変えれば TO-BE が成立するか）を実ソースへ踏み込んで確かめ、成立しない主張を否決する。
  - 漏れ検出: 変更予定箇所の列挙から漏れた影響先（呼び出し元、依存先、同じ前提を持つ別箇所）がないかを確かめる。
  - 記述検証: `spec.md` の対象集合が肯定形で列挙されているか、主語と目的語が書かれているか、TO-BE の期待値の語尾が「〜こと」で終わっているか、各行が 1 通りにしか読めないかを判定する。読み手が範囲を広くも狭くも取れる行は、該当行を挙げて否決する。
  - 整合検証: `design.md` の TO-BE と `spec.md` の期待値が同じ振る舞いを指しているかを確かめる。`design.md` にあって期待値表に無い振る舞いは漏れ候補として返す。
- 判断範囲: 成立可否の判定、記述違反の指摘、漏れの指摘だけを行う。設計の好み、代替案の選定、実装、`design.md` と `spec.md` の書き換えは行わない。
- 出力: 判定（通過または否決）、否決理由（照合したソースの場所または該当行つき）、漏れ候補の一覧。戻し先は本モジュール。
- 完了: `design.md` の全 AS-IS 根拠と全 TO-BE 実現主張、および `spec.md` の対象集合と全期待値行に判定が付いている。
- 停止: `design.md` に根拠ソースの場所または変更予定箇所が書かれていない、`spec.md` に対象集合または期待値表がない、または対象ソースを読めない。停止時は不足箇所を本モジュールへ返す。
- 否決時: Claude 本体が同じ文脈で `design.md` または `spec.md` を書き直し、`design-review` を再実行する。否決理由を解消できない場合は人間へ論点として上げる。

## 人間設計レビュー

- `design-review` 通過後に人間へレビューを依頼する。人間が仕様を短時間で確認できるよう `spec.md` の対象集合と期待値表を先に示し、`design.md` は理由と変え方の説明として添える。`design-review` の判定結果（漏れ候補を含む）もレビュー材料に添える。画面表示の視覚レビューは `storybook-module` の Storybook 人間レビューループで行う。
- レビュー材料は `presentation` skill に従い、人間がわかりやすい構成で書く。材料は対象 task に当たる論点だけ残す。
- 差し戻しまたは追加質問の場合は、Claude 本体が同じ文脈で `design.md` と `spec.md` を書き直す。
- `検討が必要なこと` に未解決の論点が残る場合は、人間の回答を得るまで承認へ進めない。

## 下流への引き継ぎ

- 画面表示の変更（layout、文言、style、表示構造、svelte 表示コンポーネント、props、story、fixture のいずれか）がある場合は `storybook-module` を経由する。
- 画面表示の変更がない場合は `storybook-module` を bypass して `implementation-module` へ進む。
- 実装系 task の経路: `feature-workflow` →（画面表示の変更があれば `storybook-module`）→ `implementation-module` → `finalization-module`。

## 不変条件

- `design-review` 通過なしで `人間設計レビュー` を依頼しない。
- `人間設計レビュー` 承認なしで下流の段へ進めない。
- `検討が必要なこと` に未解決の論点が残るまま下流の段へ進めない。
- `spec.md` の対象集合と期待値表を固定せずに下流の段へ進めない。
- 差し戻し時は Claude 本体が同じ文脈で書き直す。`fresh` に分割しない。
- プロダクトコード、プロダクトテスト、docs 正本本文を本モジュールでは変更しない。
- remote repository を変更しない。

## 返す成果物

- 作業 branch 名、統合先 branch、分岐元 commit。
- 作業計画フォルダのパス、`plan.md` のパス、`design.md` のパス、`spec.md` のパス。
- 確定仕様の要約: `spec.md` の対象集合と TO-BE の期待値。
- `design-review` の判定結果、否決理由、漏れ候補。
- 人間設計レビューの承認状態、差し戻し記録。
- 下流への引き継ぎ: `task-id`、`design.md` の実装方針、`spec.md` の対象集合と期待値表、画面表示の変更の有無、`storybook-module` へ進むかどうか。

## 作業を止める条件

- 依頼要約または `task-id` が渡されていない。
- 既存の作業 branch と `task-id` が衝突する。
- 統合先 branch が存在しない、または local に取得できない。
- goal とやらないことが矛盾し、実装方針を矛盾なく固定できない。
- 対象集合または期待値を肯定形で固定できない（何を含めるかが人間回答なしに確定しない）。
- `検討が必要なこと` の未解決の論点に人間回答が得られない。
- `design-review` の否決理由を `design.md` または `spec.md` の書き直しで解消できない。
- `人間設計レビュー` 承認が得られない、または差し戻しを解消できない。
- 設計判断が AI 単独で確定できる範囲を越え、人間判断が要るのに得られない。
- 停止時は不足項目、衝突箇所、固定できない判断、戻し先を返す。
