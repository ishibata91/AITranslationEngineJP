# narration-line-mention-linking（叙述文・台詞本文中の固有名詞の言及検出とテーブル実装）

## 分岐情報

- 作業 branch: `claude/narration-line-mention-linking`
- 分岐元 branch: `master`
- 分岐元 commit: `cf5d4038`

## 背景（現状の仕組み）

固有名一貫性は、現状は横断マスター辞書（`master_term`）の機械置換注入だけで担保している。叙述文（narration）・台詞（line）の本文をAI翻訳へ渡す直前に、本文中に含まれる既知の固有名詞（`master_term`∪`proper_noun`に確定済みの語）を文字列一致で検出し、確定訳語をプロンプトへ注入する（`internal/engine`の`LoadDictionary`・`dict.Apply`）。

この検出は、結果一覧の取得処理（`ListResultsPage`、`internal/api/app.go`）が結果を返すたびに`dict.Apply`で都度実行し直す再構成であり、DBへ記録としては残らない。再構成した内訳は`ResultView.Terms`（機械置換内訳）として画面へ出るが、これは叙述文・台詞のどちらの結果行にも同じ形で出ており（`app.go`の叙述文ループ・台詞ループはどちらも同じ`termViews(used)`を設定する）、表示の有無に非対称は無い。

非対称があるのは表示ではなく概念モデル上の関連の種類である。`docs/concept-model.md`のe3/e4/e5によると、`narration_mention`（e4、叙述文→固有名の本文中言及）と`line_mention`（e5、台詞→固有名の本文中言及）は対称な関連だが、`narration.described_proper_noun_id`（e3、叙述文→固有名の説明対象。武器の説明文→その武器自身の名前など）は叙述文だけが持つ、単なる言及より強い関連であり、台詞に対応する概念は無い。

## 依頼要約

`docs/known-issues.md` 1番「固有名一貫性の後続 task（言及テーブル未実装）」に着手する。上記の都度再構成による機械置換注入を、本文中の言及関連としてDBへ永続化する連関テーブルへ広げる。概念モデル上の対応する未実装要素は次のとおり（`docs/concept-model.md`、`docs/er.md`）。

- `narration_mention`（e4）: 叙述文 → 固有名（本文中の言及）
- `line_mention`（e5）: 台詞 → 固有名（本文中の言及）
- `narration.described_proper_noun_id`（e3）: 叙述文 → 固有名（説明対象）の未実装 FK

## 本taskの位置づけ（機能追加であり、既存出力のリファクタではない）

既存の翻訳出力（叙述文・台詞のdest文、`Terms`機械置換内訳、実プロンプト再構成）は変えない。既存の機械置換注入・都度再構成のロジックはそのまま残す。

新規に、本文中の言及をDBへ永続化する関連（`narration_mention`・`line_mention`、および`narration.described_proper_noun_id`）を追加する。追加した関連を今回のtaskの中で画面表示や他機能へ接続する予定は無い。将来の事後検証（known-issues.md 2番「注入語の保持確認」）や言及ベース機能のための土台を、先に作る位置づけである。

## 完了定義

### 動かす範囲

- 叙述文・台詞の本文中に、`master_term`・`proper_noun`へ確定済みの固有名詞が出現する箇所を検出し、言及レコード（`narration_mention`・`line_mention`相当のテーブル）へ実データに対して記録する。
- 叙述文の説明対象FK（e3、`narration.described_proper_noun_id`相当）を、検出した言及の結果で埋める状態にする。

### 観測点

- 単体テスト: 言及検出ロジックの純粋IO部分（`internal/core`配下に新設するpackageを想定）。守るべき不変ルールとしてユニットテスト100%カバレッジを基準にする（`feedback-pure-io-rule-100-coverage`）。
- 実データ: 実際に抽出したnarration・lineに対して検出を実行し、記録された言及レコードをDB内容で確認する。

### 非劣化検証

検出ステップは既存の翻訳手続きへの追加であり、既存部分の出力を変えないことを次の観点で確認する。

- `internal/harness`の合成golden: 検出ステップ追加の前後で、叙述文・台詞のdest文、`Terms`（機械置換内訳）、実プロンプトの再構成結果が一致すること。
- `npm run verify:backend`（go test全体＋arch-lint＋境界走査）が通過すること。
- 検出ステップが読み書きする対象を新規テーブル（`narration_mention`・`line_mention`、および`narration.described_proper_noun_id`列）へ限定し、既存テーブル（`narration`・`line`・`proper_noun`・`master_term`）の既存列（dest・status等）を書き換えないこと。

