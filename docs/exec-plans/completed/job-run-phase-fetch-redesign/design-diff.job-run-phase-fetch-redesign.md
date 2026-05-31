# job-run-phase-fetch-redesign review diff

## 概要

- 図化目的: 取得・表示フローの作り直しによる構造差分を人間設計レビューで確認する。起動主体と取得回数の変更、処理対象一覧の反映取りこぼし防止（summary は独立反映）、開き直し時の再取得、初回取得中の操作排他、reactive 反映経路、可否判断の責務移動（backend→frontend application 層）、body 取得経路の統合（専用取得廃止）、等価性条件を対比する。
- 根拠参照: `./detail-spec-diff.md`、`./screen-design-diff.job-run.md`、`./plan.md`（スコープ拡大節、第2回設計レビューの追加決定節）、`frontend/src/application/usecase/term-translation-phase/term-translation-phase.usecase.ts`、`frontend/src/ui/screens/job-run/JobRunPage.svelte`、`internal/service/term_translation_phase_service.go`
- 範囲: `JobRunPage`（取得起動経路）、各段階 `*PhaseUseCase`（`fetchSummaryAndReadiness` と反映ガード）、`store`（状態保持）、`viewModel`（derived 変換）、`JobRunPage`（`currentProcessingTargetPageState` の derived 評価）、backend service（`readinessFromState` 等の可否導出関数）、gateway-contract（`*ReadinessResponse` / `*ActionEnablement` DTO）、body 専用取得（`GetBodyTranslationOutputReadiness`）

---

## 図 1: 取得起動経路の before/after（コンポーネント図）

### before（従来）

```mermaid
flowchart TB
    JRP_effect["JobRunPage\n$effect / onMount\nsetJobId を term/persona/body\n全段階に同時発火"]:::removed
    TUC["TermPhaseUseCase\nfetchSummaryAndReadiness\nsummary + readiness +\nprocessingTarget を\nPromise.all で同時取得"]:::removed
    PUC["PersonaPhaseUseCase\nfetchSummaryAndReadiness\n同上"]:::removed
    BUC["BodyPhaseUseCase\nfetchSummaryAndReadiness\n同上"]:::removed
    Bridge["Wails bridge\n最大 9 本の呼び出しが\n同時発火"]:::removed
    Store["各段階 store\n取得結果を個別に保持"]:::unchanged
    VM["viewModel\n各段階の store 購読"]:::unchanged
    Page["JobRunPage\ncurrentProcessingTargetPageState\n$derived.by で評価"]:::unchanged

    JRP_effect -->|jobId 設定と同時に全段階を並列起動| TUC
    JRP_effect -->|同上| PUC
    JRP_effect -->|同上| BUC
    TUC -->|Promise.all: summary + readiness +\nprocessingTarget 計 3 本| Bridge
    PUC -->|Promise.all: 同上 計 3 本| Bridge
    BUC -->|Promise.all: 同上 計 3 本| Bridge
    Bridge -->|各段階の取得結果を返す| Store
    Store -->|subscribe| VM
    VM -->|$derived.by 再評価| Page

    classDef removed fill:#ffebee,stroke:#c62828,color:#7f1d1d
    classDef unchanged fill:#fff8e1,stroke:#f9a825,color:#4e342e
```

### after（改定）

```mermaid
flowchart TB
    JRP_effect_new["JobRunPage\n$effect\n表示中段階のみに\nsetJobId を発火"]:::added
    PhaseSwitch["JobRunPage\n段階切り替え操作\n切り替え先段階だけ取得開始"]:::added
    ActiveUC["表示中段階の PhaseUseCase\nfetchSummaryAndReadiness\nsummary + processingTarget を取得\n（可否はフロント導出のため取得物ではない）"]:::added
    OtherUC["非表示段階の PhaseUseCase\n切り替えまで取得しない"]:::added
    Bridge_new["Wails bridge\n最大 2 本\n（表示中段階の summary + processingTarget）"]:::added
    Reopen["JobRunPage\n画面を開き直した時\n旧遅延取得の sequence を\n無効化して再取得"]:::added
    Store_new["各段階 store\n取得結果と\ninitialFetchDone を保持"]:::added
    VM_new["viewModel\ninitialFetchDone を含む\n各段階の store 購読"]:::added
    Page_new["JobRunPage\ncurrentProcessingTargetPageState\ninitialFetchDone が true の時\nだけ derived 評価"]:::added

    JRP_effect_new -->|表示中 1 段階だけ| ActiveUC
    PhaseSwitch -->|切り替え先 1 段階| ActiveUC
    JRP_effect_new -->|起動しない| OtherUC
    ActiveUC -->|summary + processingTarget 計 2 本\n（readiness 等の可否 DTO は廃止）| Bridge_new
    Reopen -->|sequence 無効化後に再起動| ActiveUC
    Bridge_new -->|取得結果| Store_new
    Store_new -->|subscribe| VM_new
    VM_new -->|$derived.by 再評価| Page_new

    classDef added fill:#e8f5e9,stroke:#2e7d32,color:#1b5e20
```

