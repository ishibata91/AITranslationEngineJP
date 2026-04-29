# Codex ワークフロー補助図

この file は補助図である。
live workflow の説明本文と判断基準の正本は [README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/README.md) とする。

Codex は設計を担当します。
Codex implementation lane は実装を担当します。

```mermaid
flowchart TD
    A[implement-lane]
    P[task folder plan.md]
    B[distill]
    C[investigate]
    U[ui-design.md]
    S[scenario-design.md]
    F[human review]
    G[implementation-scope.md]
    H[Codex implementation lane implement-lane]
    K[Codex implementation lane implementation-distill]
    M[human relay or close judgment]
    I[updating-docs]
    J[close]

    A --> P
    P --> B --> S
    P --> C --> S
    S --> U
    U --> S
    S --> F
    F --> G
    G --> H
    H --> K --> H
    H --> M
    M --> A
    M --> J
    A --> I --> J
    A --> J
```

## 位置づけ

この file は全体の向きを素早く確認するための補助図である。
live の role 境界、handoff、stop 条件、docs 正本化判断は [README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/README.md) を使う。

## 参照先

- Codex workflow 正本: [README.md](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/README.md)
- Codex implementation lane 実装入口: [implement-lane](/Users/iorishibata/Repositories/AITranslationEngineJP/.codex/skills/implement-lane/SKILL.md)
- docs 仕様入口: [docs/index.md](/Users/iorishibata/Repositories/AITranslationEngineJP/docs/index.md)
