# Storybook レビューループ記録（スライス 4: 結果行に口調メタデータ）

## 方針（確定）

専用のキャラ管理画面・編集はやめ、既存の翻訳結果行を開いたときの詳細に、台詞の話者の生成済み基底口調を「口調」メタデータとして出す（表示のみ）。既存の到達可能な画面へ統合し、独立コンポーネントを宙に浮かせない。

## 確定した表示

- 判定結果を強調する: 基底口調セル（例: 物腰やわ）を大きく、性質文（口調を普通の言葉で説明した一文）を続ける。
- 根拠は小さく: 決定経路・対人段階・感情段階・印 を 1 行の薄い注記に畳む。経路の根拠は注記の tooltip で読める。

## 反映先

- コンポーネント: `frontend/src/ui/screens/translation-run/TranslationResultRow.svelte`（展開詳細に口調メタ節を追加）
- 型: `translation-run-view.ts`（`PersonaMeta`〔cell・trait・段階・印・経路〕・段階型・`NarrationResultRow.persona?` を追加）
- 表示定数: `translation-run-presentation.ts`（段階ラベル・決定経路ラベル・補足）
- story: `TranslationResultRow.stories.ts`（本文・声質・保留・叙述文・口調なし・畳む）
- fixture: 実データで確認した代表値（Inigo 本文物腰やわ・Nazeem 声質横柄・Galmar 保留）

## 承認と分類

- 人間レビュー承認: 済み（「UI はこれでいい」）。
- 分類: 作業中分類 `Review/Changed Components/TranslationResultRow` を廃し、`UI Components/TranslationResultRow`（通常分類）へ統合。レビュー story ファイルは削除。
- 検証: `npm --prefix frontend run build-storybook` 成功、`python3 scripts/harness/run.py --suite frontend-local`（lint・test）通過。

## 表示範囲外（implementation-module へ渡す）

- 口調メタの取得と結果行への合流: api binding が persona_character を読み、結果行 DTO（cell・trait・段階・印・経路）へ載せる。frontend gateway/controller で `NarrationResultRow.persona` へ写す。
- 性質文（trait）は基底口調セルから引く（backend の性質文カタログ）。
