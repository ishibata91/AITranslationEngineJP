# Scenario Candidates: ai-provider-settings-management / failure

- `generator`: `failure`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `AIPSM`
- `candidate_count`: 9

## Generator Scope

- `viewpoint`: 異常系。失敗入力、参照不能、設定不整合、保存失敗、回復動作を候補化する。
- `included_sources`:
  - `./plan.md`
  - `../../../spec.md`
  - `../../completed/translation-job-setup-phase-provider-settings/scenario-design.md`
  - `../../completed/2026-04-16-master-persona-gap-closure.implementation-scope.md`
- `excluded_sources`: product code、product test、docs 正本、implementation-scope、採否決定、統合判断。
- `generation_notes`: 採否、統合、競合解消は `designer` が行う。候補は失敗開始条件、拒否される操作、期待エラー、観測点に限定する。

## Candidate Scenarios

### CAND-AIPSM-F001 APIキー保存失敗時に設定完了扱いにしない

- `source requirement`: `plan.md:37` APIキーとエンドポイントは翻訳フェーズやマスターペルソナ生成とは別の永続仕様で管理する。`plan.md:38` APIキー平文は UI、DTO、log、エラー要約へ出さない。`docs/spec.md:57-58` APIKey は再入力不要で保存し、暗号化保存する。
- `viewpoint`: 保存失敗 / 回復動作
- `candidate scenario id`: `CAND-AIPSM-F001`
- `actor`: プロバイダ設定を保存するユーザー。
- `trigger`: Gemini または xAI の APIキーを入力し、保存時に secret store への書き込みが失敗する。
- `rejected operation`: APIキー保存成功として provider 設定を complete にする操作を拒否する。
- `expected error`: APIキー保存に失敗したことを示す。エラー文、通知、保存要約に APIキー平文を含めない。
- `expected outcome`: provider 設定は「APIキー未保存」または「保存失敗」の状態として残る。後続の model list 取得、provider validation、翻訳フェーズ実行は secret 未解決として扱う。
- `observable point`: 設定画面の保存状態、secret store spy、settings read DTO、structured log、error summary。
- `related detail requirement type`: `failure_handling_requirement`, `security_requirement`, `data_requirement`, `recovery_requirement`
- `acceptance viewpoint`: secret 書き込み失敗後に、APIキー平文が観測面へ出ず、再保存可能な状態が表示される。
- `adoption hint`: provider settings 保存シナリオの異常系として採用候補。
- `conflict hint`: endpoint/model/batch の非 secret 設定だけを保存するか、全体を失敗にするかは設計判断が必要になる可能性がある。

### CAND-AIPSM-F002 不正 endpoint を保存しない

- `source requirement`: `plan.md:8` プロバイダ別にエンドポイントと APIキーを設定できるようにする。`plan.md:86` provider settings の保存単位と endpoint 更新を designer が含める。
- `viewpoint`: 失敗入力
- `candidate scenario id`: `CAND-AIPSM-F002`
- `actor`: endpoint を入力するユーザー。
- `trigger`: endpoint 欄へ空文字、URL として解釈できない値、provider の許可範囲外の scheme、末尾制御文字を含む値を入力して保存する。
- `rejected operation`: 不正 endpoint を provider 設定として保存する操作を拒否する。
- `expected error`: endpoint が不正であることを provider 別に表示する。APIキー、raw request、raw response は表示しない。
- `expected outcome`: 既存の有効 endpoint は上書きされない。未設定 provider の場合は endpoint 未設定のまま残る。
- `observable point`: endpoint input validation、save response、settings read DTO、DB write spy、error summary。
- `related detail requirement type`: `boundary_requirement`, `failure_handling_requirement`, `data_requirement`, `security_requirement`
- `acceptance viewpoint`: 不正 endpoint では保存完了にならず、既存の有効設定が破壊されない。
- `adoption hint`: UI と API の両方で観測できる入力異常候補。
- `conflict hint`: LM Studio の `http://localhost` を許可するか、Gemini / xAI に HTTPS を要求するかは provider 別 validation 仕様と統合が必要である。

### CAND-AIPSM-F003 provider 設定DB保存失敗時に途中状態を完了表示しない

