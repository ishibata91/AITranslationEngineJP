# CodeGraph によるコード探索の評価

## 評価 1 で固定する内容

- 変更前の状態: Codex の LSP 連携と CodeGraph の両方を使わない。
- 適用するアイデア: `internal` と `frontend` を CodeGraph へ登録し、コード探索で CodeGraph を使える状態にする。
- 期待する効果: 単語置き換えと辞書に関する探索で、回答の根拠を保ったまま、使用 token 数と経過時間が減る。
- 効果がないと判断できる差: 使用 token 数または経過時間の差が変更前の反復測定の揺れに収まる場合、または増える場合。
- 予想する副作用: 関係する file、関数、呼び出し順、実装上の条件を取りこぼす可能性がある。回答の根拠とソースを照合して確認する。

## 標本

同じ質問を変更前と適用後で使う。

1. 登録済みの翻訳辞書が本文の単語置き換えに使われるまでの処理経路を調べ、関係する file、主要な関数、呼び出し順を根拠付きで示す。
2. 固有名の表記揺れが単語置き換えと本文中の言及検出でどのように扱われるかを調べ、実装とテストの根拠を示す。

1 は通常の利用範囲を代表する。2 は backend の複数 package をまたぐ探索を含む。

## 揃える条件

- 実行場所は repository root とする。
- Codex CLI は実行時に確認した同じ版を使う。
- model は `gpt-5.6-terra` とする。
- 推論設定、sandbox、入力文、作業領域、外部接続の可否を揃える。
- 各実行は履歴を引き継がない `fresh` とし、ソースを変更しない。
- 変更前と適用後は、それぞれ合計 3 回実行する。質問 1 は 2 回、質問 2 は 1 回とする。
- CodeGraph の有無以外に差が生じた場合は記録し、その差だけで支持または棄却しない。

## 測る値

- Codex CLI が記録する入力 token 数、出力 token 数、合計 token 数。
- command の開始から終了までの経過時間。
- tool 呼び出し回数。
- 回答に必要な file、関数、呼び出し順、実装上の条件が根拠付きで含まれる割合。
- 根拠とソースが食い違う記述の件数。

品質の確認項目は変更前の測定前にソースから作り、実行する `fresh` へ渡さない。Codex 本体が固定した確認項目で出力を照合する。

### 品質の確認項目

各項目を 1 点とし、各回答を 5 点満点で採点する。根拠となる file と関数を回答から確認できない項目は 0 点とする。ソースと食い違う記述の件数は別に数える。

質問 1 は次の 5 項目を確認する。

1. `Engine.LoadDictionary` が `translationVocabulary` を呼ぶ。
2. `translationVocabulary` が `ListMasterTerms` と `ListProperNouns` を読み、共通の stoplist を適用する。
3. `master_term`、`proper_noun` の順で `dictionary.Pair` を追加し、`dictionary.NewDictionary` を呼ぶ。
4. 翻訳処理が固有名抽出後に辞書を作り、本文翻訳へ渡す。
5. `prepareSource` が runtime tag を退避し、`Dictionary.Apply` を呼び、tag を戻した本文を prompt の組み立てへ渡す。

質問 2 は次の 5 項目を確認する。

1. `dictionary.Dictionary` と `mention.Detector` の照合は大文字と小文字を区別する。
2. 両方が単語境界と長い source を優先する正規表現を使う。
3. 両方が同一 source の先勝ちを使う。
4. 両方が `translationVocabulary` の同じ読み込み順と stoplist を使う。
5. `dictionary`、`mention`、`engine` の実装と test を根拠として示す。

### 実行条件

- Codex CLI は `codex-cli 0.146.0` とする。
- 推論設定は `medium` とする。
- sandbox は `read-only` とする。
- `--ephemeral` と `--json` を使う。
- Graphify、CodeGraph、LSP、MCP、Web 検索を使わない。
- 変更前と適用後の合計 6 回を直列に実行し、同時実行による経過時間への影響を避ける。

### CodeGraph の導入結果

- CodeGraph CLI は `1.5.0` を nvm の Node.js `v24.14.1` に global 導入した。
- Codex の global MCP 設定へ `codegraph serve --mcp` を登録した。
- `internal` を独立した project として初期化した。106 file、1,646 node、4,270 edge を索引した。
- `frontend` を独立した project として初期化した。87 file、880 node、2,357 edge を索引した。
- 各 `.codegraph/.gitignore` だけを repository に残し、machine ごとの database は commit 対象から外した。
- 適用後の測定は `/Users/iorishibata/Repositories/AITranslationEngineJP/internal` を `projectPath` に指定した。

## 判定の境界

