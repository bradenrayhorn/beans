<script lang="ts">
  import type { ComponentType } from "svelte";

  export let isActive: boolean;
  export let path: string;
  export let icon: ComponentType;
  export let name: string;
</script>

<a
  class="flex flex-row items-center gap-3 px-4 py-3 hover:bg-base-content/10"
  class:text-accent={isActive}
  href={path}
>
  <svelte:component this={icon} />
  <span class="text-xs font-semibold">{name}</span>
</a>