---

## 図 2: 処理対象一覧の反映取りこぼし防止の before/after（コンポーネント図）

### before（従来 — 非対称ガード）

```mermaid
flowchart TB
    FetchResp["fetchSummaryAndReadiness\nsummary / processingTarget\nを Promise.all で受け取る"]:::removed
    StoreSummary["store.update\nsummary を無条件反映"]:::removed
    SeqCheck["processingTargetListRequestSequence\nガード\n（処理対象だけ適用）"]:::removed
    StoreTarget["store.update\nprocessingTarget を反映"]:::unchanged
    note1["注意: summary はガードなしで反映される。\n後発取得が先発応答より前に\nstore に入ると古い値で上書く。\n先行取得だけ完了すると一覧が空のまま残る"]:::removed

    FetchResp --> StoreSummary
    FetchResp --> SeqCheck
    SeqCheck -->|sequence 一致の時だけ| StoreTarget
    note1 -.非対称.-> StoreSummary

    classDef removed fill:#ffebee,stroke:#c62828,color:#7f1d1d
    classDef unchanged fill:#fff8e1,stroke:#f9a825,color:#4e342e
```

### after（改定 — 独立反映 + 一覧は連番ガード）

```mermaid
flowchart TB
    FetchSummary["summary 取得経路\nsummary を独立して取得・反映"]:::added
    FetchTarget["processingTarget 取得経路\nprocessingTarget を独立して取得"]:::added
    SeqCheck_new["processingTargetListRequestSequence\nガード\nprocessingTarget の反映に適用\n先行取得の完了で一覧を空に残すことを防ぐ"]:::added
    StoreSummary_new["store.update\nsummary を反映\n（一覧の反映取りこぼしを理由に止めない）"]:::added
    StoreTarget_new["store.update\nprocessingTarget を反映\n（ガード外の場合は破棄）"]:::added
    note2["根拠: detail-spec-diff\nterm-REQ-007\n進捗の要約と処理対象一覧は別経路で独立反映する\n先行取得だけ完了して一覧が空のまま残る状態を許可しない"]:::added

    FetchSummary --> StoreSummary_new
    FetchTarget --> SeqCheck_new
    SeqCheck_new -->|sequence 一致の時だけ| StoreTarget_new
    note2 -.仕様根拠.-> SeqCheck_new

    classDef added fill:#e8f5e9,stroke:#2e7d32,color:#1b5e20
```

---

## 図 3: 開き直し時の再取得と旧取得破棄（シーケンス図）

```mermaid
sequenceDiagram
    participant User as 利用者
    participant JRP as JobRunPage
    participant UC as 表示中段階 PhaseUseCase
    participant Bridge as Wails bridge

    User->>JRP: 翻訳実行画面を開く（初回）
    JRP->>UC: setJobId(jobId) + sequence=1 発行
    UC->>Bridge: summary / processingTarget 取得（sequence=1）
    note over Bridge: 約 15 秒の遅延（IPC 飽和の疑い）

    User->>JRP: 翻訳実行画面を閉じる
    JRP->>UC: dispose / sequence をリセット（sequence=2 発行）

    User->>JRP: 翻訳実行画面を開き直す
    JRP->>UC: setJobId(jobId) + sequence=3 発行
    UC->>Bridge: 新規取得（sequence=3）

    Bridge-->>UC: sequence=1 の遅延応答が返る
    UC->>UC: sequence=1 ≠ 現在 sequence=3\n→ processingTarget を破棄
    Bridge-->>UC: sequence=3 の応答が返る
    UC->>UC: sequence=3 = 現在 sequence=3\n→ processingTarget を store に反映
    UC-->>JRP: 処理対象一覧を件数分表示
```

---