- `source requirement`: `plan.md:39` 各プロバイダ設定では利用モデルとバッチ API 利用可否だけを provider 別の実行設定として変更できる。`plan.md:40` DB 変更が必要かは repository、migration、secret store の責務を分けて確定する。
- `viewpoint`: 保存失敗 / データ整合性
- `candidate scenario id`: `CAND-AIPSM-F003`
- `actor`: endpoint、model、batch API 切り替えを保存するユーザー。
- `trigger`: secret 保存は成功したが、endpoint、model、batch API 設定の DB 書き込みが失敗する。
- `rejected operation`: DB 保存失敗後に provider 設定全体を保存済みとして表示する操作を拒否する。
- `expected error`: provider 設定の保存に失敗したことを示す。secret 値や DB 内部エラー詳細は出さない。
- `expected outcome`: 読み直し後の provider 設定は、最後に正常保存された状態と一致する。新しい設定値は「保存済み」として扱わない。
- `observable point`: save response、settings reload result、DB transaction spy、secret store state、redacted log。
- `related detail requirement type`: `data_requirement`, `consistency_requirement`, `failure_handling_requirement`, `recovery_requirement`
- `acceptance viewpoint`: DB 保存失敗時に、画面表示、読み取り結果、後続参照の設定断面が食い違わない。
- `adoption hint`: DB 変更候補の異常系として採用候補。
- `conflict hint`: secret 保存成功後に DB 保存が失敗した場合の secret 削除または未参照化は、人間判断が必要になる可能性がある。

### CAND-AIPSM-F004 DB migration 失敗時に旧 schema で設定を書かない

- `source requirement`: `plan.md:40` DB 変更が必要かは scenario-design と implementation-scope で repository、migration、secret store の責務を分けて確定する。`completed implementation-scope.md:39-51` は SQLite concrete、migration、repository tests の分担を示す。
- `viewpoint`: 参照不能 / 保存失敗 / 回復動作
- `candidate scenario id`: `CAND-AIPSM-F004`
- `actor`: アプリを起動して provider settings を開くユーザー。
- `trigger`: provider settings 用 DB migration が失敗する。例として、新 table 作成失敗、既存列追加失敗、schema version 更新失敗がある。
- `rejected operation`: 旧 schema に対して provider settings の保存または読み取りを継続する操作を拒否する。
- `expected error`: provider settings の保存領域が利用できないことを示す。secret 値、raw SQL、内部 stack trace は表示しない。
- `expected outcome`: provider settings route は利用不可または保存不可として扱う。翻訳フェーズや master-persona の旧設定を fallback 保存先にしない。
- `observable point`: startup migration result、settings route state、repository open result、DB schema version、redacted log。
- `related detail requirement type`: `failure_handling_requirement`, `data_requirement`, `compatibility_requirement`, `recovery_requirement`
- `acceptance viewpoint`: migration 失敗時に旧 schema へ silent write せず、利用者が復旧不能な誤保存状態に入らない。
- `adoption hint`: DB migration candidate の必須異常系として採用候補。
- `conflict hint`: migration 失敗時に app 全体を止めるか、provider settings だけを止めるかは設計判断が必要である。

### CAND-AIPSM-F005 secret 非露出を保存失敗と検証失敗でも守る

- `source requirement`: `plan.md:38` APIキーは UI、DTO、log、エラー要約へ平文表示しない。`completed scenario-design.md:30` は API key 平文、secret 本体、raw request / response、raw prompt を UI、DTO、error summary、structured log、fake transport log、保存要約へ出さないと固定している。
- `viewpoint`: 設定不整合 / セキュリティ失敗
- `candidate scenario id`: `CAND-AIPSM-F005`
- `actor`: APIキーを入力し、保存または接続検証を実行するユーザー。
- `trigger`: APIキーが不正、期限切れ、secret store 参照不能、provider validation 失敗のいずれかになる。
- `rejected operation`: APIキー平文、復号可能値、raw request、raw response を表示またはログ保存する操作を拒否する。
- `expected error`: 失敗理由は分類名と provider 名だけで示す。credential の存在状態または参照状態だけを表示する。
- `expected outcome`: UI、DTO、保存要約、error summary、structured log、fake transport log に secret 本体が出ない。
- `observable point`: UI rendering、Wails DTO、gateway response、error summary、structured log、fake transport log、redaction assertion。
- `related detail requirement type`: `security_requirement`, `failure_handling_requirement`, `observability_requirement`, `testability_requirement`
- `acceptance viewpoint`: 失敗時ほど secret が露出しやすいため、全観測面で redaction を確認する。
- `adoption hint`: trust boundary の hard gate 候補として採用候補。
- `conflict hint`: 監査ログにどの失敗分類まで残すかは operation-audit 観点と統合が必要である。

