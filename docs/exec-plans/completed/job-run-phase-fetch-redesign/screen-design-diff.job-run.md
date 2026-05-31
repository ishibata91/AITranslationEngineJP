# 画面設計差分: job-run-phase-fetch-redesign

- `skill`: design-bundle
- `status`: approved-with-impl-alignment（第4回。第3回までの差分本文は人間設計レビュー承認済み。Storybook 人間レビュー中に承認された frontend 実装が、初回取得中ローディングの覆う範囲を「処理対象一覧領域」から「フェーズ画面全体」へ拡大したため、本書を実装の事実へ整合させた。整合のため変更したのは差分 1・差分 2 のローディング範囲と操作排他範囲の記述である。固定 selector の値は据え置く）
- `screen_id`: `term-translation-phase`, `persona-generation-phase`, `body-translation-phase`
- `apply_target`: `docs/screen-design/screens/term-translation-phase.md`, `docs/screen-design/screens/persona-generation-phase.md`, `docs/screen-design/screens/body-translation-phase.md`
- `source_plan`: `./plan.md`
- `detail_spec_diff`: `./detail-spec-diff.md`（表示規則、次段階開始可否のフロント導出、操作可否のフロント導出、body 成果物出力確認の取得経路統合を根拠にする）
- `revision`: 第4回（2026-05-31 Storybook 人間レビューで承認された frontend 実装へ整合）。整合点は、初回取得中ローディングの覆う範囲を「処理対象一覧領域の最前面」から「フェーズ画面全体を覆うオーバーレイ」へ広げ、操作排他の範囲も「処理対象一覧領域の操作」から「フェーズ画面全体の操作」へ広げたことである。固定 selector の値（`<phase-prefix>-processing-target-loading`）は実装に合わせて据え置く。差分 6・差分 7（フロント導出表示）の記述は変更しない。
- `revision`: 第3回（2026-05-31 未決 `Q-004`、`Q-005` の人間回答を反映）。差分 6 を、次段階開始可否に加えて操作可否（開始・一時停止・再開・再試行・取り消し・成果物出力確認の各活性と理由）もフロント導出になる画面影響へ拡張した。body の成果物出力確認の表示が専用取得ではなく段階要約取得の事実から導出される取得経路統合の画面影響も反映した。

## この差分の対象と前提

- 対象: 翻訳実行画面（`JobRunPage`）が表示する単語翻訳段階、NPC ペルソナ生成段階、本文翻訳段階の処理対象一覧領域である。
- 前提: 3 段階は処理対象一覧領域と固定 selector を同型で扱う。本差分は単語翻訳段階を基準に書き、persona / body も同型で適用する。
- 前提: 処理対象一覧領域（差分 1〜5）の画面表示規則と固定 selector は backend / bridge の仕様変更に依存しない。
- 前提: 次段階開始可否と操作可否の表示（差分 6）は、backend が可否判断を返さず段階データ事実状態だけ返す責務再配置（詳細仕様差分の `*-REQ-008` / `body-...-REQ-007`、`Q-004` 回答）を前提にする。画面は事実状態から application 層が導出した可否と理由を受け取って表示する。
- 前提: body の成果物出力確認の表示は、専用取得（`GetBodyTranslationOutputReadiness`）を廃止し段階要約取得（`BodyTranslationPhaseSummaryResponse`）の事実から application 層が導出した可否・理由を受け取って表示する取得経路統合（詳細仕様差分の `body-...-REQ-007`、`Q-005` 回答）を前提にする。
- 注意: 本差分は恒久的な画面内容だけを書く。実装方式、テスト手順、agent 引き継ぎは書かない。次段階開始可否と操作可否の活性・理由が「props 経由でどこから来るか」の取得元責務だけを示し、導出ロジックの配置や DTO 形は実装範囲で扱う。

## 用語の定義（本差分で固定する日本語の対象）

