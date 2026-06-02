# E2E テスト観点差分: fix-phase-ai-settings-pill-update-after-model-select

## 判断サマリ

- 判定: `追加候補あり`
- 停止要否: 不要。

## 参照した正本

- E2E テスト観点正本: `docs/e2e-test-design/test-design.csv`
- E2E テスト規約: `docs/e2e-test-guidelines.md`
- 画面設計正本: `docs/screen-design/screens/term-translation-phase.md`、`persona-generation-phase.md`、`body-translation-phase.md`
- 修正方針: `docs/exec-plans/active/fix-phase-ai-settings-pill-update-after-model-select/fix-decision.md`

## 既存観点との照合結果

### 単語翻訳 / NPC ペルソナ生成 / 本文翻訳 の「翻訳段階を開始する」観点

| 既存 ID | 対象画面 | 前提条件の AI 設定状態 | モデル選択 → pill 変化の操作を含むか |
| --- | --- | --- | --- |
| `E2E-UC-045` | 単語翻訳 | `AI 設定準備済み` が既に表示されている前提 | 含まない |
| `E2E-UC-046` | NPC ペルソナ生成段階 | `AI 設定準備済み` が既に表示されている前提 | 含まない |
| `E2E-UC-047` | 本文翻訳段階 | `AI 設定準備済み` が既に表示されている前提 | 含まない |

判断: 3 観点はいずれも AI 設定が保存済みの状態からテストが始まる。AI サービスとモデルを選択して保存し、pill が「設定未完了」から「固定済み」へ切り替わり、開始ボタンが有効化されるまでの導線を証明する観点が存在しない。

### 単体テスト観点の扱い

依頼内容で「単体テスト観点（bootstrap repository 注入の DI テスト、presenter `buildModelOptions` の placeholder 重複解消テスト）は test-design skill の対象外なら範囲外と明示する」と指定されている。E2E テスト規約は UI 人間操作 E2E を対象とし、単体テストは `coding-guidelines-tests.md` を正本にする。従って、次の単体テスト観点は本成果物の範囲外とする。

- `internal/bootstrap/app_controller.go` の `WithTermTranslationJobPhaseAISettings` 注入を証明する DI テスト
- presenter `buildModelOptions` の `value=""` placeholder 重複解消を証明するプレゼンターテスト

これらは `implementation-module` の `tests-unit` skill で扱う。

## E2E テスト観点差分一覧

| No | 分類 | 関連 UC | 対象画面 | 不足観点 | 追加候補 ID |
| --- | --- | --- | --- | --- | --- |
| 1 | 追加候補あり | 翻訳段階を開始する | 単語翻訳 | AI サービス選択 → モデル一覧更新 → モデル選択 → `saveAISettings` 成功 → pill「固定済み」→ 開始ボタン有効化の正常導線を証明する観点がない | `E2E-UC-056` |
| 2 | 追加候補あり | 翻訳段階を開始する | NPC ペルソナ生成段階 | 同型の正常導線（単語翻訳と同じ操作軸）を証明する観点がない | `E2E-UC-057` |
| 3 | 追加候補あり | 翻訳段階を開始する | 本文翻訳段階 | 同型の正常導線（単語翻訳と同じ操作軸）を証明する観点がない | `E2E-UC-058` |

## 追加候補テスト観点（CSV 形式）

既存 test-design.csv の末尾に追加することを想定する。

