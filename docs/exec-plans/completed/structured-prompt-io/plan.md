# Plan: structured-prompt-io

`task_id`: `structured-prompt-io`

`分岐元 branch`: `master`

`分岐元 commit`: `8b88c04f`

## 依頼要約

プロンプトの構造化入力と構造化出力（`response_format` の `json_schema`）に対応する。まず LM Studio（唯一の同期実装 `OpenAICompatible`）で動かす。openAI 式・xAI 式の仕様を調べ、同期・バッチで OpenAI 互換が使えるかを確認し、構造化プロンプトと構造化出力をどこまで抽象化するかを決める。

## 完了定義

`OpenAICompatible` 経由の単一行翻訳を、自由文字列受けからスキーマ準拠の構造化出力に置き換え、LM Studio 実接続で動くところまでを完了とする。

- 動かす範囲: 翻訳リクエストに `response_format`（`json_schema`）を載せ、構造化レスポンスから訳文を取り出し、既存の `dest` 書き戻し経路へ流す。入力は現状の文字列プロンプト（system・user）のまま変えず、構造化は出力側だけに限る。1 行ずつ同期のまま変えない。
- 観測点: backend 単体テスト（リクエストに `response_format` が載る、構造化レスポンスの解析、訳文抽出の失敗分類の純粋関数）と、実画面（LM Studio 実接続で翻訳を回し、構造化出力で訳文が `dest` に入る）。
- goal 整合: 「構造化出力に対応する」を、実際に `OpenAICompatible` の 1 経路で構造化リクエスト送信と構造化レスポンス解析が動くことで満たす。パースの堅牢化を空実装で満たしたことにしない。

## scope

- 含む: 単一行訳文のスキーマ堅牢化。出力スキーマは `translation` 1 フィールド（`strict`・required）。`response_format`（`json_schema`）の付与、構造化レスポンス解析、訳文抽出の失敗分類の純粋関数。構造化パース失敗の番兵エラー（`ErrStructuredParse`）と、engine が失敗行を skip して翻訳を継続する失敗ハンドリング（実測で判明し人間判断で追加）。入力は現状の文字列プロンプトのまま。対象実装は `OpenAICompatible`（LM Studio 接続）と `engine` の本文翻訳ループ。
- 含まない（別の新規 plan とする）: 複数対象を 1 リクエストにまとめる一括構造化。コンテキスト長でのチャンク化、クエスト・会話ツリーでのグループ化、依存関係のないデータのまとめ。バッチ API（`gemini-xai-batch-translation` が別に扱う）。クラウドプロバイダ（Gemini / xAI / Claude）の新規実装。

## 調査で確定した事実

外部仕様と現状コードを調べ、抽象化範囲の判断材料を固定した。

**外部仕様は 3 経路とも同一形式に収束する。**

| 対象 | 構造化出力の指定 | バッチ |
| --- | --- | --- |
| OpenAI | `response_format` の `json_schema`（`name`・`strict`・`schema`）を chat/completions | Batch API 対応 |
| xAI（Grok） | 同上。OpenAI 互換、`base_url` は `api.x.ai/v1` | `/v1/batches` あり。jsonl 各行の `body` が chat/completions ペイロード |
| LM Studio | 同上。OpenAI 互換、`/v1/chat/completions` | 同期のみ |

- openAI 式と xAI 式は別物でなく、xAI が OpenAI 互換のため `response_format` の `json_schema` で同一。LM Studio も同じ。
- xAI のバッチは各行の `body` が chat/completions ペイロードそのもの。同期で使う `response_format` をバッチの `body` にそのまま載せられる。同期・バッチで構造化指定の形は共通。
- 制約: LM Studio は 7B 未満モデルで構造化出力が効かないことがある。実 LLM は gemma-12b 等のため実害は小さい。

**現状コードは構造化出力の痕跡がない。**

- provider の interface 境界は `Translator`（`Translate`・`ListModels`）1 つ、実装は `OpenAICompatible` のみ。LM Studio は同じ実装に URL を変えて使う。
- リクエストは chat/completions に messages 配列、レスポンスは `choices[0].message.content` を自由文字列で受けて `dest` へ書き戻す。`response_format`・`json_schema` は未使用。
- プロンプト組立は `internal/core/prompt` の `ComposePrompt`（純粋関数、自由文字列）。JSON 期待はない。
- 翻訳は全経路 1 行ずつ同期。バッチとクラウドプロバイダは未実装。

## 抽象化範囲の見立て

今回は単一行のため、入出力スキーマを差し替える大掛かりな port は作らず、最小の追加で足りる見立て。純粋モジュール化の線引きも下記で固定する。

