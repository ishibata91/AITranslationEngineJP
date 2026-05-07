# Scenario Candidates: 2026-05-07-model-settings-card-controller / failure

- `generator`: `failure`
- `source_plan`: `./task-frame.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `MSCF`

## Generator Scope

- `viewpoint`: 失敗
- `included_sources`: `./task-frame.md`, `./light-change-planning.md`, `../../../detail-specs/translation-job-setup.md`, `../../../detail-specs/ai-provider-settings-management.md`, `../../../architecture.md`
- `excluded_sources`: プロダクトコード変更、プロダクトテスト変更、docs 正本本文変更、採否判断、統合判断、最終シナリオ表
- `generation_notes`: モデル設定カード制御の provider、model、model list、保存、取得、選択状態について、失敗入力、参照不能、設定不整合、保存失敗、回復動作だけを候補化する。

## Candidate Scenarios

### CAND-MSCF-001 APIキー未設定の provider で model list 更新が拒否される

- `source requirement`: [task-frame.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-settings-card-controller/task-frame.md:12), [translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md:40), [translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md:62)
- `viewpoint`: failure / 失敗入力
- `candidate scenario id`: `CAND-MSCF-001`
- `actor`: 利用者
- `failure start condition`: APIキーが必要な provider を選び、provider settings に APIキー状態がない。
- `rejected operation`: model list 更新を実行しようとする。
- `expected error`: APIキー未設定で更新不可であることを短い日本語文で表示する。
- `expected outcome`: 外部 model list 取得は起動されず、選択中の model は保存済み状態または未選択状態から勝手に変わらない。
- `observable point`: 更新ボタンの無効状態、APIキー未設定表示、外部取得未実行、model 選択 UI の非表示または操作不能状態。
- `related detail requirement type`: 入力不備、参照不能、UI 状態表示
- `adoption hint`: 採否は designer が判断する。共有カード制御の APIキー gate と Job Setup の既存 gate を同じ期待にできる候補である。
- `conflict hint`: LM Studio は APIキーを要求しないため、この候補を全 provider 共通 gate にすると仕様と衝突する。

### CAND-MSCF-002 APIキー不要 provider の参照不能を APIキー不足に誤分類しない

- `source requirement`: [translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md:41), [translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md:62), [ai-provider-settings-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/ai-provider-settings-management.md:26), [ai-provider-settings-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/ai-provider-settings-management.md:61)
- `viewpoint`: failure / 設定不整合
- `candidate scenario id`: `CAND-MSCF-002`
- `actor`: 利用者
- `failure start condition`: LM Studio を選び、endpoint 未設定または endpoint 参照不能の状態である。
- `rejected operation`: 参照不能の状態で model list 更新または model 保存を実行しようとする。
- `expected error`: endpoint 参照不能または model list 取得失敗として分類する。APIキー未設定エラーは表示しない。
- `expected outcome`: APIキー不足 gate は発火せず、APIキー入力、APIキー未設定 warning、credential select は表示されない。
- `observable point`: LM Studio 選択時の警告種別、credential UI 非表示、model list 更新可否、保存前エラー分類。
- `related detail requirement type`: 設定不整合、参照不能、UI 状態表示
- `adoption hint`: 採否は designer が判断する。provider 種別ごとの失敗分類を固定する候補である。
- `conflict hint`: provider settings の endpoint 未設定を同時に扱う場合、APIキー不足ではなく endpoint 参照不能として分ける必要がある。

### CAND-MSCF-003 user-facing provider list へ fake provider が混入した状態を拒否する

- `source requirement`: [task-frame.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-settings-card-controller/task-frame.md:13), [task-frame.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-settings-card-controller/task-frame.md:14), [ai-provider-settings-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/ai-provider-settings-management.md:26)
- `viewpoint`: failure / 失敗入力
- `candidate scenario id`: `CAND-MSCF-003`
- `actor`: 利用者
- `failure start condition`: provider list または保存済み表示状態に `fake` provider ID が混入している。
- `rejected operation`: `fake` provider ID を利用者が選択または保存する。
- `expected error`: provider が利用者向け選択肢ではないことを表示する。fake mode の内部詳細は表示しない。
- `expected outcome`: 利用者向け provider list は `gemini`、`lm_studio`、`xai` だけを表示し、保存 payload に `fake` provider ID を入れない。
- `observable point`: provider list、選択値、保存 payload、エラー文、frontend に fake mode 判定や `fake-model` 固有分岐がないこと。
- `related detail requirement type`: 失敗入力、設定不整合、信頼境界
- `adoption hint`: 採否は designer が判断する。fake mode を通常 provider ID の model list 応答として扱う前提を守る候補である。
- `conflict hint`: 既存保存値に `fake` provider ID が残る可能性がある場合、移行時の扱いは人間判断候補である。

### CAND-MSCF-004 model list 取得失敗時に raw data と secret を露出しない

- `source requirement`: [translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md:63), [translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md:64), [ai-provider-settings-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/ai-provider-settings-management.md:29), [ai-provider-settings-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/ai-provider-settings-management.md:61)
- `viewpoint`: failure / 参照不能
- `candidate scenario id`: `CAND-MSCF-004`
- `actor`: 利用者
- `failure start condition`: provider settings は参照できるが、model list API が失敗応答、timeout、不正応答のいずれかを返す。
- `rejected operation`: 取得失敗状態のまま model を選択または保存する。
- `expected error`: model list 取得失敗を分類と短い要約で表示する。APIキー、raw request、raw response、raw prompt は表示しない。
- `expected outcome`: model 選択 UI は表示しない、または操作不能にする。既存の選択済み model を上書きしない。
- `observable point`: 取得失敗表示、secret 非表示、raw payload 非表示、model 選択可否、保存拒否状態。
- `related detail requirement type`: 参照不能、失敗表示、情報非露出
- `adoption hint`: 採否は designer が判断する。外部応答失敗と UI 表示の境界を固定する候補である。
- `conflict hint`: provider settings 画面の接続確認失敗表示と、モデル設定カードの model list 失敗表示の文言粒度が競合しうる。

### CAND-MSCF-005 空の model list 取得成功を model 未選択として扱う

- `source requirement`: [translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md:39), [translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md:43), [translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md:63)
- `viewpoint`: failure / 参照不能
- `candidate scenario id`: `CAND-MSCF-005`
- `actor`: 利用者
- `failure start condition`: model list API は成功扱いだが、選択可能な model が 0 件である。
- `rejected operation`: model 未選択のまま保存または job 作成を実行する。
- `expected error`: model 候補がないため model を選択できないことを表示する。
- `expected outcome`: model list 取得失敗とは別の状態として扱うか、取得失敗に含めるかを designer が判断できる候補として残す。
- `observable point`: model list 状態、model 選択 UI、保存可否、job 作成可否、表示文言。
- `related detail requirement type`: 参照不能、入力不備、UI 状態表示
- `adoption hint`: 採否は designer が判断する。空配列応答を failure として分類するかどうかの判断候補である。
- `conflict hint`: 正本は「取得成功した場合だけ model 選択を表示」と書くが、空配列成功時の表示状態は未確定である。

### CAND-MSCF-006 provider 変更後の遅延 model list 応答を破棄する

- `source requirement`: [task-frame.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-settings-card-controller/task-frame.md:20), [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md:112), [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md:114)
- `viewpoint`: failure / 回復動作
- `candidate scenario id`: `CAND-MSCF-006`
- `actor`: 利用者
- `failure start condition`: provider A の model list 更新中に provider B へ変更し、その後 provider A の応答が遅れて返る。
- `rejected operation`: 遅れて返った provider A の model list を provider B の選択状態へ適用する。
- `expected error`: 原則として利用者向けエラーは出さず、古い応答を破棄する。必要な場合だけ更新済み provider の状態を表示する。
- `expected outcome`: provider B の model list 状態と model 選択状態だけが画面状態へ残る。
- `observable point`: request token または provider snapshot による破棄、Store の provider/model/list 整合、古い応答で UI が戻らないこと。
- `related detail requirement type`: 回復動作、状態整合性、遅延応答破棄
- `adoption hint`: 採否は designer が判断する。共有 controller / usecase / store に必要な遅延応答破棄の候補である。
- `conflict hint`: 破棄時にユーザーへ通知するかどうかは、正常操作体験との競合候補である。

### CAND-MSCF-007 provider 変更後に旧 provider の model 保存を拒否する

- `source requirement`: [task-frame.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-settings-card-controller/task-frame.md:20), [translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md:38), [translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md:43)
- `viewpoint`: failure / 設定不整合
- `candidate scenario id`: `CAND-MSCF-007`
- `actor`: 利用者
- `failure start condition`: provider を変更した後、旧 provider の model が選択状態に残っている。
- `rejected operation`: 旧 provider の model を現在 provider の model として保存する、または job 作成に使う。
- `expected error`: 現在 provider で model 未選択であることを表示する。
- `expected outcome`: provider 変更時に model 選択と model list 状態を現在 provider と整合する状態へ戻す。
- `observable point`: provider 変更直後の model 表示、保存 payload、job 作成可否、設定済み判定。
- `related detail requirement type`: 設定不整合、入力不備、状態整合性
- `adoption hint`: 採否は designer が判断する。provider/model/list をカード側へ集約する時の基本不変条件候補である。
- `conflict hint`: 保存済み model を provider ごとに保持して復元する設計を採る場合、単純な未選択化と競合しうる。

### CAND-MSCF-008 取得失敗時に画面ごとの fallback へ逃げない

- `source requirement`: [light-change-planning.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-settings-card-controller/light-change-planning.md:26), [light-change-planning.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-settings-card-controller/light-change-planning.md:27), [ai-provider-settings-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/ai-provider-settings-management.md:34)
- `viewpoint`: failure / 参照不能
- `candidate scenario id`: `CAND-MSCF-008`
- `actor`: 利用者
- `failure start condition`: 共有 model 設定取得口または provider settings 参照が失敗する。
- `rejected operation`: マスターペルソナまたは Job Setup の個別 fallback 値を使って保存済み設定として扱う。
- `expected error`: 設定取得または provider settings 参照不能として表示する。
- `expected outcome`: Job Setup と master-persona は provider settings を参照し、個別 secret や endpoint に fallback しない。
- `observable point`: 取得失敗表示、fallback 値の不使用、保存済み扱いにならないこと、raw payload 非表示。
- `related detail requirement type`: 参照不能、設定不整合、責務境界
- `adoption hint`: 採否は designer が判断する。共有 controller と画面別 usecase の責務境界を固定する候補である。
- `conflict hint`: 既存画面が持つ初期値と共有取得失敗時の表示規約が衝突しうる。

### CAND-MSCF-009 保存失敗時に未保存の選択状態を確定扱いにしない

- `source requirement`: [task-frame.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-settings-card-controller/task-frame.md:5), [task-frame.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-settings-card-controller/task-frame.md:20), [ai-provider-settings-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/ai-provider-settings-management.md:61)
- `viewpoint`: failure / 保存失敗
- `candidate scenario id`: `CAND-MSCF-009`
- `actor`: 利用者
- `failure start condition`: provider と model を選んだ後、保存処理が失敗する。
- `rejected operation`: 保存失敗後の未保存状態を保存済み設定として表示または利用する。
- `expected error`: 保存失敗を分類と短い要約で表示する。APIキー、raw request、raw response、raw prompt は表示しない。
- `expected outcome`: 画面は未保存変更を区別し、再試行または再読込で回復できる状態を維持する。
- `observable point`: 保存済み表示、未保存表示、保存 payload、再試行操作、secret と raw payload の非表示。
- `related detail requirement type`: 保存失敗、回復動作、情報非露出
- `adoption hint`: 採否は designer が判断する。保存失敗後の UI 状態と回復操作を固定する候補である。
- `conflict hint`: 楽観更新で表示を先に変える設計を採る場合、rollback 表示または dirty state 表示の判断が必要である。

### CAND-MSCF-010 Job Setup で model 未選択または更新中のまま job 作成を拒否する

- `source requirement`: [translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md:43), [translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md:64), [translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md:65)
- `viewpoint`: failure / 失敗入力
- `candidate scenario id`: `CAND-MSCF-010`
- `actor`: 利用者
- `failure start condition`: 3 つの翻訳段階のいずれかで model 未選択、model list 更新中、model list 取得失敗の状態が残っている。
- `rejected operation`: 翻訳 job を作成する。
- `expected error`: 対象翻訳段階の model 未選択、更新中、取得失敗を分けて表示する。
- `expected outcome`: Ready job は作成されず、作成前確認は失敗理由を表示する。
- `observable point`: job 作成ボタンの可否、段階別状態表示、作成要求の未送信、エラー文の折り返し。
- `related detail requirement type`: 入力不備、UI 状態表示、状態整合性
- `adoption hint`: 採否は designer が判断する。Job Setup 既存仕様と共有カード制御の失敗状態を接続する候補である。
- `conflict hint`: マスターペルソナ側には job 作成操作がないため、共有候補として扱う範囲を Job Setup 限定にする必要がある。

### CAND-MSCF-011 共有カード制御が master-persona 設定を Job Setup の既定値にしない

- `source requirement`: [task-frame.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/active/2026-05-07-model-settings-card-controller/task-frame.md:18), [translation-job-setup.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/translation-job-setup.md:14), [ai-provider-settings-management.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/detail-specs/ai-provider-settings-management.md:34)
- `viewpoint`: failure / 設定不整合
- `candidate scenario id`: `CAND-MSCF-011`
- `actor`: 利用者
- `failure start condition`: master-persona で保存済み model 設定があり、Job Setup を開く。
- `rejected operation`: master-persona の保存済み設定を Job Setup の既定値または保存元として適用する。
- `expected error`: 原則としてエラーではなく、Job Setup 側の phase runtime settings 未設定または保存済み状態を表示する。
- `expected outcome`: Job Setup は master-persona の AI 設定を既定値または保存元として扱わない。
- `observable point`: Job Setup 初期表示、phase runtime settings、保存 payload、master-persona 設定の非混入。
- `related detail requirement type`: 設定不整合、責務境界、状態整合性
- `adoption hint`: 採否は designer が判断する。共有カード部品と共有状態制御を混同しないための候補である。
- `conflict hint`: 共有 controller / store を単一 namespace にすると、画面横断の状態混入が起きる可能性がある。

## Open Notes

- `human decision candidate`: 空の model list 成功を「取得済み 0 件」と表示するか、「取得失敗」に寄せるかは未確定である。
- `human decision candidate`: 保存失敗後に、選択値を rollback するか、未保存変更として残すかは未確定である。
- `human decision candidate`: 既存保存値に `fake` provider ID がある場合の移行表示と拒否文言は未確定である。
- `human decision candidate`: 遅延応答破棄を利用者に通知するか、無音で破棄するかは未確定である。
- `human decision candidate`: 共有 controller / store の状態 namespace を画面単位、用途単位、カード instance 単位のどれにするかは未確定である。
- `merge candidate`: CAND-MSCF-001 と CAND-MSCF-002 は provider 種別別の APIキー gate として統合される可能性がある。
- `merge candidate`: CAND-MSCF-006 と CAND-MSCF-007 は provider 変更時の状態整合性として統合される可能性がある。
- `merge candidate`: CAND-MSCF-004、CAND-MSCF-005、CAND-MSCF-010 は model list 失敗状態と job 作成拒否の検証として統合される可能性がある。
- `rejection candidate`: 正常系の裏返しだけで、失敗開始条件、拒否操作、期待エラー、観測点を持たない候補は採否判断前に除外対象になりうる。
