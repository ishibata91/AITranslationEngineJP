# 人間観測記録

## 観測対象

- 対象: モデル設定カードを使う画面全体
- 主な経路: AIサービス設定、翻訳ジョブ設定、マスターペルソナ
- 関連: fake mode の model list 取得と model 選択可否

## 観測事実

- AIサービス設定が設定済み表示でも、credential / secret 判定で「AIサービス設定が未完了です。」へ入る。
- API key へ `hogefuga` のような正しくない値を入れても、fake mode なら実 provider へ出ない。
- fake mode なら、credential と endpoint の中身に関係なく、AI model を使える状態になるべきである。

## 期待との差分

- 期待: fake mode の許可判定は backend の test-safe loader / transport で完結する。
- 期待: frontend は fake mode や `fake-model` を知らず、backend の model list 結果だけを表示する。
- 差分: frontend の credential preflight や backend の endpoint preflight が、fake mode の model list 取得前に止める可能性がある。

## 禁止事項

- frontend へ fake mode 判定を追加しない。
- frontend へ `fake-model` 固有分岐を追加しない。
- `fake` provider ID を user-facing provider として追加しない。
- credential や endpoint の取得可否を fake 判定に使わない。