### CAND-AIPSM-F006 旧 Job Setup / master-persona 設定を provider settings の代替にしない

- `source requirement`: `plan.md:8` APIキーとエンドポイントの永続仕様は翻訳フェーズやマスターペルソナ生成とは別に管理する。`plan.md:9` 翻訳ジョブ設定やマスターペルソナ生成が個別に secret や endpoint を持たない構造へ寄せる。`completed scenario-design.md:20` は Job Setup が master-persona provider 設定を参照しないと固定している。
- `viewpoint`: 参照不能 / 設定不整合
- `candidate scenario id`: `CAND-AIPSM-F006`
- `actor`: provider settings 未設定の状態で既存 Job Setup または master-persona 設定が残るユーザー。
- `trigger`: provider settings は未保存だが、旧 Job Setup phase 設定または master-persona provider 設定に provider、model、credential 参照が残っている。
- `rejected operation`: 旧設定を provider settings の既定値、保存元、secret 解決元として使う操作を拒否する。
- `expected error`: provider settings が未設定であることを示す。旧設定の secret 参照を使って成功扱いにしない。
- `expected outcome`: provider settings 画面では未設定として表示する。翻訳フェーズや master-persona の旧設定は provider settings の保存完了条件を満たさない。
- `observable point`: settings load result、setting source summary、secret namespace、route state、validation summary。
- `related detail requirement type`: `compatibility_requirement`, `consistency_requirement`, `security_requirement`, `failure_handling_requirement`
- `acceptance viewpoint`: 永続仕様分離により、旧設定参照で新設定が成立したように見えない。
- `adoption hint`: 既存 Job Setup / master-persona との境界シナリオとして採用候補。
- `conflict hint`: 旧設定を明示 migration で取り込むか、手動再設定にするかは designer の人間判断候補になる可能性がある。

### CAND-AIPSM-F007 provider capability と保存値が矛盾する場合に実行設定へ混入させない

- `source requirement`: `plan.md:39` 各プロバイダ設定ではモデルとバッチ API 切り替えだけを設定できる。`docs/spec.md:49-51` LMStudio、Gemini、xAI を扱い、Gemini と xAI は BatchAPI を利用できる。`completed implementation-scope.md:13-15` は real provider id を `gemini`、`lm_studio`、`xai` のみにし、fake は test-only DI とする。
- `viewpoint`: 設定不整合 / 失敗入力
- `candidate scenario id`: `CAND-AIPSM-F007`
- `actor`: provider settings を保存または読み込むユーザー。
- `trigger`: `lm_studio` に batch API 有効、未知 provider id、fake provider id、provider capability に存在しない model / batch 組み合わせが保存入力または既存DBにある。
- `rejected operation`: capability と矛盾する provider 設定を有効な実行設定として保存または表示する操作を拒否する。
- `expected error`: provider capability 不整合として表示する。fake provider はユーザー向け provider list に出さない。
- `expected outcome`: 矛盾値は provider settings の complete 条件を満たさない。後続の翻訳フェーズや master-persona は矛盾値を参照しない。
- `observable point`: provider list、capability fixture、settings validation summary、settings save response、read DTO。
- `related detail requirement type`: `consistency_requirement`, `boundary_requirement`, `security_requirement`, `compatibility_requirement`
- `acceptance viewpoint`: provider capability と保存値が矛盾する場合に、stale 値や fake provider が実行経路へ入らない。
- `adoption hint`: provider capability 不整合の中心候補として採用候補。
- `conflict hint`: 不整合値を保存時に拒否するか、読み込み時に無効化して再保存を促すかは状態遷移観点と統合が必要である。

### CAND-AIPSM-F008 endpoint 参照不能時に APIキー再入力や provider fallback を要求しない

