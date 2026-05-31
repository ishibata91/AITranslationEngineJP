# 詳細仕様差分: job-run-phase-fetch-redesign

- `skill`: detail-spec-design
- `status`: approved-with-impl-alignment（第4回。第3回までの差分本文は人間設計レビュー承認済み。`Q-004`、`Q-005` は 2026-05-31 に人間回答済みで未回答の未決は 0 件。第4回は Storybook 人間レビューで承認された frontend 実装に合わせ、`term-translation-phase-REQ-007` の初回取得中操作排他の範囲を「処理対象操作」からフェーズ画面全体の操作へ整合させた。要件の主張（同時進行させない排他）は不変で、排他の覆う範囲だけを実装の事実へ合わせた）
- `source_plan`: `./plan.md`
- `detail_spec_target`: `docs/detail-specs/translation-job-management.md`, `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`
- `revision`: 第3回（2026-05-31 第2回設計レビューで残した未決 `Q-004`、`Q-005` の人間回答を反映）。第2回反映点は (a) 次段階開始可否判断の責務再配置、(b)「進捗の要約・次段階の準備状態・処理対象一覧の 3 つを同一連番ガードで揃える」表現の見直しである。第3回反映点は (c) 操作可否（開始・一時停止・再開・再試行・取り消し・成果物出力確認の各活性と理由）も backend から外し application 層で導出する責務再配置（`Q-004` 回答）、(d) body の成果物出力確認専用取得を廃止し段階要約取得へ一本化する取得経路統合（`Q-005` 回答）である。
- `screen_design_diff`: `./screen-design-diff.job-run.md`（後続で作成。表示規則の仕様点は本書で固定し、固定 selector の値確定は画面設計差分で扱う）
- `component_diagram`: 後続の `diagrammer` 成果物（`./design-diff.job-run-phase-fetch-redesign.md` を想定。本書時点では `N/A`）

## 用語の定義（本書で固定する日本語の対象）

- 対象: 翻訳実行画面とは、選択した翻訳ジョブの現在段階を開いて操作する画面である。`JobRunPage` が実体である。
- 対象: 翻訳段階とは、単語翻訳、NPC ペルソナ生成、本文翻訳、翻訳完了確認の各段階である。本書では `term`、`persona`、`body`、`complete` の 4 段階を同じ表示・取得の規則で扱う。
- 対象: 表示中段階とは、翻訳実行画面で利用者が今表示している 1 つの翻訳段階である。
- 対象: 処理対象一覧とは、表示中段階で処理、生成、確認する実体の一覧である。件数分の行と、その総件数を持つ。
- 対象: 進捗母数とは、表示中段階の進行状況に出る対象件数である。単語翻訳段階では AI 翻訳対象語件数（`summary.aiTargetCount` が表す件数）が該当する。
- 注意: 進捗母数と処理対象一覧の総件数は、別経路で取得・集計する独立した値である。値が一致する保証は元々ない。
- 対象: 段階データ事実状態とは、表示中段階について判断材料となる事実値の集まりである。具体的には、フェーズ状態（実行前、実行中、一時停止、完了、失敗などの状態値）、対象件数、完了件数、失敗件数、エラー種別を指す。これらは「今のデータがどうなっているか」を表す事実であり、「次へ進めてよいか」という判断を含まない。
- 対象: 次段階開始可否とは、表示中段階の段階データ事実状態から導く「次の翻訳段階を開始してよいか」という遷移判断である。単語翻訳段階では次段階の本文系処理開始可否、NPC ペルソナ生成段階では本文翻訳段階開始可否、本文翻訳段階では成果物出力確認の可否が該当する。本書では、この判断を段階データ事実状態から導出する責務をフロント（application 層）へ置く。
- 対象: 操作可否とは、表示中段階で利用者が実行できる各操作について「その操作を今行ってよいか」という判断と、行えない場合の理由である。具体的には、開始（`canStart`）、一時停止（`canPause`）、再開（`canResume`）、再試行（`canRetry`）、取り消し（`canCancel`）、本文翻訳段階の成果物出力確認（`canCheckOutputReadiness`）の各活性と、それぞれの不可理由（`startBlockedReason`、`pauseBlockedReason`、`resumeBlockedReason`、`retryBlockedReason`、`cancelBlockedReason`、`outputReadinessBlockedReason` などの理由文字列）を指す。本書では、この操作可否も段階データ事実状態から導出する責務をフロント（application 層）へ置く。次段階開始可否と同じ責務境界で扱う。
- 注意: 従来は backend が次段階開始可否（応答 DTO の `canStartNextPhase` などの真偽値、`blockedReason` などの理由文字列）と操作可否（`actionEnablement` の各 `can*` 真偽値と各 `*BlockedReason` 理由文字列）まで算出して返していた。本書はこの責務配置を見直し、backend は段階データ事実状態だけ返し、次段階開始可否と操作可否はともにフロントが導出する形へ変更する。
- 対象: 段階データ事実状態は、操作可否導出の入力としても使う。操作可否導出に使う事実は、フェーズ状態（実行前、実行中、一時停止、回復可能失敗、完了などの状態値）、ジョブが終端状態かどうか、対象件数、完了件数、失敗件数、エラー種別、実行設定が構成済みかどうか（翻訳サービス・モデル・実行方式が設定済みか）とする。これらは「今のデータと設定がどうなっているか」という事実であり、「その操作を行ってよいか」という判断を含まない。

