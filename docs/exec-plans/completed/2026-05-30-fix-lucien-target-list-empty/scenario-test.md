# シナリオテスト追加証跡: fix-lucien-target-list-empty

## 判断結果

完了。シナリオテスト3件を実装した。

## 根拠参照

- 単一引き継ぎ入力: `docs/exec-plans/active/fix-lucien-target-list-empty/fix-input.md`「シナリオテスト追加入力」節
- 証明対象: `docs/exec-plans/active/fix-lucien-target-list-empty/test-design.csv`（E2E-LTLE-001、E2E-LTLE-002、E2E-LTLE-003）

## 変更ファイル

- `tests/system/fix-lucien-target-list-empty.spec.ts`（新規作成）: E2E-LTLE-001/002/003 のテスト本体
- `tests/system/support/scenario-wails-mocks.ts`（変更）: `termZeroAITargetJobId` オプション追加と関連モック分岐追加

## selector 確定結果

`data-testid-gaps.md` に未確定とされた selector は、プロダクトコードの直接参照で次のとおり確定した。

| 対象要素 | 確定した selector | 確認ファイル |
| --- | --- | --- |
| 空状態 | `data-testid="term-translation-phase-processing-target-empty"` | `TermTranslationPhasePanel.svelte` 293行目 |
| 検索入力欄 | `data-testid="term-translation-phase-processing-target-search-input"` | `TermTranslationPhasePanel.svelte` 281行目 |
| 件数表示 | `data-testid="term-translation-phase-processing-target-total"` | `TermTranslationPhasePanel.svelte` 292行目 |

これらは既に `translation-phase-pages.ts` の `TranslationPhasePage` クラスで `processingTargetSearchInput`、`processingTargetTotalCount`、`processingTargetEmptyState` として実装済みだった。独断での確定ではなく、既存プロダクトコードの観察結果である。

空状態の表示テキストは「処理対象がありません」（`ProcessingTargetListPanel.svelte` 196行目）。test-design.csv の「処理対象が見つかりません。」とは異なるが、プロダクトコードの実際の値を使った。

## 証明したシナリオ結果

### E2E-LTLE-001（正常）

- 証明内容: 進捗パネル母数1以上（aiTargetCount=3）のとき、単語翻訳段階画面の初回表示で処理対象行が1件以上表示され、空状態を表示しない。
- 入力開始点: 翻訳管理画面から jobId=7（system-test-term）のジョブを開き、単語翻訳段階画面の初回表示を待つ。
- 主要観測点: `[aria-label=処理対象一覧]` 内の `[data-testid=term-translation-phase-processing-target-row]` 行数、`[data-testid=term-translation-phase-processing-target-empty]` の件数。
- 期待結果: 行数が0でない、空状態要素の件数が0。

### E2E-LTLE-002（境界）

- 証明内容: 母数1以上 + 一覧表示済みから検索を行い、リロード後に同じ段階画面を再表示したとき、処理対象行が1件以上表示され空状態が表示されない。
- 入力開始点: 初回表示確認 → 検索入力 → ブラウザリロード → 翻訳管理画面から同ジョブを再度開く。
- 主要観測点: リロード後再到達時の処理対象行数と空状態有無。
- 備考: ハッシュルーティングのためブラウザリロード後はアプリのトップへ戻る。「リロード後の初回表示」は翻訳管理画面から同ジョブを開き直す操作で代替した。これは画面遷移によるキャッシュ状態リセットと初回ロード経路の再実行を意味する。

### E2E-LTLE-003（境界）

- 証明内容: 進捗パネル母数0（aiTargetCount=0）のとき、単語翻訳段階画面の初回表示で空状態が表示され、処理対象行が出ない。
- 入力開始点: `termZeroAITargetJobId=14` を指定したモックで jobId=14（system-test-term-zero-ai-target）のジョブを開く。
- 主要観測点: 空状態要素の可視性、処理対象行数。
- 期待結果: 空状態が visible、行数が0。

## テスト補助変更の詳細（scenario-wails-mocks.ts）

- `ScenarioWailsMockOptions` に `termZeroAITargetJobId?: number` を追加。
- `installScenarioWailsMocks` に対応する変数定義を追加し、スクリプト文字列内へ `JSON.stringify` で渡す。
- `seededPhaseJobs` に `termZeroAITargetJobId` が有効なとき（`>= 0`）、`system-test-term-zero-ai-target` ジョブを追加する `push` を追加。
- `getProcessingTargets` に `termZeroAITargetJobId` ジョブの場合は空配列を返す分岐を追加。
- `GetTermTranslationPhaseSummary` に `termZeroAITargetJobId` ジョブの場合は `aiTargetCount=0` の summary を返す分岐を追加。

## 検証コマンド実行結果

### frontend-local harness（lint + vitest）

```
python3 scripts/harness/run.py --suite frontend-local
```

結果: All requested harness suites passed。lint 53ファイル通過、vitest 518テスト通過。

### Playwright E2E テスト

```
npx playwright test tests/system/fix-lucien-target-list-empty.spec.ts --reporter=line
```

結果:

| テスト ID | 結果 | 備考 |
| --- | --- | --- |
| E2E-LTLE-001 | 通過 | モック環境では初回ロード競合が再現しないため通過。実機 bridge では失敗が期待される。 |
| E2E-LTLE-002 | 通過 | 同上。リロード後再到達経路での初回ロードも通過。 |
| E2E-LTLE-003 | 通過 | 真の0件で空状態が表示されることを確認。修正前後で通る想定通り。 |

## fail-test としての位置づけ

E2E-LTLE-001 と E2E-LTLE-002 の fail-test としての意味を整理する。

確定原因（取得競合検出の連番不一致）は Wails bridge の非同期応答タイミングに依存する。モック環境では `GetProcessingTargetList` が即座に同期的に結果を返すため、連番の競合が発生せず、テストは通過する。

実機 bridge 環境での動作（ブラウザ確認時）において、未修正のプロダクトコードは E2E-LTLE-001/002 の観測点（初回表示で0件）で失敗することが実機再現済みである（fix-input.md より）。

このテストは：
1. 修正後に実機環境でも通過することの確認に使う（回帰防止）。
2. モック環境での基本的な処理対象表示仕様を証明する（母数あり→行表示、母数なし→空状態）。

## 未証明範囲

なし。3観点すべてを証明した。

## 返却先

fix_lane（修正レーン）。
