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
4. 編集前に gateguard の事実確認を行う

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

## 文のルール
- 抽象語より具体語を使う。
- 「これ」「それ」「今回」「この場合」は参照先が曖昧なら使わない。
- 因果は明示する。少なくとも「理由」「影響」「対応」のどれかを入れる。
- 同列項目は同じ抽象度・同じ文型で並べる。
- 読点が3つを超える文は分割を検討する。
- 修飾語は被修飾語の近くに置く。
- 強い断定には根拠を添える。
- 推測は「可能性がある」と明示する。

## 箇条書き
- 真に列挙である場合だけ使う。
- 各項目は1論点にする。
- 名詞止めと文を混在させない。
- 実装、設計、運用、文書を同じ箇条書きに混ぜない。
- 2項目で済むなら過剰に分解しない。

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
