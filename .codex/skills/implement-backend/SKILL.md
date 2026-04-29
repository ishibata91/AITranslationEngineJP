---
name: implement-backend
description: Codex implementation lane 側の backend 実装作業プロトコル。layer 責務、lane-local validation の判断基準を提供する。
---
# Implement Backend

## 目的

この skill は作業プロトコルである。
`implementation_implementer` agent が backend owned_scope を実装する時に、usecase、service、repository、adapter の責務整合と dependency direction を守る判断基準を提供する。

## 対応ロール

- `implementation_implementer` が使う。
- 呼び出し元は `implement_lane` とする。
- 返却先は `implement_lane` とする。
- owner artifact は `implement-backend` の出力規約で固定する。

## 入力規約

- backend package を変更する時
- lane-local validation と error path を実装する時
- 入力に source_ref、owner、承認状態が不足する場合は推測で補わない。

## 外部参照規約

- エージェント実行定義とツール権限は [implementation_implementer.toml](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/agents/implementation_implementer.toml) の `allowed_write_paths` / `allowed_commands` とする。
- 外部 artifact が不足または衝突する場合は停止し、衝突箇所を返す。

## 内部参照規約

## 判断規約

- layer 責務と依存方向を守る
- error path と validation を owned_scope 内で閉じる
- single_handoff_packet と owned_scope を確認して プロダクトコード だけを変更する
- `APIテスト` 先行時だけ implementation_tester output も確認する
- lane-local validation result または未実行理由を返す
- `lint:backend` の format、vet、static、arch、module で落ちる境界違反を事前に避ける

- usecase / service / repository / adapter の責務を確認する
- [lint-policy.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/lint-policy.md) の backend lint 内訳を確認する
- [architecture.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/architecture.md) の backend 依存方向に従い、usecase、service、repository、adapter concrete の境界を跨がない
- usecase から repository concrete、runtime concrete、driver API を直接参照しない
- lane-local validation を実行した場合は結果を closeout に残す

## 出力規約

- 出力は判断結果、根拠 source_ref、不足情報、次 agent が判断できる材料を含む。
- 出力にツール権限、エージェント実行定義、プロダクトコードの変更義務を含めない。

## 完了規約

- 承認済み owned_scope 内の成果だけが返却されている。
- validation、未実行項目、residual risk が source_ref 付きで整理されている。
- layer 責務と dependency direction を確認した。
- backend lint の format、static、arch、module 観点を確認した。
- validation と error path を owned_scope 内で確認した。
- single_handoff_packet と lane-local validation を確認した。
- `APIテスト` 先行時だけ implementation_tester output を確認した。

## 停止規約

- frontend だけの変更を実装する時
- UI check を行う時
- backend 境界を設計し直す時
- owned_scope 外の layer refactor を混ぜない
- controller、usecase、service で concrete 実装を new しない
- service core から filesystem、Wails runtime、DB driver の concrete API を直接呼ばない
- プロダクトテスト、fixture、snapshot、test helper を変更しない
- docs や workflow 文書を変更しない
- 停止時は不足項目、衝突箇所、戻し先を返す。
- owned_scope 外の layer refactor を混ぜなかった場合は停止する。
- usecase、service、controller で concrete 実装を new しなかった場合は停止する。
- プロダクトテスト、fixture、snapshot、test helper を変更しなかった場合は停止する。
- docs / workflow 文書を変更しなかった場合は停止する。
