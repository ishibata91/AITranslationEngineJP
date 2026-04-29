---
name: codex-review-behavior
description: Codex 実装後 review の挙動正しさグループ作業プロトコル。
---
# Codex Review Behavior

## 目的

変更後のコードが PR の目的どおりに振る舞うかを見る。
diff から取得した実コードを、正解の挙動ベクトルにどの程度近いかで score 化する。

## 対応ロール

- `review_behavior` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` の review aggregation とする。
- owner artifact は `codex-review-behavior` の出力規約で固定する。

## 入力規約

- 正常系の挙動
- 条件分岐と境界値
- 例外系
- 既存挙動との差分
- bug 修正の場合の原因対応
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。
- 必須入力: review_target_diff, implementation_scope_path, implementation_result
- 任意入力: final_validation_result, touched_files

## 外部参照規約

- agent runtime と tool policy は [review_behavior.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/review_behavior.toml) の `allowed_write_paths` / `allowed_commands` とする。
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-behavior/SKILL.md

## 内部参照規約

## 判断規約

`score > 0.85` を pass とする。
仕様にない入力や不明な期待値は confidence を下げ、score と混同しない。

## 出力規約

- `observed_scope`: 確認した実コード、経路、未確認範囲
- `violated_invariant`: PR 目的、既存挙動、受け入れ条件のどれが破られたか
- `root_cause_hypotheses`: 症状を生む原因候補と根拠
- `local_patch_assessment`: 局所修正で閉じるか、他層へ波及するか
- `exploration_scope`: 追加で読むべき範囲と読まない範囲
- `remediation_considerations`: 修正者が考慮すべき支配点、risk、invariant tests
- 出力に tool policy、agent runtime、product code の変更義務を含めない。
- 必須出力: group, decision, score, confidence, observed_scope, violated_invariant, root_cause_hypotheses, local_patch_assessment, exploration_scope, remediation_considerations, invariant_tests, findings, evidence_paths, not_reviewed
- 出力 field 要件: {"group": "behavior_correctness", "decision": "pass / fail / blocked のいずれか。score > 0.85 の場合だけ pass", "score": "PR 目的に対する実コードの挙動一致度を 0.0 から 1.0 で返す", "confidence": "根拠十分性を 0.0 から 1.0 で返す", "observed_scope": "確認した実コード、主要経路、未確認範囲を分ける", "violated_invariant": "PR目的、既存挙動、受け入れ条件のどれが破られたかを書く。未特定なら不明理由を書く", "root_cause_hypotheses": "症状を生む原因候補、根拠、反証余地を分ける", "local_patch_assessment": "局所修正で閉じるか、他層へ波及するか、局所修正が危険な理由を書く", "exploration_scope": "追加で読むべき範囲と読まない範囲を書く", "remediation_considerations": "修正者が考慮すべき candidate control points、durable fix signal、local patch risk を返す。修正範囲を命令しない", "invariant_tests": "破れた挙動 invariant を固定する test 観点を返す", "findings": "正常系、条件分岐、境界値、例外系、既存挙動差分、bug 原因対応だけを書く"}

## 完了規約

- 対象 review 観点の finding、score、根拠、残留 risk が返却されている。
- 権限・信頼境界系の hard gate は score 相殺せず明示されている。
- PR 目的と実コードの主要経路を照合した。
- 正常系、条件分岐、境界値、例外系を確認した。
- 既存挙動との差分を確認した。
- bug 修正の場合は原因対応を確認した。
- violated invariant と root cause hypothesis を分けた。
- local patch assessment と invariant tests を返した。
- 命名、関数分割、テスト網羅性を主判定にしなかった。
- completion signal: 挙動正しさの score、破られた invariant、root cause hypothesis、local patch assessment、根拠が返っている
- residual risk key: not_reviewed

## 停止規約

- 命名
- 関数分割
- 読みやすさ
- テスト網羅性
- コードスタイル
- 停止時は不足項目、衝突箇所、reroute 先を返す。
