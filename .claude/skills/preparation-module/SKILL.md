---
name: preparation-module
description: "プロダクト変更を伴う task の入口で、作業 branch と active plan folder の plan.md を固定するモジュール skill。完了定義（システム上どこまで動かすか）は本モジュールで固定し、想定 Y/N、設計成果物、file 単位の触り方、人間確認は後続モジュールへ渡す。TRIGGER when: プロダクト変更を伴う task の入口で、作業 branch と plan.md を固定する必要がある。"
---
# Preparation Module

## 目的

`preparation-module` は、プロダクト変更を伴う task の入口で、作業 branch と active plan folder の `plan.md` を固定するモジュール skill である。
本モジュールは、システム上どこまで動かすか（task 後に観測できる振る舞いと、その観測点）を `完了定義` として固定する。これは「どこまでで完了か」を実装前に 1 つに決め、goal と除外範囲の矛盾を実装前に潰すためである。
想定 Y/N の評価、設計成果物の決定、file 単位の触り方、人間確認は本モジュールでは扱わず、`design-module`（実装系 task）または `investigation-module`（修正系 task）の入口に委ねる。

## 呼び出し関係

- 呼び出し元: 人間または agent。
- 返却先: 呼び出し元。
- モジュールが呼ぶ最下層 skill: なし。
- モジュールが呼ぶ最下層 agent: なし。

## 入口条件

- 呼び出し元から依頼要約と `task-id` が渡されている。
- 統合先 branch（既定 `master`）が local に存在する。

## 出口条件

- 作業 branch が `claude/<task-id>` として local に存在する。
- active plan folder `docs/exec-plans/active/<task-id>/` が存在する。
- `plan.md` に依頼要約と分岐元 branch、`完了定義`（動かす範囲と観測点）、軽 / 重判定結果が記録されている。

## 担当 artifact

| 成果物ID | 担当 | 依存対象 | 起動先 |
| --- | --- | --- | --- |
| `branch 準備` | 呼び出し元 agent | `[]` | なし |
| `plan.md 初期化` | 呼び出し元 agent | `branch 準備` | なし |
| `完了定義` | 呼び出し元 agent | `plan.md 初期化` | なし |
| `軽 / 重判定` | 呼び出し元 agent | `完了定義` | なし |

## 完了定義

`plan.md` 初期化の直後、システム上どこまで動かすかを 1 つ固定する。固定結果を `plan.md` に記録する。

- 動かす範囲: task 後に利用者または検証者が観測できる振る舞いを書く。差込点や接続境界を置くだけで振る舞いが観測できない状態を「動く」と書かない。
- 観測点: その振る舞いをどこで確かめるかを、単体テスト・実画面・実データのいずれかで書く。
- goal 整合: `完了定義` は goal と矛盾させない。goal がある振る舞いを「効く」と言うなら、`完了定義` はその振る舞いが実際に動くことを要求する。最小実装・仮実装・空テーブルだけで goal を満たしたことにしない。
- close_conditions 整合: `plan.md` の close_conditions は、`完了定義` の振る舞いを観測点で検証できる形にする。観測点を書かない close_conditions を残さない。
- goal と除外範囲（含まない）が矛盾する場合（goal が要る手段を 含まない が除外している）は、本モジュールで止め、goal を狭めるか 含まない を緩めるかを人間へ返す。
- 非対象: file 単位の触り方、設計差分図、テスト設計は本モジュールで決めない。重 task は `design-module`、軽 task は実装時に Claude 本体が決める。

## 軽 / 重判定

`完了定義` を書いた直後、次の 2 軸で task の重さを判定する。判定結果と根拠を `plan.md` に 1 行ずつ記録する。

- 画面が動くか（layout、文言、style、表示構造、svelte 表示コンポーネント、props、story、fixture のいずれかを変える）
- `docs/architecture.md` への反映が要るか（層構成、依存方向、Bootstrap、Wails 境界、強い制約のいずれかを変える）

判定結果による後続モジュール:

- 両方 N → 軽 task。`design-module` と `storybook-module` を bypass。`preparation-module` → `implementation-module` → `finalization-module`
- 片方でも Y → 重 task。`design-module` 経由。画面が動く場合は `storybook-module` も経由。`preparation-module` → `design-module` →（画面が動くなら `storybook-module`）→ `implementation-module` → `finalization-module`

修正系 task（バグ修正、refactor）は通常 軽 task。`investigation-module` 経由の場合は、`investigation-module` の入口で重さを再評価する。

## branch 準備

- 作業 branch を既定名 `claude/<task-id>` として local に作成または切り替える。
- 統合先 branch（既定 `master`）からの分岐元 commit を控える。
- remote repository を変更しない（push、tag push、remote branch delete は行わない）。

## plan.md 初期化

- `docs/exec-plans/active/<task-id>/plan.md` を作成または既存を確認する。
- 依頼要約と分岐元 branch、分岐元 commit を記録する。
- 後続モジュールが書き加える想定 Y/N、artifact 集合、設計成果物の参照、検証結果は本モジュールでは書かない。

## 不変条件

- プロダクトコード、プロダクトテスト、docs 正本本文を本モジュールでは変更しない。
- remote repository を変更しない。
- `完了定義`（動かす範囲と観測点）は本モジュールで固定する。
- 想定 Y/N、artifact 集合、設計成果物（実装範囲、設計差分図、テスト設計）の決定、file 単位の触り方、人間確認は本モジュールで行わない。後続モジュールがそれぞれの入口で扱う。

## 返す成果物

- 作業 branch 名、分岐元 branch、分岐元 commit。
- active plan folder のパス、`plan.md` のパス。
- `完了定義`: 動かす範囲、観測点。
- 後続モジュールへの引き継ぎ: 依頼要約、`task-id`、`完了定義`。

## 作業を止める条件

- 依頼要約または `task-id` が渡されていない。
- 既存の作業 branch と `task-id` が衝突する。
- 統合先 branch が存在しない、または local に取得できない。
- goal と除外範囲（含まない）が矛盾し、`完了定義` を矛盾なく固定できない。
- 停止時は不足項目、衝突箇所、戻し先を返す。
