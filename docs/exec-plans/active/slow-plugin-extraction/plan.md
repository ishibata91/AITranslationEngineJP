# Plan: slow-plugin-extraction

`task_id`: `slow-plugin-extraction`

`分岐元 branch`: `master`

`分岐元 commit`: `9c993d09`

## 依頼要約

一部プラグインで抽出フェーズ（esp から翻訳対象を抽出する段）が異常に遅い。原因を特定し高速化する。

## 観測記録

`structured-prompt-io` task の実画面確認中に観測した。

- Outfit Recognition Framework.esp（633KB、翻訳対象 1805 行）: 抽出フェーズが約 6 分。app 本体プロセスが CPU 約 108%（STAT=R）で解析し続けた。ハングではなく実処理。
- Innocence Lost - Quest Expansion.esp（124KB、197 件表示）: 抽出フェーズが約 60 秒。
- 抽出は Go backend（app 本体）が実行する（C# 別プロセスでない）。翻訳フェーズ（LLM）とは別段で、遅いのは抽出段。
- サイズ比は約 5 倍（124KB→633KB）だが、抽出時間比は約 6 倍以上（60 秒→約 6 分）。行数・レコード数に対して非線形に遅くなる可能性がある。

## 完了定義

抽出フェーズの律速箇所を特定し、非効率を解消して、大きめプラグイン（Outfit Recognition Framework.esp）の抽出時間を実測で有意に短縮するところまでを完了とする。

- 動かす範囲: 抽出処理の律速箇所（毎回の master フル解析、レコード数に対する非線形処理などの候補）を特定し、解消する。律速を特定するだけ・差込点を置くだけで実測が速くならない状態を「動く」としない。
- 観測点: 実データ（Outfit Recognition Framework.esp, 1805 行）の抽出時間を before / after で実測する。backend のベンチ、または実画面のいずれかで確かめる。
- goal 整合: 「抽出が異常に遅い件を高速化」を、Outfit 規模の抽出時間が有意に短縮することで満たす。律速特定の報告だけで完了としない。

## scope

- 含む: 抽出フェーズ（esp 解析 → 翻訳対象の抽出・DB 取込）の性能。律速箇所の特定と解消。
- 含まない: 翻訳フェーズ（LLM 呼び出し）の性能。構造化出力（`structured-prompt-io` で完了）。抽出結果の正しさ変更（抽出する行の増減）。

## 調査方針（後続モジュールで詳細化）

原因が未特定のため、`investigation-module` で観測ログ駆動に律速を特定してから修正方針を固定する。

- 抽出処理のコード経路を特定する（どの package・関数が esp を解析するか）。
- master（Skyrim.esm 249MB 等）を毎回フル解析するか、キャッシュ・部分読みがあるかを確認する。
- レコード数・行数に対する計算量（線形か、O(n^2) 等の非線形か）を pprof 等で測る。

## 軽 / 重判定

- 画面が動くか: N。抽出は backend 処理で、layout・文言・style・表示構造・story を変えない見込み。
- `docs/architecture.md` 反映が要るか: 未確定（暫定 N）。既存層内の性能改善なら不要。抽出処理の構造変更（キャッシュ層の新設など）が層・依存・Wails 境界を変える場合は後続モジュールで再評価する。

判定: 暫定 軽 task（性能改善）。ただし原因未特定のため `investigation-module` で観測記録から律速を特定し、修正方針を固定してから `implementation-module` へ進む。architecture 反映が要る構造変更が判明したら `design-module` へ迂回する。