- 対象: 処理対象一覧領域とは、表示中段階で処理する実体の一覧、総件数、検索欄、ページ操作をまとめて表示する領域である。`ProcessingTargetListWrapper` と `ProcessingTargetListPanel` が実体である。
- 対象: 初回取得とは、翻訳ジョブを開いた初回表示で、表示中段階の処理対象一覧を取得する処理である。
- 対象: 初回取得中ローディングレイヤーとは、初回取得が完了するまでフェーズ画面全体（進行状況兼操作 `E-02` などの上部区画と処理対象一覧領域 `E-04` を含むフェーズ画面の全体）を覆うオーバーレイとして出し、取得中であることを区別させ利用者操作を受け付けなくする表示である。本差分で新規に追加する。承認済み実装では `section.job-run-shell` 全体を覆う絶対配置のオーバーレイとして実装している。
- 対象: 進捗母数とは、進行状況領域（`E-02`）が表示する対象件数である。単語翻訳段階では AI 翻訳対象語件数が該当する。
- 注意: 進捗母数と処理対象一覧の総件数は別経路で取得・集計する独立した値であり、一致する保証はない。

## 差分 1: 処理対象一覧の件数表示・空状態表示の規則を固定する（`E-04` 状態別表示）

`docs/screen-design/screens/term-translation-phase.md` の `E-04 処理対象一覧` の「状態別表示」へ次を追加・明確化する。persona / body の対応エレメントにも同型で適用する。

| 状態 | 表示 |
| --- | --- |
| 取得中 | 初回取得中ローディングレイヤーをフェーズ画面全体を覆うオーバーレイとして表示し、取得中であることを区別させる。オーバーレイがフェーズ画面全体を覆うため、検索・件数・ページ操作・処理対象一覧の操作を含むフェーズ画面全体の操作を受け付けない。 |
| 件数あり | 処理対象一覧の総件数が 1 以上の場合、現在ページの表示範囲を件数として示し、処理対象行を表示する。空状態は表示しない。 |
| 空状態 | 処理対象一覧の総件数が 0 の場合だけ、`処理対象がありません` を表示する。 |
| 行展開 | 開いた処理対象行に metadata を表示する。 |

規則の固定点:

- 規則1: 処理対象一覧の総件数が 1 以上なら、初回取得完了後の初回表示で処理対象行を件数分（現在ページの表示範囲分）表示する。
- 規則2: 処理対象一覧の総件数が 0 の場合だけ空状態 `処理対象がありません` を表示する。総件数が 1 以上のとき空状態を表示する状態は、画面表示上の不具合として扱う。
- 規則3: 進捗母数（`E-02` の対象件数）と処理対象一覧の総件数は別々の値として表示する。値が一致しないことだけを理由に不具合とは扱わない。
- 注意: 件数表示の文言形式（現在ページの表示範囲を示す `pageRangeLabel`）は現状を維持し、本差分では変更しない。

## 差分 2: 初回取得中ローディングレイヤーを新規エレメントとして追加する（`E-04` へ追記）

`E-04 処理対象一覧` の「表示項目」へ次を追加する。persona / body の対応エレメントにも同型で適用する。

| 表示項目 | 取得元 | 表示形式 | 表示条件 | セレクタ（`data-testid`） |
| --- | --- | --- | --- | --- |
| 初回取得中ローディングレイヤー | props | 取得中であることを区別できる表示（フェーズ画面全体を覆うオーバーレイ） | 初回取得が完了するまで表示する | `<phase-prefix>-processing-target-loading` |

表示・操作排他の規則:

- 規則4: 初回取得中ローディングレイヤーは、翻訳ジョブを開いた初回表示で初回取得が完了するまで、フェーズ画面全体（進行状況兼操作 `E-02` などの上部区画と処理対象一覧領域 `E-04` を含むフェーズ画面の全体）を覆うオーバーレイとして表示する。
- 規則5: 初回取得中ローディングレイヤーの表示中は、検索欄入力、ページ切り替え、処理対象行を開く操作に加え、フェーズ画面全体の操作（開始・一時停止・再開・再試行・取り消し・成果物出力確認などの操作行を含む）を受け付けない。オーバーレイがフェーズ画面全体を覆うことで利用者操作を実質不可にし、初回取得と利用者操作による取得を同時に進行させない。
- 規則6: 初回取得が完了した時点でローディングレイヤーを外し、検索・ページ操作・行展開とフェーズ画面全体の操作を受け付ける。
- 規則7: 初回取得中ローディングレイヤーは term / persona / body の各段階で同型に表示する。表示条件、操作排他の範囲（フェーズ画面全体）、外す条件を 3 段階で揃える。
- 注意: 既存の `busy` 表示（処理対象一覧パネルの取得中表示）とローディングレイヤーの関係（どちらを正とするか、併用するか）は実装方式であり、本差分では `busy` ではなく初回取得完了状態（実装の `initialFetchDone=false`）を初回取得中ローディングレイヤーの根拠状態として扱うことを示す。承認済み実装では、フェーズ画面全体を覆う絶対配置オーバーレイ（フェーズ画面の親要素を相対配置にし、オーバーレイを最前面の絶対配置とする）と、処理対象一覧の検索・ページ操作ハンドラを未提供にする入力無効化を併用している。具体的な前面化方式と入力無効化方式の詳細は実装範囲で扱う。

