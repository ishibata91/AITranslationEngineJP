# 辞書に無い漏れ語の抽出方法研究 — 評価結果

実施日: 2026-07-04〜2026-07-05。実装は `internal/core/mention/candidate.go`（純粋ルール）と
`internal/core/mention/prose.go`（prose の固有表現・品詞解析 adapter）、評価ハーネスは
`cmd/poc-missing-term`。

## 評価方法（既知語 held-out 方式）

1. C#/Mutagen 抽出器で評価 DB を作る（`tools/extractor` の `--sqlite` 出力。`tmp/missing-term-eval/`）。
2. 取込段と同じ振り分け（`engine.Dispatch`）で narration・line の本文を得て、辞書
   （master_term ∪ proper_noun、stoplist 選別後）の言及が出る本文を言及検出（`mention.Detector`）で特定する。
3. 安定キー整列＋固定シード（seed=1）の決定的サンプリングで N=1000 件を採る。
4. サンプル本文の言及語（正解ラベル）を辞書から隠し、候補検出へ「未知語」として見せる。
5. recall（隠した語のうち検出できた割合）・代理精度（候補のうち隠した語だった割合）・重複を測る。

正解ラベルの選別（ハーネス側）: 英字を含み固有名詞の表記形（大文字始まりの語と接続語の並び）の
辞書行だけを隠す対象にする。記号だけの行（"."）・文の形の行（"Let's go"）は固有名詞でないため
正解にしない。動的タグ（`<Alias=...>`）内だけの言及と、本文全体が名前そのものの行も対象外。

## 結果

| 指標 | dev（Skyrim.esm） | held-out（Dawnguard.esm＋Dragonborn.esm＋inigo.esp） |
| --- | --- | --- |
| サンプル | 1000 件（言及付き 14,473 件から） | 1000 件（内訳 270／454／276） |
| 正解（隠した語） | 550 語 | 395 語 |
| recall | **96.7%**（532/550） | **95.4%**（377/395） |
| 重複（正規化後同一の複数出力） | **0** | **0** |
| 決定性（3 回実行のブレ） | **±0 ポイント** | **±0 ポイント** |
| 候補数 | 1,989 語 | 1,267 語 |
| 代理精度 | 26.7% | 29.8% |
| 真の精度（推定） | **約 63%** | **約 54%** |

- 実行条件: `--n 1000 --seed 1 --ner --runs 3`。held-out の plugin は開発中一切参照していない
  （dev 調整はすべて Skyrim.esm で行い、held-out は最終評価で初めて実行した）。
- 達成基準の判定: 1（recall 95%）✓、3（重複 0）✓、5（stoplist 除外）✓（供給選別で機構的に成立）、
  6（held-out 2 plugin 以上）✓（3 plugin、Bethesda DLC 2 本＋community mod 1 本）、
  7（決定性 ±3pt 以内）✓（完全一致）、4（汎化）recall・重複は held-out でも成立 ✓。
  2（精度）は下記のとおり当初閾値 90% に届かず、閾値を調整した（goal.md の調整条項に基づく）。

## 精度の実測方法と内訳

代理精度（候補のうち隠した語だった割合）は保守的な下限値になる。辞書に本当に無い固有名
（本 task の本来の獲物。書籍内の人名・地名・lore 用語）も「誤検出」側に数えるため。

真の精度は、誤検出の系統抽出サンプル（dev: 1/7 で 209 件、held-out: 1/5 で 178 件）を
人手ラベリング（固有名詞か否か。goal の除外 4 分類 = 一般語・助詞・記号・語の断片、に該当すれば否）
して推定した。

- dev: 誤検出 1,457 語中、標本 209 件の固有名詞率 50.2% → 真の精度 ≈（532＋0.502×1457）/1989 ≈ **63%**
- held-out: 誤検出 890 語中、標本 178 件の固有名詞率 34.8% → 真の精度 ≈（377＋0.348×890）/1267 ≈ **54%**

誤検出（固有名詞でない側）の主な内訳:

- 大きい固有名の断片（Raven Rock の Rock、Tyranny of the Sun の Sun など）。全体と部分の両建て
  出力は recall 確保（辞書が全体と部分の両方を語彙に持つ形への対応）と表裏一体。
- 書籍本文の見出し・整形由来の複合（Volume Three・Living Room 等）。見出し選別で削ると
  一般語複合の正解（Fire Salts・Giant's Toe 等の item 名）が巻き添えになることを実測で確認し、撤回した。
- 文頭・NER 救済で拾った一般語 1 語（Apologies・Observation 等）。NER 救済を外すと recall が
  約 1.5 ポイント下がることを確認し、維持した。

## recall と精度の両立に関する結論

LLM 不使用の決定的手法（大文字ヒューリスティック＋用法分布＋prose の NER・品詞）では、
recall 95% 以上を保ったまま真の精度 90% に到達しない。精度を上げる選別（見出し選別・
機能語先頭の句全体抑止の拡大・複合と断片の一方だけの出力）はいずれも recall を 95% 未満へ
下げることを dev の実測で確認した。到達した運用点は recall ≈95〜97%・真の精度 ≈54〜63%。

この運用点を採る製品上の理由: 候補は人間レビューを経て辞書へ入る前提であり、見逃した語は
後工程で回収できない（recall が硬い制約）のに対し、誤検出はレビューで捨てられる（精度は
レビュー工数の問題）。

## 既知の限界（取りこぼしの残り）

- 上位句に埋まった語（Gray Fox の Fox、Dragon's Tongue Flower の Dragon's Tongue）。
  言及検出の照合粒度と句の粒度の不一致で、全部分列の出力なしには回収できない。
- 一般名詞形の辞書語が文頭・一般語用法で出る場合（Enemies are displayed…の Enemies）。
  文脈上は固有名詞の用法でないため、正解ラベル側のノイズに近い。
- 一人称代名詞との同綴り（Chapter I の I）、数字始まりの固有名（7,000 Steps）。
- `&` 込みの名前は対応済み（Kolb & the Dragon）。二重ハイフンはダッシュとして語から分離する。

## 再現手順

```
dotnet run --project tools/extractor -- --data dictionaries/Data --plugin Skyrim.esm --sqlite tmp/missing-term-eval/dev-skyrim.sqlite3
go run ./cmd/poc-missing-term --db tmp/missing-term-eval/dev-skyrim.sqlite3 --n 1000 --seed 1 --ner --runs 3 --dump 40
go run ./cmd/poc-missing-term --db tmp/missing-term-eval/heldout-dawnguard-esm.sqlite3,tmp/missing-term-eval/heldout-dragonborn-esm.sqlite3,tmp/missing-term-eval/heldout-inigo-esp.sqlite3 --n 1000 --seed 1 --ner --runs 3
```

誤検出の全件書き出しは `--fp-out <path>`（人手ラベリング用 TSV）。
