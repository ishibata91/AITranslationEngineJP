# 単体テスト観点: fix-term-translation-model-settings-empty-fixed

- `status`: confirmed
- `source_plan`: `./plan.md`
- `detail_spec_ref`: `./detail-spec-diff.md`

## 前提

- 本文書は単体テスト観点の固定のみを扱う。テスト実装は `implementation-module` が担当する。
- 観点の根拠は `./detail-spec-diff.md` の各要件（`er-REQ-001`、`er-REQ-002`、`term-translation-phase-REQ-002`、`term-translation-phase-REQ-007`、`persona-generation-phase-REQ-002`、`body-translation-phase-REQ-002`）とする。
- 3 フェーズ共通設計（`JOB_PHASE_AI_SETTINGS`、upsert、phase_type 主キー）に由来する観点は「共通」として記載する。phase_type 別の差異がある場合だけ個別に記載する。

---

## 1. `JOB_PHASE_AI_SETTINGS` repository の upsert / get

根拠: `er-REQ-001`

### 共通観点

| # | 観点 | 分類 | 証明対象 |
| --- | --- | --- | --- |
| U-REPO-001 | 存在しない phase_type に対して upsert を呼ぶと record が新規作成される | 正常 | upsert が INSERT として動作する |
| U-REPO-002 | 既存の phase_type に対して upsert を呼ぶと record が更新される | 正常 | upsert が UPDATE として動作する（重複作成しない） |
| U-REPO-003 | upsert の対象列は ai_provider / model_name / execution_mode / batch_mode の 4 値に限定され、credential_ref を含まない | 正常 | upsert 対象外の列が変更されない |
| U-REPO-004 | 存在する phase_type に対して get を呼ぶと保存された 3 値（ai_provider / model_name / execution_mode）が返る | 正常 | get が正しい record を返す |
| U-REPO-005 | 存在しない phase_type に対して get を呼ぶと ErrNotFound または nil 相当が返る | 例外 | record 不在を空値 record ではなく不在として返す |
| U-REPO-006 | upsert が job_id 列を持たない（入力に job_id を受け取らない） | 境界 | ER 正本「主キーは phase_type のみ、ジョブ無関連」を満たす |
| U-REPO-007 | 同一 phase_type への複数回 upsert 後、最大 1 件しか存在しない | 境界 | 重複 record が作成されない |

### phase_type 別の境界確認

| # | 観点 | phase_type | 分類 |
| --- | --- | --- | --- |
| U-REPO-008 | word_translation / npc_persona_generation / text_translation の 3 種が独立して upsert・get できる | 全 3 種 | 正常 |
| U-REPO-009 | word_translation record の upsert が npc_persona_generation record に影響しない | 種別間分離 | 境界 |

---

## 2. SaveAISettings service の入力検証

根拠: `term-translation-phase-REQ-002`（SaveAISettings 入力構造）

| # | 観点 | 分類 | 証明対象 |
| --- | --- | --- | --- |
| U-SAVE-001 | 入力が phase_type と provider / model / executionMode / batchMode を含む場合、upsert が成立する | 正常 | 正常入力経路 |
| U-SAVE-002 | 入力に job_id が含まれない（job_id を受け取る引数が存在しない） | 境界 | 「入力から job_id を抜く」仕様の確認 |
| U-SAVE-003 | 入力に credential_ref が含まれない（credential_ref を受け取る引数が存在しない） | 境界 | credential は provider-settings 都度解決の責務であることの確認 |
| U-SAVE-004 | phase_type が空文字または未知値の場合、保存が拒否される | 例外 | 入力検証 |
| U-SAVE-005 | 同一 phase_type への再保存は上書き（upsert）として処理され、旧値を消さない追記にならない | 境界 | 「再保存は上書き」の確認 |

---

## 3. フェーズ開始時の転写ロジック

根拠: `er-REQ-002`、`term-translation-phase-REQ-002`（開始時の転写仕様）

### 共通観点

| # | 観点 | 分類 | 証明対象 |
| --- | --- | --- | --- |
| U-START-001 | 対応 phase_type の JOB_PHASE_AI_SETTINGS record が存在する場合、フェーズ開始時に record の 3 値（provider / model / executionMode）が JOB_PHASE_RUN の対応列へ転写される | 正常 | 転写ロジックの正常経路 |
| U-START-002 | 対応 phase_type の record が存在しない場合、フェーズ開始が拒否される | 例外 | AI 設定不足（record 不在）での開始拒否 |
| U-START-003 | record は存在するが provider-settings 側の credential 解決が失敗する場合、フェーズ開始が拒否される | 例外 | credential 解決失敗での開始拒否 |
| U-START-004 | 転写時、JOB_PHASE_RUN.credential_ref には provider-settings から解決した認証参照値が設定され、JOB_PHASE_AI_SETTINGS の値は設定されない | 正常 | credential を Ready record ではなく provider-settings から都度解決する責務分担 |
| U-START-005 | JOB_PHASE_RUN はフェーズ開始成立時にだけ新規作成され、Ready 期には事前作成されない | 境界 | ER 正本「Ready 中は JOB_PHASE_RUN を事前作成しない」の確認 |

---

## 4. backend 応答 DTO の execution 有無分岐

根拠: `term-translation-phase-REQ-002`（backend 応答の構造）、同型で `persona-generation-phase-REQ-002`、`body-translation-phase-REQ-002`

### 共通観点