- 変更前の質問 1 における 2 回の合計 token 数の差と経過時間の差を、同じ条件で生じる揺れとして使う。
- 適用後の中央値が、合計 token 数と経過時間の両方で変更前の中央値を下回り、減少幅が変更前の最大の揺れを超えた場合に期待する効果を確認する。
- 必要な根拠の割合が変更前を下回る場合、またはソースと食い違う記述が増える場合は副作用として扱う。
- 期待する効果があり、副作用が増えない場合は支持とする。
- 期待する効果がない場合、反対方向の差がある場合、または品質低下が明らかな場合は棄却とする。
- 効果と品質低下の両方がある場合は人間判断とする。
- 変更前の反復測定で揺れを定められない場合は判断不能とする。

境界の出どころは、同じ条件で行う変更前の反復測定である。値を見た後に判定式を変更しない。

## 終了条件

変更前と適用後を合計 3 回ずつ測り、使用 token 数、経過時間、tool 呼び出し回数、回答の根拠を記録した時点で評価 1 を終える。

## 結果

### 変更前

| 質問 | 回 | 入力 token | 出力 token | 合計 token | 経過時間 | tool 呼び出し | 品質 | 食い違い |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| 1 | 1 | 520,443 | 2,819 | 523,262 | 69.58 秒 | 8 | 5/5 | 0 |
| 1 | 2 | 434,943 | 2,757 | 437,700 | 68.36 秒 | 7 | 5/5 | 0 |
| 2 | 1 | 304,414 | 2,189 | 306,603 | 52.29 秒 | 7 | 5/5 | 0 |

質問 1 の中央値は合計 480,481 token、68.97 秒である。質問 1 の反復間の差は 85,562 token、1.22 秒である。3 回全体の中央値は合計 437,700 token、68.36 秒、tool 呼び出し 7 回である。

変更前の生データと回答は `/private/tmp/codegraph-exploration-benchmark/baseline/` に保存した。質問 2 は 1 回だけなので、質問 2 単独の揺れは算出できない。

### 適用後

| 質問 | 回 | 入力 token | 出力 token | 合計 token | 経過時間 | tool 呼び出し | 品質 | 食い違い |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| 1 | 1 | 322,573 | 2,107 | 324,680 | 51.41 秒 | 6 | 4/5 | 0 |
| 1 | 2 | 266,245 | 1,728 | 267,973 | 55.14 秒 | 5 | 5/5 | 0 |
| 2 | 1 | 252,236 | 2,211 | 254,447 | 54.13 秒 | 5 | 4/5 | 0 |

質問 1 の中央値は合計 296,327 token、53.28 秒である。変更前から合計 token 数は 38.3%、経過時間は 22.8%減った。減少幅は 184,155 token、15.70 秒であり、変更前の揺れである 85,562 token、1.22 秒を超えた。

3 回全体の中央値は合計 267,973 token、54.13 秒、tool 呼び出し 5 回である。変更前の全体中央値から合計 token 数は 38.8%、経過時間は 20.8%、tool 呼び出しは 28.6%減った。

質問 2 は合計 token 数が 17.0%減り、経過時間が 3.5%増えた。質問 2 は両状態とも 1 回だけなので、経過時間の差を揺れと分離できない。

品質は変更前の 15/15 から適用後の 13/15 へ下がった。適用後の質問 1 の 1 回目は、`master_term` を `proper_noun` より先に追加する順序を明示しなかった。適用後の質問 2 は、`engine/mention_test.go` の統合 test を根拠として示さなかった。ソースと食い違う記述はなかった。

適用後の回答は CodeGraph を使ったことを本文に記載したため、採点時に変更前と適用後の別を伏せられなかった。品質項目は測定前に固定したが、採点者の盲検化はできていない。

適用後の生データと回答は `/private/tmp/codegraph-exploration-benchmark/applied/` に保存した。

### 判定

合計 token 数と経過時間の全体中央値は減った。質問 1 の減少幅も変更前の揺れを超えた。一方、品質は 2 点下がった。固定した判定の境界に従い、評価 1 は人間判断とする。

## 評価 2: 公開 skill の検索手順

### 固定する内容

- 比較対象 1 は CodeGraph Project Reader の段階的な探索とする。`explore` で領域を把握し、`node` で特定 symbol を読み、必要な場合だけ `callers`、`callees` を使う。`maxFiles` は 3 から 8 とする。
- 比較対象 2 は MOOSE CodeGraph の二方向検索とする。名前と類義語を `query` で探し、機能の短い英語説明を `explore` へ渡し、両方の結果を照合する。曖昧な symbol だけ `node` で確認する。
- 比較対象 3 は Ataraxis Codebase Exploration の一括 query とする。関係する symbol を一つの `explore` へまとめ、返された source を再取得せず、blast radius と test の参照を使う。
- 3 種類を混ぜない。各 `fresh` には担当する 1 種類の手順だけを渡す。
- 期待する効果は、評価 1 と同じ品質を保ちながら、合計 token 数、経過時間、tool 呼び出し回数が減ることである。
- 予想する副作用は、検索回数を抑えた結果として、読み込み順、stoplist、runtime tag、test の根拠を取りこぼすことである。

