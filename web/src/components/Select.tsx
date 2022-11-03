import { ChevronDownIcon } from "@chakra-ui/icons";
import {
  Box,
  chakra,
  Flex,
  Input,
  InputGroup,
  InputRightElement,
  Popover,
  PopoverAnchor,
  PopoverContent,
  Spinner,
  Text,
  useMultiStyleConfig,
} from "@chakra-ui/react";
import { useCombobox } from "downshift";
import { useEffect, useState } from "react";
import { useController } from "react-hook-form";

interface Props<ItemType> {
  name: string;
  itemToString: (item: ItemType | undefined | null) => string;
  itemToID: (item: ItemType | undefined) => string;
  isLoading: boolean;
  isOpen?: boolean;
  setIsOpen?: (isOpen: boolean) => void;
  items: Array<ItemType>;
}

export const useAsyncSelect = () => {
  const [isOpen, setIsOpen] = useState(false);
  return { isOpen, selectProps: { isOpen, setIsOpen } };
};

const Select = <T extends unknown>({
  name,
  itemToString,
  itemToID,
  items: providedItems,
  isLoading: parentIsLoading = false,
  isOpen: parentIsOpen,
  setIsOpen: parentSetIsOpen,
}: Props<T>) => {
  const {
    field: { onChange, onBlur, value, ref },
  } = useController({ name });

  const [selectedItem, setSelectedItem] = useState(value ?? null);
  const [isLoading, setIsLoading] = useState(parentIsLoading);
  const [items, setItems] = useState(providedItems);

  useEffect(() => {
    setItems(providedItems);
    setIsLoading(parentIsLoading);
  }, [providedItems, isLoading]);

  const {
    isOpen,
    getInputProps,
    getMenuProps,
    getItemProps,
    closeMenu,
    openMenu,
    setInputValue,
  } = useCombobox({
    isOpen: parentIsOpen,
    selectedItem,
    onInputValueChange: ({ inputValue, ...r }) => {
      if (!inputValue || inputValue === itemToString(r.selectedItem)) {
        setItems(providedItems);
      } else {
        setItems(
          providedItems.filter((item) =>
            itemToString(item)
              ?.toLowerCase()
              ?.startsWith((inputValue ?? "").toLowerCase())
          )
        );
      }
    },
    items,
    itemToString,
    onIsOpenChange: (stateChange) => {
      if (parentSetIsOpen) {
        parentSetIsOpen(stateChange.isOpen ?? false);
      }
      if (!stateChange.isOpen) {
        setInputValue(itemToString(stateChange.selectedItem));
      }
    },
    onSelectedItemChange: (changes) => {
      setSelectedItem(changes.selectedItem);
      onChange(changes.selectedItem);
    },
  });

  const styles = useMultiStyleConfig("ComponentSelect");

  return (
    <>
      <Popover
        isLazy
        isOpen={isOpen}
        onClose={closeMenu}
        autoFocus={false}
        matchWidth
        placement="bottom"
        closeOnBlur={false}
      >
        <Box w="full">
          <PopoverAnchor>
            <InputGroup>
              <Input
                {...getInputProps({
                  onClick: () => openMenu(),
                  onBlur,
                  ref,
                })}
              />
              <InputRightElement pointerEvents="none">
                <ChevronDownIcon />
              </InputRightElement>
            </InputGroup>
          </PopoverAnchor>

          <PopoverContent {...styles.wrapper} {...getMenuProps()}>
            {isLoading ? (
              <Flex w="full" p={2} justifyContent="center">
                <Spinner />
              </Flex>
            ) : (
              <>
                {items.map((item, key) => (
                  <chakra.button
                    tabIndex={-1}
                    key={`${key}.${itemToID(item)}`}
                    type="button"
                    aria-checked={itemToID(selectedItem) === itemToID(item)}
                    __css={styles.item}
                    {...getItemProps({ item, index: key })}
                  >
                    {itemToString(item)}
                  </chakra.button>
                ))}

                {items.length === 0 && (
                  <Text as="i" p={2} textAlign="center" w="full">
                    No options
                  </Text>
                )}
              </>
            )}
          </PopoverContent>
        </Box>
      </Popover>
    </>
  );
};

export default Select;
