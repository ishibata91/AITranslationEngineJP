# Shared UI component

- 複数画面で同じ意味と操作を持つ表示部品だけを置く。
- 値は props で受け、操作は callback prop で返す。
- gateway、store、generated binding を直接参照しない。
- 画面固有の条件を増やして汎用化しない。
- 画面固有の部品は対象 screen の配下へ置く。
- component の状態と主要 variant を Storybook story で示す。
