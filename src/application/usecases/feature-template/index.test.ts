import { describe, expect, it } from "vitest";
import { createFeatureTemplateScreenUsecase } from "./index";
import type { FeatureScreenStorePort } from "@application/ports/input/feature-screen";
import type {
  FeatureTemplateData,
  FeatureTemplateQuery,
} from "@shared/contracts/feature-template";

function createStore(
  initialFilters: FeatureTemplateQuery,
): FeatureScreenStorePort<FeatureTemplateData, string, FeatureTemplateQuery> {
  let state = {
    data: null as FeatureTemplateData | null,
    error: null as string | null,
    filters: initialFilters,
    loading: false,
    selection: null as string | null,
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

describe("createFeatureTemplateScreenUsecase", () => {
  it("Given filters and valid selection When initialize runs Then request and selection are preserved", async () => {
    const store = createStore({
      query: "job",
    });

    store.setSelection("job-2");

    const requests: string[] = [];
    const usecase = createFeatureTemplateScreenUsecase({
      gateway: {
        async load(request) {
          requests.push(request.query);

          return {
            items: [
              {
                detail: "queued",
                id: "job-1",
                status: "queued",
                title: "Import job",
              },
              {
                detail: "running",
                id: "job-2",
                status: "running",
                title: "Translate job",
              },
            ],
          };
        },
      },
      store,
    });

    await usecase.initialize();

    expect(requests).toEqual(["job"]);
    expect(store.getState()).toEqual({
      data: {
        items: [
          {
            detail: "queued",
            id: "job-1",
            status: "queued",
            title: "Import job",
          },
          {
            detail: "running",
            id: "job-2",
            status: "running",
            title: "Translate job",
          },
        ],
      },
      error: null,
      filters: {
        query: "job",
      },
      loading: false,
      selection: "job-2",
    });
  });

  it("Given loaded data When refresh fails Then error updates without clearing data or selection", async () => {
    const store = createStore({
      query: "persona",
    });

    store.setLoaded({
      data: {
        items: [
          {
            detail: "ready",
            id: "persona-a",
            status: "ready",
            title: "Persona observation",
          },
        ],
      },
      selection: "persona-a",
    });

    const usecase = createFeatureTemplateScreenUsecase({
      gateway: {
        async load() {
          throw new Error("template gateway failed");
        },
      },
      store,
    });

    await usecase.refresh();

    expect(store.getState()).toEqual({
      data: {
        items: [
          {
            detail: "ready",
            id: "persona-a",
            status: "ready",
            title: "Persona observation",
          },
        ],
      },
      error: "template gateway failed",
      filters: {
        query: "persona",
      },
      loading: false,
      selection: "persona-a",
    });
  });
});
