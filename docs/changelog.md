# 変更・判断履歴

正本（`requirements.md`、`system_requirements.md`、`architecture.md` など）には現在の状態だけを書く。
「なぜ変えたか」「何を落としたか」などの判断履歴は本ファイルに残し、正本へ混ぜない。
新しい entry を上に追加する。1 entry は date 見出しで区切る。

## 2026-06-14 T1 後の architecture.md との構造差異を整合

### 変更

- keyring secret store を `internal/repository/` から `internal/store/secret/`（package `secret`）へ移動。`architecture.md` §3・§7 の「secret 子に置く」に合わせた。`.go-arch-lint.yml` の component を `repository` → `secret` に更新。
- `db.Apply` に schema version の読み込み時検査を追加。DB の `user_version` がアプリの想定 migration 数より新しければ適用せずエラーにする（`architecture.md` §6「Go は読み込み時に version を検査する」を実装）。
- `architecture.md` §4: 多態の port は `provider` 1 つのみと明記し、`store` 用の狭い interface（consumer 側・実装 1 つ・単体テスト用）は port ではない切り離しとして許容を追記。
- `architecture.md` §5: runtime の閉じ込め先を `bootstrap` と `api`（Bind 公開面）に明記。`api` が runtime を進捗 push とファイル選択ダイアログに使うことを許容し、下位層へは漏らさないと固定。
- `architecture.md` §7: `db/` に migration 適用（`db.Apply`）を追記。`store` が起動時に委譲する旨を記載。

### 判断

- T1 実装が `architecture.md` と食い違った 3 点を、コード修正と doc 改訂に振り分けて整合した。
  - keyring 場所: doc が明示指定（secret 子）。コードを doc に合わせて移動。
  - store の狭い interface: テスト容易性の利益が大きく、design-module のテスト設計（engine を mock 越しに試す）と整合する。doc を実態に合わせて改訂。
  - migration 適用の場所: ユーザー指示「migration とリポジトリは分けて」で `db` パッケージへ分離済み。doc に明記。
- Wails runtime の `api` 直接利用は、§2 図と §3 が「`api` が runtime events を push」と示し、§5 の閉じ込め先 adapter に Bind 公開面（`api`）が含まれるため、乖離ではないと確定。§5 を明文化して曖昧さを除いた。
- 残る差異は意図的な未実装（provider 3 系統・engine の重複排除/辞書/ペルソナ/XML・進捗 push）で、後続タスクで埋める。§8「現在の状態と移行」の陳腐化は別途更新する。

### 検証

- Go test 緑、backend lint（format/vet/static/arch/module）0 issues。store の version 検査込みで store test 緑。

## 2026-06-14 T1 最小縦切り（抽出 → 翻訳 → DB → 画面）を実装

### 変更

- backend（Go）を初実装。`internal/model`（Narration）、`internal/store`（sqlx ＋ modernc.org/sqlite、migration 適用）、`internal/provider`（Translator port ＋ OpenAI 互換実装）、`internal/engine`（未訳を翻訳し仮訳で書き戻す手続き）、`internal/api`（Wails Bind 公開面）、`internal/bootstrap`（composition root）、`main.go`（Wails entry）。
- `db/migrations/0001_init.sql`：narration テーブルの DDL（C#↔Go 契約 1 本）。`db/migrations.go`：embed して公開。
- `tools/extractor`（C#）に `NarrationSqliteWriter` と `--sqlite` モードを追加。BOOK:DESC を narration へ UPSERT。
- frontend を daisyUI で再構築。Tailwind v4 ＋ daisyUI v5 の独自テーマ `dovahkael`、汎用部品（Field/TextField/SelectField/FileSelectField/StatusBadge）、画面 `TranslationRunScreen` ＋ container、gateway、`main.ts`/`App.svelte`。
- lint 整備：`.go-arch-lint.yml` 新設（新層の依存方向）、`.golangci.yml` の static 違反解消、frontend eslint で生成 `wailsjs` を除外、`wails-boundary.test.mjs` で Wails 境界を検証。
- `.gitignore`：`db` 全体無視を `db/*.sqlite3*` に絞り、`db/` の source を追跡。

### 判断

