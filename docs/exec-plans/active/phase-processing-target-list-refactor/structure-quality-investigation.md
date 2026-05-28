# 構造品質調査: phase-processing-target-list-refactor

## 調査結果

- 判断結果: 完了
- 調査 mode: `構造品質調査`
- 引き継ぎ先: `designer`
- 対象範囲: 単語翻訳、NPC ペルソナ生成、本文翻訳の `処理対象一覧` と、その件数、検索、read model 境界

## 根拠参照

### plan / 既存調査

- `docs/exec-plans/active/phase-processing-target-list-refactor/plan.md`
- `docs/exec-plans/active/phase-processing-target-list-refactor/spec-drift-investigation.md`

### frontend

- `frontend/src/ui/components/ProcessingTargetListPanel.svelte`
- `frontend/src/ui/components/ProcessingTargetListWrapper.svelte`
- `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte`
- `frontend/src/ui/screens/persona-generation-phase/PersonaGenerationPhasePanel.svelte`
- `frontend/src/ui/screens/body-translation-phase/BodyTranslationPhasePanel.svelte`
- `frontend/src/application/presenter/term-translation-phase/term-translation-phase.presenter.ts`
- `frontend/src/application/presenter/persona-generation-phase/persona-generation-phase.presenter.ts`
- `frontend/src/application/presenter/body-translation-phase/body-translation-phase.presenter.ts`
- `frontend/src/application/usecase/term-translation-phase/term-translation-phase.usecase.ts`
- `frontend/src/application/usecase/persona-generation-phase/persona-generation-phase.usecase.ts`
- `frontend/src/application/usecase/body-translation-phase/body-translation-phase.usecase.ts`
- `frontend/src/application/store/term-translation-phase/term-translation-phase.store.ts`
- `frontend/src/application/store/persona-generation-phase/persona-generation-phase.store.ts`
- `frontend/src/application/store/body-translation-phase/body-translation-phase.store.ts`
- `frontend/src/application/gateway-contract/processing-target/processing-target-gateway-contract.ts`
- `frontend/src/ui/screens/job-run/JobRunPage.svelte`

### backend

- `internal/service/term_translation_phase_service.go`
- `internal/service/persona_generation_phase_service.go`
- `internal/service/body_translation_phase_service.go`
- `internal/service/body_translation_input_snapshot.go`
- `internal/service/processing_target_read_model_service.go`
- `internal/repository/processing_target_sqlite_repository.go`
- `internal/controller/wails/processing_target_controller.go`
- `internal/controller/wails/term_translation_phase_controller.go`
- `internal/controller/wails/persona_generation_phase_controller.go`
- `internal/controller/wails/body_translation_phase_controller.go`
- `internal/usecase/processing_target_read_model_usecase.go`
- `internal/usecase/persona_generation_phase_contract_test.go`

## 観測事実

- `ProcessingTargetListPanel.svelte` は `totalCount` を優先してページ範囲と件数を表示する。件数の意味付けは受け取らない。表示責務に留まる。
- `ProcessingTargetListWrapper.svelte` は `countText` を表示できるが、3 フェーズ画面はいずれも `countText` を渡していない。一覧ヘッダ側では主要件数の意味を補っていない。
- 3 フェーズの usecase は、`setJobId` で空の `processingTargetPageState` を先に入れ、その後 `fetchSummaryAndReadiness` で summary と一覧を並列取得する。単語翻訳と NPC ペルソナ生成は単一 page state、本文翻訳は `processingTargetPageStatesByPhase` を持つ。
- `ProcessingTargetReadModelService` は page、pageSize、phase 名の正規化だけを行い、件数意味や母集団の整合は扱わない。
- `ProcessingTargetController` と frontend `ProcessingTargetListResponse` は、top-level `metadata` 配列を持つが、controller は常に空配列を返す。UI は item ごとの metadata を使うが、top-level `metadata` は clone されるだけで使われていない。
- `JobRunPage.svelte` は本文翻訳だけ `body_translation` と `translation_complete` の page state を phase 別に持つ。`translation_complete` は今回の必要観点 3 件の対象外である。

## フェーズ別の構造品質

### 単語翻訳

#### observed

- summary 側の母集団は `TermTranslationPhaseService.ReadSummary` で `len(candidates)`、`len(snapshotHits)`、差分の `aiTargetCount` から作られる。
- 一覧側の母集団は `processing_target_sqlite_repository.go` の `processingTargetTermCountSQL` / `processingTargetTermListSQL` で `DICTIONARY_ENTRY` の `dictionary_lifecycle = 'job'` を数える。
- presenter は進行状況詳細に `progress.totalCount` と `progress.aiTargetCount` を出し、画面メトリクスは `summary.totalTermCount` を `対象` として出す。
- panel は `summaryProcessingTargetItems` で `AI 訳語候補` を補助表示できるが、実際の表示は `processingTargetPageState` がある限り backend 一覧を優先する。
- usecase は一覧取得 phase を固定で `term_translation` にしており、summary と一覧の意味合わせを frontend で補正しない。