## 差分 3: 段階切り替え時の表示中段階だけ取得する画面挙動を固定する（`E-04` 状態別表示の補足）

- 規則8: 翻訳ジョブを開いた初回表示では、表示中の 1 段階の処理対象一覧だけを取得する。他段階の処理対象一覧は、その段階を表示へ切り替えた時に取得する。
- 規則9: 表示中段階を切り替えた時、切り替え先段階の処理対象一覧は切り替え時点で取得を開始する。取得完了までは取得中であることを区別できる表示（ローディングレイヤーまたは `busy` 表示）を出す。
- 規則10: 段階切り替え直後に取得中表示が一時的に出ることを許容する。取得完了後に切り替え先段階の処理対象一覧か空状態を表示する。

## 差分 4: 検索後リロードで初回 0 件へ戻らない表示挙動を固定する（`E-04` 状態別表示の補足）

- 規則11: 初回取得が完了して処理対象一覧を一度表示した後は、利用者の検索操作・ページ操作なしに、表示中段階の処理対象一覧が初回表示前の空状態へ戻らない。
- 規則12: 翻訳実行画面を一度閉じて再び開いた場合、開き直し時に表示中段階の処理対象一覧の初回取得をやり直す。総件数が 1 以上の段階では、開き直し後の初回表示でも処理対象行を件数分表示し、初回 0 件へ戻らない。
- 規則13: 利用者の検索操作・ページ操作は、初回取得が完了した後にだけ行える（差分 2 の操作排他による）。検索・ページ操作の取得結果は、利用者が要求した検索条件とページ位置を反映する。

## 差分 5: E2E 固定 selector を確定する（`E2E 固定 selector` 表へ追加）

`docs/screen-design/screens/term-translation-phase.md` の `E2E 固定 selector` 表へ次を追加する。persona / body の対応する画面設計正本にも同型で追加する。値は現行 frontend 実装の `data-testid` 命名に一致させて確定する。

| 対象 | `data-testid` | 関連テスト |
| --- | --- | --- |
| 処理対象件数表示 | `term-translation-phase-processing-target-total` | `E2E-LTLE-001`, `E2E-LTLE-002` |
| 処理対象一覧の空状態 | `term-translation-phase-processing-target-empty` | `E2E-LTLE-001`, `E2E-LTLE-002`, `E2E-LTLE-003` |
| 処理対象一覧の検索入力欄 | `term-translation-phase-processing-target-search-input` | `E2E-LTLE-002` |
| 初回取得中ローディングレイヤー | `term-translation-phase-processing-target-loading` | 残置 E2E の前提状態確認に使う |

persona / body の同型 selector（差分の適用先で確定する）:

| 対象 | term | persona | body |
| --- | --- | --- | --- |
| 件数表示 | `term-translation-phase-processing-target-total` | `persona-generation-phase-processing-target-total` | `body-translation-phase-processing-target-total` |
| 空状態 | `term-translation-phase-processing-target-empty` | `persona-generation-phase-processing-target-empty` | `body-translation-phase-processing-target-empty` |
| 検索入力欄 | `term-translation-phase-processing-target-search-input` | `persona-generation-phase-processing-target-search-input` | `body-translation-phase-processing-target-search-input` |
| ローディングレイヤー | `term-translation-phase-processing-target-loading` | `persona-generation-phase-processing-target-loading` | `body-translation-phase-processing-target-loading` |
| 処理対象行 | `term-translation-phase-processing-target-row.<target-id>` | `persona-generation-phase-processing-target-row.<target-id>` | `body-translation-phase-processing-target-row.<target-id>` |

