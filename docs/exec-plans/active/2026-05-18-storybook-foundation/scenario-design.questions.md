# Human Decision Questionnaire

## 状態

- `status`: answered
- `source`: `./scenario-design.md`
- `reason`: `Q-SBF-001` から `Q-SBF-004` までの人間回答を反映済みである。
- `answered_by`: human

## 回答済み決定

- `Q-SBF-001`: Storybook build は Storybook 専用 gate に分ける。理由は、厳しく harness で制限する必要を感じないためである。
- `Q-SBF-002`: Storybook fixture は component 横の `frontend/src/ui/**/__fixtures__` に置く。理由は、component 単位の見た目検証用であり、業務データを Storybook で扱わないためである。
- `Q-SBF-003`: Storybook review URL と確認状態は、専用 `storybook-review.md` に記録する。
- `Q-SBF-004`: Storybook 運用は、後続 task の plan と POC task が成功した後に、skill、agent、docs へ反映する。

### [Q-SBF-001] Storybook build gate の扱い

決める仕様:
Storybook build を既存 `frontend-local` gate に含めるか、Storybook 専用 gate として分けるかを決める。

決定済み:
`npm --prefix frontend run build-storybook` は Storybook 最小基盤の検証入口である。

未確定:
`plan.md` が `frontend-local` に含めるか別 gate にするかを未決事項にしている。

選択肢:
1. `frontend-local` に含める。通常 frontend 検証で Storybook build も必ず確認する。導入直後の失敗は frontend-local 全体を止める。
2. Storybook 専用 gate に分ける。Storybook 基盤の失敗を専用 command で追える。通常 frontend-local とは別に実行記録が必要になる。
3. 当面 task-local 検証だけにする。この task では build-storybook を実行するが、恒久 gate への組み込みは docs 正本化後に回す。
4. その他

AI 推奨:
2. Storybook 専用 gate に分ける。Storybook 導入直後は通常 frontend gate と失敗原因を分けて観測する方が追跡しやすい。恒久 gate へ入れる判断は、build が安定してから正本化できる。

人間回答:
2. Storybook 専用 gate に分ける。理由は、厳しく harness で制限する必要を感じないためである。

### [Q-SBF-002] fixture 配置 root

決める仕様:
Storybook 用 fixture を UI 近傍に置くか、Storybook 専用 directory に置くかを決める。

決定済み:
fixture は fixed props または view model fixture とし、gateway mock、backend DTO mock、実行フロー再現は使わない。

未確定:
`plan.md` が `frontend/src/ui/**/__fixtures__` と Storybook 専用 directory のどちらに置くかを未決事項にしている。

選択肢:
1. `frontend/src/ui/**/__fixtures__` に置く。component 近傍で fixture を探せる。UI 階層の移動時に fixture も一緒に移動する。
2. Storybook 専用 directory に置く。Storybook 用データを一箇所で管理できる。component から fixture までの距離は伸びる。
3. 最小 story だけ近傍に置き、共通化は後続 task で決める。今回の基盤は小さく閉じられる。後続 Master Persona task で配置再判断が必要になる。
4. その他

AI 推奨:
1. `frontend/src/ui/**/__fixtures__` に置く。この repo の frontend 規約は画面専用 component と shared component の境界を重視している。fixture も表示単位の近くに置く方が、backend mock へ広がりにくい。

人間回答:
1. `frontend/src/ui/**/__fixtures__` に置く。理由は、component 単位の見た目検証用であり、業務データを Storybook で扱わないためである。

### [Q-SBF-003] review URL 記録 artifact

決める仕様:
Storybook review URL と確認状態を、task-local のどの artifact に記録するかを決める。

決定済み:
この task では `ui-design.md` と `screen-design-diff.*.md` は N/A であり、review URL は Storybook 表示確認先だけを指す。

未確定:
候補成果物では `ui-design.md`、専用 artifact、実装結果 artifact の候補が混在している。

選択肢:
1. 専用 `storybook-review.md` に記録する。Storybook URL と確認状態を独立して追跡できる。task-local artifact が 1 つ増える。
2. 実装結果 artifact に記録する。実行結果と review URL を同じ場所で読める。Storybook URL だけを探す時の入口は弱くなる。
3. `plan.md` の実装結果欄に記録する。既存 plan だけで完結する。plan が実行証跡で膨らみやすい。
4. その他

AI 推奨:
1. 専用 `storybook-review.md` に記録する。`ui-design.md` が N/A の task でも Storybook review URL を見失いにくい。後続 Master Persona task も同じ形式を参照しやすい。

人間回答:
1. 専用 `storybook-review.md` に記録する。

### [Q-SBF-004] Storybook 運用の docs 正本化先

決める仕様:
Storybook 運用を後でどの docs 正本へ反映するかを決める。

決定済み:
この task では docs 正本を変更しない。docs 正本化が必要な場合は人間承認後に `updating-docs` へ戻す。

未確定:
`plan.md` が Storybook 運用をどの docs 正本へ反映するかを未決事項にしている。

選択肢:
1. `docs/tech-selection.md` に採用技術として反映する。Storybook を frontend tooling の採用技術として明示できる。運用手順の詳細は別文書が必要になる可能性がある。
2. `docs/lint-policy.md` に gate として反映する。Storybook build の検証責務を明示できる。採用技術や review URL 運用は別途説明が必要になる。
3. まず task-local に留め、後続 task 後に正本化する。基盤導入直後の判断を急がない。恒久ルール化は Master Persona story 追加後に回る。
4. その他

AI 推奨:
3. まず task-local に留め、後続 task 後に正本化する。Storybook 基盤だけでは、採用技術、gate、review URL 運用のどこまで恒久化するか判断材料が薄い。Master Persona story 追加後の運用実績を見て正本化先を決める方がよい。

人間回答:
4. 後続 task の plan と POC task が成功したら、skill、agent、docs に反映する形にする。
