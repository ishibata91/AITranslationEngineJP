class GatewayResponseShapeError extends Error {
  constructor(
    public readonly userFacingMessage: string,
    public readonly internalDiagnostic: string
  ) {
    super(userFacingMessage)
    this.name = "GatewayResponseShapeError"
  }
}

type RuntimeShapeIssue = {
  path: string
  expected: string
}

export function createGatewayResponseShapeError(
  bindingName: string,
  issues: RuntimeShapeIssue[]
): GatewayResponseShapeError {
  const details = issues
    .map((issue) => `${issue.path}: expected ${issue.expected}`)
    .join("; ")
  return new GatewayResponseShapeError(
    "Gateway response shape is invalid.",
    `${bindingName} returned invalid DTO shape. ${details}`
  )
}

export function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value)
}

export function isString(value: unknown): value is string {
  return typeof value === "string"
}

export function isNumber(value: unknown): value is number {
  return typeof value === "number" && Number.isFinite(value)
}

export function isBoolean(value: unknown): value is boolean {
  return typeof value === "boolean"
}

export function isOptionalString(value: unknown): value is string | undefined {
  return value === undefined || isString(value)
}

export function isOptionalNumber(value: unknown): value is number | undefined {
  return value === undefined || isNumber(value)
}

export function isArrayOf<T>(
  value: unknown,
  isItem: (item: unknown) => boolean
): value is T[] {
  return Array.isArray(value) && value.every((item) => isItem(item))
}
