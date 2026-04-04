import type {
  FeatureScreenState,
  FeatureScreenStorePort,
  FeatureScreenUsecase,
} from "@application/ports/input/feature-screen";

export type DictionaryCandidate = {
  destText: string;
  sourceText: string;
};

export type DictionaryCandidateGroup = {
  candidates: DictionaryCandidate[];
  sourceText: string;
};

type DictionaryObserveRequest = {
  sourceTexts: string[];
};

export type DictionaryObserveResult = {
  candidateGroups: DictionaryCandidateGroup[];
};

export type DictionaryObserveFilters = {
  lastSubmittedRequest: DictionaryObserveRequest | null;
  sourceTexts: string[];
};

export type DictionaryObserveScreenState = FeatureScreenState<
  DictionaryObserveResult,
  number,
  DictionaryObserveFilters
>;

type DictionaryObserveSubscriber = (
  state: DictionaryObserveScreenState,
) => void;

export interface DictionaryObserveScreenStore extends FeatureScreenStorePort<
  DictionaryObserveResult,
  number,
  DictionaryObserveFilters
> {
  subscribe(run: DictionaryObserveSubscriber): () => void;
}

export interface DictionaryObserveScreenInput extends FeatureScreenUsecase<
  number,
  DictionaryObserveFilters
> {
  initialize(): Promise<void>;
  observe(): Promise<void>;
  refresh(): Promise<void>;
  retry(): Promise<void>;
  select(selection: number | null): void;
}

type CreateDictionaryObserveScreenUsecaseOptions = {
  executor: (
    request: DictionaryObserveRequest,
  ) => Promise<DictionaryObserveResult>;
  store: DictionaryObserveScreenStore;
  toErrorMessage?: (error: unknown) => string;
};

function defaultToErrorMessage(): string {
  return "Dictionary observation failed. Try again.";
}

function createInitialState(): DictionaryObserveScreenState {
  return {
    data: null,
    error: null,
    filters: {
      lastSubmittedRequest: null,
      sourceTexts: [],
    },
    loading: false,
    selection: null,
  };
}

function cloneRequest(
  request: DictionaryObserveRequest,
): DictionaryObserveRequest {
  return {
    sourceTexts: [...request.sourceTexts],
  };
}

function cloneFilters(
  filters: DictionaryObserveFilters,
): DictionaryObserveFilters {
  return {
    lastSubmittedRequest:
      filters.lastSubmittedRequest === null
        ? null
        : cloneRequest(filters.lastSubmittedRequest),
    sourceTexts: [...filters.sourceTexts],
  };
}

function createObserveRequest(
  state: DictionaryObserveScreenState,
): DictionaryObserveRequest | null {
  if (state.filters.sourceTexts.length === 0) {
    return null;
  }

  return {
    sourceTexts: [...state.filters.sourceTexts],
  };
}

function reconcileSelection(args: {
  currentSelection: number | null;
  data: DictionaryObserveResult;
  mode: "observe" | "refresh";
}): number | null {
  if (args.data.candidateGroups.length === 0) {
    return null;
  }

  if (args.mode === "observe") {
    return 0;
  }

  if (
    args.currentSelection !== null &&
    args.currentSelection >= 0 &&
    args.currentSelection < args.data.candidateGroups.length
  ) {
    return args.currentSelection;
  }

  return 0;
}

export function createDictionaryObserveScreenStore(): DictionaryObserveScreenStore {
  let state = createInitialState();
  const subscribers = new Set<DictionaryObserveSubscriber>();

  function notify(): void {
    subscribers.forEach((subscriber) => subscriber(state));
  }

  return {
    subscribe(run) {
      subscribers.add(run);
      run(state);

      return () => {
        subscribers.delete(run);
      };
    },
    getState() {
      return state;
    },
    setError(message) {
      state = {
        ...state,
        error: message,
        loading: false,
      };
      notify();
    },
    setFilters(filters) {
      state = {
        ...state,
        filters: cloneFilters(filters),
      };
      notify();
    },
    setLoaded(payload) {
      state = {
        ...state,
        data: payload.data,
        error: null,
        loading: false,
        selection: payload.selection,
      };
      notify();
    },
    setLoading() {
      state = {
        ...state,
        error: null,
        loading: true,
      };
      notify();
    },
    setSelection(selection) {
      state = {
        ...state,
        selection,
      };
      notify();
    },
  };
}

export function createDictionaryObserveScreenUsecase({
  executor,
  store,
  toErrorMessage = defaultToErrorMessage,
}: CreateDictionaryObserveScreenUsecaseOptions): DictionaryObserveScreenInput {
  async function loadRequest(
    request: DictionaryObserveRequest,
    mode: "observe" | "refresh",
  ): Promise<void> {
    store.setLoading();

    try {
      const data = await executor(request);
      const currentState = store.getState();
      const selection = reconcileSelection({
        currentSelection: currentState.selection,
        data,
        mode,
      });

      store.setLoaded({
        data,
        selection,
      });
    } catch (error) {
      store.setError(toErrorMessage(error));
    }
  }

  async function rerunLastSubmittedRequest(): Promise<void> {
    const state = store.getState();

    if (state.filters.lastSubmittedRequest === null) {
      return;
    }

    await loadRequest(
      cloneRequest(state.filters.lastSubmittedRequest),
      "refresh",
    );
  }

  return {
    async initialize() {},
    async observe() {
      const state = store.getState();
      const request = createObserveRequest(state);

      if (request === null) {
        return;
      }

      store.setFilters({
        ...state.filters,
        lastSubmittedRequest: cloneRequest(request),
      });

      await loadRequest(request, "observe");
    },
    refresh() {
      return rerunLastSubmittedRequest();
    },
    retry() {
      return rerunLastSubmittedRequest();
    },
    select(selection) {
      store.setSelection(selection);
    },
    async updateFilters(filters) {
      store.setFilters(filters);
    },
  };
}
