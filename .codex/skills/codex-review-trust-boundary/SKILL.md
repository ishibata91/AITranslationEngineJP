---
name: codex-review-trust-boundary
description: Codex 実装後 review の権限・信頼境界グループ作業プロトコル。hard gate として扱う。
---
# Codex Review Trust Boundary

## 目的

ユーザー、tenant、role、外部入力、secret の境界を越えていないかを見る。
他観点の高 score で相殺してはいけないため hard gate とする。

## 対応ロール

- `review_trust_boundary` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` の review aggregation とする。
- owner artifact は `codex-review-trust-boundary` の出力規約で固定する。

## 入力規約

- 認証と認可
- tenant isolation と admin 権限
- user-controlled input
- secret 漏洩と PII
- SQL injection、XSS、SSRF、file upload、外部 URL
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。
- 必須入力: review_target_diff, implementation_scope_path, implementation_result
- 任意入力: final_validation_result, touched_files

## 外部参照規約

- agent runtime と tool policy は [review_trust_boundary.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/review_trust_boundary.toml) の `allowed_write_paths` / `allowed_commands` とする。
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-trust-boundary/SKILL.md

## 内部参照規約

## 判断規約

`score > 0.85` を pass とする。
この group は `hard_gate: true` を返し、fail を average score で相殺しない。

## 出力規約

- `observed_scope`: 確認した trust boundary と未確認範囲
- `violated_invariant`: 破られた auth、secret、tenant、外部入力の invariant
- `root_cause_hypotheses`: trust boundary 破壊を生む原因候補と根拠
- `local_patch_assessment`: 局所ガードで足りるか、境界の再固定が必要か
- `exploration_scope`: 追加で読むべき認可、入力、secret、外部接続範囲
- `remediation_considerations`: 修正者が考慮すべき支配点、hard gate risk、invariant tests
- 出力に tool policy、agent runtime、product code の変更義務を含めない。
- 必須出力: group, decision, score, confidence, hard_gate, observed_scope, violated_invariant, root_cause_hypotheses, local_patch_assessment, exploration_scope, remediation_considerations, invariant_tests, findings, evidence_paths, not_reviewed
- 出力 field 要件: {"group": "trust_boundary", "decision": "pass / fail / blocked のいずれか。score > 0.85 の場合だけ pass", "score": "権限・信頼境界の安全性を 0.0 から 1.0 で返す", "confidence": "根拠十分性を 0.0 から 1.0 で返す", "hard_gate": "true。低 score は他観点の高 score で相殺しない", "observed_scope": "確認した trust boundary、外部入力、secret、認可経路、未確認範囲を分ける", "violated_invariant": "破られた auth、secret、tenant、外部入力の invariant を書く。未特定なら不明理由を書く", "root_cause_hypotheses": "trust boundary 破壊を生む原因候補、根拠、反証余地を分ける", "local_patch_assessment": "局所ガードで足りるか、境界の再固定が必要か、局所修正が危険な理由を書く", "exploration_scope": "追加で読むべき認可、入力、secret、外部接続範囲と読まない範囲を書く", "remediation_considerations": "修正者が考慮すべき candidate control points、hard gate risk、durable fix signal を返す。修正範囲を命令しない", "invariant_tests": "破れた trust boundary invariant を固定する test 観点を返す", "findings": "認証、認可、tenant isolation、外部入力、secret、SQL injection、XSS、SSRF、外部URL、admin権限、PII だけを書く"}

## 完了規約

- 対象 review 観点の finding、score、根拠、残留 risk が返却されている。
- 権限・信頼境界系の hard gate は score 相殺せず明示されている。
- 認証、認可、tenant isolation を確認した。
- user-controlled input と外部 URL を確認した。
- secret、admin 権限、PII を確認した。
- SQL injection、XSS、SSRF、file upload を確認した。
- violated invariant と root cause hypothesis を分けた。
- local patch assessment と invariant tests を返した。
- hard gate failure を他観点で相殺しなかった。
- completion signal: 権限・信頼境界の hard gate score、破られた invariant、root cause hypothesis、local patch assessment、根拠が返っている
- residual risk key: not_reviewed

## 停止規約

- 実装の短さ
- 読みやすさ
- 性能
- テスト妥当性
- 停止時は不足項目、衝突箇所、reroute 先を返す。
