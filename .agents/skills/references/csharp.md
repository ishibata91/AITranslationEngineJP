# C#

- project が有効にする nullable reference types と implicit usings を前提にする。
- nullable value は使用前に絞り込み、null-forgiving operator で検査を省略しない。
- `IDisposable` と `IAsyncDisposable` は `using` で確実に破棄する。
- 文字列の識別子、並び順、重複判定では用途に合う `StringComparer` を明示する。
- 非同期処理は `Task` を返し、同期的な待機で塞がない。
- 例外は処理と対象が分かる文脈を付けて上位へ返し、失敗を握りつぶさない。
- collection query は意味が明確な場合に LINQ を使い、副作用を混ぜない。
- xUnit test は外部から観測できる結果を検証し、共有状態と実行順へ依存させない。
