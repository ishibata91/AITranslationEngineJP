# Task Plan: dialogue-tone-naturalness

会話の訳から機械翻訳感を消す実験 task。`/loop` で仮説と実測を繰り返す形で進める。
何をどこまで達成したら終わりかは `criteria.md`、各回の判断は `loop-log.md`、測った値は `measurements.csv` が持つ。

## 引き継ぎ

準備の段（`experiment-workflow`）の途中。作業する機械を Mac から Windows（RTX4090 と 9800X3D）へ移した。`tmp/` の道具と凍結した標本、中心 DB `db/aitranslation.dev.sqlite3` は人間が手で運んだ。

**済んでいること**

- `Skyrim.esm` の抽出、`line` への取り込み、話者の基底口調の生成。
- 標本の凍結。開発用 598 件・評価用 598 件を `tmp/dialogue-tone-naturalness/dataset/` へ書いた。以後選び直さない。
- 運んだ DB で凍結標本 1,196 件すべてが `line.id` で引け、原文が一致し、基底口調が存在することの確認。
- `criteria.md` の人間承認。達成条件の表・選択肢の表・標本の表の全ての行が 2026-07-27 に承認済み。
- 測る道具 4 つ。判別器（`discriminator/discriminate.py`）、診断と破損行（`main.go`）、読む 60 件を選ぶ道具（`pickreview`）、読んだ結果を数える道具（`tallyreview`）。

**次の一手**

1. LM Studio（`http://192.168.0.226:1234`、この Windows マシン自身）を上げ、`gemma-4-12B-it-QAT` を載せる。
2. `go run ./tmp/dialogue-tone-naturalness/translate -set dev -model <API の model 名> -round 1` で開発用 598 件を訳す。続けて `-suffix -b` を付けてもう 1 度訳し、生成の揺れを測る。
3. 判別器と診断の道具を回し、回 1 を `measurements.csv` へ書く。`pickreview` が選んだ 60 件を `fresh` が読み、`tallyreview` で数える。
4. `criteria_reviewer` の達成条件レビューを通す。通過するまでループの段へ入らない。

**砂場は commit に入らない。** 道具と凍結した標本と訳した出力は `tmp/` にあり `.gitignore` の 14 行目で除外される。中心 DB も `db/*.sqlite3` で除外される。別マシンで続ける場合は、`tmp/` と中心 DB を手で運ぶ。凍結した標本は `line.id` で DB を引くので、DB を作り直すと標本が引けなくなる。

## 事象

人間が見たことだけを書く。

- 会話の訳が全体的に丁寧すぎて、機械で訳したような硬さが残る。
- 意訳が足りず、英語の構文をそのまま日本語へ写した訳になっている。
- クエストログと叙述系（物品説明、効果説明、書物、世界観断片）の訳は現状でも良い。

## 対象

会話だけを対象にする。`INFO:NAM1`（NPC の応答）で、標本を取る plugin は `Skyrim.esm`。

事象は `inigo.esp` の訳を見て起きたが、標本は `Skyrim.esm` から取る。理由を 2 つ書く。

- `inigo.esp` の会話は話者が adult male 5,101 件と adult female 143 件に偏り、年配と子供を 1 件も持たない。口調テンプレートは年齢と性別で引くので、この標本では区分ごとの効き目を確かめられない。
- `inigo.esp` の会話は公式訳を 8 件しか持たない。`Skyrim.esm` は全ての行に公式訳が 1 対 1 で付く。

口調の対象（例文・指示文・種族注記）は plugin をまたいで共通なので、`Skyrim.esm` で効いた変更は `inigo.esp` の訳へも届く。

変える対象。

- `assets/role-speech-examples.tsv`: 口調の例文 57 行。区分ごとに原文と訳例の対を持ち、翻訳の指示文へ差し込まれる。
- `assets/role-speech.tsv`: 一人称と言い回しの傾向。区分ごとに文字列で持つ。上の例文と 1 つの要因として一緒に振る。理由は `criteria.md` の要因の振り方が持つ。
- `internal/core/personatone`: 区分から口調指示を組む処理。種族に応じた注記を足す分岐を含む。
- 口調の指示文（`record_type_translation` の `口調`）: app が使う DB の実値。画面編集済みで migration のどの文面とも一致しないため、変えるなら画面から入れる。

変えない対象。

- 叙述系とクエストログの訳。人間が見て既に良いと判断しているため、変える理由が無い。
- `DIAL:FULL`（プレイヤーの選択肢文）。話者が NPC でなく、口調の指示が別の経路で組まれるため、実験の要因が増える。
- 翻訳に渡す文脈の量（会話の場面、手前の連鎖）。やらないことへ回す。
- 凍結した標本。2026-07-27 に 1 度だけ選んで `tmp/dialogue-tone-naturalness/dataset/` へ書いた。ループは選び直せない。

## 砂場

- 場所: `tmp/dialogue-tone-naturalness/`。実験の道具、測定結果、標本を読んだ結果（`sample-review.jsonl`）を置く。
- 除外の根拠: `.gitignore` の 14 行目に `tmp` があり、`git check-ignore -v tmp/dialogue-tone-naturalness/main.go` で 14 行目に当たることを確かめた。

## 接続先

- 外部の接続先: `http://localhost:11434`（ollama）。翻訳を回す。使うモデルは `translategemma:12b` に固定し、回の間で変えない。
- 機械: 翻訳も訳文らしさの判定も RTX4090 と 9800X3D の Windows マシン 1 台で回す。測る手段の詳細は `criteria.md` の選択肢の表が持つ。
- 使わないモデルと理由: `hy-mt2-7b` は空応答と記号混入を起こす。`google/gemma-4-12b-qat` と `qwen3.5:9b` は応答の前に推論を書くため 1 件に 7 秒前後かかり、598 件を 20 回ほど繰り返す量に合わない。`gemma4:e2b` は同じ 1 件で 3.7 秒かかり `translategemma:12b` より遅い。
- データ: 中心 DB `db/aitranslation.dev.sqlite3`。標本を取るため 2026-07-27 に `Skyrim.esm` を 2 つ目の対象 plugin として抽出し、`line` へ 40,833 件を取り込み、話者 909 人の基底口調を生成した。抽出前の DB は `db/aitranslation.dev.pre-skyrim-backup.sqlite3` に控えてある。公式日本語訳は `reference_translation` の `INFO:NAM1` 41,159 件。

## branch 情報

- `execution_branch`: `claude/dialogue-tone-naturalness`
- `target_branch`: `master`
- `source_commit`: `ca4b6ba32226f85d3255b4eae24166e8defbc00a`

## やらないこと

- 会話の場面を翻訳へ渡すこと（会話の種類、クエストの場面、応答が答えている選択肢文、同一応答の続き行）。後続 task で扱う。
- 会話の手前の連鎖を翻訳へ渡すこと（`LinkTo` の 1 段逆引き、作成順による並びの近似）。後続 task で扱う。
- 訳文の破損の原因を直すこと。件数を診断として出すところまでとする。
- 叙述系とクエストログの訳の書き方を変えること。
- プロダクト正本（`docs/architecture.md` など）への反映。実験の採否が決まった後に、本 task の外で人間が決める。
