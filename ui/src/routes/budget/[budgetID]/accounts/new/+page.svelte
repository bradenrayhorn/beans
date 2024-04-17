<script lang="ts">
  import { paths, withParameter } from "$lib/paths";
  import IconError from "~icons/mdi/alert-circle-outline";
  import { page } from "$app/stores";
  import SubmitButton from "$lib/components/SubmitButton.svelte";
  import type { PageData } from "./$types";
  import { superForm } from "sveltekit-superforms";
  import { zodClient } from "sveltekit-superforms/adapters";
  import { schema } from "./schema";

  export let data: PageData;
  const { form, errors, message, submitting, enhance } = superForm(data.form, {
    validators: zodClient(schema),
    invalidateAll: true,
  });
</script>

<div class="w-full bg-base-200 p-8 md:pt-12">
  <div class="max-w-md w-full mx-auto">
    <h1 class="mb-8 text-2xl font-bold">New Account</h1>

    <form class="flex flex-col gap-8" method="POST" action="?/save" use:enhance>
      {#if $message}
        <div
          role="alert"
          class="rounded px-4 py-2 flex gap-2 items-center bg-error/20 text-error font-semibold text-sm"
        >
          <IconError class="text-lg" />
          {$message}
        </div>
      {/if}

      <label class="form-control">
        <span class="label label-text">Name</span>
        <input
          name="name"
          type="text"
          class="input input-bordered w-full"
          autocomplete="off"
          class:input-error={$errors.name}
          aria-invalid={$errors.name ? "true" : undefined}
          bind:value={$form.name}
        />
        {#if $errors.name}
          <div class="text-xs pt-2">
            <span class="label-text-alt text-error">{$errors.name}</span>
          </div>
        {/if}
      </label>

      <label class="label cursor-pointer max-w-32">
        <span class="label-text">Off Budget</span>
        <input
          name="offBudget"
          type="checkbox"
          class="checkbox"
          bind:checked={$form.offBudget}
        />
      </label>

      <SubmitButton
        class="btn btn-primary btn-sm w-full"
        isLoading={$submitting}
      >
        Save
      </SubmitButton>
    </form>

    <a
      href={withParameter(paths.budget.accounts.base, $page.params)}
      class="btn btn-secondary btn-sm w-full mt-4"
    >
      Cancel
    </a>
  </div>
</div>