```csv
E2E-UC-056,翻訳段階を開始する,単語翻訳,"画面表示: [data-testid=term-translation-phase-screen] が表示されている。画面表示: [data-testid=term-translation-phase-ai-model-lock-state] が 設定未完了 を表示している。画面表示: [data-testid=term-translation-phase-start-button] が disabled 状態で表示されている。","[aria-label=AIサービス] を select して Gemini を選択する; [aria-label=モデル一覧を更新] を click する; モデル一覧が更新されるまで待つ; モデル select で gemini-test を select する","[data-testid=term-translation-phase-ai-model-lock-state] が 固定済み を表示する; [data-testid=term-translation-phase-start-button] が enabled になる","正常: AI 設定保存後に pill と開始ボタンが切り替わる。AI モデル固定状態 selector は画面設計 E2E 固定 selector から取得"
E2E-UC-057,翻訳段階を開始する,NPC ペルソナ生成段階,"画面表示: [data-testid=persona-generation-phase-screen] が表示されている。画面表示: [data-testid=persona-generation-phase-ai-model-lock-state] が 設定未完了 を表示している。画面表示: [data-testid=persona-generation-phase-start-button] が disabled 状態で表示されている。","[aria-label=AIサービス] を select して Gemini を選択する; [aria-label=モデル一覧を更新] を click する; モデル一覧が更新されるまで待つ; モデル select で gemini-test を select する","[data-testid=persona-generation-phase-ai-model-lock-state] が 固定済み を表示する; [data-testid=persona-generation-phase-start-button] が enabled になる","正常: NPC ペルソナ生成段階で同型の AI 設定保存後 pill 変化を証明する。3 phase の観点網羅を担保する"
E2E-UC-058,翻訳段階を開始する,本文翻訳段階,"画面表示: [data-testid=body-translation-phase-screen] が表示されている。画面表示: [data-testid=body-translation-phase-ai-model-lock-state] が 設定未完了 を表示している。画面表示: [data-testid=body-translation-phase-start-button] が disabled 状態で表示されている。","[aria-label=AIサービス] を select して Gemini を選択する; [aria-label=モデル一覧を更新] を click する; モデル一覧が更新されるまで待つ; モデル select で gemini-test を select する","[data-testid=body-translation-phase-ai-model-lock-state] が 固定済み を表示する; [data-testid=body-translation-phase-start-button] が enabled になる","正常: 本文翻訳段階で同型の AI 設定保存後 pill 変化を証明する。3 phase の観点網羅を担保する"
```

## selector 確認事項

追加候補観点が参照する `data-testid` は、各 phase の画面設計書「E2E 固定 selector」表で定義済みである。

| selector | 定義場所 | 状態 |
| --- | --- | --- |
| `term-translation-phase-ai-model-lock-state` | `term-translation-phase.md` E2E 固定 selector | 定義済み |
| `term-translation-phase-start-button` | `term-translation-phase.md` E2E 固定 selector | 定義済み |
| `persona-generation-phase-ai-model-lock-state` | `persona-generation-phase.md` E2E 固定 selector | 定義済み |
| `persona-generation-phase-start-button` | `persona-generation-phase.md` E2E 固定 selector | 定義済み |
| `body-translation-phase-ai-model-lock-state` | `body-translation-phase.md` E2E 固定 selector | 定義済み |
| `body-translation-phase-start-button` | `body-translation-phase.md` E2E 固定 selector | 定義済み |

3 phase 全ての `ai-model-lock-state` と `start-button` selector は定義済みのため、`data-testid-gaps.md` に記録する不足 selector はない。

## 3 phase 網羅性の確認

- 単語翻訳: `E2E-UC-056` で追加候補あり
- NPC ペルソナ生成: `E2E-UC-057` で追加候補あり
- 本文翻訳: `E2E-UC-058` で追加候補あり

既存の観点（`E2E-UC-045` / `E2E-UC-046` / `E2E-UC-047`）は「AI 設定準備済み前提で開始を実行する」観点であり、追加候補（`E2E-UC-056` / `E2E-UC-057` / `E2E-UC-058`）は「未設定状態からモデルを選択して AI 設定保存 → pill 変化を証明する」観点であるため、重複しない。

注記: 承認時の ID は `E2E-UC-053/054/055` だったが、`translation-phases.spec.ts`（`E2E-UC-053`）および `master-dictionary-management.spec.ts`（`E2E-UC-054`）に既存 ID が存在していたため、`E2E-UC-056/057/058` に採番し直した。`test-design.csv` も同様に更新済み。