## 詳細仕様差分

### `translation-job-management-REQ-006` 利用者向け情報を主要目的へ絞る

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/translation-job-management.md`

親要件:
利用者は未完了候補、操作可否、実行対象を主要情報として判断できる。

仕様:
- 翻訳実行画面で翻訳ジョブを開いた時、表示中段階に処理対象が実在する場合は、初回表示で処理対象一覧を件数分確認できる。
- 表示中段階に処理対象が実在しない場合だけ、初回表示で空状態を確認できる。
- 処理対象が実在するのに初回表示で空状態になる状態は、本要件を満たさない不具合として扱う。
- 取得が完了する前は、利用者は取得中であることを区別でき、取得完了後は取得した処理対象一覧か空状態のどちらかを確認できる。
- 翻訳ジョブを開いた初回表示では、表示中段階の処理対象一覧だけを取得する。表示していない他段階の処理対象一覧は、その段階を表示へ切り替えた時に取得する。翻訳ジョブを開いた直後に全段階の処理対象一覧をまとめて取得する従来の挙動は廃止する。
- 表示中段階を切り替えた時、切り替え先段階の処理対象一覧は切り替え時点で取得を開始する。切り替え直後に取得中であることを区別できる表示が一時的に出ることは許容する。
- 取得が完了して処理対象一覧を一度確認できた後は、利用者の検索操作やページ操作なしに、表示中段階の処理対象一覧が初回表示前の空状態へ戻らない。
- 表示中段階を切り替えた時も、切り替え先段階に処理対象が実在する場合は、その段階の処理対象一覧を件数分確認できる。
- 取得が遅延している間に翻訳実行画面が一度閉じられ再び開かれた場合、開き直し時に表示中段階の処理対象一覧の取得をやり直す。閉じる前に開始して遅延していた取得の結果は、開き直し後の表示へ反映せず破棄する。
- 進捗母数と処理対象一覧の総件数は、独立した取得・集計の結果として扱う。利用者は両方を別々の値として判断でき、値が一致しないことだけを理由に不具合とは判断しない。

未決:
- なし

回答:
- `Q-001`: 表示中段階だけ取得する規則へ変更する。初回表示は表示中の 1 段階だけ処理対象一覧を取得し、他段階は表示へ切り替えた時に取得する。全段階をまとめて取得する従来挙動は廃止する。段階切り替え直後に取得中表示が一時的に出ることは許容する。
  - `回答者`: 人間
  - `回答日`: 2026-05-31
  - `反映先`: 本要件の仕様「初回表示で表示中段階だけ取得」「段階切り替え時に切り替え先段階を取得」、`screen-design-diff.job-run.md`（段階切り替え時の取得中表示）。
- `Q-002`: 取得遅延中に画面が閉じて再び開かれた場合は、開き直し時に取得をやり直す。閉じる前の遅延中取得の結果は破棄する。
  - `回答者`: 人間
  - `回答日`: 2026-05-31
  - `反映先`: 本要件の仕様「開き直し時に取得をやり直し、遅延していた旧取得の結果を破棄する」。

### `term-translation-phase-REQ-007` 操作と状態を利用者が判別できる

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/term-translation-phase.md`

親要件:
利用者は開始、一時停止、再開、再試行、取り消しの可否と理由を判断できる。

