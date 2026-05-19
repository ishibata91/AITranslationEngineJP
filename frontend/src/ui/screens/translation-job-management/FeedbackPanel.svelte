<script lang="ts">
  import type { TranslationJobManagementFeedbackState } from "@application/contract/translation-job-management/translation-job-management-screen-types"

  interface Props {
    feedback: TranslationJobManagementFeedbackState
  }

  let { feedback }: Props = $props()

  function toneClass(tone: string): string {
    return `tone-${tone}`
  }
</script>

<section
  class={`panel feedback-panel ${toneClass(feedback.tone)}`}
  data-testid="translation-job-management-feedback-notification"
>
  <p class="feedback-title">{feedback.title}</p>
  <p class="feedback-message">{feedback.message}</p>
</section>

<style>
  .panel {
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    background: rgba(33, 27, 24, 0.88);
    padding: 1.1rem 1.2rem;
  }

  .feedback-title {
    margin: 0;
    color: #fff6ea;
  }

  .feedback-message {
    color: rgba(236, 223, 205, 0.78);
  }

  .feedback-panel,
  .tone-danger,
  .tone-warning {
    border-color: rgba(240, 180, 100, 0.3);
  }

  .tone-warning {
    border-color: rgba(255, 196, 92, 0.4);
  }

  .tone-danger {
    border-color: rgba(255, 128, 102, 0.38);
  }
</style>
