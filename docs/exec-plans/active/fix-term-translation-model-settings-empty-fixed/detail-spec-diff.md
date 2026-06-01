# 詳細仕様差分: fix-term-translation-model-settings-empty-fixed

- `skill`: detail-spec-design
- `status`: ready-for-human-review
- `source_plan`: `./plan.md`
- `detail_spec_target`: `docs/detail-specs/term-translation-phase.md`, `docs/detail-specs/persona-generation-phase.md`, `docs/detail-specs/body-translation-phase.md`, `docs/er.md`, `docs/diagrams/er/combined-data-model-er.puml`
- `screen_design_diff`: N/A
- `component_diagram`: `./design-diagram.md`（後続で固定）

## 方向比較と採用案

### 比較対象

| 方向 | 内容 | Ready 期の AI 設定保持 | 実行中の固定値保持 | 「設定なし」表現 |
| --- | --- | --- | --- | --- |
| A | snapshot テーブルを ER 正本へ正式追加し、Ready 期の保存値と実行時固定値を同居させる | snapshot 行 | 同じ snapshot 行 | snapshot 行に空フィールドを許容 |
| B | Ready 期の AI 設定を保持する 3 フェーズ種別共通の独立テーブル `JOB_PHASE_AI_SETTINGS`（主キー `phase_type` のみ）を ER 正本に新設し、実行中固定値は `JOB_PHASE_RUN` の既存列で持つ | 独立テーブル（`phase_type` ベース、ジョブと無関連） | `JOB_PHASE_RUN.ai_provider` ほか既存列 | 該当 row 不在 |
| C | 専用テーブルを廃止し、provider-settings 正本から都度解決し、フェーズ開始時に `JOB_PHASE_RUN` を作成して固定する | provider-settings 正本（フェーズ単位の選択保持なし） | `JOB_PHASE_RUN` 既存列 | フェーズ単位 record が存在しない |

### 採用案: B

採用理由:
- 人間レビューが「snapshot は『動いてるぶん』だけ、設定は別の共有ストアで保持」「『設定なし』状態は record 不在で表現」を明示している。Ready 期の AI 設定保持場所を独立 record として置けば、record の有無で「設定済み／未設定」を判別でき、空文字フィールドを「未設定」の代理表現にする経路を作らずに済む。
- `docs/er.md:67-68`「`Ready` job には `JOB_PHASE_RUN` を事前作成しない」を維持できる。`JOB_PHASE_RUN` は ER 正本のとおり「フェーズが開始許可された時だけ作成する実行情報」だけを保持し、Ready 期の保存値を持たない。
- `docs/er.md:26`「フェーズ別 AI 設定、指示構成、最終 AI 実行情報は `JOB_PHASE_RUN` に保持する」は、フェーズ実行中の固定値を保持する規約として維持する。Ready 期の保存値は別 record で扱い、フェーズ開始時に `JOB_PHASE_RUN` の対応列へ転写する。
- persona-generation / body-translation も同じ Ready 期 AI 設定の保持を必要としているため、独立テーブルを共通設計として 1 箇所に固定できる。

不採用理由:
- 方向 A: Ready 期保存値と実行時固定値を同居させる構造を維持する。人間レビューが拒否した「空文字 record で『未設定』を表す」構造を残し、`applyTermTranslationRuntimeSnapshot` 相当の空文字上書きの誘因が消えない。
- 方向 C: Ready 期に「フェーズ単位で利用者が選択した AI 設定」を保持できない。`docs/screen-design/screens/term-translation-phase.md:102-104` で「AI サービス / モデル / 処理方式を Ready 期に選択できる」が定義済みで、フェーズ単位の選択結果が `Ready` の間も残っている必要があるため、provider-settings 正本（利用者単位の登録）だけでは満たせない。

## 詳細仕様差分

### `er-REQ-001` フェーズ別 AI 設定の Ready 期保存値を `JOB_PHASE_AI_SETTINGS` で保持する

- `変更種別`: 追加
- `要件扱い`: 追加要件
- `正本反映先`: `docs/er.md`, `docs/diagrams/er/combined-data-model-er.puml`

親要件:
利用者は、フェーズ種別ごとに、フェーズ単位の AI 選択値（AI サービス、モデル、処理方式）を事前に保存し、保存済みかどうかを record の有無で判別できる。保存値は翻訳ジョブを横断して使い回せる。