仕様:
- 単語翻訳段階を初回表示した時、AI 翻訳対象語が 1 件以上ある場合は、処理対象一覧を件数分確認できる。
- AI 翻訳対象語が 0 件の場合だけ、単語翻訳段階の処理対象一覧で空状態を確認できる。
- 単語翻訳段階の進捗母数（AI 翻訳対象語件数）と処理対象一覧の総件数は独立した値である。進捗母数が 1 以上でも処理対象一覧が空状態のままになる状態は、本要件を満たさない不具合として扱う。
- 進捗の要約（進行状況の集計）と処理対象一覧は、別経路で取得する独立した結果として扱う。両者の取得・反映の新旧判断は、後述の取得・反映規則に従う。
- 処理対象一覧の取得結果は、後から起動した取得が先行取得より新しいと判断できる時だけ、新しい取得結果へ更新する。先行取得だけが完了して新しい取得結果の反映を取りこぼし、処理対象一覧が空状態のまま残る状態は許可しない。これが本不具合（初回表示で処理対象一覧だけ 0 件残留）の是正の主軸である。
- 進捗の要約は、処理対象一覧の取得・反映とは別の独立した経路で取得・反映する。進捗の要約の反映可否を理由に、処理対象一覧の反映を取りこぼしてはならない。
- 利用者は、単語翻訳段階の各操作（開始、一時停止、再開、再試行、取り消し）の可否と、行えない場合の理由を確認できる。
- 操作可否（各操作の活性と不可理由）は、段階データ事実状態から application 層で導出する。各操作の可否と理由を backend から取得する反映対象に含めない（具体的な導出条件と等価性条件は `term-translation-phase-REQ-008` で固定する）。
- 次段階開始可否は、本要件では backend から取得する反映対象に含めない。次段階開始可否は、段階データ事実状態から application 層で導出する（`term-translation-phase-REQ-008` で固定する）。
- 初回表示の取得が完了するまでは、利用者は単語翻訳段階の検索操作とページ操作を行えない。初回表示の取得中であることを区別できる表示（フェーズ画面全体を覆うローディングレイヤー）を出し、取得完了後に検索操作とページ操作を受け付ける。取得中はローディングレイヤーがフェーズ画面全体を覆うことで、検索・ページ操作だけでなくフェーズ画面全体の操作を実質受け付けない。これにより初回表示の取得と利用者操作による取得が同時に進行する状態自体を起こさせない。
- 利用者の検索操作やページ操作で処理対象一覧を取得し直す経路は、初回表示の取得完了後に行う別の取得として扱い、利用者が要求した検索条件とページ位置の結果を確認できる。

未決:
- なし

回答:
- `Q-003`: 初回表示の取得中は、取得中であることを区別できる表示（ローディングレイヤー）を上に出して実質操作不可にし、利用者の検索操作とページ操作を受け付けない。これにより初回表示の取得と利用者操作による取得を同時進行させず、競合自体を起こさせない。取得完了後に利用者操作を受け付ける。よって「同時進行時の優先規則」は不要となり、「同時進行させない（排他）」で確定する。
  - `回答者`: 人間
  - `回答日`: 2026-05-31
  - `反映先`: 本要件の仕様「初回表示の取得完了まで検索・ページ操作を行えない」「初回表示の取得完了後に利用者操作を受け付ける」、`screen-design-diff.job-run.md`（初回表示の取得中ローディングレイヤーの表示規則と操作排他範囲）。
  - `実装整合（2026-05-31）`: Storybook 人間レビューで承認された frontend 実装は、`Q-003` の「ロード中レイヤーを上に出して実質操作不可にする」意図に沿い、ローディングレイヤーをフェーズ画面全体を覆うオーバーレイにした。これにより排他範囲は「処理対象操作（検索・ページ操作）」からフェーズ画面全体の操作へ広がった。検索・ページ操作の取得起動を起こさせない仕様（同時進行させない排他）は満たしたままで、覆う範囲だけが画面全体へ広がった整合である。詳細は `screen-design-diff.job-run.md` 差分 1・差分 2 を正とする。term / persona / body 同型。
- `Q-004`: 操作可否（開始・一時停止・再開・再試行・取り消しの各活性と不可理由）も backend から外し application 層で導出する。次段階開始可否と同じ責務境界で扱う。
  - `回答者`: 人間
  - `回答日`: 2026-05-31
  - `反映先`: 本要件の仕様「操作可否を段階データ事実状態から application 層で導出する」、`term-translation-phase-REQ-008` の責務境界と等価性条件、`persona-generation-phase-REQ-007` / `body-translation-phase-REQ-006` の同型仕様。

