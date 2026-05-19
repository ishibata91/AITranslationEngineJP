<script lang="ts">
  import ActionButton from "./ActionButton.svelte"
  import ButtonGroup from "./ButtonGroup.svelte"

  type Align = "start" | "end" | "stretch"

  interface Props {
    ariaLabel: string
    align?: Align
  }

  let {
    ariaLabel,
    align = "end"
  }: Props = $props()

  const noop = (): void => {}
</script>

<ButtonGroup {ariaLabel} {align}>
  <ActionButton label="戻る" variant="secondary" onClick={noop} />
  <ActionButton label="保存する" variant="primary" onClick={noop} />
  <ActionButton label="削除する" variant="danger" onClick={noop} />
</ButtonGroup>
