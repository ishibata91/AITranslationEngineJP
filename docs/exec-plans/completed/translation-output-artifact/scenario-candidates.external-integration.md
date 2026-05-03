# Scenario Candidates: translation-output-artifact / external-integration

- `generator`: `external-integration`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TOA-EI`
- `candidate_count`: 8

## Generator Scope

- `viewpoint`: `external-integration`
- `included_sources`:
  - `docs/exec-plans/active/translation-output-artifact/plan.md`
  - `tasks/usecases/translation-output-artifact.yaml`
  - `docs/spec.md`
  - `docs/er.md`
  - `docs/references/xtranslator_ref.md`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md`
- `excluded_sources`:
  - プロダクトコード。
  - プロダクトテスト。
  - docs 正本本文の変更。
  - `.codex` の変更。
- `generation_notes`:
  - 対象 usecase は完了済み job から成果物を出力するため、外部 AI provider 呼び出しは主経路に含めない。
  - xTranslator 互換仕様、XML file、status adapter、secret 非露出を外部連携境界として候補化する。
  - 採否、統合、競合解消、最終 scenario-design への反映は `designer` に残す。

## Candidate Scenarios

### CAND-TOA-EI-001 完了済み job から xTranslator XML row を生成する

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:9-19`
  - `tasks/usecases/translation-output-artifact.yaml:22-26`
  - `docs/spec.md:43-43`
  - `docs/spec.md:65-67`
  - `docs/er.md:71-77`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TOA-EI-001`
- `external boundary`: xTranslator XML row adapter。
- `actor`: Output Review を使うユーザー。
- `trigger`: 完了済み job と output artifact を選び、xTranslator 互換成果物の出力を開始する。
- `start condition`: 本文翻訳フェーズ完了後の job が選択されている。
- `expected outcome`: `JOB_TRANSLATION_FIELD` から `EDID`、`REC`、`FIELD`、`FORMID`、`Source`、`Dest`、`Status` を持つ出力 row を再構成できる。
- `observable point`: output summary、diff preview、生成 XML、`TRANSLATION_ARTIFACT`、`XTRANSLATOR_OUTPUT_ROW`。
- `fake_or_stub`: completed job fixture、field result fixture、temp DB、XML parser。
- `related detail requirement type`: `success_requirement`, `data_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: 成果物出力の主経路候補として扱える。
- `conflict hint`: job 完了条件と output readiness は lifecycle または state-transition 候補と統合する必要がある。

### CAND-TOA-EI-002 内部出力ステータスを xTranslator Status へ写像する

- `source requirement`:
  - `docs/spec.md:29-32`
  - `docs/spec.md:43-43`
  - `docs/spec.md:65-67`
  - `docs/references/xtranslator_ref.md:65-73`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TOA-EI-002`
- `external boundary`: 内部出力ステータスから xTranslator `Status` への adapter。
- `actor`: Output Review を使うユーザー。
- `trigger`: output row 生成時に field 単位の出力ステータスを XML へ変換する。
- `start condition`: field result が訳文と内部出力ステータスを保持している。
- `expected outcome`: 明示済みの `cached` は xTranslator `Status=1` に写像され、未確定の内部ステータス写像は固定済み扱いにしない。
- `observable point`: status mapping result、row-level `Status`、mapping error summary。
- `fake_or_stub`: status fixture、cached field fixture、unknown status fixture。
- `related detail requirement type`: `compatibility_requirement`, `data_requirement`, `failure_handling_requirement`
- `adoption hint`: output artifact と xTranslator 互換性の接続点として扱える。
- `conflict hint`: `cached` 以外の内部ステータス対応は人間判断候補として `designer` が質問票へ出す可能性がある。

### CAND-TOA-EI-003 入力ファイルの出自と record 識別子を XML row に保持する

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:10-16`
  - `docs/spec.md:12-15`
  - `docs/spec.md:43-43`
  - `docs/er.md:18-24`
  - `docs/er.md:27-35`
  - `docs/references/xtranslator_ref.md:53-63`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TOA-EI-003`
