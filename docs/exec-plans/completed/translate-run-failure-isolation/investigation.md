# Investigation: translate-run-failure-isolation

`investigation.md` は不具合の再現と原因究明だけを持つ。修正フロー（`fix-workflow`）だけが作る。
どう直すかの設計は持たない（`design.md` が持つ）。原因究明は `fix-decision` skill を読んで適用する。

## 観測済み問題

- 台詞・叙述文の翻訳中に 1 リクエストの応答が失敗すると、その行だけを飛ばさず翻訳 run 全体が中断する。
- 対象の失敗は次の 4 系統。応答エンベロープの decode 失敗、応答に `choices` が無い、非 200 応答、通信失敗。
- 応答 1 件の失敗が、残りの未訳行の翻訳をすべて止める。中断までに書き戻した仮訳（`statusProvisional`）は DB に残る。残りは未訳のまま残り、UI は「翻訳に失敗」を表示する。

## 画面再現確認

本 issue は backend の失敗処理経路の設計不足であり、UI で観測された欠陥報告ではない。人間観測記録は known-issues #7 の本文（該当コード経路を明記済み）である。原因は制御フローの事実で、後述のとおりソース読解で決定的に確定する。

UI からの再現には限界がある。狙う症状「1 リクエストの mid-run 失敗で残り全行が止まる」を UI 操作で誘発するには、特定の 1 リクエストだけを失敗させる故障注入が要る。故障注入はプロダクトコードの変更を要し、本モジュール（`fix-workflow`）の安全境界（一時観測ログ以外のプロダクトコード変更禁止）に反するため、UI 単独では再現できない。

粗い症状「基盤失敗で run が止まる」は、接続先を無効な endpoint にして翻訳を実行すれば UI で観測できる（最初の行の通信失敗で run が中断し「翻訳に失敗」が出る）。ただし、この粗い再現は「全行が同じ理由で失敗する」ケースと「1 行だけ失敗して残りが道連れになる」ケースを区別しない。区別に要る mid-run 単発失敗の再現は、上記のとおり UI では作れない。

## 原因仮説

観測事実（症状）と、`internal/engine/engine.go`・`internal/provider/openai_compatible.go` のソースから、原因候補を層ごとに立てる。

- 仮説1（provider 層のエラー分類）: `provider.Translate` が、モデル起因の構造化出力失敗だけを番兵エラー `ErrStructuredParse` でラップし、通信・非 200・decode・`choices` 無しは番兵化していない。
- 仮説2（engine 層の loop 継続条件）: `engine.translateNarrations`・`engine.translateLines` の loop が、`errors.Is(err, provider.ErrStructuredParse)` の場合だけ `continue` し、それ以外の全エラーで `return` して `Run` を中断する。
- 検証順序: 仮説2 を先に確認する（run 全体を止めているのは engine の `return`）。次に仮説1 を確認する（engine が区別できる情報を provider が渡していない）。

## 観測ログ検証

一時観測ログを追加しない。原因は制御フローの静的事実であり、ログ観測より強い根拠（ソース上の分岐）で確定するため、`fix-decision` の「画面操作結果だけで原因仮説を固定しない」要請はソース確認で満たす。以下はソース上の該当箇所。

- 仮説2 の確定（engine の loop）:
    - `engine.go:264-272`（`translateNarrations`）: `provider.Translate` のエラーが `ErrStructuredParse` なら `parseFailures++; continue`、それ以外は `return fmt.Errorf("叙述文の翻訳: %w", err)`。
    - `engine.go:313-321`（`translateLines`）: 同じ構造で、`ErrStructuredParse` 以外は `return fmt.Errorf("台詞の翻訳: %w", err)`。
    - `engine.go:216-222`（`Run`）: `translateNarrations`・`translateLines` の返すエラーをそのまま `return` し、run を終える。
    - `internal/api/app.go:466-474`（`RunExtractAndTranslate`）: `engine.Run` の runErr を「翻訳に失敗: %w」で包み frontend へ返す。UI は失敗表示になる。
- 仮説1 の確定（provider のエラー分類）:
    - `openai_compatible.go:127-129`（通信失敗）: `return "", fmt.Errorf("翻訳要求: %w", err)`。番兵なし。
    - `openai_compatible.go:132-134`（非 200）: `return "", fmt.Errorf("翻訳要求: status %d", resp.StatusCode)`。番兵なし。
    - `openai_compatible.go:142-144`（エンベロープ decode 失敗）: `return "", fmt.Errorf("翻訳応答のデコード: %w", err)`。番兵なし。
    - `openai_compatible.go:145-147`（`choices` 無し）: `return "", fmt.Errorf("翻訳応答に choices が無い")`。番兵なし。
    - `openai_compatible.go:148-152` → `extractTranslation`（`openai_compatible.go:207-221`）: content 解析失敗だけ `ErrStructuredParse` でラップ（`engine` が唯一飛ばせる失敗）。

## 確定原因

- engine の本文フェーズ loop（`translateNarrations`・`translateLines`）は、`provider.ErrStructuredParse` の 1 種類だけを未訳のまま飛ばし、それ以外の全エラーで `return` して `Run` を中断する。
- provider の `Translate` は、通信失敗・非 200・エンベロープ decode 失敗・`choices` 無しを番兵エラーで分類しておらず、engine が「その行だけ飛ばせる失敗」と「run を止めるべき失敗」を区別する手掛かりを持たない。
- 2 つが重なり、モデル起因以外の 1 リクエスト失敗が run 全体を止める。これが known-issues #7 の症状の原因である。
