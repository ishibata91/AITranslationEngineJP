import type {
  FeatureScreenState,
  FeatureScreenStorePort,
  FeatureScreenUsecase,
} from "@application/ports/input/feature-screen";

export type TranslationPreviewReusableTerm = {
  destText: string;
  sourceText: string;
};

export type TranslationPreviewJobPersona = {
  npcFormId: string;
  personaText: string;
  race: string;
  sex: string;
  voice: string;
} | null;

export type TranslationPreviewEmbeddedElement = {
  elementId: string;
  rawText: string;
};

export type TranslationPreviewItem = {
  embeddedElementPolicy: {
    descriptors: TranslationPreviewEmbeddedElement[];
    unitKey: string;
  };
  jobId: string;
  jobPersona: TranslationPreviewJobPersona;
  reusableTerms: TranslationPreviewReusableTerm[];
  translatedText: string;
  translationUnit: {
    editorId: string;
    extractionKey: string;
    fieldName: string;
    formId: string;
    recordSignature: string;
    sortKey: string;
    sourceEntityType: string;
    sourceText: string;
  };
  unitKey: string;
};

type TranslationPreviewRequest = {
  jobId: string;
};

export type TranslationPreviewResult = {
  items: TranslationPreviewItem[];
  jobId: string;
};

export type TranslationPreviewFilters = {
  jobId: string;
  lastSubmittedRequest: TranslationPreviewRequest | null;
};

export type TranslationPreviewScreenState = FeatureScreenState<
  TranslationPreviewResult,
  string,
  TranslationPreviewFilters
>;

type TranslationPreviewSubscriber = (
  state: TranslationPreviewScreenState,
) => void;

export interface TranslationPreviewScreenStore extends FeatureScreenStorePort<
  TranslationPreviewResult,
  string,
  TranslationPreviewFilters
> {
  subscribe(run: TranslationPreviewSubscriber): () => void;
}

export interface TranslationPreviewScreenInput extends FeatureScreenUsecase<
  string,
  TranslationPreviewFilters
> {
  initialize(): Promise<void>;
  observe(): Promise<void>;
  refresh(): Promise<void>;
  retry(): Promise<void>;
  select(selection: string | null): void;
}

type CreateTranslationPreviewScreenUsecaseOptions = {
  executor: (
    request: TranslationPreviewRequest,
  ) => Promise<TranslationPreviewResult>;
  store: TranslationPreviewScreenStore;
  toErrorMessage?: (error: unknown) => string;
};

function defaultToErrorMessage(): string {
  return "Translation preview failed. Try again.";
}

function createInitialState(): TranslationPreviewScreenState {
  return {
    data: null,
    error: null,
    filters: {
      jobId: "",
      lastSubmittedRequest: null,
    },
    loading: false,
    selection: null,
  };
}

function cloneRequest(
  request: TranslationPreviewRequest,
): TranslationPreviewRequest {
  return {
    jobId: request.jobId,
  };
}

function cloneFilters(
  filters: TranslationPreviewFilters,
): TranslationPreviewFilters {
  return {
    jobId: filters.jobId,
    lastSubmittedRequest:
      filters.lastSubmittedRequest === null
        ? null
        : cloneRequest(filters.lastSubmittedRequest),
  };
}

function createObserveRequest(
  state: TranslationPreviewScreenState,
): TranslationPreviewRequest | null {
  if (state.filters.jobId.length === 0) {
    return null;
  }

  return {
    jobId: state.filters.jobId,
  };
}

function reconcileSelection(args: {
  currentSelection: string | null;
  data: TranslationPreviewResult;
  mode: "observe" | "refresh";
}): string | null {
  if (args.data.items.length === 0) {
    return null;
  }

  if (args.mode === "observe") {
    return args.data.items[0].unitKey;
  }

  if (
    args.currentSelection !== null &&
    args.data.items.some((item) => item.unitKey === args.currentSelection)
  ) {
    return args.currentSelection;
  }

  return args.data.items[0].unitKey;
}

export function createTranslationPreviewScreenStore(): TranslationPreviewScreenStore {
  let state = createInitialState();
  const subscribers = new Set<TranslationPreviewSubscriber>();

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

export function createTranslationPreviewScreenUsecase({
  executor,
  store,
  toErrorMessage = defaultToErrorMessage,
}: CreateTranslationPreviewScreenUsecaseOptions): TranslationPreviewScreenInput {
  async function loadRequest(
    request: TranslationPreviewRequest,
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
