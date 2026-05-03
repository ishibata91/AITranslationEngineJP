# 詳細仕様: 翻訳成果物出力

- `upper_scenario_id`: `translation-output-artifact`
- `status`: `approved`
- `source_plan`: `docs/exec-plans/completed/translation-output-artifact/plan.md`
- `scenario_source`: `docs/exec-plans/completed/translation-output-artifact/scenario-design.md`
- `ui_source`: `docs/exec-plans/completed/translation-output-artifact/ui-design.md`
- `implementation_source`: `docs/exec-plans/completed/translation-output-artifact/implementation-scope.md`, `work_history/runs/2026-05-03-translation-output-artifact-run/README.md`
- `review_source`: `docs/exec-plans/completed/translation-output-artifact/reviewback.behavior.yaml`, `docs/exec-plans/completed/translation-output-artifact/reviewback.contract.yaml`, `docs/exec-plans/completed/translation-output-artifact/reviewback.trust-boundary.yaml`, `docs/exec-plans/completed/translation-output-artifact/reviewback.state-invariant.yaml`, `docs/exec-plans/completed/translation-output-artifact/reviewback.responsibility-boundary.yaml`

## 要約

- 利用者は完了済み翻訳 job から出力候補を選び、翻訳結果を確認して xTranslator 互換 XML を出力できる。
- Output Review は出力準備状態、拒否理由、結果要約、差分 preview、成果物状態、再出力状態を同じ上位シナリオとして扱う。
- 成果物出力は既存の翻訳結果から XML と出力行を再構成する。本文再翻訳、provider 実行、xTranslator 本体の自動操作は扱わない。

## 対象

- 対象利用者は、Skyrim Mod 翻訳成果物を確認して配布形式へ出したい利用者である。
- 開始条件は、body phase が Completed であり、job-level 状態も `Completed` である翻訳 job が存在することである。
- 完了状態は、出力結果、row count、file path、generated_at、compatibility summary、artifact status を確認できることである。
- 主要データは `TRANSLATION_ARTIFACT`、`XTRANSLATOR_OUTPUT_ROW`、`JOB_TRANSLATION_FIELD`、入力ファイル出自、出力ステータス分布である。

## 仕様

- Output Review は completed job list、selected job summary、input provenance summary、output readiness、拒否理由、result summary、output status distribution を表示する。
- 出力候補は body phase Completed、job-level `Completed`、field result 整合、output status 整合を満たす job だけにする。
- 未完了、失敗中、`Canceled`、field result 不整合、status 不整合の job では、artifact 生成を開始せず、成功状態の `TRANSLATION_ARTIFACT` と `XTRANSLATOR_OUTPUT_ROW` を作らない。
- xTranslator row は `EDID`、`REC`、`FIELD`、`FORMID`、`Source`、`Dest`、`Status` を row 単位で持つ。各 row は 1 つの `JOB_TRANSLATION_FIELD` に対応する。
- 内部 `cached` は xTranslator `Status=1` へ写像する。辞書置換である事実は xTranslator `Status` とは別の内部 summary に残す。
- 未定義 status、必須列欠損、重複 row 候補、致命的な構造違反は、成功 row または成功 artifact に混入させない。
- XML は UTF-8 として parse 可能であり、対象ゲームに対応した root element と `<String>` 子要素を持つ。
- Skyrim SE は `SSETranslator`、Skyrim LE は `TESVTranslator` を root element にできる。
- XML 特殊文字と日本語 text は壊さず、local parser で互換性を検査する。real xTranslator 起動は必須にしない。
- diff preview は translation unit 単位で Source、Dest、Status、row 反映内容、stale reason、再出力可否を表示する。
- 参照不能または古い row は正しい preview として表示しない。差分から開始する操作は artifact 再出力であり、本文再翻訳ではない。
- 同じ job の再出力では現行 artifact を更新または置換し、同一 field の row を重複作成しない。
- row validation、XML serialization、file write、artifact 保存に失敗した場合、artifact は成功状態にならない。失敗理由、failed stage、retryable flag を確認できる。
- target count 0 は output readiness false の理由にしない。row count 0 の completed job でも artifact summary と output status を確認できる。
- 文字列サイズ上限超過、翻訳非推奨 field、RACE 先頭スペース、末尾スペースなどの xTranslator 互換上の危険値は compatibility summary に出す。
- 出力処理は AI provider、network、secret store を必須経路にしない。
- UI、DTO、summary、structured log、debug log、runtime event へ secret、API key、復号可能値、provider raw payload、過剰な本文全文を出さない。
- 監査要約は artifact id、operation kind、row count、digest、error kind、status を中心にする。

## 受け入れ根拠

- `SCN-TOA-001`: completed job と result summary を確認する。
- `SCN-TOA-002`: 未完了または不整合 job から artifact を作らない。
- `SCN-TOA-003`: xTranslator row 必須列と status mapping を生成する。
- `SCN-TOA-004`: UTF-8 の xTranslator XML を生成して parser で検査する。
- `SCN-TOA-005`: diff preview と再出力必要状態を確認する。
- `SCN-TOA-006`: 同じ job の再出力で artifact と row を重複作成しない。
- `SCN-TOA-007`: 生成失敗を成功成果物として扱わず回復可能にする。
- `SCN-TOA-008`: 本文翻訳対象 0 件の completed job でも出力結果を確認する。
- `SCN-TOA-009`: xTranslator 互換上の危険値を検出して summary に出す。
- `SCN-TOA-010`: 出力処理で provider、network、secret を使わず redaction を守る。
- human decision は `approved` として plan に記録されている。
- 最終検証は backend-local、frontend-local、scenario-gate、system-test、UI 証跡が pass として plan に記録されている。
- 5 観点 reviewback はすべて `review_status: no_issue`、`must_fix_open: false`、`max_level: none` である。

## UI 契約由来の恒久仕様

- 表示項目は completed job list、selected job summary、input provenance summary、output readiness、拒否理由、result summary、output status distribution、diff preview、artifact status、row count、generated_at、file path、re-output state、compatibility summary、redacted error summary、operation summary である。
- 主要操作は completed job 選択、diff preview 表示、xTranslator XML 出力、出力済み artifact の再出力、summary から対象 unit への移動である。
- 出力 action は output readiness true、row validation pass、出力先 path valid の時だけ有効にする。
- 再出力 action は existing artifact があり、output readiness true で、stale または再出力可能状態の時だけ有効にする。
- invalid job、`Canceled` job、field result 不整合、status 不整合では出力 action を無効にし、理由を隣接表示する。
- loading、empty、not-ready、ready、preview-ready、generating、success、failed、stale、re-output-needed を状態差分として扱う。
- target count 0、row count 0、readonly path、XML parse failure、compatibility warning は区別して表示する。
- 長い plugin 名、file path、Source、Dest、error reason は折り返し、省略表示、詳細展開で overflow を防ぐ。
- mobile では job list、summary、diff preview を縦順にし、action は該当 section 内に置く。
- 出力不可理由は色だけで伝えず、短い日本語文言と status icon を併用する。
- disabled action の理由は tooltip または隣接文言で確認できる。
- diff preview は追加、変更、欠損をテキスト label で区別する。

## 対象外

- 本文翻訳フェーズの再翻訳、field result 編集、provider 実行、Job Run UI の再設計。
- xTranslator 本体の自動操作と、real xTranslator import smoke の必須化。
- プロダクトコード、プロダクトテスト、`.codex`、implementation-scope 自体の正本昇格。
- 複数世代の artifact revision 履歴の正本化。
