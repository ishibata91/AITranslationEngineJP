<script lang="ts">
  // 翻訳対象プラグイン画面。汎用の表示部品を組み立てるだけにし、state・API・Wails・validation は持たない。
  // すべて props で受け、操作 event は callback prop で上へ返す。
  // 新しいプラグインを選ぶ入口と、翻訳したプラグインの一覧（進捗・削除・結果を開く導線）を出す。
  // AI 設定・実行・進捗・結果一覧は、選んだプラグインの実行・結果画面（別画面）が持つ。
  import StatusBadge from "@ui/components/StatusBadge.svelte"
  import FileSelectField from "@ui/components/FileSelectField.svelte"
  import { pluginProgressBadge } from "./target-plugins-presentation"
  import type { TargetPluginRow } from "./target-plugins-view"

  interface Props {
    // 新規翻訳の入口。選んだプラグインのパス。
    newPluginPath: string
    // プラグイン選択ダイアログの起動。
    onSelectNewPlugin: () => void
    // プラグインパスの直接入力。
    onNewPluginPathInput: (value: string) => void
    // 選んだプラグインの実行・結果画面へ進む。
    onProceedToRun: () => void
    // 表示するプラグイン一覧。作成時刻の新しい順で並べて渡す想定。
    plugins: TargetPluginRow[]
    // 一覧取得中フラグ。true の間はローディングを出す。
    loading?: boolean
    // 削除確認中のプラグイン。指定した行だけ確認 UI を出す。未指定なら通常の削除ボタン。
    pendingDeletePlugin?: string
    // 削除実行中フラグ。pendingDeletePlugin の行で確定ボタンをスピナー・無効化する。
    deleting?: boolean
    // エラー文言。空なら出さない。
    errorMessage?: string
    // 一覧の行を選ぶ。そのプラグインの実行・結果画面へ進む。
    onOpenPlugin: (plugin: string) => void
    // 行の削除ボタン押下。確認 UI を出すのは container が担う。
    onRequestDelete: (plugin: string) => void
    // 削除確認の確定。実削除は callback で上へ返す。
    onConfirmDelete: (plugin: string) => void
    // 削除確認の取消。
    onCancelDelete: () => void
  }

  let {
    newPluginPath,
    onSelectNewPlugin,
    onNewPluginPathInput,
    onProceedToRun,
    plugins,
    loading = false,
    pendingDeletePlugin,
    deleting = false,
    errorMessage = "",
    onOpenPlugin,
    onRequestDelete,
    onConfirmDelete,
    onCancelDelete
  }: Props = $props()
</script>

<div class="min-h-screen w-full px-6 py-12 flex justify-center">
  <section class="w-full max-w-4xl flex flex-col gap-8" aria-labelledby="targets-title">
    <header class="flex flex-col gap-3">
      <span class="u-mono text-xs tracking-[0.32em] text-accent">
        翻訳の管理
      </span>
      <h1 id="targets-title" class="u-display text-4xl font-semibold text-base-content">
        翻訳対象プラグイン
      </h1>
      <p class="max-w-2xl text-base-content/70 leading-relaxed">
        翻訳するプラグインを選び、翻訳したプラグインを一覧します。削除するとそのプラグインの翻訳成果を消し、再実行で最初からやり直します。ほかのプラグインの成果と共有辞書は残ります。
      </p>
      <div
        class="h-px w-full bg-gradient-to-r from-transparent via-primary/50 to-transparent"
      ></div>
    </header>

    <div class="card bg-base-200/55 border border-base-300/70 shadow-xl u-edge-top">
      <div class="card-body gap-4">
        <h2 class="u-display text-sm tracking-widest text-base-content/60">
          新しいプラグインを翻訳
        </h2>
        <FileSelectField
          label="プラグイン"
          buttonLabel="プラグインを選択"
          path={newPluginPath}
          placeholder="/path/to/Data/Mod.esp"
          emptyText="プラグインファイルが未選択です。"
          hint="master と Strings がある Data フォルダ内のプラグインを選ぶか、フルパスを直接入力します。"
          onSelect={onSelectNewPlugin}
          onPathInput={onNewPluginPathInput}
        />
        <div class="card-actions justify-end pt-1">
          <button
            type="button"
            class="btn btn-primary px-8"
            disabled={newPluginPath.length === 0}
            onclick={onProceedToRun}
          >
            翻訳へ進む
          </button>
        </div>
      </div>
    </div>

    <div class="card bg-base-200/40 border border-base-300/60">
      <div class="card-body gap-5">
        <div class="flex items-center justify-between gap-3">
          <h2 class="u-display text-sm tracking-widest text-base-content/60">
            プラグイン一覧
          </h2>
          {#if plugins.length > 0}
            <span class="badge badge-outline badge-sm u-mono">{plugins.length} 件</span>
          {/if}
        </div>

        {#if errorMessage.length > 0}
          <div role="alert" class="alert alert-error alert-soft">
            <span>{errorMessage}</span>
          </div>
        {/if}

        {#if loading}
          <div class="flex flex-col items-center gap-3 py-12 text-center text-base-content/50">
            <span class="loading loading-ring loading-lg text-primary"></span>
            <p>プラグインを読み込んでいます。</p>
          </div>
        {:else if plugins.length === 0}
          <div class="flex flex-col items-center gap-3 py-12 text-center text-base-content/50">
            <span class="u-display text-3xl text-base-content/25">ᚱ</span>
            <p>まだ翻訳したプラグインはありません。上でプラグインを選んで翻訳すると、ここに残ります。</p>
          </div>
        {:else}
          <ul class="flex flex-col gap-2">
            {#each plugins as row (row.plugin)}
              {@const progress = pluginProgressBadge(row)}
              <li
                class="rounded-lg border border-base-300/60 bg-base-100/40 flex flex-wrap items-center justify-between gap-x-4 gap-y-3"
              >
                <button
                  type="button"
                  class="flex grow flex-col gap-1 min-w-0 px-4 py-3 text-left transition-colors hover:text-primary"
                  onclick={() => onOpenPlugin(row.plugin)}
                >
                  <span class="u-mono text-sm truncate">{row.plugin}</span>
                  <span class="text-xs text-base-content/50">作成 {row.createdLabel}・結果を開く</span>
                </button>
                <div class="flex items-center gap-3 px-4 py-3">
                  {#if pendingDeletePlugin === row.plugin}
                    <span class="text-xs text-base-content/70">
                      成果を消して再実行でやり直します。削除しますか？
                    </span>
                    <button
                      type="button"
                      class="btn btn-sm btn-error"
                      disabled={deleting}
                      onclick={() => onConfirmDelete(row.plugin)}
                    >
                      {#if deleting}
                        <span class="loading loading-spinner loading-xs"></span>
                      {/if}
                      {deleting ? "削除中…" : "削除する"}
                    </button>
                    <button
                      type="button"
                      class="btn btn-sm btn-ghost"
                      disabled={deleting}
                      onclick={onCancelDelete}
                    >
                      取消
                    </button>
                  {:else}
                    <StatusBadge label={progress.label} tone={progress.tone} />
                    <button
                      type="button"
                      class="btn btn-sm btn-outline btn-error"
                      onclick={() => onRequestDelete(row.plugin)}
                    >
                      削除
                    </button>
                  {/if}
                </div>
              </li>
            {/each}
          </ul>
        {/if}
      </div>
    </div>
  </section>
</div>
