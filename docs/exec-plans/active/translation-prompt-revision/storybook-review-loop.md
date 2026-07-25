# Storybook Review Loop: translation-prompt-revision

`storybook-review-loop.md` は確定した story、変更された画面仕様、反映先、現在分類、承認状態だけを持つ。人間コメントの履歴は残さない。

## 起動と確認先

- 起動 command: `npm --prefix frontend run storybook -- --port 6008`
- URL: `http://localhost:6008/`

## 対象 story

| story | 現在分類 | 見るもの | 承認状態 |
| --- | --- | --- | --- |
| DirectiveEditor / 指示文一覧（9 種） | `UI Components/DirectiveEditor` | 指示文が 9 種並ぶこと、各指示文の対象 REC:FIELD が分割後の割り当てで出ること | 承認済み（2026-07-26） |
| プロンプトテンプレート / ベースタブ | `Screens/プロンプトテンプレート` | base 指示文が 4 段落の新既定文で出ること | 承認済み（2026-07-26） |
| プロンプトテンプレート / レコード別タブ | `Screens/プロンプトテンプレート` | 9 種の指示文と対象一覧が画面内で崩れないこと | 承認済み（2026-07-26） |

## 変更された画面仕様

- 指示文の並びは design.md 層2 の表と同じ順（物品説明・効果説明・世界観断片・書物体・日記体・固有名・操作名・語義・口調）にする。
- 対象 REC:FIELD の `logical_name` は `db/migrations/0006_record_type_translation.sql` の seed 値に揃える。
- 表示構造は変えない。`DirectiveEditor.svelte` は指示文を `{#each}` で並べるため、9 種化は行数の増加として現れる。

## 反映先

- `frontend/src/ui/screens/template-editor/template-editor.fixtures.ts`
- `frontend/src/ui/screens/template-editor/DirectiveEditor.stories.ts`
- `frontend/src/ui/screens/template-editor/TemplateEditorScreen.stories.ts`

## 表示範囲外の残課題（implementation-module へ渡す）

- `PromptTemplateForm.personaTemplate` の削除。view 型・gateway 型・container の state と連動するため表示範囲で扱わない。
- 指示文の割り当ての正本は seed 側にあり、`db/migrations/0006_record_type_translation.sql` の書き直しと `internal/store/seed_test.go` の 9 種化は backend 側に残る。

## 合意済み frontend 保護

- 承認済み画面: プロンプトテンプレート画面（ベースタブ、レコード別タブ）と `DirectiveEditor`。
- 表示規則: 指示文は design.md 層2 の表の順で縦に並べる。対象 REC:FIELD の `logical_name` は seed の値に揃える。指示文ごとに「本文 textarea → 差し込める変数（あれば）→ 対象一覧」の順で出す。
- 変更禁止範囲: `DirectiveEditor.svelte` と `TemplateEditorScreen.svelte`・`TemplateBasePane.svelte` の表示構造、props の形、style。後続の実装で表示を変えない。
- 通常分類へ戻した story: `UI Components/DirectiveEditor`、`Screens/プロンプトテンプレート`。
