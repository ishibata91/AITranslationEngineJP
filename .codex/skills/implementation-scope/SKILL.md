---
name: implementation-scope
description: Codex 側の実装スコープ知識 package。human review 後に、人間が Codex implementation lane へ渡せる handoff packet を owned_scope、依存、検証単位へ分ける判断基準を提供する。
---

# Implementation Scope

## 目的

`implementation-scope` は知識 package である。
`designer` agent が human review 後に、Codex implementation handoff packet を固定するための、分割粒度、依存、validation、completion signal の見方を提供する。

実行境界、source of truth、handoff、stop / reroute は [design-bundle](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md) を参照する。

## いつ参照するか

- design bundle が human review 済みになった時
- 人間が `implement_lane` に渡せる owned_scope を作る時
- handoff ごとの depends_on と validation を固定する時

## 参照しない場合

- human review 前に実装 scope を決める時
- 承認済み implementation-scope なしで `implement_lane` の implementation execution へ handoff する時
- product code を直接実装する時
- 実装時の再現、trace、review 補助を扱う時

## 知識範囲

- `implementation-scope.md` の構成
- contract freeze の固定条件と証跡
- owned_scope、depends_on、validation_commands、completion_signal
- validation ownership gate
- parallel execution gate
- Codex implementation handoff packet の構成
- docs 正本化を handoff に混ぜない境界

## 原則

- human review 後にだけ作る
- 下流 handoff が依存する public seam は、implementation handoff より先に contract freeze として固定する
- 1 handoff は独立検証可能な粒度にする
- 1 handoff は原則として `1 受け入れユースケース × 1 validation intent` に収める
- 用語体系は `受け入れテスト > システムテスト > UI人間操作E2E / APIテスト` を正本にする
- `E2E` は UI 人間操作起点だけを指す
- `APIテスト` は public seam 起点の system-level test として扱う
- UI が入口の機能では、裏側の直接呼び出しや fixture 直接投入だけを `UI人間操作E2E` の完了条件にしない
- handoff が大きいかどうかは、論理境界に加えて想定 file 数と想定変更行数で判定する
- scope、依存、first_action、validation、done condition を必ず揃える
- 並列実行可能性は task 出し時に明示する
- human review 済みの詳細要求タイプと質問票回答だけを handoff source にする
- validation command は handoff の owned_scope と completion_signal だけで pass できるものにする
- backend と frontend は必ず別 handoff に分ける
- frontend handoff は contract freeze 済みの backend contract / DTO / gateway 境界に depends_on する
- 必要な場合だけ `本番経路` を notes に書き、必須 artifact や domain 固有欄にはしない
- `本番経路` は実行時に通る public API / DTO / controller / UI entry / persistence path を指す
- `本番経路` は domain 名や画面名の知識ではなく、handoff の補助語として扱う
- Codex は承認済み implementation-scope に基づいて handoff packet を作る
- Codex implementation lane に docs 正本化や workflow 変更を渡さない

## Handoff Split Rule

implementation-scope の handoff は、token 量の事前計算ではなく、論理境界と規模の目安で分割する。
1 handoff は原則として `1 受け入れユースケース × 1 validation intent` に収める。
受け入れユースケースは domain 名や画面名ではなく、人間または system が開始する処理単位として扱う。

受け入れユースケースとは、1 つの操作または system process が、永続化、backend contract、frontend state / UI まで必要範囲を通って成立し、完了後にその use case を原則として再編集しなくてよい単位である。
ただし implementation handoff では backend と frontend を同一 handoff に含めず、backend 側の contract / DTO / gateway 境界を `contract_freeze` として固定してから frontend 側を別 handoff として切る。
layer 単位の分割は、単体では完了判定できない中間状態を増やすため、最初の分割根拠にしない。
APIテストを先に固める場合は、public seam、request / response contract、外部入力開始、主要観測点を固定する。
UI人間操作E2Eをあとに固める場合は、開始操作、入力方法、主要操作列、UI-visible outcome を固定する。
裏側 API、service、fixture への直接投入は補助検証であり、UI人間操作E2E の完了判定にはしない。

## Contract Freeze Rule

contract freeze は、下流 handoff が参照してよい public seam を固定する段階である。
実装完了ではなく、依存先を増やしてよい境界が確定したことを意味する。

contract freeze として固定してよい対象:

- public API request / response
- DTO shape と field obligation
- gateway interface
- controller entry と error surface
- frontend が参照してよい state contract

contract freeze は次を満たす必要がある。

- downstream handoff の `depends_on` に書ける completion signal がある
- field 名、nullability、error shape、識別子、永続化 key などの境界差分が列挙されている
- `validation_commands` が public seam の固定を直接確認できる
- notes に freeze 根拠 artifact を書ける

contract freeze を固定できない場合は、frontend handoff や並列 handoff を開かない。
この場合は backend 側の探索または replan を優先し、見込み contract を completion signal にしない。

## Size Gate

handoff を作る前に、既存 code map、類似変更、owned_scope からおおよその touched files と changed lines を見積もる。
changed lines は、生成物、snapshot、lockfile、docs 正本化を除いた product code / product test の追加行と削除行の合計として扱う。

