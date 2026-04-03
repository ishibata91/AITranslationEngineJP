import { describe, expect, it } from "vitest";
import { createFeatureScreenStore } from "./index";

describe("createFeatureScreenStore", () => {
  it("tracks loading, data, error, and selection in one screen state", () => {
    const store = createFeatureScreenStore<
      { items: string[] },
      string,
      { query: string }
    >({
      filters: { query: "" },
    });

    store.setLoading();
    store.setLoaded({
      data: { items: ["job-1", "job-2"] },
      selection: "job-2",
    });
    store.setError("refresh failed");

    expect(store.getState()).toEqual({
      data: { items: ["job-1", "job-2"] },
      error: "refresh failed",
      filters: { query: "" },
      loading: false,
      selection: "job-2",
    });
  });

  it("updates filters and selection without touching loaded data", () => {
    const store = createFeatureScreenStore<
      { items: string[] },
      string,
      { query: string }
    >({
      filters: { query: "old" },
    });

    store.setLoaded({
      data: { items: ["persona-a"] },
      selection: null,
    });
    store.setFilters({ query: "new" });
    store.setSelection("persona-a");

    expect(store.getState()).toEqual({
      data: { items: ["persona-a"] },
      error: null,
      filters: { query: "new" },
      loading: false,
      selection: "persona-a",
    });
  });
});
