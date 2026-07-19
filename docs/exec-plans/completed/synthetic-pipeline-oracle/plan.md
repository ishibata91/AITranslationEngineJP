# Plan: synthetic-pipeline-oracle

- `task_id`: `synthetic-pipeline-oracle`
- `working_branch`: `claude/synthetic-pipeline-oracle`（base: `master`、分岐元 commit: `37fb6b03`）
- `target_branch`: `master`

## 依頼要約

単体テストでは届かない「エンティティ・クラスをまたいだ合成結果が仕様どおりか」を、C# 抽出機と Go メインシステム（翻訳）の双方で確かめる統合テストを作る。E2E は重く、AI の実画面確認が E2E 相当の確信を担うため、単体とは別に置く。判定基準はテストオラクルに持ち、C# と Go が同じ JSON を読んで照合する。

## オラクルの方針

- 中心は「システムが実現する仕様」。ドメイン語彙で書き、実装座標（table・列・status・訳値・prompt 文字列・file・method）を入れない。
- 粒度は処理段（stage）中心。UC は現状 1 本なので割る軸にしない。1 エントリ = stage × 属性。
- `given` は入力 esp のリッチさの一断面。`category`（正常/異常）は given から従属する。
- 判定基準を共有 JSON に持ち、C# 抽出と Go 翻訳の双方が同じ 1 ファイルを読む（どちらもパースできる形＝ JSON）。`id` が両テストの対応確認の join キー。

## オラクルのテンプレ（形）

1 エントリ 6 フィールド。

