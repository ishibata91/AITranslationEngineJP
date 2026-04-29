---
name: implementation-scope
description: Codex 側の実装スコープ作業プロトコル。人間レビュー 後に、人間が Codex implementation レーン へ渡せる 引き継ぎ入力 を 承認済み実装範囲、依存、検証単位へ分ける判断基準を提供する。
---
# Implementation Scope

## 目的

`implementation-scope` は作業プロトコルである。
`designer` agent が 人間レビュー 後に、Codex implementation 引き継ぎ入力 を固定するための、分割粒度、依存、検証、完了条件 の見方を提供する。

実行境界、正本、引き継ぎ、stop / 戻し は [design-bundle](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md) を参照する。

## 対応ロール

- `designer` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- 担当成果物は `implementation-scope` の出力規約で固定する。

## 入力規約

- design bundle が 人間レビュー 済みになった時
- 人間が `implement_lane` に渡せる 承認済み実装範囲 を作る時
- 引き継ぎ ごとの 依存対象 と 検証 を固定する時
- 入力に 根拠参照、担当者、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- エージェント実行定義とツール権限は [designer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.toml) の 書き込み許可 / 実行許可 とする。
- 雛形: [implementation-scope.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/templates/task-folder/implementation-scope.md)
- Codex implementation レーン 入口: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-lane/SKILL.md)
- 実行定義 skill: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md)
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

### 拘束観点

- `implementation-scope.md` の構成
- 契約固定 の固定条件と証跡
- 承認済み実装範囲、依存対象、検証コマンド、完了条件
- 検証 ownership 判定条件
- parallel execution 判定条件
- Codex implementation 引き継ぎ入力 の構成
- docs 正本化を 引き継ぎ に混ぜない境界

### 引き継ぎ分割規約

implementation-scope の 引き継ぎ は、token 量の事前計算ではなく、論理境界と規模の目安で分割する。
1 引き継ぎ は原則として `1 受け入れユースケース × 1 検証 intent` に収める。
受け入れユースケースは domain 名や画面名ではなく、人間または system が開始する処理単位として扱う。

受け入れユースケースとは、1 つの操作または system process が、永続化、backend 契約、frontend 状態 / UI まで必要範囲を通って成立し、完了後にその use case を原則として再編集しなくてよい単位である。
ただし implementation 引き継ぎ では backend と frontend を同一 引き継ぎ に含めず、backend 側の 契約 / DTO / gateway 境界を `contract_freeze` として固定してから frontend 側を別 引き継ぎ として切る。
層 単位の分割は、単体では完了判定できない中間状態を増やすため、最初の分割根拠にしない。
APIテストを先に固める場合は、公開接点、要求 / 応答契約、外部入力開始、主要観測点を固定する。
UI人間操作E2Eをあとに固める場合は、開始操作、入力方法、主要操作列、UI-visible 結果 を固定する。
裏側 API、service、検証データ への直接投入は補助検証であり、UI人間操作E2E の完了判定にはしない。

### 契約固定規約

契約固定 は、下流 引き継ぎ が参照してよい 公開接点 を固定する段階である。
実装完了ではなく、依存先を増やしてよい境界が確定したことを意味する。

契約固定 として固定してよい対象:

- public API request / response
- DTO 形 と field obligation
- gateway interface
- controller 入口 と エラー 表面
- frontend が参照してよい 状態 契約

契約固定 は次を満たす必要がある。

- downstream 引き継ぎ の `依存対象` に書ける 完了条件 がある
- field 名、null 許容、エラー 形、識別子、永続化 key などの境界差分が列挙されている
- `検証コマンド` が 公開接点 の固定を直接確認できる
- 補足 に freeze 根拠 成果物 を書ける

契約固定 を固定できない場合は、frontend 引き継ぎ や並列 引き継ぎ を開かない。
この場合は backend 側の探索または replan を優先し、見込み 契約 を 完了条件 にしない。