### 標本と実行条件

質問は「登録済みの翻訳辞書が本文の単語置き換えに使われるまでの処理経路を調べ、関係する file、主要な関数、呼び出し順を根拠付きで示す。」とする。通常の利用範囲を代表し、複数 package、呼び出し経路、実装条件を含むため選んだ。

- model は `gpt-5.6-terra`、推論設定は `medium`、Codex CLI は `codex-cli 0.146.0` とする。
- CodeGraph CLI は `1.5.0` とし、`internal` の同じ index を使う。
- 各方式を履歴のない `fresh` で 1 回ずつ直列実行する。
- sandbox は `read-only` とし、`--ephemeral` と `--json` を使う。
- Graphify、LSP、Web 検索を使わない。通常の file 検索と読み取りは、担当する公開 skill が許す不足確認に限る。
- 3 種類は検索手順と、それに伴って選ばれる CodeGraph MCP または CLI のcommandだけが異なる。

### 測定と判定

入力 token 数、出力 token 数、合計 token 数、経過時間、tool 呼び出し回数を測る。品質は評価 1 の質問 1 と同じ 5 項目で採点し、ソースと食い違う記述を別に数える。

品質が 5/5 で食い違いがなく、合計 token 数と経過時間の両方が他の品質 5/5 の方式を下回る方式を暫定優位とする。token 数と経過時間が反対方向へ動く場合、または品質が下がる場合は人間判断とする。各方式 1 回なので、暫定優位は再現性を示さず、差の確証が必要な場合は同じ条件で反復する。

### 終了条件

3 種類を各 1 回測り、効率、品質、食い違いを記録した時点で評価 2 を終える。

### 結果

| 方式 | 入力 token | 出力 token | 合計 token | 経過時間 | tool 呼び出し | 品質 | 食い違い |
|---|---:|---:|---:|---:|---:|---:|---:|
| CodeGraph Project Reader | 467,792 | 2,608 | 470,400 | 71.08 秒 | 12 | 4/5 | 0 |
| MOOSE CodeGraph | 224,463 | 2,505 | 226,968 | 57.10 秒 | 6 | 5/5 | 0 |
| Ataraxis Codebase Exploration | 93,566 | 1,644 | 95,210 | 37.48 秒 | 1 | 3/5 | 0 |

CodeGraph Project Reader は、同じ最初の探索を CLI と MCP で重複して実行した。複数 symbol を一度に `codegraph node` へ渡して失敗した後、symbol ごとに 9 回の `node` を実行した。`master_term` を `proper_noun` より先に追加する順序を回答へ明示しなかった。

MOOSE CodeGraph は、3 回の `query` と 3 回の `explore` を CLI で実行した。名前検索と機能検索の結果を照合し、5 項目をすべて回答した。

Ataraxis Codebase Exploration は、1 回の `codegraph_explore` だけを実行した。合計 token 数、経過時間、tool 呼び出し回数は 3 方式で最小だった。一方、固有名の確定後に辞書を作って本文翻訳へ渡す順序と、`prepareSource` による runtime tag の退避と復元を回答へ含めなかった。

MOOSE CodeGraph は CodeGraph Project Reader より、合計 token 数が 51.8%、経過時間が 19.7%、tool 呼び出し回数が 50.0%少なく、品質は 1 点高かった。

Ataraxis Codebase Exploration は MOOSE CodeGraph より、合計 token 数が 58.1%、経過時間が 34.4%、tool 呼び出し回数が 83.3%少なかった。品質は 2 点低かった。

生データと回答は `/private/tmp/codegraph-exploration-benchmark/public-skills/` に保存した。

### 判定

純粋な探索効率は Ataraxis Codebase Exploration が最も高かった。品質を保った方式は MOOSE CodeGraph だけだった。最大の効率と品質低下が同時に確認されたため、評価 2 は人間判断とする。各方式 1 回の結果なので、差の再現性は未確認である。

## 評価 3: Semble によるアンカー発見との組み合わせ

### 固定する内容

- 比較対象 1 は CodeGraph を使わず、Semble でアンカーを発見した後に通常の file 読み取りで処理経路を確認する方式とする。
- 比較対象 2 は Semble でアンカーを発見した後、`codegraph-ataraxis-explore` の一括 query で処理経路を確認する方式とする。
- 比較対象 3 は Semble でアンカーを発見した後、`codegraph-moose-search` の名前検索と機能検索を照合して処理経路を確認する方式とする。
- 各方式で Semble の最初の query を必須とし、Semble の結果を受け取る前に symbol 名を推測して検索しない。
- 3 種類を混ぜない。各 `fresh` には担当する 1 種類の手順だけを渡す。
- 期待する効果は、アンカーを事前に渡さなくても評価 2 と同じ品質を保ち、合計 token 数と経過時間を抑えることである。
- 予想する副作用は、Semble の上位候補に引かれ、読み込み順、stoplist、runtime tag、test の根拠を取りこぼすことである。

