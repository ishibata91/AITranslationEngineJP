# Task Plan: prebuilt-translation-dictionary

## branch情報

- `execution_branch`: `codex/prebuilt-translation-dictionary`
- `target_branch`: `master`
- `source_commit`: `5e971fcc9ffda72a69b4e52737af94fc6d2063f1`

## やること

- **R-1 事前作成した辞書成果物を生成する**: repository直下の`dictionary/`で、中心DBの`master_term`にある公式strings由来の英日対と`category`を移植し、翻訳実行より前に再生成できる辞書成果物を作る。pluginごとの既訳またはAI仮訳を持つ`proper_noun`は移植しない。
- **R-2 MCPから移植した辞書を操作できるようにする**: Codexが移植した辞書項目の検索、取得、追加、修正を頻繁に行える最小のMCPの道具を用意する。原語、訳語、`category`の通常indexと、英語と日本語の部分一致に使う検索用indexを持たせる。分類、候補抽出、レビューの機能は後続の要求で同じMCPへ足す。

## やらないこと

- 翻訳対象の文章を辞書訳へ事前置換する処理を追加または改善すること。
- 文章の各部分から辞書候補を取得し、翻訳エージェントへ渡す処理を実装すること。
- 辞書候補を翻訳エージェントへ渡す指示を設計すること。
- 辞書成果物そのものをGitへ追加すること。
- Nexus Modsへのupload、公開、更新操作を行うこと。
- 辞書項目を原語だけで一意にすること。
- `openai-batch-dictionary-file-token-usage`で棄却した辞書全文のfile参照またはVector store検索を再評価すること。
- `mod-npc-name-derivation`が扱うNPC名以外の要求を本taskへ統合すること。
- `proper_noun`を辞書成果物へ移植すること。第三者Mod由来の訳を配布対象へ混ぜないため、pluginを問わずtable全体を対象外にする。
- Open English WordNetによる一般語判定、意味と品詞の分類、辞書項目を作る規則を策定すること。移植した実データをMCPで確認した後のtaskで扱う。
- 会話、書籍、その他の英語原文から未収録の候補を抽出すること。分類の規則を確定した後のtaskで扱う。
- AIレビュー、人間レビュー、配布対象の判定を実装すること。MCPによる基本操作を確定した後のtaskで扱う。
- 辞書成果物をNexus Modsへ配置できるarchiveへまとめること。レビューと配布可否を確定した後のtaskで扱う。
