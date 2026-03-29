export interface FeatureScreenGateway<TRequest, TData> {
  load(request: TRequest): Promise<TData>;
}
