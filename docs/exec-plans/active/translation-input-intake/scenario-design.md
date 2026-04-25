# Scenario Design: translation-input-intake

- `skill`: scenario-design
- `status`: approved
- `source_plan`: `./plan.md`
- `ui_source`: `./ui-design.md`
- `final_artifact_path`: `docs/scenario-tests/translation-input-intake.md`
- `topic_abbrev`: `TII`

## Fixed Requirements

- `must_pass_requirements`:
  - xEdit 抽出 JSON を入力データとして登録できる。
  - 登録済み input file の一覧を参照できる。
  - 1 入力データを 1 翻訳ジョブ候補として識別できる。
  - 入力データから翻訳レコードと翻訳フィールドを展開し、件数とカテゴリを確認できる。
  - 抽出 JSON を正本とし、SQLite 側の入力キャッシュを削除しても再構築できる。
- `non_goals`:
  - 翻訳ジョブ作成、翻訳フェーズ実行、AI API 実行は含めない。
  - 訳文、出力ステータス、出力成果物の生成は含めない。
  - docs 正本化、product code、product test、implementation-scope は扱わない。
  - 後続の `translation-job-setup`、各翻訳 phase、`translation-output-artifact` の要件を先取りしない。

## Detail Requirement Coverage

各抽象要件について、必要な詳細要求タイプを `explicit`、`derived`、`not_applicable`、`deferred`、`needs_human_decision` に分類する。
`needs_human_decision` は 0 件である。

`scenario-design.requirement-coverage.json` を正本にする。


### `REQ-TII-001` xEdit 抽出 JSON の入力データ登録

- `source_requirement`: xEdit 抽出 JSON を入力データとして登録し、1 入力データを 1 翻訳ジョブ候補として識別できる。
- `requirement_kind`: operation
- `needs_human_decision`: なし
- `fixed_decisions`: 1 file = 1 入力データ。同一 hash は拒否。不正 JSON、非 xEdit JSON、必須 field 欠落は登録前に全体拒否する。

### `REQ-TII-002` 翻訳レコードと翻訳フィールドへの展開

- `source_requirement`: 入力データから翻訳レコード、翻訳フィールド、再構築可能な入力キャッシュを作る。
- `requirement_kind`: persistence
- `needs_human_decision`: なし
- `fixed_decisions`: 未定義 RecordType + SubrecordType は異常として警告し、入力の正本性確認のため非翻訳対象として観測可能に保持する。

### `REQ-TII-003` 入力キャッシュの再構築

- `source_requirement`: 入力キャッシュを削除しても抽出 JSON から再構築できる。
- `requirement_kind`: workflow
- `needs_human_decision`: なし
- `fixed_decisions`: 初期受け入れは小 fixture のみ固定する。

### `REQ-TII-004` Input Review UI での入力確認

- `source_requirement`: app-shell から Input Review を開き、取り込み結果を確認する。
- `requirement_kind`: display
- `needs_human_decision`: なし
- `fixed_decisions`: Input Review はページ内で完結させる。app-shell 導線の詳細は dashboard-and-app-shell 側へ deferred とする。

## Human Decision Questionnaire

未回答質問はない。回答済み判断は `scenario-design.requirement-coverage.json` を正本にする。

## Risks

- `implementation_risks`:
  - `translation unit` は usecase 由来の語であり、scenario では正本語彙の翻訳レコードと翻訳フィールドへ寄せる必要がある。
  - 入力 cache と抽出 JSON 正本の境界を曖昧にすると、後続 job setup で 1 job = 1 input の対応が崩れる。
  - `JOB_TRANSLATION_FIELD` の訳文、出力ステータスを入力取り込みで作ると、翻訳 phase の責務を先取りする。
- `test_data_risks`:
  - xEdit 抽出 JSON の最小 fixture schema は小 fixture 前提で固定する必要がある。
  - 未定義 field、重複 import、不正 JSON は fixture を分ける必要がある。
  - browser file input は環境により absolute path ではなく bare filename だけを渡すため、実ランタイム経路を fixture だけで代替しない。
  - frontend response の `warnings`、`categories`、`sampleFields` などの配列項目は null で返る可能性を guard する。
  - app-shell の screen-design 正本がなく、UI system test の入口導線は dashboard-and-app-shell 側に残る。

## Rules

- ケース ID は `SCN-TII-NNN` 形式にする。
- Markdown table は使わず、1 ケースごとの縦型ブロックで書く。
- `期待結果` は観測可能な結果にする。
- `needs_human_decision` が残る scenario matrix を human review へ進めない。
- `not_applicable` と `deferred` は理由なしで通さない。
- paid な real AI API を前提にしない。

