# Scenario Candidates: translation-job-setup / operation-audit

- `generator`: `operation-audit`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `TJS`

## Generator Scope

- `viewpoint`: 翻訳ジョブ作成前後の運用確認、validation 履歴、設定再現性、秘密情報の非保存を operation-audit 観点で候補化する。
- `included_sources`: `./plan.md`、`../../../../tasks/usecases/translation-job-setup.yaml`、`../../../spec.md`、`../../../er.md`、`../../completed/translation-input-intake/scenario-design.md`、`../../completed/translation-input-intake/ui-design.md`
- `excluded_sources`: 最終 scenario matrix、candidate 採否判断、product code、product test、docs 正本、他 viewpoint generator の担当判断
- `generation_notes`: audit log の保存形式は固定しない。保存対象と非保存対象の境界、validation 理由の観測粒度、再現に必要な設定断面を候補として明示する。

## Candidate Scenarios

### CAND-TJS-001 ジョブ作成時の設定断面を後から再確認できる

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の「1 入力データごとに 1 翻訳ジョブを作成できる」「基盤参照、AI 基盤、実行方式を選択できる」、`docs/spec.md` の「1つの翻訳ジョブは1つのxEdit抽出データを対象とし、入力ファイルの出自を失わずに保持できること」「各フェーズのAPI選択、APIKeyは再入力不要で保存ができること」、`docs/er.md` の「TRANSLATION_JOB は 1 つの X_EDIT_EXTRACTED_DATA だけを参照する」「フェーズ別 AI 設定、指示構成、最終 AI 実行情報は JOB_PHASE_RUN に保持する」
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-TJS-001`
- `actor`: 運用者
- `trigger`: 作成済み翻訳ジョブの設定根拠を後から確認したい。
- `expected outcome`: ジョブに紐づく入力データ、選択した共通辞書、共通ペルソナ、AI runtime、実行方式、validation 通過時点を再確認できる。再現に不要な秘密値は表示も保存もしない。
- `observable point`: ジョブ詳細 UI、設定要約、入力データ出自、`TRANSLATION_JOB` と `JOB_PHASE_RUN` の参照整合。
- `related detail requirement type`: `workflow`, `display`, `persistence`, `security`
- `adoption hint`: job setup 画面の作成完了要約、または job 詳細の監査セクションに統合しやすい。
- `conflict hint`: lifecycle 観点の作成フロー候補と重なりやすい。監査観点では「いつでも再確認できる保存断面」に限定して統合対象にする。

### CAND-TJS-002 validation failure 理由を監査用に再確認できる

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の「create job の前に validation を通せる」「validation failure の理由を確認できる」、`docs/spec.md` の「翻訳に利用する翻訳補助メタデータ､辞書､共通基盤データは､実行前､実行後ともにUIからユーザーが観測可能であること」
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-TJS-002`
- `actor`: 運用者
- `trigger`: validation failure 後に、何が不足し、どの設定が失敗理由だったかを後から追跡したい。
- `expected outcome`: 直近 validation の結果種別、失敗理由、対象入力または参照不足の区分を再確認できる。失敗理由は再試行や設定見直しに必要な粒度で残り、秘密情報や過剰な内部実装詳細は含まれない。
- `observable point`: validation 結果表示、error summary、job setup 画面の activity/history 領域、監査用の validation status 要約。
- `related detail requirement type`: `operation`, `display`, `workflow`, `security`
- `adoption hint`: failure 観点の validation 異常系候補と merge しつつ、こちらは「履歴として残す最小断面」を担当させると整理しやすい。
- `conflict hint`: failure 観点では失敗処理そのものを扱う可能性がある。operation-audit では failure 後の再観測性だけを残し、復旧操作は他 viewpoint へ譲る。

### CAND-TJS-003 AI 基盤設定の監査で secret を露出しない