## 図 4: 初回取得中ローディングレイヤーによる操作排他（シーケンス図）

```mermaid
sequenceDiagram
    participant User as 利用者
    participant JRP as JobRunPage
    participant Layer as 初回取得中ローディングレイヤー
    participant UC as 表示中段階 PhaseUseCase
    participant Bridge as Wails bridge

    JRP->>UC: 初回取得開始（initialFetchDone=false）
    JRP->>Layer: initialFetchDone=false → フェーズ画面全体を覆うオーバーレイを表示
    note over Layer: フェーズ画面全体の操作を受け付けない\n（検索・ページ・行展開・上部状態区画の操作を含む）

    User->>Layer: 検索操作を試みる
    Layer-->>User: 操作をブロック（オーバーレイがフェーズ画面全体を覆っている）

    UC->>Bridge: summary / processingTarget 取得
    Bridge-->>UC: 取得完了
    UC->>UC: store に summary と processingTarget を反映\ninitialFetchDone=true に更新

    JRP->>Layer: initialFetchDone=true → オーバーレイを外す
    note over JRP: フェーズ画面全体の操作を受け付ける
    JRP-->>User: 処理対象一覧を件数分表示
```

---

## 図 5: bridge → store → viewModel → 画面の reactive 反映経路（コンポーネント図）

```mermaid
flowchart TB
    Bridge_path["Wails bridge\nsummary / processingTarget を返す\n（可否 DTO は応答に含まない）"]:::added
    UC_path["表示中段階 PhaseUseCase\nsequence ガードを経て\ninitialFetchDone=true を含む\n取得値を store に反映"]:::added
    Store_path["表示中段階 store\nsummary / processingTargetPageState /\ninitialFetchDone を保持\n（nextPhaseReadiness / actionEnablement は保持しない）"]:::added
    VM_path["viewModel\nstore を subscribe して\ninitialFetchDone を含む\n各フィールドを変換"]:::added
    Page_derived["JobRunPage\ncurrentProcessingTargetPageState\n$derived.by\ninitialFetchDone=true かつ\nitems.length=0 の時だけ\npageState を評価"]:::added

    Summary_path["進捗母数経路\nsummary.aiTargetCount\n（独立した集計）"]:::unchanged
    Target_path["処理対象一覧経路\nprocessingTargetPageState.items\n（独立した取得）"]:::added

    Bridge_path -->|取得値を返す| UC_path
    UC_path -->|store.update| Store_path
    Store_path -->|subscribe callback| VM_path
    VM_path -->|$derived.by 再評価| Page_derived

    Bridge_path -->|summary.aiTargetCount| Summary_path
    Bridge_path -->|processingTargetPageState.items| Target_path

    Summary_path -.独立した値\n一致保証なし.-> Target_path

    classDef added fill:#e8f5e9,stroke:#2e7d32,color:#1b5e20
    classDef unchanged fill:#fff8e1,stroke:#f9a825,color:#4e342e
```

---

## 図 6: 可否判断の責務移動 before/after（コンポーネント図）

この図はスコープ拡大分（`plan.md` 「スコープ拡大」節・「第2回設計レビューの追加決定」節、`detail-spec-diff.md` の `*-REQ-008` / `body-...-REQ-007`）に対応する。before は backend service が事実状態から可否を導出して応答 DTO に含める。after は backend が事実状態だけ返し、frontend の application 層が同じ導出を行う。

### before（従来 — backend で可否を導出）

```mermaid
flowchart TB
    FactState_be["backend service\n段階データ事実状態を取得\nフェーズ状態、対象件数、完了件数、\nジョブ終端状態、エラー種別、\n実行設定の構成有無などの事実値"]:::removed
    CanStart_be["backend service\nreadinessFromState 等の関数\n次段階開始可否を導出\ncanStartNextPhase / blockedReason"]:::removed
    ActionEnable_be["backend service\ntermTranslationStartBlockedReason 等\n各操作可否を導出\ncanStart / canPause / canResume /\ncanRetry / canCancel +\n各 *BlockedReason\npersona / body も同型"]:::removed
    DTO_be["backend controller\nDTO に可否判断結果を含めて返す\nTermTranslationNextPhaseReadinessResponse\nPersonaGenerationBodyReadinessResponseDTO\nBodyTranslationOutputReadinessResponseDTO\nActionEnablement 各 can* と *BlockedReason"]:::removed
    AppLayer_be["frontend application 層\n取得した可否判断をそのまま利用\n（自前導出しない）"]:::removed

    FactState_be -->|事実値| CanStart_be
    FactState_be -->|事実値| ActionEnable_be
    CanStart_be -->|canStartNextPhase + blockedReason| DTO_be
    ActionEnable_be -->|actionEnablement 各 can* + *BlockedReason| DTO_be
    DTO_be -->|bridge 経由で frontend へ| AppLayer_be

    classDef removed fill:#ffebee,stroke:#c62828,color:#7f1d1d
```

