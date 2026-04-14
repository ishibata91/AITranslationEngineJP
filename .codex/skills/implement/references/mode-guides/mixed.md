# Implement: mixed

## Focus

- scope freeze 済みの frontend-backend 横断変更を扱う

## Rules

- `implementation-scope` に従って `owned_scope` を守る
- 片側だけで閉じないことを scope artifact で確認する
- validation は両側の証跡を返す
