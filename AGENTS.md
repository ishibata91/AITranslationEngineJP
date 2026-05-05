# AGENTS.md

会話と作業は日本語を基本にする。
英語の key、既存名、command は必要な時だけ使う。
tmp/code-map/index.jsonにコード地図がある

## 会話ルール

- 見出しなしの長文を避け、2〜4 個の短い見出しに分ける
- 見出しは `##` を使う
- 箇条書きは 3〜6 件に抑える
- 1 段落は 3 文以内にする
- 変更報告は 1 行 1 ファイルにする
- 長い出力の末尾には `SUMMARY` を付ける
- 指示語の使用を原則禁止する。

## 目的

AITranslationEngineJp は Skyrim Mod 向け翻訳エンジンです。
この repo は agent-first で進めます。

## 参照マップ

最初に `.codex/README.md` と使う skill の `SKILL.md` を読む。
agent の `permissions.json` と contract は、skill 本文から agent-owned reference を辿る。

- workflow 正本: `.codex/README.md`
- Codex 入口: `.codex/skills/propose-plans/SKILL.md`
- Codex implementation lane 実装入口: `.codex/skills/implementation-orchestrate/SKILL.md`
- 仕様入口: `docs/index.md`
- 長期原則: `docs/core-beliefs.md`
- 恒久要件: `docs/spec.md`
- architecture: `docs/architecture.md`
- 作業計画: `docs/exec-plans/`

## 強い制約

- Codex は設計、計画、handoff、docs 正本化を担当する
- Codex implementation lane は承認済み `implementation-scope` から実装する
- AI design review は行わず、人間が design bundle を review する
- Codex implementation lane は `docs/`、`.codex/`、`.codex/skills`、`.codex/agents` を変更しない
- docs 正本化は Codex の `updating-docs` だけが扱う

## 実装前に確認すること

1. `.codex/README.md` と使う skill の `SKILL.md` を読み、必要なら agent-owned `permissions.json` と contract を確認する
2. `docs/index.md` から関係する文書だけ読む
3. active / completed plan に同種 task がないか確認する
4. 手をつける前に、影響範囲、実行計画、検証方法を短く固定する

## 実装後にやること

1. 必要な follow-up を plan か issue に記録する
2. docs 正本更新は human 承認済みの時だけ行う
3. 完了した plan を `docs/exec-plans/completed/` へ移す

## 補足

- library の書き方は `npx ctx7 library` / `npx ctx7 docs` で Context7 を確認する
- wails は `npm run dev:wails:agent-browser` で起動する
- ブラウザ操作は `agent-browser` CLI を使う
- UI 証跡は `agent-browser open http://localhost:34115` から取得する
- Sonar project は `ishibata91_AITranslationEngineJP`


# 日本語出力規約 v1

## 優先順位
1. 正確性
2. 可読性
3. 検証可能性
4. 簡潔さ

## 基本文体
- 技術説明は常体で書く。
- 結論を最初の1〜2文で述べる。
- 1文1論点にする。
- 主語・対象・作用を省略しすぎない。
- 事実、推測、提案を分けて書く。
- 同じ概念には同じ用語を使う。
- 回答文は、解釈の幅が狭い日本語で構成する。
- ファイル名、キー、既存名、コマンドだけは固定名として残す。
- 既存名は、出典に完全一致で存在する名前だけを指す。
- 説明のために作った英語の名詞句は、既存名として扱わない。
- 固定名を残す場合も、本文の説明は日本語だけで意味が通るように書く。
- 状態名や判定値は、固定名を残し、日本語で状態の意味を補う。

## 回答前の自己検査
- 回答前に、日本語出力規約へ違反している文を書き直す。
- 英語の名詞句は、出典に完全一致で存在する固定名だけ残す。
- 固定名を残す場合は、日本語で意味、状態、理由、影響を補う。
- 状態、原因、対応、注意を一文に詰め込まない。
- 悪い例に近い文は、良い例の `対象:` `結果:` `原因:` `対応:` `注意:` の形へ分ける。

## 文のルール
- 抽象語より具体語を使う。
- 「これ」「それ」「今回」「この場合」は参照先が曖昧なら使わない。
- 因果は明示する。少なくとも「理由」「影響」「対応」のどれかを入れる。
- 同列項目は同じ抽象度・同じ文型で並べる。
- 読点が3つを超える文は分割を検討する。
- 修飾語は被修飾語の近くに置く。
- 強い断定には根拠を添える。
- 推測は「可能性がある」と明示する。
- 参照、状態、差分、判断を一文に詰め込まない。
- 固定名を複数使う場合は、各固定名の意味と関係を日本語で分けて説明する。

## 箇条書き
- 真に列挙である場合だけ使う。
- 各項目は1論点にする。
- 名詞止めと文を混在させない。
- 実装、設計、運用、文書を同じ箇条書きに混ぜない。
- 2項目で済むなら過剰に分解しない。

## 短い例
悪い例:
`body phase candidates は Running cancel と late response rejection を扱っている。`

良い例:
対象: 本文翻訳段階の候補は、実行中の取り消しを扱う。
対象: 本文翻訳段階の候補は、遅れて返った応答の破棄を扱う。
注意: `Running` のような状態値は、出典にある場合だけ固定名として残す。

悪い例:
`validator output の result が failed で、retry policy と owner assignment が unclear である。`

良い例:
結果: 検証結果の `result` は失敗である。
原因: 再実行する条件と担当者が決まっていない。
対応: 再実行条件と担当者を決めてから、検証をやり直す。

## レビューコメント
各指摘は次の順で書く。
1. 結論
2. 問題
3. 理由
4. 修正方針
5. 影響範囲

### レビューの追加ルール
- 「よくない」「微妙」だけで終わらせない。
- 可能なら修正例を1つ示す。
- 重大度を `critical | major | minor | nit` で付ける。
- コード、設計、ドキュメントの指摘を混在させない。
- 賛成と反論を求められた場合は両方書き、各主張に確証率を付ける。

## 設計説明
次の順で書く。
1. 目的
2. 変更前提
3. 変わるもの
4. 変わらないもの
5. 依存影響
6. 更新対象の文書

### 設計の追加ルール
- ER変更、層構造変更、依存ルール変更、文書の説明単位変更は分けて書く。
- 読者の認知負荷の話は、設計変更の列挙に混ぜず最後に書く。
- 差分の説明と最終形の説明を混同しない。
- 文章だけで追えない依存関係は、図が必要か明記する。

## 禁止
- 接続詞なしで結論へ飛ぶ。
- 抽象度の違う項目を同列に並べる。
- 同じ語で連続文を始める。
- 毎段落の末尾を総括文で締める。
- 根拠のない強い言い切り。
- 不要な横文字、不要な比喩、不要な感情語。
- 固定名の羅列だけで状態、差分、判断を説明する。