### after（改定 — frontend application 層で可否を導出）

```mermaid
flowchart TB
    FactState_fe["backend service\n段階データ事実状態だけ返す\nフェーズ状態、対象件数、完了件数、\nジョブ終端状態、エラー種別、\n実行設定の構成有無などの事実値"]:::added
    DTO_fe["backend controller\nDTO に事実状態だけ含める\ncanStartNextPhase / blockedReason /\ncan* / *BlockedReason を含まない\nbody は BodyTranslationPhaseSummaryResponse に\ncompletedFieldCount / statusConsistent / outputCount を含める"]:::added
    AppLayer_fe["frontend application 層\n段階データ事実状態を受け取り\n可否を自前で導出する\n（term / persona / body 同型）"]:::added
    CanStart_fe["frontend application 層\n次段階開始可否を導出\n（term: 次段階開始可否\n persona: 本文翻訳段階開始可否\n body: 成果物出力確認可否）"]:::added
    ActionEnable_fe["frontend application 層\n各操作可否を導出\ncanStart / canPause / canResume /\ncanRetry / canCancel /\nbody: canCheckOutputReadiness\n各 *BlockedReason\nterm / persona / body 同型"]:::added

    FactState_fe -->|事実状態だけ含む DTO| DTO_fe
    DTO_fe -->|bridge 経由で frontend へ| AppLayer_fe
    AppLayer_fe -->|事実値を入力| CanStart_fe
    AppLayer_fe -->|事実値を入力| ActionEnable_fe

    classDef added fill:#e8f5e9,stroke:#2e7d32,color:#1b5e20
```

---

## 図 7: body 取得経路の統合 before/after（コンポーネント図）

この図はスコープ拡大分（`detail-spec-diff.md` の `body-translation-phase-REQ-007`、`Q-005` 回答）に対応する。before は body 段階で段階要約取得と専用取得の 2 本が存在する。after は専用取得を廃止し段階要約取得 1 本へ統合、出力可否はフロント導出。

### before（従来 — 2 経路）

```mermaid
flowchart TB
    SummaryReq_be["frontend\nBodyPhaseUseCase\n段階要約取得\nGetBodyTranslationPhaseSummary\n→ BodyTranslationPhaseSummaryResponse"]:::removed
    ReadinessReq_be["frontend\nBodyPhaseUseCase\n成果物出力確認専用取得\nGetBodyTranslationOutputReadiness\n→ BodyTranslationOutputReadinessResponse\n（ready / blockedReason を含む）"]:::removed
    Bridge_body_be["Wails bridge\n2 本の取得を同時発火"]:::removed
    AppUse_be["frontend application 層\n専用取得の ready / blockedReason\nをそのまま利用"]:::removed

    SummaryReq_be -->|段階要約を取得| Bridge_body_be
    ReadinessReq_be -->|出力確認可否を取得| Bridge_body_be
    Bridge_body_be -->|BodyTranslationPhaseSummaryResponse| AppUse_be
    Bridge_body_be -->|BodyTranslationOutputReadinessResponse\n（ready / blockedReason）| AppUse_be

    classDef removed fill:#ffebee,stroke:#c62828,color:#7f1d1d
```

### after（改定 — 1 経路へ統合）

