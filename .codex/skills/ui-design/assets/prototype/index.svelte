<script>
  const screens = [
    {
      id: "primary",
      label: "主要画面",
      description: "最初に確認する画面",
    },
    {
      id: "secondary",
      label: "関連画面",
      description: "画面間導線で到達する画面",
    },
  ];

  let activeScreenId = screens[0].id;

  $: activeScreen = screens.find((screen) => screen.id === activeScreenId) ?? screens[0];
</script>

<script context="module">
  export const prototypeContract = {
    routeNavigation: "in_prototype_only",
    productionReference: "product_code_must_not_reference_ui_prototype",
  };
</script>

<section class="prototype-shell" data-ui-prototype-sample-data-root>
  <header class="prototype-header">
    <div>
      <p class="eyebrow">UIプロトタイプ</p>
      <h1>確認対象画面</h1>
    </div>
    <nav class="screen-nav" aria-label="確認対象画面">
      {#each screens as screen}
        <button
          class:active={activeScreenId === screen.id}
          type="button"
          on:click={() => (activeScreenId = screen.id)}
        >
          {screen.label}
        </button>
      {/each}
    </nav>
  </header>

  <main class="screen-frame" aria-live="polite">
    <section class="screen-summary">
      <p class="eyebrow">選択中</p>
      <h2>{activeScreen.label}</h2>
      <p>{activeScreen.description}</p>
    </section>

    {#if activeScreen.id === "primary"}
      <section class="placeholder-panel">
        <h3>主要画面の確認区画</h3>
        <p>表示項目、主要操作、状態差分を置く。</p>
      </section>
    {:else}
      <section class="placeholder-panel">
        <h3>関連画面の確認区画</h3>
        <p>画面間導線で変わる表示と操作を置く。</p>
      </section>
    {/if}
  </main>
</section>

<style>
  :global(body) {
    color: #2f2926;
    font-family:
      "Noto Sans JP",
      system-ui,
      sans-serif;
  }

  .prototype-shell {
    min-height: 100vh;
    padding: 32px;
  }

  .prototype-header {
    display: flex;
    align-items: end;
    justify-content: space-between;
    gap: 24px;
    max-width: 1120px;
    margin: 0 auto 24px;
  }

  h1,
  h2,
  h3,
  p {
    margin: 0;
  }

  h1 {
    font-size: 32px;
    font-weight: 700;
  }

  h2 {
    font-size: 24px;
    font-weight: 700;
  }

  h3 {
    font-size: 18px;
    font-weight: 700;
  }

  .eyebrow {
    margin-bottom: 6px;
    font-size: 13px;
    font-weight: 700;
  }

  .screen-nav {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .screen-nav button {
    min-height: 40px;
    padding: 0 14px;
    border: 1px solid #a79b92;
    border-radius: 6px;
    color: #2f2926;
    background: #fffaf5;
    font: inherit;
    cursor: pointer;
  }

  .screen-nav button.active {
    color: #ffffff;
    background: #3f5f58;
    border-color: #3f5f58;
  }

  .screen-frame {
    display: grid;
    grid-template-columns: 280px minmax(0, 1fr);
    gap: 16px;
    max-width: 1120px;
    margin: 0 auto;
  }

  .screen-summary,
  .placeholder-panel {
    min-height: 220px;
    padding: 20px;
    border: 1px solid #d4c8be;
    border-radius: 8px;
    background: #fffaf5;
  }

  .placeholder-panel {
    display: grid;
    align-content: start;
    gap: 12px;
  }

  @media (max-width: 720px) {
    .prototype-shell {
      padding: 20px;
    }

    .prototype-header,
    .screen-frame {
      grid-template-columns: 1fr;
    }

    .prototype-header {
      align-items: stretch;
      flex-direction: column;
    }
  }
</style>
