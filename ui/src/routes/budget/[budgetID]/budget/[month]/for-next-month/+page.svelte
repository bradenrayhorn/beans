<script lang="ts">
  import { page } from "$app/stores";
  import type { PageData } from "./$types";
  import { enhance } from "$app/forms";
  import { navigating } from "$app/stores";
  import FormError from "$lib/components/FormError.svelte";
  import SubmitButton from "$lib/components/SubmitButton.svelte";
  import IconClose from "~icons/mdi/Close";
  import { paths, withParameter } from "$lib/paths";

  export let data: PageData;

  let isSubmitting = false;
  $: isLoading = !!$navigating || isSubmitting;
</script>

<div class="flex items-center justify-between">
  <h2 class="font-bold">For Next Month</h2>

  <a
    href={withParameter(paths.budget.budget.month, $page.params)}
    class="btn btn-ghost btn-sm me-1"
    aria-label="Close form"
  >
    <IconClose />
  </a>
</div>

<FormError />

<form
  class="flex flex-col gap-2"
  method="POST"
  action="?/save"
  use:enhance={() => {
    isSubmitting = true;

    return async ({ update }) => {
      await update({ reset: false });

      isSubmitting = false;
    };
  }}
>
  <input name="monthID" type="text" class="hidden" value={data.month.id} />

  <label>
    <span class="label label-text">Carryover</span>
    <input
      name="carryover"
      type="text"
      class="input input-sm input-bordered w-full"
      value={data.month.carryover.rawDisplay ?? ""}
    />
  </label>

  <div class="w-full flex flex-row justify-between">
    <SubmitButton class="btn btn-primary btn-sm" {isLoading}>Save</SubmitButton>
  </div>
</form>