### `term-translation-phase-REQ-008` 次段階開始可否と操作可否を段階データ事実状態から判断できる

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/term-translation-phase.md`

親要件:
利用者は次の翻訳段階を開始してよいかを判断できる。利用者は開始、一時停止、再開、再試行、取り消しの可否と理由を判断できる。

仕様:
- 利用者は、単語翻訳段階の次段階開始可否（次段階を開始してよいか）と、開始できない場合の理由を確認できる。
- 次段階開始可否は、段階データ事実状態から導く。導出に使う事実は、単語翻訳段階のフェーズ状態（完了かどうか）、ジョブが終端状態かどうか、対象件数、確認済み件数、エラー種別とする。
- 次段階を開始してよいと判断できる成立条件は、ジョブが終端状態でなく、単語翻訳段階のフェーズ状態が完了であり、確認済み件数が対象件数以上であることとする。
- 上記の成立条件を満たさない場合は次段階を開始できない。開始できない理由は、ジョブが終端状態である、単語翻訳段階が未完了である、のいずれかを段階データ事実状態から区別して示す。
- 利用者は、各操作（開始、一時停止、再開、再試行、取り消し）の可否と、行えない場合の理由を確認できる。各操作の可否と理由も、段階データ事実状態から application 層で導出する。導出に使う事実は、フェーズ状態（実行前、実行中、一時停止、回復可能失敗、完了などの状態値）、ジョブが終端状態かどうか、実行設定が構成済みかどうか（翻訳サービス・モデル・実行方式が設定済みか）、対象件数、確認済み件数、エラー種別とする。
- 各操作可否の不可理由は、段階データ事実状態から区別して示す。区別する理由は、ジョブが終端状態である、実行中フェーズが存在する、ジョブが開始可能状態でない、実行設定が未構成である、フェーズが実行中でない、フェーズが再開可能でない、フェーズが再試行可能でない、のいずれかとする。
- 段階データ事実状態を返す責務は backend に置く。次段階開始可否と各操作可否の真偽と理由を導出する責務は application 層に置く。backend は次段階開始可否の真偽値（`canStartNextPhase` など）と理由文字列（`blockedReason` など）、および操作可否の各真偽値（`canStart`、`canPause`、`canResume`、`canRetry`、`canCancel` など）と各理由文字列（`startBlockedReason`、`pauseBlockedReason`、`resumeBlockedReason`、`retryBlockedReason`、`cancelBlockedReason` など）を応答に含めない。
- 等価性条件: backend が返す段階データ事実状態から application 層が導出する次段階開始可否および各操作可否は、再配置前に backend が返していた可否・理由と、同じ入力事実に対して同じ結果になる。次段階開始可否の真偽・理由、各操作（開始・一時停止・再開・再試行・取り消し）の真偽・理由のいずれについても、再配置前後で同じ事実入力に同じ結果が導けることを満たす。

未決:
- なし

回答:
- `Q-004`: 操作可否（開始・一時停止・再開・再試行・取り消しの各活性と不可理由）も backend から外し application 層で導出する。backend は段階の事実状態だけ返す。各操作の可否・理由は事実状態から application 層が導出する。再配置前に backend が返していた可否・理由と、再配置後に application 層が導出する可否・理由は、同じ事実入力に対して同じ結果になることを等価性条件として固定する。term/persona/body 同型で扱う。
  - `回答者`: 人間
  - `回答日`: 2026-05-31
  - `反映先`: 本要件の責務境界と等価性条件、`persona-generation-phase-REQ-008`、`body-translation-phase-REQ-007` の同型仕様、実装範囲。

### `persona-generation-phase-REQ-007` 操作と状態を利用者が判別できる

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/persona-generation-phase.md`

親要件:
利用者はフェーズの進行、生成対象、生成結果、ペルソナ参照状態、本文翻訳フェーズの開始可否を判断できる。

仕様:
- NPC ペルソナ生成段階を表示した時、生成対象が 1 件以上ある場合は、処理対象一覧を件数分確認できる。
- 生成対象が 0 件の場合だけ、処理対象一覧で空状態を確認できる。
- 処理対象一覧の取得結果は、後から起動した取得が先行取得より新しいと判断できる時だけ更新する。先行取得だけが完了して処理対象一覧の反映を取りこぼし、空状態のまま残る状態は許可しない。
- 進捗の要約は、処理対象一覧の取得・反映とは別の独立した経路で取得・反映する。進捗の要約の反映可否を理由に、処理対象一覧の反映を取りこぼしてはならない。
- 進捗母数と処理対象一覧の総件数は独立した値として扱う。
- 利用者は、NPC ペルソナ生成段階の各操作（開始、一時停止、再開、再試行、取り消し）の可否と、行えない場合の理由を確認できる。操作可否は段階データ事実状態から application 層で導出する。各操作の可否と理由を backend から取得する反映対象に含めない（具体的な導出と等価性条件は `persona-generation-phase-REQ-008` で固定する）。
- 本文翻訳段階開始可否は、本要件では backend から取得する反映対象に含めない。本文翻訳段階開始可否は段階データ事実状態から application 層で導出する（`persona-generation-phase-REQ-008` で固定する）。

