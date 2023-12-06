<script lang="ts">
  export let isLoading: boolean;
  let className = "";
  export { className as class };
</script>

<button class={`relative ${className}`} type="submit" disabled={isLoading}>
  {#if isLoading}
    <span class="loading loading-spinner absolute" />
    <span class="invisible"><slot /></span>
  {:else}
    <slot />
  {/if}
</button>
