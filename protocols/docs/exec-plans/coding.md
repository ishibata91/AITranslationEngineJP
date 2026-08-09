# Exec plan

`design-workflow` と `fix-workflow` が作る `docs/exec-plans/**/<task-id>/` の文書は、現在の task を実装へ渡す契約として分ける。

- `plan.md` は人間と合意した要求だけを持つ。設計、仕様、未決定事項、参照、判断履歴を書かない。
- `design.md` は要求ごとの as-is と to-be の方針だけを持つ。source の path、symbol、外部資料、未決定事項、仕様を書かない。
- `spec.md` は要求ごとの観測可能で検証可能な仕様だけを持つ。方針、参照、未決定事項、満たさない部分を書かない。
- `pending.md` は task 固有の未決定事項とブロッカーだけを持つ。各項目には解決すべき問い、ブロックする段、解決に必要な入力を書く。解決した項目は結論を正本へ反映してから削除する。
- `references.md` は source の path、symbol、外部資料の URL または識別子だけを持つ。各参照に `REF-<番号>` を振る。解釈、設計理由、結論を書かない。
- `log.jsonl` は監査用の追記履歴だけを持つ。人間が task の履歴確認を明示的に依頼した場合だけ読む。通常の設計、仕様、実装、レビュー、コンテキスト再取得、状況確認のために読まない。

`pending.md` が空でない task は design-review、設計HITL、implementation-workflow へ進めない。未決定事項を `design.md` または `spec.md` に残したまま進めない。

`design.md`、`spec.md`、`implementation.md` に source の path、symbol、外部資料を直接書かない。必要な参照は `references.md` の `REF-<番号>` で示す。

## log.jsonl

1 行を 1 件の状態遷移 event とする。情報共有、経過メモ、結論本文は event に書かない。

event の field は次だけとする。未定義の field を追加しない。

| field | 責務 | 形式 |
| --- | --- | --- |
| `at` | event を記録した日時 | 秒まで持ち、UTC offset を持つ ISO 8601 文字列 |
| `actor` | event を記録した主体 | `human`、`codex`、`agent:<role>` のいずれか |
| `type` | 起きた状態遷移 | 下表のいずれか |
| `documents` | event が対象とする task folder 内の正本 | 空でない `plan.md`、`design.md`、`spec.md`、`pending.md`、`references.md`、`implementation.md`、`investigation.md`、`storybook-review-loop.md` の file 名配列 |
| `pending_id` | 追加または解決した pending 項目 | `P-<正の整数>`。`pending_added` と `pending_resolved` だけで必須 |

`type` は次だけとする。

| type | 記録する状態遷移 |
| --- | --- |
| `pending_added` | 未決定事項またはブロッカーを `pending.md` へ追加した |
| `pending_resolved` | pending 項目の結論を正本へ反映して `pending.md` から削除した |
| `design_review_passed` | design-review が通過した |
| `design_review_rejected` | design-review が否決した |
| `design_approved` | 人間が設計を承認した |
| `implementation_review_passed` | implementation-review が通過した |
| `implementation_review_rejected` | implementation-review が否決した |
| `implementation_approved` | 人間が実装を承認した |

設計、仕様、要求の結論は event へ複製せず、対応する正本へ書く。`documents` は読んだ file ではなく、event が対象とする正本だけを指す。

通常の追記は `jq -cn` で 1 行を作る。追記前に schema を検査し、`log.jsonl` を読まずに末尾へ追加する。

```sh
event="$(jq -cn \
  --arg at "$(date -Iseconds)" \
  --arg actor 'codex' \
  --arg type 'pending_resolved' \
  --arg pending_id 'P-1' \
  --argjson documents '["plan.md", "pending.md"]' \
  '{at: $at, actor: $actor, type: $type, pending_id: $pending_id, documents: $documents}')"

jq -e '
  (.at | type == "string" and test("^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(\\.[0-9]+)?(Z|[+-][0-9]{2}:[0-9]{2})$")) and
  (.actor | test("^(human|codex|agent:[a-z0-9_-]+)$")) and
  (.type | IN("pending_added", "pending_resolved", "design_review_passed", "design_review_rejected", "design_approved", "implementation_review_passed", "implementation_review_rejected", "implementation_approved")) and
  (.documents | type == "array" and length > 0 and all(.[]; IN("plan.md", "design.md", "spec.md", "pending.md", "references.md", "implementation.md", "investigation.md", "storybook-review-loop.md"))) and
  (if (.type == "pending_added" or .type == "pending_resolved")
   then ((keys | sort) == ["actor", "at", "documents", "pending_id", "type"] and (.pending_id | test("^P-[1-9][0-9]*$")))
   else (keys | sort) == ["actor", "at", "documents", "type"]
   end)
' <<< "$event" >/dev/null && print -r -- "$event" >> log.jsonl
```

人間が task の履歴確認を明示的に依頼した場合だけ、次のように読む。履歴確認を依頼されていない場合は、task が長期間続いた後でも `log.jsonl` を読まない。

```sh
jq -c 'select(.type == "pending_resolved")' log.jsonl
```
