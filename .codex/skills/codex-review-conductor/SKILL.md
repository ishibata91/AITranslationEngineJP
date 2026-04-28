---
name: codex-review-conductor
description: Codex 実装後 review conductor 知識 package。観点グループ別 review agent を並列 spawn し、score、hard gate、remediation aggregation を集約する判断基準を提供する。
---

# Codex Review Conductor

## 目的

`codex-review-conductor` は知識 package である。
Copilot 実装完了後に人間が `codex exec` で起動した `review_conductor` が、観点グループ別 review agent を並列 spawn し、結果を remediation aggregation として集約する判断基準を提供する。

## いつ参照するか

- Copilot completion packet の review request payload から Codex review を行う時
- diff から取得した実コードを観点別に score 化する時
- hard gate、overall decision、修正判断材料を集約する時

## 知識範囲

- 観点グループ別 review agent の spawn
- score threshold `score > 0.85`
- conflict priority `trust_boundary > behavior > contract > state_invariant`
- hard gate の非相殺
- blocking findings、confidence notes、remediation aggregation の集約
- Copilot が再解釈なしで受け取れる `copilot_action`

## 原則

- review agent は context を引き継がず並列 spawn する
- 各 agent には同じ diff、scope、validation result を渡す
- overall pass は `strict_pass` または `priority_override_pass` の場合だけ返す
- 権限・信頼境界は hard gate とし、平均 score で相殺しない
- 観点間 finding が競合する時は `trust_boundary > behavior > contract > state_invariant` の順で裁定する
- 優先度で退けた finding は削除せず、`priority_overrides` と `residual_risks` に残す
- review agent は修正範囲を命令せず、observed scope、violated invariant、root cause hypothesis、local patch risk、candidate control point を返す
- conductor は各 review agent の raw result を落とさず `reviewer_result_bundle` に保持する
- conductor の aggregation は、どの観点 result のどの field から作ったかを `aggregation_trace` に残す
- conductor は観点ごとの局所最適を調停し、primary failure mode、dominant invariant、minimum durable fix boundary を返す
- Codex review は実装修正、docs 正本化、design review を行わない

## 標準パターン

1. Copilot payload に diff、implementation-scope、implementation result、final validation result があることを確認する。
2. diff、scope、implementation result、final validation result が不足している場合は、観点別 review を起動せず `copilot_action` を返す。
3. `review_behavior`、`review_contract`、`review_trust_boundary`、`review_state_invariant` を context 継承なしで並列 spawn する。
4. 各 agent の score、confidence、findings、evidence、violated invariant、root cause hypothesis、local patch assessment、remediation considerations を raw result として受け取る。
5. raw result を `reviewer_result_bundle` にそのまま保持し、不足 field があれば `information_loss_notes` に記録する。
6. `overall_score` は group score の最小値にする。
7. すべての group score が 0.85 より大きい場合は `decision_basis: strict_pass` にする。
8. 非 trust_boundary の finding が上位観点と競合する場合は、優先度で上位観点を採用し `decision_basis: priority_override_pass` にできる。
9. 複数観点の local patch assessment を統合し、primary failure mode、dominant invariant、final exploration scope、minimum durable fix boundary を返す。
10. 統合時に採用、統合、退けた観点 signal を `aggregation_trace` に残す。
11. blocking findings、priority overrides、confidence notes、residual risks、copilot action、remediation handoff を分けて返す。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は agent contract に従う。

## Copilot 返答

`review_conductor` は、人間が Copilot へ戻せる形で次 action を返す。
Copilot は戻された `copilot_action` を再解釈せず分岐する。
`copilot_action` は次のいずれかにする。

- `close`: `strict_pass` で residual がない場合に返す
- `report_residual`: `priority_override_pass` または低 confidence residual が残る場合に返す
- `fix`: Copilot が remediation aggregation を基に chosen strategy と chosen scope を決めるべき場合に返す
- `rerun_validation`: validation 不足で観点別 review を起動しない場合に返す
- `rerun_codex_review`: diff、scope、implementation result など payload 不足で観点別 review を起動しない場合に返す