### 標本と実行条件

質問は「登録済みの翻訳辞書が本文の単語置き換えに使われるまでの処理経路を調べ、関係する file、主要な関数、呼び出し順を根拠付きで示す。」とする。質問に file 名、symbol 名、package 名は含めない。

- model は `gpt-5.6-terra`、推論設定は `medium`、Codex CLI は `codex-cli 0.146.0` とする。
- Semble CLI と MCP は `0.5.4` とし、同じ repository の自動更新される索引を使う。
- CodeGraph CLI は `1.5.0` とし、CodeGraph を使う 2 方式は `internal` の同じ index を使う。
- 各方式を履歴のない `fresh` で 1 回ずつ直列実行する。
- sandbox は `read-only` とし、`--ephemeral` と `--json` を使う。
- LSP、Graphify、Web 検索を使わない。
- 比較対象 1 は CodeGraph を使わない。比較対象 2 と 3 は指定した CodeGraph skill 以外の探索手順を混ぜない。

### 測定と判定

入力 token 数、出力 token 数、合計 token 数、経過時間、tool 呼び出し回数を測る。品質は評価 1 の質問 1 と同じ 5 項目で採点し、ソースと食い違う記述を別に数える。

品質が 5/5 で食い違いがなく、合計 token 数と経過時間の両方が他の品質 5/5 の方式を下回る方式を暫定優位とする。token 数と経過時間が反対方向へ動く場合、または品質が下がる場合は人間判断とする。各方式 1 回なので、暫定優位は再現性を示さない。

### 終了条件

3 種類を各 1 回測り、効率、品質、食い違いを記録した時点で評価 3 を終える。

### 結果

| 方式 | 入力 token | 出力 token | 合計 token | 経過時間 | tool 呼び出し | 品質 | 食い違い |
|---|---:|---:|---:|---:|---:|---:|---:|
| CodeGraph なし + Semble | 290,556 | 2,078 | 292,634 | 100.63 秒 | 6 | 5/5 | 0 |
| Ataraxis + Semble | 230,970 | 1,803 | 232,773 | 51.79 秒 | 5 | 3/5 | 0 |
| MOOSE + Semble | 183,248 | 2,036 | 185,284 | 67.80 秒 | 5 | 0/5 | 0 |

CodeGraph なし + Semble は、Semble の `search` を 1 回、`find_related` を 2 回実行した後、必要な source を通常の file 読み取りで確認した。固定した 5 項目をすべて回答した。

Ataraxis + Semble は、Semble の `search` を 1 回実行した後、`codegraph_explore` を 1 回実行した。`LoadDictionary`、`translationVocabulary`、読み込み順、stoplist、`NewDictionary` は回答した。固有名確定後に辞書を作って本文へ渡す順序と、`prepareSource` による runtime tag の退避・復元は回答に含まれなかった。

MOOSE + Semble は、Semble の `search` を 1 回実行した後、`codegraph query` と `codegraph explore` を各 1 回実行した。`NewDictionary` と `Dictionary.Apply` の内部条件は回答したが、`LoadDictionary`、`translationVocabulary`、読み込み順、stoplist、固有名確定後の本文経路、`prepareSource` を確認できなかった。

Ataraxis + Semble は CodeGraph なし + Semble より、合計 token 数が 20.5%、経過時間が 48.5%少なかった。品質は 2 点低かった。

MOOSE + Semble は CodeGraph なし + Semble より、合計 token 数が 36.7%、経過時間が 32.6%少なかった。品質は 5 点低かった。

MOOSE + Semble は合計 token 数が最小だった。Ataraxis + Semble は経過時間が最短で、MOOSE + Semble より品質が 3 点高かった。

測定前の起動確認で、Ataraxis は `projectPath` の指定漏れにより repository root に index がないと判断して終了した。この結果は方式を実行していないため除外し、固定済みの `internal` を明示して有効な 1 回を測った。

MOOSE は最初の実行で `codegraph explore` を 6 回呼び、固定した手順に違反した。この結果は除外した。違反した実行は合計 507,346 token、90.75 秒、品質 5/5 であり、追加探索によって品質を回復できる一方、今回固定した MOOSE 方式ではないことを示す。次の実行は Semble を外部サービスと誤認して拒否したため除外した。Semble が利用者の許可したローカル MCP であることを明記した後に、有効な 1 回を測った。

生データと回答は `/private/tmp/codegraph-exploration-benchmark/semble-combinations/results/` に保存した。除外した起動も `invalid-*` の名前で残した。

### 判定

