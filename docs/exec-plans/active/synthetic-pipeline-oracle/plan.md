# Plan: synthetic-pipeline-oracle

- `task_id`: `synthetic-pipeline-oracle`
- `working_branch`: `claude/synthetic-pipeline-oracle`（base: `master`、分岐元 commit: `37fb6b03`）
- `target_branch`: `master`

## 依頼要約

単体テストでは届かない「エンティティ・クラスをまたいだ合成結果が仕様どおりか」を、C# 抽出機と Go メインシステム（翻訳）の双方で確かめる統合テストを作る。E2E は重く、AI の実画面確認が E2E 相当の確信を担うため、単体とは別の統合テストを置く。判定基準は共有のテストオラクル（`test-oracle/`）に持ち、C# と Go が同じ JSON を読んで照合する。

## オラクルの設計

判定基準を `test-oracle/specs.json` に共有 JSON で持つ。スキーマと規約は `test-oracle/README.md` が正本。要点だけ再掲する。

- 粒度: 処理段（stage）中心。1 エントリ = stage × 属性。UC は 1 本しかないので割る軸にしない。
- stage: `extraction`(C#) / `ingest`(Go) / `prompt-build`(Go) / `post-process`(Go) / `end-to-end`(両方)。適用系は stage から一意に決まる。
- `given` は入力 esp のリッチさの一断面。`category`（正常/異常）は given から従属。
- 6 フィールド: `id`・`stage`・`attribute`・`category`・`given`・`spec`。`id` が C#/Go 両テストの対応確認の join キー。
- ドメイン語彙で書き、table・列・訳値・prompt 文字列を入れない。

## 継ぎ目と決定性

- 継ぎ目はテストダブルで跨ぐ。C# は master 無しの自己完結 plugin を Mutagen で構築して抽出、Go は抽出後の DB を seed する。1 プロセスで C# と Go を結合しない。
- 翻訳 provider は決定的 fake で固定する。AI 訳の妥当性（注入語をモデルが保持するか等）はオラクルの範囲外で、実画面確認へ委ねる。
- 実 .esm・実 LLM・実行時実データを使わない。出力を入力とコードだけに依存させる。

## 検証するシナリオ

`test-oracle/specs.json`（現状 54 件）が正本。stage × 属性で正常/異常を対にして網羅する。C#/Go の各テストは spec の `id` を参照し、その `given` の fixture を入口から叩いて `spec` を assert する。

## 完了定義

- fixture: リッチさの違うレコードを 1 本に混載した synthetic esp（各 spec の `given` が名指す実体）。
- Go 側: 合成入力で翻訳の入口を叩き、`stage` が go・両方 の spec を送信プロンプト列と最終 DB から assert する。
- C# 側: 合成 plugin を構築して抽出の入口を叩き、`stage` が C#・両方 の spec を抽出後 staging から assert する。
- 既存の golden 文字列比較（`TestSyntheticNonRegression`・`testdata/synthetic.golden`）はオラクル照合へ置き換える。決定性テスト（`TestSyntheticDeterministic`）は残す。
- 差込点を置くだけで観測できない状態を「動く」と書かない。空 assert・仮実装で通したことにしない。

## 含まない（スコープ外）

- 永続ライフサイクル（plugin 単位取得・削除カスケード）。既存単体テストに委ねる。
- 単段で純粋に検証できるルール。既存単体テストに委ね、重複させない。
- AI 訳の妥当性。実画面確認へ委ねる。

## 軽 / 重判定

- 画面が動くか: **N**。テスト追加・差し替えのみ。
- `docs/architecture.md` 反映が要るか: **N**。層・依存・Wails 境界は不変。
- 判定: 両方 N → **軽 task**。

## close_conditions

- C#/Go 双方の統合テストが、`test-oracle/specs.json` の各 spec を入口から叩いて期待どおり assert して通る。
- 既存の golden 文字列比較が撤去され、決定性テストは残る。
- `npm run verify:backend`（`go test`・arch-lint・boundary 走査）と C# の `dotnet test tools/extractor.Tests` が通る。

## 成果物

- `test-oracle/README.md`・`test-oracle/specs.json`（済）: 共有 JSON オラクルとスキーマ。
- fixture（synthetic esp）・Go 側統合テスト・C# 側統合テスト（未）。

## 現状

- `test-oracle/`（README ＋ specs 54 件）を作成し commit 済み。
- fixture と C#/Go のテスト結線は未着手。
- 以前の作業で入っていた「共有ファイルを作らない素な 2 本」方式の diff（`integration_test.go`・`integration-test-guidelines.md`・C# 統合テスト・golden 削除・docs 参照）は、共有オラクル方針と衝突するため戻した。
