---
name: codex-review-contract
description: Codex 実装後 review の契約・互換性グループ作業プロトコル。
---
# Codex Review Contract

## 目的

既存利用者、外部 API、内部 API、DB schema、event payload を壊していないかを見る。
コード自体が動いても契約破壊が利用者側障害になるため、diff の public boundary を score 化する。

## 対応ロール

- `review_contract` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` の review aggregation とする。
- owner artifact は `codex-review-contract` の出力規約で固定する。

## 入力規約

- API request / response
- GraphQL schema と DB migration
- public method と event payload
- queue message と webhook
- error code、nullable / required、versioning
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。
- 必須入力: review_target_diff, implementation_scope_path, implementation_result
- 任意入力: final_validation_result, touched_files

## 外部参照規約

- agent runtime と tool policy は [review_contract.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/review_contract.toml) の `allowed_write_paths` / `allowed_commands` とする。
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-contract/SKILL.md

## 内部参照規約

## 判断規約

`score > 0.85` を pass とする。
既存 field の意味変更や nullable / required の変更は高い減点対象にする。

## 出力規約

- `observed_scope`: 確認した public boundary と未確認範囲
- `violated_invariant`: 破られた API、schema、nullable / required、versioning の invariant
- `root_cause_hypotheses`: 契約破壊を生む原因候補と根拠
- `local_patch_assessment`: 局所 shim で足りるか、public seam 固定が必要か
- `exploration_scope`: contract 影響確認に必要な読む範囲
- `remediation_considerations`: 修正者が考慮すべき支配点、互換 risk、invariant tests
- 出力に tool policy、agent runtime、product code の変更義務を含めない。
- 必須出力: group, decision, score, confidence, observed_scope, violated_invariant, root_cause_hypotheses, local_patch_assessment, exploration_scope, remediation_considerations, invariant_tests, findings, evidence_paths, not_reviewed
- 出力 field 要件: {"group": "contract_compatibility", "decision": "pass / fail / blocked のいずれか。score > 0.85 の場合だけ pass", "score": "既存利用者、外部API、内部API、DB schema、event payload への互換性を 0.0 から 1.0 で返す", "confidence": "根拠十分性を 0.0 から 1.0 で返す", "observed_scope": "確認した public boundary、schema、payload、未確認範囲を分ける", "violated_invariant": "破られた API、schema、nullable / required、versioning の invariant を書く。未特定なら不明理由を書く", "root_cause_hypotheses": "契約破壊を生む原因候補、根拠、反証余地を分ける", "local_patch_assessment": "局所 shim で足りるか、public seam 固定が必要か、局所修正が危険な理由を書く", "exploration_scope": "contract 影響確認に追加で読むべき範囲と読まない範囲を書く", "remediation_considerations": "修正者が考慮すべき candidate control points、互換 risk、durable fix signal を返す。修正範囲を命令しない", "invariant_tests": "破れた contract invariant を固定する test 観点を返す", "findings": "API、schema、public method、event payload、nullable / required、versioning だけを書く"}

## 完了規約

- 対象 review 観点の finding、score、根拠、残留 risk が返却されている。
- 権限・信頼境界系の hard gate は score 相殺せず明示されている。
- API request / response の互換性を確認した。
- DB schema、event payload、queue message を確認した。
- public method と error code を確認した。
- nullable / required と versioning を確認した。
- violated invariant と root cause hypothesis を分けた。
- local patch assessment と invariant tests を返した。
- 内部実装の綺麗さを主判定にしなかった。
- completion signal: 契約・互換性の score、破られた invariant、root cause hypothesis、local patch assessment、根拠が返っている
- residual risk key: not_reviewed

## 停止規約

- 内部実装の綺麗さ
- テストの十分性
- 可読性
- パフォーマンス最適化
- 停止時は不足項目、衝突箇所、reroute 先を返す。
