# backend database

- `db/migrations/` は SQL schema の正本とする。
- migration は repository-owned とする。
- migration の適用は `db.Apply` が担当する。
- `store/` は起動時の migration 適用を `db.Apply` に委譲する。
