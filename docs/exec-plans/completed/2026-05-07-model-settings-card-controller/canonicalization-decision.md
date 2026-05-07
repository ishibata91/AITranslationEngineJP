# 正本化判断

## 状態

- `artifact`: `正本化判断`
- `status`: `completed`
- `detail_spec_sync`: `stopped`
- `source_plan`: `./plan.md`

## 判断

- モデル設定カードの provider、model、model list、保存、取得、選択状態の集約は恒久仕様候補である。
- `docs/detail-specs/translation-job-setup.md` は既に Job Setup の model list API 取得を記録している。
- `docs/detail-specs/ai-provider-settings-management.md` は Job Setup と master-persona が provider settings を参照することを記録している。
- master persona 側の model list 更新経路と共有モデル設定カード集約は、詳細仕様へ追記候補である。

## 停止理由

- `docs/index.md` は、`docs/` 正本は human が先に更新し、agent は human が直接起動した `updating-docs` でだけ同期すると定めている。
- 現在の呼び出しは `implement-lane` であり、`updating-docs` の直接起動ではない。
- そのため、詳細仕様正本反映はこの task では停止する。

## 次に更新する候補

- `docs/detail-specs/ai-provider-settings-management.md`: master-persona が保存済み provider settings を参照して model list を取得すること。
- `docs/detail-specs/translation-job-setup.md`: 共有モデル設定カードの責務境界が UI 部品ではなく application 側にあること。
- `docs/frontend-fake-api.md`: 必要なら fakeAPI success / error 状態の人間レビュー利用例。

