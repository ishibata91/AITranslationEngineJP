# dialogue-meaning-accuracy 実験計画

会話の訳から意味の取り違えを減らす実験。前の task `dialogue-tone-naturalness` が口調（敬体、一人称、文字数）を公式訳の水準へ寄せた後に残った、意味そのものの誤りを対象にする。

## 事象

人間が回 25 の訳文 598 件を読んで見つけた、口調の対象を変えても直らなかった誤りを 1 件 1 行で書く。
いずれも `dialogue-tone-naturalness` の回 1 から回 25 まで一度も直らなかった。

- 前の行を見ないと訳せない台詞で、代名詞の指す先が落ちる。`Don't you remember? His name is Hroggar.` の `His` が消え、訳文が「え？フロガーって名前だよね」になる。誰の名前かが伝わらない。
- 英語の語順と品詞をなぞった訳文が出る。`Those who are ignorant find the Way of the Voice easy to ridicule.` が「声の道というものを嘲笑することに容易さを感じるだろう」になる。公式訳は「無知なる者は、声の道を嘲りやすいようだ」。
- 前の行を見ないと訳せない台詞で、応答が噛み合わない。回 25 の訳文には、質問へ答える形の台詞が、質問を知らない状態で訳されたものがある。

3 つのうち英語の語順については、`dialogue-tone-naturalness` の後に人間が 3 つのモデル（`kaelri/hy-mt2:7b`、`gemma4:e2b`、`gemma4:e4b`）で同じ 10 件を訳し、3 つとも `translategemma:12b` より公式訳に近い語順を出した。
記録は `tmp/dialogue-tone-naturalness/hy-mt2-probe/report.md` が持つ。
この事象がモデル固有の弱点である可能性があり、指示文で動かせるかは分かっていない。

## 対象

### 変える対象

要因を 2 つに絞る。回ごとにどちらを変えたかを 1 つに保ち、原因を分けられる形にする。

| 要因 | 変える対象 | 変える中身 |
| --- | --- | --- |
| A | 翻訳のプロンプトへ渡す前の行の台詞 | 渡さない（現状）／直前の 1 行／直前の 2 行 |
| B | 中心 DB の `directive` テーブル key `口調` の指示文 | 英語の語順と品詞をなぞらせない書き方を足す／足さない |

要因 A は engine の変更を伴う。会話の連鎖（どの台詞の次にどの台詞が来るか）を引く実装が要る。
現在の翻訳経路は 1 行を単独で投げる（`internal/engine/engine.go` に前の行を引く実装が無いことを 2026-07-28 に確認した）。

要因 B は `dialogue-tone-naturalness` で既に 1 行入っている（「原文の内容をすべて訳す。そのうえで英語の語順と品詞をなぞらず、説明を足さずに短く言い切る。」）。
この行があってなお語順の問題が残るので、書き方を変える形で振る。

### 変えない対象

変えない理由を添えて書く。

| 変えない対象 | 理由 |
| --- | --- |
| 例文（`assets/role-speech-examples.tsv`） | `dialogue-tone-naturalness` で口調のために作り込んだ。変えると敬体・一人称・文字数が動き、意味の取り違えの効果と分けられない |
| 役割語（`assets/role-speech.tsv`） | 同上 |
| 基底口調の性質文（`internal/core/personatone` の `toneTraits`） | 同上 |
| モデルと接続先 | `translategemma:12b`、ollama の `http://localhost:11434` に固定する。回の間で変えると比べられない |
| 標本 | `dialogue-tone-naturalness` の凍結標本 1,196 件（開発用 598・評価用 598）をそのまま使う。前の task の値と地続きで読める |
| 固有名の辞書（`master_term`、`proper_noun`） | 別 task `dictionary-lookup-coverage` が扱う |

## 砂場

`tmp/dialogue-meaning-accuracy/` に道具と出力を置く。
`.gitignore` の 14 行目が `tmp` を除外することを、2026-07-28 に `git check-ignore -v tmp/dialogue-meaning-accuracy/x.txt` で確かめた。

`dialogue-tone-naturalness` の道具のうち次を流用する。書き換えずに `-data` と `-out` で場所だけ差し替える。

- `tmp/dialogue-tone-naturalness/translate`: 標本を訳して jsonl へ出す
- `tmp/dialogue-tone-naturalness/tallyreview`: 人が読んだ結果を集計する
- `tmp/dialogue-tone-naturalness/dataset/`: 凍結標本

前の行を渡す形は `translate` の変更を伴う。変更した版は `tmp/dialogue-meaning-accuracy/translate` へ置き、前の task の道具を書き換えない。

## 接続先

実験の条件として固定するものを書く。

| 項目 | 値 |
| --- | --- |
| モデル | `translategemma:12b` |
| 接続先 | ollama の OpenAI 互換エンドポイント `http://localhost:11434` |
| 中心 DB | `db/aitranslation.dev.sqlite3`（`line` 49,378 件、`reference_translation` 72,599 件、`speaker` 926 人） |
| 基本指示文 | `prompt_template.base_directive` 270 文字 |
| 固有名の辞書 | `master_term` 17,346 件、`proper_noun` 13,471 件のうち訳あり 146 件 |

## branch 情報

- 作業 branch: `claude/dialogue-meaning-accuracy`
- 統合先 branch: `master`
- 分岐元 commit: `ca4b6ba3`

## やらないこと

- 口調（敬体、一人称、二人称、文字数比）を良くすること。`dialogue-tone-naturalness` が達成済み。本 task では悪化していないことだけを診断として見る。
- 固有名の崩れを直すこと。原因が辞書の引き当てにあると分かっている。別 task `dictionary-lookup-coverage` へ回す。
- モデルを替えること。`translategemma:12b` に固定する。モデルの比較は `tmp/dialogue-tone-naturalness/hy-mt2-probe/report.md` で済んでいる。
- 会話の連鎖を引く実装をプロダクト正本へ入れること。実験の砂場で試し、採否が決まった後に本 task の外で人間が決める。
- プロダクト正本（`docs/architecture.md` など）への反映。
