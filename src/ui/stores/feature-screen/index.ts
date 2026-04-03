import type { Readable } from "svelte/store";
import { get, writable } from "svelte/store";
import type {
  FeatureScreenState,
  FeatureScreenStorePort,
} from "@application/ports/input/feature-screen";

export interface FeatureScreenStore<
  TData,
  TSelection,
  TFilters,
> extends FeatureScreenStorePort<TData, TSelection, TFilters> {
  subscribe: Readable<
    FeatureScreenState<TData, TSelection, TFilters>
  >["subscribe"];
}

type CreateFeatureScreenStoreOptions<TFilters> = {
  filters: TFilters;
};

export function createFeatureScreenStore<TData, TSelection, TFilters>({
  filters,
}: CreateFeatureScreenStoreOptions<TFilters>): FeatureScreenStore<
  TData,
  TSelection,
  TFilters
> {
  const screenState = writable<FeatureScreenState<TData, TSelection, TFilters>>(
    {
      data: null,
      error: null,
      filters,
      loading: false,
      selection: null,
    },
  );

  return {
    subscribe: screenState.subscribe,
    getState() {
      return get(screenState);
    },
    setError(message) {
      screenState.update((current) => ({
        ...current,
        error: message,
        loading: false,
      }));
    },
    setFilters(nextFilters) {
      screenState.update((current) => ({
        ...current,
        filters: nextFilters,
      }));
    },
    setLoaded(payload) {
      screenState.update((current) => ({
        ...current,
        data: payload.data,
        error: null,
        loading: false,
        selection: payload.selection,
      }));
    },
    setLoading() {
      screenState.update((current) => ({
        ...current,
        error: null,
        loading: true,
      }));
    },
    setSelection(selection) {
      screenState.update((current) => ({
        ...current,
        selection,
      }));
    },
  };
}
