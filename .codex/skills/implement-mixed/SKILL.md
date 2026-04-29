---
name: implement-mixed
description: Codex implementation lane 側の API / Wails / DTO / gateway など frontend と backend の接合点実装作業プロトコル。
---
# Implement Mixed

## 目的

この skill は作業プロトコルである。
`implementation_implementer` agent が scope freeze 済みの API、Wails binding、DTO、gateway、adapter contract など frontend と backend の接合点 owned_scope を実装する時の判断基準を提供する。

mixed は広い frontend / backend 同時変更の許可ではない。
片側だけで閉じる UI 実装や backend 実装は、それぞれ `implement-frontend` または `implement-backend` を使う。

## 対応ロール

- `implementation_implementer` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- owner artifact は `implement-mixed` の出力規約で固定する。

## 入力規約

- implementation-scope が API、Wails binding、DTO、gateway、adapter contract の接合点変更を明示している時
- 片側だけでは contract 整合を証明できない owned_scope を扱う時
- validation を frontend / backend の接合点 evidence として返す時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_implementer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_implementer.toml) の `allowed_write_paths` / `allowed_commands` とする。
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

## 判断規約

- implementation-scope の owned_scope を守る
- mixed の対象を API、Wails binding、DTO、gateway、adapter contract の接合点だけに限定する
- 片側だけで閉じない理由を scope artifact で確認する
- lane_context_packet を確認して プロダクトコード だけを変更する
- `APIテスト` 先行時だけ implementation_tester output も確認する
- validation は frontend、backend、接合点 contract の証跡を分ける

- API / Wails / DTO / gateway / adapter contract のどれを接合点として変更したか closeout に残す
- 両側の touched files を handoff と対応づける
- frontend / backend / 接合点 contract の lane-local validation evidence を分ける
- lane-local validation command の不足を residual risk にする

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

## 完了規約

- 承認済み owned_scope 内の成果だけが返却されている。
- validation、未実行項目、residual risk が source_ref 付きで整理されている。
- API / Wails / DTO / gateway / adapter contract の接合点 scope が承認済みであることを確認した。
- 両側の touched files を handoff と対応づけた。
- lane_context_packet と lane-local validation evidence を分けた。
- `APIテスト` 先行時だけ implementation_tester output を確認した。

## 停止規約

- frontend または backend の片側だけで閉じる時
- API / Wails / DTO / gateway / adapter contract の接合点を変更しない時
- 横断範囲が未承認の時
- 追加設計で横断 scope を広げる時
- mixed を広い frontend / backend 同時変更の口実にしない
- 片側の都合で scope を広げない
- API 接合点を変えずに UI と backend を同時に触らない
- プロダクトテスト、fixture、snapshot、test helper を変更しない
- docs や workflow 文書を変更しない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- mixed を広い frontend / backend 同時変更の口実にしなかった場合は停止する。
- API 接合点を変えずに UI と backend を同時に触らなかった場合は停止する。
- プロダクトテスト、fixture、snapshot、test helper を変更しなかった場合は停止する。
- docs / workflow 文書を変更しなかった場合は停止する。
