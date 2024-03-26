<script lang="ts">
  import { paths } from "$lib/paths";
  import IconBack from "~icons/mdi/navigate-before";
  import { enhance } from "$app/forms";
  import { navigating } from "$app/stores";
  import FormError from "$lib/components/FormError.svelte";
  import SubmitButton from "$lib/components/SubmitButton.svelte";

  let isSubmitting = false;
  $: isLoading = !!$navigating || isSubmitting;
</script>

<div class="max-w-md w-full mx-auto pt-8 px-4">
  <h1 class="mb-8 text-2xl font-bold">New Budget</h1>

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

    <div class="w-full flex flex-row justify-between">
      <a href={paths.budgets.list} class="btn btn-sm"><IconBack />Back</a>
      <SubmitButton class="btn btn-primary btn-sm" {isLoading}
        >Save</SubmitButton
      >
    </div>
  </form>
</div>
