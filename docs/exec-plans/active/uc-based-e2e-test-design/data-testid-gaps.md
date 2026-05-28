# data-testid 不足一覧

## 目的

この資料は、`test-design.csv` を E2E 実装へ渡す前に固定が必要な selector を整理する。
対象は、現時点で親領域 selector と表示名に依存している操作、入力、行、状態表示である。

## 判断基準

- `data-testid` が固定済みの親領域は、現状のまま使う。
- 子要素の button、input、select、row、状態表示は、E2E 実装で直接指定する必要がある場合に不足として扱う。
- `aria-label` だけで定義されている翻訳フェーズ画面は、E2E 実装前に `data-testid` の固定方針を決める必要がある。

## ダッシュボード

| Gap ID | 対象 | 必要な固定 selector | 関連テスト |
| --- | --- | --- | --- |
| DTG-DASH-001 | 主要ページ card 内の `開く` button | `dashboard-primary-page-open-button` など | `E2E-UC-001` |
| DTG-DASH-002 | モバイルの主要ページ toggle | `dashboard-mobile-navigation-toggle` など | `E2E-UC-002` |
| DTG-DASH-003 | navigation item | `dashboard-global-navigation-item` など | `E2E-UC-002`, `E2E-UC-026` |

## AIサービス設定

| Gap ID | 対象 | 必要な固定 selector | 関連テスト |
| --- | --- | --- | --- |
| DTG-PROVIDER-001 | AIサービス row | `provider-settings-ai-service-row.<provider-id>` 相当 | `E2E-UC-003` から `E2E-UC-006` |
| DTG-PROVIDER-002 | APIキー設定 button | `provider-settings-api-key-open-button` など | `E2E-UC-003` |
| DTG-PROVIDER-003 | APIキー password input | `provider-settings-api-key-input` など | `E2E-UC-003` |
| DTG-PROVIDER-004 | APIキー保存 button | `provider-settings-api-key-save-button` など | `E2E-UC-003` |
| DTG-PROVIDER-005 | 接続確認 button | `provider-settings-connection-check-button` など | `E2E-UC-004`, `E2E-UC-027` |
| DTG-PROVIDER-006 | 設定保存 button | `provider-settings-save-button` など | `E2E-UC-005`, `E2E-UC-028` |
| DTG-PROVIDER-007 | リセット button | `provider-settings-reset-button` など | `E2E-UC-006` |
| DTG-PROVIDER-008 | 入力不正エラー表示 | `provider-settings-validation-error` など | `E2E-UC-028` |

## マスター辞書

| Gap ID | 対象 | 必要な固定 selector | 関連テスト |
| --- | --- | --- | --- |
| DTG-DICT-001 | 検索 input | `master-dictionary-search-input` など | `E2E-UC-007`, `E2E-UC-029` |
| DTG-DICT-002 | カテゴリ select | `master-dictionary-category-select` など | `E2E-UC-007` |
| DTG-DICT-003 | 辞書 row | `master-dictionary-entry-row` と entry id | `E2E-UC-007`, `E2E-UC-010`, `E2E-UC-011`, `E2E-UC-031` |
| DTG-DICT-004 | 新規登録 button | `master-dictionary-create-button` など | `E2E-UC-008`, `E2E-UC-009`, `E2E-UC-030` |
| DTG-DICT-005 | 作成・更新 modal 入力 | `master-dictionary-entry-source-input`, `master-dictionary-entry-category-select`, `master-dictionary-entry-origin-input`, `master-dictionary-entry-translation-input` など | `E2E-UC-008` から `E2E-UC-010`, `E2E-UC-030` |
| DTG-DICT-006 | 保存 button | `master-dictionary-entry-save-button` など | `E2E-UC-008` から `E2E-UC-010` |
| DTG-DICT-007 | 必須項目エラー表示 | `master-dictionary-entry-validation-error` など | `E2E-UC-009` |
| DTG-DICT-008 | 詳細 panel の更新 button | `master-dictionary-detail-edit-button` など | `E2E-UC-010` |
| DTG-DICT-009 | 詳細 panel の削除 button | `master-dictionary-detail-delete-button` など | `E2E-UC-011`, `E2E-UC-031` |
| DTG-DICT-010 | 削除確認 modal の確定・取消 button | `master-dictionary-delete-confirm-button`, `master-dictionary-delete-cancel-button` など | `E2E-UC-011`, `E2E-UC-031` |
| DTG-DICT-011 | XML file input | `master-dictionary-xml-file-input` など | `E2E-UC-012` |
| DTG-DICT-012 | XML 取り込み開始 button | `master-dictionary-xml-import-button` など | `E2E-UC-012`, `E2E-UC-032` |

## マスターペルソナ

| Gap ID | 対象 | 必要な固定 selector | 関連テスト |
| --- | --- | --- | --- |
| DTG-PERSONA-001 | JSON file input | `master-persona-json-file-input` など | `E2E-UC-013` |
| DTG-PERSONA-002 | AIサービス select | `master-persona-ai-service-select` など | `E2E-UC-013` |
| DTG-PERSONA-003 | モデル select | `master-persona-model-select` など | `E2E-UC-013` |
| DTG-PERSONA-004 | 実行方法 select | `master-persona-execution-mode-select` など | `E2E-UC-013` |
| DTG-PERSONA-005 | ペルソナ作成 button | `master-persona-generate-button` など | `E2E-UC-013`, `E2E-UC-033` |
| DTG-PERSONA-006 | ペルソナ検索 input | `master-persona-search-input` など | `E2E-UC-014`, `E2E-UC-034` |
| DTG-PERSONA-007 | プラグイン select | `master-persona-plugin-select` など | `E2E-UC-014` |
| DTG-PERSONA-008 | ペルソナ row | `master-persona-row` と persona id | `E2E-UC-014` から `E2E-UC-016`, `E2E-UC-035` |
| DTG-PERSONA-009 | 詳細 panel の編集・削除 button | `master-persona-edit-button`, `master-persona-delete-button` など | `E2E-UC-015`, `E2E-UC-016`, `E2E-UC-035` |
| DTG-PERSONA-010 | 編集 modal 入力 | `master-persona-summary-input`, `master-persona-speech-style-input`, `master-persona-body-input` など | `E2E-UC-015`, `E2E-UC-035` |
| DTG-PERSONA-011 | 編集保存・キャンセル button | `master-persona-edit-save-button`, `master-persona-edit-cancel-button` など | `E2E-UC-015`, `E2E-UC-035` |
| DTG-PERSONA-012 | 削除確定 button | `master-persona-delete-confirm-button` など | `E2E-UC-016` |

