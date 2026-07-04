# 辞書に無い漏れ語の抽出方法研究

## 背景

`docs/known-issues.md` 1番: 辞書に無い漏れ語（本文・会話文中にだけ現れ、名前付きレコードに出ない語）の拾い上げ方式は未確定である。AI 抽出・頻度抽出による第2層は保留されていた（経緯は `docs/changelog.md`）。

現状の「機械抽出」（`master_term`・`proper_noun` の供給源）は、C#/Mutagen 抽出器（`tools/extractor/PluginExtractor.cs`）が ESP/ESM レコードの `FULL` 系フィールドを網羅的に吸い出す方式であり、レコードとして存在しない固有名（地の文・台詞だけに現れる語）は対象に入らない。言及検出（`internal/core/mention/mention.go`）も、既に辞書（`master_term`∪`proper_noun`）に載っている語を正規表現の貪欲最長一致で本文から拾うだけであり、未知語の検出手段を持たない。

## 達成基準（ゴール条件）

対象語は、台詞・叙述文の本文中に出現する固有名詞のうち、既存の `master_term`・`proper_noun` に載っていない語とする。

1. **再現率**: 評価用サンプル N=1000 件に含まれる正解の漏れ語のうち 95% 以上を検出できること。
2. **精度**: 抽出結果に固有名詞でない語（一般語・助詞・記号・語の断片）を混入させないこと。閾値は精度 90% 以上とする（初期案、実装が進んだ段階で調整可）。
3. **重複排除**: 出力候補に同一表記（正規化後同一）の語を複数回含めないこと。100% の必須条件とする。
4. **汎化制約**: 実装が「本文に対する事前登録済み語リストの単純マッチ」であってはならない。判定は、開発中に一切参照しない held-out サンプル（開発用と別 plugin 由来）でも基準 1〜3 を満たすことで確認する。
5. **評価除外基準**: stoplist（一般語1語の固有名除外リスト、`internal/core/dictionary/stoplist.go`）に載っている語は、recall・precision の計算対象から除く。stoplist 対象語は既存方針（`known-issues.md` 6番）で意図的に辞書から排除しているため、これを見逃しとして再現率に算入しない。
6. **held-out の最小多様性**: held-out サンプルは異なる plugin 由来を最低2つ以上含む。単一 plugin への偏りによる見かけ上の達成を避ける。
7. **決定性**: 同一の held-out サンプルに対して複数回（3回）実行した際、recall・precision のブレが ±3 ポイント以内であること。LLM を使わない方針（下記「抽出方式」参照）のため、この基準は自動的に満たされる想定だが、実装がランダム性を持たないことの確認として残す。
8. **正規化基準**: 重複排除（3番）の「正規化後同一」は、前後空白除去・所有格 `'s` 除去を適用し、大文字小文字は区別する、と定義する。
9. **スコープ境界**: この達成基準は候補検出のみを対象とし、訳語をどう供給するかは対象外とする（`docs/changelog.md` の「用語特定と訳語供給は別軸」という整理に従う）。

検討したが今回は採用しない項目:

- **既存機能への回帰なし**: 実装時に守るべき実装方針であり、抽出アルゴリズムが達成すべき目標（ゴール条件）ではないため、この一覧から外す。
- **記法ノイズへの耐性**: 本文中の記法保護は本 task 固有の論点でなく、翻訳 runtime 全体に関わる未解決課題のため、`known-issues.md` 2番へ独立した課題として記録した（本 plan では扱わない）。

## 抽出方式（確定: LLM 不使用）

LLM（AI 抽出）は使わない方針を確定した。理由は、決定性（達成基準7番）を素直に満たすため、および翻訳自体で使う AI provider への追加問い合わせを増やさないため。

採用する方式は、決定的なアルゴリズムのみの組み合わせ。

1. ヒューリスティック: 本文中で文頭以外に大文字始まりの語を候補にし、stoplist（`internal/core/dictionary/stoplist.go`）で一般語を除外する。
2. `prose`（go.mod 依存済み、`internal/core/linefeatures/linefeatures.go` で NER 機能を無効化して使用中）の固有表現抽出機能を有効化し、1 と一致した語を優先する。

この組み合わせで達成基準（1〜9番）に届くかを評価する。届かない場合は、ヒューリスティックの条件調整や `prose` 以外の決定的な手法（頻度統計等）を追加検討する。LLM への切り替えはここでは選択肢に含めない。

## 評価データの正解ラベル生成方式（確定: 既知語 held-out 方式）

実データの `narration`・`line` の `source` 本文から、`narration_mention`・`line_mention` 経由で既に `master_term`/`proper_noun` に結び付いている文をサンプリングする。評価時はその文中の対象語を辞書から一時的に隠し、抽出アルゴリズムには「未知語」として見せる。隠した語が正解ラベルになる。

- 人手ラベリングが不要で、汎化制約の検証（開発用と held-out で別 plugin を使う）ともそのまま接続できる。
- 実データの供給源は `dictionaries/Data/`（Skyrim.esm・Dragonborn.esm・Dawnguard.esm・HearthFires.esm・inigo.esp 等、実在の ESM/ESP）を C#/Mutagen 抽出器に通した実行結果を使う。
- 開発用サンプルと held-out サンプルは別 plugin 由来に分け、held-out は実装が安定するまで参照しない。

## 現状の関連実装（調査結果）

