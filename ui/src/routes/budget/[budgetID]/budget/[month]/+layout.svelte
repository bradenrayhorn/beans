<script lang="ts">
  import dayjs from "dayjs";
  import type { PageData } from "./$types";
  import IconNext from "~icons/mdi/ChevronRight";
  import IconBack from "~icons/mdi/ChevronLeft";
  import { paths, withParameter } from "$lib/paths";
  import { page } from "$app/stores";
  import BudgetTable from "./BudgetTable.svelte";
  import BudgetCards from "./BudgetCards.svelte";
  import { goto } from "$app/navigation";
  import { keybind } from "$lib/actions/keybind";

  export let data: PageData;

  $: isList =
    $page.url.pathname ===
    withParameter(paths.budget.budget.month, $page.params);

  $: date = dayjs(data.month.date);
  $: month = date.format("MMMM");
  $: year = date.format("YYYY");

  $: pastMonthLink = withParameter(paths.budget.budget.month, {
    ...$page.params,
    month: date.add(-1, "month").format("YYYY-MM"),
  });
  $: nextMonthLink = withParameter(paths.budget.budget.month, {
    ...$page.params,
    month: date.add(1, "month").format("YYYY-MM"),
  });

  function toPastMonth() {
    goto(pastMonthLink);
  }
  function toNextMonth() {
    goto(nextMonthLink);
  }
</script>

<div class="flex grow h-full">
  <div class="grow flex flex-col">
    {#if !isList}
      <div class="flex flex-col md:hidden">
        <slot />
      </div>
    {/if}

    <div
      class="flex flex-col md:flex-row md:justify-between items-center shrink-0 px-4"
      class:hidden={!isList}
      class:md:flex={!isList}
      class:flex={isList}
    >
      <div class="py-2 md:py-6 flex items-center gap-2">
        <a
          class="btn btn-xs btn-ghost"
          href={pastMonthLink}
          aria-label="Previous month"
          use:keybind={{ key: "h", action: toPastMonth }}
        >
          <IconBack />
        </a>

        <div class="flex gap-2">
          <h1 class="text-lg md:text-3xl font-bold">{month}</h1>
          <h2 class="text-lg font-bold md:text-sm md:font-normal">{year}</h2>
        </div>

        <a
          class="btn btn-xs btn-ghost"
          href={nextMonthLink}
          aria-label="Next month"
          use:keybind={{ key: "l", action: toNextMonth }}
        >
          <IconNext />
        </a>
      </div>

      <div class="flex w-full flex-col items-center md:items-end">
        <div class="dropdown dropdown-bottom md:dropdown-end">
          <button class="btn btn-sm btn-ghost">
            <b>To Budget:</b>
            <span>{data.month.budgetable.display}</span>
          </button>
          <div
            class="card compact dropdown-content z-10 shadow bg-base-100 rounded w-64 -ml-10 md:ml-0 p-4 text-sm"
            role="dialog"
            aria-label="To budget breakdown"
          >
            <div class="flex justify-between w-full">
              <b>Income</b>
              <span>{data.month.income.display}</span>
            </div>

            <div class="flex justify-between w-full">
              <b>From last month</b>
              <span>{data.month.carriedOver.display}</span>
            </div>

            <div class="flex justify-between w-full">
              <b>Assigned this month</b>
              <span>-{data.month.assigned.display}</span>
            </div>

            <div class="flex justify-between w-full">
              <b>For next month</b>
              <span>-{data.month.carryover.display}</span>
            </div>
          </div>
        </div>

        <a
          href={withParameter(paths.budget.budget.forNextMonth, $page.params)}
          class="text-xs md:mr-3 link link-primary"
        >
          Save for next month
        </a>
      </div>
    </div>

    <div class="hidden md:block grow bg-base-100">
      <BudgetTable month={data.month} categoryGroups={data.categoryGroups} />
    </div>

    {#if isList}
      <div class="md:hidden grow">
        <BudgetCards month={data.month} categoryGroups={data.categoryGroups} />
      </div>
    {/if}
  </div>

  {#if !isList}
    <div class="hidden md:block shrink-0 bg-base-100 h-full shadow p-4 w-64">
      <slot />
    </div>
  {/if}
</div>