#### 構造設計不整合

- summary 境界は「原語候補全体」と「AI 送信対象」を返すが、read model 境界は「翻訳ジョブ内辞書」を返す。phase summary と processing target read model の主語が一致していない。
- view は `ProcessingTargetListPanel` に `totalCount` だけを渡すため、`0 件` でも progress 側の `4930 件` と並び得る。件数の意味が DTO 境界で保持されていない。

#### 責務分離不足

- 単語翻訳画面は phase metrics、progress detail、summary fallback、backend page state の 4 か所に件数主語を分散している。件数表示の正本が 1 か所に固定されていない。
- read model が job dictionary 一覧を返す一方、summary は candidate 全体を返すため、どちらを処理対象一覧の正本にするかが backend 境界で未分離のまま frontend に流れている。

#### 変更候補

- summary と一覧で同じ母集団を参照するように、単語翻訳専用の processing target read model 主語を固定する候補。
- その主語に合わせて、term phase の presenter / panel が参照する主要件数ラベルを 1 系統へ寄せる候補。
- summary fallback と backend page state の優先順位は維持しつつ、fallback が示す対象と backend 一覧の対象を同じ意味へそろえる候補。

#### 変更不要範囲

- `ProcessingTargetListPanel.svelte` のページング計算自体は、受け取った `totalCount` を正しく表示している。今回の不一致原因ではない。
- 単語翻訳 usecase の検索イベント配線、ページ移動配線、busy 更新は一覧取得責務として妥当である。

### NPC ペルソナ生成

#### observed

- `PersonaGenerationPhaseService` は `Progress.TargetCount` と `TargetSummary.TargetCount` に同じ `snapshot.targetCount` を入れる。
- `persona_generation_phase_contract_test.go` は `TargetSummary.TargetCount == Progress.TotalCount` を検証している。
- 一覧 read model は `PERSONA` と `NPC_PROFILE` を join して返し、検索対象は `display_name`、`form_id`、`editor_id`、`record_type`、`race`、`sex`、`npc_class`、`voice_type` である。
- presenter は主要件数に `summary.progress.targetCount` を使い、panel の summary fallback では `targetSummary.targetCount` を `NPC 件数` として表示する。
- panel の検索 placeholder は `名前で検索` だが、repository query は NPC 属性も検索対象に含める。

#### 構造設計不整合

- backend では件数母集団が `snapshot.targetCount` にそろっている一方、検索主語は UI 文言より広い。画面文言と query 責務が一致していない。
- panel fallback は `NPC 件数` と `生成対象` を並べるが、backend 一覧 DTO 自体は件数の意味を返さないため、view がどの count を主表示に採るべきかを境界で説明できない。

#### 責務分離不足

- presenter は `progress.targetCount` を主要件数にし、panel fallback は `targetSummary.targetCount` を別名で重ねる。意味が近い count を 2 か所で管理している。
- search の責務が「名前 UI」と「属性検索 query」で分かれており、どちらが仕様かを境界で固定していない。

#### 変更候補

- 件数主語は現行の `snapshot.targetCount` を基準に維持し、frontend の主要表示ラベルと fallback 表示語をその主語へそろえる候補。
- 検索責務は、UI 文言を query 実態へ寄せるか、query を UI 文言へ寄せるかのどちらかへ統一する候補。

#### 変更不要範囲

- phase summary と processing target read model を結ぶ backend の page 取得経路自体は共通 read model として分離されている。単語翻訳のような母集団分裂は観測していない。
- `PersonaGenerationPhaseUseCase` の page state 管理は単一 phase 前提として十分である。

### 本文翻訳

#### observed

- `body_translation_input_snapshot.go` は snapshot 全 field 数を `TargetCount` に入れ、完全一致辞書除外で provider へ送らない項目だけ `ProviderTargetCount` を別で増減する。
- `BodyTranslationPhaseService.buildProgress` は `TotalCount` と `TargetCount` に `snapshot.TargetCount` を使い、percent の分母は `snapshot.ProviderTargetCount` を使う。
- 一覧 read model は `JOB_TRANSLATION_FIELD` を参照し、`output_status != 'dictionary_exact_match'` の行だけを返す。query の母集団は provider target 側に近い。
- presenter は主要件数に `summary.progress.targetCount` を使い、別途 `summary.requestSummary.providerTargetCount` を `AI 送信対象` として持つ。
- panel fallback でも `本文翻訳対象` と `AI 送信対象` を併記する。
- body usecase / store は `body_translation` と `translation_complete` の page state を phase 別に保持する。

#### 構造設計不整合

- progress / metrics は snapshot 全体件数、一覧は `dictionary_exact_match` 除外後件数を示す。件数の主語が progress と一覧で一致していない。
- percent 計算の分母だけ `ProviderTargetCount` を使い、表示 count は `TargetCount` を使うため、同じ progress block 内でも count の意味が 2 系統ある。

#### 責務分離不足

