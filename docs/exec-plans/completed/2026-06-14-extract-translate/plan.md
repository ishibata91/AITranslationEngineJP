# Task Plan: 2026-06-14-extract-translate

- `workflow`: work
- `status`: done（実装・検証完了。finalization 待ち）
- `task_id`: 2026-06-14-extract-translate
- `task_mode`: 重 task（画面が動く）
- `request_summary`: 最小縦切り。1 つの plugin を抽出 → SQLite → engine が読む → OpenAI 互換 provider で叙述文 1 種を翻訳 → dest を SQLite へ戻し → 最小実行画面に結果を表示するまでを 1 本通す。
- `goal`: 抽出 → 翻訳 → DB → 画面表示の end-to-end を 1 本通し、GUI から翻訳を走らせて結果を見られる最小の価値を出す。
- `constraints`: ペルソナ・マスター辞書・重複排除の作り込み・provider 4 系統・全レコード種別・xTranslator XML 出力・設定/監視画面・進捗 push の作り込みは対象外。
- `close_conditions`: GUI で plugin を指定し実行すると、翻訳された叙述文の原文・訳文が画面に出る。SQLite の `narration.dest` に訳文が入る。
- `execution_branch`: `claude/2026-06-14-extract-translate`（作成済み）
- `source_branch`: `master`
- `branch_point`: `06503b66641bd5793feae700040cf442e3d73e2a`
- `target_branch`: `master`

## 軽 / 重判定

- 画面が動くか: Y（最小実行画面を作る）。
- `docs/architecture.md` 反映が要るか: N（既存骨格の初実装で、層構成・依存方向・Wails 境界・C#↔Go 契約は `architecture.md` §1〜§7 に固定済み。T1 はそれを実装するだけで doc を変えない）。
- 判定: 重 task。フロー = `preparation → design → storybook →（画面表示）→ implementation → finalization`。

## 設計再評価（design-module 入口）

- `docs/architecture.md` 反映: N 確定。§7 ディレクトリ正本（`internal/api`・`engine`・`store`・`provider`・`model`・`bootstrap`、`db/`、`frontend/src/ui`・`gateway`）と §6 C#↔Go 契約（SQLite・repo-owned migration・schema version）に従い実装する。doc 本文は変えない。
- 設計差分図: 不要。decision table（architecture 反映 要→差分図 要）に照らし、反映 N のため差分図を作らない。代わりに `summary.md` の「図」へ T1 縦切りのデータ流れ図を置き、人間設計レビューの対象にする。
- lint/arch-lint config: `.golangci.yml` depguard と `.go-arch-lint.yml` は旧層名（controller/usecase/service 等）で、`arch` mode は不在の `.go-arch-lint.yml` を参照して現状失敗する。新層名へ config を書き換える。これは `architecture.md` §4 依存方向・§7 正本へ tooling を合わせる実装作業で、doc 変更ではない。
- 叙述文 1 種: `BOOK:DESC`（書物本文）を採る。散文の叙述文で概念モデルの叙述文に素直、1 record type・1 field・ordinal 常に 0、`dictionaries/Data/Dawnguard.esm` に実在し実機検証できる。装備 DESC への拡張は `TranslationCounts.Enumerate` のフィルタ追加だけで済む後続作業。

## Scope（含む / 含まない）

含む:
- `tools/extractor`（C#/Mutagen）に SQLite writer を追加し、1 plugin から叙述文 1 種（候補: 武器・防具の DESC、または BOOK 本文）を抽出して `er.md` の `narration` テーブルへ書く。
- `db/` に初期 SQL migration（最小: `narration`、必要なら識別カラムだけ）。
- `internal/store`（sqlx 薄アクセス）: `narration` の読み書き。
- `internal/engine`: `narration` を読み、provider で翻訳し `dest` を書く最小手続き。重複排除なし、文体無視で素直に投げる。
- `internal/provider`: AI クライアント interface ＋ OpenAI 互換 1 実装（`net/http`）。
- `internal/api` ＋ `internal/bootstrap`: Bind 公開（抽出+翻訳実行 command、結果取得 query）、composition root。
- `internal/model`: `narration` 等のデータ構造。
- `frontend/src/ui`: 最小実行画面（plugin パス入力、provider 接続設定 endpoint/key/model、実行ボタン、結果一覧）。Store。
- `frontend/src/gateway`: Wails bindings ラッパ。

含まない:
- 固有名の重複排除（訳の単位化）、`proper_noun`/`set_phrase`/`placement` の本格運用。
- ペルソナ（`speaker`/`race`/`faction`/`voice_type`、口調指示）。
- マスター辞書。
- provider の Gemini / xAI / Claude 実装。
- xTranslator XML 出力。
- 設定画面・監視画面・進捗 push の作り込み。

## Routing Notes

- `required_reading`:
  - `docs/concept-model.md`（概念正本。叙述文の扱い）
  - `docs/er.md`（テーブル定義。今回は `narration` 中心）
  - `docs/architecture.md`（§3 各責務、§6 C#↔Go 境界、§7 ディレクトリ正本）
  - `docs/system_requirements.md`（§1 AI 翻訳、対応 provider）
  - `docs/tech-selection.md`（採用技術）
  - `docs/references/xtranslator_ref.md`（識別キー・status 値域）
  - `docs/coding-guidelines-backend.md` / `-frontend.md` / `-tests.md`
