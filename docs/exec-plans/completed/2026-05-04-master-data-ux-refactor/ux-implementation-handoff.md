# UX実装修正入力

## 状態

- `handoff_status`: 承認済み
- `implementation_skill`: `implement-frontend`
- `source_ui_contract`: [ui-design.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/ui-design.md)
- `human_review`: [human-ui-review.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/human-ui-review.md)

## 実装目的

マスターペルソナ生成画面を、承認済み UI改善契約に合わせて作り直す。
画面目的は、翻訳前の準備としてベースゲームや大型 Mod の NPC ペルソナを事前に作成し、作成後に一覧と詳細で確認できるようにすることである。

## 実装範囲

- `MasterPersonaPage.svelte` の表示構造、文言、配置、画面専用部品化を変更する。
- 既存 `AIModelSelectionCard.svelte` は変更せず、そのまま利用する。
- 一覧は大量件数を前提に、プラグイン名と NPC 名を中心にした細い行へ変更する。
- 詳細は識別情報、声、話し方、ペルソナ本文を中心に表示する。
- 編集と削除は、詳細区画の補助操作として扱う。

## 禁止変更範囲

- backend 実装を変更しない。
- Wails gateway contract を変更しない。
- `AIModelSelectionCard.svelte` を変更しない。
- docs 正本本文を変更しない。
- task-local UIプロトタイプ、mock data を product code へ参照させない。
- 生成前の見積もり時間、料金目安、生成対象サンプルを追加しない。
- 既存スキップ理由一覧、次に確認するペルソナ推薦、プロンプト内容表示を追加しない。

## 読むファイル

- [implement-frontend/SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-frontend/SKILL.md)
- [ui-design.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/ui-design.md)
- [human-ui-review.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/human-ui-review.md)
- [prototype/index.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/prototype/index.svelte)
- [MasterPersonaPage.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/MasterPersonaPage.svelte)
- [AIModelSelectionCard.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/components/AIModelSelectionCard.svelte)

## 確認観点

- 画面冒頭のタイトルが `マスターペルソナ作成` になっている。
- 主要 CTA が `ペルソナを作成` として表示される。
- `Gateway`、`preview 状態`、プロンプトテンプレート説明が通常表示に出ない。
- 一覧の初期表示がプラグイン名と NPC 名を中心に読める。
- モデル選択カードが既存部品由来のまま表示される。
- 390px 幅でモデルカード、JSON ファイル名、長い NPC 名、ペルソナ本文が横にはみ出さない。

## 期待成果物

- frontend 実装結果ファイルを作成する。
- 変更ファイル、確認コマンド、未確認理由を記録する。
- 実装中に公開契約変更が必要になった場合は停止する。
