# task 枠

## 人間依頼

- 既存の大量修正と中間成果物を破棄する。
- fake provider の理解を、通常 provider interface を DI で差し替える偽物として固定する。
- fake mode の model list は `fake-model` 1 件を返す。
- Job Setup とマスターペルソナのモデル設定で `fake-model` を選べる状態にする。

## 変更禁止範囲

- AIサービス設定の provider 一覧に fake を追加しない。
- Job Setup の user-facing provider 一覧に fake を追加しない。
- provider catalog へ user-facing provider として fake を追加しない。
- 翻訳管理タブ初期表示を変更しない。
- fake 専用 UI 分岐を増やして、通常 model list 契約を迂回しない。

## 確認したい結果

- fake mode で `getModels` を呼ぶと `fake-model` 1 件が返る。
- Job Setup のモデルカードで `fake-model` を選べる。
- マスターペルソナのモデルカードで `fake-model` を選べる。
- UI で provider を `fake` に切り替える必要がない。
- fake mode は外部 secret、実 provider、外部 HTTP に出ない。

## 判断メモ

既存 HEAD では、生成実行の fake mode は test-safe transport 差し替えで動いている。
一方で `ProviderFake` という provider ID も存在し、DI fake と provider ID fake が混ざっている。
本 task では user-facing provider ID としての `fake` を広げず、DI fake の model list 契約へ寄せる。

## 既存のズレ

- `ProviderFake` は通常 provider registry の 1 entry として存在する。
- fake mode は provider registry entry を差し替えず、transport を test-safe 実装へ差し替えている。
- model list 契約は fake mode の DI 実装を持たず、通常 provider の loader 経路を使う。
- そのため、現状は「fake provider ID」と「DI で差し替える fake mode」が同居している。
- 修正では `fake` ID を UI へ露出させず、fake mode の DI 契約として `getModels` を固定応答にする。
