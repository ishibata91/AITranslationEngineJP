# 人間観測記録

## 観測

- 対象: ジョブセットアップのモデル選択。
- 失敗: APIサービス設定値に関わらず、モデル一覧が表示されない。
- 期待: fake provider では、一覧取得が走れば `fake-model` が表示される。
- 判断: 人間観測済み不具合として、既存 `fix-lane` task の再発確認に扱う。

## 境界

- frontend は fake provider ID や `fake-model` を特別扱いしない。
- service 層へ fake mode 判定を戻さない。
- provider catalog へ fake provider を追加しない。
- 修正はモデル一覧取得が backend に到達する経路へ限定する。
