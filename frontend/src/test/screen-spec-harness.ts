import { describe, it } from "vitest"

export interface ScreenSpecCheck {
  id: string
  verify: () => void | Promise<void>
}

function duplicateIds(ids: readonly string[]): string[] {
  const seen = new Set<string>()
  const duplicates = new Set<string>()

  for (const id of ids) {
    if (seen.has(id)) duplicates.add(id)
    seen.add(id)
  }

  return [...duplicates].sort()
}

export function validateScreenSpecCoverage(
  specIds: readonly string[],
  checks: readonly ScreenSpecCheck[]
): void {
  const checkIds = checks.map((check) => check.id)
  const duplicateSpecIds = duplicateIds(specIds)
  const duplicateCheckIds = duplicateIds(checkIds)
  const specIdSet = new Set(specIds)
  const checkIdSet = new Set(checkIds)
  const missing = [...specIdSet].filter((id) => !checkIdSet.has(id)).sort()
  const extra = [...checkIdSet].filter((id) => !specIdSet.has(id)).sort()
  const problems: string[] = []

  if (missing.length > 0) problems.push(`不足: ${missing.join(", ")}`)
  if (extra.length > 0) problems.push(`余分: ${extra.join(", ")}`)
  if (duplicateSpecIds.length > 0) {
    problems.push(`画面仕様側の重複: ${duplicateSpecIds.join(", ")}`)
  }
  if (duplicateCheckIds.length > 0) {
    problems.push(`単体テスト側の重複: ${duplicateCheckIds.join(", ")}`)
  }

  if (problems.length > 0) {
    throw new Error(`画面仕様IDの対応が一致しない。${problems.join(" / ")}`)
  }
}

export function runScreenSpecHarness(
  screenName: string,
  specIds: readonly string[],
  checks: readonly ScreenSpecCheck[]
): void {
  describe(`${screenName}の画面仕様`, () => {
    it("画面仕様IDと検証関数が一対一で対応する", () => {
      // 仕様または検証関数の追加漏れと、IDの重複を画面単位で検出する。
      validateScreenSpecCoverage(specIds, checks)
    })

    for (const check of checks) {
      it(`[${check.id}]`, async () => {
        // Autodocsに表示する同じ仕様IDの公開表示または操作可否を検証する。
        await check.verify()
      })
    }
  })
}
