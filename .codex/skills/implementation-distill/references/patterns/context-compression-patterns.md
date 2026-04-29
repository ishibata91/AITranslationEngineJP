# Context Compression Patterns

## 目的

`implementation_distiller` が single_handoff_packet 1 件を実装可能な lane_context_packet へ圧縮するための判断パターンをまとめる。
agent TOML の tool policy と contract の output obligation は上書きしない。

## 採用する考え方

- entry point、execution flow、architecture layer、dependency を分けて読む。
- 周辺 code を読む前に、single_handoff_packet と owned_scope を固定する。
- facts、inferred、gap を混ぜず、実装者が使う順に圧縮する。
- similar implementation を探し、既存 pattern を優先する。
- planning 情報は first_action、file path、symbol/type/function、line number、risk、validation entry を含める。
- 要件、実装方針、決定事項は source と implementation_implementer impact 付きで要約する。
- compression は token 削減ではなく、patch 生成に必要な fix_ingredients の保存を優先する。
- 追加が必要そうな method、interface、field は、実 code の present / absent 確認が終わるまで inferred として扱う。
- first_action は 1 completion_signal clause だけに対応させる。
- validation entry は broad gate ではなく、最初に実行する cheap check を優先する。

## 適用ルール

- Wails binding、frontend gateway、backend service / infra 境界を分けて記録する。
- code context は file / function / block の構造単位で扱い、構造を壊す断片化を避ける。
- patch 生成に必要な symbol、型、式、境界、call site を fix_ingredients に残す。
- 類似していても修正に不要な context は distracting_context に分け、required_reading に混ぜない。
- repository method、interface、field の新設を示す時は、該当定義を読んで absent fact を related_code_pointers に残す。
- first_action が 1 edit で clause を閉じられない場合は、同じ clause の最小 closure chain を change_targets に上流から leaf まで残す。
- `partial`、`multiple clauses`、`advance boundary` のような clause closure が曖昧な表現は使わない。
- existing_patterns が見つからない場合は、searched scope、searched layer、none の実装影響を添える。
- validation_entry は validation_commands から最小の lane-local command を選び、広い command しかない時はその理由を書く。
- `docs/architecture.md` と `docs/coding-guidelines.md` は必要な判断だけに圧縮する。
- handoff の文章を写すのではなく、implementation handoff に必要な制約へ変換する。
- 要件文書、実装方針、決定事項、out of scope、禁止事項は requirements_policy_decisions に要約する。
- required_reading に non-code 文書を残す場合は、要約済みでも原文確認が必要な理由を書く。
- required_reading の先頭は、implementation_implementer が最初に触る実 code の path、symbol、line number にする。
- related_code_pointers は path、symbol/type/function、line number、読み取った事実を 1 組で記録する。
- 実 code を読めず変更点が特定できない場合は、実装可能 packet にせず gaps に blocker として返す。
- 不足情報は `gaps` に残し、実装案で埋めない。

## 赤旗

- `required_reading` が広すぎて実装者の最初の一手が分からない。
- `required_reading` がファイル名の列挙で終わっている。
- `related_code_pointers` に symbol/type/function や line number がない。
- `fix_ingredients` がなく、なぜその context が patch に必要か分からない。
- 類似 context を distracting_context に分けず、implementation_implementer の読む対象に混ぜている。
- repository method が必要、と推測だけで書いている。
- first_action の clause_closed が partial 相当で、1 手目が何を閉じるか曖昧である。
- first_action が leaf contract だけ、または複数 clause をまとめて触っている。
- existing_patterns が none だけで、探索範囲と実装影響がない。
- validation entry が broad command だけで、cheap check の検討がない。
- 実 code を読まず handoff の文章を言い換えている。
- 要件、実装方針、決定事項の文書を required_reading に丸投げしている。
- implementation_implementer がどの file のどの関数を変更すべきか再調査する必要がある。
- inferred を fact として書いている。
- owned_scope 外の architecture tour が長い。
- validation entry がないまま implementation_implementer へ渡している。
