---
name: phase-6.5-ui-check
description: 実装完了後に `chrome-devtools` で主要導線と画面状態を確認し、UI 逸脱の証跡を返す。
---

# Phase 6.5 UI Check

## Goal

- `wails dev` 起動後の `http://localhost:34115` を MCP 経由の `chrome-devtools` から操作し、主要導線と画面状態を確認する
- 承認済み HTML モック artifact、承認済み Scenario テスト一覧 artifact、承認済み task_id、review 用差分図、受け入れ確認と実装結果を照合する
- UI 逸脱、console error、network failure、設計差分を切り分けて返す
- 承認済み HTML モック artifact と実装画面の視覚構造を照合し、layout、主要情報ブロック、表示状態切替、主要導線の配置が一致しているかを確認する

## Rules

- 第6段階の完了後に進める
- UI確認前に `wails dev` が起動済みで、`http://localhost:34115` を開ける状態を確認する
- MCP 経由の `chrome-devtools` を使った確認と証跡整理に限定する
- 新しい仕様解釈や見た目の好みを追加しない
- 恒久修正や test 追加は行わない
- UI 逸脱は第6段階へ戻し、設計差分だけを上流へ戻す
- `implementation_required_reading` を読まずに設計差分判定へ進まない
- 再現操作は主要導線と高リスク状態に絞る
- 視覚構造の確認では、承認済み HTML モック artifact を source of truth とし、見た目の好みではなく layout、情報ブロック、ラベル、状態表示、導線配置の一致だけを判定する
- 承認済み HTML モック artifact に存在する主要構造が実装で欠落、統合、簡略化されている時は `reroute` とする

## Reference Use

- 着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.phase-6.5-ui-check.json` を参照して入力契約を確認する。
- `orchestrating-implementation` へ返す時は `references/phase-6.5-ui-check.to.orchestrating-implementation.json` を返却契約として使う。
