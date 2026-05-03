# Scenario Candidates: translation-output-artifact / failure

- `generator`: `failure`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TOA`
- `candidate_count`: 11

## Generator Scope

- `viewpoint`: 失敗観点。出力開始不可、xTranslator 行再構成失敗、整合性違反、保存失敗、再出力回復を候補化する。
- `included_sources`:
  - `docs/exec-plans/active/translation-output-artifact/plan.md`
  - `tasks/usecases/translation-output-artifact.yaml`
  - `docs/spec.md`
  - `docs/er.md`
  - `docs/references/xtranslator_ref.md`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md`
- `excluded_sources`: プロダクトコード、プロダクトテスト、docs 正本変更、`.codex` 変更、他観点候補、最終 `scenario-design.md` の採否判断。
- `generation_notes`: 候補は `designer` が採否、統合、競合解消を行う前提で書く。エラー名、状態名、実行段階は最終固定しない。

## Candidate Scenarios

### CAND-TOA-001 完了済みでない job では成果物出力を開始しない

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:20-26`
  - `docs/spec.md:149-150`
  - `docs/spec.md:197-198`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:386-400`
- `viewpoint`: failure
- `candidate scenario id`: `CAND-TOA-001`
- `actor`: Output Review を操作するユーザー
- `failure start condition`: 本文翻訳フェーズが未完了、失敗中、または field result 不整合で output readiness が成立しない。
- `trigger`: completed job として扱えない job で output artifact 生成を試行する。
- `rejected operation`: xTranslator 互換成果物生成の開始。
- `expected error`: 完了済み job ではない理由、または field result 不整合の理由を表示し、`TRANSLATION_ARTIFACT` と `XTRANSLATOR_OUTPUT_ROW` を成功状態で作成しない。
- `expected outcome`: output artifact の状態は未生成または失敗扱いになり、result summary と再出力導線は成功として表示されない。
- `observable point`: Output Review UI、output readiness result、result summary、`TRANSLATION_ARTIFACT` 件数、`XTRANSLATOR_OUTPUT_ROW` 件数。
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `consistency_requirement`
- `adoption hint`: state-transition 候補の開始可能条件と統合できる。
- `conflict hint`: 正常系候補が未完了 job の result summary 表示を許す場合、出力開始可否と競合する。

### CAND-TOA-002 Canceled の途中成功結果を出力成果物へ使わない

