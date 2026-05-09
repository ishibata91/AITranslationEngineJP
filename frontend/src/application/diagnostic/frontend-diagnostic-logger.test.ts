import { afterEach, describe, expect, it, vi } from "vitest"

import {
  createNoopFrontendDiagnosticLogger,
  createPinoFrontendDiagnosticLogger
} from "./frontend-diagnostic-logger"

afterEach(() => {
  vi.restoreAllMocks()
})

describe("createPinoFrontendDiagnosticLogger", () => {
  it("writes sanitized diagnostic fields to browser console", () => {
    const info = vi.spyOn(console, "info").mockImplementation(() => undefined)

    const logger = createPinoFrontendDiagnosticLogger({
      source: "frontend",
      ignored: undefined
    }).with({ where: "frontend.runtime" })
    logger.info("diagnostic message", {
      event: "runtime_event_dropped",
      result: "skipped",
      id: "job:1",
      reason: "stale_event"
    })

    expect(info).toHaveBeenCalledOnce()
    expect(info.mock.calls[0][0]).toMatchObject({
      source: "frontend",
      where: "frontend.runtime",
      event: "runtime_event_dropped",
      result: "skipped",
      id: "job:1",
      reason: "stale_event",
      msg: "diagnostic message"
    })
    expect(info.mock.calls[0][0]).not.toHaveProperty("ignored")
  })
})

describe("createNoopFrontendDiagnosticLogger", () => {
  it("ignores diagnostic messages without writing to console", () => {
    const debug = vi.spyOn(console, "debug").mockImplementation(() => undefined)
    const info = vi.spyOn(console, "info").mockImplementation(() => undefined)
    const warn = vi.spyOn(console, "warn").mockImplementation(() => undefined)
    const error = vi.spyOn(console, "error").mockImplementation(() => undefined)

    const logger = createNoopFrontendDiagnosticLogger().with({ id: "job:1" })
    logger.debug("debug message", { where: "frontend.runtime" })
    logger.info("info message")
    logger.warn("warn message")
    logger.error("error message")

    expect(debug).not.toHaveBeenCalled()
    expect(info).not.toHaveBeenCalled()
    expect(warn).not.toHaveBeenCalled()
    expect(error).not.toHaveBeenCalled()
  })
})