## 未完了ジョブ一覧

| Gap ID | 対象 | 必要な固定 selector | 関連テスト |
| --- | --- | --- | --- |
| DTG-JOB-001 | 状態フィルタ option | `translation-job-management-state-filter-option` など | `E2E-UC-017` |
| DTG-JOB-002 | job card の識別 selector | `translation-job-management-job-card.<job-id>` 相当 | `E2E-UC-017` から `E2E-UC-020`, `E2E-UC-036` から `E2E-UC-039` |
| DTG-JOB-003 | 停止 button | `translation-job-management-stop-button` など | `E2E-UC-018`, `E2E-UC-037` |
| DTG-JOB-004 | 再開 button | `translation-job-management-resume-button` など | `E2E-UC-019`, `E2E-UC-038` |
| DTG-JOB-005 | 削除 button | `translation-job-management-delete-button` など | `E2E-UC-020`, `E2E-UC-039` |
| DTG-JOB-006 | 削除確認 modal の確定・戻る button | `translation-job-management-delete-confirm-button`, `translation-job-management-delete-cancel-button` など | `E2E-UC-020`, `E2E-UC-039` |

## 翻訳入力レビュー

| Gap ID | 対象 | 必要な固定 selector | 関連テスト |
| --- | --- | --- | --- |
| DTG-INPUT-001 | JSON file input | `translation-input-review-json-file-input` など | `E2E-UC-021` |
| DTG-INPUT-002 | 登録または読み込み button | `translation-input-review-register-button` など | `E2E-UC-021`, `E2E-UC-040` |
| DTG-INPUT-003 | 登録失敗理由表示 | `translation-input-review-register-error` など | `E2E-UC-040` |

## 翻訳フェーズ画面

| Gap ID | 対象 | 必要な固定 selector | 関連テスト |
| --- | --- | --- | --- |
| DTG-PHASE-001 | phase root | `term-translation-phase-screen`, `persona-generation-phase-screen`, `body-translation-phase-screen` など | `E2E-UC-045` から `E2E-UC-047` |
| DTG-PHASE-002 | 状態概要 | `<phase>-status-panel` など | `E2E-UC-045` から `E2E-UC-047` |
| DTG-PHASE-003 | 開始 button | `<phase>-start-button` など | `E2E-UC-045` から `E2E-UC-047` |
| DTG-PHASE-004 | 進捗表示 | `<phase>-progress-bar`, `<phase>-progress-counts` など | `E2E-UC-045` から `E2E-UC-047` |
| DTG-PHASE-005 | AI モデル固定状態 | `<phase>-ai-model-lock-state` など | `E2E-UC-045` から `E2E-UC-047` |
| DTG-PHASE-006 | 処理対象 row | `<phase>-processing-target-row` と target id | `E2E-UC-045` から `E2E-UC-047` |
| DTG-PHASE-007 | 後続導線 | `body-translation-phase-complete-next-action` など | `E2E-UC-047` |
| DTG-PHASE-008 | 次の作業フッターの主要操作 button | `job-run-next-action-primary-button` など | `E2E-UC-048` から `E2E-UC-050` |
| DTG-PHASE-009 | 次へ進めない理由表示 | `job-run-next-action-blocked-reason` など | `E2E-UC-050` |

## 出力管理

| Gap ID | 対象 | 必要な固定 selector | 関連テスト |
| --- | --- | --- | --- |
| DTG-OUTPUT-001 | 出力候補 job row | `output-management-output-candidate-row` と job id | `E2E-UC-023` から `E2E-UC-025`, `E2E-UC-042` から `E2E-UC-044` |
| DTG-OUTPUT-002 | 対象ゲーム select | `output-management-target-game-select` など | `E2E-UC-023` |
| DTG-OUTPUT-003 | 出力先 path input | `output-management-output-path-input` など | `E2E-UC-023`, `E2E-UC-024`, `E2E-UC-042` |
| DTG-OUTPUT-004 | XML 出力 button | `output-management-export-button` など | `E2E-UC-023`, `E2E-UC-042` |
| DTG-OUTPUT-005 | 再出力 button | `output-management-reexport-button` など | `E2E-UC-024`, `E2E-UC-043` |
| DTG-OUTPUT-006 | 差分 row | `output-management-diff-row` と row id | `E2E-UC-025`, `E2E-UC-044` |
| DTG-OUTPUT-007 | 出力先 path エラー表示 | `output-management-output-path-error` など | `E2E-UC-042` |

## 優先度

最初に固定する対象は、E2E 実装で locator が不安定になりやすい row と button である。
優先順は次の通りである。

1. 操作 button の単体 selector。
2. list row の識別 selector。
3. form input と select の単体 selector。
4. エラー表示と無効理由の状態 selector。
5. 翻訳フェーズ画面の root、開始 button、進捗、処理対象 row。
