import { afterEach, describe, expect, it } from "vitest";
import { configureBootstrapStatusPort, loadBootstrapStatus } from "./load-bootstrap-status";

afterEach(() => {
  configureBootstrapStatusPort(null);
});

describe("loadBootstrapStatus", () => {
  it("returns the dto from the supplied gateway port", async () => {
    const result = await loadBootstrapStatus({
      async getBootstrapStatus() {
        return {
          backendVersion: "0.1.0",
          boundaryReady: true,
          frontendEntry: "src/main.ts"
        };
      }
    });

    expect(result).toEqual({
      backendVersion: "0.1.0",
      boundaryReady: true,
      frontendEntry: "src/main.ts"
    });
  });

  it("returns the dto from the configured default port", async () => {
    configureBootstrapStatusPort({
      async getBootstrapStatus() {
        return {
          backendVersion: "0.2.0",
          boundaryReady: true,
          frontendEntry: "src/main.ts"
        };
      }
    });

    await expect(loadBootstrapStatus()).resolves.toEqual({
      backendVersion: "0.2.0",
      boundaryReady: true,
      frontendEntry: "src/main.ts"
    });
  });

  it("throws when no bootstrap status port is configured", async () => {
    await expect(loadBootstrapStatus()).rejects.toThrow("BootstrapStatusPort is not configured.");
  });
});
