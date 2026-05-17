# Scenario Candidates: 2026-05-18-storybook-foundation / failure

- `generator`: `failure`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `SBF`

## Generator Scope

- `viewpoint`: 失敗観点
- `included_sources`: `plan.md`, `docs/tech-selection.md`, `docs/coding-guidelines-frontend.md`, `docs/lint-policy.md`, `frontend/package.json`, `frontend/vite.config.ts`
- `excluded_sources`: Master Persona の部品化、画面再設計、gateway mock、backend DTO mock、実行フロー再現
- `generation_notes`: Storybook 最小基盤の起動、build、fixture 配置、review URL 記録、scope 逸脱だけを候補化する。候補の採否、統合、最終シナリオ表は `designer` が扱う。

## Candidate Scenarios

### CAND-SBF-001 Storybook dev 入口が起動できない

- `source requirement`: `plan.md` の goal と close_conditions。Storybook dev 入口が動くこと。
- `viewpoint`: 参照不能、設定不整合
- `candidate scenario id`: `CAND-SBF-001`
- `actor`: 実装者
- `trigger`: `frontend/package.json` に Storybook dev script がない、または script が `.storybook/` 設定を参照できない。
- `rejected operation`: Storybook dev server を起動して review 用 URL を得る操作。
- `expected error`: command failure として終了し、起動 URL を記録できない。
- `observable point`: `npm --prefix frontend run storybook` 相当の実行結果、process exit code、起動 URL の有無。
- `related detail requirement type`: `failure_handling_requirement`, `testability_requirement`
- `adoption hint`: 起動入口そのものを close condition に含める場合に採用候補になる。
- `conflict hint`: dev 起動を検証対象外にして build だけを必須にする正常系候補と競合する可能性がある。

### CAND-SBF-002 Storybook build が失敗する

- `source requirement`: `plan.md` の validation_commands。`npm --prefix frontend run build-storybook` が検証入口として挙げられている。
- `viewpoint`: 設定不整合、保存失敗
- `candidate scenario id`: `CAND-SBF-002`
- `actor`: 実装者
- `trigger`: Storybook build script、builder 設定、Vite 設定、Svelte 設定のいずれかが不整合になる。
- `rejected operation`: 静的 Storybook を build して基盤確認を完了する操作。
- `expected error`: build command が失敗し、Storybook 静的出力が作成されない。
- `observable point`: `npm --prefix frontend run build-storybook` の実行結果、出力 directory、frontend-local との関係記録。
- `related detail requirement type`: `failure_handling_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: build-storybook を merge 前 gate に含める場合に採用候補になる。
- `conflict hint`: Storybook build を `frontend-local` に含めるか別 gate にするかが未決のため、検証段階が他候補と競合する可能性がある。

### CAND-SBF-003 Vite alias または Svelte plugin を Storybook が解決できない

- `source requirement`: `frontend/vite.config.ts` の alias と plugin。Storybook 関連設定が対象であるという `plan.md` の影響範囲。
- `viewpoint`: 参照不能、設定不整合
- `candidate scenario id`: `CAND-SBF-003`
- `actor`: 実装者
- `trigger`: Storybook 側の設定が `@ui`、`@application`、`@controller` の alias または Svelte plugin 設定を取り込めない。
- `rejected operation`: 最小 story を Storybook 上で render する操作。
- `expected error`: module resolution failure または Svelte compile failure が発生する。
- `observable point`: Storybook dev / build の error log、対象 import path、最小 story の render 結果。
- `related detail requirement type`: `failure_handling_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: Storybook config と既存 Vite config の整合を証明対象に入れる場合に採用候補になる。
- `conflict hint`: Storybook 専用 config を最小にする方針と、既存 Vite config を再利用する方針の境界で競合する可能性がある。

### CAND-SBF-004 fixture 配置が決まらず後続 story が参照できない