- `source requirement`: `docs/spec.md` の「各フェーズのAPI選択、APIKeyは再入力不要で保存ができること」「APIKeyは暗号化して保存すること」、`docs/er.md` の「credential_ref は暗号化済み API key そのものではなく、secret store への参照だけを保持する」
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-TJS-003`
- `actor`: 運用者
- `trigger`: AI runtime 設定が保存済みかを確認したいが、監査表示や永続化で API key の実値を露出したくない。
- `expected outcome`: どの provider と model を選んだか、保存済み credential を参照しているかは確認できる。一方で API key の平文、復号可能な値、過剰な secret metadata は監査表示にも保存要約にも出ない。
- `observable point`: runtime 設定表示、保存済み credential 状態表示、永続化の `credential_ref` 利用、ログや history 断面の redaction。
- `related detail requirement type`: `security`, `display`, `persistence`
- `adoption hint`: external-integration 観点の provider 設定候補と競合しやすいため、こちらは redaction rule 専用候補として切り出すと扱いやすい。
- `conflict hint`: external-integration 観点で provider 接続確認まで広がりやすい。operation-audit では secret を保存しない監査境界に限定する。

### CAND-TJS-004 validation 通過後に設定変更があれば再検証要否を追跡できる

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の「create job の前に validation を通せる」、`docs/spec.md` の「翻訳ジョブの中断､再開､失敗回復が継続的に行えること」「翻訳ジョブ､APIの実行進捗を確認できること」、`docs/er.md` の「ジョブ状態は JOB_PHASE_RUN 群から集約する」
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-TJS-004`
- `actor`: 運用者
- `trigger`: validation 通過後に入力データや基盤参照や runtime を変更したため、作成済み pass 判定がまだ有効か確認したい。
- `expected outcome`: 現在の設定断面が最後に通過した validation 断面と一致しているかを追跡できる。不一致なら再 validation 必要であることを観測でき、古い pass を誤って監査根拠に使わない。
- `observable point`: validation 状態バッジ、最終 validation 対象の要約、変更後の dirty 状態表示、job 作成可否との整合。
- `related detail requirement type`: `workflow`, `state`, `display`, `operation`
- `adoption hint`: state-transition 観点の禁止遷移候補と関連が深い。監査側では「pass 判定の失効可視化」を独立論点として残すと良い。
- `conflict hint`: state-transition では Ready 化条件や再実行条件へ統合される可能性がある。operation-audit では pass 判定の履歴有効性だけに絞る。

### CAND-TJS-005 作成されたジョブが 1 input = 1 job の出自を失わない

- `source requirement`: `tasks/usecases/translation-job-setup.yaml` の goal と completion criteria、`docs/spec.md` の「複数の入力データを登録し、それぞれの入力データを独立した翻訳ジョブとして管理できること」「1つの翻訳ジョブは1つのxEdit抽出データを対象とし、入力ファイルの出自を失わずに保持できること」、`docs/er.md` の「TRANSLATION_JOB は 1 つの X_EDIT_EXTRACTED_DATA だけを参照する」
- `viewpoint`: `operation-audit`
- `candidate scenario id`: `CAND-TJS-005`
- `actor`: 運用者
- `trigger`: 複数入力が存在する環境で、どの入力からどのジョブが作られたかを後から監査したい。
- `expected outcome`: 各ジョブがちょうど 1 つの入力データへ対応し、入力名、出自、作成日時、関連 validation 結果を突合できる。複数入力の混線や、同一 job への複数 input 紐付けは監査上検出できる。
- `observable point`: job 一覧、job 詳細、入力データ参照、`TRANSLATION_JOB` と `X_EDIT_EXTRACTED_DATA` の関連、作成履歴要約。
- `related detail requirement type`: `persistence`, `display`, `operation`
- `adoption hint`: lifecycle 観点の作成成功候補と merge しやすいが、こちらは出自追跡可能性の監査証跡として別軸で残す価値がある。
- `conflict hint`: actor-goal や lifecycle では「作成できること」が主目的になりやすい。operation-audit では「後から混線なく説明できること」を採用判断材料として渡す。

## Open Notes

- `human decision candidate`: validation 履歴に残す失敗理由の粒度。ユーザー向け短文だけにするか、参照不足カテゴリや対象設定断面まで保持するかは人手判断が必要。
- `merge candidate`: `CAND-TJS-001` と `CAND-TJS-005` は job 作成監査の統合候補である。`CAND-TJS-002` と `CAND-TJS-004` は validation 監査として統合候補である。
- `rejection candidate`: 監査形式を特定のログ実装や attempt 履歴テーブル前提で固定する案。`docs/er.md` は attempt 履歴テーブルを持たないため、保存形式の固定は候補から外す余地がある。
