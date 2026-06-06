/**
 * boundary-golden-loader.ts
 *
 * frontend 側の共有 golden 読み込み helper。
 * vitest 実行時に Node.js fs を使い、backend 側と同一の物理ファイル
 * (internal/apitest/testdata/boundary/master_dictionary/*.golden.json) を読む。
 *
 * 両側が同一ファイルを参照することで、片側のみ golden を書き換えた場合に
 * もう片側のテストが失敗する検出経路が成立する。
 */

import { readFileSync } from "node:fs"
import { fileURLToPath, URL } from "node:url"

/**
 * BoundaryGoldenInput は golden ファイルの入力相当の説明を表す。
 */
export interface BoundaryGoldenInput {
  description: string
  method: string
  case: string
  request: unknown
}

/**
 * BoundaryGoldenFile は golden JSON ファイルのトップレベル構造。
 * input は入力相当の説明、expected は期待される応答 DTO field 値。
 */
export interface BoundaryGoldenFile<TExpected = unknown> {
  input: BoundaryGoldenInput
  expected: TExpected
}

/**
 * goldenBasePath は frontend/src/controller/wails/ から
 * internal/apitest/testdata/boundary/master_dictionary/ への相対パスを
 * import.meta.url を使って絶対パスに解決した基底ディレクトリ。
 */
const goldenBasePath = fileURLToPath(
  new URL(
    "../../../../internal/apitest/testdata/boundary/master_dictionary",
    import.meta.url
  )
)

/**
 * loadBoundaryGolden は指定した名前の golden ファイルを読み込み、
 * TExpected 型にキャストした BoundaryGoldenFile を返す。
 *
 * @param name - ファイル名のみ（例: "list_normal.golden.json"）
 */
export function loadBoundaryGolden<TExpected = unknown>(
  name: string
): BoundaryGoldenFile<TExpected> {
  const filePath = `${goldenBasePath}/${name}`
  const raw = readFileSync(filePath, "utf-8")
  return JSON.parse(raw) as BoundaryGoldenFile<TExpected>
}
