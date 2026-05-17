# Scenario Candidates: 2026-05-18-storybook-foundation / external-integration

- `generator`: `external-integration`
- `source_plan`: `./plan.md`
- `scenario_design_target`: `./scenario-design.md`
- `topic_abbrev`: `STORYBOOK-EXT`
- `candidate_count`: 5

## Generator Scope

- `viewpoint`: Storybook 最小基盤の外部境界、静的表示境界、browser review URL、backend 非連携境界を扱う。
- `included_sources`: `plan.md`, `docs/tech-selection.md`, `docs/coding-guidelines-frontend.md`, `docs/lint-policy.md`, `frontend/package.json`, `frontend/vite.config.ts`, `docs/spec.md`, `docs/architecture.md`, `docs/er.md`, `docs/screen-design/README.md`, `docs/detail-specs/README.md`, `docs/scenario-tests/README.md`
- `excluded_sources`: Master Persona 固有の部品分割、画面再設計、gateway mock、backend DTO mock、実行フロー再現、real AI provider、secret store、network provider 連携。
- `generation_notes`: 最終シナリオ表、採否、統合判断は行わない。Storybook は backend 連携検証の代替にしない。

## Candidate Scenarios

### CAND-STORYBOOK-EXT-001 Storybook dev server が Vite と Svelte の表示だけを公開する

- `source requirement`: `plan.md` は Storybook dev / build の入口が動くことを close condition にする。`docs/tech-selection.md` は frontend build tool を `Vite`、UI framework を `Svelte 5` とする。
- `viewpoint`: external-integration / adapter 境界
- `candidate scenario id`: `CAND-STORYBOOK-EXT-001`
- `external boundary`: Storybook dev server と browser の間のローカル表示境界。
- `actor`: 実装者
- `trigger`: Storybook dev server を起動し、browser review URL を開く。
- `expected outcome`: browser で最小 story が表示される。Storybook は Wails runtime、generated `wailsjs`、backend Controller、AI provider を呼ばない。
- `fake_or_stub`: 最小 story の固定 props。gateway mock と backend DTO mock は使わない。
- `observable point`: Storybook canvas または docs view に最小 story が表示される。表示確認の URL は task-local artifact に記録できる。
- `related detail requirement type`: `testability_requirement`, `compatibility_requirement`
- `adoption hint`: Storybook dev 起動と browser 表示確認を、backend 連携確認とは別の受け入れ観点にできる。
- `conflict hint`: lifecycle 観点が実行フロー確認を Storybook に含める場合、この候補の非連携境界と衝突する。

### CAND-STORYBOOK-EXT-002 Storybook 静的 build が backend なしで成果物を作る

- `source requirement`: `plan.md` は `npm --prefix frontend run build-storybook` を validation command に置く。`frontend/package.json` は現状 Storybook script を持たないため、最小基盤で script 追加が対象になる。
- `viewpoint`: external-integration / ファイル境界
- `candidate scenario id`: `CAND-STORYBOOK-EXT-002`
- `external boundary`: Storybook static build output と filesystem の間の生成物境界。
- `actor`: 実装者
- `trigger`: `npm --prefix frontend run build-storybook` を実行する。
- `expected outcome`: 静的 Storybook build が成功し、生成物は frontend tooling の成果物として作られる。Wails build、backend 実行、AI provider 通信は発生しない。
- `fake_or_stub`: 最小 story、固定 props、静的 asset。temp DB、fake provider、fake secret store は使わない。
- `observable point`: build command の成功、生成 directory の存在、backend 接続を前提にしない static asset 構成。
- `related detail requirement type`: `testability_requirement`, `compatibility_requirement`
- `adoption hint`: `frontend-local` に含めるかどうかは別判断にし、単独 command の通過条件として候補化できる。
- `conflict hint`: failure 観点が backend 未起動を Storybook build failure とみなす場合、この候補の前提と衝突する。

### CAND-STORYBOOK-EXT-003 story は UI Component へ固定 props を渡し、gateway 境界を越えない