### 規模判定条件

引き継ぎ を作る前に、既存 code map、類似変更、承認済み実装範囲 からおおよその touched files と changed lines を見積もる。
changed lines は、生成物、スナップショット、lockfile、docs 正本化を除いた プロダクトコード / プロダクトテストの追加行と削除行の合計として扱う。

規模の目安:

- 通常: `15 files` 以下、かつ `800 changed lines` 以下なら 1 受け入れユースケース 引き継ぎ として扱える
- 注意: `16-25 files` または `801-1500 changed lines` なら、完了条件 が 1 つに閉じ、検証 検証データ が限定できる場合だけ 1 引き継ぎ にしてよい
- 分割必須: `26 files` 以上、または `1501 changed lines` 以上が見込まれるなら、引き継ぎ 前に分割する
- 強制停止: `40 files` 以上、または `2500 changed lines` 以上が見込まれるなら、1 引き継ぎ として渡さず、人間に replan 要求を返す

規模で分割する時は、次の順で切る。

1. 別 use case に分けられるなら use case で切る。
2. 同じ use case 内でも、契約固定 と backend persistence / implementation、frontend 状態 / UI は必ず切る。
3. それでも大きい場合は、parse、preview、generation、settings save など 失敗種別 が違う処理で切る。

### 境界規約

import、generation、settings save、preview、create / update / delete、export のように use case が違う処理は、同じ 層 でも分割する。
失敗種別 が違う処理も、可能なら分割する。

同じ受け入れユースケースでも、backend と frontend は 1 引き継ぎ に含めない。
契約固定 引き継ぎ は backend 実装全体ではなく、公開接点 の固定だけを扱う。
backend 引き継ぎ は永続化、service / usecase、controller、DTO / gateway 境界までを扱う。
frontend 引き継ぎ は確定済み 契約固定 に依存して 状態 / UI を扱う。

backend 側の 引き継ぎ に含めてよい 層:

- repository / SQLite concrete
- service / usecase
- controller / bootstrap
- gateway 契約 / DTO mapping

frontend 側の 引き継ぎ に含めてよい 層:

- frontend 状態 / presenter / usecase / controller
- frontend UI screen

ただし 注意 以上の規模なら、`補足` に 1 引き継ぎ とする理由、想定 file 数、想定 changed lines、分割しない理由を書く。
実行時に通る経路が誤読されやすい場合は、`補足` に `本番経路` として public API / DTO / controller / UI 入口 / persistence path だけを書く。
特定 domain の処理名や業務知識は skill へ持ち込まず、task 内成果物 側へ置く。

禁止例:

- domain 名や画面名だけを根拠に複数 use case を同じ 引き継ぎ にする
- 通常 / 注意 を超える規模なのに、file 数と changed lines の見積もりを書かずに 1 引き継ぎ にする
- backend 契約 と frontend UI を同じ 引き継ぎ に含める
- 契約固定 を置かずに frontend 引き継ぎ を開始する
- migration、import、generation、settings save のような 失敗種別 の違う処理を「同じ画面だから」という理由だけでまとめる

### 並列実行規約

引き継ぎ 作成時は、まず `依存対象` から依存 DAG を作る。
次に、同じ段階で依存が解消できる 引き継ぎ を `実行グループ` にまとめる。
`実行グループ` は `wave-1`、`wave-2`、`wave-3` のように必要な数だけ連番で作る。
`実行グループ` は Codex implementation レーン 側の 着手可能 wave であり、同じ wave 内でも `並列可能対象` に列挙されない 引き継ぎ は並列実行しない。
`ready_wave` は `実行グループ` と同じ値を 引き継ぎ ごとに明示し、着手可能 wave 表で 引き継ぎ 一覧、開始前に完了している依存、並列 pair、阻害要因 を確認できる形にする。

並列実行可能な 引き継ぎ は、次をすべて満たす必要がある。

