<script lang="ts">
  import IconLight from "~icons/mdi/weather-sunny";
  import IconDark from "~icons/mdi/moon-and-stars";
  import { onMount } from "svelte";
  import { setTheme } from "$lib/theme-switch";

  let isDark = false;

  onMount(() => {
    isDark =
      document.documentElement.getAttribute("data-theme") === "beansDark";
  });
</script>

<svelte:head>
  <script lang="ts">
    isDark = false;
  </script>
</svelte:head>

<button
  class="btn btn-sm btn-circle btn-ghost"
  on:click={() => {
    setTheme(!isDark);
    isDark = !isDark;
  }}
>
  {#if isDark}
    <IconDark />
  {:else}
    <IconLight />
  {/if}
</button>
