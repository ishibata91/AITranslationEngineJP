interface MasterPersonaRuntimePollingAdapterOptions {
  intervalMs?: number
  scheduler?: (callback: () => void, delay: number) => number
  canceller?: (handle: number) => void
}

export class MasterPersonaRuntimePollingAdapter {
  private timerHandle: number | null = null

  private readonly intervalMs: number

  private readonly scheduler: (callback: () => void, delay: number) => number

  private readonly canceller: (handle: number) => void

  constructor(options: MasterPersonaRuntimePollingAdapterOptions = {}) {
    this.intervalMs = options.intervalMs ?? 1500
    this.scheduler =
      options.scheduler ??
      ((callback, delay) => window.setInterval(callback, delay))
    this.canceller =
      options.canceller ?? ((handle) => window.clearInterval(handle))
  }

  start(onTick: () => void): boolean {
    this.stop()
    if (typeof window === "undefined") {
      return false
    }
    this.timerHandle = this.scheduler(onTick, this.intervalMs)
    return true
  }

  stop(): void {
    if (this.timerHandle === null) {
      return
    }
    this.canceller(this.timerHandle)
    this.timerHandle = null
  }
}
