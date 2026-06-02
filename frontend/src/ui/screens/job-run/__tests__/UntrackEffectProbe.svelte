<script lang="ts">
  import { untrack } from "svelte"

  interface Props {
    // $effect が実行されるたびに呼ばれるコールバック。テスト側で実行回数を計測する。
    onEffectRun: (capturedUntrackedValue: string) => void
  }

  let { onEffectRun }: Props = $props()

  // コンポーネント内部の $state として管理する。
  // テスト側から exported 関数で直接 $state を変化させることで、
  // @testing-library/svelte の rerender による props 全体差し替えを避ける。
  let trackedTrigger = $state("initial")
  let untrackedValue = $state("value-a")

  $effect(() => {
    // trackedTrigger を直接読む → $effect の依存に登録される
    void trackedTrigger

    // untrackedValue を untrack でラップして読む → $effect の依存に登録されない
    const captured = untrack(() => untrackedValue)

    onEffectRun(captured)
  })

  // テスト側から $state を直接変化させるための exported 関数
  export function setTrackedTrigger(value: string): void {
    trackedTrigger = value
  }

  export function setUntrackedValue(value: string): void {
    untrackedValue = value
  }
</script>

<span data-testid="probe-tracked">{trackedTrigger}</span>
<span data-testid="probe-untracked">{untrackedValue}</span>
