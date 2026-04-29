# Test Patterns

## 目的

`implementation_tester` が 単一引き継ぎ入力 1 件を プロダクトテスト で証明するための判断パターンをまとめる。
agent TOML の ツール権限 と skill の出力規約は上書きしない。

## 採用する考え方

- 実装前 implementation_tester は、承認済み受け入れ条件、公開接点、入力開始点、主要観測点、期待 結果 が固定済みの `APIテスト` だけで使う。
- 単体 test と原因未確定の 回帰 test は実装後に追加または更新する。
- 単一引き継ぎ入力 は完了条件、公開接点、test 対象、検証コマンド の順に読む。
- test は behavior を証明し、implementation detail を固定しない。
- null、空、invalid、境界、エラー経路、concurrency をリスクに応じて含める。
- E2E は `UI人間操作E2E` として critical user flow と browser 表面 の証跡に絞る。
- flaky test は arbitrary wait ではなく、明確な condition wait へ直す。

## 適用ルール

- paid real AI API は test で呼ばない。fake provider、DI seam、test bootstrap を使う。
- 承認済み実装範囲 外を直接読んで test 対象範囲 を広げない。
- test 小範囲 が渡された場合は、その 小範囲 の 完了合図 clause、公開接点、test 対象 file、検証コマンド だけを証明する。
- 文脈不足基準 は 構造判定条件 とし、証明対象の挙動、公開接点、test 対象、検証主眼、fixture/helper 方針、対象限定検証 の不足時だけ 文脈不足 を返す。
- test 小範囲 が 完了合図 clause、公開接点、test 対象 file、検証コマンド のいずれにも対応しない場合は 文脈不足 を返す。
- not_insufficient_context: 承認済み シナリオ を元に期待どおり fail する test、局所的 import 修正、既存 test file 内の軽微な確認は停止理由にしない。
- 原因未確定の 回帰 test を実装前に書く必要がある場合は停止し、post-implementation test として orchestrator へ返す。
- backend service / usecase / controller と frontend gateway / screen controller を主戦場にする。
- UI 操作証跡は `agent-browser` CLI を使う。
- プロダクトテスト runner としての Playwright は、既存 test が必要とする場合だけ使う。
- coverage と `python3 scripts/harness/run.py --suite all` は 最終検証 レーン へ defer する。
- AAA を守り、1 test は 1 behavior / 分岐 / シナリオ 結果 を証明する。

## 赤旗

- test 名前 が `works` など曖昧で、何を証明したか分からない。
- test body に条件分岐がある。
- external AI / network / clock / random に依存している。
- coverage / harness all の未実行理由を 最終検証 レーン へ渡していない。
- assertion が弱く、壊れても 通過 する。
- 列挙済み file / symbol 外を探索して 文脈 を膨らませている。
- 文脈不足 を返さず 広域 investigation を始めている。
- 判定基準 mismatch になる不安や通常の局所確認を 文脈不足 にしている。
