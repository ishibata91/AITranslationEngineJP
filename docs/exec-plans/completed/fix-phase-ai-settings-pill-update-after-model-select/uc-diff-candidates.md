# UC 差分候補: fix-phase-ai-settings-pill-update-after-model-select

## 判断サマリ

- 判定: `記述不足` を含む。`新規判断必要` は含まない。
- 停止要否: 不要。investigation-module は継続できる。

## 参照した正本

- ユースケース正本: `docs/usecases/uc-translation-management.md`
- 画面設計正本: `docs/screen-design/screens/term-translation-phase.md`、`persona-generation-phase.md`、`body-translation-phase.md`
- 修正方針: `docs/exec-plans/active/fix-phase-ai-settings-pill-update-after-model-select/fix-decision.md`

## UC 差分一覧

| No | 対象 UC | 分類 | 不足箇所 | 内容 | 判断根拠 |
| --- | --- | --- | --- | --- | --- |
| 1 | 翻訳段階を開始する（`uc-translation-management.md`） | 記述不足 | 主シナリオ前段フローの不足 | AI サービスとモデルを選択して `saveAISettings` が成功した後、summary の `aiSettings` が更新され、AI 設定 pill が「設定未完了」から「固定済み」へ切り替わり、開始ボタンが有効化される段階的な状態遷移が、UC の主シナリオにも代替シナリオにも記述されていない。 | 画面設計 E-03「モデル変更」操作は「成功時: 選択中モデルを更新する」と記述するが、pill 状態変化と開始ボタン有効化までの連鎖を説明する UC の記述がない。本 task の確定原因（`phaseAISettingsRepository` 未注入）が修正されると、この遷移フローが初めて正常動作するため、UC としての境界を記述する必要がある。 |
| 2 | 翻訳段階を開始する（`uc-translation-management.md`） | 記述不足 | 例外フローの境界不足 | backend の `JobPhaseAISettingsRepository` が未設定（`phase ai settings repository is not configured`）エラーを返した場合に、frontend が pill を「設定未完了」固定のまま開始ボタンを無効にし続ける経路が、例外シナリオとして記述されていない。本 task の恒久修正後はこのエラーが発生しなくなるが、UC として「AI 設定保存が失敗した場合に pill 状態が変化しない」境界を明示するかを判断する必要がある。 | 既存の例外 E1「開始条件が不足する」は開始操作を拒否する経路を説明するが、`saveAISettings` 自体が失敗した時の前段エラーを説明しない。ただし修正後はこの経路は発生しない。境界 UC として記述するかは `記述不足` として扱い、`新規判断必要` ではない。 |
| 3 | 3 phase 共通（単語翻訳 / NPC ペルソナ生成 / 本文翻訳） | 差分なし | 3 phase 全てで AI 設定保存の正常フローが共通である | 同一の `AIModelSelectionCard` コンポーネントと同型の presenter を 3 phase が共有する。UC の「翻訳段階を開始する」は 3 phase 共通で 1 つの UC として記述されており、phase 別の UC 分割は不要。 | `uc-translation-management.md` 「翻訳段階を開始する」が 3 phase を統合して記述している。記述不足の内容（No.1、No.2）は 3 phase に共通して適用される。 |

## 「新規判断必要」が含まれない根拠

- No.1 の状態遷移フロー（モデル選択 → pill「固定済み」→ 開始ボタン有効化）は、画面設計書 E-03 と「翻訳段階を開始する」UC 条件 2（「現在段階の AI 設定が開始条件を満たしている」）が期待挙動を間接的に定義している。新規仕様判断は不要で、記述不足の補完で対応できる。
- No.2 のエラー経路は、本 task の恒久修正後に発生しなくなる。UC 正本への恒久追記を要するかは、`finalization-module` で `updating-docs` skill を呼ぶ際に判断する。現時点では `記述不足` と分類する。