仕様:
- フェーズ別 AI 設定の Ready 期保存値は、`JOB_PHASE_RUN` とは別の独立テーブル `JOB_PHASE_AI_SETTINGS` で保持する。
- `JOB_PHASE_AI_SETTINGS` は 3 フェーズ種別（単語翻訳、ペルソナ生成、本文翻訳）共通の汎用テーブルとして 1 枚で定義し、各 record は `phase_type` で区別する。
- `JOB_PHASE_AI_SETTINGS` の主キーは `phase_type` のみとし、ジョブとは関連を持たない。`job_id` 列、`user_id` 列、ジョブへの外部キーは持たない。
- `JOB_PHASE_AI_SETTINGS` には 3 フェーズ種別の最大 3 件の record だけが存在し得る。
- 同一 `phase_type` の record が存在しない状態は「利用者がそのフェーズ種別の AI 選択値を保存していない状態」を表す。
- 同一 `phase_type` の record が存在する状態は「利用者がそのフェーズ種別の AI 選択値を保存済みの状態」を表し、保存値は以後のすべての翻訳ジョブで使い回せる。
- `JOB_PHASE_AI_SETTINGS` の record は、AI サービス、モデル、処理方式の 3 値を保持する。認証参照（credential 識別子）は保持しない。
- 認証参照、認証状態、利用可能モデル一覧は AIサービス設定（provider-settings 正本）の責務であり、`JOB_PHASE_AI_SETTINGS` は保持しない。
- 値の変更経路は上書き（upsert）に限る。明示的な削除 API は持たない。
- ジョブの削除や終端遷移に連動する cascade は持たない。`JOB_PHASE_AI_SETTINGS` の record はジョブの生存期間と独立に維持する。

未決:
- なし

回答:
- `Q-001`: 回答済み
  - `回答`: 3 フェーズ種別（単語翻訳、ペルソナ生成、本文翻訳）すべてを同時に独立 record 設計へ移行する。`JOB_PHASE_AI_SETTINGS` は 3 種別共通の汎用テーブルとし、phase_type で区別する。本 task の実装範囲も 3 フェーズ同時とする。本回答により `fix-decision.md` の禁止修正 5（影響範囲を単語翻訳に限定する規約）は覆る。`fix-decision.md` 自体の書き換えは investigation-module 範囲のため本差分では指示しない。
  - `回答者`: 人間レビュー
  - `回答日`: 2026-06-01
  - `反映先`: `er-REQ-001`（3 種別共通テーブルとして定義）、波及範囲表（3 フェーズ同時実装）、影響する正本ファイル表（detail-specs 3 件）
- `Q-002`: 回答済み
  - `回答`: テーブルの正式名は `JOB_PHASE_AI_SETTINGS` とする。仮称ではなく確定名として採用する。
  - `回答者`: 人間レビュー
  - `回答日`: 2026-06-01
  - `反映先`: `er-REQ-001`、本差分全体（仮称表現を `JOB_PHASE_AI_SETTINGS` に統一）

### `er-REQ-002` `JOB_PHASE_RUN` の AI 設定列の責務を実行中固定値に限定する

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/er.md`

親要件:
利用者は、フェーズ実行中に固定された AI 接続情報（AI サービス、モデル、処理方式、認証参照）を、フェーズ実行情報の一部として判断できる。

仕様:
- `JOB_PHASE_RUN` の AI サービス、モデル、処理方式、認証参照は、フェーズ開始が許可されて `JOB_PHASE_RUN` を作成する時点で確定し、以後同じ `JOB_PHASE_RUN` の継続中は維持する。
- `JOB_PHASE_RUN.credential_ref` は「フェーズ開始時に AIサービス設定（provider-settings 正本）から解決した認証参照を実行時固定値として記録する」責務に限定する。Ready 期の利用者保存値の保持には使わない。
- `JOB_PHASE_RUN` の AI サービス、モデル、処理方式は、フェーズ開始時に同じ `phase_type` を持つ `JOB_PHASE_AI_SETTINGS` record の選択値から転写する。`JOB_PHASE_AI_SETTINGS` record が該当 `phase_type` に対して存在しない場合はフェーズ開始が成立しないため、転写も発生しない。
- `JOB_PHASE_RUN` の AI 設定列は、Ready 期の利用者保存値の保持には使わない。
- `JOB_PHASE_RUN` は引き続きジョブ単位で扱い、Ready 期には作成せず、フェーズ開始時にだけ作成する規約を維持する。
- フェーズ再実行は同じ `JOB_PHASE_RUN` を継続する規約を維持し、AI 設定列の値も継続する。

未決:
- なし

回答:
- なし

### `term-translation-phase-REQ-002` AI 設定を開始時に再解決する（変更）

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/term-translation-phase.md`