未決:
- なし

回答:
- なし

### `persona-generation-phase-REQ-008` 本文翻訳段階開始可否と操作可否を段階データ事実状態から判断できる

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/persona-generation-phase.md`

親要件:
利用者は本文翻訳フェーズの開始可否を判断できる。利用者は開始、一時停止、再開、再試行、取り消しの可否と理由を判断できる。

仕様:
- 利用者は、NPC ペルソナ生成段階の本文翻訳段階開始可否と、開始できない場合の理由を確認できる。
- 本文翻訳段階開始可否は、NPC ペルソナ生成段階の段階データ事実状態（生成対象件数、生成完了件数、ペルソナ参照状態などの事実値）から application 層で導出する。
- 利用者は、NPC ペルソナ生成段階の各操作（開始、一時停止、再開、再試行、取り消し）の可否と、行えない場合の理由を確認できる。各操作の可否と理由も、段階データ事実状態（フェーズ状態、ジョブの終端状態、実行設定の構成有無、生成対象件数、生成完了件数、エラー種別などの事実値）から application 層で導出する。
- 段階データ事実状態を返す責務は backend に置く。本文翻訳段階開始可否と各操作可否の真偽と理由を導出する責務は application 層に置く。backend は本文翻訳段階開始可否の真偽値（`bodyReadiness` など）と理由文字列（`blockedReason` など）、および操作可否の各真偽値（`canStart`、`canPause`、`canResume`、`canRetry`、`canCancel` など）と各理由文字列を応答に含めない。
- 等価性条件: backend が返す段階データ事実状態（生成対象件数や生成完了件数を含む `inputSummary` 相当の事実値）から、application 層が再配置前と同じ本文翻訳段階開始可否および各操作可否を導けることを満たす。本文翻訳段階開始可否、各操作の可否・理由のいずれについても、再配置前後で同じ事実入力に同じ結果が導けることを満たす。

未決:
- なし

回答:
- `Q-004`: 操作可否も backend から外し application 層で導出する。NPC ペルソナ生成段階の各操作の可否・理由を、本文翻訳段階開始可否と同じ責務境界で事実状態から導出する。term/persona/body 同型で扱う。等価性条件は本要件の仕様で固定する。
  - `回答者`: 人間
  - `回答日`: 2026-05-31
  - `反映先`: 本要件の責務境界と等価性条件、`term-translation-phase-REQ-008` / `body-translation-phase-REQ-007` の同型仕様、実装範囲。

### `body-translation-phase-REQ-006` 進行と結果を判断できる

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/body-translation-phase.md`

親要件:
利用者は本文翻訳フェーズの進行、結果、失敗理由、成果物出力条件を判断できる。

仕様:
- 本文翻訳段階を表示した時、翻訳項目が 1 件以上ある場合は、処理対象一覧を件数分確認できる。
- 翻訳項目が 0 件の場合だけ、処理対象一覧で空状態を確認できる。
- 処理対象一覧の取得結果は、後から起動した取得が先行取得より新しいと判断できる時だけ更新する。先行取得だけが完了して処理対象一覧の反映を取りこぼし、空状態のまま残る状態は許可しない。
- 進捗の要約は、処理対象一覧の取得・反映とは別の独立した経路で取得・反映する。進捗の要約の反映可否を理由に、処理対象一覧の反映を取りこぼしてはならない。
- 本文翻訳段階の処理対象一覧と翻訳完了確認の処理対象一覧は、表示中段階ごとに独立した取得・表示として扱う。
- 利用者は、本文翻訳段階の各操作（開始、一時停止、再開、再試行、取り消し、成果物出力確認）の可否と、行えない場合の理由を確認できる。操作可否は段階データ事実状態から application 層で導出する。各操作の可否と理由を backend から取得する反映対象に含めない（具体的な導出と等価性条件は `body-translation-phase-REQ-007` で固定する）。
- 成果物出力確認可否は、本要件では backend から取得する反映対象に含めない。成果物出力確認可否は段階データ事実状態から application 層で導出する（`body-translation-phase-REQ-007` で固定する）。

未決:
- なし

回答:
- なし