品質を維持した方式は CodeGraph なし + Semble だけだった。Ataraxis + Semble は速度を大きく減らしたが品質が 2 点下がった。固定回数の MOOSE + Semble は token 数が最小だったが、処理経路の根拠を取得できず品質が 0/5 だった。効率と品質低下が同時に確認されたため、評価 3 は人間判断とする。各方式 1 回の結果なので、差の再現性は未確認である。

## 評価 4: 領域を限定した Semble アンカー発見

### 固定する内容

- 評価 3 は Semble に repository root を渡したため、`internal` 以外の候補が検索結果へ混ざった。
- 適用するアイデアは、backend の質問では Semble の `repo` と CodeGraph の `projectPath` の両方を repository の `internal` に固定することである。
- CodeGraph なし + Semble、Ataraxis + Semble、MOOSE + Semble の 3 方式を各 1 回測り直す。
- Ataraxis と MOOSE は、Semble を local Model2Vec と local cache を使う許可済み MCP として必ず呼ぶ。Semble の結果を受け取る前に symbol 名を作らない。
- Ataraxis は Semble で得た実在 symbol を 1 回の `codegraph_explore` へまとめる。
- MOOSE は Semble で得た実在 symbol を使い、`codegraph query` と機能を表す `codegraph explore` を各 1 回実行して照合する。
- 期待する効果は、評価 3 より関連するアンカーが上位へ入り、CodeGraph を使う方式の品質が上がることである。
- 予想する副作用は、領域を限定しても 1 回または 2 回の CodeGraph 検索では呼び出し順と実装条件を取得できないことである。

### 標本と実行条件

質問、品質の 5 項目、model、推論設定、Codex CLI、sandbox、禁止する道具は評価 3 と同じにする。各方式を履歴のない `fresh` で 1 回ずつ直列実行する。

- Semble の `repo` は `/Users/iorishibata/Repositories/AITranslationEngineJP/internal` とする。
- CodeGraph の `projectPath` は `/Users/iorishibata/Repositories/AITranslationEngineJP/internal` とする。
- 測定前に `internal` 専用の Semble 索引を作り、3 方式が同じ作成済み索引から開始する。
- CodeGraph なし + Semble は Semble の追加検索と通常の file 読み取りを許可する。
- Ataraxis と MOOSE は CodeGraph が返した source を通常の file 読み取りで再取得しない。

### 測定と判定

入力 token 数、出力 token 数、合計 token 数、経過時間、tool 呼び出し回数を測る。品質と食い違いは評価 3 と同じ数え方を使う。

評価 3 の同じ方式より品質が上がり、合計 token 数または経過時間の増加と同時に起きた場合は人間判断とする。品質が上がらない場合は領域限定の効果を棄却する。各方式 1 回なので、確認した差は再現性を示さない。

### 終了条件

3 種類を各 1 回測り、評価 3 との差、品質、食い違いを記録した時点で評価 4 を終える。

### 結果

| 方式 | 入力 token | 出力 token | 合計 token | 経過時間 | tool 呼び出し | 品質 | 食い違い |
|---|---:|---:|---:|---:|---:|---:|---:|
| CodeGraph なし + Semble | 396,493 | 2,474 | 398,967 | 74.33 秒 | 7 | 5/5 | 0 |
| Ataraxis + Semble | 186,046 | 2,050 | 188,096 | 58.95 秒 | 3 | 1/5 | 0 |
| MOOSE + Semble | 199,618 | 2,055 | 201,673 | 54.19 秒 | 5 | 3/5 | 0 |

3 方式とも Semble の `repo` に `/Users/iorishibata/Repositories/AITranslationEngineJP/internal` を渡した。Ataraxis と MOOSE は CodeGraph にも同じ path を渡した。

CodeGraph なし + Semble は、辞書の読み込み順、stoplist、固有名確定後の辞書作成、`prepareSource` による runtime tag の退避・復元をすべて回答した。評価 3 より合計 token 数は 36.3%増え、経過時間は 26.1%減り、品質は 5/5 のままだった。

Ataraxis + Semble は、固有名を本文より先に確定し、部分形を派生してから辞書を作り、本文処理へ渡す順序を回答した。`LoadDictionary` が `translationVocabulary` を呼ぶこと、両供給元と stoplist、追加順、`prepareSource` と `Dictionary.Apply` は回答に含まれなかった。評価 3 より合計 token 数は 19.2%減り、経過時間は 13.8%増え、品質は 3/5 から 1/5 へ下がった。

MOOSE + Semble は、`LoadDictionary`、`translationVocabulary`、両供給元と stoplist、追加順、`NewDictionary`、`Dictionary.Apply` を回答した。固有名確定後に辞書を作って本文処理へ渡す順序と、`prepareSource` による runtime tag の退避・復元は回答に含まれなかった。評価 3 より合計 token 数は 8.8%増え、経過時間は 20.1%減り、品質は 0/5 から 3/5 へ上がった。