- `source requirement`: `plan.md` の goal と未決事項。fixture 配置方針を固定すること。
- `viewpoint`: 参照不能、人間判断候補
- `candidate scenario id`: `CAND-SBF-004`
- `actor`: 後続 Master Persona task の実装者
- `trigger`: fixture を `frontend/src/ui/**/__fixtures__` に置くか、Storybook 専用 directory に置くかが未確定のまま基盤 task が完了する。
- `rejected operation`: 後続 story が基盤の fixture 配置規則に従ってサンプル data を import する操作。
- `expected error`: import path が定まらず、後続 task が fixture 配置を再判断する。
- `observable point`: Storybook 最小 story の fixture import path、task-local 記録、後続 task から参照できる配置方針の有無。
- `related detail requirement type`: `data_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: fixture 配置を Storybook 基盤の必須成果物にする場合に採用候補になる。
- `conflict hint`: fixture 配置は未決事項であり、AI だけで最終決定すると人間判断を飛ばす可能性がある。

### CAND-SBF-005 review URL が task 成果物へ記録されない

- `source requirement`: `plan.md` の goal。review URL 記録方針を固定すること。
- `viewpoint`: 保存失敗、回復動作
- `candidate scenario id`: `CAND-SBF-005`
- `actor`: 人間レビュー担当者
- `trigger`: Storybook dev server は起動するが、review URL、確認状態、未確認理由を task 成果物へ残さない。
- `rejected operation`: 人間レビュー担当者が同じ URL と確認状態を後から追跡する操作。
- `expected error`: review URL 証跡が欠落し、後続 task が見た目レビューの入口を再作成する。
- `observable point`: task 成果物内の review URL 記録、起動 command、確認状態、未確認理由。
- `related detail requirement type`: `observability_requirement`, `testability_requirement`, `recovery_requirement`
- `adoption hint`: Storybook を人間見た目レビュー入口にする場合に採用候補になる。
- `conflict hint`: review URL を `ui-design.md` に残すか別 artifact に残すかが後続 task 側でも未決のため、記録場所が競合する可能性がある。

### CAND-SBF-006 Storybook に禁止された mock または実行フローを入れる

- `source requirement`: `plan.md` の constraints。gateway mock、backend DTO mock、実行フロー再現を Storybook に入れない。
- `viewpoint`: 設定不整合、競合候補
- `candidate scenario id`: `CAND-SBF-006`
- `actor`: 実装者
- `trigger`: 最小 story のために gateway mock、backend DTO mock、または翻訳実行フロー再現を Storybook 側へ追加する。
- `rejected operation`: Storybook 基盤 task の範囲内で外部連携や業務実行状態を再現する操作。
- `expected error`: task scope 逸脱として扱い、基盤確認の成果物に含めない。
- `observable point`: 追加 file の import 境界、story の props、mock の種類、backend DTO 参照の有無。
- `related detail requirement type`: `compatibility_requirement`, `consistency_requirement`, `failure_handling_requirement`
- `adoption hint`: Storybook 基盤を props で閉じた UI 部品確認に限定する場合に採用候補になる。
- `conflict hint`: 後続 Master Persona POC が richer fixture を要求する場合、基盤 task で扱う範囲との境界確認が必要になる。

### CAND-SBF-007 Storybook 運用の docs 正本化先が未確定のまま残る

- `source requirement`: `plan.md` の未決事項と canonicalization_targets。Storybook 運用をどの docs 正本へ反映するかが未決。
- `viewpoint`: 人間判断候補、保存失敗
- `candidate scenario id`: `CAND-SBF-007`
- `actor`: `designer`
- `trigger`: Storybook の script、build gate、review URL 記録方針は作られるが、docs 正本化先が決まらない。
- `rejected operation`: 実装 agent が未承認の docs 正本を更新して Storybook 運用を固定する操作。
- `expected error`: docs 正本化は未決として残し、`updating-docs` または人間判断へ戻す。
- `observable point`: task-local 成果物の canonicalization note、`docs/tech-selection.md` と `docs/lint-policy.md` への反映要否、未決事項の残り。
- `related detail requirement type`: `observability_requirement`, `consistency_requirement`, `testability_requirement`
- `adoption hint`: docs 正本化先の未決を scenario-design questions に流す場合に採用候補になる。
- `conflict hint`: この task の禁止事項は docs 正本変更を禁じているため、実装完了条件に docs 正本化を含める候補と競合する。

## Open Notes

- `human decision candidate`: fixture 配置場所、Storybook build を `frontend-local` に含めるか、review URL の記録先、Storybook 運用の docs 正本化先。
- `merge candidate`: `CAND-SBF-001` と `CAND-SBF-002` は Storybook 入口確認として統合候補になる。
- `rejection candidate`: `CAND-SBF-006` は禁止事項の検出候補であり、採用時も実装内容としては拒否動作を固定する候補になる。

## Completion Summary

- `viewpoint`: 失敗観点
- `candidate_count`: 7
- `artifact_path`: `docs/exec-plans/active/2026-05-18-storybook-foundation/scenario-candidates.failure.md`
- `task_artifact_root`: `docs/exec-plans/active/2026-05-18-storybook-foundation/`
- `target_diff`: Storybook 最小基盤。起動、build、fixture 配置、review URL 記録、scope 逸脱。
- `remaining_risk`: 採否、統合、競合解消、質問票化は `designer` が行う。
