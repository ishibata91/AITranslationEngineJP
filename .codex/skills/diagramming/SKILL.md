---
name: diagramming
description: Codex 側の図作成作業プロトコル。PlantUML と structure diff の source と review artifact の扱いを提供する。
---
# Diagramming

## 目的

`diagramming` は作業プロトコルである。
`designer` agent が diagram を必要資料として作る時に、source of truth、review artifact、validation の見方を提供する。

標準 `implement_lane` flow には含めない。
人間が明示した時、または `wall-discussion` の結論として図が必要になった時だけ参照する。

## 対応ロール

- `diagrammer` が使う。
- 返却先は caller または次 agent とする。
- owner artifact は `diagramming` の出力規約で固定する。

## 入力規約

- 人間が明示的に図を求めた時
- `wall-discussion` の結論として図が必要になった時
- structure diff を review 可能な図にする時
- PlantUML source を作成、更新する時
- diagram source と rendered artifact の責務を分ける時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。
- 必須入力: caller, diagram_goal, source_of_truth
- 任意入力: diagram_kind, target_task_folder, validation_commands
- selectors: {"diagram_kind": ["structure-diff", "plantuml"]}
- 必須 artifact: diagram source or target task folder

## 外部参照規約

- エージェント実行定義とツール権限は [diagrammer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/diagrammer.toml) の `allowed_write_paths` / `allowed_commands` とする。
- standard caller: [design-bundle](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md)
- explicit helper binding: [diagrammer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/diagrammer.toml)
- エージェント実行定義: [diagrammer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/diagrammer.toml)
- ツール権限: エージェント実行定義の `allowed_write_paths` / `allowed_commands` に従う
- primary runtime skill: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md)
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/diagramming/SKILL.md

## 内部参照規約

### 拘束観点

- PlantUML を使った structure diff と設計補助図
- source of truth と review artifact の分離
- render、validate、差分確認の扱い
- 一時 PNG と AI 目視による可読性確認
- 図を分割すべき大きさの判断

### Readability Gate

- 一時 PNG は `/tmp` へ生成し、repo の正本にしない。
- AI が画像を確認し、線、文字、余白、線長、交差、legend の干渉を確認する。
- 可読性が悪い場合は、配置調整ではなく図の分割を先に検討する。

### Readability Patterns

- 正本図の用語、package 名、layer 名は、対象 docs の構造主語と節名に合わせる。
- architecture 図では [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) の構造主語、依存方向、節名を優先する。
- 正本図では可読性のために node を勝手に畳まず、まず edge、legend、note の粒度で調整する。
- package / group は意味単位または layer 単位で作り、上から下に主依存を読める順序に固定する。
- 長距離の wiring や bootstrap の concrete 生成線は、全 edge を描くより legend へ要約する。
- ER 図は table の意味単位で group を作り、関係線を最小限にして読む順序を安定させる。
- docs 正本図は neutral な構造図にし、差分図 style は `docs/exec-plans/` 配下の review artifact に限定する。
- `skinparam linetype ortho` のような角ばった線指定は、読みやすさが下がる場合があるため既定では使わない。

### Split Rule

- primary node が 12 個を超える時は分割する。
- attribute 付き class / table が 8 個を超える時は分割する。
- visible edge が 20 本を超える時は分割する。
- package / boundary が 4 個を超える時は分割する。
- package をまたぐ長距離 edge が複数ある時は overview と detail に分ける。

- 参照 template は [review-diff-style-legend.puml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/diagramming/references/templates/review-diff-style-legend.puml) とする。
- 参照 example は [review-diff-template.puml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/diagramming/references/examples/review-diff-template.puml) とする。

## 判断規約

- 図の source of truth を先に明示する
- review 用 SVG、PNG、screenshot を正本にしない
- PlantUML や diagram library の書き方は `npx ctx7 library` / `npx ctx7 docs` で確認する
- 図の方向は常に上から下にする
- validation できない diagram は gap として返す
- 原則日本語で書くが、必要に応じて英語も使う

- source と review artifact を別物として扱う
- validation 結果と未確認 gap を残す
- 図が補う設計判断を明示する
- 大きすぎる図は overview と detail に分ける
- active 規約 は agent に対して 1 ファイルだけ置く。diagram kind は selector で扱う。

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。
- 図化対象判断: どの diagram artifact を扱うかを返す。
- source target: 図の根拠にした source path を返す。
- review diff diagram: review に使う差分図を返す。
- 確認結果: render または validation の結果を返す。
- 未決事項: 採否判断に必要な open question を返す。

## 完了規約

- 出力規約を満たし、次の actor が再解釈なしで判断できる。
- 不足情報または停止理由がある場合は明示されている。
- 図の source of truth と review artifact を分けた。
- diagram kind と読者を明示した。
- PlantUML や library が関係する場合は `npx ctx7 library` / `npx ctx7 docs` で確認した。
- 一時 PNG を生成し、AI が画像で可読性を確認した。
- 図の分割条件に当たらないことを確認した。
- 正本 docs の用語、構造主語、layer 名と diagram の package 名を揃えた。
- 正本図では node を勝手に畳まず、edge、legend、note の粒度で調整した。
- 必須 evidence: source path, render or validation result
- 完了判断材料: designer が diagram artifact の採否を判断できる。
- 残留リスク: 採否判断に必要な未決事項が返っている。

## 停止規約

- UI 要件契約だけで diagram が不要な時
- プロダクトコードの構造を実装で変更する時
- docs 正本化だけが目的の時
- review 用 artifact を正本にしない
- プロダクトコード 変更を diagramming に混ぜない
- source 不明のまま図を更新しない
- PlantUML 以外の diagram source を新規作成しない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- review 用 SVG / PNG を正本にしなかった場合は停止する。
- validation なしで完了扱いにしなかった場合は停止する。
- プロダクトコード 変更を diagramming に混ぜなかった場合は停止する。
- 線、文字、legend、note の重なりを放置しなかった場合は停止する。
- docs 正本図を差分図 style にしなかった場合は停止する。
- 可読性向上を理由に正本 node を消さなかった場合は停止する。
- 角ばった線指定で読みやすさを下げなかった場合は停止する。
- 拒否条件: source of truth が不明
- 拒否条件: product implementation が必要