Ataraxis の Semble 検索は `Engine.Run` と固有名処理を主なアンカーとして選び、MOOSE の Semble 検索は `LoadDictionary`、`Dictionary.Apply`、`runtimetag.Mask` をアンカーとして選んだ。領域を同じにしても、1 回の Semble 結果から選ぶアンカーと CodeGraph query の組み立てによって取得範囲が変わった。

生データと回答は `/private/tmp/codegraph-exploration-benchmark/semble-scoped/results/` に保存した。

### 判定

領域限定は MOOSE の品質を 3 点上げ、経過時間を減らした。一方、合計 token 数は増え、品質は 5/5 に届かなかった。Ataraxis は合計 token 数を減らしたが、品質が 2 点下がった。領域限定の効果と方式ごとの品質低下が同時に確認されたため、評価 4 は人間判断とする。各方式 1 回の結果なので、差の再現性は未確認である。

## 評価 5: as-is 仕様のシーケンス図

### 固定する内容

- 評価 4 までの品質基準は実装項目を個別に数えたため、利用者が求める処理の物語を再構成できるかを直接測っていなかった。
- 評価対象を「辞書置き換えの as-is 仕様を Mermaid シーケンス図にする」へ変更する。
- 比較対象は MOOSE + Semble と Ataraxis + Semble の 2 方式とする。CodeGraph なし + Semble は比較しない。
- 期待する効果は、アンカーを事前に渡さなくても、辞書の準備から本文置換と翻訳までの呼び出し順とデータの受け渡しを図として再構成できることである。
- 予想する副作用は、検索結果に含まれた一部の symbol だけを並べ、前後の呼び出しまたは分岐がつながらない図になることである。

### 標本と実行条件

依頼は「辞書置き換えの as-is 仕様を Mermaid のシーケンス図にする。参加者、呼び出し順、辞書データの受け渡し、本文へ置換する前後、処理を通らない分岐を根拠付きで示す。」とする。file 名、package 名、symbol 名は渡さない。

- model は `gpt-5.6-terra`、推論設定は `medium`、Codex CLI は `codex-cli 0.146.0` とする。
- Semble の `repo` と CodeGraph の `projectPath` は `/Users/iorishibata/Repositories/AITranslationEngineJP/internal` とする。
- 各方式を履歴のない `fresh` で 1 回ずつ直列実行する。
- sandbox は `read-only` とし、`--ephemeral` と `--json` を使う。
- LSP、Graphify、Web 検索を使わない。
- MOOSE と Ataraxis は各 skill の検索回数と再読禁止を守る。

### 品質の確認項目

各項目を 1 点とし、5 点満点で採点する。図または根拠説明から順序を追えない項目は 0 点とする。

1. 実行入口から本文処理までの呼び出しを示す。
2. 固有名の確定と派生を終えてから本文用辞書を作る順序を示す。
3. 登録語彙の読み込み、stoplist、`master_term` 優先、辞書構築までを一続きで示す。
4. 本文処理から prompt 組み立て、runtime tag の退避、辞書置換、tag 復元、翻訳器呼び出しまでを一続きで示す。
5. 参照訳が完全一致する場合は辞書置換と翻訳器呼び出しを通らない分岐を示す。

図に source と食い違う順序、呼び出し、条件がある場合は別に件数を数える。Mermaid の `sequenceDiagram` として構文を読めない場合は品質を 0 点とする。

### 測定と判定

入力 token 数、出力 token 数、合計 token 数、経過時間、tool 呼び出し回数、品質、食い違いを測る。

品質が高い方式を暫定優位とする。品質が同じ場合は、合計 token 数と経過時間の両方が少ない方式を暫定優位とする。token 数と経過時間が反対方向へ動く場合は人間判断とする。各方式 1 回なので、暫定優位は再現性を示さない。

### 終了条件

2 方式を各 1 回測り、図、効率、品質、食い違いを記録した時点で評価 5 を終える。

### 結果

| 方式 | 入力 token | 出力 token | 合計 token | 経過時間 | tool 呼び出し | 品質 | 食い違い |
|---|---:|---:|---:|---:|---:|---:|---:|
| MOOSE + Semble | 177,812 | 2,015 | 179,827 | 63.20 秒 | 4 | 2/5 | 0 |
| Ataraxis + Semble | 182,383 | 2,107 | 184,490 | 54.44 秒 | 4 | 0/5 | 0 |

MOOSE + Semble は、`LoadDictionary` から登録語彙の読み込み、stoplist、`NewDictionary` までを一続きで示した。参照訳が完全一致する場合に `composeBodyPrompt` を通らない分岐も示した。実行入口から本文処理までの接続、固有名の確定と派生、`composeBodyPrompt` から `Dictionary.Apply` へ至る呼び出しは取得できなかった。`Mask` と `Restore` は根拠へ挙げたが、実際の呼び出し順を図へ含めなかった。

