<script>
  import "../app.postcss";
  import { initTheme } from "$lib/theme-switch";
</script>

<svelte:head>
  <!-- The \u{73} syntax is needed to avoid a parse error: https://github.com/sveltejs/svelte-eslint-parser/issues/346. -->
  <!-- eslint-disable-next-line svelte/no-at-html-tags -->
  {@html `<\u{73}cript>${initTheme.toString()} initTheme(); </script>`}
</svelte:head>

<slot />
