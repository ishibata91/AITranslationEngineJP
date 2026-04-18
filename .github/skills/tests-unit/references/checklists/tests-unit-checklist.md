# Tests Unit Checklist

## Knowledge Check

- [ ] 1 test で 1 public behavior / branch / error path だけを証明した
- [ ] setup の clock、random、ID、repository 応答順序を固定した
- [ ] implementation_task_ids の外へ広げなかった

## Common Pitfalls

- [ ] test body に条件分岐を入れなかった
- [ ] test のためだけの product code 変更を広げなかった
- [ ] 新しい要件解釈を足さなかった
