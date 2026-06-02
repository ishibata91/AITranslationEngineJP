import { render } from "@testing-library/svelte"
import { describe, expect, test, vi } from "vitest"
import { tick } from "svelte"

import UntrackEffectProbe from "./UntrackEffectProbe.svelte"

// テスト専用コンポーネントの export 関数の型を定義する
type UntrackEffectProbeExports = {
  setTrackedTrigger: (value: string) => void
  setUntrackedValue: (value: string) => void
}

// untrack() でラップした $state の読み取りが $effect の依存に登録されないことを証明する回帰テスト。
// 証明対象: setCurrentPhasePage 内で untrack を使って currentPhasePage を読む変更により、
// ユーザー操作起点での currentPhasePage 変化が $effect を再実行しないことの不変条件。
// 検証方法: テスト専用の小さな Svelte コンポーネント UntrackEffectProbe を使い、
// onEffectRun コールバックの呼び出し回数で $effect 再実行の有無を観測する。
describe("untrack() による $effect 依存非登録", () => {
  test("untrack でラップした値の変化は $effect を再実行しない", async () => {
    // Arrange: onEffectRun の呼び出し回数と引数を記録する
    // trackedTrigger を固定し、untrackedValue だけを変化させる
    const onEffectRun = vi.fn()

    const { component } = render(UntrackEffectProbe, {
      props: { onEffectRun }
    })

    const probe = component as unknown as UntrackEffectProbeExports

    // 初回マウント時の $effect 実行を確認する
    await tick()
    expect(onEffectRun).toHaveBeenCalledTimes(1)
    expect(onEffectRun).toHaveBeenLastCalledWith("value-a")

    // Act: untrack でラップした値 (untrackedValue) だけを変化させる
    // trackedTrigger は変えない。$effect の依存に登録されていなければ再実行されない。
    probe.setUntrackedValue("value-b")
    await tick()

    // Assert: $effect は再実行されていない
    expect(onEffectRun).toHaveBeenCalledTimes(1)
  })

  test("tracked な値の変化は $effect を再実行する", async () => {
    // $effect の依存として登録されている trackedTrigger の変化は $effect を再実行する
    // untrack の効果を対照として確認する
    const onEffectRun = vi.fn()

    const { component } = render(UntrackEffectProbe, {
      props: { onEffectRun }
    })

    const probe = component as unknown as UntrackEffectProbeExports

    await tick()
    expect(onEffectRun).toHaveBeenCalledTimes(1)

    // Act: $effect の依存として登録された trackedTrigger を変化させる
    probe.setTrackedTrigger("changed")
    await tick()

    // Assert: $effect が再実行される
    expect(onEffectRun).toHaveBeenCalledTimes(2)
    // $effect 実行時の untrackedValue は変化していないので "value-a" のまま
    expect(onEffectRun).toHaveBeenLastCalledWith("value-a")
  })

  test("tracked な値の変化時に untrack で読んだ値は変化後の最新値を返す", async () => {
    // $effect が tracked な変化で再実行された場合、
    // untrack() は最新の untrackedValue を返すことを証明する
    const onEffectRun = vi.fn()

    const { component } = render(UntrackEffectProbe, {
      props: { onEffectRun }
    })

    const probe = component as unknown as UntrackEffectProbeExports

    await tick()

    // Act: untracked を先に変化させてから tracked で $effect を起動する
    probe.setUntrackedValue("value-b")
    probe.setTrackedTrigger("changed")
    await tick()

    // Assert: $effect が 1 回だけ追加実行され、untracked は最新値 "value-b" を返す
    expect(onEffectRun).toHaveBeenCalledTimes(2)
    expect(onEffectRun).toHaveBeenLastCalledWith("value-b")
  })
})
