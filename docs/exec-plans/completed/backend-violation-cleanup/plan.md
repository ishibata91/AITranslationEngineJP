# backend-violation-cleanup（複雑度ゲートのリファクタと static baseline の解消）

## 依頼要約

backend-refactor-test-foundation（完了）が安全網として整えた検査が表面化させた残債を解消する。具体的には、認知複雑度ゲート（gocognit）で baseline 除外した 6 関数を分割して `//nolint:gocognit` を外し、golangci-lint static の既存 baseline 11 件を解消する。テスト基盤（golden harness・arch-lint・境界走査）は不変に保ち、非劣化を保証したうえで内部を直す。

- 分岐元 branch: master（backend-refactor-test-foundation の merge 後）
- 分岐元 commit: 877077b5
- 作業 branch: claude/backend-violation-cleanup
- 起動条件: backend-refactor-test-foundation が完了・merge 済みであること。本書はその done を前提にした backlog。

## 完了定義

### 動かす範囲

1. 認知複雑度ゲートの解除。実測 >15 の 6 関数を分割し、`//nolint:gocognit` を 1 つも残さず gocognit を緑にする。各関数の外部から観測できる振る舞い（戻り値・golden）は変えない。
2. static baseline の解消。golangci-lint static の既存 9 件（wrapcheck 4・gosec 2・errcheck 1・staticcheck 1・errorlint 1）を直し、`lint:backend:static` を緑にする。
   注: 当初 11 件だったが、`engine-pure-rule-split`（純粋ルールの core/ 別 package 化）で `engine.LoadRoleSpeech` を廃止したため、`role_speech.go` 由来の gosec G304 1 件と errcheck 1 件が消え 9 件になった。残りの違反の所在も core/ へ移った（下表）。
3. 非劣化の保証。上記の内部変更後も、②の合成入力 golden（`internal/harness`）が一致し、`npm run verify:backend` が緑であること。

### 観測点

- 複雑度: `npm run lint:backend:static` で gocognit 指摘 0、かつコードに `//nolint:gocognit` が 0 件（`grep -rn "nolint:gocognit" internal` が空）。
- static: `sh ./scripts/lint/run-go-backend-lint.sh static` が exit 0（issues 0）。
- 非劣化: `go test ./internal/harness/...` の golden 一致、`npm run verify:backend` が exit 0。

## 軽 / 重判定

| 軸 | 判定 | 根拠 |
| --- | --- | --- |
| 画面が動くか | N | 変更対象は backend のみ（gocognit の関数分割、golangci-lint static 違反の解消）。layout・文言・style・表示構造・svelte 表示コンポーネント・props・story・fixture のいずれも変えない。frontend 変更は除外範囲。 |
| `docs/architecture.md` 反映が要るか | N | 関数内部の段階・分岐の helper 切り出しと lint 違反の局所修正に限る。層構成・依存方向・Bootstrap・Wails 境界・強い制約のいずれも変えない。`.go-arch-lint.yml` の依存方向ルールは不変に保つ。 |

判定結果: 両方 N のため軽 task。`design-module` と `storybook-module` を bypass し、`preparation-module` → `implementation-module` → `finalization-module` で進める。

## 1. 認知複雑度ゲートの解除（gocognit `//nolint` の撤去）

各関数の現状複雑度と分割方針。外から見える契約は変えず、内部の段階・分岐を関数へ切り出す。

| 関数 | 場所 | 複雑度 | 分割方針 |
| --- | --- | --- | --- |
| `(*Engine).Run` | `internal/engine/engine.go`（`//nolint:gocognit` 付き） | 22 | 翻訳手続きの段階（固有名フェーズ・叙述文/定型句フェーズ・台詞フェーズ）を私的メソッドへ切り出し、`Run` は段階の連結と進捗集計に絞る。 |
| `DeriveTerms` | `internal/core/termderive/termderive.go` | 26 | 派生規則（姓名分割 two・二つ名前部 byname・短名）ごとに helper を切り出し、`add` クロージャの分岐を規則関数へ移す。純粋部は既に core/termderive へ別 package 化済み。 |
| `ParseTermXML` | `internal/core/termxml/termxml.go` | 23 | XML トークン走査の状態分岐（StartElement の種別判定・CharData 収集）をハンドラ関数へ切り出す。`engine-pure-rule-split` で `parseTermXML`→`ParseTermXML` に公開・移設済み。 |
| `BuildUsage` | `internal/core/termusage/termusage.go` | 16 | 文分割の外側ループと、文内トークンの大小・文頭判定を別関数へ切り出す。 |
| `(*App).pageRows` | `internal/api/app.go:569` | 26 | 叙述文・台詞・固有名の区間ごとの取得とカーソル境界判定を区間別 helper へ切り出し、`pageRows` は連結に絞る。 |
| `captureDBState` | `internal/harness/golden.go:40` | 21 | テーブル別ダンプを `{table, orderNote, query, format}` の spec スライスへ表化し、1 ループで回す。golden の出力順は spec の並びで保つ（テスト基盤側の負債）。 |

