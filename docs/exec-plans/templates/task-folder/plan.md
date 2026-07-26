# Task Plan: <task-id>

`plan.md` は branch 情報と、人間が見た事象、そこから起こした要求を持つ。
設計は `design.md`、確定仕様は `spec.md`、恒久的に残す判断は `docs/changelog.md` が持つ。

## 事象

人間が見たことだけを書く。何をしたいかは要求欄が持つ。

- <人間が見た画面、操作、ログ、失敗、期待との差分を 1 件 1 行で書く。>

## 要求

事象から起こした、何をどうするかを書く。人間の合意を得るまで空のままにし、空の状態で `design.md` と `spec.md` へ進まない。

`design.md` と `spec.md` は、ここに書いた番号ごとに節を分ける。

- **R-1 <要求の見出し文>**: <何を、どうするかを 1 文で書く。事象の言い換えにしない。>

## branch 情報

- `execution_branch`: `claude/<task-id>`
- `target_branch`: `master`
- `source_commit`: <分岐元 commit hash>

## やらないこと

- <この task で扱わない範囲を書く。要求に要る手段が対象として残っているかを確かめる。>
