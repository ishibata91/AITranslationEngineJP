# ER Feasibility Gaps

- workflow: heavy
- status: completed
- architect: codex
- research: none
- coder: codex
- reviewer: architect
- scope: `docs/er-draft.md`, `docs/spec.md`, `docs/architecture.md`, `docs/executable-specs.md`

## Request Summary

- ER 図の実現性レビューで判明した不足を修正し、`extractData.pas`、要件、最終出力の成立条件に整合するようにする

## Investigation Summary

- Facts:
- `extractData.pas` は `dialogue_groups`, `quests`, `items`, `magic`, `locations`, `system`, `messages`, `load_screens`, `npcs`, `cells` を JSON 出力する
- `extractData.pas` の `INFO` 応答は `voicetype` をまだ JSON 出力しているが、ER では `NPC.voice` 正本へ寄せる方針になっている
- `TRANSLATION_JOB` は単一 `PLUGIN_EXPORT` 参照しか持たず、要件の複数入力ファイルを直接表現できない
- `JOB_RECORD` はカテゴリ横断の 1 行モデルで、xTranslator XML が要求する `REC` / `FIELD` / `Source` / `Dest` 単位の出力を lossless に保持できない
- mod 追加 NPC のジョブ内ペルソナを保持し、UI から観測するためのテーブルが ER にない
- Options:
- 既存 `JOB_RECORD` を拡張してフィールド単位へ寄せる
- `TRANSLATION_UNIT` 系の新エンティティを導入して翻訳単位、出力単位、フェーズ観測単位を分離する
- Risks:
- ER だけ直して `spec.md` / `architecture.md` の前提を更新しないと文書間で整合が崩れる
- `extractData.pas` の現行 JSON と DB 正規化後モデルを区別して書かないと入力契約が曖昧になる
- Unknowns:
- 標準配布形式の詳細な出力フィールドは未定義のままなので、今回は xTranslator 互換出力に必要な最小契約までを固定する

## Implementation Plan

- 入力 ER を `extractData.pas` の現行 JSON と DB 正規化後モデルの差分が分かる形に修正する
- 翻訳ジョブ ER を複数 `PLUGIN_EXPORT`、ジョブ内ペルソナ、フィールド単位の翻訳単位、xTranslator 出力メタデータを表現できる構造へ更新する
- `spec.md` と `architecture.md` に複数入力ジョブと出力生成責務を補強し、`executable-specs.md` に検証観点を追加する

## Delegation Map

- Research: なし
- Coder: codex が実装
- Worker: なし

## Acceptance Checks

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`

## Record Updates

- `docs/er-draft.md`
- `docs/spec.md`
- `docs/architecture.md`
- `docs/executable-specs.md`

## Outcome

- `docs/er-draft.md` を raw JSON 互換と DB canonical モデルの差分が分かる形に更新した
- 翻訳ジョブ ER を `JOB_PLUGIN_EXPORT`、`TRANSLATION_UNIT`、`JOB_TRANSLATION_UNIT`、`JOB_PERSONA_ENTRY`、`JOB_OUTPUT_ARTIFACT` を含む構成へ更新した
- `spec.md`、`architecture.md`、`executable-specs.md` に複数入力ジョブ、job-local persona、xTranslator 出力の成立条件を反映した

## Verification

- `powershell -File scripts/harness/run.ps1 -Suite structure`
- `powershell -File scripts/harness/run.ps1 -Suite design`
- `powershell -File scripts/harness/run.ps1 -Suite all`
