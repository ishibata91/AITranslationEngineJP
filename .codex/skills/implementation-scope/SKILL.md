---
name: implementation-scope
description: Codex 側の実装スコープ知識 package。human review 後に、人間が Copilot へ渡せる handoff packet を owned_scope、依存、検証単位へ分ける判断基準を提供する。
---

# Implementation Scope

## 目的

`implementation-scope` は知識 package である。
`designer` agent が human review 後に、人間向け Copilot handoff packet を固定するための、分割粒度、依存、validation、completion signal の見方を提供する。

実行権限、agent contract、handoff、stop / reroute は [designer.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.agent.md) が持つ。

## いつ参照するか

- design bundle が human review 済みになった時
- 人間が Copilot に渡せる owned_scope を作る時
- handoff ごとの depends_on と validation を固定する時

## 参照しない場合

- human review 前に実装 scope を決める時
- Codex から Copilot へ直接 handoff する時
- product code を直接実装する時
- 実装時の再現、trace、review 補助を扱う時

## 知識範囲

- `implementation-scope.md` の構成
- owned_scope、depends_on、validation_commands、completion_signal
- validation ownership gate
- 人間向け Copilot handoff packet の構成
- docs 正本化を handoff に混ぜない境界

## 原則

- human review 後にだけ作る
- 1 handoff は独立検証可能な粒度にする
- 1 handoff は原則として `1 e2e use case × 1 validation intent` に収める
- handoff が大きいかどうかは、論理境界に加えて想定 file 数と想定変更行数で判定する
- scope、依存、validation、done condition を必ず揃える
- validation command は handoff の owned_scope と completion_signal だけで pass できるものにする
- Codex は Copilot へ直接渡さず、人間へ handoff packet を返す
- Copilot に docs 正本化や workflow 変更を渡さない

## Handoff Split Rule

implementation-scope の handoff は、token 量の事前計算ではなく、論理境界と規模の目安で分割する。
1 handoff は原則として `1 e2e use case × 1 validation intent` に収める。
use case は domain 名や画面名ではなく、人間または system が開始する処理単位として扱う。

e2e use case とは、1 つの操作または system process が、永続化、backend contract、frontend state / UI まで必要範囲を通って成立し、完了後にその use case を原則として再編集しなくてよい単位である。
layer 単位の分割は、単体では完了判定できない中間状態を増やすため、最初の分割根拠にしない。

## Size Gate

handoff を作る前に、既存 code map、類似変更、owned_scope からおおよその touched files と changed lines を見積もる。
changed lines は、生成物、snapshot、lockfile、docs 正本化を除いた product code / product test の追加行と削除行の合計として扱う。

規模の目安:

- normal: `15 files` 以下、かつ `800 changed lines` 以下なら 1 e2e use case handoff として扱える
- caution: `16-25 files` または `801-1500 changed lines` なら、completion_signal が 1 つに閉じ、検証 fixture が限定できる場合だけ 1 handoff にしてよい
- split required: `26 files` 以上、または `1501 changed lines` 以上が見込まれるなら、handoff 前に分割する
- hard stop: `40 files` 以上、または `2500 changed lines` 以上が見込まれるなら、1 handoff として渡さず propose-plans へ戻す

規模で分割する時は、次の順で切る。

1. 別 use case に分けられるなら use case で切る。
2. 同じ use case 内なら、backend persistence / contract と frontend state / UI の安定した contract 境界で切る。
3. それでも大きい場合は、parse、preview、generation、settings save など failure mode が違う処理で切る。

## Boundary Rule

import、generation、settings save、preview、create / update / delete、export のように use case が違う処理は、同じ layer でも分割する。
failure mode が違う処理も、可能なら分割する。

同じ e2e use case を成立させるために必要なら、次の layer を 1 handoff に含めてよい。

- repository / SQLite concrete
- service / usecase
- controller / bootstrap
- frontend gateway contract / DTO mapping
- frontend state / presenter / usecase / controller
- frontend UI screen

ただし caution 以上の規模なら、`notes` に 1 handoff とする理由、想定 file 数、想定 changed lines、分割しない理由を書く。

禁止例:

- domain 名や画面名だけを根拠に複数 use case を同じ handoff にする
- normal / caution を超える規模なのに、file 数と changed lines の見積もりを書かずに 1 handoff にする
- backend contract と frontend UI を同じ handoff に含めたのに、completion_signal が UI または API 動作として検証できない
- migration、import、generation、settings save のような failure mode の違う処理を「同じ画面だから」という理由だけでまとめる

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
3. handoff を e2e use case と validation intent で分割する。
4. 各 handoff の想定 file 数と想定 changed lines を見積もる。
5. normal / caution / split required / hard stop を判定する。
6. 各 handoff に owned_scope、depends_on、validation_commands、completion_signal を書く。
   - validation_commands は validation ownership gate を通し、completion_signal を直接検証するものだけを残す。
7. 人間が Copilot に渡す entry、禁止事項、期待される完了報告を明示する。
8. Copilot 修正完了後に正本化が必要なら `propose_plans` へ戻す前提を残す。

この手順は知識上の標準例である。
実行順、必須 input、完了条件は `designer` agent contract に従う。

## DO / DON'T

DO:
- 承認済み artifact だけを source にする
- implementation handoff を e2e use case で分ける
- file 数と changed lines の目安で大きすぎる handoff を事前に切る
- layer をまたぐ時は、completion_signal が e2e 動作として検証できるようにする
- validation command と completion signal を揃える
- validation command が owned_scope と解消済み depends_on だけで pass できることを確認する
- broad validation を途中 handoff に置く場合は、必要な downstream scope と理由を notes に書く
- 人間がそのまま Copilot に渡せる packet にする

DON'T:
- human review 前に owned_scope を確定しない
- layer だけを根拠に、単体では完了判定できない micro handoff を量産しない
- file 数と changed lines が split required を超える handoff を 1 件として渡さない
- Codex から Copilot へ直接 handoff しない
- docs 正本化を Copilot handoff に含めない
- 未実装の後続 handoff を必要とする validation command を途中 handoff に入れない
- final validation で見るべき broad command を lane-local validation として扱わない
- implementation-time investigation を Codex 側へ戻さない

## Checklist

- [implementation-scope-checklist.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implementation-scope/references/checklists/implementation-scope-checklist.md) を参照する。
- checklist は知識確認用であり、実行義務は `designer` agent contract が決める。

## References

- template: [implementation-scope.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/exec-plans/templates/task-folder/implementation-scope.md)
- Copilot entry: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/implementation-orchestrate/SKILL.md)
- agent spec: [designer.agent.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/designer.agent.md)
- agent contract: [designer.contract.json](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/references/designer/contracts/designer.contract.json)

## Maintenance

- 権限、write scope、output obligation を skill 本体へ戻さない。
- Copilot 実装 workflow の詳細は [.github/skills](/Users/iorishibata/Repositories/AITranslationEngineJP/.github/skills/) に置く。
- handoff 粒度の長い例は references に分離する。
