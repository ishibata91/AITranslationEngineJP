# 詳細仕様: 翻訳成果物出力

- `detail_spec_id`: `translation-output-artifact`
- `status`: `approved`
- `source_artifacts`: `docs/exec-plans/completed/translation-output-artifact/plan.md`
- `implementation_artifacts`: `docs/exec-plans/completed/translation-output-artifact/implementation-scope.md`, `work_history/runs/2026-05-03-translation-output-artifact-run/README.md`
- `review_artifacts`: `docs/exec-plans/completed/translation-output-artifact/reviewback.behavior.yaml`, `docs/exec-plans/completed/translation-output-artifact/reviewback.contract.yaml`, `docs/exec-plans/completed/translation-output-artifact/reviewback.trust-boundary.yaml`, `docs/exec-plans/completed/translation-output-artifact/reviewback.state-invariant.yaml`, `docs/exec-plans/completed/translation-output-artifact/reviewback.responsibility-boundary.yaml`

## 親要件と仕様

### `translation-output-artifact-REQ-001` 完了済み翻訳ジョブから出力候補を選べる

親要件:
利用者は完了済み翻訳ジョブから出力候補を選び、出力準備状態と拒否理由を判断できる。

仕様:
- 利用者は完了済み翻訳ジョブ、選択中の翻訳ジョブ、入力出自、成果物出力条件、拒否理由、翻訳結果の分布を判断できる。
- 出力候補は、本文翻訳フェーズ `Completed`、翻訳ジョブ全体 `Completed`、翻訳結果整合、出力状態整合を満たす翻訳ジョブだけにする。
- 未完了、失敗中、`Canceled`、翻訳結果不整合、出力状態不整合の翻訳ジョブは、成果物生成の拒否理由にする。
- 未完了、失敗中、`Canceled`、翻訳結果不整合、出力状態不整合の翻訳ジョブは、成果物生成の失敗対象にする。
- 出力対象件数 0 でも成果物出力条件は成立できる。
- 出力行 0 の完了済み翻訳ジョブでも、成果物生成結果と出力状態を判断できる。

### `translation-output-artifact-REQ-002` xTranslator 互換 XML を生成できる

親要件:
利用者は翻訳結果を xTranslator 互換 XML として出力できる。

仕様:
- xTranslator 互換 XML の行は、外部形式の列名として `EDID`、`REC`、`FIELD`、`FORMID`、`Source`、`Dest`、`Status` を持つ。
- 各行は 1 つの翻訳項目に対応する。
- 内部 `cached` は xTranslator `Status=1` へ写像する。
- 辞書置換である事実は xTranslator `Status` とは別の内部情報として残す。
- 未定義の状態、必須列欠損、重複行候補、致命的な構造違反は、行または成果物の失敗理由にする。
- XML は UTF-8 として解析可能であり、対象ゲームに対応したルート要素と `<String>` 子要素を持つ。
- Skyrim SE は `SSETranslator`、Skyrim LE は `TESVTranslator` をルート要素にできる。
- XML 特殊文字と日本語テキストは、xTranslator 互換 XML の文字列として保持する。
- xTranslator 互換性は、成果物の XML 構造と行内容で満たす。

### `translation-output-artifact-REQ-003` 差分確認と再出力を扱える

親要件:
利用者は出力前後の差分、古い状態、再出力可否を判断できる。

仕様:
- 差分確認は翻訳単位ごとに原文、訳文、xTranslator 状態、反映内容、古い理由、再出力可否を判断できる。
- 参照不能または古い行は、古い差分状態として扱う。
- 差分から開始する操作は成果物再出力として扱う。
- 同じ翻訳ジョブの再出力では現行成果物を更新または置換し、同一翻訳項目の行を一意に扱う。
- 再出力は既存成果物があり、成果物出力条件が成立し、古い状態または再出力可能状態の時だけ成立する。

### `translation-output-artifact-REQ-004` 生成失敗を回復可能な失敗にする

親要件:
行検証、XML 生成、ファイル書き込み、成果物保存の失敗は回復可能な失敗として扱う。

仕様:
- 行検証、XML 生成、ファイル書き込み、成果物保存に失敗した場合、成果物は失敗対象にする。
- 失敗理由、失敗段階、再試行可否を判断できる。
- 出力は成果物出力条件、行検証成功、出力先妥当性が成立した時だけ成立する。
- 無効な翻訳ジョブ、`Canceled` の翻訳ジョブ、翻訳結果不整合、出力状態不整合では出力が拒否結果になり、理由を判断できる。
- 未準備、準備完了、差分確認可能、生成中、成功、失敗、古い状態、再出力必要を状態差分として扱う。
- 出力対象件数 0、出力行 0、読み取り専用の出力先、XML 解析失敗、互換性警告は別状態として扱う。

### `translation-output-artifact-REQ-005` xTranslator 互換上の危険値を判断できる

親要件:
利用者は成果物の xTranslator 互換上の注意点を判断できる。

仕様:
- 文字列サイズ上限超過、翻訳非推奨の翻訳項目、RACE 先頭スペース、末尾スペースなどの xTranslator 互換上の危険値を判断できる。
- 利用者は成果物状態、出力行数、生成時点、出力先、再出力状態、互換性上の注意、保護対象を含まない失敗要約を判断できる。
- 出力不可理由、再出力可否、差分確認の追加、変更、欠損は状態情報として判断できる。

### `translation-output-artifact-REQ-006` 出力処理は AIサービスと秘密値から分離する

親要件:
成果物出力は既存の翻訳結果から XML と出力行を再構成し、AIサービスや秘密値から分離する。

仕様:
- 出力処理は AIサービス、ネットワーク、秘密値管理から分離する。
- 成果物出力は既存の翻訳結果から XML と出力行を再構成する。
- 本文再翻訳、AIサービス実行、xTranslator 本体の自動操作は成果物出力の対象外にする。
- 秘密値、認証キー、復号可能値、外部サービスとの生データ、過剰な本文全文は利用者向け情報の対象外にする。
- 運用上必要な要約は、保存済みの状態事実から導出する。
- 監査用に導出できる要約は、成果物識別子、操作種別、出力行数、同一性、失敗分類、状態を中心にする。

## 根拠

- human decision は `approved` として plan に記録されている。
- 最終検証は pass として plan に記録されている。
- 5 観点 reviewback はすべて `review_status: no_issue`、`must_fix_open: false`、`max_level: none` である。