### 含まない（除外範囲）

- 辞書に無い漏れ語（名前付きレコードに出ない語）の抽出。`known-issues.md`でAI抽出・頻度抽出による第2層として明示的に保留済みのため対象外とする。
- `line_sequence`（e7、会話の流れ）。対話木context整備（`known-issues.md` 2番）と重なる別テーマのため対象外とする。
- `speaker_name`（e8）・`faction_name`（e14、話者・勢力が名乗る名）。本文中の言及とは別概念であり、話者名は現状`master_term`・人名派生で代替できているため対象外とする。
- 検出結果を翻訳結果表示画面へ露出する機能（フィルタ・編集等）。`known-issues.md` 5番のscopeであり別taskとする。

## 軽 / 重判定

| 軸 | 判定 | 根拠 |
| --- | --- | --- |
| 画面が動くか | N | 変更対象は検出ロジックとテーブルの新設であり、layout・文言・style・表示構造・svelte表示コンポーネント・props・story・fixtureのいずれも変えない。 |
| `docs/architecture.md`反映が要るか | N | `internal/core/<name>`への新規pure package追加（既存の`dictionary`・`termderive`等と同型）、`engine`の取込段への検出ステップ追加、`store`・DB migrationへのテーブル追加はいずれも既存層内の追加であり、層構成・依存方向・Bootstrap・Wails境界・強い制約を変えない（`feedback-architecture-reflection-structural-only`の判断基準）。新規Wails Bindは想定しない。 |

判定結果: 両方Nのため軽task。`design-module`と`storybook-module`をbypassし、`preparation-module` → `implementation-module` → `finalization-module`で進める。
検出方式（文字列一致の粒度、重複候補の解決など）は`known-issues.md`で「実装判断に委ねる」と明示されており、軽taskの設計判断として`implementation-module`でClaude本体が固定する。

## 実装記録（implementation-module）

### 実装判断（known-issues.md が実装判断に委ねた点の固定）

- 検出方式: 機械置換（`internal/core/dictionary`）と同じ照合規則（貪欲最長一致・語境界・大小区別・同一原語は先勝ち）を新設の純粋package `internal/core/mention` に持つ。言及レコードと注入語がずれると後続の事後検証（known-issues 2番）が成り立たないため、規則の一致を不変条件とする。
- 語彙と優先順: `LoadDictionary` と同じ供給源（`master_term` 全件 → `proper_noun` 全件）を同じ順で積む。同一原語は master_term が先勝ちし、注入と同じ側が言及の相手になる。
- 訳状態への非依存: 言及は原語（source）の出現で決まる関連のため、dest の確定状態に依存させない。これにより検出を取込段（`Engine.Ingest` の最終ステップ）へ置ける（固有名フェーズを待たない）。
- 言及の相手: `narration_mention`・`line_mention` は排他2列（`proper_noun_id` / `master_term_id`、CHECK で片方だけ非NULL）で両供給源を指せる形にした。概念の固有名箱は `proper_noun` だが、横断辞書 `master_term` だけに載る語（base ゲーム由来の名前）の言及も事後検証に要るため。
- 対象行: `narration` テーブル全行（定型句収容行を含む）と `line` 全行。機械置換が同様に全行へ当たるため対称にした。
- e3（説明対象）の導出: 本文でなくレコード構造から決める。同一レコード（plugin, form_id, rec）の `FULL`（`extracted_field`）を `proper_noun`（category, source）へ結ぶ SQL join（`LinkNarrationDescribed`）。`box='叙述文'` の行だけ対象（定型句は説明対象を持たない）。

### 計画からの逸脱（1件）

- `narration.described_proper_noun_id` 列は追加せず、専用テーブル `narration_described`（narration_id PK → proper_noun_id、0..1）で持つ。理由: C# 抽出器が全 migration SQL を毎回 ensure する契約のため `ALTER TABLE` は再実行で失敗する（migration 0007 の注記と同じ制約）。`line_condition` と同型の冪等な 0..1 テーブルにした。

### 観測

