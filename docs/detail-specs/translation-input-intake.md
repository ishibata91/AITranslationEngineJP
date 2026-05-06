# 詳細仕様: 翻訳入力取り込み

- `upper_scenario_id`: `translation-input-intake`
- `status`: `approved`
- `source_plan`: `docs/exec-plans/completed/translation-input-intake/plan.md`
- `scenario_source`: `docs/exec-plans/completed/translation-input-intake/scenario-design.md`
- `ui_source`: `docs/exec-plans/completed/translation-input-intake/ui-design.md`
- `implementation_source`: `docs/exec-plans/completed/translation-input-intake/implementation-scope.md`, `work_history/runs/2026-04-25-translation-input-intake.yaml-tasks-usecases-translation-input-intake.yaml-propose-plans-U-run/codex.md`, `work_history/runs/2026-04-25-translation-input-intake.yaml-tasks-usecases-translation-input-intake.yaml-propose-plans-U-run/copilot.md`
- `review_source`: reviewback なし。古い task のため、`plan.md` の `workflow_state: completed`、`human_review_status: approved`、完了済み作業単位、closeout 材料を完了根拠にする。

## 要約

- 利用者は xEdit 抽出 JSON を登録し、翻訳開始前に入力データ、件数、カテゴリ、代表 field を確認できる。
- 入力取り込みは登録操作ごとに入力データを作る。同じファイルまたは同じ入力内容から複数の翻訳 job を作成できる。
- 抽出 JSON を入力正本として扱い、SQLite 側の入力キャッシュは抽出 JSON から再構築できる。
- 入力取り込みは job を自動作成しない。登録結果を確認した利用者は、Job Setup へ進める。

## 対象

- 対象利用者は、Skyrim Mod 翻訳の開始前に xEdit 抽出 JSON の内容を確認したい利用者である。
- 開始条件は、基盤データ管理が成立し、xEdit 抽出 JSON を登録できる状態である。
- 完了状態は、登録済み入力データの一覧、出自情報、翻訳レコード件数、翻訳フィールド件数、カテゴリ別件数、代表 field を確認できることである。
- 主要データは入力データ、抽出 JSON 正本、翻訳レコード、翻訳フィールド、カテゴリ集計、入力キャッシュである。

## 仕様

- 登録対象は xEdit 抽出 JSON である。1 回の登録操作は 1 入力データとして登録する。
- 入力データは file path、file name、file hash、import timestamp を出自情報として保持する。
- 同一 hash の再登録は、重複入力として全体拒否しない。新しい入力データとして扱い、既存入力データと区別できる。
- Data Load の新規登録は入力データ作成だけを行い、翻訳 job は自動作成しない。
- 登録済み input は、利用者が Job Setup で選択した時に翻訳 job の対象になる。
- 不正 JSON、非 xEdit JSON、必須 field 欠落は登録前に全体拒否する。
- 登録した入力データから翻訳レコードと翻訳フィールドを展開し、件数とカテゴリを確認できる。
- 未定義 RecordType と SubrecordType の組み合わせは警告として扱い、非翻訳対象として観測可能に保持する。
- 訳文、出力ステータス、出力成果物は入力取り込みで表示または保存しない。
- 入力キャッシュを削除しても、抽出 JSON 正本から再構築できる。
- 再構築後の翻訳レコード件数、翻訳フィールド件数、カテゴリは再構築前と一致する。
- 再構築は同じ入力データを新規入力データとして重複登録しない。
- 初回登録では、browser file input 由来の bare filename を OS path として読まない。
- 初回登録 request が bare filename だけで content も source handle もない場合は invalid request として拒否する。
- `source_file_missing` は、cache rebuild 時に保存済み正本が見つからない場合だけ使う。
- backend response の `warnings`、`categories`、`sampleFields` が null の場合でも、frontend は空配列として扱う。
- 登録済みまたは警告ありの selected input では、Job Setup へ進む導線を表示できる。
- 失敗または再構築が必要な selected input では、Job Setup へ進む導線を表示しない。

## 受け入れ根拠

- `SCN-TII-001`: xEdit 抽出 JSON を登録して input file 一覧に表示する。
- `SCN-TII-002`: 翻訳レコードと翻訳フィールドの件数とカテゴリを確認する。
- `SCN-TII-003`: 入力キャッシュ削除後に抽出 JSON から再構築する。
- `SCN-TII-004`: 不正 JSON または非 xEdit JSON を拒否する。
- `SCN-TII-005`: 同一抽出 JSON の再登録を扱う。
- `SCN-TII-006`: Data Load の状態差分を確認する。
- `SCN-TII-007`: browser file input が bare filename でも source file missing にしない。
- `SCN-TII-008`: response の配列項目が null でも frontend は落ちない。
- `translation-job-management` の 2026-05-06 修正で、同じ入力データから複数の翻訳 job を作成できる方針が承認済みである。
- `translation-job-management` の 2026-05-06 修正で、Data Load 登録成功後に Job Setup へ進む導線が承認済みである。
- human decision は `approved` として plan に記録されている。
- 実装完了は backend、frontend、final validation の 3 作業単位が完了した状態として plan に記録されている。
- 最終検証は structure と scenario-gate が PASS として plan に記録されている。

## UI 契約由来の恒久仕様

- 表示項目は input file 一覧、入力データ概要、sample field、error summary である。
- input file 一覧は表示名、登録状態、登録日時、出自情報、再構築可否を表示する。
- 入力データ概要は翻訳レコード件数、翻訳フィールド件数、カテゴリ別件数を表示する。
- sample field は RecordType、SubrecordType、FormID、EditorID、原文を表示する。
- error summary は失敗種別、対象 file、再試行可否を表示する。
- 主要操作は xEdit 抽出 JSON の登録、登録済み入力データの選択、取り込み結果の再読込、入力キャッシュの再構築、失敗時の再試行、Job Setup への移動である。
- 登録 action は基盤データ管理が成立し、登録中でない時だけ有効にする。
- 再構築 action は抽出 JSON 正本の出自情報があり、cache 欠落または再構築可能状態の時だけ有効にする。
- retry action は失敗状態と retry 可能な error の時だけ有効にする。
- Job Setup へ進む action は、登録済みまたは警告ありの selected input だけで有効にする。
- loading、empty、progress、success、error、disabled、retry を状態差分として扱う。
- invalid JSON、source file missing、cache missing は error variant として区別する。
- desktop では一覧と詳細を並べて確認できる。mobile では一覧の後に詳細を積む。
- file path、FormID、EditorID、error message は折り返しまたは省略表示で overflow を防ぐ。
- file registration、rebuild、retry は keyboard 操作可能にする。
- error は色だけで伝えず text でも伝える。
- progress は screen reader が読める text を持つ。

## 対象外

- 翻訳フェーズ実行、AI API 実行、訳文生成、出力成果物生成。
- app-shell navigation 上の位置詳細。
- 実装時の一時成果物と実装運用情報。
