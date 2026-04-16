# 作業計画

- workflow: work
- status: completed
- lane_owner: implement
- scope: ai-provider-boundary-refactor
- task_id: ai-provider-boundary-refactor
- task_mode: refactor

## 依頼要約

- `internal/infra/ai` の AI provider 実装を persona から切り離す。
- AI provider の責務は provider-agnostic な AI request / response と外部 AI API 呼び出しだけにする。
- provider 実装を 1 ファイルに集約せず、client、transport、Gemini、OpenAI-compatible provider へ分割する。

## 判断根拠

- 人間から、AI provider が master persona に完全依存している点を修正するよう明示された。
- `internal/infra/ai` は infra adapter であり、persona prompt や persona body は service 側の責務である。
- `internal/infra/ai/master_persona_provider.go` は client、transport、provider concrete、JSON DTO、response parser を 1 ファイルに集約している。

## 実装スコープ

- `internal/infra/ai` から `MasterPersona` 命名と persona body 責務を除去する。
- `internal/infra/ai` に provider-agnostic な `ProviderRequest`、`ProviderResponse`、`ProviderClient` を置く。
- Gemini と OpenAI-compatible request / response 実装を分割する。
- service 側の master persona prompt 組み立てと body 生成 port は維持し、infra AI client への adapter だけを更新する。
- fake は provider option に戻さず、test-safe transport と env mode に閉じる。

## 非スコープ

- provider の種類追加。
- frontend の挙動変更。
- prompt template の内容変更。
- docs 正本の恒久仕様更新。

## 検証

- `go test ./internal/infra/ai ./internal/service ./internal/bootstrap`
- `npm run lint:backend`
- `python3 scripts/harness/run.py --suite structure`

## クローズ条件

- `internal/infra/ai` の public API と file 名から `MasterPersona` 責務が除去される。
- AI provider concrete は共通 response を返す。
- master persona 固有の body / prompt 責務は service または bootstrap adapter 側に残る。
- backend targeted tests と structure harness が pass する。

## 結果

- `internal/infra/ai/master_persona_provider.go` を廃止し、`provider.go`、`provider_client.go`、`transport.go`、`gemini.go`、`openai_compatible.go` へ分割した。
- `internal/infra/ai` の request / response / provider client API から `MasterPersona` 命名を除去した。
- master persona 固有の prompt と body 変換は `internal/service` と `internal/bootstrap` の adapter 境界に残した。
- coverage script は Go file rename 後の stale coverage path を避けるため、旧 coverage file 削除と `go test -count=1` を追加した。

## 検証結果

- `go test ./internal/infra/ai ./internal/service ./internal/bootstrap`: pass
- `npm run lint:backend`: pass
- `python3 scripts/harness/run.py --suite structure`: pass
- `python3 scripts/harness/run.py --suite coverage`: pass
- `python3 scripts/harness/run.py --suite all`: pass
- Sonar MCP: HIGH/BLOCKER open 0、reliability open 0、security open 0
- Sonar quality gate status: `NONE`