- provider は「完成プロンプトを送るだけ、文面構築しない」を維持する。base 指示・口調指示の合成は従来どおり engine が担う。
- 出力スキーマは静的定義（`translation` 1 フィールド固定）のため、独立モジュールへ切り出さず、provider が `chatRequest` に固定 struct として載せて送る。`response_format` は HTTP ペイロードの一部で、文面構築には当たらない。
- 訳文抽出（content 文字列から訳文を取り出し失敗を分類する）だけを純粋関数へ分離する。判断の分岐（構文不正・必須欠落/空値・正常）を持つため。100% カバレッジの対象は本関数の分岐網羅に限る。
- 将来の複数対象一括（別 plan）を、同じ `response_format` の口に載せられる素地を壊さない。訳文抽出の純粋関数が配列版（件数一致・id 対応・部分失敗）へ育つ入口を兼ねる。ただし今回は配列化を実装しない。

## パース失敗時の扱い

段2（content から訳文抽出）で訳文が取れない場合は、provider が番兵エラー `ErrStructuredParse` でラップして返し、engine がその行を未訳のまま skip して翻訳を継続する（既存のタグ欠落と同じ扱い、再実行で再翻訳）。1 行の構造化パース失敗で翻訳全体を止めない。フォールバック（content 全体を訳文とみなす）はしない。

- provider: 訳文抽出の失敗（構文不正・必須欠落・空値）を `ErrStructuredParse` でラップする。空応答は `json.Unmarshal` の "unexpected end of JSON input"（構文不正）に落ちる。
- engine: `errors.Is(err, provider.ErrStructuredParse)` で識別し、`parseFailures` を数えて `continue`（skip）。フェーズ末に集約観測ログ（`structured_parse_failed`, result=skipped）を出す。通信・HTTP 非 200・choices 空は従来どおり Run 全体停止。
- 当初は「失敗として上位へ返す」設計だったが、engine が 1 行失敗で Run を停止する既存挙動と組み合わさると、実 LLM の非決定的な空応答で翻訳全体が止まる回帰が判明した。人間判断で「タグ欠落と同じ skip」に決定し、scope を engine の失敗ハンドリングまで広げた。

## 軽 / 重判定

- 画面が動くか: N。単一行のスキーマ堅牢化は backend 中心で、layout・文言・style・表示構造・story を変えない見込み。
- `docs/architecture.md` 反映が要るか: N。層構成・依存方向・Bootstrap・Wails 境界は不変。強い制約（port は `provider` 1 つ、provider は文面構築しない）は維持する（`response_format` の付与は HTTP ペイロード組立で文面構築でない）。engine の失敗ハンドリング追加も層・依存を変えない。

判定: 両方 N のため軽 task。`design-module` と `storybook-module` を bypass し、`preparation-module` → `implementation-module` → `finalization-module` で進める。

## close_conditions

- backend 単体テストが、リクエストへの `response_format`（`json_schema`）付与、構造化レスポンスからの訳文抽出、パース失敗の型分け（構文不正・空応答・必須欠落/空値）、engine が構造化パース失敗行を skip して継続することを検証し、緑になる。
- 訳文抽出の失敗分類の純粋関数は、ユニットテスト 100% カバレッジ（構文不正・空応答・必須欠落/空値・正常の各分岐）を満たし、失敗は `ErrStructuredParse` でラップされる。
- 実画面で LM Studio 実接続の翻訳を回し、構造化出力で訳文が `dest` に入ること、および構造化パース失敗行が全体を止めず skip されることを目視で確認する。

## 最終検証結果

harness 検証と実画面確認を通過した。

- backend テスト緑（`npm run test:backend`）。訳文抽出の純粋関数はカバレッジ 100%。lint 全通過（`npm run lint:backend`、format/vet/static/arch/boundary/module）。
- 実画面（LM Studio 実接続）: 構造化リクエスト送信・構造化レスポンス解析・`dest` 書き戻しが動作。構造化パース失敗行を skip して翻訳が継続することを確認（前回停止点 19 を超えて 46 まで進行）。
- モデル差の観測: 実環境のロードモデルは `hy-mt2-7b`（7B）と `google/gemma-4-12b-qat`。7B は特定行で空応答を返し skip 対象になる。gemma-12b は構造化出力が安定し、訳文に記号混入もない。plan の「実 LLM は gemma-12b 等のため実害は小さい」は妥当と確認。
- 副次観測（別軸、本 plan の scope 外）: engine の Run は選んだ plugin 単位でなく DB 全体の未訳行を翻訳する挙動。

## 正本化判断

- `docs/architecture.md` 反映: 不要。engine への失敗ハンドリング追加と provider の番兵エラー追加は、層構成・依存方向・Wails 境界・強い制約を変えない。恒久仕様の新規正本化なし。
- 後続課題（別 plan）: 一部プラグイン（Outfit Recognition Framework, 1805 行）で抽出フェーズが約 6 分と異常に遅い件を別 plan で扱う。
