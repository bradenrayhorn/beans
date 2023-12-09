<script lang="ts">
  import IconGeneral from "~icons/mdi/settings";
  import IconCategories from "~icons/mdi/category";
  import IconPayees from "~icons/mdi/person";
  import { paths, withParameter } from "$lib/paths";
  import { page } from "$app/stores";
</script>

<ul class="menu menu-vertical bg-base-200 rounded-box">
  <li>
    <a href={withParameter(paths.budget.settings.general, $page.params)}>
      General <IconGeneral />
    </a>
  </li>
  <li>
    <a
      href={withParameter(paths.budget.settings.categories.base, $page.params)}
    >
      Categories <IconCategories />
    </a>
  </li>
  <li>
    <a href={withParameter(paths.budget.settings.payees.base, $page.params)}>
      Payees <IconPayees />
    </a>
  </li>
</ul>