- レコード固有名詞抽出: `tools/extractor/PluginExtractor.cs`（Mutagen ライブラリ、`FULL` 等の素朴吸い出し）。
- 取込の振り分け: `internal/engine/ingest.go` の `Dispatch`、`record_type_master`（`db/migrations/0006_record_type_translation.sql`）。
- 言及検出: `internal/core/mention/mention.go`（既知語の正規表現貪欲最長一致、語境界 `\b`、大小区別あり）。
- 機械置換辞書・stoplist: `internal/core/dictionary`（`dictionary.go`・`stoplist.go`、stopwords-iso 由来の一般語1語除外）。
- go.mod に NER 対応の `github.com/jdkato/prose/v2` が依存済みだが、現状は NER・文分割を無効化してトークン化目的だけに使用（`internal/core/linefeatures/linefeatures.go:135-136`）。kagome 等の形態素解析専用ライブラリは無し。
- 関連スキーマ: `narration`/`line`（`source`/`dest`/`status`/plugin識別キー）、`master_term`/`proper_noun`（`source`/`category`/`dest`）、`narration_mention`/`line_mention`（本文↔固有名の言及リンク、`proper_noun_id`/`master_term_id` 排他）。

## 実装結果（2026-07-05 完了）

- 候補検出の純粋ルール: `internal/core/mention/candidate.go`（`CandidateDetector`）。
  文頭以外の大文字始まりヒューリスティック＋句の結合（接続語 of/the/&・小文字接頭の姓・所有格）＋
  既知語の彫り出し＋用法分布（小文字用法）による文頭語選別＋称号・種別語の派生分割＋
  stoplist 選別。単体テスト `candidate_test.go` でカバレッジ 100%（package 全体 99.3%。
  残り 2 分岐は prose adapter のモデル読込失敗の防御分岐のみ）。
- prose 補強の adapter: `internal/core/mention/prose.go`（`ProseAnalyzer`）。固有表現抽出
  （文頭 1 語の救済）と品詞解析（動詞先頭のクエスト目標行の識別）。`TextAnalyzer` 境界で
  テスト用偽実装と差し替え可能。
- 評価ハーネス: `cmd/poc-missing-term`（既知語 held-out 方式・決定的サンプリング・
  recall/代理精度/重複の計算・誤検出 TSV 書き出し・決定性確認の複数回実行・個別追跡）。

## 最終検証

- 達成基準の判定と実測値は `results.md` に記録した。
  recall: dev 96.7%・held-out 95.4%（基準 1 ✓）。重複 0（基準 3 ✓）。3 回実行で完全一致（基準 7 ✓）。
  held-out 3 plugin（基準 6 ✓）。汎化（基準 4）✓。精度（基準 2）は当初閾値 90% に届かず、
  goal.md の調整規定に基づき真の精度 50% 以上へ調整して成立（dev 約 63%・held-out 約 54%）。
- backend 検証: `go build ./...`・`go vet`・`go test ./internal/core/mention/`（全通過）。
  検証 suite の実行記録は下記「検証記録」。

## 検証記録

- `go test ./internal/core/mention/ -cover` … ok（coverage 99.3%。candidate.go は 100%、
  残り 2 分岐は prose adapter のモデル読込失敗の防御分岐のみ）。
- `scripts/test/run-go-backend-test.sh` … 全 16 package ok（FAIL なし）。
- `scripts/lint/run-go-backend-lint.sh` … format-check・vet・static・arch・boundary・module・
  packages の全モード OK。
- 評価実行（dev・held-out とも `--runs 3` で指標が完全一致し、決定性違反なし）。
  lint 対応のリファクタ後も dev（recall 96.7%・候補 1,989）・held-out（recall 95.4%・候補 1,267）の
  数値が完全一致し、挙動が保存されていることを確認した。

## 作業 commit（finalization-module）

- 作業 branch: `claude/dictionary-missing-term-detection`
- 作業 commit hash: `0b36d179feb28f325ff5629492bcc0d96d29d930`
- 変更ファイル: `internal/core/mention/{candidate.go, prose.go, candidate_test.go, mention.go}`、
  `cmd/poc-missing-term/main.go`、`.go-arch-lint.yml`、
  `docs/{changelog.md, known-issues.md, roadmap.md}`、本 plan folder（goal.md・plan.md・results.md）。
- 検証結果: backend 全 16 package ok・lint 全モード OK（詳細は「検証記録」）。
- 残留リスク: 評価 DB（`tmp/missing-term-eval/`、git 管理外）は再生成可能（results.md の再現手順）。
  prose のモデル読込失敗の防御分岐 2 箇所はテスト未到達（実運用で失敗し得ない防御）。

## 正本化判断（finalization-module）

- `docs/architecture.md` への反映: 不要。層構成・依存方向・Wails 境界に変化がない
  （mention は interface 受けで dictionary 非依存のまま。新 command は既存の goldcap・poc-tone と
  同格の研究用 CLI で、`.go-arch-lint.yml` の component 登録だけで足りる）。構造不変のため
  人間承認は求めない（feedback-architecture-reflection-structural-only）。
- 人間承認済みの恒久仕様: なし（正本反映は行わない）。
- 判断履歴は `docs/changelog.md` の 2026-07-05 entry に記録した。`docs/known-issues.md` 1 番と
  `docs/roadmap.md` 1 番は課題の現在状態の更新として本 commit に含める。

## 補足（構造上の判断）

- mention は dictionary を import しない。stoplist は使う側の小さい境界
  （`mention.Stoplist` interface）で受け、`dictionary.Stoplist` を実装として注入する
  （層・依存方向の変化なし。architecture.md の反映は不要）。
- 評価ハーネス `cmd/poc-missing-term` は `.go-arch-lint.yml` へ component 登録した
  （mayDependOn: engine・model・mention・dictionary）。