- `source requirement`: `docs/architecture.md` は `UI Component` が backend DTO、generated binding、`Store`、`Gateway` を直接扱わないと定義する。`docs/coding-guidelines-frontend.md` は generated `wailsjs` と backend DTO の import を gateway 境界に閉じ込める。
- `viewpoint`: external-integration / adapter 境界
- `candidate scenario id`: `CAND-STORYBOOK-EXT-003`
- `external boundary`: UI Component props と backend adapter の分離境界。
- `actor`: 実装者
- `trigger`: 最小 story が表示対象 component を import し、固定 props を渡す。
- `expected outcome`: story は UI 表示を確認できる。story は Gateway、generated `wailsjs`、backend DTO、RuntimeEventAdapter を import しない。
- `fake_or_stub`: component props fixture、view model fixture。gateway mock と backend DTO mock は使わない。
- `observable point`: story source と import 境界で、backend 由来の依存が混入していないことを確認できる。
- `related detail requirement type`: `security_requirement`, `testability_requirement`, `compatibility_requirement`
- `adoption hint`: 後続 Master Persona task の story 追加時にも、props 境界の確認候補として再利用できる。
- `conflict hint`: operation-audit 観点が Storybook に runtime event や backend response の再現を求める場合、この候補の props-only 境界と衝突する。

### CAND-STORYBOOK-EXT-004 browser review URL は Storybook の表示確認 URL として記録する

- `source requirement`: `plan.md` は review URL 記録方針を固定することを goal にする。`docs/screen-design/README.md` は visual design と画面別設計を正本として扱うが、この task は画面再設計を対象にしない。
- `viewpoint`: external-integration / browser 境界
- `candidate scenario id`: `CAND-STORYBOOK-EXT-004`
- `external boundary`: local browser と Storybook server または Storybook static preview の間の確認境界。
- `actor`: 実装者
- `trigger`: 実装後に Storybook review URL を task-local artifact へ記録する。
- `expected outcome`: review URL は Storybook の表示確認先を指す。fakeAPI URL、Wails runtime URL、backend API URL を Storybook review URL として扱わない。
- `fake_or_stub`: Storybook story の固定 props。fakeAPI scenario は使わない。
- `observable point`: task-local artifact に review URL、確認対象 story、未確認理由を記録できる。
- `related detail requirement type`: `testability_requirement`, `observability_requirement`
- `adoption hint`: 人間見た目レビューの入力を Storybook URL に寄せる候補として使える。
- `conflict hint`: operation-audit 観点が review URL を恒久 docs へ直接正本化する場合、この task の task-local 記録境界と衝突する。

### CAND-STORYBOOK-EXT-005 Storybook は backend 連携検証の代替にならない

- `source requirement`: `plan.md` は gateway mock、backend DTO mock、実行フロー再現を Storybook に入れないと制約する。`docs/spec.md` は AI provider、secret、翻訳 job、DB、xTranslator 出力を含む業務要件を定義するが、この task は Storybook 最小基盤だけを対象にする。
- `viewpoint`: external-integration / provider 境界、network 境界、secret 境界
- `candidate scenario id`: `CAND-STORYBOOK-EXT-005`
- `external boundary`: Storybook 表示確認と backend / AI provider / secret store / DB / filesystem business flow の非連携境界。
- `actor`: designer
- `trigger`: 最終シナリオ統合時に、Storybook 確認の証明対象を分類する。
- `expected outcome`: Storybook は UI 部品の表示確認に限定される。AI provider 通信、API key 保存、翻訳 job 実行、DB 永続化、xTranslator XML 出力の成立証明には使わない。
- `fake_or_stub`: props fixture、story fixture。fake provider、fake secret store、temp DB、XML parser は使わない。
- `observable point`: scenario-design で Storybook 確認と backend 連携検証の実行段階を分けて記録できる。
- `related detail requirement type`: `security_requirement`, `data_requirement`, `testability_requirement`
- `adoption hint`: Storybook 基盤 task の受け入れ条件から、backend 連携の成功条件を除外する候補として使える。
- `conflict hint`: external-integration 以外の候補が Storybook で gateway mock や backend DTO mock を使う前提を置く場合、この候補と競合する。

## Open Notes

- `human decision candidate`: Storybook build を `frontend-local` に含めるか、別 gate にするかは `plan.md` の未決事項である。
- `human decision candidate`: fixture を `frontend/src/ui/**/__fixtures__` に置くか、Storybook 専用 directory に置くかは `plan.md` の未決事項である。
- `human decision candidate`: Storybook 運用をどの docs 正本へ反映するかは `plan.md` の未決事項である。
- `merge candidate`: `CAND-STORYBOOK-EXT-001` と `CAND-STORYBOOK-EXT-004` は、browser review URL の扱いで統合される可能性がある。
- `merge candidate`: `CAND-STORYBOOK-EXT-003` と `CAND-STORYBOOK-EXT-005` は、backend 非連携境界として統合される可能性がある。
- `rejection candidate`: Storybook に gateway mock、backend DTO mock、実行フロー再現を入れる候補は、引き継ぎ入力の禁止事項により不採用候補になる。