- 進め方: 1 関数ずつ分割 → `go test ./...`（harness golden 含む）緑 → `//nolint:gocognit` 撤去 → static で gocognit 0 を確認、を繰り返す。
- 不変: 切り出した helper を増やしても、各 helper 自身が 15 を超えないこと。golden を更新しないこと（更新が要るなら振る舞い変化なので停止して原因を見る）。

## 2. static baseline の解消

| linter | 件数 | 箇所 | 方針 |
| --- | --- | --- | --- |
| wrapcheck | 4 | `internal/engine/engine.go`（`InsertDerivedTerms`）、`internal/engine/persona_generate.go`（`UpsertPersonaCharacter`）、`internal/core/termxml/termxml.go`（xml `Token`）、`internal/core/termxml/termxml.go`（xml `DecodeElement`） | 返すエラーを `fmt.Errorf("...: %w", err)` で文脈づけて wrap する。store 呼び出しと外部 xml 呼び出しのいずれも日本語の文脈文を足す。 |
| gosec G304 | 2 | `internal/engine/engine.go`（`readXMLDir` の `os.ReadFile`）、`internal/lexicon/nrc.go`（`os.Open`） | 入力パスを `filepath.Clean` で正規化し、許容範囲を明示するか、意図的な外部入力であることを `#nosec G304` ＋ 理由コメントで限定許可する。読み込み元は利用者が選ぶ参照データのため、限定許可が妥当な見込み（要確認）。 |
| errcheck | 1 | `internal/lexicon/nrc.go`（`f.Close` 未チェック） | `defer func() { _ = f.Close() }()` へ替え、読み取り専用 open の Close エラーを明示的に無視する（harness の `dumpRows` と同じ作法）。 |
| staticcheck QF1001 | 1 | `internal/core/termderive/termderive.go`（De Morgan） | 条件式を De Morgan 則で書き換え、可読な肯定形にする。真理値は保つ。 |
| errorlint | 1 | `internal/core/termxml/termxml.go`（`==` での error 比較） | `err == io.EOF` を `errors.Is(err, io.EOF)` へ替える。 |

- 進め方: linter 種別ごとにまとめて直し、各種別の修正後に static を回して件数が減ることを確認する。
- 不変: 振る舞いを変えない修正に限る。`io.EOF` 判定や Close 無視は意味を変えないことを確認する。

## 含まない（除外範囲）

- engine 内部構造の作り変え・store/DB の作り直し（包括的リファクタ本体の別 scope）。本書は「複雑度の分割」と「lint 違反の解消」に限る。
- まだ純粋でないルールの純粋部括り出しと新規ユニット追加。これは backend-refactor-test-foundation ⑤の [`unit-test-plan.md`](../backend-refactor-test-foundation/unit-test-plan.md)（finalize 後は completed 配下）を入力にする別 scope。ただし §1 の `DeriveTerms`・`parseTermXML` の分割は⑤の括り出しと動きが重なるため、着手時に整合させる。
- frontend・C# 抽出器の変更。
- 汎用ボイス NPC（衛兵など）の話者未解決台詞への口調 fallback。本タスクの実画面確認（Innocence Lost - Quest Expansion.esp）で観測した別 scope。`docs/exec-plans/active/generic-voice-tone-fallback/plan.md` に backlog 化した。着手時に preparation-module で正式化する。

## 関連

- 前提タスク: backend-refactor-test-foundation（安全網と検証ハーネスの整備。完了）。
- 安全網: `.go-arch-lint.yml`（依存方向）、`scripts/lint/run-boundary-scan.sh`（責務違反走査）、`internal/harness`（合成 golden）。

## 実装・検証結果

### 変更ファイル（1 行 1 ファイル）

