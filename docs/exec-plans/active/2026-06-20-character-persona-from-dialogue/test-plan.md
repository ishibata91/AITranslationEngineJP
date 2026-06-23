# テスト計画（ESP 抽出結果での網羅回帰、2026-06-22）

本書は、口調・ペルソナ本実装が PoC で示した分類を実 ESP 抽出データで再現することを保証するテスト計画である。単体テスト（純粋 IO・100%、`implementation-scope.md`）とは別の、抽出結果 fixture を使う backend 特性テストを定める。根拠は `poc-tone-report.md`（golden の出典）と `persona-signal-map.md`。

## 目的

production の ToneClassifier・融合・語彙マーカーが、PoC で示した結果（特定キャラの基底口調、経路選択、修正の効果、横断 mod）を再現することを、回帰テストで保証する。較正や改修で品質が落ちたら検知する。

## テストの種類

- 単体（純粋・100%）: ToneClassifier の採点・融合・帯分け、性質文写像、合成順。`implementation-scope.md` に定義済み。本書では扱わない。
- 網羅回帰（本書）: ESP 抽出結果の fixture DB に production を通し、PoC の golden と照合する。DB を読むが LLM・UI を含まない backend テスト。prose は line_analysis のキャッシュ経由で 1 度だけ走る。

## fixture（ESP 抽出結果の subset）

- ESP 本体（skyrim.esm 等）は大容量・著作物のため repo に置けない。よって抽出結果から golden 対象の話者とその台詞だけを抜いた小さな fixture DB を `internal/engine/tone/testdata/` に置く。
- 作り方: skyrim.esm と inigo.esp を extractor で抽出した DB（`db/poc-skyrim.sqlite3`・`db/poc-inigo.sqlite3`）から、golden の話者 id に紐づく line・line_speaker・speaker・voice_type・race だけを抜き出して書き出す。再現用の抜き出しクエリを testdata 生成スクリプトに残す。
- fixture は 2 つ: skyrim 由来（本文・voice・保留・修正・マーカーを覆う）と inigo 由来（横断 mod）。

## golden（PoC で示した、保証する結果）

帯は基底口調の対人帯（尊大・中立・丁寧）。較正で動きうるスコアの厳密値は縛らず、明快なキャラは帯で断定、境界キャラは符号で断定する。

| 話者 | 由来 | 経路 | 保証する判定 | 狙い |
|---|---|---|---|---|
| Ulfric | skyrim | 本文 | 尊大（ぞんざい） | 命令的指導者 |
| GeneralTullius | skyrim | 本文 | 尊大 | 帝国軍将軍 |
| Astrid | skyrim | 本文 | 尊大（冷然・見下し） | 暗殺者ギルド長 |
| Tolfdir | skyrim | 本文 | 丁寧（物腰やわ） | 温厚な学院教師 |
| ArnielGane | skyrim | 本文 | 丁寧 | 臆病な学者 |
| Paarthurnax | skyrim | 本文 | 丁寧（慇懃） | 龍の賢者 |
| Maven | skyrim | 本文 | 中立（既知の許容誤差） | 含意の尊大は取れない |
| Nazeem | skyrim | voice | 尊大（横柄 prior） | 台詞少でも voice で見下し |
| Edorfin | skyrim | voice | 尊大（高慢 prior） | 1 台詞、voice が唯一の信号 |
| Grelod | skyrim | voice | 尊大＋高い感情 | 気難老＋激情 |
| Hadvar | skyrim | 保留 | 尊大でない | 誘導命令の除外（修正1） |
| Ralof | skyrim | 保留 | 尊大でない | 同上 |
| Delphine | skyrim | 本文 | Khajiit マーカー無し | 種族 gate（修正2、Breton） |
| Jzargo | skyrim | — | Khajiit マーカー有り | 種族 gate（実 Khajiit 保持） |
| Inigo | inigo | 本文 | 丁寧（物腰やわ）＋Khajiit マーカー | 追加 NPC を対話で分類 |
| Nazeem | inigo | voice | 尊大（横柄） | 既存 NPC が skyrim 版と一致 |

## 網羅する場面（漏れの無い分け方）

- 本文経路（印 10 以上）: 尊大・中立・丁寧の各代表に加え、印 10 以上の話者集合で対人帯の分布（尊大 4 / 中立 5 / 丁寧 16）と中立潰れ率（20% 前後）を範囲で検証する。
- voice prior 経路（印 10 未満）: 気質を持つ voice で対人帯を決める（Nazeem・Edorfin・Grelod）。気質を持たない voice は中立に落ちる（Constance）。
- 保留経路: 固有 voice ×台詞少（Hadvar 等）で、薄い本文値を低信頼で保持し、尊大に誤判定しない。
- 修正の保証: 誘導命令の除外（Hadvar・Ralof が尊大でない）、Khajiit 三人称が実 Khajiit だけ（Delphine 除外・Jzargo 保持）。
- 横断 mod: 追加 NPC Inigo を分類、既存 NPC（Nazeem）が skyrim 由来と一致、creature（Cr 声型）が対象外。
- 行キャッシュ: 同一本文の line が複数あるとき、prose 呼び出しが本文ユニーク数に等しい（重複本文で再実行しない）。

## 判定の粒度（較正に頑健）

- 明快なキャラ: 対人帯（尊大／中立／丁寧）の離散値で断定する。
- 境界キャラ: 対人スコアの符号（負＝尊大寄り、正＝丁寧寄り）で断定し、厳密値で縛らない。
- 分布: count の範囲（例 尊大 3〜5）で判定し、厳密一致にしない。較正で 1 つ動いても通る幅にする。
- マーカー: 有無を断定する（Khajiit 三人称の有無など）。

## 期待値の根拠と更新

- golden の出典は `poc-tone-report.md` の結果表である。本実装の較正でスコアが動く項目は、帯・符号・分布で固定し、スコアの厳密値を golden にしない。
- 較正でしきい値を意図的に変える時は、golden の帯・分布が変わる理由を `poc-tone-report.md` と本書へ追記してから golden を更新する。理由なく golden を緩めない。