- `source requirement`: `docs/spec.md:51-52` Gemini と xAI の BatchAPI 失敗時に別 provider fallback は不要。`plan.md:86` endpoint 更新と既存 Job Setup / master-persona からの参照境界を designer が含める。`completed scenario-design.md:74` は endpoint 不通や model list failure を API key missing とは別カテゴリで表示すると固定している。
- `viewpoint`: 参照不能 / 回復動作
- `candidate scenario id`: `CAND-AIPSM-F008`
- `actor`: endpoint を保存済みの provider で model list または接続検証を実行するユーザー。
- `trigger`: endpoint が DNS 失敗、接続拒否、timeout、HTTP 失敗、provider 形式不一致のいずれかになる。
- `rejected operation`: endpoint 参照不能を APIキー未設定として扱う操作、別 provider へ自動 fallback する操作、APIキー再入力を必須にする操作を拒否する。
- `expected error`: endpoint 参照不能として表示する。APIキー状態とは別カテゴリにする。
- `expected outcome`: 保存済み APIキー参照は保持される。利用者は endpoint の修正または再検証へ進める。
- `observable point`: validation summary、request spy、endpoint error category、secret reference state、settings route state。
- `related detail requirement type`: `failure_handling_requirement`, `recovery_requirement`, `security_requirement`, `testability_requirement`
- `acceptance viewpoint`: 接続失敗時に secret を再入力させず、失敗原因を endpoint 側として観測できる。
- `adoption hint`: endpoint 不正とは別の参照不能候補として採用候補。
- `conflict hint`: 外部連携観点の provider validation failure と統合される可能性がある。

### CAND-AIPSM-F009 endpoint 更新後に旧 model / 旧 capability を再利用しない

- `source requirement`: `plan.md:8` provider 別に endpoint、APIキー、モデル、バッチ API 切り替えを設定する。`completed scenario-design.md:88` は遅延 model list 結果を現在の provider / APIキー状態へ混入させないと固定している。
- `viewpoint`: 設定不整合 / 旧設定参照
- `candidate scenario id`: `CAND-AIPSM-F009`
- `actor`: endpoint を変更するユーザー。
- `trigger`: endpoint 変更前に選択した model が、変更後 endpoint の model list または provider capability と一致しない。
- `rejected operation`: endpoint 変更前の model、batch capability、遅延 model list 結果を変更後 endpoint の設定として再利用する操作を拒否する。
- `expected error`: 現在 endpoint で model 未選択または model 再取得待ちであることを示す。
- `expected outcome`: endpoint 変更後は model 選択が未確定になり、現在 endpoint の model list 成功まで provider settings の complete 条件を満たさない。
- `observable point`: endpoint field、model selector、model list source endpoint、save button state、settings read DTO。
- `related detail requirement type`: `state_requirement`, `consistency_requirement`, `failure_handling_requirement`, `compatibility_requirement`
- `acceptance viewpoint`: endpoint 更新で旧 model が残り、後続実行が存在しない model を参照する失敗を防ぐ。
- `adoption hint`: 旧設定参照と遅延結果混入の provider settings 版として採用候補。
- `conflict hint`: state-transition 観点の endpoint 更新シナリオと統合される可能性がある。

## Open Notes

- `human decision candidate`: `Q-AIPSM-F001` secret 保存成功後に DB 保存が失敗した場合、secret を削除するか、未参照 secret として残して再試行へ使うか。
- `human decision candidate`: `Q-AIPSM-F002` provider 別 endpoint validation の厳しさ。例として、LM Studio の local HTTP 許可、Gemini / xAI の HTTPS 必須、空 endpoint の既定扱いがある。
- `human decision candidate`: `Q-AIPSM-F003` DB migration 失敗時に app 起動全体を止めるか、provider settings 画面だけを利用不可にするか。
- `merge candidate`: `CAND-AIPSM-F008` は external-integration 観点の接続検証失敗候補と統合される可能性がある。
- `merge candidate`: `CAND-AIPSM-F009` は state-transition 観点の endpoint 更新候補と統合される可能性がある。
- `rejection candidate`: exact table 名、migration 実装方式、keyring backend の具体エラー文、retry 実装方式は本 generator では確定しない。
