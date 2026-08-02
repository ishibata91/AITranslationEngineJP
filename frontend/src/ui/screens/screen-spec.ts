export interface ScreenSpec {
  id: string
  statement: string
}

export interface ScreenState<TArgs> {
  storyName: string
  precondition: string
  args: TArgs
  specs: readonly ScreenSpec[]
}

export function defineScreenState<
  const TArgs,
  const TSpecs extends readonly ScreenSpec[]
>(
  state: ScreenState<TArgs> & { specs: TSpecs }
): ScreenState<TArgs> & {
  specs: TSpecs
} {
  return state
}

export function screenStateDescription<TArgs>(
  state: ScreenState<TArgs>
): string {
  const specs = state.specs
    .map((spec) => `- \`${spec.id}\`: ${spec.statement}`)
    .join("\n")

  return `### 前提条件

${state.precondition}

### 画面仕様

${specs}`
}

export function screenSpecIds<TArgs>(
  states: readonly ScreenState<TArgs>[]
): string[] {
  return states.flatMap((state) => state.specs.map((spec) => spec.id))
}