## Scenario Matrix

### SCN-TII-001 xEdit 抽出 JSON を登録して input file 一覧に表示する

- `分類`: 正常系
- `観点`: xEdit 抽出 JSON を入力データとして登録し、job 候補として識別できる。
- `事前条件`: 基盤データ管理が成立している。小さな xEdit 抽出 JSON fixture がある。
- `手順`:
  1. Input Review で 1 つの xEdit 抽出 JSON file を登録する。
  2. input file 一覧を開く。
  3. 登録済み入力データの詳細を確認する。
- `期待結果`:
  1. 入力データが 1 件増える。
  2. 入力データが 1 翻訳ジョブ候補として識別できる。
  3. 出自情報として file path、file name、file hash、import timestamp を確認できる。
- `観測点`: UI 一覧、backend query、入力データ ID。
- `fake_or_stub`: fixed xEdit JSON fixture、temp DB。
- `責務境界メモ`: job 作成、translation phase、output artifact は実行しない。

### SCN-TII-002 翻訳レコードと翻訳フィールドの件数とカテゴリを確認する

- `分類`: 正常系
- `観点`: 入力データから翻訳レコードと翻訳フィールドを展開し、件数とカテゴリを観測できる。
- `事前条件`: 複数 RecordType と SubrecordType を含む xEdit 抽出 JSON fixture がある。
- `手順`:
  1. fixture を入力データとして登録する。
  2. Input Review で取り込み結果の概要を開く。
  3. sample field を確認する。
- `期待結果`:
  1. 翻訳レコード件数が表示される。
  2. 翻訳フィールド件数が表示される。
  3. カテゴリ別件数が表示される。
  4. 未定義 RecordType + SubrecordType は異常として警告され、非翻訳対象として観測できる。
  5. 訳文と出力ステータスは表示、保存されない。
- `観測点`: UI 概要、repository query、record / field count。
- `fake_or_stub`: fixed xEdit JSON fixture、field definition fixture。
- `責務境界メモ`: `JOB_TRANSLATION_FIELD` は入力取り込みの保存対象にしない。

### SCN-TII-003 入力キャッシュ削除後に抽出 JSON から再構築する

- `分類`: 正常系
- `観点`: SQLite 側の入力キャッシュを削除しても、抽出 JSON 正本から再構築できる。
- `事前条件`: 登録済み入力データと抽出 JSON 正本がある。未完了 job 参照はない。
- `手順`:
  1. 入力キャッシュを削除する。
  2. 抽出 JSON から入力キャッシュを再構築する。
  3. 再構築前後の件数とカテゴリを比較する。
- `期待結果`:
  1. 入力データの job 候補識別が維持される。
  2. 翻訳レコード件数、翻訳フィールド件数、カテゴリが一致する。
  3. 再構築は新しい入力データを増やさない。
- `観測点`: cache row count、UI 件数、再構築結果。
- `fake_or_stub`: fixed xEdit JSON fixture、temp DB。
- `責務境界メモ`: cache 削除対象の job 参照条件は後続 job lifecycle と衝突しない範囲だけ確認する。

### SCN-TII-004 不正 JSON または非 xEdit JSON を拒否する

- `分類`: 主要失敗系
- `観点`: 不正な入力を登録済み job 候補にしない。
- `事前条件`: 壊れた JSON fixture と、xEdit 抽出形式ではない JSON fixture がある。
- `手順`:
  1. 不正 JSON を登録する。
  2. 非 xEdit JSON を登録する。
  3. input file 一覧と error 表示を確認する。
- `期待結果`:
  1. 登録前に全体拒否される。
  2. 壊れた入力が後続 job 候補として誤表示されない。
  3. UI に拒否理由の種別が表示される。
- `観測点`: UI error、backend error、input file 一覧件数。
- `fake_or_stub`: invalid JSON fixture、non-xEdit JSON fixture。
- `責務境界メモ`: error 文言の粒度は Q-001 の回答後に固定する。

### SCN-TII-005 同一抽出 JSON の再登録を扱う

- `分類`: 境界条件
- `観点`: 同一入力の再登録で重複 job 候補を作るかどうかを一貫して扱う。
- `事前条件`: 同一内容の xEdit 抽出 JSON fixture がある。
- `手順`:
  1. fixture を登録する。
  2. 同じ fixture を再登録する。
  3. input file 一覧と件数を確認する。
