# Staged Apply Flow

## 目的

`.codex` へ直接書けない環境では、反映済み file を `tmp/codex/files/` に置く。
人間が VSCode task または `scripts/codex/apply_tmp_codex.py` を実行して、正本へ上書き、追加、削除する。
通常 apply は基本的に人間が実行する。
Codex は staged file の作成と `--check-only` による final gate 確認までを担当する。
人間から明示指示がある場合だけ、Codex が通常 apply を試行できる。

この flow は緊急回避ではなく、最終反映前に差分を明示するための安全手順である。

## 配置

staging directory は task ごとに分けない。
`tmp/codex` を 1 回の apply queue として使う。

```text
tmp/codex/
├── README.md
├── deletion-rationale.md        # 削除がある時だけ置く
├── delete-paths.txt             # file 削除がある時だけ置く
└── files/
    └── .codex/
        └── <target-file>
```

`files/` 配下は、repo root からの相対 path をそのまま再現する。
script は `tmp/codex/files/<repo-relative-path>` を `<repo-root>/<repo-relative-path>` へ copy する。
現在の apply 対象は `.codex/` 配下だけに制限する。

## VSCode task

`.vscode/tasks.json` には汎用 task を 1 個だけ置く。
task の command は repo 側の Python script を直接指す。

```json
{
  "label": "codex: apply tmp/codex to .codex",
  "type": "shell",
  "command": "python3 ${workspaceFolder}/scripts/codex/apply_tmp_codex.py",
  "problemMatcher": []
}
```

## Final Gate

script は copy 前に次を確認する。

- `tmp/codex/files/` または `tmp/codex/delete-paths.txt` が存在する
- staged source の hash が copy 前後で変わらない
- 反映先と staged source の diff を表示する
- diff に削除行がある場合は停止する
- file 削除は `tmp/codex/delete-paths.txt` に列挙された path だけ許可する
- 記載削除または file 削除が必要な場合は削除対象、理由、削除後の正本確認先を記録する
- apply 成功後に `tmp/codex` を全削除する

削除行の自動検出は、空行以外の `diff -u` の `-` 行を対象にする。
単なる追記や新規 file の追加では停止しない。

## Check-only

`scripts/codex/apply_tmp_codex.py` は `--check-only` を持つ。
`--check-only` は copy せず、Final Gate と構文確認だけを実行する。

Codex は `.codex` へ直接書けない場合、この mode で staging の妥当性を確認する。
人間は確認後に通常実行する。
Codex は人間から明示指示がない限り、通常実行へ進まない。

## Deletion Rationale

記載を消す必要がある時は、`tmp/codex/deletion-rationale.md` を置く。
少なくとも次を書く。

- 消す file path
- 消す記載または範囲
- 消す理由
- 削除後の正本確認先、または不要になった根拠

file 自体を削除する時は、`tmp/codex/delete-paths.txt` に repo-relative path を 1 行ずつ書く。
script は削除行または file 削除を検出した時、`deletion-rationale.md` が対象 path を含まなければ停止する。