- `canonicalization_targets`: `docs/architecture.md`（必要時のみ。原則 N）
- `validation_commands`: backend `go test ./...` / frontend `Vitest` / `svelte-check`（実装時に確定）

## 着手前メモ（compact 跨ぎ）

- これは 4 タスク計画の T1。依存関係と T2〜T4 の輪郭は各 plan.md 参照。
- T1 の最小縦切りの確定事項: 叙述文 1 種、provider は OpenAI 互換、画面込み（重 task）。
- 次アクション: master の未コミット docs 変更（ER 設計・掃除）を整理し、`preparation-module` で branch `claude/2026-06-14-extract-translate` を切って着手。実装本体は Claude 本体が 1 文脈で縦通し（`coding-protocol` / `implementation-module`）。

## Branch Status

- `branch_ready`: yes（`claude/2026-06-14-extract-translate`、分岐元 `06503b66`）
- `remote_operation`: `not-performed`

## HITL Status

- `detail_spec_hitl`: `required-after-design-bundle`
- `frontend_human_review`: `承認済み`（Storybook 人間レビューで「UI は OK」。記録は `storybook-review-loop.md`）

## storybook-module 結果（合意済み frontend 保護）

- 承認済み画面: `翻訳実行`（`Screens/翻訳実行`、6 状態）。
- 承認済み部品（`UI Components/*`）: `StatusBadge`・`TextField`・`SelectField`・`FileSelectField`・`TranslationResultRow`・`ResultsPanel`。
- 表示規則: Tailwind v4 ＋ daisyUI v5 の独自テーマ `dovahkael`。入力は plugin ファイル選択、モデルは `getModels` 取得の選択、結果は見開き対訳。文言・状態網羅は `storybook-review-loop.md` 参照。
- 変更禁止範囲（後続実装で表示を変えない境界）: 承認済み svelte 表示コンポーネントの構造、props 形、daisyUI ベースの style、画面の状態網羅。
- 反映先ファイル・通常分類 story 一覧: `storybook-review-loop.md` に記載。
- 検証: `build-storybook` 成功、`eslint`・`tsc` 通過。`knip`・`boundaries` は production entry 未配線の greenfield 構造状態で未達（implementation-module の配線で解消）。

## Outcome

- preparation 完了（branch 作成・docs コミット）。
- storybook-module 完了（画面・部品の表示実装と Storybook 人間レビュー承認、通常分類復帰）。
- implementation-module 完了（TDD で backend 縦通し・C# SQLite writer・frontend production 配線）。`status = done`。
- 検証：Go test 緑、backend lint 緑、frontend lint 緑、C# 17 テスト緑、build-storybook 緑、実 app で end-to-end 目視確認（OpenAI 互換モックに対し getModels・抽出 65 件・翻訳・対訳表示）。詳細は `docs/changelog.md` の T1 entry。
- 次：finalization-module（branch のレビュー・master への統合）。実 LM Studio（`127.0.0.1:1234`、API キーなし）での実訳確認はユーザー側。

## finalization-module

### 正本化判断

- `docs/architecture.md` 反映: **不要（正本反映 N）**。T1 は §1〜§7（層構成・依存方向・強い制約・Wails 境界）に**従って実装**しただけで、構造仕様を変えていない。新規の恒久仕様は生じていないため正本反映の対象なし。
- 注: §8「現在の状態と移行」は T1 前の greenfield 状態（"backend 削減済み・writer 未実装" 等）のままで、骨格実装後は古い。ただし §8 は正本反映対象（§1〜§7 の構造仕様）ではなく、`architecture.md` 本文変更は人間承認を要するため、本 finalization では変更しない。§8 の現状追従は後続で扱う。
- 正本反映: 実施なし（対象なし）。

### 作業 commit

- commit: `652642f0`（feat: 抽出→翻訳→DB→画面の縦切り（T1）を実装）。docs commit `ffe6db3b` も branch に含む。
- 変更: 60 ファイル（+3045/-153）。生成物（dist/wailsjs/storybook-static）と `.DS_Store`・dev sqlite は gitignore で除外。
- 検証: Go test 緑、backend lint（format/vet/static/arch/module）0 issues、frontend lint（eslint/tsc/knip/boundaries）緑、C# 17 テスト緑。
- 残留リスク: 大量同期翻訳の進捗 UI なし、書物本文の HTML 様タグ未整理、フォント CDN（`docs/changelog.md` の残課題参照）。

### local merge

- command: `git merge --no-ff claude/2026-06-14-extract-translate`（source `claude/2026-06-14-extract-translate` → target `master`）。
- 結果: merge commit `0c85500e`。conflict なし。
- remote 操作なし（push / tag / remote delete 未実施）。

### merge 後検証

- master 上で Go test 緑、backend lint 0 issues、frontend lint 緑。

### completed 移動

- `docs/exec-plans/active/2026-06-14-extract-translate/` → `docs/exec-plans/completed/2026-06-14-extract-translate/`。