Ataraxis + Semble は、`ListProperNouns` と固有名翻訳の分岐を示した。辞書置き換えの処理経路として必要な実行入口、固有名確定後の辞書構築、登録語彙と stoplist、本文置換、参照訳の分岐は示さなかった。Semble の上位候補が辞書置き換えの中心から外れ、一括した `codegraph_explore` も固有名処理を中心に取得した。

Ataraxis + Semble は MOOSE + Semble より 8.76 秒短かった。一方、合計 token 数は 4,663 token 多く、品質は 2 点低かった。tool 呼び出し回数は同じだった。

Ataraxis の最初の実行は、Semble MCP がローカル検索であることを承認処理が判定できず、`search` を拒否したため除外した。Semble `0.5.4` は local Model2Vec、BM25、local cache を使い、API key と外部検索サービスを必要としない。利用者が `search` を tool 単位で許可した後に、有効な 1 回を測った。

生データと回答は `/private/tmp/codegraph-exploration-benchmark/sequence-story/results/` に保存した。除外した実行は `ataraxis-invalid-refusal`、過去の結果を参照した MOOSE の実行は `moose-invalid-memory` の名前で残した。

### 判定

MOOSE + Semble は品質で Ataraxis + Semble を上回ったため、固定した判定に従うと暫定優位である。ただし、MOOSE + Semble も 2/5 であり、辞書置き換えの as-is 仕様を一続きのシーケンス図として再構成できなかった。両方式とも目的を満たさないため、評価 5 だけを根拠に採用する方式は決めない。各方式 1 回の結果なので、差の再現性は未確認である。

## 評価 6: rg で入口を特定する CodeGraph skill

### 固定する内容

- 適用するアイデアは、`rg`を1回だけ使って実在symbolの入口を特定し、CodeGraphを原則1回、不足が1点だけ残る場合は追加1回使う`codegraph-rg-explore`である。
- 期待する効果は、`rg`後に複数回行っていたfile読み取りをCodeGraphの一括取得へ置き換え、評価5の品質項目を保ちながら合計token数、経過時間、tool呼び出し回数を減らすことである。
- 予想する副作用は、最初の`rg`で選んだsymbolにCodeGraphの取得範囲が偏り、実行入口、辞書構築、本文置換、参照訳の分岐のいずれかを取りこぼすことである。
- Semble、MOOSE、Ataraxisは使わない。

### skillの動作確認

依頼は評価5と同じas-is仕様のシーケンス図とする。modelは`gpt-5.6-terra`、推論設定は`medium`、Codex CLIは`codex-cli 0.146.0`とする。`--disable memories`、`--ephemeral`、`--json`、`read-only`を使い、履歴と過去の回答を渡さない。

動作確認は1回だけ実行する。品質は評価5の5項目で採点する。1回の結果はskillが呼び出し上限を守って必要な道具を呼べるかだけを確認し、CodeGraphの効果の支持または棄却には使わない。

### 確証の比較条件

動作確認後の比較では、同じ依頼、model、推論設定、Codex CLI、sandbox、履歴無効化を揃える。変更前は`rg`後に通常のコード調査規約で関係を追跡し、適用後は`codegraph-rg-explore`で関係を追跡する。各状態を3回測り、合計token数、経過時間、tool呼び出し回数、品質、食い違いを記録する。

品質5/5かつ食い違い0件を保った適用後の中央値が、変更前の中央値より合計token数と経過時間の両方で少ない場合に支持する。品質が下がる場合、または合計token数と経過時間の両方が減らない場合は棄却する。片方だけが減る場合は人間判断とする。

### 動作確認の結果

| skill版 | 入力token | 出力token | 合計token | 経過時間 | tool呼び出し | 品質 | 食い違い |
|---|---:|---:|---:|---:|---:|---:|---:|
| 初版 | 128,699 | 2,341 | 131,040 | 67.42秒 | 3 | 0/5 | 0 |
| 近接検索版 | 118,190 | 2,097 | 120,287 | 53.25秒 | 3 | 3/5 | 0 |

初版は`rg`で依頼文の`辞書置き換え`を完全一致検索して0件になった後、実在symbolを持たずにCodeGraphを呼んだ。CodeGraphは`backend`という語から秘密情報保存処理を取得した。初版はskillの入口取得規約が不足したため無効な動作確認とする。

近接検索版は`辞書`と`置換または適用`を同じ行から検索し、`submitBodyBatch`、`LoadDictionary`、`Dictionary.Apply`を最初のCodeGraphへ渡した。辞書構築と参照訳による迂回は取得できた。2回目も`codegraph explore`を使ったため、`composeBodyPrompt`からruntime tagの退避、辞書適用、復元へ進む実装を取得できなかった。

