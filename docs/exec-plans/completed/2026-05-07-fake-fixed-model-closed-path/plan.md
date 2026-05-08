# fake provider DI と fake-model 取得契約

## 状態

- レーン: `light-change-lane`
- 現在成果物: `作業レポート入力`
- 次成果物: `人間確認`
- 作業計画フォルダ: `docs/exec-plans/active/2026-05-07-fake-fixed-model-closed-path/`

## 依頼要約

fake provider は user-facing provider ID ではない。
fake provider は通常 provider interface の実装を DI で差し替える、テスト用の偽物である。
Job Setup とマスターペルソナのモデル設定では、通常の model list 契約を通じて `fake-model` を選べる状態にする。

## 修正前提

- AIサービス設定の provider 一覧へ `fake` を追加しない。
- Job Setup の provider 一覧へ `fake` を追加しない。
- `translationJobSetupUserFacingProviderIDs` へ `fake` を追加しない。
- provider catalog へ user-facing provider として `fake` を追加しない。
- 翻訳管理タブ初期表示は本 task で変更しない。

## 既存のズレ

- 生成実行の fake mode は、通常 provider 実装を直接差し替えず、test-safe transport 差し替えで動いている。
- `ProviderFake` という provider ID が存在し、通常 provider interface の DI fake と user-facing provider ID の概念が混ざっている。
- model list 側には fake mode の DI 契約がなく、通常 provider ID の `getModels` を fake 実装で置き換える構造になっていない。
- 既存のズレを放置して `fake` provider ID を UI や Job Setup catalog へ広げると、利用者が `fake` を選ぶ構造になってしまう。
- 本 task は、`fake` provider ID を広げるのではなく、DI された fake 実装が通常契約へ `fake-model` を返す方向へ寄せる。

## 変更したい境界

- fake mode では、通常の `getModels` 契約が `fake-model` 1 件を返す。
- fake mode では、通常の生成 provider 契約が外部 HTTP へ出ず固定応答を返す。
- UI は provider ID を `fake` に切り替えず、受け取った model list から `fake-model` を選べる。
- マスターペルソナと Job Setup のモデルカードは、通常の model list 結果として `fake-model` を扱う。

## 成果物DAG

- `task 枠`: 完了。
- `軽量変更計画`: 完了。
- `実装証跡`: 完了。
- `人間確認`: 未着手。
- `テスト修正証跡`: 不要。プロダクトテスト変更なしで検証通過。
- `レビュー通過根拠`: 完了。
- `正本化判断`: 完了。human 承認済み恒久仕様は未確認。
- `詳細仕様正本反映`: 未実施。
- `作業レポート入力`: 完了。
- `作業計画完了移動`: 未実施。

## 停止中成果物

- `人間確認`: モデルカード実表示と操作確認が未入力。
- `詳細仕様正本反映`: human 承認済み恒久仕様が未確認。
- `作業計画完了移動`: 未完了成果物が残っているため未実施。

## 検証予定

- fake mode の `getModels` が `fake-model` 1 件を返すことを確認する。
- Job Setup のモデルカードで `fake-model` を選べることを確認する。
- マスターペルソナのモデルカードで `fake-model` を選べることを確認する。
- fake mode が外部 secret、実 provider、外部 HTTP に出ないことを確認する。

## 停止条件

fake を provider ID として UI や user-facing catalog に追加する必要が出た場合は停止する。
新しい公開 DTO、Wails method、DB schema が必要になった場合は停止する。
`getModels` 契約ではなく UI 専用分岐でしか実現できない場合は停止する。
