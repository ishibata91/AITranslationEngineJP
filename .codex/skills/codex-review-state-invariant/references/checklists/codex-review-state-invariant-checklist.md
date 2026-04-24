# Codex Review State Invariant Checklist

- [ ] transaction、lock、idempotency、retry を確認した
- [ ] race condition と DB 更新順序を確認した
- [ ] cache invalidation と event 発行を確認した
- [ ] partial failure、soft delete、集計値を確認した
- [ ] 二重作成または二重課金の可能性を確認した
