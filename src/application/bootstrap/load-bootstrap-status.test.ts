import { describe, expect, it } from "vitest";
import { loadBootstrapStatus } from "./load-bootstrap-status";

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
});
