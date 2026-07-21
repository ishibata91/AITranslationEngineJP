# 変更・判断履歴

正本（`requirements.md`、`system_requirements.md`、`architecture.md` など）には現在の状態だけを書く。
「なぜ変えたか」「何を落としたか」などの判断履歴は本ファイルに残し、正本へ混ぜない。
新しい entry を上に追加する。1 entry は date 見出しで区切る。

## 2026-07-21 xAI batch 翻訳の UI 導線を接続

### 変更

- `frontend/src/gateway/translation-gateway.ts`: 生成済み Wails 公開面 `GetXAIModels`・`SubmitBatchTranslation`・`RefreshBatchTranslations` のラッパ（`fetchXaiModels`・`submitBatchTranslation`・`refreshBatchTranslations`）を追加。generated 呼び出しは gateway 境界に閉じる既存方針を踏襲。
- `frontend/src/ui/screens/translation-run/TranslationRunScreen.svelte`（と view・presentation・fixtures・stories）: 翻訳実行画面へ配送方式選択（同期 / xAI（batch））を追加。xAI 選択時はエンドポイント欄と取得補足を xAI 用へ切替、実行ボタンを「送信」「反映」の 2 ボタンへ差し替え、送信後の案内と反映注記を表示。Storybook 人間レビュー承認済み。
- `frontend/src/ui/screens/translation-run/TranslationRunContainer.svelte`: 配送方式 state と xAI 用ハンドラ（送信・反映・xAI モデル取得）を追加。方式切替で既定エンドポイント（同期 `http://127.0.0.1:1234` / xAI `https://api.x.ai`）とモデル取得先を切り替える。接続情報は永続化せず、送信・反映のたびに渡す。

### 判断

- 進行状況の可視化は最小（frontend のみ）で確定（人間確定）。送信・反映と送信後案内だけを出し、batch の進行段や pending 表示のための backend 状態取得公開面は足さない。進行は反映で結果一覧が変わるかで判断する。
- 反映は plugin を取らない global 操作（進行中の全 batch をまとめて確認・反映する backend 公開面の非対称）。画面文言で「進行中の batch をまとめて確認」と示し、反映後は開いている plugin の結果一覧を読み直す。
- 反映は手動ボタンのみ。接続情報を永続化しないため起動時の自動反映はせず、backend の常駐ポーリング非採用にそろえてポーリングもしない。
- backend の翻訳ロジック・provider・永続・Wails 境界は変えず、既存公開面を frontend gateway から消費するだけのため、`docs/architecture.md` の層・依存・境界は不変（正本反映なし）。
- 課金を避けるため、送信・反映の実 xAI 疎通と `json_schema` の batch 実挙動は手動 e2e（人間）に限定。自動検証は `npm run test:frontend` と実 app での配送方式切替・ボタン出し分けの目視までとし、送信・反映ボタンは押していない。
- 根拠となる作業計画: `docs/exec-plans/completed/xai-batch-ui/`（`design.md` の AS-IS/TO-BE 図・検討事項の解消記録）。

## 2026-07-20 xAI batch 翻訳（非同期の大量翻訳）に対応

### 変更

- `db/migrations/0013_batch_translation.sql`: batch 進行本体 `batch_translation`（plugin と 1 対 1・進行段・固有名/本文の外部 batch ID）と送信行対応 `batch_request`（custom_id `種別:id` で外部 batch と翻訳対象行を結ぶ）を追加。
- `internal/provider/batch.go`・`xai_batch.go`: 2 つ目の多態 port `BatchTranslator`（送信・状態確認・結果取得）と xAI 実装 `XAIBatch`（JSONL file upload 方式、`/v1/batches` 系）を追加。同期 `Translator`（`OpenAICompatible`）は不変。失敗種別は同期と同じ番兵で表す。
- `internal/core/batchplan/`: batch 管理の純粋決定規則（custom_id 割り当て・進行段遷移・結果適用・再送信可否）を新設。カバレッジ 100%。
- `internal/engine/engine.go`・`batch.go`: 同期本文フェーズの「リクエスト構築」と「1 件適用判断」を振る舞い不変で共有関数（`composeBodyPrompt`・`batchplan.DecideApply`）へ抽出し、batch 反映と共有。`BatchRunner`（送信/反映のオーケストレーション・薄いシェル）を追加。
- `internal/store/batch_translation.go`・`target_plugin.go`: batch 永続アクセスと、plugin 削除の手続き DELETE への batch テーブル追記（子 → 親順）。
- `internal/api/app.go`・`bootstrap.go`: Wails 公開面へ `SubmitBatchTranslation`・`RefreshBatchTranslations`・`GetXAIModels` を追加し、`XAIBatch`・`BatchRunner` を配線。
- `docs/architecture.md`（§2 図・§3・§4・§4.1・§7・§8）・`docs/er.md` §2・`.go-arch-lint.yml`: 2 つ目の port と `batchplan` の構造を反映。

### 判断

- xAI は batch 専用の非同期配送とし、同期経路を持たせない（人間確定）。同期 `Translator` の振る舞いは変えず、batch は 2 つ目の port として並置する。「同期の翻訳本体は変えない」は振る舞い不変の抽出（behavior-preserving refactor）を許容する意味とし、既存 engine テストで回帰を担保した。
- 固有名は本文へ機械置換で注入するため本文より先に確定する必要がある。よって固有名 batch → 本文 batch の 2 段逐次とした（設計レビュー回答 B。案 A 固有名同期・案 C 横断辞書のみ注入より一貫性を優先）。
- 状態確認は常駐プロセスもバックグラウンドポーリングも作らず、起動時・画面操作の時点だけで行う（plan 確定）。接続情報は永続化せず、送信・反映のたびに UI から渡す。
- batch 管理の決定規則を単一の純粋モジュール `batchplan`（100% カバレッジ）へ寄せ、結果適用を同期と共有することで、外から見て同期と batch の結果が区別できないことを不変条件にした（依頼の「シングルモジュール 100%・外から変わらない」に対応）。
- 課金を避けるため、自動テストは実 xAI API に触れない。単体（純粋核・fake HTTP）と結合（fake batch プロバイダで送信→永続→反映→書き戻しの 2 段連鎖・再起動相当・sync との dest 一致）で網羅し、実 API 疎通と `json_schema` の batch 実挙動は手動 e2e に委ねた（人間確定）。
- 期限切れ（`expires_at`）・一部失敗の未反映行は未訳のまま残し、利用者の再送信で回収する（同期の skip と同思想）。送信 HTTP 失敗で外部 ID が空のまま残った半端な進行は、再送信で reset して復旧する（バグレビュー指摘の対応。再送信拒否は「現段の外部 batch ID が非空」の待機中進行に限る）。
- `batch_request.kind`/`row_id` は custom_id と重複する保持で、現状の反映は custom_id から行を再導出する（対応表は冗長）。列を落とすか反映を対応表ベースへ寄せるかは後続判断とした。
- UI 導線（xAI 選択・送信・反映）は本 task の scope 外（backend 先行、人間確定）。Wails 公開面までを足し、frontend は変更しない。
- 根拠となる作業計画: `docs/exec-plans/completed/gemini-xai-batch-translation/`（`design.md` の実装方針・AS-IS/TO-BE 図・検討事項の解消記録）。

## 2026-07-20 翻訳 run の失敗を種別で切り分け（known-issues #7 解決）

### 変更

- `internal/provider/openai_compatible.go`: 失敗分類を追加。番兵エラー `ErrResponseUnreadable`（応答エンベロープの decode 失敗・`choices` 無し）と `ErrServerTransient`（非 200 の 429・5xx）、および非 200 status を分ける純粋関数 `statusSkippable`（429 と 5xx を skippable、その他 4xx を fatal）を新設。`Translate` の各失敗経路をこの分類でラップ。
- `internal/engine/engine.go`: 本文フェーズ loop（`translateNarrations`・`translateLines`）の継続条件を、`ErrStructuredParse` 限定から skippable 集合の判別（`lineSkips.record`）へ拡張。skippable は該当行を未訳のまま飛ばし、fatal・未知の失敗は従来どおり `return` して `Run` を中断。飛ばした行を理由別（`structured_parse_failed`・`response_unreadable`・`server_transient`）にフェーズ末で集約ログ。
- `internal/provider/openai_compatible_test.go`・`internal/engine/engine_test.go`: 分類の純粋網羅、非 200 分類、応答読み取り失敗分類、skippable は skip して続行・fatal は以降を処理せず中断、の各テストを追加。
- `docs/known-issues.md`: 課題7「翻訳リクエストの失敗で翻訳 run 全体が止まる」を解決済みとして削除。

### 判断

