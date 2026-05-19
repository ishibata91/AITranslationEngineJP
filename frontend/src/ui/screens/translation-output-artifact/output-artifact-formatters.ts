export function formatCount(value: number | undefined): string {
  if (typeof value !== "number") {
    return "-"
  }

  return `${value.toLocaleString("ja-JP")} 件`
}

export function formatDistribution(
  distribution: Record<string, number> | undefined
): string {
  if (!distribution || Object.keys(distribution).length === 0) {
    return "-"
  }

  return Object.entries(distribution)
    .map(([key, value]) => `${key}: ${value}`)
    .join(" / ")
}

export function formatStatus(value: string | undefined): string {
  if (!value) {
    return "-"
  }

  return value.replaceAll("_", " ")
}
