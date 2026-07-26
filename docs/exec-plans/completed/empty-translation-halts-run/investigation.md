# Investigation: empty-translation-halts-run

`investigation.md` は不具合の再現と原因究明だけを持つ。どう直すかの設計は `design.md` が持つ。

## 観測済み問題

- 固有名 1 件の応答が空だと、その実行が画面のエラー表示で終わり、残り全件が未訳のまま残る。エラー文言は `翻訳に失敗: 固有名の翻訳: 翻訳応答解析: 構造化出力の解析失敗: translation が空`。
- 本文（叙述文・台詞）は 1 件も翻訳要求が送られない。固有名フェーズが本文フェーズより前に走り、そこで実行が終わるため。

## 画面再現確認

**再現の作り方**: 実 LLM（LM Studio の `hy-mt2-7b`）の空応答は非決定的なので、OpenAI 互換の応答を返す観測用スタブサーバ（`127.0.0.1:18234`）を立て、1 件目の翻訳要求だけ `{"translation":""}` を返し 2 件目以降は訳文を返す形に固定した。画面のエンドポイント欄がスタブを指せるため、実画面の操作だけで再現できる。

**再現手順**（`chrome-devtools` MCP ツールで実施。画面に `data-testid` は無いため、見出しとボタン文言のセレクタで操作した）

1. 翻訳対象プラグイン画面で `inigo.esp` の成果を削除し、未訳状態へ戻す（固有名 146 件、叙述文 112 件、台詞 8545 件がすべて未訳）。
2. プラグイン欄へ `/Users/iorishibata/Repositories/AITranslationEngineJP/dictionaries/Data/inigo.esp` を入れ、`翻訳へ進む` を押す。
3. 翻訳実行画面でエンドポイント欄へ `http://127.0.0.1:18234/v1` を入れ、`取得` でモデル `stub-empty-first` を選ぶ。
4. `実行` を押す。

**観測した画面状態**

- 状態表示が `失敗` になり、`role="alert"` の領域へ `翻訳に失敗: 固有名の翻訳: 翻訳応答解析: 構造化出力の解析失敗: translation が空` が出た。人間観測記録と同一文言。
- 結果一覧は `まだ結果はありません。` のままで、訳文が 1 件も並ばない。

**観測した外部境界と DB**

- スタブが受けた翻訳要求は 1 件だけ（`user='Inigo Companions Commentary'`、空応答）。2 件目以降の要求が来ていない＝実行が続いていない。
- 実行後の DB は固有名 146 件中 146 件、叙述文 112 件中 112 件、台詞 8545 件中 8545 件が未訳。空応答 1 件で全件が未訳のまま残ることを DB で確認した。

## 原因仮説

観測事実（要求 1 件で止まる、失敗分類の文言が provider 由来）から 3 つ立て、外側から順に確かめた。

- 仮説1: 画面側が失敗を受けて実行を打ち切っている。検証順序は 1 番目（画面に一番近い層のため）。根拠は、停止が画面のエラー表示として現れること。
- 仮説2: provider が空応答を「その行だけ飛ばせる失敗」として分類していない。検証順序は 2 番目。根拠は、エラー文言が provider の失敗分類（`構造化出力の解析失敗`）を含むこと。
- 仮説3: 固有名フェーズが失敗分類を見ずに、失敗をそのまま呼び出し元へ返している。検証順序は 3 番目。根拠は、本文フェーズが同じ失敗分類で行を未訳据え置きにしている一方、エラー文言の先頭が `固有名の翻訳` であること。

## 観測ログ検証

**追加した一時ログ**: `internal/engine/proper_noun.go` の `translateProperNouns` の翻訳失敗地点へ、失敗分類（`provider.ErrStructuredParse` かどうか）と、本文フェーズが使う純粋規則 `batchplan.DecideApply` の判定値を出す WARN ログを 1 個だけ足した。

**観測結果**（`tmp/logs/wails-dev.log`）

- `{"msg":"tmp: proper noun translate failed","where":"translateProperNouns","result":"halt","id":1,"structured_parse":true,"outcome_kind":2}`
- `structured_parse: true` は、provider が空応答を番兵エラー `ErrStructuredParse` で包んでいることを示す。よって仮説2 は否定される。
- `outcome_kind: 2` は `batchplan.ApplySkipStructuredParse`（未訳据え置き）である。同じ失敗を本文フェーズと同じ規則へ通せば「その固有名だけ未訳で据え置き」と判定されるのに、固有名フェーズはその判定を取らずに失敗を返している。よって仮説3 が支持される。
- スタブへの要求が 1 件で止まっており、画面操作の前に backend 側で実行が終わっている。よって仮説1 は否定される。

**削除確認**: 一時ログと、そのために足した import を削除した。`grep -rn 'tmp: proper noun' internal/` が空、`git diff --stat` が空（プロダクトコードの差分なし）であることを確認した。

## 確定原因

固有名フェーズ（`internal/engine/proper_noun.go` の `translateProperNouns`）が、翻訳失敗を失敗分類で振り分けずに、そのまま呼び出し元（`internal/engine/engine.go` の `Run`）へ返している。`Run` は固有名フェーズの失敗で戻るため、後続の叙述文・台詞フェーズへ進まない。

同じ失敗分類に対する扱いが、次の 3 経路で揃っていない。

- 同期の本文フェーズ（`translateNarrations`、`translateLines`）: `batchplan.DecideApply` で振り分け、構造化出力の空はその行だけ未訳据え置きにして実行を続ける。
- batch の固有名反映（`internal/engine/batch.go` の `applyOne`）: 同じ `batchplan.DecideApply` を通し、未訳据え置きにする。
- 同期の固有名フェーズ: 振り分けを通さず、実行を止める。
