# Tests: scenario-implementation

## Goal

- Scenario artifact または fix 再現条件をそのまま証明する

## Rules

- happy path と主要 failure path を含める
- fixture や helper は scenario を支える範囲に限定する
- runtime event 完了は completion event を観測点にする