```mermaid
flowchart TB
    SummaryReq_fe["frontend\nBodyPhaseUseCase\n段階要約取得 1 本のみ\nGetBodyTranslationPhaseSummary\n→ BodyTranslationPhaseSummaryResponse\n（completedFieldCount / statusConsistent /\n outputCount などを含む事実状態に集約）"]:::added
    Bridge_body_fe["Wails bridge\n1 本の取得"]:::added
    AppUse_fe["frontend application 層\n段階要約の事実状態から\n成果物出力確認可否（ready相当）と\nblockedReason 相当を導出"]:::added
    Screen_fe["画面（viewModel / store 経由）\n導出した可否・理由を受け取り表示"]:::added
    Abolished["GetBodyTranslationOutputReadiness\nBodyTranslationOutputReadinessResponse\n廃止"]:::removed

    SummaryReq_fe -->|段階要約を取得| Bridge_body_fe
    Bridge_body_fe -->|BodyTranslationPhaseSummaryResponse\n（事実状態を集約）| AppUse_fe
    AppUse_fe -->|導出した可否・理由を渡す| Screen_fe

    classDef added fill:#e8f5e9,stroke:#2e7d32,color:#1b5e20
    classDef removed fill:#ffebee,stroke:#c62828,color:#7f1d1d
```

---

## 図 8: 等価性条件の確認（コンポーネント図）

この図はスコープ拡大分（`detail-spec-diff.md` の `term-REQ-008`、`persona-REQ-008`、`body-...-REQ-007` の等価性条件）に対応する。責務再配置の前後で、同じ事実入力に同じ可否・理由が得られることを示す。

```mermaid
flowchart LR
    FactIn["段階データ事実状態\n（共通入力）\nフェーズ状態、対象件数、完了件数\nジョブ終端状態、エラー種別\n実行設定の構成有無など"]:::unchanged

    Before["before\nbackend service が\n同じ事実状態から\n可否・理由を導出\ncanStartNextPhase / blockedReason\nactionEnablement 各 can* + *BlockedReason\nbody: ready / blockedReason（専用取得）"]:::removed

    After["after\nfrontend application 層が\n同じ事実状態から\n可否・理由を導出\n導出ロジックを application 層に移設"]:::added

    Equiv["等価性条件（仕様で固定）\n同じ事実入力に対して\nbefore の backend 導出結果と\nafter の frontend 導出結果が\n一致すること\n（term / persona / body 全段階で成立）"]:::added

    FactIn -->|入力| Before
    FactIn -->|入力| After
    Before -.再配置前後で同結果.-> Equiv
    After -.再配置前後で同結果.-> Equiv

    classDef unchanged fill:#fff8e1,stroke:#f9a825,color:#4e342e
    classDef removed fill:#ffebee,stroke:#c62828,color:#7f1d1d
    classDef added fill:#e8f5e9,stroke:#2e7d32,color:#1b5e20
```

---

## 差分凡例

- 赤: 削除する要素または経路を示す。
- 緑: 追加する要素または経路を示す。
- 黄色: 変更しない要素または経路を示す。

---

## 各箱の説明

### 図 1（取得起動経路）

- `JobRunPage $effect / onMount（before）`: ジョブ選択時に term / persona / body の全 3 段階に `setJobId` を同時発火する既存の起動経路。削除対象。
- `JobRunPage $effect（after）`: ジョブ選択時に表示中の 1 段階だけに `setJobId` を発火する改定後の起動経路。追加する。
- `JobRunPage 段階切り替え操作`: 利用者が段階タブを切り替えた時に、切り替え先 1 段階だけ取得を開始する経路。追加する。
- `表示中段階 PhaseUseCase（after）`: 表示中の 1 段階の `fetchSummaryAndReadiness` を呼ぶ。取得本数は summary と processingTarget の最大 2 本。可否 DTO は取得物ではなくなるため、従来の readiness 取得本数は 0 本になる。追加する。
- `非表示段階 PhaseUseCase`: 切り替えまで取得を行わない。段階切り替え時に初めて起動する。追加する。
- `Wails bridge（before）`: 最大 9 本（3 段階 × 3 本）が同時発火する既存経路。削除対象。
- `Wails bridge（after）`: 最大 2 本（表示中 1 段階の summary + processingTarget）に削減する。readiness 等の可否 DTO 取得はなくなる。追加する。
- `JobRunPage 開き直し`: 翻訳実行画面を閉じて再び開いた時に、旧取得の sequence を無効化して再取得を開始する経路。追加する。
- `各段階 store`: 取得結果を保持する。`initialFetchDone` フラグを新規に追加する。`nextPhaseReadiness` や `actionEnablement` は保持しなくなる。
- `viewModel`: 各段階の store を購読して変換する。`initialFetchDone` を含む点が変更。可否判断値は store に保持しないため、viewModel の変換対象からも外れる。
- `JobRunPage currentProcessingTargetPageState`: `$derived.by` で処理対象ページ状態を評価する。`initialFetchDone` を評価条件に加える点が変更。