親要件:
フェーズ開始と再試行は最新の AIサービス設定を参照し、秘密値を利用者向け情報から分離する。

仕様（既存仕様への追加・変更分のみ列挙）:
- 単語翻訳フェーズの Ready 期 AI 設定の保存、参照、変更は、`JOB_PHASE_AI_SETTINGS` の `phase_type = word_translation` record として扱う。
- `JOB_PHASE_AI_SETTINGS` record は AI サービス、モデル、処理方式の 3 値を保持する。認証参照、認証状態、利用可能モデル一覧は AIサービス設定（provider-settings 正本）の責務であり、`JOB_PHASE_AI_SETTINGS` は保持しない。
- 単語翻訳フェーズ用の `JOB_PHASE_AI_SETTINGS` record（`phase_type = word_translation`）が存在しない状態は「AI 設定不足」として判断できる。
- 単語翻訳フェーズ用 record が存在し、保持する AI サービス、モデル、処理方式の 3 値と、AIサービス設定（provider-settings 正本）から都度解決した認証参照の 1 値のすべてが解決済みである状態は「AI 設定準備済み」として判断できる。
- 単語翻訳フェーズ用 record が存在し、保持値または provider-settings 正本側の解決結果のいずれかが解決できない状態は「AI 設定不足」として判断できる。
- 単語翻訳フェーズ開始時は、`phase_type = word_translation` の `JOB_PHASE_AI_SETTINGS` record の 3 値と、AIサービス設定の最新解決結果（認証参照）を合わせて、実行時固定値を `JOB_PHASE_RUN` の AI 設定列へ転写する。
- 単語翻訳フェーズ実行中の表示は、`JOB_PHASE_RUN` の AI 設定列を判断根拠にする。
- 単語翻訳フェーズ未開始（`JOB_PHASE_RUN` が未作成）の表示は、`JOB_PHASE_AI_SETTINGS` record と AIサービス設定の解決結果を判断根拠にする。
- `JOB_PHASE_AI_SETTINGS` record の不在を、空値の record で代理表現しない。

AI 設定保存操作（`SaveAISettings` 相当）の入力構造に関する仕様:
- 入力は `phase_type` と AI 選択値（AI サービス、モデル、処理方式、処理方式の補助値 `batchMode`）のみで構成する。
- 入力に `job_id` を含めない。保存操作はジョブと無関連で、`JOB_PHASE_AI_SETTINGS` を `phase_type` を主キーとして upsert する。
- 同一 `phase_type` への再保存は上書きとして処理する。明示的な削除入力は持たない。

backend 応答の構造に関する仕様:
- 単語翻訳フェーズの状態を返す backend 応答は、フェーズ単位の「設定全体の情報」を返すか、「該当状態が存在しないこと」を示すかの二択でだけ表現する。
- 「設定不足の理由文字列」「blocked 理由文字列」など、状態判断から派生する説明文字列を backend 応答に含めない。状態判断の表現は frontend 表示の責務とする。
- 単語翻訳フェーズ用 `JOB_PHASE_AI_SETTINGS` record が存在しない状態は、応答の AI 設定 field を不在（field を含めない、または null）にすることで表現する。空文字値を持つ AI 設定 field で代理表現しない。
- 実行情報（`execution` 相当 field）は、`JOB_PHASE_RUN` が存在する場合だけ応答へ含める。`JOB_PHASE_RUN` が未作成の状態を、空値を持つ `execution` field で代理表現しない。
- frontend は、応答の AI 設定 field の有無で「AI 設定準備済み／設定不足」を判断し、応答の `execution` field の有無で「実行情報あり／未開始」を判断する。

未決:
- なし

