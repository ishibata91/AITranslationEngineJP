// 指示文（directive）の表示用の派生関数。型は directive-view.ts に置く。
import type {
  Directive,
  RecordAssignment,
  DirectiveSection
} from "./directive-view"

// 指示文ごとに、割り当てた REC:FIELD（対象）を束ねる。指示文の順序を保つ。
export function buildDirectiveSections(
  directives: Directive[],
  assignments: RecordAssignment[]
): DirectiveSection[] {
  return directives.map((d) => ({
    key: d.key,
    instruction: d.instruction,
    variables: d.variables,
    targets: assignments.filter((a) => a.directive === d.key)
  }))
}