- `external boundary`: xEdit 抽出由来の record 識別子から xTranslator row 識別子への adapter。
- `actor`: Output Review を使うユーザー。
- `trigger`: 出力対象 field を XML row へ変換する。
- `start condition`: 1 job が 1 つの xEdit 抽出データに紐づき、翻訳 record と field が参照できる。
- `expected outcome`: 入力ファイルの出自を失わず、`EDID`、`REC`、`FIELD`、`FORMID`、`Source` が元 record と対応する。
- `observable point`: row identifier summary、generated XML、input provenance summary。
- `fake_or_stub`: xEdit extracted data fixture、translation record fixture、field fixture。
- `related detail requirement type`: `data_requirement`, `consistency_requirement`, `compatibility_requirement`
- `adoption hint`: xTranslator 互換出力の識別子欠落を防ぐ候補として扱える。
- `conflict hint`: 出自情報の監査保存は operation-audit 候補と重なる可能性がある。

### CAND-TOA-EI-004 Skyrim 向け XML file を UTF-8 で書き出す

- `source requirement`:
  - `docs/spec.md:65-67`
  - `docs/references/xtranslator_ref.md:24-43`
  - `docs/references/xtranslator_ref.md:45-52`
  - `docs/references/xtranslator_ref.md:53-63`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TOA-EI-004`
- `external boundary`: ローカル file system と xTranslator XML serializer。
- `actor`: Output Review を使うユーザー。
- `trigger`: xTranslator 互換成果物の出力先を指定し、file 書き出しを実行する。
- `start condition`: XML row が生成済みで、出力先 file path が指定されている。
- `expected outcome`: UTF-8 XML として保存され、Skyrim 向け root element と `<String>` child 要素を確認できる。
- `observable point`: file path、file write result、XML root、XML encoding、artifact state。
- `fake_or_stub`: temp filesystem、read-only path fixture、XML parser。
- `related detail requirement type`: `success_requirement`, `failure_handling_requirement`, `compatibility_requirement`
- `adoption hint`: file 出力の外部境界候補として扱える。
- `conflict hint`: `SSETranslator` と `TESVTranslator` のどちらを Skyrim 正本にするかは人間判断候補になる。

### CAND-TOA-EI-005 再出力で row と artifact state を重複させない

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:15-19`
  - `tasks/usecases/translation-output-artifact.yaml:22-26`
  - `tasks/usecases/translation-output-artifact.yaml:31-35`
  - `docs/er.md:71-77`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TOA-EI-005`
- `external boundary`: 既存 artifact file と再出力 file の境界。
- `actor`: Output Review を使うユーザー。
- `trigger`: 既に出力済みの output artifact に対して再出力を実行する。
- `start condition`: 完了済み job と既存 artifact state が存在する。
- `expected outcome`: 再出力結果の artifact state を観測でき、同じ field から重複 row が作られない。
- `observable point`: artifact state、row count、file path、re-output result summary。
- `fake_or_stub`: existing artifact fixture、completed job fixture、temp filesystem。
- `related detail requirement type`: `冪等性_requirement`, `data_requirement`, `consistency_requirement`
- `adoption hint`: 再出力導線を xTranslator file 境界として検証できる候補である。
- `conflict hint`: 上書き、別名保存、version 追加のどれを採用するかは lifecycle 候補または人間判断と競合する。

### CAND-TOA-EI-006 生成 XML を xTranslator import 互換として検査する

- `source requirement`:
  - `docs/spec.md:65-67`
  - `docs/references/xtranslator_ref.md:24-43`
  - `docs/references/xtranslator_ref.md:53-63`
  - `docs/references/xtranslator_ref.md:139-149`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TOA-EI-006`
- `external boundary`: xTranslator import 互換性検査境界。
- `actor`: Output Review を使うユーザー。
- `trigger`: 生成済み XML の互換性を検査する。
- `start condition`: xTranslator 互換 XML file が生成されている。
- `expected outcome`: real xTranslator を必須にせず、required tag、row count、FormID reference による照合可能性を検査できる。
- `observable point`: XML parse result、required tag check、FormID reference check、compatibility summary。
- `fake_or_stub`: local XML parser、xTranslator minimal fixture、known valid XML fixture。
- `related detail requirement type`: `compatibility_requirement`, `testability_requirement`, `failure_handling_requirement`
- `adoption hint`: CI または scenario validation で real 外部 tool なしに互換性を観測できる。
- `conflict hint`: real xTranslator を使う手動確認を必須にするかは designer または人間レビューの判断に残す。