selector 確定の根拠と注意:

- 注意1: 件数表示・空状態・検索入力欄の 3 つの `data-testid` は現行 frontend 実装に既に存在する。`data-testid-gaps.md` が未確定としていたのは画面設計正本「E2E 固定 selector」表への未登録を指す。本差分でこの 3 つを画面設計正本へ登録して値を確定する。
- 注意2: ローディングレイヤーの `data-testid`（`<phase-prefix>-processing-target-loading`）は本差分で新規に確定した。承認済み実装でこの値（例: `term-translation-phase-processing-target-loading`）の要素が追加済みである。Storybook 人間レビューでローディングの覆う範囲がフェーズ画面全体へ広がった後も、selector の値は処理対象一覧由来の名称のまま据え置く。据え置く根拠は、(a) 承認済み実装と残置 E2E（`E2E-LTLE-004` ほか）がこの値を参照しており、改名すると実装・E2E・`data-testid-gaps.md` への破壊的波及が生じる、(b) 実体がフェーズ画面全体オーバーレイへ広がっても、初回取得中の取得中表示という役割は同一であり、名称の変更で得られる整合上の利得がない、ことである。要素の `aria-label` は「処理対象一覧を取得中」、表示テキストは「処理対象一覧を取得中...」、`aria-busy="true"` を承認済み実装の事実として確認した。
- 注意3: 空状態の文言は現行実装の `処理対象がありません` を正とする。`data-testid-gaps.md` の `処理対象が見つかりません。` は誤記として扱い、画面設計正本へは `処理対象がありません` を残す。
- 注意4: 残置 E2E `tests/system/fix-lucien-target-list-empty.spec.ts` は `term-translation-phase-processing-target-empty` を直接参照し、page object（`tests/system/support/translation-phase-pages.ts`）経由で `-total`、`-search-input`、`-row` を参照する。本差分の 4 値はこれらの参照と一致する。

## 差分 6: 次段階開始可否と操作可否の活性・理由表示がフロント導出になる画面影響を固定する（`E-02` 進行状況兼操作）

差し戻し入力（`Q-004` 回答）に従い、次段階開始可否（次の翻訳段階を開始してよいか）と各操作可否（開始・一時停止・再開・再試行・取り消し・成果物出力確認の各活性と理由）の表示を、backend の判断結果ではなく application 層が段階データ事実状態から導出した値で表示する形へ変更する。`docs/screen-design/screens/term-translation-phase.md` の `E-02 進行状況兼操作` に次を追加・明確化する。persona / body の対応エレメントにも同型で適用する。

次段階開始可否と操作可否の表示項目（`E-02` の「表示項目」へ追加・明確化する）:

| 表示項目 | 取得元 | 表示形式 | 表示条件 | 備考 |
| --- | --- | --- | --- | --- |
| 次段階開始可否の活性 | props（application 層が段階データ事実状態から導出した可否） | 次段階を開始する操作の活性・非活性 | 常に判断できる | backend の `canStartNextPhase` を直接表示しない |
| 次段階を開始できない理由 | props（application 層が段階データ事実状態から導出した理由） | 開始できない理由の短文 | 開始できない場合 | backend の `blockedReason` を直接表示しない |
| 各操作の活性（開始・一時停止・再開・再試行・取り消し） | props（application 層が段階データ事実状態から導出した可否） | 各操作ボタンの活性・非活性 | 常に判断できる | backend の `actionEnablement` の各 `can*` を直接表示しない |
| 各操作の不可理由（開始・一時停止・再開・再試行・取り消し） | props（application 層が段階データ事実状態から導出した理由） | 各操作が行えない理由の短文 | その操作が行えない場合 | backend の各 `*BlockedReason` を直接表示しない |

規則の固定点:

