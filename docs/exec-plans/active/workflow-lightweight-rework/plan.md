# plan: workflow-lightweight-rework

## 依頼要約

workflow 全体を本質から問い直し、軽量化する。

直前 task `boundary-integration-test-rollout`（rollout 親 task、active で design-module 出口完了済み）の検討中に、user が次を提示:

1. この製品（Skyrim Mod 翻訳エンジン、desktop、user 1 人）では E2E 全網羅で十分。
2. 「frontend / backend の意思ずれ検出」は組織的事情（agent 分割）に起因し、技術 layer で守るのは対症療法。
3. workflow 全体（特に設計と harness、agent 分割）を本質から問い直す。

議論結果として次が確定:

- 設計は task 性格でスケールする（軽 / 重 2 段階、I1）
- harness は触った範囲だけ動かす（B 維持）
- 実装 agent 分割は廃止。Claude 本体が文脈を持ったまま縦通しで実装する（E1）
- skill は「着替える先」として残す
- 単体テスト守備範囲は限定（key 保存、プロンプト構築、純粋 class、ビジネスロジック層、バグ再発防止のみ）
- `docs/usecases/`、`docs/detail-specs/` は廃止（agent 分割消滅で認識合わせ docs が不要になる）
- 中核 docs: 画面設計 + `docs/architecture.md`
- module skill 構成は現行 5 維持、中身軽量化（H3）
- 実装 / test / 観測 skill は 1 skill `implement` に統合（F1）

## 作業 branch

- 作業 branch: `claude/workflow-lightweight-rework`
- 分岐元 branch: `master`
- 分岐元 commit: `946946d6`

## 直前 task との関係

直前 task の active plan folder（`docs/exec-plans/active/boundary-integration-test-rollout/`）は untracked のまま破棄した（design-module 出口で正本反映前に方針転換）。

直前 task の master 上 pilot 実装（`master-dictionary-boundary-test`）と種別仕様（`docs/detail-specs/boundary-integration-test.md`）は本 task で削除する。

## 検証戦略

テスト追加なし、新規 production code 追加なし。検証は次の確認のみ:

- `npm run lint:frontend`
- `npm run lint:backend`
- `npm run test:frontend`
- `npm run test:backend`
- `docs/architecture.md` の手動確認
- `.claude/skills/implement/SKILL.md` の手動確認

## 検証結果（2026-06-06）

- `python3 scripts/harness/run.py --suite backend-local`: PASS（lint:backend、test:backend、coverage:backend すべて通過）
- `python3 scripts/harness/run.py --suite frontend-local`: PASS（lint:frontend、test:frontend 54 file / 636 test 通過）
- `python3 scripts/harness/run.py --suite structure`: PASS（docs index 整合通過）

## 実施内容（2026-06-06）

### Phase A: rollout 親 task 破棄

- `docs/exec-plans/active/boundary-integration-test-rollout/` を folder ごと削除（untracked のまま）

### Phase B: master 上 pilot 実装削除

- `internal/apitest/master_dictionary_boundary_test.go` 削除
- `internal/apitest/boundary_golden_loader.go` 削除
- `internal/apitest/testdata/boundary/master_dictionary/` 配下 11 file 削除
- `internal/apitest/testdata/` 空 directory 削除
- `frontend/src/controller/wails/master-dictionary.boundary.test.ts` 削除
- `frontend/src/controller/wails/boundary-golden-loader.ts` 削除
- `frontend/package.json` の knip ignore 1 行削除

### Phase C: docs 階層整理

- `docs/usecases/` 配下 7 file 削除
- `docs/detail-specs/` 配下 10 file 削除（`boundary-integration-test.md` を含む）
- `docs/index.md` から detail-specs / usecases 参照を削除
- `docs/references/term-translation-target-record-candidates.md` の detail-specs 参照を architecture.md に変更

業務ルール（13 種別の REC、`IsTermTarget`）は code（`internal/recclassification/`）に encode 済みのため救出統合不要と判断。

### Phase D: module skill 中身軽量化（H3）

- `.claude/skills/preparation-module/SKILL.md`: 入口に軽 / 重判定 step を追加
- `.claude/skills/design-module/SKILL.md`: agent 起動廃止、Claude 本体が直接書く形式に書き換え、detail-spec-diff 廃止、単体テスト守備範囲（論点 6）を明示
- `.claude/skills/implementation-module/SKILL.md`: 層別 implementer 起動を廃止、Claude 本体が `implement` skill を読んで縦通しで実行する形式に書き換え
- `.claude/skills/storybook-module/SKILL.md`: `frontend_implementer` / `designer` 起動を廃止、Claude 本体直接実装に変更
- `.claude/skills/finalization-module/SKILL.md`: 詳細仕様正本反映を廃止、正本反映対象を `docs/architecture.md` と `docs/screen-design/` に限定

### Phase E: skill 統合（F1）

- `.claude/skills/implement/SKILL.md` 新規作成（backend + frontend ロジック + 統合境界 + test + 観測ログを統合した「着替える先」skill）
- `.claude/skills/implement-frontend/`、`implement-backend/`、`implement-integration/`、`tests-scenario/`、`tests-unit/`、`observability-implementer/`、`detail-spec-design/` の 7 skill folder を削除

### Phase F: agent 整理

- `.claude/agents/frontend_implementer.md`、`backend_implementer.md`、`integration_implementer.md`、`implementation_tester.md` の 4 agent 削除
- `designer`、`test_designer`、`conflict_resolver`、`fix_decider` は維持（明示的な user 起動時のみ使う）

### Phase G: CLAUDE.md / メモリ反映

- `CLAUDE.md`: 「実装フローの原則」section を追加（Claude 本体が 1 文脈で書く原則）
- memory `feedback-implementer-no-agent-split` 新規保存
- memory `feedback-boundary-responsibility-separation` を本見直し後の責務分離（形式 = code、意味と表示 = 画面設計）に書き換え
- `MEMORY.md` index を更新
