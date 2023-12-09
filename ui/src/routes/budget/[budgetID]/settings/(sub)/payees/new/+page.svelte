<script lang="ts">
  import { paths, withParameter } from "$lib/paths";
  import { enhance } from "$app/forms";
  import { navigating, page } from "$app/stores";
  import FormError from "$lib/components/FormError.svelte";
  import SubmitButton from "$lib/components/SubmitButton.svelte";

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
      <a href={withParameter(paths.budget.settings.payees.base, $page.params)}>
        Payees
      </a>
    </li>
    <li>New Payee</li>
  </ul>
</div>

<div class="max-w-md w-full mx-auto p-4">
  <h2 class="font-bold mb-4 text-lg">New Payee</h2>

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