- `source requirement`:
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:30-32`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:354-378`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:397-400`
- `viewpoint`: failure
- `candidate scenario id`: `CAND-TOA-002`
- `actor`: Output Review を操作するユーザー
- `failure start condition`: `Canceled` 後に途中成功した訳文と出力ステータスが残っている。
- `trigger`: Canceled job を completed job と同じ一覧または選択導線から出力対象にしようとする。
- `rejected operation`: Canceled job の成果物生成、diff preview 生成、再出力開始。
- `expected error`: Canceled job は成果物出力対象外である理由を表示し、途中成功結果を output readiness に使わない。
- `expected outcome`: Canceled job の `JOB_TRANSLATION_FIELD` は参照されても、成功成果物として保存されない。
- `observable point`: Output Review UI、output readiness 不可理由、`JOB_TRANSLATION_FIELD`、`TRANSLATION_ARTIFACT`。
- `related detail requirement type`: `state_requirement`, `failure_handling_requirement`, `consistency_requirement`
- `adoption hint`: CAND-TOA-001 と統合して terminal job guard の失敗系にできる。
- `conflict hint`: lifecycle 候補が Canceled job の再出力を回復導線として扱う場合、designer が業務ルールを分離する必要がある。

### CAND-TOA-003 xTranslator 必須行項目が欠ける row を出力しない

- `source requirement`:
  - `docs/spec.md:43`
  - `docs/spec.md:65-67`
  - `docs/er.md:73-77`
  - `docs/references/xtranslator_ref.md:53-64`
- `viewpoint`: failure
- `candidate scenario id`: `CAND-TOA-003`
- `actor`: Output Review を操作するユーザー
- `failure start condition`: `EDID`、`REC`、`FIELD`、`FORMID`、`Source`、`Dest`、`Status` のいずれかを再構成できない field がある。
- `trigger`: 必須項目欠損を含む job で xTranslator 互換 XML を生成する。
- `rejected operation`: 欠損 row を含む output artifact の成功保存。
- `expected error`: 欠損した項目名と対象 field を result summary で確認できる。
- `expected outcome`: 欠損 row は xTranslator 互換成果物へ混入せず、artifact 全体は成功完了にならない。
- `observable point`: row build result、result summary、`XTRANSLATOR_OUTPUT_ROW`、生成 XML の row count。
- `related detail requirement type`: `data_requirement`, `failure_handling_requirement`, `compatibility_requirement`
- `adoption hint`: xTranslator row 生成の契約固定ケースにできる。
- `conflict hint`: 正常系が欠損 row をスキップして部分成果物を成功扱いにする場合、成功条件と競合する。

### CAND-TOA-004 内部出力ステータスを xTranslator `Status` へ写像できない

- `source requirement`:
  - `docs/spec.md:29-32`
  - `docs/spec.md:43`
  - `docs/spec.md:65-67`
  - `docs/references/xtranslator_ref.md:63-73`
- `viewpoint`: failure
- `candidate scenario id`: `CAND-TOA-004`
- `actor`: Output Review を操作するユーザー
- `failure start condition`: 内部出力ステータスが未定義、または xTranslator の `Status` 値へ写像できない。
- `trigger`: 写像不能な status を含む field で成果物生成を実行する。
- `rejected operation`: `Status` 欠損または範囲外の XML row 出力。
- `expected error`: status 写像不能の件数、対象 field、辞書置換由来かどうかの内部観測情報を確認できる。
- `expected outcome`: 内部 `cached` は `Status=1` に写像される。写像不能な status は成功成果物へ混入しない。
- `observable point`: output status summary、row build result、`XTRANSLATOR_OUTPUT_ROW.Status`、result summary。
- `related detail requirement type`: `consistency_requirement`, `data_requirement`, `compatibility_requirement`
- `adoption hint`: status mapping の境界値ケースにできる。
- `conflict hint`: `cached` と xTranslator `Status=1` を同一の意味として表示する候補がある場合、内部観測情報の分離要件と競合する。

### CAND-TOA-005 `JOB_TRANSLATION_FIELD` と output row が一対一に対応しない

- `source requirement`:
  - `docs/er.md:23`
  - `docs/er.md:37-40`
  - `docs/er.md:76-77`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:386-400`
- `viewpoint`: failure
- `candidate scenario id`: `CAND-TOA-005`
- `actor`: Output Review を操作するユーザー
- `failure start condition`: 1 つの `JOB_TRANSLATION_FIELD` に対して output row が重複する、または対応 row が欠落する。
- `trigger`: 重複または欠落を含む row set で result summary、diff preview、XML 出力を生成する。
- `rejected operation`: 不整合 row set の成功保存と成功表示。
- `expected error`: 対象 field 件数、row 件数、重複件数、欠落件数を result summary で確認できる。
- `expected outcome`: artifact は成功完了にならず、field と row の対応を修正するまで再出力待ちになる。
- `observable point`: result summary、diff preview、row count、`JOB_TRANSLATION_FIELD`、`XTRANSLATOR_OUTPUT_ROW`。
- `related detail requirement type`: `consistency_requirement`, `data_requirement`, `failure_handling_requirement`
- `adoption hint`: CAND-TOA-003 と同じ row validation group に統合できる。
- `conflict hint`: diff preview が欠落 row を表示対象外として扱う正常系候補と競合する可能性がある。

### CAND-TOA-006 Skyrim 向け XML root を決められない設定では XML を確定しない

- `source requirement`:
  - `docs/spec.md:65-67`
  - `docs/references/xtranslator_ref.md:45-51`
  - `tasks/usecases/translation-output-artifact.yaml:14-16`
- `viewpoint`: failure
- `candidate scenario id`: `CAND-TOA-006`
- `actor`: Output Review を操作するユーザー
- `failure start condition`: 入力ファイルの出自情報または出力設定から Skyrim LE / Skyrim SE の root 要素を決められない。
- `trigger`: root 要素未確定のまま xTranslator 互換 XML を生成する。
- `rejected operation`: root 要素が曖昧な XML 成果物の成功保存。
- `expected error`: Skyrim 向け root 要素を決められない理由を確認できる。
- `expected outcome`: `SSETranslator` または `TESVTranslator` のどちらを使うか未確定の XML は成功成果物にならない。
- `observable point`: output configuration summary、生成 XML header、result summary、入力ファイルの出自情報。
- `related detail requirement type`: `compatibility_requirement`, `data_requirement`, `failure_handling_requirement`
- `adoption hint`: external-integration 候補または actor-goal 候補の出力形式選択と統合できる。
- `conflict hint`: `docs/spec.md` は対象ゲームを Skyrim とするが、root 要素は Skyrim LE / SE で分かれるため、人間判断候補になる。