| key | 役割 | 値の例 |
|---|---|---|
| `id` | 対応テスト存在確認の join key。C# と Go のテストが同じ id を参照する | `persona-sex-from-female-flag` |
| `stage` | どちらのツールの継ぎ目か。担当と観測点が一意に決まる | `extraction`(C# 抽出) / `integration`(Go 翻訳の継ぎ目) |
| `attribute` | その継ぎ目で運ばれる属性 | 話者 / ペルソナ / 台詞感情 / 固有名 / 箱振り分け / 翻訳対象フィールド |
| `category` | 分類。`given` のリッチさから従属的に決まる | `正常` / `異常` |
| `given` | 入力のリッチさの一断面（どの属性が揃う/欠ける） | 話者 NPC に女性フラグ（Female flag）が立つ |
| `spec` | `given` のもとで観測できる、実現する仕様 | 性別を女性フラグから確定する |

JSON の形（stage 代表 2 例）:

```json
{
  "specs": [
    {
      "id": "persona-sex-from-female-flag",
      "stage": "extraction",
      "attribute": "ペルソナ",
      "category": "正常",
      "given": "話者 NPC に女性フラグ（Female flag）が立つ",
      "spec": "性別を女性フラグから確定する（喋るオブジェクト TACT は性別を持たない）"
    },
    {
      "id": "line-emotion-injected",
      "stage": "integration",
      "attribute": "台詞感情",
      "category": "正常",
      "given": "台詞が非 Neutral の感情型（TRDT）を持つ",
      "spec": "抽出した台詞感情が、その台詞のプロンプトへ乗る（感情 staging→プロンプトの合流）"
    }
  ]
}
```

書き方の規約:

- ドメイン語彙で書く。table・列・status・訳値・prompt 文字列・file・method を `spec` と `given` に入れない。
- 固有名は箱（叙述文・固有名・台詞・定型句・無訳片）に畳み、REC:FIELD 種別（WEAP・BOOK…）は列挙しない。
- 固定名（Female flag・TRDT・ElderRace など）は残してよい。残すときは日本語で意味を補う。

## 継ぎ目と入力

結合オラクルはツール単位で入口→出口を取る。守るのは単体で守れない継ぎ目だけ。

- C# 抽出ツール（`extraction` 段）: 入口＝合成 esm を実抽出（`PluginExtractor.Extract`）、出口＝抽出結果／staging。合成 esm は独立生成スクリプト（`tools/synthetic-fixture`、Mutagen で master 無し自己完結 plugin）が作り、`test-oracle/fixture/Synthetic.esm` へ commit。生成はテスト・build から切り離し fixture 変更時だけ手動で回す。
- Go 翻訳ツール（`integration` 段）: 入口＝既存の合成 fixture（`SyntheticFixture()`）を seed に `SyntheticRun`、出口＝`Capture`（送信プロンプト列・翻訳件数・最終 DB）。C# 抽出出力を Go へ流し込む結合はしない（1 プロセスで連結しない）。
- 守る対象は継ぎ目だけに絞る。単段で純粋に閉じるルール（口調閾値・stoplist 判定・固有名派生・役割語引き等）は core package の既存単体が守るため、オラクルに入れない。
- master や localized が要る 3 spec（`override-identical-not-counted`・`container-owns-child-counted`・`localized-header-excluded`）は master 無し esm で再現できないため、既存の実 esm 単体テスト（`CountParityTests` ほか）へ委ねる。specs.json では `coverage: "existing-unit"` を付け、C# 網羅番人は担保済みとして扱う。
- 翻訳 provider は決定的 fake（`RecordingProvider`）で固定。LLM 忠実性はオラクル範囲外（実画面確認へ委ねる）。
- 製品の実 .esm・実 LLM・実行時実データを使わない（合成 esm は使う）。出力を入力とコードだけに依存させる。

## 完了定義

- specs: C# 抽出の継ぎ目（`extraction` 16 件＝assert 13 ＋ 委任 3）と、Go 翻訳の継ぎ目（`integration` 6 件）を持つ共有 JSON（`test-oracle/specs.json`）。
- C# 側: 合成 esm を実抽出の入口で読み、`extraction` 段 spec を抽出結果／staging から assert（1 オラクル 1 関数、`[Oracle]` id）。
- Go 側: `SyntheticRun` を 1 回通し、read-only の出口（`Capture`・最終 DB）へ `integration` 段 6 件を assert（1 オラクル 1 関数、登録 map）。
- 網羅番人: 各側に、自側 stage の spec id と登録テスト id の完全一致を見るメタテストを置く。spec を書いてテストを書き忘れたら落ちる。
- golden 文字列比較（`TestSyntheticNonRegression`・`testdata/synthetic.golden`）を撤去。決定性テスト（`TestSyntheticDeterministic`）は残す。
- テストの書き方は各テストフォルダの `CLAUDE.md` に固定（1 オラクル 1 関数・id を引ける・AAA・独立・given は入力側）。
- 差込点を置くだけで観測できない状態を「動く」と書かない。空 assert・仮実装で通したことにしない。

## テスト設計

各テスト = 入口へ入力を通し、出口へ spec を assert する。ツール単位で入口→出口を取り、途中の package 継ぎ目を外から単体のように確かめない。

| stage | 担当 | 入口 | 出口（観測点） |
|---|---|---|---|
| `extraction` | C# | 合成 esm → `PluginExtractor.Extract` | 抽出結果 `ExtractionResult` / `extracted_info_emotion`・`extracted_info_condition` staging |
| `integration` | Go | 合成 fixture → `SyntheticRun` | `Capture`（送信プロンプト列・翻訳件数）・最終 DB（narration/proper_noun/line 等） |

- 網羅番人: 各側で specs.json を stage 絞りし、登録テスト id 集合との完全一致を見る。
- 守らない対象: 単段で純粋に閉じるルール（core 既存単体に委任）、state derive・Wails 周辺（実画面確認）。オラクルは継ぎ目だけを見る。
- provider は決定的ダミー固定。訳品質は対象外（実画面確認）。

## 含まない（スコープ外）

- 永続ライフサイクル（plugin 単位取得・削除カスケード）。既存単体テストに委ねる。
- 単段で純粋に検証できるルール。既存単体テストに委ね、重複させない。
- AI 訳の妥当性。実画面確認へ委ねる。

## 軽 / 重判定

- 画面が動くか: **N**。テスト追加・差し替えのみ。
- `docs/architecture.md` 反映が要るか: **N**。層・依存・Wails 境界は不変。
- 判定: 両方 N → **軽 task**。

## 実装進捗

verified（実測済み）:
- 基盤: Mutagen 0.53.1 で master 無しの自己完結 esm を組み disk へ書き、実抽出（`PluginExtractor.Extract`）で読み戻せることを実測。
- 生成器: 独立 console `tools/synthetic-fixture`（`SyntheticEsmBuilder` ＋ Program）が合成 esm を書き出す。合成 esm を `test-oracle/fixture/Synthetic.esm` へ commit 済み。
- builder の given を extraction 段の非委任 13 spec 全てへ拡張し、実抽出で観測できることを確認:
    - 条件由来話者（GetIsID/GetIsRace/GetInFaction/GetIsVoiceType の肯定条件）、否定条件の除外（GetIsID==0）、声型 VTYP のみ採用（FLST はプール扱いで除外）。
    - 継承元の名（TPLT 連鎖）で話者名を解決（`InheritedGuard`→`Aventus Aretino`、種別「形態」）。
    - 条件性別: GetIsSex(Male)==0 を女へ畳む、男女混在 FLST は性別を定めない（行なし）。
    - 複数応答の ordinal 採番、感情型 Fear/Anger 記録・Neutral 除外。
- C# 側オラクルテスト（`tools/extractor.Tests/OracleExtractionTests`）が green: `dotnet test tools/extractor.Tests` で合計 34 件通過。内訳は extraction 段 13 spec の照合＋委任確認番人。
    - loader `OracleSpecs`（specs.json を parse）＋ `[Theory]` の default 分岐が網羅番人（spec 追加で assert 未登録なら落ちる）。委任 [Fact] が existing-unit の 3 件を明示台帳にする。
- 既存 C# 単体テスト・Go は未変更で無影響。

Go 側（integration 段）verified:
- 設計の訂正: オラクルは「段 × 属性」の全列挙でなく、単体で守れない継ぎ目だけを見る。単段で純粋に閉じるルール（口調閾値・stoplist 判定・固有名派生・役割語引き等）は core package の既存単体が守るため、specs.json から外した。go 段を継ぎ目 6 件（stage=`integration`）へ刈り込み、単段ルール specs を削除した。
- 6 件の継ぎ目: 件数保存（`count-parity`）・未知除外（`box-routing-unknown-skipped`）・固有名一貫（`proper-noun-consistent`）・感情結線（`line-emotion-injected`）・話者結線（`speaker-tone-injected`）・stoplist 一貫（`stoplist-preserved-in-body`）。
- Go 結合テスト（`internal/harness/oracle_test.go`）: 入口 `SyntheticRun` を 1 回通し、read-only の出口（`Capture` と最終 DB）を spec ごとに照合。1 spec 1 関数、id は登録 map、網羅番人は登録 id 集合と specs.json の integration 段 id の完全一致。
- fixture 最小追加: `SyntheticFixture` に感情 staging（`extracted_info_emotion` の Fear）と、フルネーム固有名を含む台詞（0x520）を追加。感情結線と固有名一貫の given を通す。
- golden 文字列比較（`TestSyntheticNonRegression`）と `testdata/synthetic.golden` を撤去。決定性テスト（`TestSyntheticDeterministic`）は残す。
- テストの書き方は各テストフォルダの `CLAUDE.md`（`tools/extractor.Tests/CLAUDE.md`・`internal/harness/CLAUDE.md`）に固定: 1 オラクル 1 関数・id を引ける・AAA・独立・given は入力側。
- 検証: `npm run verify:backend`（go test・arch-lint・boundary）通過。C# `dotnet test tools/extractor.Tests` 34 件通過。

## close_conditions

- C#/Go 双方の統合テストが、共有 JSON の各 spec を入口から叩いて期待どおり assert して通る。
- 既存の golden 文字列比較が撤去され、決定性テストは残る。
- `npm run verify:backend`（`go test`・arch-lint・boundary 走査）と C# の `dotnet test tools/extractor.Tests` が通る。

## finalize

### 正本化判断

- `docs/architecture.md` 反映: **不要**。層・依存・Wails 境界は不変（軽 task）。テスト追加・差し替えと合成 fixture 生成器のみで、層構成に触れない。
- `docs/er.md` 反映: **不要**。schema 変更なし（既存 `extracted_info_emotion` へ seed するのみ、新 table・列なし）。
- 人間承認済みの恒久仕様: **なし**。正本反映は行わない。
- 判断・履歴は `docs/changelog.md` へ記録した。

### 作業 commit / merge

- 作業 commit hash: `d7ca5df0`（branch `claude/synthetic-pipeline-oracle`）
- local merge: `master` へ `--no-ff`（merge commit `0f764fb1`）、conflict なし
- merge 後検証: `npm run verify:backend`（go test・arch-lint・boundary）通過、C# `dotnet test tools/extractor.Tests` 34 件通過
- remote 変更なし（push・tag・remote branch 削除は行わない）
