---
name: updating-docs
description: "`finalization-module` 内で Claude 本体が使う docs 正本化作業プロトコル。呼び出し元の docs 正本化判断後に、人間承認済み docs-only 成果物を正本へ反映する判断基準を提供する。"
---
# Updating Docs

## 目的

`updating-docs` は作業プロトコルである。
Claude 本体が `finalization-module` の docs 正本化判断後に、人間承認済み 成果物 を docs 正本へ反映するための、正本、承認確認、検証 の見方を提供する。

人間可読な実行境界、引き継ぎ、停止 / 戻し はこの skill を正本にする。

## 対応ロール

- `finalization-module` から呼ばれた Claude 本体が直接使う（サブエージェントを起動しない）。
- 返却先は 呼び出し元 とする。
- 担当成果物は `updating-docs` の出力規約で固定する。

## 呼び出し元から渡される情報

- 呼び出し元: docs 正本化を依頼した agent または人間。
- docs正本化起動入力: 呼び出し元が docs 正本化判断後に渡す根拠、承認、正本化対象。
- 承認記録: 人間が docs 正本化を承認した記録。
- 承認済み成果物: docs 正本へ反映してよい成果物。
- 正本化対象: 更新してよい docs 正本。
- 非必須入力: 検証コマンド、根拠 docs を受け取る。
- 必須成果物: docs 正本化起動入力、承認済み docs-only 成果物、`/Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md` を受け取る。

## 作業前に読む正本

- 実行境界は本 skill と [finalization-module/SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.claude/skills/finalization-module/SKILL.md) に従う。
- docs index: [index.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md)
- 詳細仕様正本: [detail-specs](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/README.md)
- 画面設計書正本: [screen-design](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/README.md)
- 画面設計書雛形: [template.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/screen-design/screens/template.md)
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。

## 担当ロールが判断してよい範囲

- 呼び出し元の docs 正本化判断後にだけ正本化へ進む
- 人間承認済み 成果物 だけを反映する
- docs-only 対象範囲 を超えない
- implementation-scope を docs 正本へ自動昇格しない
- `detail-specs` は詳細仕様単位で作り、画面単位または個別ユースケース単位へ独断で分割しない
- `detail-specs` へ反映する内容は、`detail-spec-diff.md`、`screen-design-diff.<screen-id>.md`、実装結果、任意のレビュー結果から恒久仕様だけを製本する
- `screen-design-diff.<screen-id>.md` は、人間承認済みの場合だけ `docs/screen-design/screens/<screen-id>.md` へ反映する
- 画面設計書正本へ反映する内容は、`docs/screen-design/screens/template.md` の項目に合う画面内容だけに限定する
- 未確定仕様を独断で補完しない

- docs 正本化起動入力を根拠として残す
- 承認 記録 を根拠として残す
- 正本 と task 内成果物 を分ける
- 検証 結果を残す

## skill が扱わない対象

- 作業流れ、skill、エージェント実行定義、プロダクトコード、プロダクトテストは変更しない。
- 呼び出し元の docs 正本化判断前の正本化と未承認 draft の正本化は扱わない。
- implementation-scope を docs 正本へ自動昇格しない。
- task 内の実画面確認結果を docs 正本へそのまま昇格しない。
- 実装指示、テスト手順、agent handoff を画面設計書正本へ昇格しない。
- スキーマ移行、DB 移行、基盤移行、cutover 手順は `detail-specs` へ昇格しない。
- `detail-specs` へ移す対象は、承認済み `detail-spec-diff.md` にある恒久仕様だけにする。
- プロダクト実装を同時に進めない。

## 返す成果物

- 判断結果: docs 正本化の完了、未完了、停止の判定を返す。
- 根拠参照: docs 更新の根拠にした承認記録と成果物を返す。
- 不足情報: docs 正本化を完了できない不足項目を返す。
- 次判断材料: 呼び出し元が次を判断できる材料を返す。
- 引き継ぎ先: 呼び出し元を返す。
- 渡す対象範囲: docs 更新結果、検証、残り 不足を返す。
- 変更 docs: 更新した docs ファイルを返す。
- 更新した正本: 反映した 正本 を返す。
- 確認結果: 実行した 検証 と未実行理由を返す。
- 残留不足: 未反映、未確認、判断待ちを返す。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコード変更の指示を含めない。

## 作業を完了できる条件

- 出力規約を満たし、次の 実行者 が再解釈なしで判断できる。
- 不足情報または停止理由がある場合は明示されている。
- docs 正本化起動入力を確認した。
- 人間承認 記録 を確認した。
- 承認済み 成果物 と 正本 対象 を対応づけた。
- 画面設計差分を反映する場合は、対象の `docs/screen-design/screens/<screen-id>.md` と対応づけた。
- 検証 結果と 残り 不足 を記録した。
- 必須 根拠: docs 正本化起動入力、承認 記録、根拠成果物パス、検証結果。
- 完了判断材料: docs 正本が 承認済み 成果物 と同期している。
- 残留リスク: 未反映、未確認、判断待ちが返っている。

## 作業を止める条件

- docs 正本化起動入力が未確認の時
- 作業流れ / skill / エージェント実行定義 や skill / agent を変更する時
- プロダクトコードやプロダクトテストの変更が必要な時
- 人間承認 が不足している時
- docs 正本化起動入力が分からない場合は停止する。
- 承認 がない場合は停止する。
- 作業流れ 変更なら呼び出し元へ戻す。
- プロダクト 実装が必要なら呼び出し元へ戻す。
- 停止時は不足項目、衝突箇所、戻し先を返す。
- docs 正本化起動入力が不足する場合は停止する。
- 承認が不足する場合は停止する。
- プロダクト実装が必要な場合は停止する。
- 作業流れ / skill / エージェント実行定義 の変更が必要な場合は停止する。
- docs-only 対象範囲ではない場合は停止する。
