---
name: codex-review-trust-boundary
description: Codex 実装後 レビュー の権限・信頼境界グループ作業プロトコル。必須判定条件 として扱う。
---
# Codex Review Trust Boundary

## 目的

ユーザー、tenant、role、外部入力、secret の境界を越えていないかを見る。
他観点の高評価で相殺してはいけないため、強制停止条件として扱う。

## 対応ロール

- `review_trust_boundary` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` の レビュー 集約 とする。
- 担当成果物は `codex-review-trust-boundary` の出力規約で固定する。

## 入力規約

- レビュー対象差分: 実装後 レビュー の対象になる差分を受け取る。
- 実装目的: レビュー対象差分が満たすべき目的を受け取る。
- implementation-scope の場所: 承認済み実装範囲の参照先を受け取る。
- 実装結果: 実装 agent が返した実装結果を受け取る。
- 最終検証結果: `implement_lane` が確認した最終検証結果を受け取る。
- 変更ファイル: レビュー対象差分に含まれる変更ファイル一覧を受け取る。
- 作業計画フォルダ: `docs/exec-plans/active/<task-id>/` を受け取る。
- レビューYAMLパス: `docs/exec-plans/active/<task-id>/reviewback.trust-boundary.yaml` を受け取る。
- 差し戻しYAMLパス: `work_history/runs/<run>/review-reject-trust-boundary.yaml` を受け取る。

## 外部参照規約

- エージェント実行定義と実行境界は [review_trust_boundary.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/review_trust_boundary.toml) に従う。
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

権限・信頼境界観点表は次を拘束する。
この観点は強制停止条件として扱い、他観点の評価で相殺しない。

| 観点 | 確認内容 |
| --- | --- |
| 認証 | 未認証利用者が保護対象へ到達しないか |
| 認可 | user、role、admin 権限の境界を越えていないか |
| tenant isolation | tenant 間の参照、更新、検索結果が混ざらないか |
| 外部入力 | user-controlled 入力が検証、正規化、境界固定なしに使われないか |
| secret / PII | secret、API key、個人情報が保存、表示、ログ、URL に漏れないか |
| injection | SQL injection、XSS、SSRF、file upload、外部 URL の危険経路がないか |

重大度指標は次を拘束する。

| レベル | 意味 |
| --- | --- |
| `blocker` | マージ・リリースを止めるべき問題 |
| `critical` | 修正必須だが、条件付きで先に進める可能性がある問題 |
| `major` | 品質・保守性・仕様整合に大きく影響する問題 |
| `minor` | 局所的な改善で済む問題 |
| `nit` | 修正してもよいが、必須ではない問題 |

## 判断規約

- レビュー問題の重大度は内部参照規約の重大度指標から選ぶ。
- `blocker`、`critical`、`major` は `fix_required: true` にする。
- `minor`、`nit` は `fix_required: false` にする。
- 未解決の `fix_required: true` がある場合は `must_fix_open: true` にする。
- `max_level` は未解決指摘の最大重大度にする。
- この観点の失敗は、他観点の平均評価で相殺しない。
- `assessment.hard_gate` は常に `true` にする。

## 非対象規約

- 実装の短さ、読みやすさ、性能は主判定にしない。
- テスト妥当性は、権限・信頼境界の直接根拠になる場合だけ扱う。
- 修正範囲の命令やプロダクトコード変更の指示は出力しない。

## 出力規約

- レビューYAML: `docs/exec-plans/active/<task-id>/reviewback.trust-boundary.yaml` を作成、追記、解決更新、削除する。
- レビューYAML基本項目: `reviewback_version`、`task_id`、`review_yaml_path`、`viewpoint`、`reviewer_agent`、`review_status`、`must_fix_open`、`max_level` を記録する。
- レビューYAML状態値: `review_status` は `no_issue`、`issues_open`、`stopped` のいずれかを記録する。
- レビューYAML評価項目: `assessment.safety`、`assessment.evidence_sufficiency`、`assessment.hard_gate` を記録する。
- レビューYAML確認範囲: `checked_scope.summary`、`checked_scope.files`、`unchecked_scope` を記録する。
- レビューYAML指摘: `issues` に `id`、`level`、`title`、`problem`、`reason`、`violated_condition`、`cause_candidate`、`evidence`、`local_fix_assessment`、`additional_check_scope`、`fix_considerations`、`verification`、`fix_required`、`status` を記録する。
- レビューYAML禁止出力: `forbidden_outputs` にツール権限変更、エージェント実行定義変更、プロダクトコード変更指示、修正範囲命令の禁止を記録する。
- 差し戻しYAML: `review_status` が `issues_open` または `stopped` の場合は、ワークフロー改善用ログを `work_history/runs/<run>/review-reject-trust-boundary.yaml` に追記する。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコード変更の指示、修正範囲の命令を含めない。

## 完了規約

- `reviewback.trust-boundary.yaml` に 対象 レビュー 観点の指摘、安全性評価、根拠、残留リスクが記録されている。
- 権限・信頼境界系の強制停止条件は、他観点の高評価で相殺せず明示されている。
- 内部参照規約の権限・信頼境界観点表を確認した。
- 破られた不変条件と原因候補を分けた。
- 局所修正評価と不変条件テスト観点を返した。
- 強制停止条件の失敗を他観点で相殺しなかった。
- 完了判断材料として、`must_fix_open`、`max_level`、安全性評価、`hard_gate: true`、破られた不変条件、原因候補、局所修正評価、根拠が記録されている。
- 残留リスクとして、未確認範囲と理由が記録されている。
- `review_status` が `issues_open` または `stopped` の場合は、差し戻し YAML が `work_history/runs/<run>/review-reject-trust-boundary.yaml` に追記されている。

## 停止規約

- `レビュー対象差分` が不足する場合は停止する。
- `実装目的` が不足する場合は停止する。
- `レビューYAMLパス` が不足する場合は停止する。
- `差し戻しYAMLパス` が不足する場合は停止する。
- 外部成果物 が不足または衝突する場合は停止する。
- 権限・信頼境界以外の観点を主判定にしそうな場合は停止する。
- 停止時は不足項目、衝突箇所、戻し先を返す。
