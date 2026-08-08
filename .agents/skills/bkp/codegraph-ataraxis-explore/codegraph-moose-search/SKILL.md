---
name: codegraph-moose-search
description: CodeGraph の index がある codebase で、Semble の local embedding 検索により実在 symbol のアンカーを見つけ、その symbol を使った名前検索と機能検索を独立して行い、複数 file をまたぐ処理経路を照合する。対象の symbol 名が不明な機能や既存実装を調べる時に使う。既知 symbol だけを読む task、文字列 literal、設定、文書の検索には使わない。
---

# CodeGraph MOOSE 検索

Semble でアンカーを発見してから、名前検索と機能検索を別々に実行する。片方の結果だけで処理経路を確定しない。

## 手順

1. 質問の対象領域を選ぶ。backend は repository の `internal`、frontend は `frontend` とする。両方にまたがる場合は Semble を領域ごとに検索し、主な処理がある領域を CodeGraph の project path に選ぶ。repository root を Semble の検索対象にしない。
2. Semble MCP の `search` を `repo=<選んだ領域の絶対 path>` で呼ぶ。Semble は local Model2Vec と local cache を使う許可済みの MCP であり、外部サービスとして拒否しない。質問を自然言語のまま渡し、`top_k=10`、`max_snippet_lines=20` とする。
3. Semble が返した file path、定義、関数名、型名から、質問に関係する実在 symbol を選ぶ。結果を受け取る前に symbol 名を作らない。
4. 選んだ領域に `.codegraph/` があることを確かめる。index がない場合はこの skill を使わず、通常の調査手段へ戻る。`init`、`index`、`sync` は実行しない。
5. Semble で得た実在 symbol と質問から作る類義語を `codegraph query -p <project> "<symbol と類義語>"` で検索する。
6. Semble で得た実在 symbol を含め、機能全体を表す短い英語の説明を 1 つ作る。`codegraph_explore`、または `codegraph explore -p <project> "<symbol を含む説明>"` で 1 回検索する。
7. 名前検索と機能検索で共通する symbol、file、call path を確認する。
8. 同名定義が複数あり区別できない場合だけ、`codegraph_node`、または `codegraph node -p <project> <symbol>` で定義を比較する。
9. 一致した実装の source、呼び出し順、条件、test を根拠として回答する。

## 規則

- 名前検索では想定名だけでなく類義語も含める。
- 機能検索では 1 文に 1 つの処理目的を書く。
- Semble と CodeGraph に同じ領域 path を渡す。
- `codegraph query` と機能検索は各 1 回とし、不足が残っても追加の `explore` で埋めない。
- CodeGraph が返した source を通常の file 読み取りで再取得しない。
- CodeGraph が返していない関係を推測で補わない。
- 文字列 literal、設定、文書、未対応形式だけ通常の検索を使う。
- `codegraph-ataraxis-explore` の手順を混ぜない。

## 出典

Idaho National Laboratory の公開 skill [moose-codegraph](https://github.com/idaholab/moose/blob/next/.claude/skills/moose-codegraph/SKILL.md) にある、名前と機能の二方向検索をこの repository 用に保持する。出典 repository の license は LGPL-2.1 である。