- known-issues #7 が未確定として残していた「基盤失敗を fail-fast で止めるか、失敗行を飛ばして続行するか」の方針を、失敗種別で 2 分する形で確定した（人間確定）。fatal（run を止める）＝通信失敗（断・timeout・context 中断）・4xx のうち 429 以外（401/403 認証・400 不正など）・リクエスト生成失敗・分類できない未知の失敗。skippable（その行を未訳のまま飛ばす）＝非 200 の 429・5xx・応答エンベロープの decode 失敗・`choices` 無し・content の空/スキーマ違反（既存 `ErrStructuredParse`）。
- 未知の失敗を既定で fatal にしたのは安全側に倒すため。将来 provider に新しい失敗経路が増えても、黙って飛ばさず run を止めて人間に見せる。
- 帰結として、通信の一時的な瞬断・timeout でも run が止まる（「通信断は止める」の選択の帰結）。リトライ機構は本 task の scope 外とした（人間確定）。飛ばした行は再実行で拾える（`Run` は未訳行だけを対象にするため冪等）。
- 観測点は単体テストに固定した。1 リクエスト単発の mid-run 失敗は UI から誘発できない（故障注入がプロダクトコード変更を要し、修正フローの安全境界に反する）ため、実画面での再現確認は対象外とした。
- 根拠となる作業計画: `docs/exec-plans/completed/translate-run-failure-isolation/`（`investigation.md` の確定原因、`design.md` の失敗分類表・AS-IS/TO-BE 図）。

## 2026-07-20 会話グラフ 3 task を棄却しバルク翻訳機構を revert

### 変更

- `dialogue-graph`・`dialogue-flow-context`・`bulk-line-translation` の 3 plan folder を `docs/exec-plans/rejected/` へ移動。前 2 者の `plan.md` を棄却記録へ書き換え、`bulk-line-translation/plan.md`（一度 completed）へ棄却注記を追記。
- 棄却根拠（PNAM/TCLT の実測）を `dialogue-graph/plan.md` に集約。
- バルク翻訳機構（commit `8599cd1a`）を `git revert` で撤去（`internal/core/chunking/` 削除、`provider`/`harness` の batch 撤去、`engine.go` のチャンク経路撤去、`app.go` の TokenBudget 撤去、frontend の予算 UI 撤去、bulk テスト削除）。追随修正 `a1e62436`（翻訳対象を選んだ plugin へ絞る）は残す。plugin 絞り込みテスト `TestRunScopesToTargetPlugin` を `engine_test.go` へ移設。
- `known-issues.md` の no7（1 リクエスト失敗で run 全体が止まる）は、bulk 用語を除き単一経路の課題として残す。

### 判断

- 会話往復グラフ構築の費用対効果を、抽出器へ使い捨ての測定を差して実 plugin で測り直した（測定コードは revert 済み）。結果、翻訳文脈は 2 層に分かれ費用対効果が非対称と判明した。
- **PNAM は会話の前後リンクではない**（解決分は 100% 同一 DIAL 内・Skyrim.esm は 0%）。会話の流れは TCLT のみが担い、循環つき合流グラフになる（合流は被参照 DIAL の約 40%、循環は inigo で 302 本・深さ 59）。木でも DAG でもない。
- **層1**（NPC 台詞が答える DIAL:FULL＝選択肢文。r1 で 100% 供給・グラフ不要）は安価で頑健。**層2**（TCLT 逆算の手前連鎖。対象 12〜25%・循環・合流）は構築コストに品質寄与が見合わない。よって層2 のグラフ構築（本 2 task）を棄却した（人間確定）。
- `bulk-line-translation`（slice1）はバルク機構を土台にする会話文脈 task 群が棄却されて存在意義を失うため、コードを revert し台詞 1 行ずつの従来動作へ戻した（人間確定）。`a1e62436` の plugin 絞り込みはバルクと独立に有用なため残す。
- e7（`line_sequence`）は未実装のまま `known-issues.md` に残す。層1 が必要になった時はグラフを作らず別途軽量に扱う。

## 2026-07-20 翻訳前区間のパフォーマンス改善（抽出子 DLL 化・進捗表示・Normalize memoize）

### 変更

- `internal/api/app.go`: C# 抽出子の起動を `dotnet run --project` から publish 済み DLL 直実行（`dotnet <extractor.dll>`）へ。`ExtractorConfig.ProjectPath`→`DLLPath`、`buildExtractorArgs` を DLL 形へ、DLL 不在時の明示エラー。`ProgressEvent` に段内サブ段 `step` を追加し、`RunExtractAndTranslate` の翻訳前区間を単一 emit から4境界（extract/derive/reference/ingest の各段直前）の step 付き emit へ。
- `internal/bootstrap/bootstrap.go`・`cmd/goldcap/main.go`: `ExtractorConfig` を DLL パスへ追随（フラグ `-extractor-project`→`-extractor-dll`）。
- `scripts/dev/run-wails.sh`: `wails dev` の前に `dotnet publish tools/extractor` を1度実行（ビルド起動点を起動 script へ一本化）。
- `tools/extractor/RecordDataIndex.cs`: `Normalize` を FormKey 単位で memoize（null 結果も cache）。`tools/extractor.Tests/RecordDataIndexTests.cs` を新設。
- frontend `translation-run-{view,presentation}.ts`・`TranslationProgress.svelte`・`TranslationProgress.stories.ts`: 抽出段の4サブ段見出し（台詞を抽出/辞書を準備/既存訳を取り込む/翻訳対象を仕分け）を表示。
- `docs/known-issues.md`: 課題6「配布 app での C# 抽出子の同梱未整備」を追加。

### 判断

- 申告（大きめ plugin の抽出が約6分）を実測で切り分けた結果、抽出処理そのもの（C# 抽出子＋Go 後段）は available データで高速（Outfit 約3.5s、USSEP 約2.9s）で、6分は再現しなかった。最有力原因は初回 `dotnet run` の NuGet restore＋ビルド崖を app 本体と誤認したケース（確定できず仮説）。よって体感の遅さと固まりを生む構造（毎回のビルド評価・崖、翻訳前区間の無音）を恒久的に直す方針にした。
- 抽出子の DLL 化はビルド起動点を dev 起動 script へ一本化する（Go フォールバック build は持たせない、人間確定）。
- `Normalize` memoize は master 依存が数十件の mod への備え（available データでは非線形が顕在化しない）として今回含める（人間確定）。抽出結果は不変（`CountParityTests` で担保）。
- 配布 app 対応（self-contained publish・同梱・配布フロー新設）は規模が大きいため今回やらず、`known-issues.md` 課題6 に残して `feature-workflow` の別 task に回す（人間確定）。
- `architecture.md` 反映は不要と判断。抽出子は DLL 化後も子プロセス起動のまま、進捗も同じ runtime events で、層・依存・Wails 境界を変えないため。

## 2026-07-20 設計と調査を責務分離し design.md を両フロー共通の 1 本へ

### 変更

- `docs/exec-plans/templates/task-folder/design.md`: 両フロー共通の 1 テンプレートへ。「どう実装/どう直すか（実装方針＝AS-IS→TO-BE）＋検討事項」だけを持つ（一度フロー別に分けた `design-feature.md`／`design-fix.md` は削除）。
- `docs/exec-plans/templates/task-folder/investigation.md` を新設（修正フローだけが作る。観測済み問題・画面再現確認・原因仮説・観測ログ検証・確定原因）。
- `.claude/skills/fix-workflow/SKILL.md`: `修正方針判断` 成果物を `調査`（investigation.md）と `設計`（design.md）へ分割。「修正実行入力」を「実装への引き継ぎ」へ改名。
- `.claude/skills/feature-workflow/SKILL.md`: design.md テンプレート参照を単一へ。
- `.claude/skills/implementation-module/SKILL.md`・`coding-protocol/SKILL.md`・`fix-decision/SKILL.md`: 「修正実行入力」を「実装への引き継ぎ」へ、投入先を investigation.md／design.md へ更新。
- `docs/exec-plans/templates/task-folder/README.md`・`docs/ai-operations-workflow.md`: investigation.md 追加と design.md 単一化を反映。

### 判断

- 再現確認・原因究明は設計（design）の責務でなく調査（investigation）の責務（人間指摘）。design.md は「どう実装/どう直すか」だけに絞る。この境界だと feature と fix の design.md は同形になるため、直前にフロー別へ分けた design テンプレートを 1 本へ戻す。
- 調査成果物は物理ファイルで分離する（案X）。task フォルダに `investigation.md` を新設し、fix だけが作る。feature は investigation を作らない。
- 「修正実行入力」は造語で用語規約（言葉の発明禁止）に反するため「実装への引き継ぎ」へ改名。過去記録（`docs/exec-plans/completed/`）は書き換えず現行契約ファイルのみ変更。
- plan.md はフロー別に割らず 1 本のまま（割るほど書くことが無い）。

## 2026-07-20 plan.md に「やること」を復活

### 変更

- `docs/exec-plans/templates/task-folder/plan.md`・`README.md`: plan.md の保持内容を「branch 情報・やらないこと」から「branch 情報・やること・やらないこと」へ戻す。
- `.claude/skills/feature-workflow/SKILL.md`・`.claude/skills/fix-workflow/SKILL.md`: plan.md 節と入口条件に「やることの要点」を追加。
- `docs/ai-operations-workflow.md`: plan.md の説明を同内容へ更新。

### 判断

- 直前の再編（2026-07-20 オーケストレーター化）で plan.md から「やること」を外した結果、plan.md 単体では何の task か分からず changelog を検索しないと内容が追えなくなった（人間指摘）。plan.md に「やること」を戻す。
- 「やること」は人間の依頼内容をそのまま要約した粒度に限定する。設計判断・手段選定・原因仮説は `design.md` に置き、plan.md には持たせない。

