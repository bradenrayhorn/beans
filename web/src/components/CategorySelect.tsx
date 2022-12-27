import Select, { useAsyncSelect } from "@/components/Select";
import { useCategories } from "@/data/queries/category";

type Props = {
  name: string;
};

const CategorySelect = ({ name }: Props) => {
  const { selectProps } = useAsyncSelect();
  const { categoryGroups, isLoading } = useCategories();
  const categories =
    categoryGroups?.flatMap((categoryGroup) => categoryGroup.categories) ?? [];

  return (
    <Select
      name={name}
      itemToString={(item) => item?.name ?? ""}
      itemToID={(item) => item?.id ?? ""}
      isLoading={isLoading}
      items={categories}
      isClearable={true}
      {...selectProps}
    />
  );
};

export default CategorySelect;