- `依存対象` が空、または同一 group 開始前に完了済みである
- `承認済み実装範囲` の想定変更 file / module / test 対象 が他 引き継ぎ と重ならない
- public 契約、DTO、schema、migration、shared 検証データ などの shared 境界 を同時に変更しない
- `検証コマンド` が 引き継ぎ内 で、失敗時に 担当引き継ぎ を特定できる
- 契約固定 確定前の frontend 引き継ぎ ではない
- 同じ 広域 判定条件 修正や同じ flaky environment 阻害要因 を解消対象にしない

並列不可の task は `並列不可理由` に理由を書く。
理由は `依存対象`、`承認済み実装範囲重複`、`共有契約変更`、`検証担当不明`、`バックエンドフロントエンド順序`、`広域判定条件共有` のいずれかに寄せる。
これ以外の理由が必要な場合は、task 内成果物 側に具体理由を書き、skill 側の共通分類は増やさない。

`実行グループ: wave-1` は即実行可能な 引き継ぎ を指す。
`実行グループ: wave-N` は、`wave-1` から `wave-(N-1)` までのうち、その 引き継ぎ の `依存対象` に必要な 完了条件 が完了した後に実行できる 引き継ぎ を指す。
backend と frontend は別 引き継ぎ のまま維持し、frontend は 契約固定 完了後の wave に置く。
最終検証、Sonar、Codex レビュー は全 wave 完了後にだけ実行する。

### 初手規約

引き継ぎ 作成時は、各 引き継ぎ に `初手` を必ず書く。
`初手` は Codex implementation レーン が最初に閉じる 1 clause だけを示す。
広い調査開始、複数 clause、partial な advance は書かない。

`初手` には次を含める。

- path
- symbol または対象単位
- 変更種別
- 対応する `完了条件` clause
- 1 手目にする理由

1 edit で clause を閉じられない場合は、同じ clause の最小 closure chain を `補足` または `完了条件` に補足する。
ただし複数 clause を 1 つの `初手` にまとめない。

### 検証担当者判定条件

引き継ぎ 作成時は、各 `検証コマンド` がその 引き継ぎ の 担当者 に属しているかを必ず確認する。
検証 担当者 は、`承認済み実装範囲` の変更だけでその コマンド を 通過 させられる 引き継ぎ である。

各 コマンド は次を満たす必要がある。

- `完了条件` を直接検証している
- `承認済み実装範囲` と解消済み `依存対象` だけで 通過 できる
- 未実装の後続 引き継ぎ を前提にしない
- 失敗した時に、その 引き継ぎ の実装不足として説明できる
- 広域 検証 は原則 `最終検証とレビュー` に寄せる

途中 引き継ぎ に 広域 検証 を置く場合は、広域 コマンド が必要な理由、必須 downstream 対象範囲、分割しない理由を `補足` に書く。
この説明を書けない場合、その コマンド は対象 引き継ぎ の 検証 ではなく 最終検証 に移す。

## 判断規約

- 人間レビュー 後にだけ作る
- 下流 引き継ぎ が依存する 公開接点 は、implementation 引き継ぎ より先に 契約固定 として固定する
- 1 引き継ぎ は独立検証可能な粒度にする
- 1 引き継ぎ は原則として `1 受け入れユースケース × 1 検証 intent` に収める
- 用語体系は `受け入れテスト > システムテスト > UI人間操作E2E / APIテスト` を正本にする
- `E2E` は UI 人間操作起点だけを指す
- `APIテスト` は 公開接点 起点の system-level test として扱う
- UI が入口の機能では、裏側の直接呼び出しや 検証データ 直接投入だけを `UI人間操作E2E` の完了条件にしない
- 引き継ぎ が大きいかどうかは、論理境界に加えて想定 file 数と想定変更行数で判定する
- 対象範囲、依存、初手、検証、完了条件 を必ず揃える
- 並列実行可能性は task 出し時に明示する
- 人間レビュー 済みの詳細要求タイプと質問票回答だけを 引き継ぎ source にする
- 検証コマンド は 引き継ぎ の 承認済み実装範囲 と 完了条件 だけで 通過 できるものにする
- backend と frontend は必ず別 引き継ぎ に分ける
- frontend 引き継ぎ は 契約固定 済みの backend 契約 / DTO / gateway 境界に 依存対象 する
- 必要な場合だけ `本番経路` を 補足 に書き、必須 成果物 や domain 固有欄にはしない
- `本番経路` は実行時に通る public API / DTO / controller / UI 入口 / persistence path を指す
- `本番経路` は domain 名や画面名の知識ではなく、引き継ぎ の補助語として扱う
- Codex は承認済み implementation-scope に基づいて 引き継ぎ入力 を作る
- Codex implementation レーン に docs 正本化や 作業流れ 変更を渡さない