## 2026-07-20 AI 運用ワークフローをオーケストレーター化・`fork` 委譲・AS-IS/TO-BE 図へ再編

### 変更

- `.claude/skills/preparation-module/` を廃止（削除）。入口責務（branch・plan.md 固定）を各フローの入口へ統合。
- `.claude/skills/design-module/` を `feature-workflow` へ改名し、新規実装フローの入口オーケストレーターへ再定義。責務を実装方針・AS-IS→TO-BE の変更内容・検討事項・設計レビューへ整理し、実装範囲とテスト設計を外した。
- `.claude/skills/investigation-module/` を `fix-workflow` へ改名し、修正フローの入口オーケストレーターへ再定義。自前で branch・plan.md を固定し、修正方針判断を design.md に持つ。
- `.claude/skills/storybook-module/`・`.claude/skills/implementation-module/`: 作業本体を、親の文脈とモデルを継承する`fork`へ委譲する形に変更。軽 task の design bypass を廃止。
- `.claude/skills/coding-protocol/`・`finalization-module/`・`conflict-resolver/`・`fix-decision/`・`.claude/agents/conflict_resolver.md`: 旧モジュール名参照の更新、plan.md への記録指示を作業結果返却または changelog へ振り替え。
- `.claude/skills/presentation/`: 色で塗り分ける差分図を廃止し、AS-IS→TO-BE の 2 図へ変更。
- `CLAUDE.md`: 実装フローの原則を「`fork` 委譲は可、`fresh` への分割は不可」へ全面改訂。agent を `fresh`（引き継ぎなし）と `fork`（引き継ぎあり）の 2 語で定義し、model 既定は `fork` が親モデル（opus 含む）継承可、`fresh` は haiku/sonnet 既定へ。言葉の発明禁止を用語規約に追加。
- `docs/exec-plans/templates/task-folder/`: plan.md を最小形（branch 情報・やらないこと）へ、design.md を新設、旧参照 detail-spec-diff.md を削除、README を更新。
- `docs/ai-operations-workflow.md`: 上記の新構成へ書き直し。

### 判断

- `fork` 委譲の目的は並列化ではなく、本体セッションの文脈を実装詳細で汚さないこと（人間指定）。`fork` は全文脈を継承するので「1 文脈で書く」原則と両立し、opus 継承も許可する。
- 軽 task も design を通す方針に変えた。入口を task の軽重で bypass しない。
- 修正フローの入口重複（feature-workflow と fix-workflow が各々 branch・plan.md を固定）は、今回は許容して運用で試す（人間判断）。
- plan.md は branch 情報とやらないことだけを持つ最小アーティファクトへ。判断履歴・検証結果は plan.md に残さず、changelog または作業結果へ寄せる。

## 2026-07-19 実行時タグ保護を「AI に生タグを見せる」方式へ変更

### 変更

- `internal/core/runtimetag/runtimetag.go`: API を刷新。`Mask`/`Restore`（退避と復元）、`Tags`/`HasTag`（生タグ抽出・有無）、`CountMissing`（AI 出力の欠落照合）、`GuardInstruction`（生タグ保持の指示文）。旧 `Unmask`/`Placeholders`/`HasPlaceholder` を廃止。
- `internal/engine/engine.go`: `prepareSource` を新設。実行時タグは機械置換（`dict.Apply`）の間だけ退避し、直後に生タグへ戻してから AI へ送る。AI 出力は `CountMissing` で照合し、タグ欠落行は未訳のまま残す（従来どおり）。
- `internal/core/prompt/prompt.go`: タグ保護指示の付与判定を退避トークンから生タグ（`HasTag`）へ変更。
- `internal/api/app.go`: 実プロンプト再構成も同順（退避→機械置換→生タグ復元）にして送信と一致させる。
- `internal/harness/provider.go`・`oracle_test.go`: fake は生タグを保持する形へ、オラクルは「user に生タグ・system に保護指示」を照合。

### 判断

- 前方式（タグを不透明トークン `⟦0⟧` へ退避したまま AI へ送る）は、トークンから意味が消えて弱いモデルが「無意味な記号」と見なし落としやすい弱点があった（人間指摘）。方針を「退避は機械置換から守る間だけにし、AI には意味の分かる生タグ `<BribeCost>` を見せる」へ変更した。
- 退避の目的（辞書機械置換がタグ内部の語を書き換えるのを防ぐ）は `Mask`→`dict.Apply`→`Restore` の一時退避で維持。AI へは生タグ＋保護指示を渡し、出力にタグが原形で残るかを照合、欠落行は未訳保留（安全網は不変）。
- 実 LLM で確認: 前方式で `hy-mt2-7b` が 2 回とも落とした `<BribeCost>` 2 行が、生タグ方式では 2 つとも原形保持で翻訳された（`（<BribeCost>ゴールド）`）。

## 2026-07-19 既存訳の完全一致置換と実行時タグ保護（known-issues 項目7・8）を実装、項目4 を close

### 変更

- `internal/core/runtimetag/`（新規）: 本文中の実行時タグ（`<Alias=...>` 等 `<...>`）を退避 `Mask`／復元 `Unmask`（欠落数を返す）する純粋ルール。保護指示 `GuardInstruction`・退避有無 `HasPlaceholder`。単体カバレッジ 100%。
- `internal/core/prompt/prompt.go`: `ComposePrompt` が退避トークンを持つ本文に system へタグ保護指示を足す（送信・再構成で自動一致）。
- `internal/core/termxml/reference.go`（新規）: xTranslator 英日 XML から record 単位の既存訳（`ParseReferences`・`ReferencesFromFiles`）を抽出。
- `db/migrations/0012_reference_translation.sql`（新規）: 既存訳の照合表 `reference_translation`（rec, field, source, dest, UNIQUE(rec,field,source)）。
- `internal/model/reference_translation.go`・`internal/store/reference_translation.go`（新規）: model と Insert/List。
- `internal/engine/reference.go`（新規）: `LoadReferenceTranslations`（XML→表）、`referenceIndex`（表→照合 map）、`ReferenceStore`。
- `internal/engine/engine.go`: 翻訳ループで既訳完全一致は AI を呼ばず `statusTranslated`(1) で流用。タグは Mask→翻訳→Unmask し、欠落行は未訳（status=0）のまま残す。欠落は `slog`（result=skipped）で phase 集約。
- `internal/api/app.go`: Run で `LoadReferenceTranslations` を結線。結果画面の実プロンプト再構成を送信と同順（Mask→dict.Apply）に合わせる。
- `main.go`: `slog` の JSON handler を composition root で設定。
- `internal/harness/`: `RecordingProvider` が退避トークンを保持。fixture にタグ入り叙述文・既訳一致台詞を追加。integration spec 2 件（`runtime-tag-preserved`・`existing-translation-reused`）とオラクル。
- `.go-arch-lint.yml`: `runtimetag` component 登録と依存追加（prompt・engine・api）。
- `docs/known-issues.md`・`docs/roadmap.md`: 実装済みの項目4（xTranslator 書き出し）と実装した項目7・8 を除き、番号と相互参照を整合。
- `docs/er.md`: `reference_translation` を実現テーブルへ追加。

### 判断

- 項目4（xTranslator 書き出し）は記述が古く、実際は backend（`termxml.MarshalStrings`・`engine.ExportXTranslator`）から frontend の書き出しボタンまで実装済みと確認し、close した。
- 項目7 の照合キーは「plugin・form_id・原文」でなく **(rec, field, source)** にした。供給源の xTranslator XML が FormID を持たず、公式訳は base ゲーム由来で対象 plugin（Mod）と plugin が一致しないため。同一原文の既訳を対象横断で流用する。
- 項目8 のタグ欠落時の対応（人間決定）: 落ちたタグの自動差し戻しは差し込み位置不明で不可のため、欠落行は壊れた訳を確定させず未訳のまま残す（再実行で再翻訳）。あわせてタグを持つ本文には保護指示をプロンプトへ注入して欠落自体を減らす。per-record の画面上の拾い上げ・修正は結果画面編集（項目4＝別課題）へ委ねる。
- 実 LLM e2e（`Innocence Lost - Quest Expansion.esp`、LM Studio `hy-mt2-7b`）で、既訳 67 件流用（status=1）・タグ保持・弱いモデルでの欠落 2 行の未訳保留と集約ログ発火を目視。
- architecture.md 反映: 不要と判断（`runtimetag` は既存 core 構造の新 leaf で層・依存方向・Wails 境界は不変。enforced な `.go-arch-lint.yml` のみ更新）。

## 2026-07-19 合成パイプラインの結合オラクル（C# 抽出＋Go 翻訳の継ぎ目）を追加

### 変更

