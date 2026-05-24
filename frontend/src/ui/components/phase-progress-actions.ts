import type { PhaseActionItem } from "./phase-panel-types"

interface PhaseProgressActionSelectionOptions<ActionId extends string> {
  hiddenActionIds?: readonly ActionId[]
  runActionIds?: readonly ActionId[]
  alwaysActionIds?: readonly ActionId[]
}

export function selectPhaseProgressActions<ActionId extends string>(
  actions: PhaseActionItem<ActionId>[],
  {
    hiddenActionIds = [],
    runActionIds = [],
    alwaysActionIds = []
  }: PhaseProgressActionSelectionOptions<ActionId>
): PhaseActionItem<ActionId>[] {
  const hiddenActionIdSet = new Set(hiddenActionIds)
  const runActionIdSet = new Set(runActionIds)
  const alwaysActionIdSet = new Set(alwaysActionIds)
  const visibleActions = actions.filter(
    (action) => !hiddenActionIdSet.has(action.id)
  )
  const enabledActions = visibleActions.filter((action) => !action.disabled)
  const runAction = enabledActions.find((action) =>
    runActionIdSet.has(action.id)
  )
  const disabledRunFallbackAction = visibleActions.find((action) =>
    runActionIdSet.has(action.id)
  )
  const selectedActionIds = new Set<ActionId>()
  const selectedActions: PhaseActionItem<ActionId>[] = []
  const selectedRunAction = runAction ?? disabledRunFallbackAction
  if (selectedRunAction) {
    selectedActions.push(selectedRunAction)
    selectedActionIds.add(selectedRunAction.id)
  }

  selectedActions.push(
    ...visibleActions.filter(
      (action) =>
        !selectedActionIds.has(action.id) &&
        !runActionIdSet.has(action.id) &&
        (!action.disabled || alwaysActionIdSet.has(action.id))
    )
  )

  return selectedActions
}