- 冪等の観測は `IngestCounts.Mentions`（追加件数）で返す。slog による観測ログは追加しない（backend logger の正本が新 architecture 確定待ちで、repo に slog 実績が無い。観測点は plan の定義どおり DB 内容で足りる）。

### 最終検証（2026-07-04 通過）

- `npm run verify:backend` 通過（go test 全体・arch-lint・境界走査）。`npm run lint:backend`（format・vet・static・arch・boundary・module）も通過。
- 合成 golden（`internal/harness` TestSyntheticNonRegression）: golden 無更新で一致。dest・Terms・実プロンプトの非劣化を確認。
- 実データ非劣化: 分岐元 commit `cf5d4038` で `scripts/golden/capture.sh`（inigo.esp）を捕獲し、本変更の作業ツリーで `goldcap -mode compare` が一致（プロンプト 8772 件）。
- 実データ DB 確認（inigo.esp、実 harness 全 pipeline）: `narration_mention` 152 件（master側66・proper側86）、`line_mention` 1881 件（master側938・proper側943）、`narration_described` 75 件（叙述文111行中、FULL を持つレコードの分）。排他2列の違反 0 件。実例: 台詞 "…Riften always smells lovely…" → Riften/リフテン、QUST:CNAM のクエストログ → クエスト名 "Conversations started by the player"。
- 単体テスト: `internal/core/mention` カバレッジ 100%（feedback-pure-io-rule-100-coverage 充足）。store の冪等（部分 UNIQUE 索引が NULL 重複を止めること）と e3 解決 SQL は実 SQLite のテストで証明。

### finalization-module への引き継ぎ

- docs 正本反映の候補: `docs/er.md`（narration_mention・line_mention・narration_described の実装済み化、e3/e4/e5 の状態更新、排他2列と専用テーブル化の正規化根拠）、`docs/known-issues.md`（1番の該当行の解決・残課題の整理）、`docs/architecture.md` §の core package 列挙（`mention` の追記。層・依存方向の変化は無し）、`docs/changelog.md`。
- 仕様変更・仕様追加: 無し（既存出力は非劣化。新規テーブルの追加のみ）。

## finalization 記録（finalization-module）

### 正本化判断

- `docs/architecture.md` への反映: 不要。`internal/core/mention` は既存の core package 群と同型の純粋 package 追加、`store`・`engine` の変更も既存層内の追加であり、層構成・依存方向・Bootstrap・Wails 境界・強い制約を変えない（軽/重判定どおり。`feedback-architecture-reflection-structural-only` に従い、構造不変のため人間承認も求めない）。import 方向の実制約は `.go-arch-lint.yml` へ `mention` component として追加済み。
- 恒久仕様の正本反映: 対象なし（仕様変更・仕様追加なし）。
- docs の現在状態反映（正本反映とは別の、実装済み状態への記録更新）: `docs/er.md`（3 テーブルの実装済み化・ER 図・正規化根拠 5/6 追加）、`docs/known-issues.md`（1 番から e3/e4/e5 を除去、2 番へ照合対象の追記）、`docs/concept-model.md`（言及 note の未実装注記を実装参照へ）、`docs/changelog.md`（2026-07-04 entry に経緯と判断を記録）。前 task（`ed540ff8` の `docs/er.md` 反映）と同じ扱いで作業 commit に含める。

### 作業 commit

- branch: `claude/narration-line-mention-linking`（分岐元 `master` @ `cf5d4038`）
- commit: `37e07076`（23 files, +1045 / -185）
- 検証結果: 実装記録の「最終検証（2026-07-04 通過）」のとおり（verify:backend・lint 一式・合成/実データ golden 一致・実データ DB 確認・カバレッジ 100%）。
- 残留リスク: 短い一般語と同綴りの固有名（例: MESG:FULL の "Yes"/"No"）が本文へ言及として当たる。これは既存の機械置換注入が持つ同じ弱点の写しで、言及テーブルは注入の実態を忠実に記録する（概念モデル弱点 1 の受容範囲）。

### local merge・merge 後検証・completed 移動

- local merge: `master` へ `git merge --no-ff claude/narration-line-mention-linking`。merge commit `d257a517`。conflict なし。
- merge 後検証: `master` 上で `npm run verify:backend` 通過（go test 全体・arch-lint・境界走査）。
- completed 移動: 本 plan folder を `docs/exec-plans/completed/narration-line-mention-linking/` へ移動。