近接検索版の結果に基づき、追加1回を最初の結果に現れたsymbolの`node`、`callers`、`callees`だけに限定する。変更後のskillは別の動作確認として測る。

`node`限定版は入力116,431 token、出力1,975 token、合計118,406 token、44.46秒、tool呼び出し3回、品質2/5、食い違い0件だった。`rg`1回、`codegraph explore`1回、`codegraph node`1回の上限を守った。追加取得先に外側の`planBodyRequests`を選び、未確認だった`composeBodyPrompt`内のデータ変換を取得しなかったため、近接検索版より品質が1点下がった。

`node`限定版の結果に基づき、最初のアンカーを入口、データ構築、実際の変換、処理を通らない条件から分散して選び、追加取得では確認済みの外側の順序より未確認のデータ変換を優先するようskillを変更する。

アンカー分散版は入力71,666 token、出力875 token、合計72,541 token、23.59秒、tool呼び出し2回、品質0/5、食い違い0件だった。terraは`dictionary`と`replaceまたはapply`だけを`rg`へ渡した。英語symbolと日本語コメントが同じ行にないため0件となり、skillの停止条件に従ってCodeGraphを呼ばなかった。

アンカー分散版の結果に基づき、1回の`rg`へ依頼文の言語とsourceのsymbolで使う英語の両方を含めるようskillを変更する。検索回数とCodeGraphの呼び出し上限は変更しない。

日英検索版は入力140,079 token、出力2,310 token、合計142,389 token、56.21秒、tool呼び出し4回、品質2/5、食い違い0件だった。`rg`1回、`codegraph explore`1回、`codegraph node`1回の上限を守った。`rg`は`prepareSource`の説明コメントへ一致したが、一致行だけでは数行下の関数定義を実在symbolとして取得できなかった。terraは実際の変換段をアンカーへ含めず、追加取得先にも外側の`planBodyRequests`を選んだ。

日英検索版の結果に基づき、1回の`rg`で一致行の前後3行以内を取得し、コメントまたは呼び出しを近接する定義へ結び付けるようskillを変更する。検索回数とCodeGraphの呼び出し上限は変更しない。

文脈取得版は入力240,115 token、出力2,088 token、合計242,203 token、52.42秒、tool呼び出し4回、品質2/5、食い違い0件だった。`Engine.Run`から辞書構築、`prepareSource`によるruntime tagの退避、`Dictionary.Apply`、tag復元までを取得した。固有名の確定と派生、本文処理から翻訳器までの接続、参照訳による迂回は取得できなかった。

文脈取得版は`rg -C 3`がtestを含む広い文脈を返し、合計token数を押し上げた。複雑な処理経路を品質5/5で再構成するには、CodeGraph 2回では入口側と本文変換側の両方を補えない。次版は`rg`をproduction sourceと前後2行へ限定し、CodeGraphを最初の`explore`1回と、異なる処理段の`node`最大2回へ変更する。

3回上限版は入力302,906 token、出力2,305 token、合計305,211 token、59.70秒、tool呼び出し5回、品質0/5、食い違い0件だった。`LoadDictionary`と`prepareSource`を追加取得したが、参照訳の完全一致分岐と翻訳器呼び出しを持つ本文処理を取得しなかった。`rg`の前後2行では、説明コメントの3行後にある`translateNarrations`の定義を実在symbolとして取得できなかった。

3回上限版の結果に基づき、`rg`の文脈を前1行と後3行へ変更する。一致した説明コメントの直後にある定義を、コメント中で言及されたsymbolより優先してアンカーへ選ぶ。CodeGraphの呼び出し上限は変更しない。

定義取得版は入力206,520 token、出力3,011 token、合計209,531 token、67.85秒、tool呼び出し5回、品質4/5、食い違い0件だった。`rg`1回、`codegraph explore`1回、`codegraph node`2回を使った。登録語彙、stoplist、`master_term`優先、辞書構築、本文の参照訳による迂回、prompt組み立て、runtime tagの退避、辞書置換、tag復元、翻訳器呼び出しを一続きで示した。固有名の確定と派生を終えてから辞書を作る順序は示さなかった。

履歴なし通常探索の1回は合計347,652 token、66.54秒、tool呼び出し9回、品質4/5だった。参考比較では、定義取得版は品質を保ち、合計token数が39.7%、tool呼び出しが44.4%減った。経過時間は1.31秒、2.0%増えた。両方式は各1回であり、通常探索には1回のtool呼び出しで1操作だけを行う制限があるため、この差だけでskillを支持または棄却しない。

`codegraph-rg-explore`の動作確認は、`rg`1回とCodeGraph最大3回の上限を守り、品質4/5で終了した。確証の比較は、固定した最終版と同じ通常探索を各3回測るまで未完了とする。
