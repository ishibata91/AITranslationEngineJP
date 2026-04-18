---
name: gateguard
description: Codex + MCP 用の編集前 fact gate 知識 package。file mutation 前に確認すべき事実と阻止条件を提供する。
---

# GateGuard

## 目的

`gateguard` は知識 package である。
MCP file mutation や destructive command の前に、対象、根拠、影響範囲を確認するための判断基準を提供する。

Codex hook は MCP tool call を完全には遮断できない。
そのため MCP file mutation 前の direct-use gate として参照する。

## いつ参照するか

- `write_file`、`move_file` などの MCP file mutation の前
- destructive command や shell file mutation の前
- user 指示、対象 file、根拠 docs、rollback が曖昧な時

## 参照しない場合

- read-only の file read や search だけを行う時
- すでに agent contract が mutation を禁止している時
- user が明示的に作業を停止した時

## 知識範囲

- mutation 前の事実確認
- missing facts の切り分け
- pass / blocked の判断
- hook と MCP direct-use gate の責務差

## 原則

- user の現在指示を最優先で確認する
- 対象 file と変更理由を具体化する
- rollback 可能性と影響範囲を確認する
- self-review だけで pass にしない

## 標準パターン

1. action type と target を確認する。
2. user instruction と source of truth を確認する。
3. confirmed facts と missing facts を分ける。
4. impact と rollback を確認する。
5. pass / blocked と next step を返す。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は呼び出し元 agent contract に従う。

## DO / DON'T

DO:
- mutation 前に対象と根拠を固定する
- missing facts を隠さず blocked にする
- MCP と hook の境界を正確に扱う

DON'T:
- gateguard 自身で変更を実行しない
- MCP を経由しない file mutation を要求しない
- hook が MCP を完全に止めると説明しない

## Checklist

- [gateguard-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/gateguard/references/checklists/gateguard-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は呼び出し元 agent contract が決める。

## References

- script: [codex-mcp-gateguard.js](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/gateguard/scripts/codex-mcp-gateguard.js)

## Maintenance

- file mutation の実行権限を skill 本体へ置かない。
- hook の適用範囲と MCP direct-use gate を混同しない。
- checklist は短い mutation 前確認に保つ。
