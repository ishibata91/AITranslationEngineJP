# Scenario Candidates: 2026-05-18-storybook-foundation / actor-goal

- `generator`: `actor-goal`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `SBF`

## Generator Scope

- `viewpoint`: Storybook 最小基盤を使う実行者の目的、開始操作、成功体験を候補化する。
- `included_sources`: `plan.md`, `docs/tech-selection.md`, `docs/coding-guidelines-frontend.md`, `docs/lint-policy.md`, `frontend/package.json`, `frontend/vite.config.ts`
- `excluded_sources`: Master Persona の部品化、画面再設計、gateway mock、backend DTO mock、実行フロー再現、プロダクトコード変更指示、最終シナリオ採否
- `generation_notes`: 候補は actor-goal 観点に限定する。採否、統合、競合解消は `designer` に残す。

## Candidate Scenarios

### CAND-SBF-001 Storybook 開発入口で最小 story を確認する

- `source requirement`: `plan.md` の `goal` と `close_conditions`。Storybook dev 入口が動き、空または最小サンプル story で基盤確認できること。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-SBF-001`
- `actor`: 後続 UI 実装者
- `trigger`: 後続 UI 実装者が `frontend/package.json` の Storybook 開発用 script を実行する。
- `expected outcome`: Storybook の画面が起動し、最小 story を見て story 追加の起点を確認できる。
- `observable point`: Storybook のブラウザ表示で最小 story が選択できる。起動 command が失敗しない。
- `related detail requirement type`: `success_requirement`, `testability_requirement`
- `adoption hint`: Storybook 基盤の最小正常系として採用候補にできる。
- `conflict hint`: Storybook に既存画面フローや backend 連携を再現させる候補と競合する。

### CAND-SBF-002 Storybook build で静的成果物を確認する

- `source requirement`: `plan.md` の `validation_commands` と `close_conditions`。`npm --prefix frontend run build-storybook` で build 検証できること。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-SBF-002`
- `actor`: 実装検証者
- `trigger`: 実装検証者が Storybook build 用 script を実行する。
- `expected outcome`: Storybook の静的 build が成功し、最小 story の構成が build 時にも破綻しない。
- `observable point`: build command の終了結果が成功で返る。Storybook build 出力先が生成される。
- `related detail requirement type`: `testability_requirement`, `compatibility_requirement`
- `adoption hint`: 最小基盤の完了条件を自動検証へつなげる候補にできる。
- `conflict hint`: Storybook build を `frontend-local` に含めるか別 gate にするかは未決のため、検証段階の統合時に競合しうる。

### CAND-SBF-003 Fixture 配置方針に沿って後続 story を追加できる

- `source requirement`: `plan.md` の `goal` と `未決事項`。fixture 配置方針を固定し、後続 Master Persona task が story を追加できること。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-SBF-003`
- `actor`: 後続 Master Persona task の実装者
- `trigger`: 後続実装者が固定された fixture 配置方針を見て、新しい story 用 fixture を追加する。
- `expected outcome`: story 用 fixture の置き場所が分かり、backend DTO mock や gateway mock を作らずに表示用データを準備できる。
- `observable point`: fixture の配置先と読み込み方法が Storybook story から確認できる。
- `related detail requirement type`: `data_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: 後続 task の story 追加しやすさを確認する候補にできる。
- `conflict hint`: fixture を `frontend/src/ui/**/__fixtures__` に置く案と Storybook 専用 directory に置く案が未決である。

### CAND-SBF-004 Human review 用 URL を記録できる

- `source requirement`: `plan.md` の `goal`。review URL 記録方針を最小単位で固定すること。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-SBF-004`
- `actor`: 人間レビュー依頼者
- `trigger`: 人間レビュー依頼者が Storybook 起動後に、レビュー用 URL を task-local 成果物へ記録する。
- `expected outcome`: 人間レビュアーが同じ Storybook 表示へ到達できる URL と起動条件を確認できる。
- `observable point`: task-local 成果物に Storybook URL、起動 command、対象 story が記録されている。
- `related detail requirement type`: `observability_requirement`, `testability_requirement`
- `adoption hint`: Storybook を人間確認へ接続する候補にできる。
- `conflict hint`: URL 記録先を task-local 成果物だけにするか、docs 正本化対象へ進めるかは未決である。

### CAND-SBF-005 既存 frontend 技術基盤と同じ import alias で story を書ける

- `source requirement`: `docs/tech-selection.md` の Svelte 5、TypeScript、Vite 採用。`frontend/vite.config.ts` の `@ui`、`@application`、`@controller` alias。
- `viewpoint`: actor-goal
- `candidate scenario id`: `CAND-SBF-005`
- `actor`: UI component story 作成者
- `trigger`: story 作成者が既存 frontend の import alias を使って最小 story を書く。
- `expected outcome`: story の import が Vite と同じ前提で解決され、Storybook 専用の別名体系を覚えずに story を追加できる。
- `observable point`: Storybook dev と Storybook build の両方で alias import が解決される。
- `related detail requirement type`: `compatibility_requirement`, `testability_requirement`
- `adoption hint`: Storybook 設定が既存 frontend build 設定とずれないことを確認する候補にできる。
- `conflict hint`: Storybook 側だけで alias を再定義する場合、Vite 正本との二重管理が発生しうる。

## Open Notes

- `human decision candidate`: Storybook build を `frontend-local` に含めるか、別 gate として扱うか。
- `human decision candidate`: fixture を `frontend/src/ui/**/__fixtures__` に置くか、Storybook 専用 directory に置くか。
- `human decision candidate`: Storybook 運用と review URL 記録方針をどの docs 正本へ反映するか。
- `merge candidate`: `CAND-SBF-001` と `CAND-SBF-002` は、Storybook 最小基盤の dev / build 正常系として統合できる可能性がある。
- `merge candidate`: `CAND-SBF-003` と `CAND-SBF-005` は、後続 story 追加者の作業しやすさとして統合できる可能性がある。
- `rejection candidate`: gateway mock、backend DTO mock、実行フロー再現を Storybook に入れる候補は、この task の禁止事項により不採用候補とする。
