---
name: diagramming-d2
description: D2 の図（.d2）を作成・分割・修正・検証し、`d2` コマンドで review 用の `.svg` を生成する。コードベースの境界、依存、フローを人間がレビューしやすい図にしたい時に使う。
---

# Diagramming D2

## ワークフロー

1. 図を独立した `.d2` ファイルにするか決める。
2. 1 枚で主題が混ざる時は、図を分割する。
3. `.d2` を書くか更新する。
4. `d2 validate <file>.d2` を実行する。
5. `d2 -t 201 <file>.d2 <file>.svg` を実行し、review 用の `.svg` を生成する。
6. validate や render の問題を先に直し、その後にレイアウトと読みやすさを整え、通るまで再実行する。

## ルール

- `.d2` を source of truth にする。
- `.svg` は `d2 -t 201` コマンドで生成する review 用成果物として必ず更新する。
- 1 つの図では 1 つの主題を扱い、無関係な境界やフローを混ぜない。
- class 図は `shape: class` を使い、属性、依存、振る舞いは member で表現する。
- class 図では node 名、member、edge だけで責務が読める状態を優先する。
- sequence 図の participant は論理名と actual name のラベルで役割を表す。
- robustness 図は shape、node 名、edge で boundary / control / entity 相当の役割が読めるようにする。
- ER 図は SQL テーブルとして扱い、`shape: sql_table` を使う。
- ER 図の列、型、制約は table member に寄せ、汎用 box や class shape で代用しない。
- ラベルは短く、意味が明確なものにする。
- 動詞、パラメータ、クラスや型のラベルは、可能な限り `論理名 (`actual-name`)` を同じラベル内に置く。
- 関数名や command 名だけを裸で並べず、何をするものかを論理名で先に読める形にする。
- パラメータや DTO も、値の役割が分かる論理名を付けてから actual name を添える。
- node の label では、まず「何者か」が分かる役割名を置く。
- edge label は短く保つが、依存違反 review ができなくなるほど削らない。
- node ID は役割が分かる名前にする。
- edge の向きは図の中で一貫させる。
- shape は図の意図に合わせるが、class 図は `class`、ER 図は `sql_table` を標準とする。
- 関連は通常 edge で表す。
- 集約は白 diamond の arrowhead で表す。
- 合成は黒 diamond の arrowhead で表す。
- diamond は集合度を読むための記号として使い、単なる装飾として使わない。
- overview 図では detail 図で説明すべき処理順や内部手順を詰め込みすぎない。
- overview 図でも依存違反 review に必要な node-to-node edge は残し、package-to-package まで削りすぎない。
- overview と detail では情報密度を分けてよい。
- 線のフォントサイズは22とする。
- package 内の node は、対応関係が追いやすい順に並べる。
- 線の交差は routing 変更より先に、分割・配置・ラベル圧縮で減らす。
- 大きい図は 1 枚に詰めず、複数ファイルに分ける。
- validate や render が失敗している状態で見た目調整を優先しない。

## 補足

- class 図の例:
  - `shape: class` を使い、フィールド、主要メソッド、必要な visibility を member 行へ置く。
  - クラス名は `論理名 (`actual-name`)` でまとめる。
- ER 図の例:
  - `shape: sql_table` を使い、列名、型、`primary_key`、`foreign_key`、`not_null` などの制約を member に載せる。
  - 中間 table も SQL table として表し、多対多を box 群でぼかさない。
- 関係記法の例:
  - 通常関連は通常 edge を使う。
  - 集約は白 diamond を使う。
  - 合成は黒 diamond を使う。
- `d2 validate` が通っても render 結果は必ずしも妥当ではないため、生成した `.svg` を目視確認する。
- 新しい layout / routing 構文は、既存図へ直接入れず最小例で検証してから使う。
- validate と render の両方が通っても、node サイズ、接続位置、label 混入が崩れていないか確認する。
- 直角線指定は常に正解ではなく、図全体の readability を下げるなら採用しない。
- 明示されない限り、D2 以外の図形式へ変換しない。
- review 対象は `.svg` でも、差分の正本は `.d2` のまま扱う。
