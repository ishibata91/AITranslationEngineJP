import { describe, expect, it, vi } from "vitest"

import { createNoopFrontendDiagnosticLogger } from "./frontend-diagnostic-logger"

describe("createNoopFrontendDiagnosticLogger", () => {
  it("ignores diagnostic messages without writing to console", () => {
    const debug = vi.spyOn(console, "debug").mockImplementation(() => undefined)
    const info = vi.spyOn(console, "info").mockImplementation(() => undefined)
    const warn = vi.spyOn(console, "warn").mockImplementation(() => undefined)
    const error = vi.spyOn(console, "error").mockImplementation(() => undefined)

    const logger = createNoopFrontendDiagnosticLogger().with({ jobId: "job-1" })
    logger.debug("debug message", { phase: "setup" })
    logger.info("info message")
    logger.warn("warn message")
    logger.error("error message")

    expect(debug).not.toHaveBeenCalled()
    expect(info).not.toHaveBeenCalled()
    expect(warn).not.toHaveBeenCalled()
    expect(error).not.toHaveBeenCalled()
  })
})
