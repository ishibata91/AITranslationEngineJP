import { describe, expect, it } from "vitest"
import { validateScreenSpecCoverage } from "./screen-spec-harness"

const check = (id: string) => ({ id, verify: () => {} })

describe("validateScreenSpecCoverage", () => {
  it("画面仕様IDと検証関数が一対一なら通す", () => {
    // 同じ二つのIDが一件ずつある正常な対応を受理する。
    expect(() =>
      validateScreenSpecCoverage(
        ["screen.one", "screen.two"],
        [check("screen.one"), check("screen.two")]
      )
    ).not.toThrow()
  })

  it("検証関数がない画面仕様IDを不足として失敗にする", () => {
    // 仕様側だけに存在するIDを取りこぼさない。
    expect(() =>
      validateScreenSpecCoverage(
        ["screen.one", "screen.missing"],
        [check("screen.one")]
      )
    ).toThrow("不足: screen.missing")
  })

  it("画面仕様にない検証関数のIDを余分として失敗にする", () => {
    // テスト側だけに存在するIDを仕様の対応として数えない。
    expect(() =>
      validateScreenSpecCoverage(
        ["screen.one"],
        [check("screen.one"), check("screen.extra")]
      )
    ).toThrow("余分: screen.extra")
  })

  it("画面仕様側の重複IDを失敗にする", () => {
    // 配列内の同じ仕様IDを集合化する前に検出する。
    expect(() =>
      validateScreenSpecCoverage(
        ["screen.duplicate", "screen.duplicate"],
        [check("screen.duplicate")]
      )
    ).toThrow("画面仕様側の重複: screen.duplicate")
  })

  it("単体テスト側の重複IDを失敗にする", () => {
    // 配列内の同じ検証関数IDを上書きせずに検出する。
    expect(() =>
      validateScreenSpecCoverage(
        ["screen.duplicate"],
        [check("screen.duplicate"), check("screen.duplicate")]
      )
    ).toThrow("単体テスト側の重複: screen.duplicate")
  })
})