`rerun_validation` と `rerun_codex_review` は早期 return である。
この場合、観点別 review agent は spawn しない。

## Handoff Format

`review_conductor` の返答は、score や blocking findings から始めない。
人間と Copilot が「何が問題で、なぜ局所修正では足りないか」を先に読める順にする。
最終返答では、次の top-level section または top-level JSON field をこの順に必ず置く。

`fix` の時は、次の順で返す。

1. `problem_statement`: 問題の中心を 2-3 文で書く。
2. `why_it_matters`: 破られる仕様、不変条件、ユーザー影響を書く。
3. `remediation_handoff`: `primary_failure_mode`、`dominant_invariant`、`root_cause_chain`、`evidence_map`、`why_not_narrower`、`minimum_durable_fix_boundary`、`why_not_wider`、`invariant_tests`、`review_decision` を含める。
4. `reviewer_result_bundle`: 4 観点の raw result を group 名ごとにそのまま置く。
5. `aggregation_trace`: 派生 field ごとの参照元 group、source field、採用理由を置く。
6. `remediation_aggregation`: 統合結果を置く。
7. `blocking_findings`: blocking 判定だけを置く。
8. `priority_overrides`: override した finding と理由を置く。
9. `confidence_notes`: confidence 低下理由を置く。
10. `residual_risks`: 残余 risk を置く。
11. `next_actions`: `copilot_action` ごとの次 action を置く。

`blocking_findings` は `evidence_map` の補助情報として扱う。
先頭に置かない。

`remediation_handoff` は Copilot への命令ではない。
Copilot が `chosen_strategy`、`chosen_scope`、`why_not_narrower`、`why_not_wider` を作るための判断材料である。

## Lossless Aggregation

`review_conductor` は、観点別 review agent の戻り値を要約だけに変換しない。
統合用の短い説明とは別に、各観点の raw result を残す。

必須の保持単位は次である。

- `reviewer_result_bundle`: 4 観点の raw result を group 名ごとに保持する。
- `aggregation_trace`: aggregation field ごとに参照した group、source field、採用理由を保持する。
- `unselected_group_signals`: 採用しなかったが残すべき signal を保持する。
- `information_loss_notes`: 欠落 field、読めなかった evidence、低 confidence の理由を保持する。

`remediation_aggregation` と `remediation_handoff` は派生結果である。
派生結果だけを返し、raw result を落とすことは禁止する。
`reviewer_result_bundle` または `aggregation_trace` を作れない場合は、`blocked` として `rerun_codex_review` を返す。

## DO / DON'T

DO:
- diff から取得した実コードを review 対象にする
- score と confidence を分ける
- 症状、破られた invariant、root cause hypothesis、局所修正 risk を分ける
- 各観点 result を raw のまま保持し、aggregation trace で派生元を示す
- Aggregator として minimum durable fix boundary を返す
- 返答は problem statement、why it matters、failure mode、invariant、root cause chain の順にする
- hard gate failure を blocking finding にする
- blocked reason には不足 payload と再実行条件を書く
- `copilot_action` と `decision_basis` を必ず分ける

DON'T:
- 観点別 review を conductor が代替しない
- Copilot の chosen scope を conductor が最終命令しない
- symptoms の列挙だけを Copilot handoff にしない
- 観点別 raw result を要約で上書きしない
- score、group result、blocking findings を返答の先頭に置かない
- low score を平均で隠さない
- priority override した finding を消さない
- 実装修正や docs 正本化をしない

## Checklist

- [codex-review-conductor-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/codex-review-conductor/references/checklists/codex-review-conductor-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は agent contract が決める。

## Agent が持つもの

- active contract: [review_conductor.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/review_conductor/contracts/review_conductor.contract.json)
- permissions: [permissions.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/review_conductor/permissions.json)

## Maintenance

- 観点グループを増減する時は conductor contract と `.codex/README.md` を同期する。
- threshold を変える時は全 review skill の pass 条件も同期する。
