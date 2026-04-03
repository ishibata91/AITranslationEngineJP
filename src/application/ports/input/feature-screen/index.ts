export type FeatureScreenState<TData, TSelection, TFilters> = {
  data: TData | null;
  error: string | null;
  filters: TFilters;
  loading: boolean;
  selection: TSelection | null;
};

export interface FeatureScreenStorePort<TData, TSelection, TFilters> {
  getState(): FeatureScreenState<TData, TSelection, TFilters>;
  setError(message: string): void;
  setFilters(filters: TFilters): void;
  setLoaded(payload: { data: TData; selection: TSelection | null }): void;
  setLoading(): void;
  setSelection(selection: TSelection | null): void;
}

export interface FeatureScreenUsecase<TSelection, TFilters> {
  initialize(): Promise<void>;
  refresh(): Promise<void>;
  retry(): Promise<void>;
  select(selection: TSelection | null): void;
  updateFilters(
    filters: TFilters,
    options?: { reload?: boolean },
  ): Promise<void>;
}
