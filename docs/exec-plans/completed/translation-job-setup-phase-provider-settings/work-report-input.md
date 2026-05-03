# Work Report Input: translation-job-setup-phase-provider-settings

## 結果

- `implementation_action`: `close`
- `final_validation`: `pass`
- `codex_review_result`: `pass`
- `docs_canonicalization`: `pending-human-approval`

## 実装要約

- Job Setup を master-persona provider 設定から切り離し、3 phase の provider / model / credential / batch mode を保存するようにした。
- provider model list は provider adapter 経由で取得し、API key 必須 provider は secret gate を通すようにした。
- LM Studio は API key 不要 provider として扱い、credential UI、secret load、credential missing を避けるようにした。
- `CredentialRef` は Job Setup が許可した参照だけに制限し、master-persona secret namespace への迂回参照を拒否するようにした。
- 古い model list 成功 / 失敗応答が現在 phase state を上書きしないようにした。

## 検証

- `python3 scripts/harness/run.py --suite backend-local`: pass
- `python3 scripts/harness/run.py --suite frontend-local`: pass
- `python3 scripts/harness/run.py --suite all`: pass
- `system test`: 5 passed
- `frontend coverage`: statements 65.5%, lines 65.5%
- `backend coverage`: statements 69.9%, lines 69.7%
- `Sonar`: coverage 71.1%, security 0, reliability 0, maintainability HIGH 0

## レビュー

- `reviewback.behavior.yaml`: `no_issue`
- `reviewback.contract.yaml`: `no_issue`
- `reviewback.trust-boundary.yaml`: `no_issue`
- `reviewback.state-invariant.yaml`: `no_issue`
- `reviewback.responsibility-boundary.yaml`: `no_issue`

## 残件

- UI は後続タスクでデザインオーバーホールする。
- 詳細仕様正本反映は human 承認後に docs-only 作業として判断する。