### CAND-TOA-EI-007 XML escaping、UTF-8、日本語 text、互換制約を守る

- `source requirement`:
  - `docs/spec.md:42-43`
  - `docs/spec.md:65-67`
  - `docs/references/xtranslator_ref.md:30-42`
  - `docs/references/xtranslator_ref.md:61-63`
  - `docs/references/xtranslator_ref.md:163-171`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TOA-EI-007`
- `external boundary`: XML serializer と xTranslator 互換文字列制約。
- `actor`: Output Review を使うユーザー。
- `trigger`: 日本語、XML meta character、保持要素、長文を含む field を出力する。
- `start condition`: source と dest に互換性確認が必要な文字列を含む field result がある。
- `expected outcome`: XML は UTF-8 として parse でき、XML meta character は壊れず、保持要素と訳文は row 単位で追跡できる。
- `observable point`: generated XML、parser result、row-level source/dest、compatibility error summary。
- `fake_or_stub`: Japanese text fixture、XML escaping fixture、embedded element fixture、large text fixture。
- `related detail requirement type`: `boundary_requirement`, `compatibility_requirement`, `failure_handling_requirement`
- `adoption hint`: xTranslator file 互換の boundary case として扱える。
- `conflict hint`: 最大文字列サイズ超過時に全体失敗、row 除外、警告つき出力のどれにするかは人間判断候補になる。

### CAND-TOA-EI-008 成果物出力では AI provider、network、secret を使わない

- `source requirement`:
  - `tasks/usecases/translation-output-artifact.yaml:9-19`
  - `docs/spec.md:45-59`
  - `docs/spec.md:61-67`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:30-37`
  - `docs/exec-plans/completed/body-translation-phase/scenario-design.md:404-426`
- `viewpoint`: `external-integration`
- `candidate scenario id`: `CAND-TOA-EI-008`
- `external boundary`: provider、network、secret store への非接続境界。
- `actor`: Output Review を使うユーザー。
- `trigger`: provider 実行履歴を持つ completed job から output artifact を出力する。
- `start condition`: body translation result と output-ready summary が保存済みである。
- `expected outcome`: 出力処理は AI provider、network、API key を要求せず、XML、diff preview、error summary、log に secret 本体や復号可能値を出さない。
- `observable point`: fake provider call count、fake secret store assertion、generated XML、output summary、log capture。
- `fake_or_stub`: fake provider that fails on call、fake secret store、log capture、completed result fixture。
- `related detail requirement type`: `security_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: external-integration 観点の secret / network 非対象境界を明示する候補である。
- `conflict hint`: provider 実行履歴の監査表示は operation-audit 候補と重なる可能性がある。

## Open Notes

- `human decision candidate`: `SSETranslator` と `TESVTranslator` のどちらを Skyrim 向け root element の既定値にするか。
- `human decision candidate`: `cached` 以外の内部出力ステータスを xTranslator `Status=0..4` のどれへ写像するか。
- `human decision candidate`: 再出力時に既存 file を上書きするか、別 artifact として保存するか。
- `human decision candidate`: 最大文字列サイズ超過時に全体失敗、row 除外、警告つき出力のどれにするか。
- `human decision candidate`: real xTranslator での手動 import smoke を最終検証に含めるか。
- `merge candidate`: file write failure と path invalid は failure 観点の候補と統合する可能性がある。
- `merge candidate`: output readiness と再出力 state は lifecycle / state-transition 観点の候補と統合する可能性がある。
- `merge candidate`: secret 非露出と provider 履歴要約は operation-audit 観点の候補と統合する可能性がある。
- `rejection candidate`: real AI provider や network 呼び出しを output artifact の必須検証前提にする候補は非対象にできる。
