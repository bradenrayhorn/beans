<script>
  import "../app.postcss";
  import { initTheme } from "$lib/theme-switch";
  import { navigating } from "$app/stores";
</script>

<svelte:head>
  <!-- The \u{73} syntax is needed to avoid a parse error: https://github.com/sveltejs/svelte-eslint-parser/issues/346. -->
  <!-- eslint-disable-next-line svelte/no-at-html-tags -->
  {@html `<\u{73}cript>${initTheme.toString()} initTheme(); </script>`}
</svelte:head>

{#if $navigating}
  <div class="fixed top-0 left-0 right-0 h-1.5 bg-primary/50 z-50">
    <div
      class="bg-primary animate-progressFast md:animate-progress w-32 lg:w-64 h-full will-change-transform"
    />
  </div>
{/if}

<slot />