- `test-oracle/specs.json`: 共有オラクルを継ぎ目のみへ刈り込んだ。`extraction` 16 件（C# 抽出、うち 3 件は master/localized が要り既存実 esm 単体へ委任＝`coverage: existing-unit`）＋ `integration` 6 件（Go 翻訳の継ぎ目）＝ 22 件。
- `tools/synthetic-fixture`: master 無し自己完結の合成 esm を Mutagen で組む独立生成器（`SyntheticEsmBuilder`＋Program）を新設。成果物 `test-oracle/fixture/Synthetic.esm` を commit。生成はテスト・build から切り離す。
- `tools/extractor.Tests`: C# 抽出オラクル（`OracleExtractionTests`＝13 件、`OracleInput`・`OracleSpecs`・`[Oracle]` 属性）を追加。合成 esm を実抽出し抽出結果／staging を照合。網羅番人は `[Oracle]` id 集合と specs.json の一致を見る。
- `internal/harness`: Go 結合オラクル（`oracle_test.go`＝6 件）を追加。`SyntheticRun` を 1 回通し `Capture`・最終 DB を照合。合成 fixture へ感情 staging（`extracted_info_emotion` の Fear）とフルネーム固有名台詞を追加。golden 文字列比較（`TestSyntheticNonRegression`・`testdata/synthetic.golden`）を撤去し、決定性テストは残した。
- `tools/extractor.Tests/CLAUDE.md`・`internal/harness/CLAUDE.md`: オラクルテストの書き方（1 オラクル 1 関数・id を引ける・AAA・独立・given は入力側）をテストフォルダ直下へ固定。

### 判断

- 結合オラクルが守るのは「単体で守れない継ぎ目」だけ、と定義を訂正した。当初は specs.json の 54 件（段×属性）を 1 件ずつテスト化し、単段で純粋に閉じるルール（口調閾値・stoplist 判定・固有名派生・役割語引き）まで結合ハーネスで照合しようとして無限に膨らんだ。単段ルールは core package の既存単体が守るためオラクルから外し、go 段を継ぎ目 6 件（件数保存・未知除外・固有名一貫・感情結線・話者結線・stoplist 一貫）へ刈り込んだ。
- 入口→出口の範囲はツール単位に取った（C#: esm→抽出 staging、Go: seed→翻訳出力）。C# 抽出出力を Go へ流し込む派生 seed 連結は、master 無し実 esm では作れない given（未知 REC:FIELD・別 plugin・master 辞書内容）があり、かつ既存の手組み fixture より複雑になるため却下した。
- e2e と呼ばない（実 LLM・UI を通らない）。stage 名は `integration` に統一した。
- テストの書き方の固定は root CLAUDE.md でなくテストフォルダ直下の `CLAUDE.md` に置いた。そのフォルダで作業する時に必ず読まれるため。

## 2026-07-05 辞書に無い漏れ語の候補検出（決定的ヒューリスティック）を研究の結果、不採用として revert

### 変更

- merge commit `b23f83cb`（候補検出の純粋ルール `internal/core/mention/candidate.go`・prose adapter・評価ハーネス `cmd/poc-missing-term`・関連 docs 更新）と completed 移動 `ca5c3e52` を revert した。実装・評価の全記録は git 履歴（作業 branch `claude/dictionary-missing-term-detection`、作業 commit `0b36d179`）に残る。

### 判断

- 研究の実測: LLM 不使用の決定的手法（大文字ヒューリスティック＋用法分布＋prose の固有表現・品詞解析）で、held-out 3 plugin に対し再現率 95.4%・重複 0・決定性ブレ 0 を確認した。ただし真の精度（人手ラベリング）は dev 約 63%・held-out 約 54% が天井で、精度を上げる選別はいずれも再現率を 95% 未満へ下げた。
- 不採用の理由: 実抽出の出力（inigo.esp 1,031 語・Skyrim.esm 12,515 語、約半分が断片・間投詞等のノイズ）を目視した結果、人間レビュー前提でも候補一覧の品質と量が製品に載せる水準にないと判断した。
- 代替の方向（未着手のメモ）: 事前の候補抽出でなく、翻訳結果側で同じ原語が別の訳になった行だけを突き合わせる「事後の訳揺れ検出」の筋が良い可能性がある。実際に揺れた語しか出ないため精度問題が構造的に消える。着手する場合は known-issues 1・2 番の統合 task の入口で再検討する。

## 2026-07-04 機械置換辞書の誤爆対策（一般語 stoplist）

### 変更

- `assets/stopwords-en.txt`: stopwords-iso の英語 stopword リスト（MIT、1297 語）を追加した。上流と byte 一致を保ち、出典・取得日・sha256・ライセンス全文は `assets/stopwords-en.LICENSE` に記録した。
- `internal/core/dictionary`: 供給源選別の純粋ルール `Stoplist`（`ParseStoplist`・`Blocks`）を追加した。判定は原語全体を小文字化した一語一致で、複数語の固有名には当てない。
- `internal/engine`: 共通供給点 `translationVocabulary` を新設し、機械置換辞書（`LoadDictionary`）と言及語彙（`mentionDetector`）が同じ選別を通る形にした。
- `internal/bootstrap`・`cmd/goldcap`・`internal/harness`: stoplist の読み込みと注入を配線した。合成 harness は最小リスト（yes・no）を使い、実ファイルの内容変化から切り離す。
- `internal/harness`: 合成 fixture へ FACT:MNAM の "Yes"・"No" と文頭 Yes/No の本文を追加し、golden の DB ダンプへ `narration_mention`・`line_mention` の観測点を追加した。
- `docs/known-issues.md`: 6 番（機械置換辞書の誤爆対策の残り）を新設した。`docs/roadmap.md`: 7 番として追記し、1 番の古い注記（言及テーブル未実装）を現状（残りの言及関連）へ訂正した。

### 判断

- stoplist の実装位置: 新 package を作らず置換コア `internal/core/dictionary` へ統合した（利用者指示）。供給源選別は辞書の意味の一部であり、独立 package は構成と arch-lint 設定を無駄に増やすため。
- 誤爆の実態の訂正: inigo.esp の "Yes"・"No" は FACT:FULL でなく FACT:MNAM（階級称号）由来だった（前 plan の記録を実測で訂正）。誤爆の経路（固有名 box → AI 訳 → 機械置換）は変わらず、対策の設計へ影響しない。
- 内部フラグ勢力の除外（第 2 層）は不成立で停止: 候補基準 FACT の Hidden from PC flag は、Yes/No の供給源勢力に立っておらず、Skyrim.esm では hidden=1 の 163 勢力中 154 件に master_term 既訳があり "Thieves Guild" 等の実在名まで落ちる。抽出器（tools/extractor）は変更しない。観測記録は completed plan（dictionary-false-positive-guard）にある。
- 意図的な出力変更の確認: 実データ（inigo.esp）で分岐元 golden と比較し、差分を単語レベルで機械分類した。すべて stoplist 語の還元（訳031→Yes 約 205 箇所、訳030→No 約 172 箇所、master_term 側の Hello→やあ・Mine→鉱山・Fire→炎・Turn→吸血鬼化・Down→下へ）で、translated_count は 8803 で不変。stoplist 外の固有名の置換・言及（Riften 等）は維持された。

## 2026-07-04 言及テーブル（e3/e4/e5）の実装

### 変更

- `db/migrations/0008_mention.sql`: `narration_mention`（e4 叙述文→固有名の言及）・`line_mention`（e5 台詞→固有名の言及）・`narration_described`（e3 叙述文→説明対象の固有名、0..1）を新設した。
- `internal/core/mention`: 言及検出の純粋 package を新設した。照合規則は機械置換（`internal/core/dictionary`）と同じ（貪欲最長一致・語境界・大小区別・同一原語は先勝ち）。
- `internal/engine`・`internal/store`: 取込段（`Ingest`）の最終ステップで、本文中の言及と叙述文の説明対象を検出して新規テーブルへ記録する。既存テーブルへは読みだけ。
- `docs/er.md`: 3 テーブルを実装済みへ反映し、正規化根拠 5（言及の相手の排他 2 列）・6（e3 の専用テーブル化）を追加した。
- `docs/known-issues.md`: 1 番から e3/e4/e5 を除き、残り（e7・e8/e14・e13、漏れ語の第 2 層）へ整理した。2 番へ「照合対象は言及テーブルが持つ」を追記した。
- `docs/concept-model.md`: 言及 note の未実装注記を実装参照へ置き換えた。

### 判断

- 検出方式（`known-issues.md` が実装判断に委ねた点）: 機械置換と同じ照合規則にした。言及レコードと注入語が一致することを、注入語の保持検証（`known-issues.md` 2 番）の前提にするため。
- 言及の相手: 概念上は固有名箱（`proper_noun`）だが、base ゲーム由来の名前は横断辞書 `master_term` にしか載らないため、排他 2 列（`proper_noun_id` / `master_term_id`）で両供給源を指せる形にした。
- e3 の物理形: 計画時の想定は `narration.described_proper_noun_id` 列の追加だったが、C# 抽出器が全 migration SQL を毎回冪等 ensure する契約で `ALTER TABLE` が使えず、`line_condition`（migration 0007）と同型の専用テーブル `narration_described` へ切り替えた。
- 検出時機: 取込段の最終ステップにした。言及は原語の出現で決まる関連で、訳の確定状態に依存しないため（固有名フェーズを待たない）。
- 非劣化の確認: 合成 golden（`internal/harness`）と実データ golden（inigo.esp、分岐元 `cf5d4038` で捕獲し本変更と比較）が一致。既存の dest・機械置換内訳・実プロンプトは変わらない。

## 2026-07-03 known-issues.md 新設と散在した既知課題の集約

### 変更