### CAND-TOA-007 XML として壊れる text を含む成果物を成功扱いにしない

- `source requirement`:
  - `docs/spec.md:42-43`
  - `docs/spec.md:65-67`
  - `docs/references/xtranslator_ref.md:30-43`
  - `docs/references/xtranslator_ref.md:55-64`
- `viewpoint`: failure
- `candidate scenario id`: `CAND-TOA-007`
- `actor`: Output Review を操作するユーザー
- `failure start condition`: `Source` または `Dest` に XML として壊れる文字列、または lossless に保持できない埋め込み要素が含まれる。
- `trigger`: XML 特殊文字、改行、埋め込み要素を含む field で XML 出力を生成する。
- `rejected operation`: parse 不能または round trip 不能な XML の成功保存。
- `expected error`: XML serialization 失敗、または round trip 検証失敗の対象 field を確認できる。
- `expected outcome`: 生成 XML が壊れる場合は artifact success にせず、保護要素と原文、訳文の lossless 性を維持する。
- `observable point`: XML serialization result、round trip parse result、result summary、対象 row の `Source` / `Dest`。
- `related detail requirement type`: `compatibility_requirement`, `failure_handling_requirement`, `data_requirement`
- `adoption hint`: XML 出力 boundary の lower-level 検証候補にできる。
- `conflict hint`: xTranslator が許容する文字列 escaping の詳細をどこまで検証するかは designer が判断する必要がある。

### CAND-TOA-008 output artifact または row の保存失敗を成功扱いにしない

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:16-19`
  - `tasks/usecases/translation-output-artifact.yaml:24-26`
  - `docs/er.md:71-77`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:247-249`
- `viewpoint`: failure
- `candidate scenario id`: `CAND-TOA-008`
- `actor`: Output Review を操作するユーザー
- `failure start condition`: `TRANSLATION_ARTIFACT`、`XTRANSLATOR_OUTPUT_ROW`、または生成 XML の保存で永続化失敗が起きる。
- `trigger`: output artifact 生成中に保存失敗を注入する。
- `rejected operation`: 保存失敗後の成功ステータス表示、成功 XML の公開、成功 result summary の保存。
- `expected error`: 保存失敗の発生箇所、保存済み件数、未保存件数、再出力可否を確認できる。
- `expected outcome`: artifact は成功完了にならず、再出力状態から回復操作を選べる。
- `observable point`: artifact status、row count、result summary、再出力状態、永続化結果。
- `related detail requirement type`: `failure_handling_requirement`, `recovery_requirement`, `data_requirement`
- `adoption hint`: recovery 候補または operation-audit 候補と統合できる。
- `conflict hint`: 保存失敗時に部分保存を残すか破棄するかは人間判断候補である。

### CAND-TOA-009 diff preview が参照不能または古い row を表示しない

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:17-18`
  - `tasks/usecases/translation-output-artifact.yaml:24-25`
  - `tasks/usecases/translation-output-artifact.yaml:31-35`
  - `docs/er.md:76-77`
- `viewpoint`: failure
- `candidate scenario id`: `CAND-TOA-009`
- `actor`: Output Review を操作するユーザー
- `failure start condition`: diff preview が参照する output row、source field、dest field、または artifact revision が欠落または古い。
- `trigger`: 欠落 row、古い artifact、または不整合 revision を選択して diff preview を開く。
- `rejected operation`: 不正確な diff preview の表示と、その preview からの再生成操作。
- `expected error`: preview を表示できない理由、または再読み込みが必要な理由を表示する。
- `expected outcome`: ユーザーは不正確な差分を正しい成果物として確認できない。
- `observable point`: diff preview UI、selected artifact revision、row count、result summary、再生成操作の enablement。
- `related detail requirement type`: `failure_handling_requirement`, `observability_requirement`, `consistency_requirement`
- `adoption hint`: actor-goal 候補の確認操作と統合できる。
- `conflict hint`: preview 失敗時に result summary だけ表示するか、画面全体を失敗表示にするかは UI 設計側の判断になる。

### CAND-TOA-010 再出力の再送や retry で成果物と row を重複作成しない

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:19`
  - `tasks/usecases/translation-output-artifact.yaml:26`
  - `tasks/usecases/translation-output-artifact.yaml:31-35`
  - `docs/er.md:73-77`