- `期待結果`:
  1. 同一 hash の再登録は拒否される。
  2. 入力データ件数は増えない。
  3. UI の状態と backend の永続化結果が一致する。
- `観測点`: input file 一覧、入力データ件数、error または confirmation。
- `fake_or_stub`: same-hash fixture。
- `責務境界メモ`: 既存 job 参照がある入力の上書きは後続 job setup の判断に回す。

### SCN-TII-006 Input Review の状態差分を確認する

- `分類`: UI 状態
- `観点`: loading、empty、error、disabled、progress、retry、success の各状態で確認作業が破綻しない。
- `事前条件`: Input Review UI が実装されている。app-shell 導線の詳細は dashboard-and-app-shell 側に委ねる。
- `手順`:
  1. 未登録状態を表示する。
  2. 登録中状態を表示する。
  3. 登録成功状態を表示する。
  4. 登録失敗状態を表示する。
- `期待結果`:
  1. 未登録状態では空状態と登録 action が表示される。
  2. 登録中は重複操作ができない。
  3. 成功時は一覧、件数、カテゴリ、sample field を確認できる。
  4. 失敗時は retry できる。
- `観測点`: browser surface、UI text、button enablement、focus。
- `fake_or_stub`: UI fixture state、fixed backend response。
- `責務境界メモ`: visual polish は実装後 human review で確認する。

### SCN-TII-007 browser file input が bare filename でも source file missing にしない

- `分類`: 実ランタイム主要失敗系
- `観点`: frontend file input から absolute path ではなく bare filename だけが渡る環境でも、初回登録を source file missing と誤判定しない。
- `事前条件`: browser surface の file input で Lucien など実ファイル相当の JSON を選択できる。frontend から backend へ渡る request には file name と file content または backend が解決可能な source handle がある。
- `手順`:
  1. Input Review の file input から JSON file を選択する。
  2. frontend が backend import request を作る。
  3. backend が request を import する。
- `期待結果`:
  1. backend は bare filename をそのまま OS path として読まない。
  2. file content または解決可能な source handle で import する。
  3. request が bare filename だけで content も source handle もない場合は invalid request として拒否し、source file missing にはしない。
  4. source file missing は cache rebuild 時に保存済み正本が見つからない場合だけ返す。
- `観測点`: browser file input、backend import request、error kind、input file 一覧。
- `fake_or_stub`: browser file input fixture、bare filename request、content付き request。
- `責務境界メモ`: OS absolute path 取得を browser に要求しない。Wails / browser surface の差分を backend contract で吸収する。

### SCN-TII-008 response の配列項目が null でも frontend は落ちない

- `分類`: UI 主要失敗系
- `観点`: backend response の `warnings`、`categories`、`sampleFields` などが null の場合でも、frontend は空配列として扱い spread error を起こさない。
- `事前条件`: Input Review UI が backend response を受け取れる。null 配列項目を含む response fixture がある。
- `手順`:
  1. `warnings: null`、`categories: null`、`sampleFields: null` の response fixture を frontend gateway へ返す。
  2. frontend usecase / presenter / store が response を view model へ変換する。
  3. Input Review を表示する。
- `期待結果`:
  1. frontend は null 配列項目を `[]` に正規化する。
  2. spread error が発生しない。
  3. UI は warning なし、カテゴリなし、sample field なしの状態を表示できる。
- `観測点`: frontend unit test、browser console error、Input Review 表示。
- `fake_or_stub`: null array response fixture。
- `責務境界メモ`: backend は配列を空配列で返すことを基本とするが、frontend も外部境界として null guard を持つ。

## Acceptance Checks

- `REQ-TII-001`: `SCN-TII-001`、`SCN-TII-004`、`SCN-TII-005`
- `REQ-TII-002`: `SCN-TII-002`
- `REQ-TII-003`: `SCN-TII-003`
- `REQ-TII-004`: `SCN-TII-006`、`SCN-TII-007`、`SCN-TII-008`

## Validation Commands

- `python3 scripts/scenario/requirement_gate.py docs/exec-plans/active/translation-input-intake/scenario-design.md --coverage docs/exec-plans/active/translation-input-intake/scenario-design.requirement-coverage.json --report-out docs/exec-plans/active/translation-input-intake/scenario-design.requirement-gate.md --questionnaire-out docs/exec-plans/active/translation-input-intake/scenario-design.questions.md`
- `python3 scripts/harness/run.py --suite scenario-gate`
- implementation 後候補: `go test ./internal/...`
- implementation 後候補: `npm run test -- --run`
- implementation 後候補: `agent-browser open http://localhost:34115`

## Open Questions

- なし。
