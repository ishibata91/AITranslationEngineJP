# integration-followup-result

## 変更内容

- `applyProviderSettingsExecutionSnapshot` の test-safe 判定から endpoint 必須条件を外した。
- test-safe loader の許可時は、endpoint が nil でも `CredentialReferenceID=nil`、`CredentialState=not_required`、`ErrorKind=nil` にする。
- real provider の endpoint missing と credential missing の gate は維持した。
- frontend production code と fake provider ID は変更していない。

## テスト

- 成功: `go test ./internal/service -run 'ProviderSettings' -count=1`
- 失敗: `python3 scripts/harness/run.py --suite backend-local`
- 失敗箇所: `internal/bootstrap`
- 失敗内容: master persona AI settings persistence 2 件で、期待値と実値の model が `fake-model` になっている既存 bootstrap 差分。

## 残留リスク

- `backend-local` 全体は非通過である。
- 非通過箇所は `internal/bootstrap` であり、今回の `internal/service` 修正範囲ではない。
- task 内の他者 dirty diff は戻していない。
