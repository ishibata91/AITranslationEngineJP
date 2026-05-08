import pino from "pino"

type PinoLogger = ReturnType<typeof pino>

type DiagnosticValue = string | number | boolean | null

type DiagnosticFields = Record<string, DiagnosticValue | undefined>

interface FrontendDiagnosticLogger {
  with(fields: DiagnosticFields): FrontendDiagnosticLogger
  debug(message: string, fields?: DiagnosticFields): void
  info(message: string, fields?: DiagnosticFields): void
  warn(message: string, fields?: DiagnosticFields): void
  error(message: string, fields?: DiagnosticFields): void
}

export function createPinoFrontendDiagnosticLogger(
  fields: DiagnosticFields = {}
): FrontendDiagnosticLogger {
  return new PinoFrontendDiagnosticLogger(
    pino({
      browser: {
        asObject: true
      }
    }),
    sanitizeDiagnosticFields(fields)
  )
}

export function createNoopFrontendDiagnosticLogger(): FrontendDiagnosticLogger {
  return new NoopFrontendDiagnosticLogger()
}

class PinoFrontendDiagnosticLogger implements FrontendDiagnosticLogger {
  constructor(
    private readonly logger: PinoLogger,
    private readonly fields: DiagnosticFields
  ) {}

  with(fields: DiagnosticFields): FrontendDiagnosticLogger {
    return new PinoFrontendDiagnosticLogger(this.logger, {
      ...this.fields,
      ...sanitizeDiagnosticFields(fields)
    })
  }

  debug(message: string, fields: DiagnosticFields = {}): void {
    this.logger.debug(this.buildPayload(fields), message)
  }

  info(message: string, fields: DiagnosticFields = {}): void {
    this.logger.info(this.buildPayload(fields), message)
  }

  warn(message: string, fields: DiagnosticFields = {}): void {
    this.logger.warn(this.buildPayload(fields), message)
  }

  error(message: string, fields: DiagnosticFields = {}): void {
    this.logger.error(this.buildPayload(fields), message)
  }

  private buildPayload(fields: DiagnosticFields): DiagnosticFields {
    return {
      ...this.fields,
      ...sanitizeDiagnosticFields(fields)
    }
  }
}

class NoopFrontendDiagnosticLogger implements FrontendDiagnosticLogger {
  with(): FrontendDiagnosticLogger {
    return this
  }

  debug(): void {
    return ignoreDiagnosticMessage()
  }

  info(): void {
    return ignoreDiagnosticMessage()
  }

  warn(): void {
    return ignoreDiagnosticMessage()
  }

  error(): void {
    return ignoreDiagnosticMessage()
  }
}

function ignoreDiagnosticMessage(): void {
  return undefined
}

function sanitizeDiagnosticFields(fields: DiagnosticFields): DiagnosticFields {
  const sanitized: DiagnosticFields = {}
  for (const [key, value] of Object.entries(fields)) {
    if (value === undefined) {
      continue
    }
    sanitized[key] = value
  }
  return sanitized
}