- `docs/known-issues.md`: 新設。docs 全体に散らばっていた「未解決の課題・未確定の設計判断・未実装の後続 task」を 1 か所へ集約した。現在開いている課題（言及テーブル未実装、辞書に無い漏れ語の拾い上げ、固有名一貫性の事後検証、Dialogue tree の context 長さ）だけを載せる。
- `docs/index.md`: Directory Contract と Choose The Right Record に `known-issues.md` を登録した。
- `docs/core-beliefs.md`: §2 記録システムへ「未解決の課題・未確定の判断・未実装の後続 task は `known-issues.md` に記録する」を追加した。
- `docs/system_requirements.md`: §2 の未確定のうち訳語供給方式（既訳 `master_term` と AI 訳 `proper_noun` の併用）と辞書適用方式（機械置換注入）を確定済みへ書き換え、漏れ語拾い上げと一貫性検証を `known-issues.md` へ移した。§3 の属性選定・衝突優先順位の「未確定」を確定済みへ書き換えた。
- `docs/concept-model.md`: 「既知の弱点」節を `known-issues.md` への集約に置き換え、同一ファイル内の弱点 1/2/3 への相互参照を現状記述へ更新した。
- `docs/skyrim-structure-model.md`: 「既知の弱点」節を `known-issues.md` への集約に置き換え、「実装の制約」節を C# / Mutagen 抽出器（`tools/extractor/`）参照へ改め、解消済みの「populate しない / 設計と差がある もの」一覧を削除した。
- `docs/er.md`: §3「未実装」の後続 task（言及・会話の流れ・名乗る名）を `known-issues.md` へ移し、畳み込みの設計判断（定型句・配置・無訳片）だけ残した。「既知の論点」を「設計の範囲」へ改めた。
- `docs/mutagen-migration-plan.md`: 削除。移行完了で全項目が解消済みのため（結果は完了 plan `mutagen-extractor` が持つ）。

### 判断

- 集約の範囲: 人間指示で、設計文書に組み込まれた「既知の弱点・未実装」節も含め、既知課題の記述を物理的に `known-issues.md` へ移す方針を採った。元文書には現在状態と `known-issues.md` への参照を残す。`core-beliefs.md` §3 の「同一責務を複数文書で別定義」を避けるため、開いている課題の正本は `known-issues.md` に一本化した。
- 解決済みの扱い: 人間指示で、各項目の現状をコード・migration・完了 plan で確認し、未解決のものだけ `known-issues.md` に載せた。解決済みは元記述と「解消済み」注記を正本から全て削除し、経緯は本 entry にだけ残した。
- 解決の確認根拠:
    - 訳語供給・辞書適用（`system_requirements.md` §2）: `master_term`（migration 0003、既訳流用）と `proper_noun`（migration 0006、実行内 AI 訳）の併用を本文へ機械置換注入する形で確定（`er.md` の正規化根拠 4）。
    - 属性選定・合成優先度（`system_requirements.md` §3、`concept-model.md` 弱点 3）: 口調段階の決定を本文優先・声型 prior・保留の 3 段階で `internal/core/tone/classifier.go` の `fuseAttitude` に実装。合成順は `internal/core/personatone/personatone.go`。
    - 重複排除の境界（`concept-model.md` 弱点 4）: `record_type_master`（migration 0006）が REC:FIELD → box を割当し、排他・網羅を engine の seed 整合テストで担保。
    - 固有名の同定単位（`concept-model.md` 弱点 1）: 種別（`proper_noun.category`）で粗く分ける決定を実装。同種別内の誤統合は受容済みの制約として `concept-model.md` の note に残す。
    - 純汎用台詞・PC 発話の口調（`skyrim-structure-model.md` 弱点 2）: `tone_default`・`line_condition`（migration 0007、generic-voice-tone-fallback）で口調付与済み。
    - 声型代表 Speaker と VoiceType の二重性（`skyrim-structure-model.md` 弱点 3）: 役割分担で解消（VoiceType は口調 prior、声型代表 Speaker は名前解決）。
    - 抽出器と構造モデルの差分（`skyrim-structure-model.md`「populate しない」）: xEdit Pascal 版の全項目（会話の 2 階層化・r2a〜d の分離・TACT/TPLT・VoiceType 独立化・master plugin 取込・Skyrim 追加 record）は Mutagen 移行（完了 plan `mutagen-extractor`）で解消。Pascal 版ファイル `extractData.v2.pas` は削除済み。
- 残す未解決: 言及テーブル（`narration_mention` e4・`line_mention` e5・`line_sequence` e7・`speaker_name` e8・`faction_name` e14）と関連 FK（e3・e13）は migration に無く未実装で、固有名一貫性は当面 `master_term` の機械置換で担保する。翻訳 runtime は台詞を 1 件ずつ翻訳し tree context を与えず、注入訳語の保持を照合する検証も無い。
- `mutagen-migration-plan.md` の扱い: 当初は移行完了の状態バナーで履歴として残す案にしたが、人間指示「解消ずみは全て消す」で削除へ切り替えた。移行の結果と検証実績は完了 plan `mutagen-extractor` が持つ。

## 2026-06-27 dev 起動ごとに中心 DB を空にする

### 変更

- `scripts/dev/run-wails.sh`: dev 起動時に中心 DB ファイル（`db/aitranslation.dev.sqlite3` と付随する `-wal`・`-shm`）を削除する処理を追加した。既存 listener の停止後、`wails dev` 起動前に消す。起動時の `store.Open`（`db.Apply` の migration）が空スキーマを作り直し、seed（`prompt_template`・`directive`）は復元される。

### 判断

- 動機: 中心 DB はファイルが永続し、取込が `INSERT OR IGNORE`（UNIQUE `plugin, form_id, rec, field, ordinal`）のため、再抽出しても既存行の `dest`・`status` を保持する。結果として再実行しても既訳が残り、新規翻訳が走らないため、開発時の動作確認の妨げになっていた。`project-db-wipe-on-launch-intent`（起動ごとに空にしたい意向、抽出・翻訳を残さない）に沿って、起動ごとに空から始める形へ確定した。
- 置き場所の選定: 本番ビルドはこのランチャを通らないため、dev 専用スクリプトでのファイル削除に限定し、本番の永続には影響させない。Go 側（bootstrap）に dev 分岐を足す案は採らず、起動入口の 1 箇所に閉じた。
- 全消去の粒度: 翻訳成果だけのリセットでなく DB ファイル全体の削除を選んだ（意向「抽出・翻訳を残さない」に一致）。`master_term`（固有名辞書 25071 行）も毎起動で消え、次回 C# 抽出で再構築される。`prompt_template` を UI で手編集していた場合は migration の既定値へ戻る。dev 専用かつ全消去の明示要望のため許容する。
- 実機確認: スクリプト経由で再起動し、`narration`・`line`・`proper_noun`・`persona_character`・`line_analysis`・`master_term`・`extracted_field` が全て 0、`prompt_template`=1・`directive`=7 が復元、実画面の結果一覧が「まだ結果はありません」の空状態であることを確認した。

## 2026-06-27 record-type-translation-expansion の AI 実走確認と正本化判断（finalization）

### 変更

- `docs/exec-plans/active/2026-06-23-record-type-translation-expansion/plan.md`: finalization-module 記録（AI 実走の実画面確認・検証・正本化判断・merge）を追加。
- `docs/architecture.md`: §3（engine 責務へ取込段・固有名フェーズを反映）・§8（現状記述を `extracted_field` 経由の振り分けへ修正、本 task の移行 entry を追加）を更新。
- `docs/er.md`: 実テーブル基準へ作り直し。採用原則を「厳密 1 対 1」から「概念由来＋実装正規化」へ緩め、テーブルを概念箱由来／実装・運用／未実装の 3 区分で記述。スコープを抽出入力から中心 DB へ拡張。`docs/index.md` の er.md 説明も合わせた。

### 判断

- パイプライン段階追加の正本化: 翻訳手続きへ「取込段（C# 素朴吸い出しの `extracted_field` を `record_type_master` で箱別へ振り分け）」と「固有名フェーズ（固有名を本文より先に確定し辞書化）」を加え、箱判定の責務を C# 抽出器から Go engine へ移した。当初は層・依存方向・Wails 境界のいずれも不変（新規 package/box なし、engine の新規 consumer interface は §4 既存パターン、Bind 追加は §5 機構を変えない）を理由に、`feedback-architecture-reflection-structural-only` に従い `docs/architecture.md` 反映を見送った。merge 後に人間指示「他の乖離を直してくれ」を受け、構造不変でも §8 の現在状態の記述が実装と乖離している点を反映へ切り替えた。§3・§8 を修正し、口調指示の供給が `prompt_template` の口調テンプレートから口調 `directive` へ移った点も §8 entry に明記した。
- `docs/er.md` の扱い: 当初「概念箱でない config・実現方式テーブルなので追記不要」と判断したが、人間指示「er.md は現テーブルと同期しなければならない」で撤回。方針は「er.md は概念モデルに基づき実装のため正規化したテーブル設計。厳密 1 対 1 は不要だが、概念モデルと無関係な恣意的 DB はだめ」。当時の er.md は概念モデル全体の論理 ER（未実装テーブルを含む）で実 DB と両方向に 8 テーブルずつ乖離していたため、実テーブル基準へ作り直し、各テーブルの概念由来または実現方式を追える 3 区分で固定した。
- AI 実走で MECE モデルの動作を確認: 固有名フェーズ先行→ペルソナ生成→本文フェーズ（叙述文・台詞）の段階順が hy-mt2-7b で成立。固有名注入が叙述文・台詞を通して一貫し、台詞の口調が話者差で出た。7B 翻訳特化モデルでも注入カタカナを保持する点を `project-injected-token-fidelity` へ追記（「弱い小型モデルは崩す」の単純一般化は不成立）。