- 叙述文 1 種は `BOOK:DESC`（書物本文）。装備 DESC への拡張は `TranslationCounts.Enumerate` のフィルタ追加だけで済む。
- provider 接続情報（endpoint/apiKey/model）は永続化せず画面から都度渡す。API キーなしの OpenAI 互換（LM Studio 等）に対応するため、キーが空のとき Authorization を付けない。base URL は `/v1` 配下へ正規化（`http://127.0.0.1:1234` でも届く）。
- 抽出は Go の api が C# extractor を `dotnet run` で子プロセス起動し、続けて engine を呼ぶ同期手続き。進捗 push は対象外。
- AI 翻訳は訳状態 3（仮訳）で書き戻す。
- 起動時に中心 DB の現状を読み込み、前回の結果を画面に出す。

### 検証

- TDD：provider（/v1 正規化・auth・getModels・翻訳）、store（migration・未訳取得・dest 更新）、engine（仮訳書き戻し・provider エラー伝播）、api（status ラベル・DTO・extractor 引数）、C# NarrationSqliteWriter（BOOK:DESC 書き込み・冪等）を失敗テスト先行で実装。
- Go test 緑、backend lint（format/vet/static/arch/module）緑、frontend lint（eslint/tsc/knip/boundaries）緑、C# 17 テスト緑、build-storybook 緑。
- 実 app（`dev:wails:run`、localhost:34115）で end-to-end を目視確認。OpenAI 互換モック（`127.0.0.1:1234`）に対し、getModels でモデル選択、Dawnguard.esm から 65 件抽出 → 翻訳 → SQLite → 見開き対訳表示まで動作。LM Studio を同 endpoint に立てれば同経路で実訳になる。

### 残課題

- ファイル選択ダイアログ（Wails OpenFileDialog）は実装済みだが、ネイティブダイアログのため自動 UI テストは未。
- 大量レコードの同期翻訳は進捗表示が無く待ち時間が長い（進捗 push は後続）。
- 書物本文の HTML 様タグ（`<font ...>`）の扱いは未整理。
- フォントは Google Fonts CDN。デスクトップ app 用の self-host は後続。
- greenfield 未配線の `diagnostic`・`shell-state`・`pino` は knip ignore で保持（将来配線で解消）。

## 2026-06-14 ER 設計（抽出入力）の正本 er.md を新設

### 変更

- `docs/er.md`: 新設。`concept-model.md` の箱（抽出入力）を `SQLite` の物理テーブルへ写す ER 設計。テーブル定義・関係・concept-model 対応・既知の論点を記述。
- `docs/index.md`: `er.md` を Read Order・Directory Contract・Choose The Right Record に登録。

### 判断

- スコープは抽出入力（`concept-model.md` の 10 箱と関連 e1〜e14）に限定。マスター辞書・ペルソナルール・翻訳ジョブ/結果キャッシュ・schema version 管理は対象外（あとで別途追加）。
- 概念モデルから外れない。テーブルは `concept-model.md` の箱と 1 対 1。箱を統合せず、属性（`人称`・`口調`・`背景`・`性質` を含む）も落とさない。
- 実現方式を ER に持ち込まない。重複排除のタイミング、属性の充填時期、永続化の有無は `concept-model.md` L7 のとおり実現方式の責務とし、ER は構造だけを固定する。
- 正規化は根拠を明示する。多対多と可変多重度（e4/e5/e6/e7/e8/e10/e14）は連関テーブル（第1正規形）、1 対多（e1/e2/e3/e9/e11/e12/e13）は FK 1 本。訳の単位の分離は更新異常の除去。レコード識別の分解は第1正規形と xTranslator 出力要件。form_id と edid の同居は出力自己完結のための意図的冗長。
- レコード識別キーは xTranslator String 行の `(plugin, form_id, rec, field, ordinal)`。`status` は xTranslator `Status`（0-4）を踏襲（`references/xtranslator_ref.md`）。
- 実 SQL DDL の正本は `db/` migration（`architecture.md` §7）。`er.md` は論理 ER 設計に限定し、DDL を二重に持たない。

### 経緯

- 初版で `配置`・`叙述文`・`台詞`・`無訳片` を `extracted_string` 1 テーブルに、`固有名`・`定型句` を `translation_unit` 1 テーブルに統合し、重複排除責務を `engine` に寄せる実現方式を持ち込んだ。これは概念モデルから外れる逸脱で、人間指摘により 10 箱 1 対 1 へ作り直した。

### 残課題

- 初版。ファイルレビューで確定する。
- 言及（e4/e5）の検出方式、純汎用台詞の話者群の口調決定は実現方式で決める（`concept-model.md` 弱点）。
- 対象外テーブル（辞書・ルール・ジョブ・schema version）は別設計。

## 2026-06-14 tech-selection.md の責務外記述を除去（採用技術へ純化）