- 承認済み 成果物 だけを source にする
- 承認済み詳細要求タイプを 検証 intent の根拠にする
- implementation 引き継ぎ を受け入れユースケースで分ける
- downstream 引き継ぎ が依存する 公開接点 を 契約固定 として先に固定する
- file 数と changed lines の目安で大きすぎる 引き継ぎ を事前に切る
- 層 をまたぐ時は、完了条件 が受け入れユースケースとして検証できるようにする
- UI人間操作E2E の証明は 最終検証 レーン に寄せる
- 検証コマンド と 完了条件 を揃える
- 初手 を 1 完了条件 clause に固定する
- 検証コマンド が 承認済み実装範囲 と解消済み 依存対象 だけで 通過 できることを確認する
- 並列実行可能な 引き継ぎ は 実行グループ、ready_wave、並列可能対象 で明示する
- 着手可能 wave 表で Codex implementation レーン が読む実行順を先に固定する
- 並列不可の 引き継ぎ は 並列不可理由 に分類済み理由を書く
- 広域 検証 を途中 引き継ぎ に置く場合は、必要な downstream 対象範囲 と理由を 補足 に書く
- backend と frontend を同一 引き継ぎ に入れず、依存対象 で接続する
- frontend 引き継ぎ は 契約固定 引き継ぎ の 完了条件 に接続する
- `本番経路` が必要な時だけ 補足 に補助情報として書く
- 人間がそのまま `implement_lane` に渡せる 入力一式 にする

## 出力規約

- 出力は判断結果、根拠参照、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

## 完了規約

- task 内成果物 が承認状態、根拠参照、未決事項を含んでいる。
- 人間レビュー が必要な判断を AI だけで完了扱いにしていない。
- 人間レビュー 承認 を確認した。
- scenario-design に `needs_human_decision` が残っていないことを確認した。
- 承認済み詳細要求タイプを 検証 intent の根拠にした。
- 引き継ぎ を 承認済み実装範囲、依存対象、検証 で分けた。
- 各 引き継ぎ が `1 受け入れユースケース × 1 検証 intent` に収まっている。
- 各 検証コマンド が `完了条件` を直接検証している。
- 各 検証コマンド が `承認済み実装範囲` と解消済み `依存対象` だけで 通過 できる。
- 各 引き継ぎ に 1 clause だけを閉じる `初手` を書いた。
- 各 引き継ぎ の想定 touched files と changed lines を見積もった。
- `15 files` / `800 changed lines` 以下を 通常 として扱った。
- `16-25 files` または `801-1500 changed lines` の 注意 引き継ぎ には、1 件にする理由を `補足` に書いた。
- `26 files` 以上または `1501 changed lines` 以上の 分割必須 引き継ぎ を 1 件として渡していない。
- `40 files` 以上または `2500 changed lines` 以上の 強制停止 引き継ぎ は implement-lane へ戻した。
- import / generation / settings save / preview / create / update / delete / export のうち、別 use case になっている処理を同一 引き継ぎ に混ぜていない。
- domain 名や画面名だけを根拠に、複数 use case を同一 引き継ぎ にまとめていない。
- 層 をまたぐ 引き継ぎ は、受け入れユースケース 完了条件 で完了判定できる。
- frontend 引き継ぎ は `UI人間操作E2E` を直接 担当者 にせず、最終検証 で証明する形にした。
- `依存対象` から依存 DAG を作り、着手可能 wave を `実行グループ` と `ready_wave` にした。
- 着手可能 wave 表に 引き継ぎ、開始前依存、並列 pair、阻害要因 を書いた。
- 並列可能な 引き継ぎ だけを `並列可能対象` に列挙した。
- 並列不可の理由を `並列不可理由` に分類済み reason で書いた。
- 並列可能な 引き継ぎ の 承認済み実装範囲、shared 境界、検証 担当者 が重なっていない。
- 広域 検証 を途中 引き継ぎ に置く場合は、必須 downstream 対象範囲 と理由を `補足` に書いた。
- 人間が Codex implementation レーン に渡す 入口、禁止事項、期待完了報告を明示した。

