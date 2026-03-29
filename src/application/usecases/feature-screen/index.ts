import type {
  FeatureScreenState,
  FeatureScreenStorePort,
  FeatureScreenUsecase
} from "@application/ports/input/feature-screen";
import type { FeatureScreenGateway } from "@application/ports/gateway/feature-screen";

type FeatureScreenSelectionResolver<TData, TSelection> = (args: {
  currentSelection: TSelection | null;
  data: TData;
  previousData: TData | null;
}) => TSelection | null;

type CreateFeatureScreenUsecaseOptions<TRequest, TData, TSelection, TFilters> = {
  createRequest: (state: FeatureScreenState<TData, TSelection, TFilters>) => TRequest;
  gateway: FeatureScreenGateway<TRequest, TData>;
  reconcileSelection?: FeatureScreenSelectionResolver<TData, TSelection>;
  store: FeatureScreenStorePort<TData, TSelection, TFilters>;
  toErrorMessage?: (error: unknown) => string;
};

function defaultToErrorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "Unknown screen failure.";
}

export function createFeatureScreenUsecase<TRequest, TData, TSelection, TFilters>({
  createRequest,
  gateway,
  reconcileSelection,
  store,
  toErrorMessage = defaultToErrorMessage
}: CreateFeatureScreenUsecaseOptions<TRequest, TData, TSelection, TFilters>): FeatureScreenUsecase<
  TSelection,
  TFilters
> {
  async function loadCurrent(): Promise<void> {
    const previousState = store.getState();

    store.setLoading();

    try {
      const nextState = store.getState();
      const data = await gateway.load(createRequest(nextState));
      const selection =
        reconcileSelection?.({
          currentSelection: nextState.selection,
          data,
          previousData: previousState.data
        }) ?? nextState.selection;

      store.setLoaded({
        data,
        selection
      });
    } catch (error) {
      store.setError(toErrorMessage(error));
    }
  }

  return {
    initialize() {
      return loadCurrent();
    },
    refresh() {
      return loadCurrent();
    },
    retry() {
      return loadCurrent();
    },
    select(selection) {
      store.setSelection(selection);
    },
    async updateFilters(filters, options) {
      store.setFilters(filters);

      if (options?.reload === true) {
        await loadCurrent();
      }
    }
  };
}