回答:
- `Q-003`: 回答済み
  - `回答`: Ready 期 AI 設定 record は利用者が選択した AI サービス、モデル、処理方式の 3 値だけを保持する。認証参照、認証状態、利用可能モデル一覧は AIサービス設定（provider-settings 正本）から都度解決する。`JOB_PHASE_AI_SETTINGS` の column からは `credential_ref` を除外する。
  - `回答者`: 人間レビュー
  - `回答日`: 2026-06-01
  - `反映先`: `er-REQ-001`（保持 3 値）、`er-REQ-002`（`JOB_PHASE_RUN.credential_ref` の責務）、`term-translation-phase-REQ-002`（解決責務分担、backend 応答の field 構造）
- `Q-004`: 回答済み
  - `回答`: `JOB_PHASE_AI_SETTINGS` の主キーは `phase_type` のみとし、ジョブとは関連を持たない。3 フェーズ種別の 3 件のみが存在し得る。`user_id` 列も持たない（本 repo は単一利用者ローカル app）。明示的な削除 API は持たず、上書き（upsert）で値を変える。ジョブ削除に連動する cascade は無い。`SaveAISettings` 入力から `job_id` を抜く。フェーズ開始時は `phase_type` ベースで `JOB_PHASE_AI_SETTINGS` record を参照して `JOB_PHASE_RUN` へ転写する。本設計により AI 設定値はジョブ間で使い回せる。
  - `回答者`: 人間レビュー
  - `回答日`: 2026-06-01
  - `反映先`: `er-REQ-001`（主キー、所有関係、削除経路、user_id 無し）、`er-REQ-002`（`phase_type` ベースの転写参照）、`term-translation-phase-REQ-002`（SaveAISettings 入力構造、`phase_type` ベース転写）

### `term-translation-phase-REQ-007` 操作と状態を利用者が判別できる（変更）

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/term-translation-phase.md`

親要件:
利用者は開始、一時停止、再開、再試行、取り消しの可否と理由を判断できる。

仕様（既存仕様への追加分のみ列挙）:
- 単語翻訳フェーズの開始操作は、`JOB_PHASE_AI_SETTINGS` の `phase_type = word_translation` record が存在し（AI サービス、モデル、処理方式の 3 値が保持されている）、かつ AIサービス設定（provider-settings 正本）から認証参照を解決できる場合だけ成立する。
- 単語翻訳フェーズ用 record が存在しない、または record の保持値と provider-settings 正本側の解決結果のいずれかが解決できない場合、開始操作は成立しない。
- 利用者は、開始操作が成立しない原因として「`JOB_PHASE_AI_SETTINGS` record の不在」によるものか、「record は存在するが、保持値または provider-settings 正本側の解決結果のいずれかが解決できない」によるものかを判断できる。
- 「AI 設定不足」「blocked 理由」など利用者向けの説明表現は frontend 表示の責務とし、backend 応答には含めない。

未決:
- なし

回答:
- なし

### `persona-generation-phase-REQ-002` AI 設定を開始時に再解決する（変更）

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/persona-generation-phase.md`

親要件:
フェーズ開始と再試行は最新の AIサービス設定を参照し、秘密値を利用者向け情報から分離する。

仕様（既存仕様への追加・変更分のみ列挙）:
- ペルソナ生成フェーズの Ready 期 AI 設定の保存、参照、変更は、`JOB_PHASE_AI_SETTINGS` の `phase_type = npc_persona_generation` record として扱う。
- `JOB_PHASE_AI_SETTINGS` record は AI サービス、モデル、処理方式の 3 値を保持する。認証参照、認証状態、利用可能モデル一覧は AIサービス設定（provider-settings 正本）の責務であり、`JOB_PHASE_AI_SETTINGS` は保持しない。
- ペルソナ生成フェーズ用 record が存在しない状態は「AI 設定不足」として判断できる。
- ペルソナ生成フェーズ用 record が存在し、保持 3 値と AIサービス設定の認証参照解決結果のすべてが解決済みである状態は「AI 設定準備済み」として判断できる。
- ペルソナ生成フェーズ開始時は、`phase_type = npc_persona_generation` の record の 3 値と AIサービス設定の最新解決結果（認証参照）を合わせて、実行時固定値を `JOB_PHASE_RUN` の AI 設定列へ転写する。
- ペルソナ生成フェーズ実行中の表示は `JOB_PHASE_RUN` の AI 設定列を、未開始の表示は `JOB_PHASE_AI_SETTINGS` record と AIサービス設定の解決結果を判断根拠にする。
- `JOB_PHASE_AI_SETTINGS` record の不在を、空値の record で代理表現しない。