## 2026-06-25 record-type-translation-expansion の Storybook レビューと MECE モデルへの収束

### 変更

- `frontend/src/ui/screens/template-editor/`: プロンプトテンプレート画面をサブタブ [ベース][レコード別] へ作り直し、レコード別タブを「種別ごとの指示文（directive）」の一律一覧にした（`TemplateEditorScreen.svelte`・`TemplateBasePane.svelte`・`DirectiveEditor.svelte`・`directive-view.ts`・`directive-presentation.ts`・`template-editor.fixtures.ts` 新設/改訂、`template-editor-view.ts` から `TemplatePlaceholder` 撤去、`TemplateEditorContainer.svelte` から不要 placeholders 受け渡し撤去）。
- `frontend/src/ui/screens/translation-run/`: 結果行へ元レコード種別バッジ（箱 ・ REC:FIELD）を追加（`TranslationResultRow.svelte`・`translation-run-view.ts`）。
- `frontend/src/ui/screens/record-type-master/`: 独立画面のレコード種別マスター一式を削除（directive へ畳んだため）。
- `frontend/package.json`: knip ignore を整理。
- `docs/exec-plans/active/2026-06-23-record-type-translation-expansion/implementation-scope.md`: 実装範囲を MECE モデルへ更新。`plan.md`: status・HITL Status・合意済み frontend 保護を記入。

### 判断

- Storybook 人間レビューで設計が段階的に収束した。経緯: 一覧に論理名を追加 → 種別ごとの文体 select は「割り当ては変えない」ため撤去（編集は指示文だけ）→ 翻訳対象列と翻訳しない種別をマスターから外す（読み込まない種別はマスターに不要）→ 独立画面でなくプロンプトテンプレートのタブへ統合 → 「MECE 感がない」を受けプロンプトの作られ方から再設計。
- 確定モデル: プロンプト = Base 指示 ＋ その REC:FIELD に割り当てた指示文（directive、変数を実行時に埋めたもの）。口調・文体・固有名・定型句を「指示文」という 1 つの形に揃え、口調は `{traits}` 変数を持つ指示文の 1 つにした。固有名・定型句も指示文を編集できる。各 REC:FIELD は指示文を 1 つだけ持ち（排他）、全 translatable REC:FIELD が指示文へ割り当たる（網羅）。
- データ面: `style_template`（文体専用）を `directive`（全種別の指示文）へ一般化。`prompt_template` の口調（persona）は directive の口調行へ畳む。`record_type_master` は REC:FIELD → directive の割り当てを持ち、翻訳対象だけを載せる（翻訳対象フラグ・無訳片の行を持たない）。
- storybook-module の表示範囲を越える最小の触り: `TemplateEditorContainer` から不要になった placeholders 受け渡しを撤去した。タブ状態・directive データ・指示文保存の本配線は implementation-module。

## 2026-06-25 record-type-translation-expansion の設計確定（人間設計レビュー承認）

### 変更

- `docs/exec-plans/active/2026-06-23-record-type-translation-expansion/implementation-scope.md`（新規）: design-module 出口の実装範囲とテスト設計を固定。
- `docs/exec-plans/active/2026-06-23-record-type-translation-expansion/plan.md`: `status` を design-module 承認済み・storybook-module 待ちへ更新。後続モジュールが埋める枠（Artifact Index・Routing Notes・HITL Status）を記入。
- `docs/exec-plans/active/2026-06-23-record-type-translation-expansion/design-review.md`: 人間設計レビュー承認後に削除（一時材料）。

### 判断

- 人間設計レビューを承認で通過した（2026-06-25、人間の「先へ進めて」）。確定済み方針（方針 A=master_term 権威訳専用、レコード種別マスターで box と style を持つ、純粋判定は供給源選別とプロンプト合成の 2 つ、C# 抽出器は素朴吸い出し、style 画面編集）に加え、残る 3 確認点を確定した。
- 確認点 1（抽出生テーブル）: `extracted_field` を新設して採用する。C# 抽出器を素朴吸い出しに徹させる方針の必然の受け皿で、箱判定を Go の取込段へ集約でき、C#/Go で箱知識が重複しない。抽出器が中心 DB のマスターを読む依存を増やさない。
- 確認点 2（定型句）: `RNAM`・`MESG:ITXT`・`WOOP:TNAM` は重複排除せず `narration` 系として位置キーで訳す簡略案を採用する。`set_phrase` 重複排除テーブルは構造が増える。定型句訳の不一致やコストが実測で問題になった時に後続 task で重複排除を足す（起動条件付き）。
- 確認点 3（style カラム方式）: 文体キー方式を採用する。`record_type_master.style` に文体キーを持ち、`style_template`（キー → 指示文）で指示文を共有する。文体は説明体・書物体・日記体・世界観断片の少数集合で複数 REC:FIELD が共有するため、指示文は文体に属する。直書き方式は同一指示文が rec 行へ重複し、文体 1 つの修正が複数行編集になる。概念モデルの文体を 1 対 1 で写すため正規化を採る。

## 2026-06-20 T4（prompt-persona-customization）の backend 本番接続と旧 E2E 残骸の削除

### 変更

- プロンプト構築を `internal/provider/openai_compatible.go` から `internal/engine/prompt.go` の純粋関数（`ComposePrompt`・`RenderPrompt`）へ移設。`provider.Translate` を完成 `Prompt`（System/User）受け取りへ変更し、base 指示の文面定数を撤去した。
- `internal/api/app.go`: `ResultView` に機械置換内訳（`Terms`）と実プロンプト（`Prompt`）を足し、`ListResultsPage` が各行へ辞書とテンプレートを当て直して取得時に再構成する。`GetPromptTemplate` / `SavePromptTemplate` を足した。
- プロンプトテンプレートを永続化した。`db/migrations/0004_prompt_template.sql`（単一行テーブル・既定値 seed）、`internal/model/prompt_template.go`、`internal/store/prompt_template.go`。`internal/engine` は翻訳実行と結果取得の両方で DB テンプレートを読む。
- 口調指示を編集可能テンプレートの `{traits}` 差し込み駆動へ見直した（`internal/engine/persona.go`）。性質文の中身（属性 → 性質文）は `persona_rule.go` のハードコードのまま使う。
- frontend: テンプレート編集 container（`TemplateEditorContainer.svelte`）、画面間ナビの本番ルーティング（`App.svelte` で AppShell＋AppNav）、gateway（`template-gateway.ts`、`translation-gateway.ts` の terms/prompt 写像）を配線した。
- `docs/architecture.md` §3・§8 にプロンプト構築の所在移設、`prompt_template` 永続化、Wails 新メソッド、`ResultView` の terms/prompt 供給を反映した。
- 旧 E2E system-test 一式（`scripts/test/seed-system-test-db/`、`scripts/test/run-system-test.sh`、`playwright.config.ts`、`package.json` の `test:system` / `test:system:install`）を削除した。

### 判断

- プロンプト構築は LLM 出力品質に直結し、翻訳実行と結果参照で同じ文面を使う必要がある。所在を `provider` から `engine` の 1 関数へ集約し、実行時と取得時の両方が同じ関数を呼ぶことで実プロンプトの一致を保証した。
- テンプレートは抽出・翻訳データと寿命が違う（編集結果を残したい）。中心 DB 内の専用テーブルに分離し、起動ごとの中心 DB 消去の対象から外せるようにした。永続化方針自体は未決のまま、置き場所だけ先に分けた。
- 口調指示の精緻化は「テンプレート構造の見直し」を完了線にした。性質列を差し込む `{traits}` 雛形を編集可能にし、属性 → 性質文の対応の編集は新 plan（`2026-06-20-character-persona-from-dialogue`）へ残した。
- 旧 E2E は greenfield で spec 本体が既に削除済みで、起動スクリプトと設定だけが残り `go build ./...` を壊していた（削除済み `internal/repository` を import）。新アーキの統合確認は実画面（chrome-devtools 手動）で行うため、連鎖一式を削除した。`@playwright/test` 依存はロックファイルを荒らさないため残置。`.claude/settings.json` の stale permission（`test:system` 系）は auto-mode により未編集で残る。

## 2026-06-20 T4（prompt-persona-customization）の scope 縮小とルール編集の切り出し

### 変更

