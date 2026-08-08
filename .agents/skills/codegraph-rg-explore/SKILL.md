---
name: codegraph-rg-explore
description: CodeGraph の index がある codebase で、rg を 1 回だけ使って実在 symbol の入口を見つけ、その symbol を起点に CodeGraph を原則 1 回、異なる処理段の根拠が不足する場合は追加 2 回まで使って処理経路を調べる。symbol 名が不明な機能について、コメントまたは定義から入口を特定した後の file 読み取りと tool 呼び出しを減らしたい時に使う。文字列 literal、設定、文書、参照元の全件確認には使わない。
---

# CodeGraph と rg による探索

`rg`で入口を特定し、CodeGraphで入口の前後にある呼び出し関係とsourceをまとめて取得する。

## 入力

- 調査対象: 人間が依頼した機能、処理経路、責務、または依存関係。
- 対象領域: backendはrepositoryの`internal`、frontendは`frontend`。人間がpathを指定した場合は指定を優先する。
- CodeGraph project: 対象領域の`.codegraph/`を持つ絶対path。

## 探索

1. 対象領域に`.codegraph/`があることを確かめる。indexがない場合はCodeGraphを呼ばず、通常の調査手段へ戻る。
2. 依頼文から、対象を表す異なる概念の語幹を2個選ぶ。各概念には依頼文の言語と、対象言語のsymbolで使われる一般的な英語の両方を含める。完全一致する文字列を探す依頼でない限り、依頼文全体を完全一致の検索語にしない。
3. 2個の概念が同じ行で近接する検索式を、出現順の両方を含めて作る。例えば「辞書」と「置換または適用」なら`((辞書|dictionary).{0,30}(置換|適用|replace|apply)|(置換|適用|replace|apply).{0,30}(辞書|dictionary))`とする。依頼文の言語または英語の片方だけへ限定しない。一般的な単語を単独で検索しない。
4. `rg`を対象領域に1回だけ実行する。対象言語のproduction sourceを検索し、test、生成物、`.codegraph/`を除く。最初から既知のsymbol名を検索語にしない。一致した説明コメントを直後の定義へ結び付けるため、前1行と後3行の文脈を同じ`rg`で取得する。
5. `rg`が0件の場合はCodeGraphを呼ばず、入口を取得できなかったと回答する。
6. `rg`の結果に現れた定義、呼び出し、またはコメントから、調査対象に直接関係する実在symbolを1個から4個選ぶ。一致した説明コメントの直後に定義がある場合は、その定義をコメント中で言及されただけのsymbolより優先する。検索語と無関係な同じfileのsymbolを選ばない。
7. 処理経路を調べる場合は、依頼された観点を入口、データ構築、実際の変換、処理を通らない条件に分け、`rg`の結果に存在する異なる段のsymbolを優先する。同じ責務のsymbolだけで4個を埋めない。
8. queryの先頭へ実在symbolを並べ、その後へ調査目的を1つ添える。`codegraph_explore`を`projectPath=<CodeGraph project>`、`maxFiles=8`で1回呼ぶ。MCPがない場合だけ`codegraph explore -p <CodeGraph project> "<query>" --max-files 8`を使う。
9. CodeGraphが返したsource、relationship map、caller、callee、testの参照を調査対象と照合する。
10. 最初のCodeGraph結果に調査対象と直接関係するsourceがある場合は、依頼された観点のうち根拠がない処理段を数える。入口側の分岐と本文のデータ変換の両方が不足する場合は、それぞれに最も近い実在symbolを1個ずつ選ぶ。同じ処理段から複数のsymbolを選ばない。
11. 追加呼び出しは`codegraph_node`または`codegraph node -p <CodeGraph project> <symbol>`を使い、選んだsymbolのsourceとcallerまたはcalleeの関係を同時に取得する。入口側または処理先側の関係だけが不足し、sourceが既に揃っている場合だけ`codegraph_callers`、`codegraph_callees`、`codegraph callers`、`codegraph callees`の対応する1つを使う。
12. 追加呼び出しでは`codegraph_explore`と`codegraph explore`を使わない。
13. 最初のCodeGraph結果が調査対象と無関係な場合は、無関係なsymbolで2回目を呼ばず、不足として回答する。
14. 取得した根拠から回答し、取得できなかった関係を不足として示す。

## 呼び出し上限

- Codexは`rg`を1回まで使う。
- CodexはCodeGraphを合計3回まで使う。
- Codexは最初のCodeGraph結果で必要な根拠が揃った場合、追加呼び出しを使わない。
- Codexは言い換えたqueryで`codegraph_explore`を繰り返さない。
- Codexは追加のCodeGraph呼び出しで機能検索を繰り返さない。
- CodexはCodeGraphが返したsourceを通常のfile読み取りで再取得しない。

## 根拠の扱い

- CodexはCodeGraphが返した実在symbolだけを追加呼び出しへ渡す。
- CodexはCodeGraphが返していない呼び出し関係を推測で補わない。
- Codexは参照元の全件確認が必要な場合、AGENTS.mdに従って`lsp_find_references`を使い、CodeGraphの結果だけで全件と確定しない。
- Codexはinterface実装先または呼び出し関係について、CodeGraphの結果が全件であると保証しない。

## 扱わない対象

- 文字列literal、md、json、設定file、環境変数、directory構造の検索。
- 参照元の全件確認、名前変更、diagnostics、型情報の確認。
- CodeGraphの`init`、`index`、`sync`によるindexの作成または更新。
- Semble、`codegraph-ataraxis-explore`、`codegraph-moose-search`との混用。

## 出力

- 処理経路: 根拠を取得できた呼び出し順とデータの受け渡し。
- 根拠: file pathとsymbol。
- 不足: 呼び出し上限内で取得できなかった関係または条件。
- 使用回数: `rg`とCodeGraphを呼び出した回数。

## 完了

Codexは、調査対象へ回答できる根拠を取得した時、または呼び出し上限へ達して不足を特定した時に探索を終える。