- 規則14: 次段階開始可否の活性・非活性は、application 層が段階データ事実状態（フェーズ状態、対象件数、完了件数、終端状態、エラー種別）から導いた値で決める。画面は backend が返した可否判断の真偽値を直接表示の根拠にしない。
- 規則15: 次段階を開始できない理由の文言も、application 層が段階データ事実状態から導いた理由で表示する。理由の区別（ジョブが終端状態、当該段階が未完了など）は事実状態から導く。
- 規則16: 段階データ事実状態が変わった場合、活性・理由の表示も事実状態に追随して更新する。事実状態と表示の可否・理由は同じ判断で対応する。
- 規則17: 次段階開始可否の表示は term / persona / body の各段階で同型に扱う。term は次段階開始可否、persona は本文翻訳段階開始可否、body は成果物出力確認可否を、いずれも段階データ事実状態からのフロント導出値として表示する。
- 規則18: 各操作（開始・一時停止・再開・再試行・取り消し）ボタンの活性・非活性は、application 層が段階データ事実状態（フェーズ状態、ジョブの終端状態、実行設定の構成有無、対象件数、完了件数、エラー種別）から導いた値で決める。画面は backend が返した操作可否の真偽値（`actionEnablement` の各 `can*`）を直接表示の根拠にしない。
- 規則19: 各操作が行えない理由の文言も、application 層が段階データ事実状態から導いた理由で表示する。理由の区別（ジョブが終端状態、実行中フェーズが存在、ジョブが開始可能状態でない、実行設定が未構成、フェーズが実行中でない、再開可能でない、再試行可能でない）は事実状態から導く。
- 規則20: 操作可否の表示は term / persona / body の各段階で同型に扱う。body は成果物出力確認の操作（`canCheckOutputReadiness` 相当）も同じくフロント導出値で活性・理由を表示する。
- 注意5: 操作可否（開始・一時停止・再開・再試行・取り消し・成果物出力確認の活性と理由を表す `E-02` の操作行）は、`Q-004` の人間回答により、次段階開始可否と同じくフロント導出へ移る。本差分は操作行の取得元を「画面は事実状態から application 層が導いた可否・理由を表示する」へ確定する。第2回時点で残していた現状維持扱いは解消する。
- 注意6: 次段階開始可否と各操作可否の活性・理由がどの component の props 名で渡るか、導出を presenter と usecase のどちらに置くかは実装方式であり本差分では固定しない。本差分は「画面は事実状態からフロントが導いた可否・理由を表示する」という表示根拠の責務だけを固定する。

## 差分 7: body の成果物出力確認の表示が段階要約取得の事実から導出される画面影響を固定する（`body-translation-phase` の `E-02`）

差し戻し入力（`Q-005` 回答）に従い、本文翻訳段階の成果物出力確認可否・理由の表示根拠を、成果物出力確認の専用取得ではなく段階要約取得の事実から application 層が導出した値へ変更する。`docs/screen-design/screens/body-translation-phase.md` の `E-02 進行状況兼操作` に次を追加・明確化する。

成果物出力確認の表示項目（`E-02` の「表示項目」へ明確化する）:

| 表示項目 | 取得元 | 表示形式 | 表示条件 | 備考 |
| --- | --- | --- | --- | --- |
| 成果物出力確認可否の活性 | props（application 層が段階要約取得の事実から導出した可否） | 成果物出力確認操作の活性・非活性 | 常に判断できる | 専用取得 `BodyTranslationOutputReadinessResponse` を直接表示の根拠にしない |
| 成果物出力確認ができない理由 | props（application 層が段階要約取得の事実から導出した理由） | 確認できない理由の短文 | 確認できない場合 | 専用取得の `blockedReason` を直接表示の根拠にしない |

規則の固定点:

- 規則21: 成果物出力確認可否の活性・理由は、段階要約取得（`BodyTranslationPhaseSummaryResponse`）が返す事実（完了項目件数 `completedFieldCount`、状態整合 `statusConsistent`、出力件数 `outputCount` など）から application 層が導いた値で表示する。専用取得の応答を表示の根拠にしない。
- 規則22: 段階要約取得の事実が変わった場合、成果物出力確認可否の活性・理由の表示も事実に追随して更新する。
- 注意7: 成果物出力確認の事実取得が段階要約取得 1 本へ統合されること自体は取得経路の変更であり、画面の見た目（配置、文言、色）を変えない。本差分は表示根拠を段階要約取得の事実へ移すことだけを確定する。

## 表示しない領域・しない判断

