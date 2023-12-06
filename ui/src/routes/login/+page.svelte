<script lang="ts">
  import { enhance } from "$app/forms";
  import { navigating } from "$app/stores";
  import FormError from "$lib/components/FormError.svelte";
  import SubmitButton from "$lib/components/SubmitButton.svelte";

  let isSubmitting = false;
  $: isLoading = !!$navigating || isSubmitting;
</script>

<div class="max-w-md w-full mx-auto pt-8 px-4">
  <h1 class="text-center mb-16 text-3xl font-bold">beans</h1>

  <FormError />

  <form
    class="flex flex-col gap-8"
    method="POST"
    action="?/login"
    use:enhance={() => {
      isSubmitting = true;

      return async ({ update }) => {
        await update();
        isSubmitting = false;
      };
    }}
  >
    <label>
      <span class="label label-text">Username</span>
      <input name="username" type="text" class="input input-bordered w-full" />
    </label>

    <label>
      <span class="label label-text">Password</span>
      <input
        name="password"
        type="password"
        class="input input-bordered w-full"
      />
    </label>

    <SubmitButton {isLoading} class="btn btn-primary">Sign In</SubmitButton>
  </form>
</div>
