import type {
  FeatureScreenState,
  FeatureScreenStorePort,
  FeatureScreenUsecase,
} from "@application/ports/input/feature-screen";

export type PersonaObserveEntry = {
  npcFormId: string;
  npcName: string;
  personaText: string;
  race: string;
  sex: string;
  voice: string;
};

type PersonaObserveRequest = {
  personaName: string;
};

export type PersonaObserveResult = {
  entries: PersonaObserveEntry[];
  personaName: string;
  sourceType: string;
};

export type PersonaObserveFilters = {
  lastSubmittedRequest: PersonaObserveRequest | null;
  personaName: string;
};

export type PersonaObserveScreenState = FeatureScreenState<
  PersonaObserveResult,
  number,
  PersonaObserveFilters
>;

type PersonaObserveSubscriber = (state: PersonaObserveScreenState) => void;

export interface PersonaObserveScreenStore extends FeatureScreenStorePort<
  PersonaObserveResult,
  number,
  PersonaObserveFilters
> {
  subscribe(run: PersonaObserveSubscriber): () => void;
}

export interface PersonaObserveScreenInput extends FeatureScreenUsecase<
  number,
  PersonaObserveFilters
> {
  initialize(): Promise<void>;
  observe(): Promise<void>;
  refresh(): Promise<void>;
  retry(): Promise<void>;
  select(selection: number | null): void;
}

type CreatePersonaObserveScreenUsecaseOptions = {
  executor: (request: PersonaObserveRequest) => Promise<PersonaObserveResult>;
  store: PersonaObserveScreenStore;
  toErrorMessage?: (error: unknown) => string;
};

function defaultToErrorMessage(): string {
  return "Persona observation failed. Try again.";
}

function createInitialState(): PersonaObserveScreenState {
  return {
    data: null,
    error: null,
    filters: {
      lastSubmittedRequest: null,
      personaName: "",
    },
    loading: false,
    selection: null,
  };
}

function cloneRequest(request: PersonaObserveRequest): PersonaObserveRequest {
  return {
    personaName: request.personaName,
  };
}

function cloneFilters(filters: PersonaObserveFilters): PersonaObserveFilters {
  return {
    lastSubmittedRequest:
      filters.lastSubmittedRequest === null
        ? null
        : cloneRequest(filters.lastSubmittedRequest),
    personaName: filters.personaName,
  };
}

function createObserveRequest(
  state: PersonaObserveScreenState,
): PersonaObserveRequest | null {
  if (state.filters.personaName.length === 0) {
    return null;
  }

  return {
    personaName: state.filters.personaName,
  };
}

function reconcileSelection(args: {
  currentSelection: number | null;
  data: PersonaObserveResult;
  mode: "observe" | "refresh";
}): number | null {
  if (args.data.entries.length === 0) {
    return null;
  }

  if (args.mode === "observe") {
    return 0;
  }

  if (
    args.currentSelection !== null &&
    args.currentSelection >= 0 &&
    args.currentSelection < args.data.entries.length
  ) {
    return args.currentSelection;
  }

  return 0;
}

export function createPersonaObserveScreenStore(): PersonaObserveScreenStore {
  let state = createInitialState();
  const subscribers = new Set<PersonaObserveSubscriber>();

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

export function createPersonaObserveScreenUsecase({
  executor,
  store,
  toErrorMessage = defaultToErrorMessage,
}: CreatePersonaObserveScreenUsecaseOptions): PersonaObserveScreenInput {
  async function loadRequest(
    request: PersonaObserveRequest,
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
