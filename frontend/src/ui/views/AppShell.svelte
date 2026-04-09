<script lang="ts">
  interface Props {
    title: string
  }

  let { title }: Props = $props()
</script>

<main class="shell">
  <header class="hero">
    <p class="eyebrow">AITranslationEngineJp</p>
    <h1>{title}</h1>
    <p class="description">
      `docs/architecture.md` に沿った Wails + Go + Svelte の初期骨格です。
    </p>
  </header>
</main>

<style>
  :global(body) {
    margin: 0;
    font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
    background: radial-gradient(
      circle at top,
      #f3efe1 0,
      #efe7d2 30%,
      #d8d4c4 100%
    );
    color: #18222e;
  }

  .shell {
    min-height: 100vh;
    display: grid;
    place-items: center;
    padding: 32px;
  }

  .hero {
    width: min(720px, 100%);
    padding: 40px;
    border: 1px solid rgba(24, 34, 46, 0.12);
    background: rgba(255, 252, 246, 0.8);
    backdrop-filter: blur(12px);
    box-shadow: 0 24px 80px rgba(24, 34, 46, 0.12);
  }

  .eyebrow {
    margin: 0 0 12px;
    font-size: 0.8rem;
    letter-spacing: 0.2em;
    text-transform: uppercase;
    color: #875b2a;
  }

  h1 {
    margin: 0;
    font-size: clamp(2.5rem, 6vw, 4.5rem);
    line-height: 0.95;
  }

  .description {
    margin: 20px 0 0;
    font-size: 1rem;
    line-height: 1.7;
  }
</style>
