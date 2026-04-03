import { describe, expect, it } from "vitest";
import { createFeatureScreenUsecase } from "./index";
import type { FeatureScreenStorePort } from "@application/ports/input/feature-screen";

type ExampleState = {
  items: string[];
};

type ExampleFilters = {
  query: string;
};

function createStore(
  initialData: ExampleState | null = null,
  initialSelection: string | null = null,
): FeatureScreenStorePort<ExampleState, string, ExampleFilters> {
  let state = {
    data: initialData,
    error: null as string | null,
    filters: { query: "bootstrap" },
    loading: false,
    selection: initialSelection,
  };

  return {
    getState() {
      return state;
    },
    setError(message) {
      state = {
        ...state,
        error: message,
        loading: false,
      };
    },
    setFilters(filters) {
      state = {
        ...state,
        filters,
      };
    },
    setLoaded(payload) {
      state = {
        ...state,
        data: payload.data,
        error: null,
        loading: false,
        selection: payload.selection,
      };
    },
    setLoading() {
      state = {
        ...state,
        error: null,
        loading: true,
      };
    },
    setSelection(selection) {
      state = {
        ...state,
        selection,
      };
    },
  };
}

describe("createFeatureScreenUsecase", () => {
  it("loads data through the gateway and preserves valid selection", async () => {
    const store = createStore(null, "job-1");
    const requests: ExampleFilters[] = [];
    const usecase = createFeatureScreenUsecase({
      createRequest(state) {
        return state.filters;
      },
      gateway: {
        async load(request) {
          requests.push(request);
          return { items: ["job-1", "job-2"] };
        },
      },
      reconcileSelection({ currentSelection, data }) {
        return currentSelection !== null &&
          data.items.includes(currentSelection)
          ? currentSelection
          : null;
      },
      store,
    });

    await usecase.initialize();

    expect(requests).toEqual([{ query: "bootstrap" }]);
    expect(store.getState()).toEqual({
      data: { items: ["job-1", "job-2"] },
      error: null,
      filters: { query: "bootstrap" },
      loading: false,
      selection: "job-1",
    });
  });

  it("keeps data and selection when a refresh fails", async () => {
    const store = createStore({ items: ["dictionary-a"] }, "dictionary-a");
    const usecase = createFeatureScreenUsecase({
      createRequest(state) {
        return state.filters;
      },
      gateway: {
        async load() {
          throw new Error("gateway failed");
        },
      },
      store,
    });

    await usecase.refresh();

    expect(store.getState()).toEqual({
      data: { items: ["dictionary-a"] },
      error: "gateway failed",
      filters: { query: "bootstrap" },
      loading: false,
      selection: "dictionary-a",
    });
  });
});
