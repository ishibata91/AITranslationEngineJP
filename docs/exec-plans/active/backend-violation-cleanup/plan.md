# backend-violation-cleanup（複雑度ゲートのリファクタと static baseline の解消）

## 依頼要約

backend-refactor-test-foundation（完了）が安全網として整えた検査が表面化させた残債を解消する。具体的には、認知複雑度ゲート（gocognit）で baseline 除外した 6 関数を分割して `//nolint:gocognit` を外し、golangci-lint static の既存 baseline 11 件を解消する。テスト基盤（golden harness・arch-lint・境界走査）は不変に保ち、非劣化を保証したうえで内部を直す。

- 分岐元 branch: master（backend-refactor-test-foundation の merge 後）
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

## 関連

- 前提タスク: backend-refactor-test-foundation（安全網と検証ハーネスの整備。完了）。
- 安全網: `.go-arch-lint.yml`（依存方向）、`scripts/lint/run-boundary-scan.sh`（責務違反走査）、`internal/harness`（合成 golden）。
