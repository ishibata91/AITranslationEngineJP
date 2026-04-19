# Scenario Design: 2026-04-19-sqlite-migration-repositories

- `skill`: scenario-design
- `status`: approved
- `source_plan`: `./plan.md`
- `requirements_source`: `./requirements-design.md`
- `ui_source`: `N/A`
- `figma_source`: `N/A`
- `final_artifact_path`: `docs/scenario-tests/sqlite-migration-repositories.md`
- `topic_abbrev`: `SMR`

## Rules

- ケース ID は `SCN-SMR-NNN` 形式にする。
- Markdown table は使わず、1 ケースごとの縦型ブロックで書く。
- `期待結果` は repository test または schema application test で観測できる結果にする。
- paid な real AI API を前提にしない。

## Scenario Matrix

### SCN-SMR-001 Fresh canonical schema

- `分類`: 正常系
- `観点`: 空の SQLite DB に canonical ER schema を作成できる。
- `事前条件`: temp directory に DB file が存在しない。
- `手順`:
  1. SQLite open 処理から schema を適用する。
  2. ER v1 の主要 table、index、unique、foreign key を検査する。
- `期待結果`:
  1. `NPC_PROFILE`, `TRANSLATION_RECORD`, `TRANSLATION_FIELD`, `TRANSLATION_JOB`, `JOB_PHASE_RUN`, `PERSONA`, `DICTIONARY_ENTRY` 相当の table が存在する。
  2. 旧 master 名 UI / service への接続変更は発生しない。
- `観測点`: schema test、`sqlite_master` query、foreign key pragma。
- `fake_or_stub`: temp DB のみ。
- `責務境界メモ`: UI、service、bootstrap は変更しない。

### SCN-SMR-002 Dictionary and persona canonical shape

- `分類`: 正常系
- `観点`: 共通 / job-local を同じ canonical table で表せる。
- `事前条件`: NPC profile、NPC record、dictionary entry、persona の test data がある。
- `手順`:
  1. 共通 dictionary と job-local dictionary を `DICTIONARY_ENTRY` 相当へ保存する。
  2. NPC profile に対する persona を `PERSONA` 相当へ保存する。
- `期待結果`:
  1. dictionary は `source_term` と `translated_term` を中心に保存される。
  2. lifecycle / scope / source により共通と job-local を区別できる。
  3. 同一 NPC profile への persona 二重保持は拒否される。
- `観測点`: repository integration test、unique constraint test。
- `fake_or_stub`: temp DB fixture のみ。
- `責務境界メモ`: legacy `REC` / `EDID` は辞書の中核情報として扱わない。

### SCN-SMR-003 Translation source persistence

- `分類`: 正常系
- `観点`: 翻訳入力元を `TranslationSourceRepository` で保存して再読込できる。
- `事前条件`: record、NPC profile、NPC 派生、field definition、field、field reference の test data がある。
- `手順`:
  1. repository で translation source を transaction 保存する。
  2. DB を reopen し、同じ translation source を読み込む。
- `期待結果`:
  1. NPC profile、NPC attributes、translation fields が欠落しない。
  2. ordered field の previous / next と別 record reference が復元される。
- `観測点`: repository integration test。
- `fake_or_stub`: small fixture。real xEdit は不要。
- `責務境界メモ`: parser、service、UI の変更は含めない。

### SCN-SMR-004 Job lifecycle and output persistence

- `分類`: 正常系
- `観点`: job / phase の状態遷移と、ジョブ内出力状態を別 repository 境界で保存できる。
- `事前条件`: translation source と translation job test data がある。
- `手順`:
  1. `JobLifecycleRepository` で job と phase run を作成する。
  2. `JobLifecycleRepository` で phase 対象 field / persona / dictionary を関連づける。
  3. `JobOutputRepository` で job translation field の訳文、出力ステータス、適用ペルソナを保存する。
- `期待結果`:
  1. job は 1 translation source だけを参照する。
  2. phase run の AI 設定、指示種別、最新外部 run ID、error が保持される。
  3. job output は `JOB_TRANSLATION_FIELD` として再読込できる。
- `観測点`: repository test、foreign key violation test。
- `fake_or_stub`: AI 実行結果は固定値。
- `責務境界メモ`: phase rerun history table は作らない。実ファイル export は repository test の対象にしない。

### SCN-SMR-005 Transaction rollback

- `分類`: 主要失敗系
- `観点`: 複数 table 保存失敗時に部分書き込みを残さない。
- `事前条件`: foreign key error を起こす invalid fixture がある。
- `手順`:
  1. repository 保存中に意図的な foreign key error を発生させる。
  2. 保存対象 table の row count を確認する。
- `期待結果`:
  1. repository は error を返す。
  2. transaction rollback により中途半端な row が残らない。
- `観測点`: repository test、row count、error wrapping。
- `fake_or_stub`: temp DB と fixture のみ。
- `責務境界メモ`: 旧 schema 互換ではなく canonical repository 境界の整合性を検証する。

## Acceptance Checks

- `SCN-SMR-001` は canonical schema 作成の acceptance を満たす。
- `SCN-SMR-002` は共通 / job-local の同一 table 表現と persona unique の acceptance を満たす。
- `SCN-SMR-003` は translation source と translation field の永続化 acceptance を満たす。
- `SCN-SMR-004` は job lifecycle と job output の永続化 acceptance を満たす。
- `SCN-SMR-005` は foreign key / transaction rollback acceptance を満たす。

## Validation Commands

- `go test ./internal/infra/sqlite ./internal/repository`
- `python3 scripts/harness/run.py --suite structure`
- optional regression: `go test ./internal/...`

## Open Questions

- 旧 master 名 UI の移行 scenario は `2026-04-19-legacy-schema-ui-migration-todo` で別途設計する。
