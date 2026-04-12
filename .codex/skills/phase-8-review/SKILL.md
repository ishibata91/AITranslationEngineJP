---
name: phase-8-review
description: 第8段階の実装レビューを担当し、実装差分が詳細設計と整合しているかだけを単発で確認する。
---

# Phase 8 Review

## Review Scope

- [ ] 承認済み design bundle と実装差分を照合し、対象 `task_id` ごとの責務、入出力、画面状態、依存方向、主要フローが一致している
- [ ] active exec-plan、承認済み UI モック artifact、承認済み Scenario テスト一覧 artifact、review 用差分図の前提を崩す差分がなく、未承認の仕様追加や仕様欠落がない
- [ ] 第7段階までの証明を見直し、主要責務と主要分岐に未証明が残っていない
- [ ] coverage は 70% 超過を前提に、`test-results/coverage-manifest.json` と関連 test 差分から、数字合わせだけの悪いテストが混入していない
- [ ] sonar MCP で open issue がなく、review 時点の品質ゲート阻害要因が残っていない

## Output

- decision: `pass` or `reroute`
- findings
- recheck
- closeout_notes

## Rules

- review は 1 回だけ行う
- active exec-plan、承認済み UI モック artifact、承認済み Scenario テスト一覧 artifact、承認済み task_id、承認済み required reading、review 用差分図を source of truth として checklist 順に照合する
- 新しい改善提案や新しい要件解釈は追加しない
- 実装差分なら第6段階へ、設計差分なら第2段階または第3段階へ差し戻す
- 承認済み design bundle にない仕様や好みで判定しない
- coverage については `python3 scripts/harness/run.py --suite coverage` と `test-results/coverage-manifest.json` を参照し、70% 超過を確認した上で次も見る
- [ ] 追加または更新された test が、承認済み design bundle にある責務、失敗条件、主要分岐の証明に結び付いている
- [ ] Wails runtime event を使う非同期処理の完了が同期 response や見かけの画面更新だけで判定されておらず、completion event の発火または受信で証明されている
- [ ] assertion が行数消化ではなく期待される振る舞い、出力、状態遷移、エラー条件を検証している
- [ ] private implementation detail、呼び出し回数、無意味に細かい内部順序、過剰な snapshot に依存していない
- [ ] mock、stub、fixture が対象責務を素通りさせておらず、line hit だけを増やす空疎な setup になっていない
- [ ] response fallback や別経路の成功が Wails runtime event 不達を隠し、完了したように見える構造になっていない
- [ ] 同種の trivial case を重複させて coverage を水増ししていない
- [ ] coverage 超過でも主要責務や主要分岐に未証明が残る場合は `pass` にしない

## Reference Use

- 着手前に `../orchestrating-implementation/references/orchestrating-implementation.to.phase-8-review.json` を参照して入力契約を確認する。
- `orchestrating-implementation` へ返す時は `references/phase-8-review.to.orchestrating-implementation.json` を返却契約として使う。