### 変更

- `docs/tech-selection.md`: §2・§3・§4・§7 から、採用技術でない記述（データ配置、C#↔Go 設計、責務・依存方向、別プロセス構成、抽出の挙動・制約、観測ログ出力先）を除去した。`SQLite`・`sqlx`・SQL migration・`pino`・`log/slog`・`C#/.NET`＋`Mutagen`・`.NET 8` などの採用技術そのものは残した。

### 判断

- `tech-selection.md` の責務は `index.md` L34 と `core-beliefs.md` §2 で「採用技術と品質基盤＝実装技術の選択」と定義されている。データ配置・内部境界・依存方向は `architecture.md`、観測ログ出力先は `observability-logging.md` の責務で、`core-beliefs.md` §3 が「同じ責務を複数文書で別定義している状態」を除去対象としている。
- §4 永続化・§7 抽出基盤の架構記述は、同日の「アーキテクチャ再構築」entry で `tech-selection.md` に追加したが、`architecture.md` と重複していた。責務分離のため `architecture.md` 側へ一本化し、`tech-selection.md` からは除去した。
- 除去した記述はすべて他の正本に既出で、情報損失はない。対応は次のとおり。
  - DB が持つ内容（抽出入力・マスター辞書・ペルソナルール・翻訳ジョブ）= `architecture.md` §3。
  - C# extractor 直書き・中間形式なし、SQL schema が C#↔Go の唯一の契約、migration 適用責務と extractor の冪等 ensure = `architecture.md` §1/§6。
  - 上位層へ driver 固有 API を漏らさない = `architecture.md` §4。
  - 別プロセス構成・構造体モデル準拠・Data folder の明示パス指定・macOS 実行可・抽出結果の SQLite 書き込み = `architecture.md` §1/§6。
  - backend/frontend 観測ログの出力先（`stderr`／browser console）= `observability-logging.md` §1。

### 残課題

- なし。採用技術の選択内容は変えていない。

## 2026-06-14 アーキテクチャ再構築（データ中心＋手続き、Go 維持＋SQLite 境界）

### 変更

- `docs/architecture.md`: 旧層構成（`UseCase` / `Service` / per-entity `Repository` / `Presenter` ＋ 厚い手動 DI）前提を破棄し、データ中心かつ手続き中心の骨格へ全面書き換え。Mermaid コンポーネント図、各箱の責務、依存方向、Wails 境界、C#↔Go 境界（SQLite 契約）、ディレクトリ正本を記述。
- `docs/tech-selection.md`: §4 永続化を SQLite 中心へ書き換え（抽出 sink を SQLite 正本に、JSON 中間形式を廃止、SQL schema を C#↔Go 契約に）。§7 として翻訳対象抽出基盤（C#/.NET ＋ Mutagen）を新設。§6 に抽出ツールの `xUnit` 検証を追記。公式参照に Mutagen を追加。

### 判断

- 概念モデルが示す実体（中心は Skyrim データ、翻訳は一本の手続き）に対し、旧層構成は過剰と判断。層を薄くして間接化を削る。
- engine の runtime は Go を維持する。Wails / Svelte / 既存 harness を温存するため。
- C#↔Go の受け渡し境界は SQLite とする（案1 ＝ C# extractor が SQLite へ直接書く）。境界専用の JSON 中間形式は持たない。理由: 旧 JSON 境界は xEdit Pascal が翻訳ロジックを載せられない制約に由来する。Mutagen は通常の .NET ライブラリ呼び出しで、tech-selection が既に「入力データを SQLite に持つ」方針のため、境界を SQLite に寄せれば境界専用形式を作らずに済む。
- 抽象（port）は AI provider の境界 1 つだけに置く。実装が 4 系統（Gemini / xAI / OpenAI 互換 / Claude）に分かれる唯一の箇所のため。
- 手動 DI は composition root 1 箇所へ集約する。
- 抽出の検証は C# テスト（`CountParityTests` / `ModelInvariantTests`）へ移管済み。Python の `validate_extraction.py`・`compare_counts.py` は重複のため削除済み。
- Mutagen は macOS でも動くことを公式 docs で確認した（`GameEnvironment.Typical.Builder().WithTargetDataFolder()` で registry 自動検出を回避する）。「macOS では Mutagen を動かせない」という当初の想定は誤りで撤回した。
- engine 内部のパッケージは approach A（`engine` / `model` / `store` / `provider`）を採用する。
- `.NET` 統合（Go を捨て Mutagen と engine を 1 プロセスへ統合する案）は今回不採用。Go 維持で進める。