AI 設定保存操作（`SaveAISettings` 相当）の入力構造に関する仕様:
- 入力は `phase_type` と AI 選択値（AI サービス、モデル、処理方式、`batchMode`）のみで構成する。`job_id` を含めない。
- 同一 `phase_type` への再保存は upsert として処理する。明示的な削除入力は持たない。

backend 応答の構造に関する仕様:
- ペルソナ生成フェーズの状態を返す backend 応答は、フェーズ単位の「設定全体の情報」を返すか、「該当状態が存在しないこと」を示すかの二択でだけ表現する。
- 設定不足理由などの説明文字列を backend 応答に含めない。状態判断の表現は frontend 表示の責務とする。
- `JOB_PHASE_AI_SETTINGS` record が存在しない状態は、応答の AI 設定 field を不在で表現する。空文字値で代理表現しない。
- 実行情報（`execution` 相当 field）は `JOB_PHASE_RUN` が存在する場合だけ応答へ含める。

未決:
- なし

回答:
- なし

### `body-translation-phase-REQ-002` AI 設定を開始時に再解決する（変更）

- `変更種別`: 変更
- `要件扱い`: 既存要件
- `正本反映先`: `docs/detail-specs/body-translation-phase.md`

親要件:
フェーズ開始と再試行は最新の AIサービス設定を参照し、秘密値を利用者向け情報から分離する。

仕様（既存仕様への追加・変更分のみ列挙）:
- 本文翻訳フェーズの Ready 期 AI 設定の保存、参照、変更は、`JOB_PHASE_AI_SETTINGS` の `phase_type = text_translation` record として扱う。
- `JOB_PHASE_AI_SETTINGS` record は AI サービス、モデル、処理方式の 3 値を保持する。認証参照、認証状態、利用可能モデル一覧は AIサービス設定（provider-settings 正本）の責務であり、`JOB_PHASE_AI_SETTINGS` は保持しない。
- 本文翻訳フェーズ用 record が存在しない状態は「AI 設定不足」として判断できる。
- 本文翻訳フェーズ用 record が存在し、保持 3 値と AIサービス設定の認証参照解決結果のすべてが解決済みである状態は「AI 設定準備済み」として判断できる。
- 本文翻訳フェーズ開始時は、`phase_type = text_translation` の record の 3 値と AIサービス設定の最新解決結果（認証参照）を合わせて、実行時固定値を `JOB_PHASE_RUN` の AI 設定列へ転写する。
- 本文翻訳フェーズ実行中の表示は `JOB_PHASE_RUN` の AI 設定列を、未開始の表示は `JOB_PHASE_AI_SETTINGS` record と AIサービス設定の解決結果を判断根拠にする。
- `JOB_PHASE_AI_SETTINGS` record の不在を、空値の record で代理表現しない。

AI 設定保存操作（`SaveAISettings` 相当）の入力構造に関する仕様:
- 入力は `phase_type` と AI 選択値（AI サービス、モデル、処理方式、`batchMode`）のみで構成する。`job_id` を含めない。
- 同一 `phase_type` への再保存は upsert として処理する。明示的な削除入力は持たない。

backend 応答の構造に関する仕様:
- 本文翻訳フェーズの状態を返す backend 応答は、フェーズ単位の「設定全体の情報」を返すか、「該当状態が存在しないこと」を示すかの二択でだけ表現する。
- 設定不足理由などの説明文字列を backend 応答に含めない。状態判断の表現は frontend 表示の責務とする。
- `JOB_PHASE_AI_SETTINGS` record が存在しない状態は、応答の AI 設定 field を不在で表現する。空文字値で代理表現しない。
- 実行情報（`execution` 相当 field）は `JOB_PHASE_RUN` が存在する場合だけ応答へ含める。

未決:
- なし

回答:
- なし

## 波及範囲

| 項目 | 扱い |
| --- | --- |
| 単語翻訳フェーズ（`word_translation`） | 本 task で `er-REQ-001`、`er-REQ-002`、`term-translation-phase-REQ-002`、`term-translation-phase-REQ-007` を反映する。 |
| ペルソナ生成フェーズ（`npc_persona_generation`） | 本 task で `persona-generation-phase-REQ-002` を同型反映する。実装範囲も本 task に含める。 |
| 本文翻訳フェーズ（`text_translation`） | 本 task で `body-translation-phase-REQ-002` を同型反映する。実装範囲も本 task に含める。 |