- body phase は summary、requestSummary、processing target read model、translation_complete page state を同一 store に抱える。対象 phase と完了 phase をまたぐ page state は妥当だが、今回の修正対象は `body_translation` と `translation_complete` で分けて扱う必要がある。
- panel fallback は 2 系統の count を説明できるが、backend 一覧 DTO は「この totalCount が本文翻訳対象か AI 送信対象か」を返さない。

#### 変更候補

- 本文翻訳は `body_translation` 一覧と一致させる主要 count を 1 つ固定し、その count を presenter / panel / read model で同じ主語に寄せる候補。
- `progress.targetCount` を残す場合は、一覧 total と別物であることを UI 境界で明示する候補。必要観点の「他の件数表示と合う」を優先するなら、summary 側 count の再設計が必要な候補。
- `processingTargetPageStatesByPhase` の設計は維持し、今回の修正対象を `body_translation` state に限定する候補。

#### 変更不要範囲

- `translation_complete` 用 page state と query は、今回の必要観点 3 件の対象外である。`JobRunPage.svelte` の complete 画面導線まで巻き込む必要は観測していない。
- 一覧 query 自体の「辞書完全一致除外」という対象選別は、画面設計の `辞書置換対象外` と整合している。

## 横断の構造品質

### 責務過多

- phase panel が、進行状況メトリクス、summary fallback 一覧、検索 placeholder、実一覧 page state の 4 役を同時に持つ。件数主語の調整責務が view 側へ漏れている。

### 責務分離不足

- 共通 `ProcessingTargetListResponse` は page 情報だけを返し、count の意味を返さない。phase ごとに件数主語が違うのに、Wails DTO と frontend contract がそれを表現できない。
- 共通 read model service は phase 正規化までを担うが、phase summary と整合する read model 主語を担保しない。単語翻訳と本文翻訳の不一致を止めるガードがない。

### 構造設計不整合

- `ProcessingTargetController` は top-level `metadata` を常に空配列で返す。frontend contract と store clone はその配列を保持するが、UI では利用していない。public seam に未使用要素が残っている。
- `ProcessingTargetListWrapper` は `countText` を持つが、3 フェーズ画面はいずれも使っていない。件数表示の補助責務を用意したまま未接続である。

### 未使用コード

- top-level `ProcessingTargetListResponse.metadata` は、controller で空配列固定、store で clone のみ、画面未使用である。今回の主問題ではないが、未使用 public shape として残っている。

### コーディング規約逸脱

- `coding-guidelines.md` の「同じ責務、同じ変更理由、同じ検証単位でそろえる」に対して、件数意味の決定が summary、read model、panel fallback に分散している。特に単語翻訳と本文翻訳は、同じ `処理対象一覧` 変更理由で複数層の別 count を同時に読む。

## 実装範囲候補

### 実装範囲に入れるべき構造項目

- 単語翻訳の母集団不一致を解消するための backend read model 主語の固定、または summary count 主語の固定。
- 単語翻訳 panel / presenter の count 主語統一。`対象`、`対象語件数`、`AI 送信対象`、一覧 total のどれを一致基準にするかを 1 系統へ寄せる項目。
- NPC ペルソナ生成の検索責務統一。UI placeholder と repository query 対象のずれを解消する項目。
- 本文翻訳の主要 count 統一。`progress.targetCount`、`requestSummary.providerTargetCount`、一覧 total の関係を 1 つの主語へ揃える項目。
- Wails DTO / frontend contract で、一覧 total の意味を phase 側が再解釈しなくて済む形へ整理する項目。少なくとも phase presenter 側で count 主語を明示できる境界が必要である。

### 実装範囲に入れない項目

- `translation_complete` 画面の page state と query。
- `ProcessingTargetListPanel.svelte` のページング UI 自体。
- 3 フェーズの検索イベント配線や page 移動配線の仕組み自体。
- body phase の field result 表示、output readiness、provider state 表示。
- phase 以外の Wails controller や job lifecycle 全体。
- top-level `ProcessingTargetListResponse.metadata` の整理だけを目的にした単独改修。今回の必要観点 3 件を満たすための優先項目ではない。

## 残り不足

- 単語翻訳で一覧の正本を `AI 送信対象語` に寄せるのか、`翻訳ジョブ内辞書` に寄せるのかは人間判断待ちである。
- 本文翻訳で一覧と一致させる count を `TargetCount` にするのか `ProviderTargetCount` にするのかは人間判断待ちである。
- NPC ペルソナ生成の検索仕様を UI 文言へ合わせるか query 実態へ合わせるかは人間判断待ちである。

## 残留リスク

- 単語翻訳は主語を fixed しないまま frontend だけを調整すると、開始前と開始後で別 count を示す可能性が残る。
- 本文翻訳は progress percent の分母と主要 count の分母が別のまま残る可能性がある。
- 共通 read model を広く変えすぎると、対象外の `translation_complete` 画面へ波及する可能性がある。

## 推奨 next step

- `designer` が、単語翻訳と本文翻訳の count 主語を人間判断で固定する。
- その判断を前提に、phase ごとの count 主語、検索主語、DTO 主語を `implementation-scope` へ分解する。