- `docs/exec-plans/active/2026-06-14-prompt-persona-customization/plan.md`: scope を「テンプレート編集・実プロンプト参照・機械置換内訳・口調指示テンプレート精緻化・画面間ナビゲーション」へ縮小。属性 → 性質文のルール編集を「含まない」へ移し、完了定義を 6 項目から 4 項目（テンプレート編集・実プロンプト参照・機械置換内訳・口調精緻化）へ整理。
- `docs/exec-plans/active/2026-06-14-prompt-persona-customization/implementation-scope.md`: 縦切りを 4 段から 3 段へ。旧 Slice 3「ルール編集・永続化・反映」を削除し、旧 Slice 4「テンプレート編集・口調精緻化」を Slice 3 へ繰り上げ、画面間ナビゲーションを同 Slice に含めた。
- `docs/exec-plans/active/2026-06-20-character-persona-from-dialogue/plan.md`（新規）: T4 から切り出したルール編集の受け皿。種族・汎用声型への固定ペルソナ割り当てと、キャラ専用声型ペルソナの台詞群からの生成を、仕様から設計し直す骨子。

### 判断

- T4 の storybook レビューで、ルールの持ち方を見直す指示が出た。属性キーは抽出データ由来とし、ユーザーは性質文を割り当てるだけにする。声型は汎用（パターングループ）とキャラ専用（1:1）の 2 層に分ける。
- キャラ専用声型の個別ペルソナは、ユーザーが手で書くのではなく、そのキャラの台詞群から生成したい。生成は仕組みが大きく、T4 の「テンプレート編集」より広い。T4 へ無理に詰めず、仕様から起こす別 task（`2026-06-20-character-persona-from-dialogue`）へ切り出した。
- T4 は「カスタマイズ（テンプレート編集）を先に終わらせる」方針に絞る。種族・声型 → 性質文の対応は現状ハードコード（`persona_rule.go`）のまま使い、その編集は新 plan で扱う。

### 残課題

- 新 plan `2026-06-20-character-persona-from-dialogue` は骨子のみ。preparation-module で仕様を起こす。

## 2026-06-20 設計説明 skill を presentation に統合（diagramming 廃止）

### 変更

- `.claude/skills/presentation/SKILL.md`（新規）: 人間が読んで判断する、わかりやすい説明 md を作る下位 skill。説明文と図（Mermaid）を 1 つの md に統合する。図作法（差分凡例 赤=削除・緑=追加・黄=変更なし、Before/After 並置、図の分割、設計差分図の範囲限定）を含む。
- `.claude/skills/presentation/references/templates/review-diff-template.md`: `diagramming` 配下から git mv で移動。
- `.claude/skills/diagramming/`（削除）: SKILL.md と references を削除。
- `.claude/skills/design-module/SKILL.md`: 下位 skill 参照を `diagramming` から `presentation` へ付け替え。人間設計レビュー材料の固定 2 セクション（概要・図）を廃止し、`presentation` の厚め構成（背景・課題、採用方針、代替案、図、影響範囲、テスト要点、確認点）へ寄せた。
- `.claude/skills/presentation/SKILL.md`（追補）: わかりやすさの規約を、確立した資料設計の原則で埋めた。構成の順序（ピラミッド原則の結論ファースト、グループ化）、文（一文一義、専門語に意味を添える）、視覚（近接・整列・反復・対比、関連情報の比率を上げる、視線の流れ）、図（色だけに依存せずラベルまたは線種を併用）を追加。design-module 固有の記述（最初の呼び出し元、Wails 境界、実装範囲・テスト設計 artifact 名）を外し、呼び出し元が固有セクションを指定する一般形にした。
- `.claude/skills/presentation/references/templates/review-diff-template.md`（削除）: 差分図の穴埋め雛形を廃止。
- `.claude/skills/presentation/SKILL.md`（追補）: 雛形参照を外し、構成と図種別を「題材に最適化する判断」へ直した。厚め構成を固定セットから論点候補へ変更。図種別はコンポーネント図前提をやめ、主張を最も伝える図を選ぶ形にした。差分凡例（色＋ラベル）の原則は残した。
- `.claude/skills/presentation/SKILL.md`（追補）: 図のサイズ基準を追加。1 枚を大きくしすぎず、16:9 のプレビューで文字が読める大きさに収める。要素は 7 個前後を目安にし、超えるか縦横に伸びて縮むなら主張の単位で分割する。完了条件に同基準の検査を追加。

### 判断

- `diagramming` の現役呼び出し元は `design-module` 1 つだけだったため、図作法を `presentation` へ統合し、独立 skill を廃止した。
- `presentation` は「人間が読んでわかりやすい説明 md（図を含む）」を担う。図は説明 md を分かりやすくする一手段として内包する。
- 成果物の寿命は `presentation` が決めず、呼び出し元に従わせる。`design-module` の人間設計レビューは一時（レビュー後削除）を維持する。
- 人間設計レビューの図は、Mermaid の色分け強調と Before/After 並置でスライド級の視認性へ寄せる。1 枚ずつめくる段階表示の演出は、承認・差し戻し判断に不要として持たない。
- わかりやすさの規約は、確立した資料設計の原則（ピラミッド原則、関連情報の比率、デザイン 4 原則、色だけで伝えない WCAG 1.4.1）を出典として埋めた。原則名は固定名として残し、意味を日本語で添えた。
- `presentation` は design-module 専用にせず、説明 md・図が要る module 全般から呼べる一般形にした。design-module 固有のセクションは design-module 側に残し、`presentation` へは呼び出し元の指定として渡す。
- 穴埋め雛形を持たせない方針にした。題材ごとに最適な図種別と論点が変わるため、雛形は 1 つの形へ寄せて関連情報の比率を下げる。残すのは判断基準（差分凡例の色＋ラベル、図種別の選び方、わかりやすさの原則）とし、型は持たない。

### 残課題

- なし。完了済み exec-plans の `diagramming` 言及は過去記録として残す（正本ではないため変更しない）。

## 2026-06-18 マスター辞書 T3 増分: 人名の部分形の派生（名のみ・短名を辞書化）

### 変更

- `internal/engine/termderive.go`（新規）: 人名の部分形を 3 種で派生する純関数 `DeriveTerms` を追加。shrt（`NPC_:SHRT` の短縮別名通過）・byname（` the ` を含む名の前部と Dest 末尾カタカナ連）・two（base ゲームの空白 2 語姓名を中黒 2 語で整列）。安全フィルタ（用法比 lc/uc・称号・種族/ハウス語・縮約素体・純カタカナ・最小長 4・base 衝突 skip・two の base ゲーム限定）を同関数に持つ。副作用なし。
- `internal/engine/termusage.go`（新規）: 会話文の英語原文から各英単語の用法分布（lc=小文字始まり一般語用法 / uc=文頭以外の大文字始まり固有名用法）を作る純関数 `BuildUsage` を追加。
- `internal/engine/termderive_test.go`・`termusage_test.go`（新規）: 純粋ルールの全分岐単体テスト。新規ルール関数のカバレッジ 100%。
- `scripts/dict/derive-master-terms/main.go`（新規）: ビルド時コマンド。xTranslator 英日 XML を解析し、用法分布を作り、純粋ルールを呼び、base 衝突を除いて `master_term` へ `category="derive:<種別>"` で追記する。`db.Apply` で schema を ensure し、`INSERT OR IGNORE` で二重追記を防ぐ。
- `scripts/dict/derive-master-terms/main_test.go`（新規）: XML 解析の単体テストと、temp DB へ XML→派生→追記・base 衝突 skip の結合テスト。
- 実行時の置換器 `internal/engine/dictionary.go`・`engine.go`、テーブル `master_term` は無改造。`loadDictionary` が category を問わず全件を読むため派生行を自動で取り込む。

### 判断

- 派生規則は副作用の無い単一の純粋ルール（`DeriveTerms`）へ分離し、ユニットテストカバレッジ 100% を基準にした（人間が固定した基準）。XML 解析・DB 書込の I/O 配線はルールの外（ビルド時コマンド）へ出し、結合テストで見る。
- 置き場所は判定が属する言語（置換器が Go なので Go）に合わせ、ビルド時生成で `master_term` へ焼く（人間承認の構成）。永続した派生行は category 印で目視・差し戻しできる。
- base 辞書は現行の C# extractor のまま。派生は base 書き込み後に走る Go コマンドが追記する 2 段構成。同じ既訳をより単純に得られ、純粋ルールを Go に置けるため。
- 由来種別は既存 `category` 列に `derive:<種別>` として持たせ、スキーマは変えない。実行時の照合は source だけを見るため影響しない。
- two（姓名分割）は base ゲーム XML 限定にした。patch/mod の Source/Dest 対応ずれ（USSEP で観測）による誤訳を避けるため。
- 安全フィルタは観測した失敗の手書き除外でなく、用法分布と語形による構造判定にした。一般語（`Master`・`Blood`・`Mine`）・種族語（`Imperial`・`Nord`）・称号（`Lord`・`Captain`）・縮約（`Aren`）を捨て、固有名（`Grelod`・`Mercer`）を残す。`Imperial` は用法比だけでは残る側だが種族/ハウス語集合で捨てる（地の文の誤置換回避）。
- two の category 形容語フィルタ（creature ラベル混入抑制）は本実装に入れない。検証済み数値（破壊型過剰置換 0・被覆 99.9%・held-out 汎用性）はフィルタ無しで得たもので、追加は未検証のため。起動条件は「ハーネスで過剰置換が減り被覆が落ちないと実測したとき」。

### 残課題

