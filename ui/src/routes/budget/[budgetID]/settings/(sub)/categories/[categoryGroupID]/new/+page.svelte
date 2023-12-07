<script lang="ts">
  import { paths, withParameter } from "$lib/paths";
  import { enhance } from "$app/forms";
  import { navigating, page } from "$app/stores";
  import FormError from "$lib/components/FormError.svelte";
  import SubmitButton from "$lib/components/SubmitButton.svelte";
  import type { PageData } from "./$types";

  export let data: PageData;

  $: categoryGroup = data.categoryGroups.find(
    (it) => it.id === $page.params.categoryGroupID,
  );

  let isSubmitting = false;
  $: isLoading = !!$navigating || isSubmitting;
</script>

<div class="text-sm breadcrumbs mb-8">
  <ul>
    <li>
      <a href={withParameter(paths.budget.settings.base, $page.params)}>
        Settings
      </a>
    </li>
    <li>
      <a
        href={withParameter(
          paths.budget.settings.categories.base,
          $page.params,
        )}
      >
        Categories
      </a>
    </li>
    <li>
      <a
        href={withParameter(
          paths.budget.settings.categories.group,
          $page.params,
        )}
      >
        {categoryGroup?.name}
      </a>
    </li>
    <li>New Category</li>
  </ul>
</div>

<div class="max-w-md w-full mx-auto p-4">
  <FormError />

  <form
    class="flex flex-col gap-8"
    method="POST"
    action="?/save"
    use:enhance={() => {
      isSubmitting = true;

      return async ({ update }) => {
        await update();
        isSubmitting = false;
      };
    }}
  >
    <input
      name="group_id"
      type="text"
      class="hidden"
      value={categoryGroup?.id}
    />

    <label>
      <span class="label label-text">Name</span>
      <input name="name" type="text" class="input input-bordered w-full" />
    </label>

    <div class="w-full flex justify-end">
      <SubmitButton class="btn btn-success btn-sm" {isLoading}>
        Save
      </SubmitButton>
    </div>
  </form>
</div>
