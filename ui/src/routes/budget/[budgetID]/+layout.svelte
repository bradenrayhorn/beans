<script lang="ts">
  import IconBudget from "~icons/mdi/wallet";
  import IconLedger from "~icons/mdi/invoice-text-outline";
  import IconAccounts from "~icons/mdi/account-balance";
  import IconSettings from "~icons/mdi/settings";
  import { paths, withParameter } from "$lib/paths";
  import { page } from "$app/stores";
  import NavItem from "$lib/components/NavItem.svelte";
  import Lightswitch from "$lib/components/Lightswitch.svelte";
  import { env } from "$env/dynamic/public";

  const routes = [
    {
      path: paths.budget.budget.base,
      name: "Budget",
      icon: IconBudget,
    },
    {
      path: paths.budget.ledger.base,
      name: "Ledger",
      icon: IconLedger,
    },
    {
      path: paths.budget.accounts.base,
      name: "Accounts",
      icon: IconAccounts,
    },
    {
      path: paths.budget.settings.base,
      name: "Settings",
      icon: IconSettings,
    },
  ];
  $: builtRoutes = routes.map((route) => {
    const path = withParameter(route.path, $page.params);
    return {
      ...route,
      path,
      isActive: $page.url.pathname.startsWith(path),
    };
  });
  $: navRoutes = builtRoutes.slice().splice(0, builtRoutes.length - 1);
  $: settingsRoute = builtRoutes[builtRoutes.length - 1]!;
</script>

<div class="flex flex-col w-full min-h-screen md:flex-row">
  <!-- Side navigation (desktop) -->
  <div
    class="hidden md:flex w-48 shadow-md z-10 flex-col bg-neutral text-neutral-content fixed top-0 bottom-0"
  >
    <div class="px-4 py-4 font-semibold">beans</div>

    <div class="flex flex-col justify-between h-full">
      <div class="flex flex-col">
        {#each navRoutes as route}
          <NavItem {...route} />
        {/each}
      </div>

      <div class="flex flex-col">
        <NavItem {...settingsRoute} />

        <div class="flex justify-between px-4 py-2 items-center">
          <div class="text-xs">v.{env.PUBLIC_VERSION ?? "local"}</div>
          <Lightswitch />
        </div>
      </div>
    </div>
  </div>

  <div class="hidden md:block w-48 h-full shrink-0"></div>

  <!-- Content -->
  <div class="grow bg-base-300">
    <slot />
  </div>

  <div class="flex md:hidden w-full shrink-0 h-20"></div>

  <!-- Bottom navigation (mobile) -->
  <div
    class="md:hidden flex shrink-0 justify-between items-center h-20 shadow-top bg-base-100 rounded-t-md fixed bottom-0 right-0 left-0"
  >
    {#each builtRoutes as route (route.name)}
      <a
        class="flex grow flex-col items-center"
        class:text-primary={route.isActive}
        class:text-base-content-light={!route.isActive}
        href={route.path}
      >
        <svelte:component this={route.icon} />
        <span class="text-sm">{route.name}</span>
      </a>
    {/each}
  </div>
</div>