### 図 2（処理対象一覧の反映取りこぼし防止）

- `fetchSummaryAndReadiness（before）`: `Promise.all` で取得し、summary は無条件反映、processingTarget だけ sequence ガード内で反映する非対称な既存実装。削除対象。
- `store.update summary（before）`: ガードなしで summary を反映する。先行取得だけが完了すると一覧が空のまま残る原因になる。削除対象。
- `processingTargetListRequestSequence ガード（before）`: 処理対象だけに適用している既存ガード。削除対象（改定後は処理対象一覧の取得に適用し、summary は独立経路で別管理）。
- `summary 取得経路（after）`: summary を独立して取得・反映する。一覧の反映取りこぼしを理由に止めない。追加する。
- `processingTarget 取得経路（after）`: processingTarget を独立して取得する。追加する。
- `processingTargetListRequestSequence ガード（after）`: processingTarget の反映に適用する改定後のガード。先行取得だけが完了して一覧が空のまま残ることを防ぐ。追加する。
- `store.update（after）`: summary と processingTarget を独立した経路でそれぞれ反映する。processingTarget は sequence が一致する時だけ反映し、ガード外の場合は破棄する。追加する。
- 注: 第1回図2の「3値を同一ガードで揃える」表現は、仕様注記3（`detail-spec-diff.md`）に従い修正済み。summary は独立反映であり、3値の束ね反映は本図では採用しない。

### 図 3（開き直し時の再取得）

- sequence=1 の遅延応答が sequence=3 で無効化されて破棄される経路が差分の核心。
- 閉じる操作で sequence をリセット（または dispose）し、開き直しで新しい sequence を発行することで、旧遅延応答が store に入ることを防ぐ。
- 改定後は summary と processingTarget の 2 本が対象。readiness 等の可否 DTO は取得物ではなくなる。

### 図 4（初回取得中ローディングレイヤー）

- `初回取得中ローディングレイヤー`: `initialFetchDone=false` の間、`phase-loading-overlay`（`position: absolute; inset: 0; z-index: 20`）としてフェーズ画面全体を覆うオーバーレイを表示し、操作を受け付けない新規エレメント。上部状態区画も含むフェーズ画面全体に排他が及ぶ。`screen-design-diff.job-run.md` 差分 2 に対応する。
- 検索操作・ページ切り替え・行展開を含むフェーズ画面全体の操作は `initialFetchDone=true` になるまでブロックする。初回取得と利用者操作の同時進行自体を起こさせない排他。
- 改定後は summary と processingTarget の両方の取得完了を `initialFetchDone=true` の判定に使う。

### 図 5（reactive 反映経路）

- `PhaseUseCase → store`: sequence ガードを経て `initialFetchDone=true` を含む取得値を入れる。`nextPhaseReadiness` や `actionEnablement` は store に入れない。変更部分。
- `store`: `initialFetchDone` フラグを新規に保持する。可否判断値は保持しない。変更部分。
- `viewModel`: `initialFetchDone` を含む変換を行う。可否判断値は変換対象から外れる。変更部分。
- `currentProcessingTargetPageState $derived.by`: `initialFetchDone=true` かつ `items.length=0` の時だけ pageState を評価する評価条件の追加。変更部分。
- `進捗母数経路（summary.aiTargetCount）` と `処理対象一覧経路（processingTargetPageState.items）`: 別 bridge 呼び出し・別集計の独立した経路であることを明示する。変更しない接続先。
- `Wails bridge`: 可否 DTO を応答に含まなくなる。変更部分。

### 図 6（可否判断の責務移動）

- `backend service（before）`: `readinessFromState`、`termTranslationStartBlockedReason` 等が事実状態から次段階開始可否と操作可否を導出する。削除対象（application 層へ移設）。
- `backend controller DTO（before）`: `TermTranslationNextPhaseReadinessResponse`（`canStartNextPhase`、`blockedReason`）、`PersonaGenerationBodyReadinessResponseDTO`、`BodyTranslationOutputReadinessResponseDTO`、`ActionEnablement` の各 `can*` と `*BlockedReason` を含む。削除対象。
- `backend service（after）`: 段階データ事実状態だけ返す。可否の導出はしない。追加する。
- `backend controller DTO（after）`: 事実状態だけ含む。可否判断結果の真偽値と理由文字列を含まない。body は `BodyTranslationPhaseSummaryResponse` に `completedFieldCount`、`statusConsistent`、`outputCount` を含める。追加する。
- `frontend application 層（after）`: 事実状態を受け取り、次段階開始可否と操作可否の両方を自前で導出する。term / persona / body 同型。追加する。
- `等価性条件`: 再配置前後で同じ事実入力に同じ可否・理由が得られることが仕様で固定済み。`detail-spec-diff.md` の `term-REQ-008`、`persona-REQ-008`、`body-...-REQ-007` が根拠。