規模の目安:

- normal: `15 files` 以下、かつ `800 changed lines` 以下なら 1 受け入れユースケース handoff として扱える
- caution: `16-25 files` または `801-1500 changed lines` なら、completion_signal が 1 つに閉じ、検証 fixture が限定できる場合だけ 1 handoff にしてよい
- split required: `26 files` 以上、または `1501 changed lines` 以上が見込まれるなら、handoff 前に分割する
- hard stop: `40 files` 以上、または `2500 changed lines` 以上が見込まれるなら、1 handoff として渡さず、人間に replan 要求を返す

規模で分割する時は、次の順で切る。

1. 別 use case に分けられるなら use case で切る。
2. 同じ use case 内でも、contract freeze と backend persistence / implementation、frontend state / UI は必ず切る。
3. それでも大きい場合は、parse、preview、generation、settings save など failure mode が違う処理で切る。

## Boundary Rule

import、generation、settings save、preview、create / update / delete、export のように use case が違う処理は、同じ layer でも分割する。
failure mode が違う処理も、可能なら分割する。

同じ受け入れユースケースでも、backend と frontend は 1 handoff に含めない。
contract freeze handoff は backend 実装全体ではなく、public seam の固定だけを扱う。
backend handoff は永続化、service / usecase、controller、DTO / gateway 境界までを扱う。
frontend handoff は確定済み contract freeze に依存して state / UI を扱う。

backend 側の handoff に含めてよい layer:

- repository / SQLite concrete
- service / usecase
- controller / bootstrap
- gateway contract / DTO mapping

frontend 側の handoff に含めてよい layer:

- frontend state / presenter / usecase / controller
- frontend UI screen

ただし caution 以上の規模なら、`notes` に 1 handoff とする理由、想定 file 数、想定 changed lines、分割しない理由を書く。
実行時に通る経路が誤読されやすい場合は、`notes` に `本番経路` として public API / DTO / controller / UI entry / persistence path だけを書く。
特定 domain の処理名や業務知識は skill へ持ち込まず、task-local artifact 側へ置く。

禁止例:

- domain 名や画面名だけを根拠に複数 use case を同じ handoff にする
- normal / caution を超える規模なのに、file 数と changed lines の見積もりを書かずに 1 handoff にする
- backend contract と frontend UI を同じ handoff に含める
- contract freeze を置かずに frontend handoff を開始する
- migration、import、generation、settings save のような failure mode の違う処理を「同じ画面だから」という理由だけでまとめる

## Parallel Execution Rule

handoff 作成時は、まず `depends_on` から依存 DAG を作る。
次に、同じ段階で依存が解消できる handoff を `execution_group` にまとめる。
`execution_group` は `wave-1`、`wave-2`、`wave-3` のように必要な数だけ連番で作る。
`execution_group` は Codex implementation lane 側の ready wave であり、同じ wave 内でも `parallelizable_with` に列挙されない handoff は並列実行しない。
`ready_wave` は `execution_group` と同じ値を handoff ごとに明示し、Ready Waves 表で handoff 一覧、開始前に完了している依存、並列 pair、blocker を確認できる形にする。

並列実行可能な handoff は、次をすべて満たす必要がある。

- `depends_on` が空、または同一 group 開始前に完了済みである
- `owned_scope` の想定変更 file / module / test target が他 handoff と重ならない
- public contract、DTO、schema、migration、shared fixture などの shared boundary を同時に変更しない
- `validation_commands` が handoff-local で、失敗時に owner handoff を特定できる
- contract freeze 確定前の frontend handoff ではない
- 同じ broad gate 修正や同じ flaky environment blocker を解消対象にしない

並列不可の task は `parallel_blockers` に理由を書く。
理由は `depends_on`、`owned_scope_overlap`、`shared_contract_change`、`validation_owner_ambiguous`、`backend_frontend_order`、`broad_gate_shared` のいずれかに寄せる。
これ以外の理由が必要な場合は、task-local artifact 側に具体理由を書き、skill 側の共通分類は増やさない。

`execution_group: wave-1` は即実行可能な handoff を指す。
`execution_group: wave-N` は、`wave-1` から `wave-(N-1)` までのうち、その handoff の `depends_on` に必要な completion signal が完了した後に実行できる handoff を指す。
backend と frontend は別 handoff のまま維持し、frontend は contract freeze 完了後の wave に置く。
final validation、Sonar、Codex review は全 wave 完了後にだけ実行する。

## First Action Rule

handoff 作成時は、各 handoff に `first_action` を必ず書く。
`first_action` は Codex implementation lane が最初に閉じる 1 clause だけを示す。
広い調査開始、複数 clause、partial な advance は書かない。

`first_action` には次を含める。

- path
- symbol または対象単位
- 変更種別
- 対応する `completion_signal` clause
- 1 手目にする理由

1 edit で clause を閉じられない場合は、同じ clause の最小 closure chain を `notes` または `completion_signal` に補足する。
ただし複数 clause を 1 つの `first_action` にまとめない。

