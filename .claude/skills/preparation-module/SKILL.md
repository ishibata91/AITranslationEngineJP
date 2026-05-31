---
name: preparation-module
description: "プロダクト変更を伴う task の入口で、作業 branch と active plan folder の plan.md を固定するモジュール skill。想定 Y/N、設計、人間確認は本モジュールで扱わず、後続モジュールへ渡す。TRIGGER when: プロダクト変更を伴う task の入口で、作業 branch と plan.md を固定する必要がある。"
---
# Preparation Module

## 目的

`preparation-module` は、プロダクト変更を伴う task の入口で、作業 branch と active plan folder の `plan.md` を固定するモジュール skill である。
想定 Y/N の評価、設計成果物の決定、人間確認は本モジュールでは扱わず、`design-module`（実装系 task）または `investigation-module`（修正系 task）の入口に委ねる。

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
- `plan.md` に依頼要約と分岐元 branch が記録されている。

## 担当 artifact

| 成果物ID | 担当 | 依存対象 | 起動先 |
| --- | --- | --- | --- |
| `branch 準備` | 呼び出し元 agent | `[]` | なし |
| `plan.md 初期化` | 呼び出し元 agent | `branch 準備` | なし |

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
- 想定 Y/N、artifact 集合、設計成果物の決定、人間確認は本モジュールで行わない。後続モジュールがそれぞれの入口で扱う。

## 返す成果物

- 作業 branch 名、分岐元 branch、分岐元 commit。
- active plan folder のパス、`plan.md` のパス。
- 後続モジュールへの引き継ぎ: 依頼要約、`task-id`。

## 作業を止める条件

- 依頼要約または `task-id` が渡されていない。
- 既存の作業 branch と `task-id` が衝突する。
- 統合先 branch が存在しない、または local に取得できない。
- 停止時は不足項目、衝突箇所、戻し先を返す。
