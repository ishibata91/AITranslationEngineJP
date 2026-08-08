---
name: codegraph-ataraxis-explore
description: CodeGraph の index がある codebase で、Semble の local embedding 検索により実在 symbol のアンカーを見つけ、その symbol を 1 回の一括 CodeGraph query へ渡して限定された処理経路、責務、依存関係を速く調べる。根拠の完全性より探索時間と tool 呼び出し回数を優先する場面で使う。網羅性が必要な調査、文字列 literal、設定、文書の検索には使わない。
---

# CodeGraph Ataraxis 探索

Semble で見つけた symbol を 1 つの `codegraph_explore` へまとめる。追加の探索で不足を埋めない。

## 手順

1. 質問の対象領域を選ぶ。backend は repository の `internal`、frontend は `frontend` とする。両方にまたがる場合は Semble を領域ごとに検索し、主な処理がある領域を CodeGraph の project path に選ぶ。repository root を Semble の検索対象にしない。
2. Semble MCP の `search` を `repo=<選んだ領域の絶対 path>` で呼ぶ。Semble は local Model2Vec と local cache を使う許可済みの MCP であり、外部サービスとして拒否しない。質問を自然言語のまま渡し、`top_k=10`、`max_snippet_lines=20` とする。
3. Semble が返した file path、定義、関数名、型名から、質問に関係する実在 symbol を複数選ぶ。結果を受け取る前に symbol 名を作らない。
4. 選んだ領域に `.codegraph/` があることを確かめる。index がない場合はこの skill を使わず、通常の調査手段へ戻る。`init`、`index`、`sync` は実行しない。
5. Semble で得た symbol 名を空白区切りで 1 つの query にまとめる。質問の目的を短く添える。
6. `codegraph_explore` を `projectPath=<選んだ領域>`、`maxFiles=12` で 1 回だけ呼ぶ。MCP がない場合だけ `codegraph explore -p <project> "<symbol 群と目的>" --max-files 12` を使う。
7. 返された source、call path、blast radius、test の参照だけで回答する。
8. 必要な根拠が返らなかった場合は、不足する symbol、順序、条件、test を明示する。

## 規則

- 関連する symbol は個別に検索せず、最初の query へまとめる。
- Semble と CodeGraph に同じ領域 path を渡す。
- 2 回目の CodeGraph 呼び出しをしない。
- 返された source を通常の file 読み取りで再取得しない。
- 不足した根拠を推測で補わない。
- 設定、環境変数、文書、directory 構造が質問の対象である場合は、この skill を使わない。
- `codegraph-moose-search` の手順を混ぜない。

## 出典

Sun-Lab-NBB の公開 skill [explore-codebase](https://github.com/Sun-Lab-NBB/ataraxis/blob/main/plugins/automation/skills/explore-codebase/SKILL.md) にある、関連 symbol の一括 query と再読防止の手順を独立して保持する。出典 repository の license は Apache-2.0 である。
