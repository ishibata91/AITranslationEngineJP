---
name: diagramming
description: Codex 側の図作成作業プロトコル。PlantUML と structure diff の source と レビュー 成果物 の扱いを提供する。
---
# Diagramming

## 目的

`diagramming` は作業プロトコルである。
`designer` agent が diagram を必要資料として作る時に、正本、レビュー 成果物、検証 の見方を提供する。

標準 `implement_lane` flow には含めない。
人間が明示した時、または `wall-discussion` の結論として図が必要になった時だけ参照する。

## 対応ロール

- `diagrammer` が使う。
- 返却先は 呼び出し元 または次 agent とする。
- 担当成果物は `diagramming` の出力規約で固定する。

## 入力規約

- 不足時の扱い: 入力に 根拠参照、担当者、承認状態が不足する場合は推測で補わない。
- 必須入力: 呼び出し元, diagram_goal, source_of_truth
- 任意入力: diagram_kind, target_task_folder, 検証コマンド
- selector: {"diagram_kind": ["structure-diff", "plantuml"]}
- 必須 成果物: diagram source or 対象 task folder

## 外部参照規約

- エージェント実行定義とツール権限は [diagrammer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/diagrammer.toml) の 書き込み許可 / 実行許可 とする。
- standard 呼び出し元: [design-bundle](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md)
- explicit helper 紐づけ: [diagrammer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/diagrammer.toml)
- エージェント実行定義: [diagrammer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/diagrammer.toml)
- ツール権限: エージェント実行定義の 書き込み許可 / 実行許可 に従う
- primary 実行定義 skill: [SKILL.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/design-bundle/SKILL.md)
- 外部成果物 が不足または衝突する場合は停止し、衝突箇所を返す。
- 関連 skill: /Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/diagramming/SKILL.md

## 内部参照規約

### 拘束観点

- PlantUML を使った structure diff と設計補助図
- 正本 と レビュー 成果物 の分離
- render、validate、差分確認の扱い
- 一時 PNG と AI 目視による可読性確認
- 図を分割すべき大きさの判断

### Readability Gate

- 一時 PNG は `/tmp` へ生成し、repo の正本にしない。
- AI が画像を確認し、線、文字、余白、線長、交差、legend の干渉を確認する。
- 可読性が悪い場合は、配置調整ではなく図の分割を先に検討する。

### Readability Patterns

- 正本図の用語、package 名、層 名は、対象 docs の構造主語と節名に合わせる。
- architecture 図では [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) の構造主語、依存方向、節名を優先する。
- 正本図では可読性のために node を勝手に畳まず、まず edge、legend、補足 の粒度で調整する。
- package / group は意味単位または 層 単位で作り、上から下に主依存を読める順序に固定する。
- 長距離の wiring や bootstrap の concrete 生成線は、全 edge を描くより legend へ要約する。
- ER 図は table の意味単位で group を作り、関係線を最小限にして読む順序を安定させる。
- docs 正本図は neutral な構造図にし、差分図 style は `docs/exec-plans/` 配下の レビュー 成果物 に限定する。
- `skinparam linetype ortho` のような角ばった線指定は、読みやすさが下がる場合があるため既定では使わない。

### Split Rule

- primary node が 12 個を超える時は分割する。
- attribute 付き class / table が 8 個を超える時は分割する。
- visible edge が 20 本を超える時は分割する。
- package / 境界 が 4 個を超える時は分割する。
- package をまたぐ長距離 edge が複数ある時は overview と detail に分ける。

- 参照 雛形 は [review-diff-style-legend.puml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/diagramming/references/templates/review-diff-style-legend.puml) とする。
- 参照 example は [review-diff-template.puml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/diagramming/references/examples/review-diff-template.puml) とする。

## 判断規約

- 図の 正本 を先に明示する
- レビュー 用 SVG、PNG、screenshot を正本にしない
- PlantUML や diagram library の書き方は `npx ctx7 library` / `npx ctx7 docs` で確認する
- 図の方向は常に上から下にする
- 検証 できない diagram は 不足 として返す
- 原則日本語で書くが、必要に応じて英語も使う

- source と レビュー 成果物 を別物として扱う
- 検証 結果と未確認 不足 を残す
- 図が補う設計判断を明示する
- 大きすぎる図は overview と detail に分ける
- active 規約 は agent に対して 1 ファイルだけ置く。diagram kind は selector で扱う。

## 非対象規約

- UI 要件契約だけで足りる作業は扱わない。
- プロダクトコード構造の実装修正は扱わない。
- docs 正本化だけの作業は扱わない。
- レビュー用 SVG、PNG、screenshot は正本にしない。
- PlantUML 以外の diagram source は新規作成しない。

## 出力規約

- 基本出力: 出力は判断結果、根拠参照、不足情報、次 agent が判断できる材料を含む。
- 禁止事項: 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。
- 図化対象判断: どの diagram 成果物 を扱うかを返す。
- source 対象: 図の根拠にした 根拠 path を返す。
- レビュー diff diagram: レビュー に使う差分図を返す。
- 確認結果: render または 検証 の結果を返す。
- 未決事項: 採否判断に必要な open question を返す。

## 完了規約

- 出力規約を満たし、次の 実行者 が再解釈なしで判断できる。
- 不足情報または停止理由がある場合は明示されている。
- 図の 正本 と レビュー 成果物 を分けた。
- diagram kind と読者を明示した。
- PlantUML や library が関係する場合は `npx ctx7 library` / `npx ctx7 docs` で確認した。
- 一時 PNG を生成し、AI が画像で可読性を確認した。
- 図の分割条件に当たらないことを確認した。
- 正本 docs の用語、構造主語、層 名と diagram の package 名を揃えた。
- 正本図では node を勝手に畳まず、edge、legend、補足 の粒度で調整した。
- 必須 根拠: 根拠 path, render or 検証結果
- 完了判断材料: designer が diagram 成果物 の採否を判断できる。
- 残留リスク: 採否判断に必要な未決事項が返っている。

## 停止規約

- UI 要件契約だけで diagram が不要な時
- プロダクトコードの構造を実装で変更する時
- docs 正本化だけが目的の時
- source が不明な場合は停止する。
- 停止時は不足項目、衝突箇所、戻し先を返す。
- render または検証ができない場合は停止する。
- 線、文字、legend、補足 の重なりを解消できない場合は停止する。
- 可読性向上を理由に正本 node を消さなかった場合は停止する。
- 角ばった線指定で読みやすさを下げなかった場合は停止する。
- 拒否条件: 正本 が不明
- 拒否条件: プロダクト実装 が必要