- 機械置換は実 DB（base 24,554＋派生 517）で観測済み: `Grelod`→`グレロッド`（名のみ、原問題解消）、`Mercer Frey`→`メルセル・フレイ`（最長一致）、`Mercer`→`メルセル`（姓のみ）、一般語 `master` は無置換。AI 無しで `engine.NewDictionary`+`Apply` を実 DB に対し実行して確認。
- 実 app の end-to-end も観測済み（実 LLM、2026-06-19）: 「Innocence Lost - Quest Expansion.esp」の台詞 121 件を gemma-4-12b で翻訳し、`Grelod` 28/28→`グレロッド`、`Riften` 4/4→`リフテン`、崩れ 0。観測前は前回（派生なし）の訳で裸 `Grelod` が `グレロド`/`グレロッド` に揺れていたのが、派生辞書で全て `グレロッド` に揃った。plugin はパス直接入力欄から手入力（人間補助不要。先の「ネイティブダイアログのため人間依頼」は誤り）。
- 所見（モデル依存・新規発見）: 機械置換は AI 前段で確定訳語を注入するが、end-to-end の一貫性は AI が注入カタカナを保持するかに依存する。shisa-v2.1-qwen3-8b は注入済み `リフテン` を `リヴェン` 等へ書き換え（`Riften` 0/4・`Grelod` 24/28 のみ保持）、gemma-4-12b は完全保持。派生のバグではなく（base 名 `Riften` も崩れる）、弱い小型モデルが注入トークンを保持しない問題。翻訳モデルは注入トークンを保持できる能力のものを選ぶ。
- 派生はビルド時 1 段（`go run ./scripts/dict/derive-master-terms --sqlite db/aitranslation.dev.sqlite3`）。中心 DB は wipe されず writer が `INSERT OR IGNORE` のため、app の再抽出で派生行は消えず runtime 無改造で base+派生が効く。
- `docs/architecture.md` 反映の要否は finalization-module で判断する（本増分は engine 内の純粋ルール追加とビルド時コマンドで、Wails 境界・実行時の層構造は変えていない）。
- `go build ./...` は既存の壊れ（`scripts/test/seed-system-test-db/main.go` が存在しない `internal/repository` を import）で失敗する。本増分の範囲外。

## 2026-06-15 結果一覧の N+1 廃止＋keyset cursor ページング（T2 の効率課題の恒久対応）

### 変更

- store: 台詞ごとに話者を引く `LoadLineSpeaker`（1 台詞 2 クエリ）を、台詞 id 群を一括で引く `LoadLineSpeakers(lineIDs)→map`（IN 句・host parameter 上限回避の chunk）へ置換。keyset 範囲取得 `NarrationsAfter`／`LinesAfter`（`WHERE id > ? ORDER BY id LIMIT ?`）と `CountNarrations`／`CountLines`、共通 helper `query.go` を追加。未使用の `ListNarrations`／`ListLines` を削除。
- engine: `LineStore` を一括取得へ差し替え、`LineDirective`／`linePersona`（per-line）を `LinePersonas(lineIDs)→map[int64]Persona` へ置換。`Run` はループ前に話者を 1 度だけ一括取得し、ループ内の個別 DB 問い合わせを廃止。
- api: `ResultPage{Total,Results,NextCursor,HasMore}` と `ListResultsPage(cursor,limit)`、連結列のページ範囲を決める `pageRows`、cursor 解析（`""`／`n:<id>`／`l:<id>`）を追加。`buildResults`（全件）を廃し、ページ内台詞ぶんだけ口調を一括生成する形へ。`RunExtractAndTranslate` は全件 `Results` を返すのをやめ件数要約だけ返す。
- frontend: gateway に `listResultsPage(cursor,limit)`、container に keyset state（cursor 履歴・ページ index・total・nextCursor・hasMore・ページサイズ 50）と前へ/次へ。ページャ表示 `ResultsPager`（順次送り・端で無効化・現在ページ番号）を追加し、件数バッジを総件数へ（Storybook 人間レビュー承認）。`TranslationResultRow` の表示形（コンパクト行・口調チップ・展開）は不変。

### 判断

- ページング方式は keyset（cursor）を採用し LIMIT/OFFSET を不採用にした（人間設計レビューで確定、ページサイズ 50）。結果一覧は閲覧中に行が増減しない静的集合だが、数万件の深い位置を `WHERE id > ?` の index 走査で取れる keyset を選んだ。仮想スクロールは全件を frontend へ載せ payload を増やすため不採用。
- 話者一括取得は map 返し（`LoadLineSpeakers(lineIDs)→map`）にした。所属勢力が話者に対し 1 対多のため `ListLines` への JOIN は行が増殖し group_concat 等が要る。map なら所属勢力を `[]string` のまま保て、T2 の責務分担（store は識別子・事実、engine が口調へ解釈）を維持できる。同じ map を表示と翻訳で使い回す。
- N+1 は per-line メソッドを削除して一括メソッドへ置換し、「N+1 を表現できない」形にした（特殊対応の追加でなく機構の置き換え）。
- `RunExtractAndTranslate` は全件結果のインライン返却をやめた。数万件を実行応答に載せず、実行後・起動時・ページ送りを `ListResultsPage` の 1 経路に統一した。
- `docs/architecture.md` への反映は不要と判断（§5 Wails 境界・§7 ディレクトリ正本・§8 現在の状態の構造を変えていない。keyset と一括 persona は層内の内部精緻化、`ListResults`→`ListResultsPage` は既存 Bind 境界内）。

### 残課題

- 固有名解決・マスター辞書は T3、口調ルールの精緻化・編集 UI は T4（いずれも対象外）。
- 翻訳実行中（`engine.Run`）の話者一括取得は実装済みだが、翻訳は AI latency 支配のため最適化の主目的は表示経路（`buildResults`）にある。
- 本 task の実画面検証はローカル LLM stub で pipeline 全体を通した（121 台詞・3 ページ）。AI 翻訳の本番実行は利用者の OpenAI 互換 provider で行う。

## 2026-06-14 T2 ペルソナ口調 pipeline（台詞抽出→話者解決→口調注入翻訳→進捗・口調差を画面で観測）

### 変更

- extractor（C#）に台詞（INFO:NAM1）と話者属性（speaker / race / faction / voice_type）の SQLite 書込を追加。INFO の話者 FormKey を LinkCache で NPC へ解決し、種族・声型・所属勢力の EditorID を書く（`LineSpeakerSqliteWriter`）。
- engine（Go）に台詞翻訳とペルソナ口調生成を追加。話者の声型/種族/勢力 EditorID から口調 traits を引く最小ルール（`persona_rule.go`）と、口調指示文を組む `buildPersonaDirective`／チップ用 `buildPersonaLabel`。provider の `Translate` に directive 引数を足し system prompt の base 指示後段へ注入。`Run` に本文翻訳の進捗 callback。
- api に本文翻訳の進捗 runtime events（extract / translate）と、結果へ口調指示文・口調要約を載せる `ResultView`／`ListResults` を追加。
- frontend に本文翻訳の進捗バー（`TranslationProgress`）と結果行のコンパクト化（口調チップ＋展開）、進捗 event 購読を追加。db migration `0002`（line / speaker / race / faction / voice_type / line_speaker / speaker_faction）。
- `architecture.md` §8「現在の状態」を現状へ更新（extractor が台詞・話者も書く、engine が台詞翻訳・ペルソナ・進捗を持つ）。

### 判断

- task の完了定義を「縦切り（観測可能な成果）」へ置き直した。当初の「翻訳プロンプトへ差込点を 1 本通す（単体テストで確認）」は seam 層を完了と呼ぶ弱い条件で、実 mod で観測できる成果が無かった。実 mod `Innocence Lost - Quest Expansion.esp` を実行画面から流し、台詞抽出・口調差・進捗・翻訳を実画面で観測することを完了条件にした。
- 固有名（辞書）解決は本 task から外し、マスター辞書 task（T3）へ移送。T3 の依存「T2 の辞書解決の差込点」は無効化（T3 が自身で差込点を作る）。
- 責務分担: 事実の抽出は extractor（C#）、口調などの解釈は engine（Go）。extractor は識別子・事実（EditorID）だけ書き、口調 traits は engine の最小ルールで与える。ルールの永続化と編集 UI は T4（対象外）。
- 結果行 UI は数万件・ページングに耐えるコンパクト 1 行＋口調チップ（展開で全文）にした（Storybook 人間レビュー承認）。口調差は一覧のままチップで観測。
- architecture 構造（§1〜7）は変えていない。engine 責務（辞書解決・ペルソナ生成）と runtime events は §3・§5 に既記載で、追加スキーマは `er.md` に既定義。§8 の現状記述だけ人間承認のうえ更新。

### 残課題

- 固有名解決・マスター辞書（proper_noun 抽出、line_mention e5、name 関連 e8/e13/e14）は T3。ルール・プロンプト編集 UI は T4。
- ペルソナ口調ルールは engine 内の固定最小 1 系統。声型 EditorID の網羅と気質テキストの精緻化、結果一覧の仮想スクロール／ページングは後続で扱う。
- AI 翻訳の本番実行は利用者の OpenAI 互換 provider で行う。本 task の検証はローカル stub で pipeline 全体を通した。

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