- `viewpoint`: failure
- `candidate scenario id`: `CAND-TOA-010`
- `actor`: Output Review を操作するユーザー
- `failure start condition`: 前回の成果物生成が失敗、タイムアウト、または結果確認前に中断している。
- `trigger`: 同じ completed job と同じ出力条件で再出力、retry、開始再送を行う。
- `rejected operation`: 同じ field set に対する成果物と output row の二重作成、または古い失敗 artifact の成功扱い。
- `expected error`: 既存 artifact の状態、再出力対象、上書きまたは新規作成の扱いを確認できる。
- `expected outcome`: 再送後も result summary と row count は一貫し、成功済み artifact と失敗 artifact の区別が残る。
- `observable point`: artifact list、artifact status、row count、再出力状態、result summary。
- `related detail requirement type`: `recovery_requirement`, `冪等性_requirement`, `consistency_requirement`
- `adoption hint`: lifecycle 候補の再出力導線と統合できる。
- `conflict hint`: 再出力を上書き扱いにするか履歴追加扱いにするかは人間判断候補である。

### CAND-TOA-011 xTranslator 互換上の危険値を無警告で成功出力しない

- `source requirement`:
  - `docs/spec.md:65-67`
  - `docs/references/xtranslator_ref.md:163-171`
  - `docs/references/xtranslator_ref.md:167-170`
- `viewpoint`: failure
- `candidate scenario id`: `CAND-TOA-011`
- `actor`: Output Review を操作するユーザー
- `failure start condition`: エンコード後の文字列サイズ上限超過、翻訳非推奨 field、RACE 先頭スペース欠落、末尾スペースなど、xTranslator 互換上の危険値がある。
- `trigger`: 危険値を含む row set で成果物生成または diff preview 確認を行う。
- `rejected operation`: 危険値を検出しないままの成功出力、または互換性警告なしの確認完了。
- `expected error`: 危険値の種別、対象 row、互換性への影響を result summary で確認できる。
- `expected outcome`: 重大な互換性違反は成功成果物にしない。警告扱いにできる危険値は designer が最終判断する。
- `observable point`: compatibility validation result、result summary、diff preview、`XTRANSLATOR_OUTPUT_ROW`。
- `related detail requirement type`: `compatibility_requirement`, `boundary_requirement`, `failure_handling_requirement`
- `adoption hint`: xTranslator 互換性検証の境界値ケースにできる。
- `conflict hint`: どの危険値を拒否し、どの危険値を警告だけにするかは人間判断候補である。

## Open Notes

- `human decision candidate`:
  - `CAND-TOA-006`: Skyrim の root 要素を `SSETranslator` 固定にするか、入力出自から `TESVTranslator` と分岐するかは未確定である。
  - `CAND-TOA-008`: 保存失敗時に部分保存を残すか、artifact 単位で破棄するかは未確定である。
  - `CAND-TOA-010`: 再出力を上書き扱いにするか、履歴追加扱いにするかは未確定である。
  - `CAND-TOA-011`: xTranslator 互換上の危険値を拒否、警告、許容のどれに分けるかは未確定である。
- `merge candidate`:
  - `CAND-TOA-001` と `CAND-TOA-002` は output readiness / terminal job guard として統合できる可能性がある。
  - `CAND-TOA-003`、`CAND-TOA-004`、`CAND-TOA-005`、`CAND-TOA-007`、`CAND-TOA-011` は xTranslator row validation として統合できる可能性がある。
  - `CAND-TOA-008` と `CAND-TOA-010` は再出力回復として統合できる可能性がある。
- `rejection candidate`:
  - 正常系の単純な裏返しだけで観測点が増えない候補は、designer 側で他観点候補へ merge できる。
  - 実装方式だけを固定する候補は、本成果物ではなく implementation-scope 以降へ送る。
