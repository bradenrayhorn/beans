<script lang="ts">
  import { paths, withParameter } from "$lib/paths";
  import IconBack from "~icons/mdi/navigate-before";
  import { enhance } from "$app/forms";
  import { navigating, page } from "$app/stores";
  import FormError from "$lib/components/FormError.svelte";
  import SubmitButton from "$lib/components/SubmitButton.svelte";

  let isSubmitting = false;
  $: isLoading = !!$navigating || isSubmitting;
</script>

<div class="max-w-md w-full mx-auto p-4">
  <h1 class="mb-8 text-2xl font-bold">New Account</h1>

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

    <label class="label cursor-pointer">
      <span class="label-text">Off Budget</span>
      <input name="off_budget" type="checkbox" class="checkbox" value="true" />
    </label>

    <div class="w-full flex flex-row justify-between">
      <a
        href={withParameter(paths.budget.accounts.base, $page.params)}
        class="btn btn-sm"
      >
        <IconBack />Back
      </a>
      <SubmitButton class="btn btn-success btn-sm" {isLoading}>
        Save
      </SubmitButton>
    </div>
  </form>
</div>
