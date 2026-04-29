---
name: updating-docs
description: Codex 側の docs 正本化作業プロトコル。implementation 完了後に、human 承認済み docs-only 成果物 を正本へ反映する判断基準を提供する。
---
# Updating Docs

## 目的

`updating-docs` は作業プロトコルである。
`docs_updater` agent が implementation 完了後に human 承認済み 成果物 を docs 正本へ反映するための、正本、承認確認、検証 の見方を提供する。

人間可読な実行境界、引き継ぎ、stop / 戻し はこの skill を正本にする。

## 対応ロール

- `docs_updater` が使う。
- 返却先は 呼び出し元 または次 agent とする。
- 担当成果物は `updating-docs` の出力規約で固定する。

## 入力規約

- Codex implementation レーン の修正完了が分かっている時
- human 承認済み docs-only 成果物 を docs 正本へ移す時
- 正本化 対象 と 検証 を整理する時
- task 内成果物 と docs 正本 の対応を確認する時
- 入力に 根拠参照、担当者、承認状態が不足する場合は推測で補わない。
- 必須入力: 呼び出し元, implementation_completion_report, 承認記録, 承認済み成果物, 正本化対象
- 任意入力: 検証コマンド, 根拠 docs
- 必須 成果物: Codex implementation 完了 レポート, 承認済み docs-only 成果物, /Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md

## 外部参照規約

- エージェント実行定義とツール権限は [docs_updater.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/docs_updater.toml) の 書き込み許可 / 実行許可 とする。
- 紐づけ: [docs_updater.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/docs_updater.toml)
- エージェント実行定義: [docs_updater.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/docs_updater.toml)
- forbidden: プロダクトコード、プロダクトテスト、作業流れ / skill / エージェント実行定義の変更
- ツール権限: エージェント実行定義の 書き込み許可 / 実行許可 に従う
- docs index: [index.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md)
- 紐づけ: [docs_updater.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/docs_updater.toml)
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/updating-docs/SKILL.md

## 内部参照規約

### 拘束観点

- Codex implementation 完了 レポート の確認
- docs 正本 の選び方
- 人間承認 記録 の確認
- 承認済み 成果物 と 正本 対象 の対応
- 検証 と 残り 不足 の記録

## 判断規約

- implementation 完了後にだけ正本化へ進む
- human 承認済み 成果物 だけを反映する
- docs-only 対象範囲 を超えない
- implementation-scope を docs 正本へ自動昇格しない
- 未確定仕様を独断で補完しない

- Codex implementation 完了 レポート を根拠として残す
- 承認 記録 を根拠として残す
- 正本 と task 内成果物 を分ける
- 検証 結果を残す

## 非対象規約

- 作業流れ、skill、エージェント実行定義、プロダクトコード、プロダクトテストは変更しない。
- implementation 完了前の正本化と未承認 draft の正本化は扱わない。
- implementation-scope を docs 正本へ自動昇格しない。
- プロダクト実装を同時に進めない。

## 出力規約

- 出力は判断結果、根拠参照、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

### Handoff

- 引き継ぎ先: `implement_lane`
- 渡す対象範囲: docs 更新結果、検証、残り 不足
- 変更 docs: 更新した docs ファイルを返す。
- 更新した正本: 反映した 正本 を返す。
- 確認結果: 実行した 検証 と未実行理由を返す。
- 残留 不足: 未反映、未確認、判断待ちを返す。

## 完了規約

- 出力規約を満たし、次の 実行者 が再解釈なしで判断できる。
- 不足情報または停止理由がある場合は明示されている。
- Codex implementation 完了 レポート を確認した。
- 人間承認 記録 を確認した。
- 承認済み 成果物 と 正本 対象 を対応づけた。
- 検証 結果と 残り 不足 を記録した。
- 必須 根拠: Codex implementation 完了 レポート, 承認 記録, 根拠成果物 path, 検証結果
- 完了判断材料: implementation 完了 後の docs 正本が 承認済み 成果物 と同期している。
- 残留リスク: 未反映、未確認、判断待ちが返っている。

## 停止規約

- Codex implementation レーン の修正完了が未確認の時
- 作業流れ / skill / エージェント実行定義 や skill / agent を変更する時
- プロダクトコードやプロダクトテストの変更が必要な時
- 人間承認 が不足している時
- Codex implementation レーン の修正完了が分からない場合は停止する。
- 承認 がない場合は停止する。
- 作業流れ 変更なら `implement_lane` へ戻す。
- プロダクト 実装が必要なら `implement_lane` へ戻す。
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 拒否条件: Codex implementation 完了 不足
- 拒否条件: 承認 不足
- 拒否条件: プロダクト実装 必須
- 拒否条件: 作業流れ / skill / エージェント実行定義 変更
- 停止条件: Codex implementation レーン の修正完了が分からない
- 停止条件: docs-only 対象範囲 ではない
- 停止条件: 人間承認 がない