## Validation Ownership Gate

handoff 作成時は、各 `validation_commands` がその handoff の owner に属しているかを必ず確認する。
validation owner は、`owned_scope` の変更だけでその command を pass させられる handoff である。

各 command は次を満たす必要がある。

- `completion_signal` を直接検証している
- `owned_scope` と解消済み `depends_on` だけで pass できる
- 未実装の後続 handoff を前提にしない
- 失敗した時に、その handoff の実装不足として説明できる
- broad validation は原則 `final-validation-and-review` に寄せる

途中 handoff に broad validation を置く場合は、broad command が必要な理由、required downstream scope、分割しない理由を `notes` に書く。
この説明を書けない場合、その command は対象 handoff の validation ではなく final validation に移す。

## 標準パターン

1. human review status と approval record を確認する。
2. source artifact を列挙する。
3. scenario-design の詳細要求タイプに `needs_human_decision` が残っていないことを確認する。
4. use case ごとに downstream が依存する public seam を contract freeze として固定する。
5. handoff を受け入れユースケースと validation intent で分割する。
6. 各 handoff の想定 file 数と想定 changed lines を見積もる。
7. normal / caution / split required / hard stop を判定する。
8. 各 handoff に owned_scope、depends_on、first_action、validation_commands、completion_signal を書く。
   - validation_commands は validation ownership gate を通し、completion_signal を直接検証するものだけを残す。
   - frontend handoff は `UI人間操作E2E` を直接 owner にせず、final validation で証明する。
   - frontend handoff は対応する contract freeze handoff の completion signal に depends_on する。
9. 各 handoff に execution_group、ready_wave、parallelizable_with、parallel_blockers を書く。
   - Ready Waves 表で wave ごとの handoff、依存完了条件、並列 pair、blocker を確認できるようにする。
10. 人間が Codex implementation lane に渡す entry、禁止事項、期待される完了報告を明示する。
11. implementation completion後に正本化が必要なら、completion packet の residual risk に残して人間が次 action を判断できるようにする。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は `designer` agent contract に従う。

## DO / DON'T

DO:
- 承認済み artifact だけを source にする
- 承認済み詳細要求タイプを validation intent の根拠にする
- implementation handoff を受け入れユースケースで分ける
- downstream handoff が依存する public seam を contract freeze として先に固定する
- file 数と changed lines の目安で大きすぎる handoff を事前に切る
- layer をまたぐ時は、completion_signal が受け入れユースケースとして検証できるようにする
- UI人間操作E2E の証明は final validation lane に寄せる
- validation command と completion signal を揃える
- first_action を 1 completion_signal clause に固定する
- validation command が owned_scope と解消済み depends_on だけで pass できることを確認する
- 並列実行可能な handoff は execution_group、ready_wave、parallelizable_with で明示する
- Ready Waves 表で Codex implementation lane が読む実行順を先に固定する
- 並列不可の handoff は parallel_blockers に分類済み理由を書く
- broad validation を途中 handoff に置く場合は、必要な downstream scope と理由を notes に書く
- backend と frontend を同一 handoff に入れず、depends_on で接続する
- frontend handoff は contract freeze handoff の completion signal に接続する
- `本番経路` が必要な時だけ notes に補助情報として書く
- 人間がそのまま `implement_lane` に渡せる packet にする

DON'T:
- human review 前に owned_scope を確定しない
- `needs_human_decision` が残る scenario-design から handoff を作らない
- layer だけを根拠に、単体では完了判定できない micro handoff を量産しない
- backend と frontend を同一 handoff に含めない
- contract freeze が未完了のまま frontend handoff を開かない
- UI 入口の handoff で、裏側の直接呼び出しだけを完了条件にしない
- file 数と changed lines が split required を超える handoff を 1 件として渡さない
- first_action がない handoff を Codex implementation lane に渡さない
- first_action に複数 clause や曖昧な調査開始を書かない
- Codex から Codex implementation lane へ直接 handoff しない
- docs 正本化を Codex implementation handoff に含めない
- 未実装の後続 handoff を必要とする validation command を途中 handoff に入れない
- final validation で見るべき broad command を lane-local validation として扱わない
- owned_scope、shared boundary、validation owner が曖昧な handoff を並列実行可能として扱わない
- 同じ execution_group という理由だけで handoff を並列実行しない
- domain 固有知識を skill や template の共通例として増やさない
- implementation-time investigation は Codex implementation lane 内で閉じ、Codex replan 前提にしない

## Checklist

- [implementation-scope-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-scope/references/checklists/implementation-scope-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は `designer` agent contract が決める。

## References

- template: [implementation-scope.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/templates/task-folder/implementation-scope.md)
- Codex implementation lane entry: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-lane/SKILL.md)
- runtime skill: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md)
- agent contract: [designer.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/designer/contracts/designer.contract.json)

## Maintenance

- tool policy と output obligation を skill 本体へ戻さない。
- Codex implementation lane 実装 workflow の詳細は  [.codex/skills](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/) に置く。
- handoff 粒度の長い例は references に分離する。
