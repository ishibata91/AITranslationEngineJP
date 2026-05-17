# Scenario Candidates: 2026-05-18-storybook-foundation / state-transition

- `generator`: `state-transition`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `STORYBOOK-STATE`

## Generator Scope

- `viewpoint`: `state-transition`
- `included_sources`: `./plan.md`, `docs/tech-selection.md`, `docs/coding-guidelines-frontend.md`, `docs/lint-policy.md`, `frontend/package.json`, `frontend/vite.config.ts`
- `excluded_sources`: product code change, product test change, docs canonicalization, `.codex/`, Master Persona componentization, screen redesign, gateway mock, backend DTO mock, execution flow reproduction
- `generation_notes`: Storybook 最小基盤の導入状態だけを候補化する。最終シナリオ表、採否、統合、競合解消は `designer` に残す。

## Candidate Scenarios

### CAND-STORYBOOK-STATE-001 未導入から設定済みへ遷移できる

- `source requirement`: Storybook scripts、config、最小 story を追加し、Master Persona 部品化前の基盤だけを作る。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-STORYBOOK-STATE-001`
- `actor`: frontend implementer
- `trigger`: Storybook 未導入の `frontend/package.json` と Vite 設定を起点に、Storybook 最小設定を追加する。
- `expected outcome`: Storybook 用 script、設定、最小 story が存在し、状態は未導入から設定済みへ遷移する。プロダクト画面の再設計、Master Persona 部品化、backend 実行再現は発生しない。
- `observable point`: `frontend/package.json` に Storybook 起動と build の入口が見える。Storybook 設定は `Svelte 5`、`TypeScript`、`Vite` 前提と矛盾しない。
- `related detail requirement type`: `state_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: 初回導入の正常遷移として採用候補になる。
- `conflict hint`: Storybook 設定が backend DTO mock や gateway mock を前提にする場合、task 制約と競合する。

### CAND-STORYBOOK-STATE-002 設定済みから dev 起動可能へ遷移できる

- `source requirement`: Storybook dev の入口が動き、空または最小サンプル story で基盤確認ができる。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-STORYBOOK-STATE-002`
- `actor`: frontend implementer
- `trigger`: 設定済み状態で Storybook dev script を実行する。
- `expected outcome`: Storybook dev server が起動し、最小 story を表示できる状態へ遷移する。Wails runtime、gateway、backend DTO は起動条件にならない。
- `observable point`: dev command が終了せずに待機し、ブラウザ確認用 URL を提示できる。表示対象は最小 story に限定される。
- `related detail requirement type`: `state_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: 実装後ブラウザ確認の前提状態として採用候補になる。
- `conflict hint`: dev 起動が Wails dev や backend 実行に依存する場合、Storybook 最小基盤の状態遷移として競合する。

### CAND-STORYBOOK-STATE-003 設定済みから build 通過へ遷移できる

- `source requirement`: build-storybook 検証を対象にする。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-STORYBOOK-STATE-003`
- `actor`: frontend implementer
- `trigger`: 設定済み状態で Storybook build command を実行する。
- `expected outcome`: Storybook の静的 build が成功し、状態は build 通過へ遷移する。frontend の既存 build、lint、型検査の責務と衝突しない。
- `observable point`: build command が成功し、Storybook の出力先が生成される。既存 `frontend` build の出力先 `dist` と混同しない。
- `related detail requirement type`: `state_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: close condition の build 検証として採用候補になる。
- `conflict hint`: Storybook build の出力先が既存 app build artifact と同じ扱いになる場合、lint と build artifact 管理で競合する。

### CAND-STORYBOOK-STATE-004 build 通過から後続追加可能へ遷移できる

- `source requirement`: 後続 Master Persona task が story を追加できる。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-STORYBOOK-STATE-004`
- `actor`: later frontend implementer
- `trigger`: Storybook dev と build が通過した状態で、後続 task が新しい story を追加する。
- `expected outcome`: story 追加場所、fixture 配置方針、review URL 記録方針が読み取れる状態へ遷移する。Master Persona 固有の状態 fixture はこの task で確定しない。
- `observable point`: 後続 task が最小 story を手本にして story を追加できる。fixture 方針は gateway mock、backend DTO mock、実行フロー再現を要求しない。
- `related detail requirement type`: `state_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: 後続作業へ渡す受け入れ条件として採用候補になる。
- `conflict hint`: fixture 方針が Master Persona の部品境界を先に固定する場合、task 制約と競合する。

### CAND-STORYBOOK-STATE-005 再実行しても Storybook 基盤状態が壊れない

- `source requirement`: Storybook scripts、config、build 検証、fixture 配置方針を最小単位で固定する。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-STORYBOOK-STATE-005`
- `actor`: frontend implementer
- `trigger`: Storybook 設定済み状態で dev command または build command を繰り返し実行する。
- `expected outcome`: 再実行しても story、設定、依存関係、出力 artifact が重複作成や設定破損を起こさない。状態は dev 起動可能または build 通過に戻る。
- `observable point`: 連続実行後も `package.json` の script、Storybook 設定、最小 story の数と役割が増殖しない。
- `related detail requirement type`: `冪等性_requirement`, `state_requirement`, `testability_requirement`
- `adoption hint`: tooling foundation の回帰防止候補として採用候補になる。
- `conflict hint`: 再実行のたびに一時成果物を source tree に残す場合、lint や未使用検出の対象範囲と競合する可能性がある。

### CAND-STORYBOOK-STATE-006 禁止された mock 導入では状態遷移しない

- `source requirement`: gateway mock、backend DTO mock、実行フロー再現を Storybook に入れない。
- `viewpoint`: `state-transition`
- `candidate scenario id`: `CAND-STORYBOOK-STATE-006`
- `actor`: frontend implementer
- `trigger`: Storybook story または fixture に gateway mock、backend DTO mock、実行フロー再現を追加しようとする。
- `expected outcome`: Storybook 基盤は後続追加可能状態へ遷移しない。最小基盤は UI component 表示確認に閉じ、backend 境界や実行フローの代替実装を持たない。
- `observable point`: story と fixture に generated binding、backend DTO、gateway fake、実行フロー再現の import がない。
- `related detail requirement type`: `state_requirement`, `compatibility_requirement`, `testability_requirement`
- `adoption hint`: 禁止遷移の受け入れ条件として採用候補になる。
- `conflict hint`: external-integration 観点が mock 境界を扱う場合、この候補と重複する可能性がある。採否と統合は `designer` が判断する。

## Open Notes

- `human decision candidate`: Storybook build を `frontend-local` に含めるか、別 gate にするかは未決である。fixture の正本配置と Storybook 運用 docs の反映先も未決である。
- `merge candidate`: `CAND-STORYBOOK-STATE-002` と `CAND-STORYBOOK-STATE-003` は Storybook 設定後の検証遷移として統合できる可能性がある。
- `rejection candidate`: backend 実行、Wails dev、gateway mock、backend DTO mock を前提にする候補は今回の state-transition 候補から除外する。