## 停止規約

- 人間レビュー 前に実装 対象範囲 を決める時
- 承認済み implementation-scope なしで `implement_lane` の 実装実行 へ 引き継ぎ する時
- プロダクトコード を直接実装する時
- 実装時の再現、trace、レビュー 補助を扱う時
- 人間レビュー 前に 承認済み実装範囲 を確定しない
- `needs_human_decision` が残る scenario-design から 引き継ぎ を作らない
- 層 だけを根拠に、単体では完了判定できない micro 引き継ぎ を量産しない
- backend と frontend を同一 引き継ぎ に含めない
- 契約固定 が未完了のまま frontend 引き継ぎ を開かない
- UI 入口の 引き継ぎ で、裏側の直接呼び出しだけを完了条件にしない
- file 数と changed lines が 分割必須 を超える 引き継ぎ を 1 件として渡さない
- 初手 がない 引き継ぎ を Codex implementation レーン に渡さない
- 初手 に複数 clause や曖昧な調査開始を書かない
- Codex から Codex implementation レーン へ直接 引き継ぎ しない
- docs 正本化を Codex implementation 引き継ぎ に含めない
- 未実装の後続 引き継ぎ を必要とする 検証コマンド を途中 引き継ぎ に入れない
- 最終検証 で見るべき 広域 コマンド を レーン内検証 として扱わない
- 承認済み実装範囲、shared 境界、検証 担当者 が曖昧な 引き継ぎ を並列実行可能として扱わない
- 同じ 実行グループ という理由だけで 引き継ぎ を並列実行しない
- domain 固有知識を skill や 雛形 の共通例として増やさない
- implementation-time investigation は Codex implementation レーン 内で閉じ、Codex replan 前提にしない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- 人間レビュー 前に implementation-scope を作らなかった場合は停止する。
- 人間判断が残る scenario-design から implementation-scope を作らなかった場合は停止する。
- 層 だけを根拠に micro 引き継ぎ を量産しなかった場合は停止する。
- file 数と changed lines の基準を超える 引き継ぎ を根拠なしに残さなかった場合は停止する。
- `初手` がない 引き継ぎ を残さなかった場合は停止する。
- Codex から Codex implementation レーン へ直接 引き継ぎ しなかった場合は停止する。
- docs 正本化を Codex implementation レーン 引き継ぎ に混ぜなかった場合は停止する。
- 検証コマンド なしで 引き継ぎ しなかった場合は停止する。
- UI 入口の 引き継ぎ で、裏側の直接呼び出しだけを完了条件にしなかった場合は停止する。
- 未実装の後続 引き継ぎ を必要とする 検証コマンド を途中 引き継ぎ に入れなかった場合は停止する。
- 最終検証 で見るべき 広域 コマンド を レーン内検証 として扱わなかった場合は停止する。
- 同じ `実行グループ` という理由だけで並列実行可能として扱わなかった場合は停止する。
