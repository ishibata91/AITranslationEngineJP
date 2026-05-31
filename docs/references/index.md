# References Index

関連文書: [`../index.md`](../index.md), [`../tech-selection.md`](../tech-selection.md), [`../coding-guidelines.md`](../coding-guidelines.md)

このディレクトリは、外部仕様、ベンダー資料、複数 skill が参照する運用規約の参照方針をまとめる。

## Rules

- 新しい参照資料は、原則として `docs/references/` 配下に追加する
- library や framework の書き方を docs に反映する前に、`npx ctx7 library` / `npx ctx7 docs` で official docs を確認する
- vendor dump を置く場合でも、用途、出典、参照先文書を短く説明する
- implementation を拘束する外部仕様は、どの source-of-truth doc から参照するかを明示する

## Current Sources

- Wails official docs:
  [`Getting Started`](https://wails.io/docs/gettingstarted/firstproject),
  [`Application Development`](https://wails.io/docs/guides/application-development),
  [`Project Config`](https://wails.io/docs/reference/project-config)
- Svelte official docs:
  [`Overview`](https://svelte.dev/docs/svelte/overview),
  [`TypeScript`](https://svelte.dev/docs/svelte/typescript),
  [`Migration Guide`](https://svelte.dev/docs/svelte/v5-migration-guide)
- Vite official docs:
  [`Guide`](https://vite.dev/guide/),
  [`Build`](https://vite.dev/guide/build),
  [`Config`](https://vite.dev/config/)
- [`./xtranslator_ref.md`](./xtranslator_ref.md): xTranslator の入出力形式整理
- [`./term-translation-target-record-candidates.md`](./term-translation-target-record-candidates.md): 単語翻訳フェーズの単語対象レコード候補と抽出有無の整理
- [`./agent-browser.md`](./agent-browser.md): サブエージェントの UI 証跡取得用 `agent-browser` CLI 利用規約
- [`./storybook.md`](./storybook.md): Storybook の作成、起動、分類、確認資源、`fixture` 種類基準の共通規約
- [`./vendor-api/README.md`](./vendor-api/README.md): vendor API の生参照とダンプ置き場
- `../everything-claude-code/rules/`: コーディング規約分割時に採用技術へ合う項目だけを輸入した参照元

## Migration Note

この repo の desktop 基盤は `Wails + Go + Svelte` を正本とする。
旧実装由来の参照は、新しい source of truth を補強しない限り採用しない。
