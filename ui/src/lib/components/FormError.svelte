<script lang="ts">
  import IconError from "~icons/mdi/error";
  import { page } from "$app/stores";
</script>

{#if $page.form?.message}
  <div role="alert" class="alert alert-error my-8">
    <IconError />
    <span>{$page.form?.message}</span>
  </div>
{/if}
