# Review Patterns

## 目的

`reviewer` が single_handoff_packet、lane_context_packet、実装差分を照合するための判断パターンをまとめる。
agent contract の権限や output obligation は上書きしない。

## 採用する考え方

- diff だけでなく surrounding code、call site、dependency を読む。
- confidence が高い finding だけを出し、好みや未承認改善を混ぜない。
- severity は security / correctness / regression / maintainability の順に見る。
- AI-generated code では hidden coupling、architecture drift、edge-case 欠落を重点的に見る。
- similar finding はまとめ、実装者が直せる単位で返す。

## 適用ルール

- single_handoff_packet、lane_context_packet、owned_scope を review の正本にする。
- docs、`.codex`、`.github` workflow 文書の変更が混ざっていないか見る。
- backend では layer boundary、secret logging、error leakage、Sonar gate を見る。
- frontend では Wails gateway、state、console error、visible UI evidence を見る。
- coverage は Sonar-compatible coverage 70% 以上を pass 条件にする。
- harness は `python3 scripts/harness/run.py --suite all` の evidence を確認する。
- paid real AI API を validation / UI check が呼ばないことを確認する。

## 赤旗

- hardcoded secret、token、API key、個人情報 log がある。
- input validation、authorization、transaction / rollback、timeout が抜けている。
- empty catch、silent fallback、generic success が failure を隠している。
- Sonar Security / Reliability issue または Maintainability HIGH/BLOCKER が残っている。
- coverage 70% 未満、または harness all 未実行のまま pass 判定している。
- review finding が scope 外の好みや将来改善になっている。