### 残課題

- engine、store、provider、api、bootstrap、frontend 配線の実体は未実装。再構築は `docs/exec-plans/active/` の active plan で進める。
- SQL schema（中心データの具体テーブル）は未設計。concept-model の箱を SQLite テーブルへ写す設計を plan で詰める。
- C# extractor の SQLite writer は未実装（現状は in-memory `ExtractionResult` の件数検証まで）。

## 2026-06-13 業務要件2・3 のシステム要件を一部確定（対応策の発散と絞り込み）

### 変更

- `docs/system_requirements.md` 業務要件2「単語の一貫性」: `TBD` から、確定部分（Mod 横断マスター辞書、用語特定はレコード固有名詞の機械抽出）と未確定部分（訳語供給、漏れ語対応、適用方式、検証）を分けて記述。
- `docs/system_requirements.md` 業務要件3「NPC の口調」: 「属性と会話履歴から AI でペルソナ生成」から「属性からルールベースで生成」へ方針変更。表現（構造化属性＋翻訳ディレクティブ）、変換（テンプレート・機械的）、永続境界（ルール集を永続・翻訳前設定／個々ペルソナ非永続）、適用（プロンプト注入）を記述。

### 判断

- 一貫性のスコープは Mod 横断（永続マスター辞書）を選択。ジョブ内・Mod 内は不採用。
- 用語特定は、誤検出が無く既訳ヒット率が高いレコード固有名詞の機械抽出のみを採用。AI 抽出・頻度抽出・辞書マッチによる漏れ語対応（第2層）は保留。
- 頻度抽出は単言語では対訳を出せず、対訳コーパスがある場合の統計アラインメントでのみ対訳を出せると整理。用語特定（どの語を揃えるか）と訳語供給（訳語をどう得るか）は別軸と確認。
- 業務要件3 は機械的抽出（属性 → ルール）をまず採用し、AI 生成と会話履歴解析は保留。
- ペルソナ表現は構造化属性（内部表現）＋翻訳ディレクティブ（適用形）の 2 段とし、変換はテンプレートで機械的に行う。
- 機械的抽出ではペルソナを NPC レコードから都度導出できるため、個々ペルソナは永続化しない。永続資産は「属性 → 翻訳指示のルール集」とし、翻訳前にユーザーが設定可能とする。過去構想の `master-persona`（個々ペルソナを永続・編集）とは性質が異なる。
- VoiceType（声タイプ。Skyrim が音声収録のため NPC を声・性格でグループ化した分類）は属性の中で口調と相関が高く、ルールの主軸候補とする。

### 残課題

- 業務要件2: 訳語供給方式（既訳流用のみか AI 併用か）、本文・会話中の漏れ語対応、辞書の適用方式、一貫性検証が未確定。
- 業務要件3: 使用属性の選定、属性の分類、ルール合成の衝突優先順位を概念モデルで整理する。
- ペルソナ属性・Skyrim 構造の概念モデルの置き場が未定（`index.md` Read Order の `skyrim-structure-model.md` は実体なし）。

## 2026-06-13 screen-design 廃止、画面の正本を Storybook へ。tech-selection に Storybook / Tailwind / daisyUI

### 変更

- `docs/screen-design/` を削除（README.md、design-system-ethereal-archive.md、code.html、screens/ 配下すべて）。
- `docs/tech-selection.md`: フロントエンドに Tailwind CSS、daisyUI、Storybook を追加。CSS framework 不採用の行を削除。公式参照 3 件を追加。
- `docs/index.md`: Read Order・Directory Contract・Choose The Right Record から screen-design を削除。画面・表示の設計は Storybook（`frontend/`）と明記。
- `.claude/skills/`: design-module、storybook-module、finalization-module、coding-protocol、fix-decision、investigation-module、implementation-module、diagramming の 8 件を「画面の正本 = Storybook」へ作り替え。
- `docs/exec-plans/active/README.md`、`templates/work-plan.md`、`templates/task-folder/{README.md, plan.md, detail-spec-diff.md}`: `screen-design-diff` と「Storybook 後画面設計差分整合」前提を Storybook 正本へ付け替え。
- memory（`feedback_boundary_responsibility_separation`、`feedback_storybook_module_trigger`、`feedback_implementer_no_agent_split`、`MEMORY.md`）: 画面設計参照を Storybook へ更新。

### 判断