- 本差分は処理対象一覧領域の見た目（配置、色、寸法）の変更を確定しない。初回取得中ローディングレイヤーの追加と、件数表示・空状態・検索欄・ローディングの固定 selector 確定だけを扱う。
- 進捗母数（`E-02`）の表示規則は変更しない。

## 根拠

- `source`: `./detail-spec-diff.md` の `translation-job-management-REQ-006`（初回表示で表示中段階だけ取得、開き直し時に再取得、旧取得結果の破棄）、`term-translation-phase-REQ-007`（初回取得完了まで検索・ページ操作を行えない、処理対象一覧の反映取りこぼし防止）。
- `source`: `./detail-spec-diff.md` の `term-translation-phase-REQ-008`、`persona-generation-phase-REQ-008`、`body-translation-phase-REQ-007`（次段階開始可否と操作可否を段階データ事実状態から application 層で導出、backend は事実だけ返す。body は成果物出力確認の事実を段階要約取得へ一本化し専用取得を廃止する）。差分 6、差分 7 の根拠。
- `source`: `docs/screen-design/screens/term-translation-phase.md` の `E-02 進行状況兼操作`（操作可否・禁止理由の表示）。
- `source`: `docs/screen-design/screens/term-translation-phase.md` の `E-04 処理対象一覧`、`E2E 固定 selector`。
- `source`: `docs/exec-plans/completed/2026-05-30-fix-lucien-target-list-empty/data-testid-gaps.md`（未登録 selector の記録、文言の誤記）。
- `source`: `frontend/src/ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte`（`processing-target-total` / `-empty` / `-search-input` / `-row` の現行命名、`busy` props）。
- `source`: `frontend/src/ui/components/ProcessingTargetListPanel.svelte`（空状態文言 `処理対象がありません`、`busy` props）。
- `source`: `tests/system/fix-lucien-target-list-empty.spec.ts`、`tests/system/support/translation-phase-pages.ts`（残置 E2E と page object が参照する selector）。
- `review`: 人間設計レビュー待ち。
- `validation`: 文書整形と現行 selector 命名との一致確認のみ。実機検証は実装後ブラウザ確認で扱う。

## 未決

- 差分 1〜5（処理対象一覧の表示規則と固定 selector）は確定済み。残る判断は実装方式であり実装範囲で扱う。
- 差分 6（次段階開始可否と操作可否のフロント導出表示）は、詳細仕様差分の未決 `Q-004` の人間回答（操作可否もフロント導出へ広げる）を反映し確定済み。`E-02` の開始・一時停止・再開・再試行・取り消しの活性・理由の取得元を、次段階開始可否と同型のフロント導出へ更新済み。
- 差分 7（body 成果物出力確認の取得経路統合）は、詳細仕様差分の未決 `Q-005` の人間回答（専用取得廃止・段階要約一本化）を反映し確定済み。
- 第4回の整合（2026-05-31）: Storybook 人間レビューで承認された frontend 実装に合わせ、初回取得中ローディングの覆う範囲を「処理対象一覧領域の最前面」から「フェーズ画面全体を覆うオーバーレイ」へ、操作排他の範囲を「処理対象一覧領域の操作」から「フェーズ画面全体の操作」へ更新した。固定 selector の値（`<phase-prefix>-processing-target-loading`）は据え置いた（注意2 の根拠）。term / persona / body の 3 段階に同型で適用する。
- テスト設計への波及（申し送り）: selector の値は据え置くため、selector に依存する観点（`E2E-LTLE-004`、`UT-LOAD-001` ほか）の参照対象は変わらない。一方で、観点の期待値文言「処理対象一覧領域の前面」は実体（フェーズ画面全体オーバーレイ）とずれる。`E2E-LTLE-004` の期待値「初回取得中はローディングレイヤーが処理対象一覧領域の前面に表示される」と、`UT-LOAD-001` 周辺の「処理対象一覧領域の前面」表現を、フェーズ画面全体オーバーレイへ整合させる文言更新が望ましい。テスト設計成果物自体の更新要否と `test_designer` への引き継ぎは呼び出し元レーンが判断する。
- 未回答の未決は 0 件。残る判断は実装方式であり実装範囲で扱う。差分本文（差分 6・差分 7 を含む）は第3回までで人間設計レビュー承認済みで、第4回は承認済み実装への整合更新である。