- `internal/core/termderive/termderive.go`: `DeriveTerms` を派生規則ごとの helper（`deriveShrt`・`deriveByname`・`deriveTwo`）へ分割し `//nolint:gocognit` を撤去。`safePair` の許可文字判定を肯定形へ書き換え QF1001 を解消。
- `internal/core/termxml/termxml.go`: `ParseTermXML` のレコード振り分けを `termAccum.collect` へ分割し `//nolint:gocognit` を撤去。`io.EOF` 比較を `errors.Is` へ（errorlint）。xml `Token`・`DecodeElement` の error を `fmt.Errorf("...: %w")` で wrap（wrapcheck 2）。
- `internal/core/termusage/termusage.go`: `BuildUsage` を `accumulateSentence`・`accumulateToken` へ分割し `//nolint:gocognit` を撤去。
- `internal/engine/engine.go`: `Run` を進捗集約 `runProgress` と段階メソッド `translateNarrations`・`translateLines` へ分割し `//nolint:gocognit` を撤去。`InsertDerivedTerms` の error を wrap（wrapcheck）。`readXMLDir` の `os.ReadFile` に G304 限定許可コメントを付与（gosec）。
- `internal/api/app.go`: `pageRows` を `pageBuilder` と区間メソッド（`countAll`・`fillNarrations`・`fillLines`・`fillPropers`）へ分割し `//nolint:gocognit` を撤去。早期 return の境界条件は不変。
- `internal/harness/golden.go`: `captureDBState` を `dbStateTables` 仕様列の 1 ループへ表化し `//nolint:gocognit` を撤去。各テーブルの format をパッケージ関数へ切り出す。出力順は仕様列の並びで維持。
- `internal/lexicon/nrc.go`: `f.Close` を `defer func() { _ = f.Close() }()` へ（errcheck）。`os.Open` に G304 限定許可コメントを付与（gosec）。
- `internal/engine/persona_generate.go`: `UpsertPersonaCharacter` の error を wrap（wrapcheck）。

### 最終検証（観測点）

- 複雑度: `sh ./scripts/lint/run-go-backend-lint.sh static` で gocognit 指摘 0。`grep -rn "nolint:gocognit" internal` は空（exit 1）。
- static: 同 static コマンドが `0 issues.` exit 0。当初 baseline 9 件（wrapcheck 4・gosec 2・errcheck 1・staticcheck 1・errorlint 1）を全解消。
- 非劣化: `go test ./internal/harness/...` の golden 一致。`npm run verify:backend`（go test ./... ＋ arch ＋ boundary）が exit 0。

### 作業中に混入し解消した違反

- termxml の wrap 文を当初「String 要素の…」と書き、先頭大文字で staticcheck ST1005 を 1 件混入させた。先頭を「XML の…」へ直して解消した（最終 static は 0 issues）。

## 実画面確認

- 確認環境: `npm run dev:wails:run`（http://localhost:34115）、実 LLM 別マシン `http://192.168.0.226:1234`（モデル hy-mt2-7b）、plugin は `dictionaries/Data/Innocence Lost - Quest Expansion.esp`（小規模 plugin で全フローを短時間確認）。
- 確認結果: 抽出 → 翻訳 → 結果一覧まで通過。叙述文・定型句・台詞が訳出され（engine.Run の段階分割が動作）、結果一覧 197 件が叙述文→定型句→台詞→固有名の順でページング表示（pageRows の区間分割とカーソル境界が動作）。固有名は本文へ一貫機械置換（termxml・termderive・termusage の派生が動作、master_term に派生反映）。主要 NPC 台詞に口調付与（persona 生成が動作）。
- 観測した別 scope: 汎用ボイス NPC（衛兵）の話者未解決台詞が口調なしになる既存挙動を観測。本タスク無関係（口調生成・話者連関・voice fallback のロジックは未変更、golden 一致）。`docs/exec-plans/active/generic-voice-tone-fallback/plan.md` に backlog 化した。

## finalization

### 正本化判断

- `docs/architecture.md` 反映: 不要。軽 task（preparation-module で画面 N・architecture N）。関数内の helper 切り出しと lint の局所修正に限り、層構成・依存方向・Bootstrap・Wails 境界・強い制約は不変。
- 人間承認: 不要（正本反映なし）。

### 作業 commit

- commit: `a055e3fc`（branch `claude/backend-violation-cleanup`）。
- 変更ファイル: コード 8（termderive・termxml・termusage・engine・persona_generate・api/app・harness/golden・lexicon/nrc）＋ 本 plan.md ＋ 新規 backlog plan。
- remote は未変更。

### local merge

- command: `git merge --no-ff claude/backend-violation-cleanup`（target: `master`）。
- merge commit: `812dfbeb`。conflict なし。

### merge 後検証

- `npm run verify:backend`: exit 0（go test ./... ＋ arch ＋ boundary、いずれも OK）。
- `sh ./scripts/lint/run-go-backend-lint.sh static`: `0 issues.` exit 0。gocognit 指摘 0。
- `grep -rn "nolint:gocognit" internal`: 空。
