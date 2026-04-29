---
name: codex-review-state-invariant
description: Codex 実装後 review の状態・データ不変条件グループ作業プロトコル。
---
# Codex Review State Invariant

## 目的

DB、キャッシュ、非同期処理、再実行、同時実行で壊れないかを見る。
本番で壊れやすい状態遷移を、diff から取得した実コードで score 化する。

## 対応ロール

- `review_state_invariant` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` の review aggregation とする。
- owner artifact は `codex-review-state-invariant` の出力規約で固定する。

## 入力規約

- transaction、lock、retry、idempotency
- race condition と DB 更新順序
- cache invalidation と event 発行
- queue consumer と partial failure
- soft delete、集計値、二重作成、二重課金
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。
- 必須入力: review_target_diff, implementation_scope_path, implementation_result
- 任意入力: final_validation_result, touched_files

## 外部参照規約

- agent runtime と tool policy は [review_state_invariant.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/review_state_invariant.toml) の `allowed_write_paths` / `allowed_commands` とする。
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-state-invariant/SKILL.md

## 内部参照規約

## 判断規約

`score > 0.85` を pass とする。
再実行不能、partial failure、二重処理の可能性は高い減点対象にする。

## 出力規約

- `observed_scope`: 確認した状態遷移、永続化、cache、queue と未確認範囲
- `violated_invariant`: 破られた状態、再実行、二重処理、整合性の invariant
- `root_cause_hypotheses`: 状態破壊を生む原因候補と根拠
- `local_patch_assessment`: 局所修正で invariant が戻るか、状態境界の再固定が必要か
- `exploration_scope`: 追加で読むべき状態遷移と永続化範囲
- `remediation_considerations`: 修正者が考慮すべき支配点、partial failure risk、invariant tests
- 出力に tool policy、agent runtime、product code の変更義務を含めない。
- 必須出力: group, decision, score, confidence, observed_scope, violated_invariant, root_cause_hypotheses, local_patch_assessment, exploration_scope, remediation_considerations, invariant_tests, findings, evidence_paths, not_reviewed
- 出力 field 要件: {"group": "state_invariant", "decision": "pass / fail / blocked のいずれか。score > 0.85 の場合だけ pass", "score": "DB、cache、非同期処理、再実行、同時実行に対する不変条件維持度を 0.0 から 1.0 で返す", "confidence": "根拠十分性を 0.0 から 1.0 で返す", "observed_scope": "確認した状態遷移、永続化、cache、queue、未確認範囲を分ける", "violated_invariant": "破られた状態、再実行、二重処理、整合性の invariant を書く。未特定なら不明理由を書く", "root_cause_hypotheses": "状態破壊を生む原因候補、根拠、反証余地を分ける", "local_patch_assessment": "局所修正で invariant が戻るか、状態境界の再固定が必要か、局所修正が危険な理由を書く", "exploration_scope": "追加で読むべき状態遷移、永続化、cache、queue 範囲と読まない範囲を書く", "remediation_considerations": "修正者が考慮すべき candidate control points、partial failure risk、durable fix signal を返す。修正範囲を命令しない", "invariant_tests": "破れた state invariant を固定する test 観点を返す", "findings": "transaction、lock、idempotency、retry、race condition、cache invalidation、event、queue、partial failure、soft delete、集計値、二重作成、二重課金だけを書く"}

## 完了規約

- 対象 review 観点の finding、score、根拠、残留 risk が返却されている。
- 権限・信頼境界系の hard gate は score 相殺せず明示されている。
- transaction、lock、idempotency、retry を確認した。
- race condition と DB 更新順序を確認した。
- cache invalidation と event 発行を確認した。
- partial failure、soft delete、集計値を確認した。
- 二重作成または二重課金の可能性を確認した。
- violated invariant と root cause hypothesis を分けた。
- local patch assessment と invariant tests を返した。
- completion signal: 状態・データ不変条件の score、破られた invariant、root cause hypothesis、local patch assessment、根拠が返っている
- residual risk key: not_reviewed

## 停止規約

- SQL の見た目
- 命名
- テストコードの構成
- UI 文言
- 内部設計の美しさ
- 停止時は不足項目、衝突箇所、reroute 先を返す。
