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

  $: category = data.categoryGroups
    .flatMap((g) => g.categories)
    .find((c) => c.id === $page.params.categoryID);

  $: monthCategory = data.month.categories.find(
    (mc) => mc.categoryID === $page.params.categoryID,
  );

  let isSubmitting = false;
  $: isLoading = !!$navigating || isSubmitting;
</script>

<div class="flex items-center justify-between">
  <h2 class="font-bold">{category?.name}</h2>

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
  <input name="category_id" type="text" class="hidden" value={category?.id} />

  <label>
    <span class="label label-text">Assigned</span>
    <input
      name="amount"
      type="text"
      class="input input-sm input-bordered w-full"
      value={monthCategory?.assigned.rawDisplay ?? ""}
    />
  </label>

  <div class="w-full flex flex-row justify-between">
    <SubmitButton class="btn btn-success btn-sm" {isLoading}>Save</SubmitButton>
  </div>
</form>
