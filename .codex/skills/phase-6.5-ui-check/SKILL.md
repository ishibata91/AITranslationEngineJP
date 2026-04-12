---
name: phase-6.5-ui-check
description: 実装完了後に Playwright MCP で主要導線と画面状態を確認し、UI 逸脱の証跡を返す。
---

# Phase 6.5 UI Check

## Goal

- `npm run dev:wails:docker-mcp` 起動後の `http://host.docker.internal:34115` を Playwright MCP から操作し、主要導線と画面状態を確認する
- implementation lane では承認済み HTML モック artifact、承認済み Scenario テスト一覧 artifact、承認済み task_id、review 用差分図、受け入れ確認と実装結果を照合する
- fix lane では active fix plan の再現条件、期待結果、accepted fix scope、実装結果を照合する
- UI 逸脱、console error、network failure、設計差分を切り分けて返す
- implementation lane では承認済み HTML モック artifact と実装画面の視覚構造を照合し、layout、主要情報ブロック、表示状態切替、主要導線の配置が一致しているかを確認する

## Rules

- 第6段階の完了後に進める
- ファイル送信が必要な場合、権限で引っかかるので `docker --context desktop-linux cp <host-path> <container-id>:/home/node/<file-name>` で Playwright MCP コンテナの `/home/node` へ先にコピーしてから `browser_file_upload` を使う
- `docker ps` でコンテナが見えない時は `docker --context desktop-linux ps` を優先し、同じ context で `docker cp` を実行する
- UI確認前に `npm run dev:wails:docker-mcp` が起動済みで、`http://host.docker.internal:34115` を開ける状態を確認する
- Playwright MCP を使った確認と証跡整理に限定する
- 新しい仕様解釈や見た目の好みを追加しない
- 恒久修正や test 追加は行わない
- UI 逸脱は第6段階へ戻し、設計差分だけを上流へ戻す
- implementation lane では `implementation_required_reading` を読まずに設計差分判定へ進まない
- fix lane では `reproduction_contract` と fix plan を読まずに pass / reroute を決めない
- 再現操作は主要導線と高リスク状態に絞る
- 視覚構造の確認では、implementation lane では承認済み HTML モック artifact を source of truth とし、見た目の好みではなく layout、情報ブロック、ラベル、状態表示、導線配置の一致だけを判定する
- 承認済み HTML モック artifact に存在する主要構造が実装で欠落、統合、簡略化されている時は `reroute` とする

## Reference Use

- implementation lane では着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.phase-6.5-ui-check.json` を参照して入力契約を確認する。
- fix lane では着手前に `../orchestrating-fixes/references/orchestrating-fixes.to.phase-6.5-ui-check.json` を参照して入力契約を確認する。
- `orchestrating-implementation` へ返す時は `references/phase-6.5-ui-check.to.orchestrating-implementation.json` を返却契約として使う。
- `orchestrating-fixes` へ返す時は `references/phase-6.5-ui-check.to.orchestrating-fixes.json` を返却契約として使う。
