<script lang="ts">
  import { page } from "$app/stores";
  import { paths, withParameter } from "$lib/paths";
  import type { PageData } from "./$types";

  export let data: PageData;
</script>

<div class="text-sm breadcrumbs mb-8">
  <ul>
    <li>
      <a href={withParameter(paths.budget.settings.base, $page.params)}>
        Settings
      </a>
    </li>
    <li>General</li>
  </ul>
</div>

<div class="flex flex-col gap-12">
  <div class="shadow bg-base-100 rounded-md p-4 md:p-8">
    <h2 class="font-bold mb-8">User Information</h2>

    <div class="text-base-content-light text-sm mb-8">
      Your username cannot be changed.
    </div>

    <div class="text-sm flex items-center">
      <div class="mr-16">Username</div>

      <div class="shrink">
        <input
          class="input input-bordered input-sm w-full"
          disabled
          value={data.user.username}
        />
      </div>
    </div>
  </div>

  <div class="shadow bg-base-100 rounded-md p-4 md:p-8">
    <h2 class="font-bold mb-8">Budget</h2>

    <div class="text-base-content-light text-sm mb-8">
      If you use more than one budget or wish to create a new budget, you can go
      back to the budget selector.
    </div>

    <div class="text-sm flex items-center">
      <div class="mr-16">Active budget</div>

      <div class="shrink">
        <input
          class="input input-bordered input-sm w-full"
          disabled
          value={data.budget.name}
        />
      </div>
    </div>

    <div class="mt-8">
      <a href={paths.budgets.list} class="text-sm link link-primary"
        >Back to budget selector</a
      >
    </div>
  </div>

  <div class="shadow bg-base-100 rounded-md p-4 md:p-8">
    <h2 class="font-bold mb-8">Active Session</h2>

    <div class="text-base-content-light text-sm">
      Logging out will require you to enter your username and password again.
    </div>

    <a
      href={paths.logout}
      class="btn btn-error btn-sm mt-8"
      data-sveltekit-reload
    >
      Log out
    </a>
  </div>
</div>
