# Tests Scenario Checklist

## Knowledge Check

- [ ] scenario outcome を 1 test 1 outcome に分けた
- [ ] UI が入口のシナリオは、ユーザー入力の模倣を開始点にした
- [ ] happy path と failure path を別 test case にした
- [ ] runtime event 完了の観測点を明示した

## Common Pitfalls

- [ ] 新しい要件解釈を足さなかった
- [ ] test body に条件分岐を入れなかった
- [ ] UI 入口のシナリオ試験を裏側の直接呼び出しだけで代替しなかった
- [ ] paid real AI API を呼ばなかった
