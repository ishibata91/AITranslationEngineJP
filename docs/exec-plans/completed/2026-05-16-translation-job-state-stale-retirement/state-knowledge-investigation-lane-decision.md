# State Knowledge Investigation Lane Decision

- `caller`: `light_change_lane`
- `status`: `stopped_for_human_decision`
- `source`: `state-knowledge-investigation.md`

## 結論

現在の削減範囲だけでは、翻訳ジョブ状態関連の stale 廃止は閉じない。
`人間確認` とレビューへ進める前に、追加判断が必要である。

## 観測済みの不足

- `pending`: `docs/spec.md` の正本 state には存在しないが、3 phase service に内部 state として残っている。
- `commonPhaseActionAvailability`: `TranslationJobPolicy` の共通操作規則と同じ知識を service 層へ重複させている。
- `StateMachine`: product code からは外れたが、active observability task-local に旧名参照が残っている。
- `JobIOService`: architecture 正本と lint component に残るが、実体 package は `doc.go` だけである。
- `cancelled`: `PersonaGenerationPhaseContractStub` に正本 spelling の `canceled` と違う fixture state が残っている。

## 今すぐ進めない成果物

- `人間確認`: 削減範囲が未確定なので停止する。
- `テスト修正証跡`: 追加変更範囲が未確定なので停止する。
- `実装後ブラウザ確認`: UI 変更なしでも、完了範囲が未確定なので停止する。
- `レビュー通過根拠`: レビュー対象差分が未確定なので停止する。

## 人間判断が必要な論点

- `pending` を canonical state へ昇格するか、内部一時 state として隔離するか。
- `JobIOService` を architecture 正本から外すか、別 task で実体化するか。
- `observability-log-addition` の `StateMachine` / `JobIOService` 旧名参照を今回の更新対象へ含めるか。
- `TranslationJobPolicy` の共通操作規則を read model の操作可否へどう再利用するか。
- `cancelled` fixture spelling を今回の stale 廃止に含めるか。

## 戻し先

`light_change_lane` はここで停止する。
追加設計判断が必要な場合は `designer` または人間判断へ戻す。
軽量変更として継続する場合は、上記論点のうち今回含める範囲を人間が指定する必要がある。