### `body-translation-phase-REQ-007` 成果物出力確認可否と操作可否を段階データ事実状態から判断でき、取得経路を段階要約に一本化する

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/body-translation-phase.md`

親要件:
利用者は成果物出力確認の可否を判断できる。利用者は開始、一時停止、再開、再試行、取り消しの可否と理由を判断できる。

仕様:
- 利用者は、本文翻訳段階の成果物出力確認可否と、確認できない場合の理由を確認できる。
- 成果物出力確認可否は、本文翻訳段階の段階データ事実状態（翻訳項目件数、完了項目件数 `completedFieldCount`、失敗件数、状態の整合 `statusConsistent`、出力件数 `outputCount`、エラー種別などの事実値）から application 層で導出する。
- 利用者は、本文翻訳段階の各操作（開始、一時停止、再開、再試行、取り消し、成果物出力確認）の可否と、行えない場合の理由を確認できる。各操作の可否と理由も、段階データ事実状態（フェーズ状態、ジョブの終端状態、実行設定の構成有無、翻訳項目件数、完了項目件数、失敗件数、状態整合、出力件数、エラー種別などの事実値）から application 層で導出する。
- 段階データ事実状態を返す責務は backend に置く。成果物出力確認可否と各操作可否の真偽と理由を導出する責務は application 層に置く。backend は成果物出力確認可否の真偽値（`canCheckOutputReadiness` や `outputReadiness.ready` など）と理由文字列（`blockedReason` など）、および操作可否の各真偽値（`canStart`、`canPause`、`canResume`、`canRetry`、`canCancel` など）と各理由文字列を応答に含めない。
- 成果物出力確認に必要な段階データ事実状態は、本文翻訳段階の段階要約取得（`BodyTranslationPhaseSummaryResponse`）に集約する。集約する事実値は、完了項目件数 `completedFieldCount`、状態整合 `statusConsistent`、出力件数 `outputCount` を含む。これらの事実値が段階要約に含まれていない場合は段階要約へ加える。
- 成果物出力確認可否を取得するための専用取得（`GetBodyTranslationOutputReadiness` と応答 `BodyTranslationOutputReadinessResponse`）は廃止する。成果物出力確認可否の事実入力は段階要約取得の 1 本へ統合し、出力可否（`ready`）と不可理由（`blockedReason`）は段階要約の事実から application 層が導出する。これにより成果物出力確認可否の取得経路を、段階要約取得と専用取得の 2 本から、段階要約取得の 1 本へ統合する。
- 等価性条件: backend が返す段階データ事実状態（完了項目件数 `completedFieldCount`、状態整合 `statusConsistent`、出力件数 `outputCount` を含む事実値）から、application 層が再配置前と同じ成果物出力確認可否および各操作可否を導けることを満たす。再配置前に backend が返していた成果物出力確認可否・理由、各操作の可否・理由のいずれについても、専用取得を廃止して段階要約取得へ一本化した後に、同じ事実入力に同じ結果が導けることを満たす。
- 注意: 成果物の最終出力可否を表す出力成果物側の判定（`translation-output-artifact` の `outputReadiness`）は、本要件の対象外とする。本要件は本文翻訳段階の画面で示す成果物出力確認可否に限定する。

未決:
- なし

回答:
- `Q-004`: 操作可否も backend から外し application 層で導出する。本文翻訳段階の各操作（成果物出力確認 `canCheckOutputReadiness` を含む）の可否・理由を、成果物出力確認可否と同じ責務境界で事実状態から導出する。term/persona/body 同型で扱う。等価性条件は本要件の仕様で固定する。
  - `回答者`: 人間
  - `回答日`: 2026-05-31
  - `反映先`: 本要件の責務境界と等価性条件、`term-translation-phase-REQ-008` / `persona-generation-phase-REQ-008` の同型仕様、実装範囲。
- `Q-005`: body の成果物出力確認専用取得（`GetBodyTranslationOutputReadiness` と `BodyTranslationOutputReadinessResponse`）を廃止する。成果物出力確認に必要な事実（`completedFieldCount`、`statusConsistent`、`outputCount` など）は段階要約取得（`BodyTranslationPhaseSummaryResponse`）に集約する。`outputCount` が段階要約に無ければ段階要約へ加える。出力可否（`ready`）と `blockedReason` は backend から外し、frontend が段階要約の事実から導出する。取得経路を 2 本から 1 本へ統合する。
  - `回答者`: 人間
  - `回答日`: 2026-05-31
  - `反映先`: 本要件の取得経路と仕様、実装範囲、設計差分図。

## 根拠

- `source`: `docs/spec.md` の翻訳ジョブ実行進捗、翻訳フロー、各フェーズ、翻訳ジョブ状態。
- `source`: `docs/detail-specs/translation-job-management.md` の `translation-job-management-REQ-006`（処理対象一覧と件数の判断、既定ページサイズ 50 件）。
- `source`: `docs/detail-specs/term-translation-phase.md` の `term-translation-phase-REQ-003`（AI 翻訳対象語 0 件の扱い）、`REQ-006`（対象語件数・AI 翻訳対象語件数の判断）、`REQ-007`（処理対象名の判断）。
- `source`: `docs/detail-specs/persona-generation-phase.md` の `persona-generation-phase-REQ-007`、`docs/detail-specs/body-translation-phase.md` の `body-translation-phase-REQ-006`。
- `source`: `docs/exec-plans/completed/translation-job-step-target-list-panel/detail-spec-diff.md`（処理対象一覧表示パネルの既存仕様差分）。
- `source`: `docs/exec-plans/active/job-run-phase-fetch-redesign/plan.md`（作り直し方針、症状、確定事実）。
- `source`: `docs/exec-plans/completed/2026-05-30-fix-lucien-target-list-empty/fix-decision.md`（追補1〜5、obs-r4〜r6 の実機観測）。
- `source`: `docs/exec-plans/completed/2026-05-30-fix-lucien-target-list-empty/uc-diff.md`（UC 差分なし、既存 UC「処理対象を確認する」で説明可能）。
- `source`: `docs/exec-plans/completed/2026-05-30-fix-lucien-target-list-empty/data-testid-gaps.md`（固定 selector 未確定は画面設計差分で扱う）。
- `source`: `frontend/src/application/usecase/term-translation-phase/term-translation-phase.usecase.ts` の `fetchSummaryAndReadiness`（処理対象一覧だけ連番ガード内反映の非対称）、`frontend/src/ui/screens/job-run/JobRunPage.svelte` の `$effect`・`onMount`・`currentProcessingTargetPageState`。
- `source`: `internal/service/term_translation_phase_service.go` の `readinessFromState`（段階データ事実状態から `canStartNextPhase` と `blockedReason` を backend で導出している現状）。`internal/usecase/term_translation_phase_usecase.go` の `GetTermTranslationNextPhaseReadiness`、`internal/controller/wails/term_translation_phase_controller.go` の `TermTranslationNextPhaseReadinessResponseDTO`（`canStartNextPhase`、`blockedReason`）。
- `source`: persona / body の同型構造。`internal/controller/wails/persona_generation_phase_controller.go` の `PersonaGenerationBodyReadinessResponseDTO`（`bodyReadiness`、`blockedReason`、`inputSummary`）。`internal/controller/wails/body_translation_phase_controller.go` の `BodyTranslationOutputReadinessResponseDTO` と `canCheckOutputReadiness`、`outputReadiness`。
- `source`: フロントの readiness 利用箇所。`frontend/src/application/presenter/body-translation-phase/body-translation-phase.presenter.ts` の `getEffectiveReadiness`、各 presenter / screen-types の `nextPhaseReadiness`・`bodyReadiness`・`outputReadiness`。
- `source`: `frontend/src/application/gateway-contract/term-translation-phase/term-translation-phase-gateway-contract.ts:139` の `TermTranslationNextPhaseReadinessResponse`（差し戻しで責務再配置の対象に指定された応答）。
- `review`: 未決 `Q-001`、`Q-002`、`Q-003` は人間が 2026-05-31 に回答済み。第2回差し戻しで追加した未決 `Q-004`（操作可否まで責務再配置を広げるか）、`Q-005`（成果物出力確認可否の取得経路の統廃合）も 2026-05-31 に人間回答済み。`Q-004` は操作可否もフロント導出へ広げる回答、`Q-005` は body 専用取得を廃止し段階要約取得へ一本化する回答である。未回答の未決は 0 件。差分本文は第3回の人間設計レビュー待ち。
- `source`: body 段階要約の出力件数。`frontend/src/application/gateway-contract/body-translation-phase/body-translation-phase-gateway-contract.ts` の `BodyTranslationPhaseFieldResultSummary.outputCount`（段階要約の `resultSummary` 内に既に任意項目として存在する出力件数）、`BodyTranslationPhaseActionEnablement`（`canStart`/`canPause`/`canResume`/`canRetry`/`canCancel`/`canCheckOutputReadiness` と各 `*BlockedReason`）、`BodyTranslationOutputReadinessResponse`（廃止対象の専用取得応答）。
- `source`: 操作可否の現状導出。`internal/service/term_translation_phase_service.go` の `termTranslationStartBlockedReason`、`termTranslationPauseBlockedReason`、`termTranslationResumeBlockedReason`、`termTranslationRetryBlockedReason`（実行設定の構成有無、フェーズ状態、ジョブ終端から各操作可否・理由を backend で導出している現状）。persona / body の同型 service も同じ構造で操作可否を導出する。
- `validation`: 文書整形確認と backend / frontend の readiness 算出箇所の存在確認のみ。実機検証は実装後ブラウザ確認で扱う。

## 設計上の注記（仕様にしない実装観点。実装範囲・図・画面設計差分へ引き継ぐ）

- 注記1: 約15秒の取得遅延は IPC 飽和の疑いであり真因未確定である。本書は「取得本数を必要分へ絞る」「反映基準を揃える」という仕様で取りこぼしを防ぐが、遅延の真因分離は観測ログ追加で行う前提を残す。観測ログの要否は設計確定後に呼び出し元レーンが再判定する。
- 注記2: 取得起動を `$effect` と `onMount` のどちらに置くか、`Promise.all` の同時起動本数、直列化の有無は実装方式であり本書では固定しない。実装範囲と設計差分図で扱う。
- 注記3: 「進捗の要約・次段階の準備状態・処理対象一覧の 3 つを同一連番ガードで揃える」という第1回の表現を見直した。3 本は粒度が異なる。処理対象一覧はページング付きの重い取得であり、進捗の要約は進行状況の集計であり、次段階開始可否は事実からの判断であって取得物ではない。よって 3 本を機械的に同一連番ガードで束ねる表現を改めた。本書の仕様は次の 2 点を主軸にする。第 1 に、処理対象一覧の反映を取りこぼさないこと（先行取得だけ完了して空状態が残る非対称の是正）。第 2 に、進捗の要約は処理対象一覧と別経路で独立に反映すること。次段階開始可否は backend からの取得本数ではなくなり、段階データ事実状態から application 層が導出する。連番ガード（`processingTargetListRequestSequence`）は手動操作の取得競合防止に必要であり撤廃しない。具体的な連番管理と取得・反映の実装方式は実装範囲・設計差分図で扱う。
- 注記4: 固定 selector（処理対象件数、空状態、検索入力欄、初回取得中ローディングレイヤー）の値確定は画面設計差分で扱う。本書は表示規則（初回件数分表示、母数0で空状態保持、検索後リロードで0件へ戻らない、初回取得中は操作排他）の仕様だけを固定する。
- 注記5: 初回表示の取得中ローディングレイヤーの見た目、配置、文言、排他の実装方式（オーバーレイ要素の前面化、入力の無効化方式）は画面設計差分と実装範囲で扱う。本書は「初回取得完了まで検索・ページ操作を行えない」「同時進行させない」という仕様だけを固定する。
- 注記6: 次段階開始可否と操作可否のフロント導出（`term-translation-phase-REQ-008`、`persona-generation-phase-REQ-008`、`body-translation-phase-REQ-007`）について、backend が応答 DTO から外す対象（`canStartNextPhase`、`blockedReason`、`bodyReadiness`、`canCheckOutputReadiness`、`actionEnablement` の各 `can*` と `*BlockedReason` など）と、その代わりに返す段階データ事実値の具体的な列・キー、application 層のどこ（presenter か usecase か）で導出するか、専用取得（`GetTermTranslationNextPhaseReadiness` など）の存廃は実装方式である。本書は「backend は事実だけ返す」「次段階開始可否と操作可否はフロント導出」「変更前後で同じ事実から同じ判断が導ける」という責務境界と等価性条件だけを固定する。応答 DTO の具体形、gateway-contract の型変更、導出関数の配置は実装範囲・設計差分図・公開契約で扱う。backend の操作可否導出関数群（`termTranslationStartBlockedReason` などと persona / body 同型）の application 層への移設と backend テストへの波及は、実装範囲で波及範囲を確定する。
- 注記7: body の成果物出力確認専用取得廃止（`body-translation-phase-REQ-007`、`Q-005` 回答）について、`GetBodyTranslationOutputReadiness` と `BodyTranslationOutputReadinessResponse` の削除手順、段階要約への `outputCount` 集約方式（既存 `resultSummary.outputCount` を使うか段階要約直下へ移すか）、専用取得を参照している frontend gateway / presenter / screen-types の経路統合は実装方式である。本書は「成果物出力確認の事実入力を段階要約取得 1 本へ統合する」「出力可否はフロント導出」「統合前後で同じ事実から同じ可否が導ける」という取得経路と等価性条件だけを固定する。削除に伴う backend controller / usecase / service の波及範囲は実装範囲で確定する。