| # | 観点 | 分類 | 証明対象 |
| --- | --- | --- | --- |
| U-DTO-001 | JOB_PHASE_AI_SETTINGS record が存在し、かつ JOB_PHASE_RUN が未作成の場合、応答に execution field が含まれない（nil または field 不在） | 正常 | Ready 期未開始の応答構造 |
| U-DTO-002 | JOB_PHASE_AI_SETTINGS record が存在しない場合、応答の AI 設定 field が不在または nil になる（空文字値を持つ field で代理表現されない） | 例外 | AI 設定不在の応答構造。空文字代理表現の禁止確認 |
| U-DTO-003 | JOB_PHASE_RUN が存在する場合、応答に execution field が含まれ、JOB_PHASE_RUN の実行時固定値が反映される | 正常 | 実行中表示の応答構造 |
| U-DTO-004 | 応答に blocked_reason / ai_settings_insufficient_reason などの派生説明文字列が含まれない | 境界 | 「派生情報は frontend 責務、backend は含めない」の確認 |

---

## 5. presenter の execution 有無判定

根拠: `design-diff.md`（採用案 図 4、Presenter / UI の説明）

| # | 観点 | 分類 | 証明対象 |
| --- | --- | --- | --- |
| U-PRES-001 | 応答の execution field が存在する場合、modelLabel が execution.model_name の値を返す（空文字や "-" を返さない） | 正常 | 実行中表示の presenter 正常経路 |
| U-PRES-002 | 応答の execution field が存在しない（nil または不在）場合、isExecutionConfigured が false を返す | 正常 | 未設定判定の presenter 正常経路 |
| U-PRES-003 | 応答の execution field が存在しない場合、modelLabel が "-" または未設定を示す値を返す（空文字 "" を返さない） | 例外 | 空文字フォールバック未対応（確定原因 2）の回帰防止 |
| U-PRES-004 | 応答の execution field が存在しない場合、aiSettingsBlockedReason が AI 設定不足を示す理由文字列を返す（空文字 "" を返さない） | 例外 | 確定原因 3（空文字を未設定と判別できない）の回帰防止 |
| U-PRES-005 | 応答の AI 設定 field が不在の場合と、execution field が不在の場合を、presenter が独立して判定できる | 境界 | 2 つの不在状態の混同を防ぐ |

---

## 6. panel の状態 pill 分岐

根拠: `design-diff.md`（採用案 図 1、UI / Presenter の説明）、`docs/screen-design/screens/term-translation-phase.md`（E-03 状態別表示）

| # | 観点 | 分類 | 証明対象 |
| --- | --- | --- | --- |
| U-PANEL-001 | viewModel.isExecutionConfigured が true の場合、AI 設定固定状態の pill（「固定済み」相当）が表示される | 正常 | 実行中または設定済みの pill 表示 |
| U-PANEL-002 | viewModel.isExecutionConfigured が false の場合、AI 設定未完了の pill（「設定未完了」または「AI 設定不足」相当）が表示され、「固定済み」pill は表示されない | 例外 | 未設定状態の pill 表示分岐 |
| U-PANEL-003 | viewModel.aiSettingsBlockedReason が存在する場合、禁止理由が開始ボタン周辺に表示される | 例外 | 禁止理由の表示 |
| U-PANEL-004 | viewModel.isExecutionConfigured が false の場合、開始ボタンが disabled になる | 例外 | AI 設定不足での開始ボタン無効化 |

---

## 7. applyTermTranslationRuntimeSnapshot 経路の廃止確認

根拠: `design-diff.md`（採用案 図 4、OldErrPath 廃止）、`fix-decision.md`（確定原因、採用方針）

| # | 観点 | 分類 | 証明対象 |
| --- | --- | --- | --- |
| U-DEPRECATED-001 | applyTermTranslationRuntimeSnapshot が ErrNotFound を受け取った時、initial の ai_provider / model_name / execution_mode / credential_ref を空文字で上書きしない | 境界 | 空文字上書き経路（確定原因 1）の廃止確認 |
| U-DEPRECATED-002 | applyTermTranslationRuntimeSnapshot 自体が実行コンテキスト組み立てから除外されているか、または同関数が ErrNotFound 時に何も変更しない（無操作）で返る | 境界 | 廃止または無害化の確認 |
| U-DEPRECATED-003 | TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT への書き込みが SaveAISettings 経路から切断されている（JOB_PHASE_AI_SETTINGS への書き込みに移行している） | 境界 | 書き込み先移動の確認 |

---

## 3 フェーズ共通化の判断

| 判断 | 内容 |
| --- | --- |
| 共通化する単位 | U-REPO、U-SAVE、U-START、U-DTO の各観点は 3 フェーズ（word_translation / npc_persona_generation / text_translation）共通。phase_type を変数として同じテストロジックを表引きで適用できる。 |
| phase_type 別に分ける単位 | U-PRES、U-PANEL は frontend 側（presenter / panel コンポーネント）の単位。各フェーズが独立した presenter / panel ファイルを持つ場合は、phase_type ごとに同型の観点をそれぞれ確認する。単一の共通コンポーネントで実装される場合は 1 つの観点で代表できる。実装時に判断する。 |
| U-DEPRECATED は term-translation 専用 | applyTermTranslationRuntimeSnapshot は単語翻訳フェーズ固有の関数名であり、廃止確認は term-translation の実装に対して行う。persona-generation / body-translation の同型関数が存在する場合は同型の廃止確認観点を追加する。 |
