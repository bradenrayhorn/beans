import Select, { SelectProps, useAsyncSelect } from "@/components/Select";
import { Category } from "@/constants/types";
import { useCategories } from "@/data/queries/category";

const CategorySelect = (
  props: Omit<
    SelectProps<Category>,
    "itemToString" | "itemToID" | "isLoading" | "items" | "isClearable"
  >
) => {
  const { selectProps } = useAsyncSelect();
  const { categoryGroups, isLoading } = useCategories();
  const categories =
    categoryGroups?.flatMap((categoryGroup) => categoryGroup.categories) ?? [];

  return (
    <Select
      itemToString={(item) => item?.name ?? ""}
      itemToID={(item) => item?.id ?? ""}
      isLoading={isLoading}
      items={categories}
      isClearable={true}
      {...props}
      {...selectProps}
    />
  );
};

export default CategorySelect;
