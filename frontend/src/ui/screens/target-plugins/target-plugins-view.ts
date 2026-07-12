// 翻訳対象 plugin 画面の表示用の型。値（表示定数・純関数）は target-plugins-presentation.ts に置く。
// 型と値を分けるのは、import 規約（値と型を同一 import で混在させない）に沿うため。
//
// 翻訳対象 plugin は、翻訳を実行した plugin 1 つを 1 行にした表示単位。plugin と 1 対 1。
// 一覧で成果の量を示し、削除でその plugin の成果を消して再実行でやり直す。

// 一覧の 1 行。翻訳を実行した plugin 1 つを表示用に整形した値。
export interface TargetPluginRow {
  // plugin のファイル名（例: Dawnguard.esm）。翻訳対象の識別子。
  plugin: string
  // 作成時刻の表示ラベル。時刻の整形は Presenter が行い、ここには表示文字列だけ入る。
  createdLabel: string
  // この plugin が束ねる翻訳対象の総行数（叙述文＋台詞＋固有名）。成果の量を示す。
  total: number
  // 訳済み行数（未訳を除く）。総行数に対する進捗を示す。
  translated: number
}
