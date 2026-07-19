# Investigation: slow-plugin-extraction

`investigation.md` は不具合の再現と原因究明だけを持つ。修正方針（どう直すか）は `design.md` が持つ。

## 観測済み問題

`structured-prompt-io` task の実画面確認中に人間が観測した記録。

- Outfit Recognition Framework.esp（633KB、翻訳対象 1805 行）: 抽出フェーズが約 6 分。app 本体プロセスが CPU 約 108%（STAT=R）で処理し続けた。ハングではなく実処理。
- Innocence Lost - Quest Expansion.esp（124KB、197 件表示）: 抽出フェーズが約 60 秒。
- 観測時の理解では「抽出は Go backend（app 本体）が実行、C# 別プロセスでない」とされた。後述のとおり実際は C# 子プロセスで、この帰属は不正確だった。

## 画面再現確認

実 app（Wails）経由の 6 分再現は本調査では行っていない。理由は 2 点。

- 抽出対象 plugin の選択がネイティブファイルダイアログで、実画面操作の代替が難しい（人間依頼が要る）。
- 後述の段別実測で、available データでは翻訳前区間が数秒に収まり、6 分が再現しないことが判明したため、実 app 再現より先に律速の有無そのものを実測で確定した。

代わりに、実 app が呼ぶのと同じ経路（`internal/api/app.go:RunExtractAndTranslate` → C# 抽出子 `dotnet run` → Go 後段 `DeriveMasterTerms`/`LoadReferenceTranslations`/`Ingest`）を、実データ（`dictionaries/Data`）に対して段別に実測した。

## 原因仮説

翻訳前区間（AI へ翻訳を投げる前の全処理）の律速候補を、検証順に立てた。

1. C# 抽出子の環境ロード: master 連鎖（Skyrim.esm 249MB 等）を毎回フル index 化する（`PluginEnvironment.Load` → `RecordDataIndex.Build`）。
2. C# 抽出子の override 判定: `OwnsRecord` が record 件数 × master 数で `Normalize`（zlib 展開・2 回 Sort）を総当たりする。
3. Go 後段の言及検出: `recordMentions` が語彙（master_term 24554 件 ∪ proper_noun）全件の巨大単一正規表現を、叙述文・台詞の各本文へ当てる（`internal/core/mention/mention.go`）。
4. 永続 dev DB の肥大: 全件 `List*` と INSERT 群が蓄積で膨らむ。
5. 本番 `dotnet run` の起動オーバーヘッド: 実行毎の build/restore 判定、bin/obj 欠落時の全ビルド、NuGet キャッシュ欠落時の全 restore。

## 観測ログ検証

C# 抽出子は既存の恒久ログ（`Program.cs` の `Stopwatch` 累積計測）で段別 ms を stdout に出す。Go 後段は使い捨ての計測テスト（`internal/bootstrap/zz_perfmeasure_test.go`、C# seed 済み DB のコピー上で 3 処理を計時）で測り、計測後に削除した。全て観測が起きた Mac 上、warm 状態。

段別実測（Outfit Recognition Framework.esp）:

| 段 | 所要 | 備考 |
| --- | --- | --- |
| C# load（master 連鎖 6 件） | 526 ms | Skyrim.esm 249MB も lazy overlay で軽い |
| C# extract（mod 本体走査） | 231 ms | |
| C# sqlite 書込群 | 約 400 ms | |
| C# master_term 24554 件 | 288 ms | |
| C# 壁時計 合計 | 約 1.8 s | |
| Go DeriveMasterTerms | 632 ms | |
| Go LoadReferenceTranslations | 929 ms | |
| Go Ingest（言及検出 24554 語 含む） | 82 ms | 巨大正規表現の言及検出も速い |
| 翻訳前区間 合計 | 約 3.5 s | 観測は約 6 分 |

補助実測:

- 本番同等 `dotnet run`（`--no-build` なし、build 判定込み）: 1.82 s。warm では起動オーバーヘッド無視できる。
- bin/obj 削除の cold `dotnet run`（NuGet はキャッシュ温存）: 4.5 s。ローカル再ビルドも速い。
- USSEP（unofficial skyrim special edition patch.esp、21.8MB、override 18838 件）: 2.9 s。available 最大 mod でも `OwnsRecord` は崩れない（master 連鎖が 6 件のため）。
- 抽出周辺コード（`tools/extractor`、`internal/engine` 言及/取込/参照、`internal/api/app.go`）は観測 commit `9c993d09`（2026-07-19 23:01）以降 変更なし。コードは観測時のまま。
- 実 dev DB（`db/aitranslation.dev.sqlite3`、41MB）の行数は控えめ（extracted_field 2006、master_term 25071、reference_translation 81358、line 1956 等）で、新規 DB 計測と同規模。DB 肥大は原因でない。

仮説の判定:

- 仮説 1〜4（環境ロード / OwnsRecord 総当たり / 言及検出 / DB 肥大）: available データで実測し、いずれも数百 ms〜数秒で否定。翻訳前区間は約 3.5 s に収まる。
- 仮説 5（`dotnet run` の初回コスト）: warm・ローカル cold（bin/obj 削除）とも数秒で否定。唯一残るのは NuGet 完全 restore（Mutagen 一式のネットワーク取得）で、初回 checkout 後・NuGet キャッシュ消去後に 1 度だけ起きる。非力な Mac ＋ネットワークでは数分に達しうる。抽出毎には繰り返さない。

## 確定原因

観測で確定したもの:

- 抽出処理そのもの（C# 抽出子 + Go 後段）は available データで高速（Outfit 約 3.5 s、USSEP 約 2.9 s）。「抽出自体が重い」は否定された。
- 観測の 6 分は available データでは再現しない。最も整合する説明は、初回 `dotnet run` が Mutagen 一式を NuGet restore ＋ビルドしていた `dotnet` プロセスを app 本体と見誤ったケース。1 度きりの事象で、確定はできず仮説にとどまる。

観測で確定できないが構造上実在する非効率（実データで顕在化しなかったため「潜在」）:

- `OwnsRecord` の `Normalize` 毎回再計算（zlib 展開・2 回 Sort）は、master 依存が数十件に及ぶ mod で record 数 × master 数に比例して効く。available データ（master 連鎖 6 件）では顕在化しない。
- 本番の `dotnet run` は実行毎の build 評価と、bin/obj 欠落時の全ビルド・NuGet 欠落時の全 restore の崖を踏む。崖の頂点が観測 6 分の最有力候補。
- 翻訳前区間は進捗イベントが `extract` の 1 回だけで、C# 子プロセス実行中と Go 後段 3 処理の間は無音。遅いとき画面が固まって見える。
