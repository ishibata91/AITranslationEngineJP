---
name: explore
description: 人間報告用の調査 skill。関連 docs、コード、差分、evidence を読み、判断材料になる事実だけを短く返す。
---

# Explore

## Output

- confirmed facts
- relevant evidence
- open questions
- recommended next step

## Rules

- 調査は read-only に保つ
- 人間報告で必要な事実と evidence だけを返す
- 推測を事実として確定しない
- 実装、fix、`updating-docs` は直接始めない
- 不足情報が大きい場合は次に誰へ handoff すべきかを明示する
