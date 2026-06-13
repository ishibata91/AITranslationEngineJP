import { describe, it, expect } from "vitest"
import { readdirSync, readFileSync, statSync } from "node:fs"
import { join } from "node:path"
import { fileURLToPath } from "node:url"

// coding-guidelines-frontend §6: generated wailsjs の import は gateway 境界にだけ閉じ込める。
// View・Container・component から generated bindings を直接参照しないことを保証する。
const srcDir = fileURLToPath(new URL("./src", import.meta.url))

function walk(dir) {
  const files = []
  for (const name of readdirSync(dir)) {
    const path = join(dir, name)
    if (statSync(path).isDirectory()) files.push(...walk(path))
    else if (/\.(ts|svelte)$/.test(name)) files.push(path)
  }
  return files
}

describe("Wails 境界", () => {
  it("generated wailsjs の import は gateway だけに閉じる", () => {
    const offenders = []
    for (const file of walk(srcDir)) {
      if (file.includes(`${join("src", "gateway")}`)) continue
      const content = readFileSync(file, "utf8")
      if (/from\s+["'][^"']*wailsjs/.test(content)) {
        offenders.push(file.slice(srcDir.length + 1))
      }
    }
    expect(offenders).toEqual([])
  })
})
