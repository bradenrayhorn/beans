<script lang="ts">
  let className = "";
  export { className as class };
</script>

<div
  class={`font-bold text-lg text-center md:text-left md:text-2xl ${className}`}
>
  <slot />
</div>
