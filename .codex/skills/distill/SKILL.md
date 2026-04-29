---
name: distill
description: Codex 側の共通文脈圧縮 skill。入口情報を facts、constraints、gaps、required_reading へ整理するための共通知識と圧縮パターンを提供する。
---

# Distill

## 目的

`distill` は、入口情報を短く整理するための共通知識である。
圧縮粒度、重複除去、facts / inferred / gap の分離、downstream が読む順番の作り方を扱う。

設計向けの観点は `distill-design` が持つ。
調査向けの観点は `distill-investigate` が持つ。
`distill` 本体は、どちらにも共通する圧縮の見方だけを持つ。

## いつ参照するか

- `implement-lane` の次判断に必要な repo 文脈を圧縮する時
- 設計向けまたは調査向けの詳細観点を読む前に、共通の圧縮粒度をそろえる時
- user request、active plan、docs、関連 skill の重複を短く整理する時

## 参照しない場合

- requirements、UI、scenario、diagram の詳細観点だけが必要な時
- 観測対象、再現条件、未観測情報の詳細観点だけが必要な時
- human review 済み `implementation-scope` から実装前 context を作る時
- fix、refactor、product code 実装のために文脈を整理する時

## 知識範囲

- `catalog`、`summary`、`full` の圧縮粒度
- canonical source への重複寄せ
- `confirmed`、`inferred`、`gap` の分離
- active / completed plan から有効な過去判断だけを抽出すること
- downstream が読む順番の整理

## 圧縮方針

- 先に対象を機械的に棚卸しし、その後で判断する
- 読む粒度は `catalog`、`summary`、`full` の順で上げる
- 重複 instruction は正本だけを残し、重複元は path で退避する
- downstream の次判断に必要な情報だけを残す
- 出力ごとに `confirmed`、`inferred`、`gap` の状態を明示する

## Runtime Boundary

- binding: [distiller.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/distiller.toml)
- agent runtime: [distiller.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/distiller.toml)
- contract: [distiller.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/distiller/contracts/distiller.contract.json)
- allowed: repo 文脈を read-only で棚卸しし、必要最小限に圧縮する
- forbidden: product code / product test / docs 正本 / workflow 正本を変更しない
- tool policy: agent runtime の `allowed_write_paths` / `allowed_commands` に従う

## 標準パターン

1. caller goal、入口 artifact、lane owner を確認する。
2. docs、plan、skill、関連 file を path catalog として棚卸しする。
3. 正本、重複、任意参照、未確認事項を分類する。
4. plan 履歴は検索と object 抽出で扱い、有効な過去判断だけを残す。
5. 必要なものだけ `summary` または `full` に展開する。
6. facts、constraints、gaps、required_reading へ圧縮する。
7. 次に読むべき情報または停止理由を明示する。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は `distiller` agent contract に従う。

## Stop / Reroute

- active work plan や関連 docs が不足している場合は停止する。
- 重要な fact の根拠 path を確認できない場合は停止する。
- 主要な設計判断が未確定で事実整理だけでは前進しない場合は `implement_lane` へ戻す。
- 実装前の文脈整理が目的なら、Codex implementation lane [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-distill/SKILL.md) を使う前提で `implement_lane` へ戻す。

## Handoff

- handoff 先: `implement_lane`
- 渡す contract: [distiller.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/distiller/contracts/distiller.contract.json)
- 渡す scope: 次の設計または調査を判断するための圧縮済み facts と gaps

## DO / DON'T

DO:
- 確認済み事実と推測を分ける
- 重要な fact には根拠 path を付ける
- 過去判断は有効なものだけを残し、出典 line を付ける
- 必要な参照先を読む順番つきで返す

DON'T:
- 設計向けまたは調査向けの詳細観点を共通 skill に戻さない
- broad な repo tour をしない
- product code、product test、docs 正本の変更に進まない

## Checklist

- [distill-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/distill/references/checklists/distill-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は `distiller` agent contract が決める。

## References

- 圧縮判断: [compression-patterns.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/distill/references/compression-patterns.md)
- 設計向け観点: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/distill-design/SKILL.md)
- 調査向け観点: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/distill-investigate/SKILL.md)
- binding: [distiller.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/distiller.toml)
- agent contract: [distiller.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/distiller/contracts/distiller.contract.json)

## Maintenance

- `distill` は共通圧縮知識だけを持つ。
- 設計向け、調査向けの観点は focused skill に分ける。
- 長い例や判断表は references に分離する。
- 実装前 context の整理は Codex implementation lane [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-distill/SKILL.md) に残す。
