# Scenario Candidates: 2026-05-18-storybook-foundation / lifecycle

- `generator`: `lifecycle`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `STORYBOOK-FOUNDATION`

## Generator Scope

- `viewpoint`: lifecycle
- `included_sources`: `./plan.md`, `docs/tech-selection.md`, `docs/coding-guidelines-frontend.md`, `docs/lint-policy.md`, `frontend/package.json`, `frontend/vite.config.ts`
- `excluded_sources`: Master Persona の部品化、画面再設計、gateway mock、backend DTO mock、実行フロー再現、プロダクトコード実装指示、プロダクトテスト実装指示、docs 正本化判断
- `generation_notes`: Storybook 最小基盤が作成され、起動、build、後続 story 追加、終了条件へ進む lifecycle 候補だけを列挙する。最終採否、統合、競合解消は `designer` が扱う。

## Candidate Scenarios

### CAND-STORYBOOK-FOUNDATION-001 Storybook 最小基盤を作成する

- `source requirement`: `plan.md` の `goal` は Storybook scripts、config、build 検証、fixture 配置方針、review URL 記録方針を最小単位で固定するとしている。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-STORYBOOK-FOUNDATION-001`
- `lifecycle stage`: 作成
- `start condition`: `frontend/package.json` に Storybook 用 script がなく、`.storybook/` がまだない。
- `actor`: Codex implementation lane
- `trigger`: Storybook 最小基盤の実装範囲が承認される。
- `expected outcome`: `frontend/` 基準で Storybook の dev と build を呼び出せる script と最小 config が作られる。
- `observable point`: `frontend/package.json` の Storybook script、`.storybook/` の config、既存 `frontend/vite.config.ts` の alias と Svelte 前提を崩していないこと。
- `related detail requirement type`: `testability_requirement`, `compatibility_requirement`
- `adoption hint`: 最小基盤そのものの作成確認として採用候補になる。
- `conflict hint`: Storybook の追加設定が既存 Vite alias や Svelte 5 前提と食い違う場合、frontend tooling 境界の統合判断が必要になる。

### CAND-STORYBOOK-FOUNDATION-002 Storybook dev 入口を起動できる

- `source requirement`: `plan.md` の `close_conditions` は Storybook dev の入口が動くことを求めている。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-STORYBOOK-FOUNDATION-002`
- `lifecycle stage`: 起動
- `start condition`: Storybook script と config が作成済みである。
- `actor`: Codex implementation lane
- `trigger`: Storybook dev script を実行する。
- `expected outcome`: Storybook dev server が起動し、空または最小サンプル story を表示できる。
- `observable point`: dev server の起動結果、表示可能な Storybook URL、最小 story の表示状態。
- `related detail requirement type`: `testability_requirement`
- `adoption hint`: 人間見た目レビューの入口確認として採用候補になる。
- `conflict hint`: review URL の記録先が `ui-design.md` か別 artifact か未確定のため、記録先は `designer` の統合判断に残す。

### CAND-STORYBOOK-FOUNDATION-003 Storybook build を検証できる

- `source requirement`: `plan.md` の `validation_commands` は `npm --prefix frontend run build-storybook` を含む。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-STORYBOOK-FOUNDATION-003`
- `lifecycle stage`: build
- `start condition`: Storybook 最小 config と最小 story が作成済みである。
- `actor`: Codex implementation lane
- `trigger`: `npm --prefix frontend run build-storybook` を実行する。
- `expected outcome`: Storybook の静的 build が完了し、基盤が CI や人間確認へ渡せる状態になる。
- `observable point`: build command の終了状態、Storybook build 出力、frontend の既存 build script と競合していないこと。
- `related detail requirement type`: `testability_requirement`, `compatibility_requirement`
- `adoption hint`: task の build 完了条件として採用候補になる。
- `conflict hint`: Storybook build を `frontend-local` に含めるか別 gate にするかは `plan.md` の未決事項である。

### CAND-STORYBOOK-FOUNDATION-004 後続 task が story を追加できる

- `source requirement`: `plan.md` の `close_conditions` は後続 Master Persona task が story を追加できることを求めている。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-STORYBOOK-FOUNDATION-004`
- `lifecycle stage`: 後続 story 追加
- `start condition`: Storybook 最小基盤が作成され、dev と build の入口が確認済みである。
- `actor`: 後続 frontend task
- `trigger`: 後続 task が UI Component または screen local component の story を追加する。
- `expected outcome`: 後続 story が backend DTO mock や gateway mock に依存せず、props と固定 fixture で表示確認できる。
- `observable point`: story 追加先、fixture 参照先、`docs/coding-guidelines-frontend.md` の UI Component 境界に沿った props 入力。
- `related detail requirement type`: `testability_requirement`, `compatibility_requirement`, `data_requirement`
- `adoption hint`: 後続 Master Persona task への引き継ぎ条件として採用候補になる。
- `conflict hint`: fixture を `frontend/src/ui/**/__fixtures__` に置くか Storybook 専用 directory に置くかは `plan.md` の未決事項である。

### CAND-STORYBOOK-FOUNDATION-005 Storybook 基盤 task を終了できる

- `source requirement`: `plan.md` の `close_conditions` は Storybook dev / build の入口、空または最小サンプル story、後続 story 追加可能性を終了条件にしている。
- `viewpoint`: lifecycle
- `candidate scenario id`: `CAND-STORYBOOK-FOUNDATION-005`
- `lifecycle stage`: 終了
- `start condition`: 作成、起動、build、後続 story 追加前提の確認が完了している。
- `actor`: implement_lane
- `trigger`: Storybook 基盤 task の完了判定を行う。
- `expected outcome`: Storybook 基盤は完了扱いになり、Master Persona 部品化や画面再設計を含まないまま後続 task へ渡せる。
- `observable point`: dev / build の検証結果、最小 story の有無、後続 story 追加に必要な配置方針、残った未決事項。
- `related detail requirement type`: `testability_requirement`, `state_requirement`
- `adoption hint`: lifecycle の終点確認として採用候補になる。
- `conflict hint`: Storybook 運用をどの docs 正本へ反映するかは未決であり、この候補では docs 正本化を確定しない。

## Open Notes

- `human decision candidate`: Storybook build を `frontend-local` に含めるか、別 gate にするか。
- `human decision candidate`: fixture を `frontend/src/ui/**/__fixtures__` に置くか、Storybook 専用 directory に置くか。
- `human decision candidate`: Storybook 運用をどの docs 正本へ反映するか。
- `merge candidate`: 作成、起動、build の 3 候補は、最終シナリオ表で Storybook 基盤の正常 lifecycle として統合される可能性がある。
- `merge candidate`: 後続 story 追加と終了条件は、後続 task への handoff scenario として統合される可能性がある。
- `rejection candidate`: gateway mock、backend DTO mock、実行フロー再現を前提にする候補は、この task の制約に反するため除外候補になる。