### 図 7（body 取得経路の統合）

- `GetBodyTranslationOutputReadiness + BodyTranslationOutputReadinessResponse（before）`: body 成果物出力確認の専用取得。削除対象。
- `BodyTranslationPhaseSummaryResponse（before）`: 事実状態の一部だけを含む既存の段階要約取得応答。修正対象（`completedFieldCount`、`statusConsistent`、`outputCount` を集約して拡張する）。
- `段階要約取得 1 本（after）`: `BodyTranslationPhaseSummaryResponse` に `completedFieldCount`、`statusConsistent`、`outputCount` を含める形で事実状態を集約する。取得経路は 2 本から 1 本になる。追加する。
- `frontend application 層（after）`: 段階要約取得の事実から成果物出力確認可否（`ready` 相当）と `blockedReason` 相当を導出する。追加する。

### 図 8（等価性条件の確認）

- `段階データ事実状態（共通入力）`: before の backend 導出と after の frontend 導出が共通して使う事実入力。変更しない（事実の取得元は backend のまま）。
- `before の backend 導出`: `readinessFromState` 等が行っていた可否算出。削除対象。
- `after の frontend 導出`: application 層が引き継ぐ可否算出。追加する。
- `等価性条件`: 同じ事実入力に対して before と after が同じ結果を返すことを仕様で固定済み。term / persona / body 全段階に成立することが条件。

---

## 検証

- Mermaid 記述確認: flowchart TB（図 1 × 2、図 2 × 2、図 5、図 6 × 2、図 7 × 2）、flowchart LR（図 8）、sequenceDiagram（図 3、図 4）の合計 12 図の Mermaid コードブロックを確認した。各図に図種別、箱または参加者、接続、凡例 classDef が記述されていることを確認した。
- Markdown 確認: 各図に見出し、説明、差分凡例、各箱の説明が揃っていることを確認した。
- 差分凡例確認: 赤（removed）、緑（added）、黄色（unchanged）の 3 区分を全図で統一している。
- 追加予定・削除予定・変更しない接続先の区別: 削除対象を赤（removed）、追加対象を緑（added）、変更しない要素を黄色（unchanged）で区別している。全図で統一している。
- detail-spec-diff.md との整合確認:
  - 図 2（after）は `term-REQ-007`（summary と処理対象一覧の独立反映、一覧の反映取りこぼし防止）に対応し、第1回図2の「3値を同一ガードで揃える」表現を修正した。
  - 図 3 は `translation-job-management-REQ-006`（開き直し時に再取得、旧取得結果を破棄）に対応する。取得対象を summary と processingTarget の 2 本に修正した。
  - 図 4 は `term-REQ-007`（初回取得完了までフェーズ画面全体の操作を行えない）および `screen-design-diff.job-run.md` 差分 2 に対応する。`phase-loading-overlay` がフェーズ画面全体（上部状態区画を含む）を覆うオーバーレイであることを反映した。
  - 図 6 は `term-REQ-008`、`persona-REQ-008`、`body-...-REQ-007`（可否判断の責務再配置、backend は事実だけ返す）に対応する。term / persona / body 同型を示した。
  - 図 7 は `body-...-REQ-007`（専用取得廃止、段階要約取得一本化、出力可否はフロント導出）に対応する。
  - 図 8 は `term-REQ-008`、`persona-REQ-008`、`body-...-REQ-007` の等価性条件を図で示した。
- screen-design-diff.job-run.md との整合確認: 差分 6（操作可否・次段階開始可否のフロント導出）は図 6 に対応する。差分 7（body 成果物出力確認の取得経路統合）は図 7 に対応する。
- backend / bridge 仕様変更の範囲確認: bridge 伝送機構自体の破壊的変更は含まない。bridge の呼び出し本数削減（before: 最大 9 本、after: 最大 2 本）と、backend 応答 DTO から可否判断値を除外することは含む。これは仕様で確定済みのスコープ拡大分である。