波及範囲の判断根拠:
- 人間レビュー Q-001 の回答により、3 フェーズ種別を同時に独立 record 設計（`JOB_PHASE_AI_SETTINGS`、`phase_type` 主キー）へ移行する。`JOB_PHASE_AI_SETTINGS` は 3 種別共通の汎用テーブルとして再利用性を担保する。
- 上記回答により `fix-decision.md` の「禁止修正 5」（persona-generation・body-translation の同一パターン同時修正の禁止）は本 task で覆る。`fix-decision.md` 自体の書き換えは investigation-module 範囲のため本差分では指示しない。後続の investigation-module 入口で覆りを反映する。
- ER 正本は 1 枚で 3 フェーズ種別すべてを表現するため、独立テーブルの導入と 3 detail-spec の同型差分は整合する。

## 影響する正本ファイル

| ファイル | 差分内容 |
| --- | --- |
| `docs/er.md` | 「フェーズとフェーズ実行」節に `JOB_PHASE_AI_SETTINGS`（3 フェーズ種別共通、主キー `phase_type` のみ、ジョブと無関連、user_id 無し、上書きのみ、cascade 無し）の責務、および `JOB_PHASE_RUN` AI 設定列の責務（実行時固定値の記録に限定）を追記する。 |
| `docs/diagrams/er/combined-data-model-er.puml` | `JOB_PHASE_AI_SETTINGS` テーブルを追加し、`phase_type` を主キーとして AI サービス、モデル、処理方式の 3 列を描く。`TRANSLATION_JOB` との関連線は描かない。 |
| `docs/detail-specs/term-translation-phase.md` | `term-translation-phase-REQ-002`、`term-translation-phase-REQ-007` を本差分に従って改訂する。`JOB_PHASE_AI_SETTINGS` の `phase_type = word_translation` record としての扱い、保持 3 値、認証関連の provider-settings 都度解決、`SaveAISettings` 入力に `job_id` を含めないこと、backend 応答の field 構造、利用者向け説明表現の frontend 責務化を明文化する。 |
| `docs/detail-specs/persona-generation-phase.md` | `persona-generation-phase-REQ-002` を本差分に従って改訂する。`phase_type = npc_persona_generation` record としての扱いと同型の差分を反映する。 |
| `docs/detail-specs/body-translation-phase.md` | `body-translation-phase-REQ-002` を本差分に従って改訂する。`phase_type = text_translation` record としての扱いと同型の差分を反映する。 |

## 根拠

- `source`:
  - `./plan.md`（想定 Y/N 再評価、人間判断「ER 仕様変更を含めて検討」、人間レビュー指摘「空文字は異常結果でしかない」「snapshot は『動いてるぶん』だけ」「設定なしは record 不在で表現」）
  - `./fix-decision.md`（確定原因: `applyTermTranslationRuntimeSnapshot` の `ErrNotFound` 時空文字上書きが仕様逸脱）
  - `docs/er.md:23-26,64-73`（フェーズ別 AI 設定の保持規約、Ready 期 `JOB_PHASE_RUN` 未作成規約）
  - `docs/diagrams/er/combined-data-model-er.puml:167-183`（`JOB_PHASE_RUN` の AI 設定列定義）
  - `docs/detail-specs/term-translation-phase.md`（`term-translation-phase-REQ-002`、`term-translation-phase-REQ-007` 既存仕様）
  - `docs/screen-design/screens/term-translation-phase.md:91-129`（AI モデル選択領域の状態別表示）
- `review`: 人間レビュー（plan.md の停止判断、design-module 入口の再評価記録、差し戻し指摘 1「backend 応答は『全体』か『不在』の二択、派生説明は禁止」、差し戻し指摘 2「credential は AIサービス設定の責務、phase AI 設定 record から除外」、`Q-003` 回答 2026-06-01、`Q-001`/`Q-002`/`Q-004` 回答 2026-06-01）。本差分は人間レビューの 3 件の確定回答を反映し、未決 0 件で `status: ready-for-human-review` を維持する。
- `validation`: 未実行。プロダクト変更は伴わず、差分は ER 正本と詳細仕様正本の文章設計差分のみ。
