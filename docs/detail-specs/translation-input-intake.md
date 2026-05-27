# 詳細仕様: 翻訳入力取り込み

- `detail_spec_id`: `translation-input-intake`
- `status`: `approved`
- `source_artifacts`: `docs/exec-plans/completed/translation-input-intake/plan.md`
- `implementation_artifacts`: `docs/exec-plans/completed/translation-input-intake/implementation-scope.md`, `work_history/runs/2026-04-25-translation-input-intake.yaml-tasks-usecases-translation-input-intake.yaml-propose-plans-U-run/codex.md`, `work_history/runs/2026-04-25-translation-input-intake.yaml-tasks-usecases-translation-input-intake.yaml-propose-plans-U-run/copilot.md`
- `review_artifacts`: `docs/exec-plans/completed/translation-input-intake/plan.md`

## 親要件と仕様

### `translation-input-intake-REQ-001` xEdit 抽出 JSON を入力データとして登録できる

親要件:
利用者は xEdit 抽出 JSON を登録し、翻訳開始前に入力データを判断できる。

仕様:
- 登録対象は xEdit 抽出 JSON である。
- 1 回の登録操作は 1 入力データとして登録する。
- 入力データは、元ファイル名、元ファイルの同一性、登録時点を出自情報として保持する。
- 同一内容の再登録は、新規入力データ登録の対象にする。
- 同一 hash の再登録は新しい入力データとして扱い、既存入力データと区別できる。
- 不正 JSON、非 xEdit JSON、必須項目欠落は登録拒否理由にする。
- 初回登録では、利用者が選択したファイル名だけの情報を invalid request として扱う。
- 初回登録要求がファイル名だけの場合は、登録拒否結果にする。

### `translation-input-intake-REQ-002` 登録結果を判断できる

親要件:
利用者は登録済み入力データの規模、分類、代表例、失敗理由を判断できる。

仕様:
- 登録した入力データから翻訳レコードと翻訳フィールドを展開し、件数とカテゴリを判断できる。
- 未定義 RecordType と SubrecordType の組み合わせは警告として扱い、非翻訳対象として観測可能に保持する。
- 利用者は、入力データの登録状態、出自、翻訳対象の規模、失敗分類を判断できる。
- 代表例は RecordType、SubrecordType、FormID、EditorID、原文を根拠として示せる。
- 警告、カテゴリ、代表例が未提供の場合は空の値として扱う。

### `translation-input-intake-REQ-003` 入力データの出自と登録結果を保持する

親要件:
利用者は登録済み入力データの出自と登録結果を、翻訳ジョブ作成前に確認できる。

仕様:
- 訳文、出力ステータス、出力成果物は入力取り込みの対象外にする。
- 登録済み入力データは、出自情報、登録状態、翻訳レコード件数、翻訳フィールド件数、カテゴリを保持する。
- 入力データの登録結果は、同じ入力データを参照する翻訳ジョブ作成で利用できる。
- `source_file_missing` は、入力ファイル参照が必要な登録処理で参照不能な場合の失敗分類として扱う。

### `translation-input-intake-REQ-004` ジョブ設定へ進む条件を登録状態に合わせる

親要件:
登録結果を確認した利用者は、登録済みまたは警告ありの入力データだけをジョブ設定へ進められる。

仕様:
- 入力登録の結果は入力データ作成である。
- 翻訳ジョブ作成はジョブ設定の結果として扱う。
- 登録済み入力データは、利用者がジョブ設定で選択した時に翻訳ジョブの対象になる。
- 登録済みまたは警告ありの入力データでは、ジョブ設定へ進む条件が成立する。
- 失敗した入力データは、ジョブ設定へ進む拒否理由にする。

### `translation-input-intake-REQ-005` 取り込み状態の差分を判断できる

親要件:
利用者は登録、再試行、失敗状態を判別できる。

仕様:
- 利用者は xEdit 抽出 JSON の登録、登録済み入力データの選択、失敗時の再試行、ジョブ設定への移行を行える。
- 登録は基盤データ管理が成立し、登録中以外の時だけ成立する。
- 再試行は失敗状態と再試行可能な失敗分類の時だけ成立する。
- 未登録、処理中、成功、失敗、利用不可、再試行可能を状態差分として扱う。
- 不正 JSON、入力ファイル参照不能は失敗分類として区別する。

## 根拠

- human decision は `approved` として plan に記録されている。
- 実装完了は backend、frontend、final validation の 3 作業単位が完了した状態として plan に記録されている。
- 旧最終検証は structure と旧設計 gate が PASS として plan に記録されている。
