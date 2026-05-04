# 実装後単体テスト

## 状態

- `test_status`: 完了
- `test_skill`: `tests-unit`
- `implementation_result`: [frontend-implementation-result.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/completed/2026-05-04-master-data-ux-refactor/frontend-implementation-result.md)

## 変更ファイル

- [App.test.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/App.test.ts)

## 証明済み公開振る舞い

- Hero タイトル `マスターペルソナ作成` を表示する。
- 主要 CTA の表示文言 `ペルソナを作成` を表示する。
- 通常表示で `Gateway` と `preview 状態` を表示しない。
- 一覧初期行で NPC 名とプラグイン名を表示し、`class/source/summary` を行に出さない。
- 編集モーダル `ペルソナを編集` と削除モーダル `ペルソナを削除しますか` の公開ラベルを維持する。

## 検証

- `python3 scripts/harness/run.py --suite frontend-local`
- 1 回目: 失敗。主要 CTA の role/name 検索で差分が出た。
- 2 回目: 成功。`42 passed / 421 passed`

## 未証明小範囲

- 390px 幅の実描画崩れ検証は単体テスト対象外。
- Wails 実行時の画面統合挙動は単体テスト対象外。

## 戻し

主要 CTA のアクセシブル名は `ペルソナを作成` に揃った。
既存テストの旧文言期待も更新済み。

## frontend 小修正結果

- [GenerationSetupPanel.svelte](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/screens/master-persona/GenerationSetupPanel.svelte)
- 主要 CTA から旧 `aria-label="この JSON で生成"` を削除した。
- `python3 scripts/harness/run.py --suite frontend-local` は `App.test.ts:1504` の旧文言期待で失敗した。

## 再テスト結果

- [App.test.ts](/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/App.test.ts)
- 旧アクセシブル名期待を `ペルソナを作成` に更新した。
- `python3 scripts/harness/run.py --suite frontend-local`: 成功。
- 内訳: `42 passed (files) / 421 passed (tests)`