- 画面・表示の設計判断の置き場を Storybook の story と svelte コンポーネントへ移す（ユーザー選択「Storybook を正本にする」）。
- `画面設計差分`（`screen-design-diff.<screen-id>.md`）doc を廃止。design-module は画面表示 doc を作らず、画面表示の設計は storybook-module が story とコンポーネントで直接行う。
- storybook-module の「Storybook 後画面設計差分整合」成果物を廃止。実装範囲を越える画面変更が要る場合は design-module へ差し戻して `実装範囲` を見直す経路に統一。
- finalization-module の docs 正本反映対象を `docs/architecture.md` のみに限定。画面の正本（Storybook）は frontend source として作業 commit に含める。
- fix-decision / investigation-module の画面再現確認の selector 正本を「実装済み画面の `data-testid` またはセレクタ」へ変更。
- AI プロバイダ指定「OpenAI API(llama, lmstudio)」は OpenAI 互換 API でローカル実行（llama、LM Studio）と OpenAI 本体を含む意味と解釈。

### 残課題

- exec-plans templates には screen-design と無関係の旧表記（Codex 名称、`docs/detail-specs/`、`agent-browser`）が残る。本変更の対象外。整理は別途判断する。
- `index.md` Read Order の `skyrim-structure-model.md`、`core-beliefs.md` の `er.md` は実体が無い既存リンク（本変更前から）。

## 2026-06-13 要件文書を業務要件とシステム要件へ分割

### 変更

- `docs/spec.md` を `docs/requirements.md` へ rename（git mv、業務要件の内容は不変）。
- `docs/system_requirements.md`: 新規作成。業務要件 1〜4 に対応するシステム要件を記述。1 = AI 利用（Gemini / xAI / OpenAI 互換 API（OpenAI、llama、LM Studio）/ Claude）、2 = TBD、3 = NPC 属性と会話履歴からペルソナ生成、4 = 機能要件＝業務要件。
- `docs/index.md`: Read Order、Directory Contract、Choose The Right Record を `requirements.md` / `system_requirements.md` へ更新。
- `docs/core-beliefs.md`: 関連文書リンクと記録方針を「業務要件 = `requirements.md`、システム要件 = `system_requirements.md`」へ更新。
- `docs/architecture.md` / `docs/tech-selection.md`: 関連文書リンクを `requirements.md` / `system_requirements.md` へ更新。
- `docs/screen-design/README.md`: `spec.md` 参照を `requirements.md` へ更新。

### 判断

- 業務要件（何をしたいか）とシステム要件（どう達成するか）を別文書に分ける。
- 単語の一貫性のシステム要件は TBD として明示的に保留（ユーザー判断）。
- AI プロバイダ指定「OpenAI API(llama, lmstudio)」は、OpenAI 互換 API でローカル実行（llama、LM Studio）と OpenAI 本体を含む意味と解釈した。

### 残課題

- `.claude/skills/coding-protocol/SKILL.md` の 2 箇所が `docs/spec.md` を参照（system 要件の参照行、docs 正本一覧）。auto mode で skill 編集が拒否されたため未修正。ユーザー承認後に `requirements.md` / `system_requirements.md` へ張り替える。
- `requirements.md` は用語集を廃止済みのため、`screen-design/README.md` の「用語」参照は形式的に古い。画面設計の用語運用を決める時に整理する。

## 2026-06-13 spec.md を業務要件専用へ書き換え

### 変更

- `docs/spec.md`: 恒久要件・用語集・状態機械を全削除し、4 つの業務要件（翻訳したい / 単語の一貫性 / NPC の口調 / xTranslator 形式出力）へ全面書き換え。各要件に目的を併記し、成功条件は不記載。

### 判断

- `spec.md` は業務要件（何をしたいか）だけにする。システム要件（どう実現するか）は別文書で扱う。
- 入力取得手段（xEdit 抽出など）はシステム要件側へ回す。業務要件側は対象を「Skyrim Mod のテキスト」とだけ書く。
- xTranslator 形式出力は、ツール固有だがユーザーの明示要望のため業務要件として残す。
- 成功条件は記載しない（ユーザー判断）。目的は記載する。

### 残課題

- システム要件の置き場が未定。入力取得手段、AI 基盤、ジョブ運用は置き場を決めてから書き起こす。
- `index.md`（`spec.md` を「恒久要件と用語集」と記述）と `core-beliefs.md`（「永続要件は `spec.md` に記録する」と記述）の文言が、業務要件専用化により古くなった。システム要件の文書を決める時に併せて直す。
