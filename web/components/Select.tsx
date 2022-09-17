import { ChevronDownIcon } from "@chakra-ui/icons";
import {
  Box,
  Button,
  Input,
  InputGroup,
  InputRightElement,
  Popover,
  PopoverAnchor,
  PopoverContent,
  Spinner,
  Text,
} from "@chakra-ui/react";
import { useCombobox } from "downshift";
import { useEffect, useState } from "react";
import { useController } from "react-hook-form";

interface Props<ItemType> {
  name: string;
  itemToString: (item: ItemType | undefined | null) => string;
  itemToID: (item: ItemType | undefined) => string;
  isLoading: boolean;
  items: Array<ItemType>;
}

const Select = <T extends unknown>({
  name,
  itemToString,
  itemToID,
  isLoading,
  items: providedItems,
}: Props<T>) => {
  const {
    field: { onChange, onBlur, value, ref },
  } = useController({ name });

  const [selectedItem, setSelectedItem] = useState(value ?? null);
  const [items, setItems] = useState(providedItems);

  useEffect(() => {
    setItems(providedItems);
  }, [providedItems]);

  const {
    isOpen,
    getInputProps,
    getMenuProps,
    getComboboxProps,
    getItemProps,
    closeMenu,
    openMenu,
    highlightedIndex,
    setInputValue,
  } = useCombobox({
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
      if (!stateChange.isOpen) {
        setInputValue(itemToString(stateChange.selectedItem));
      }
    },
    onSelectedItemChange: (changes) => {
      setSelectedItem(changes.selectedItem);
      onChange(changes.selectedItem);
    },
  });

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
        <Box w="full" {...getComboboxProps()}>
          <PopoverAnchor>
            <InputGroup>
              <Input
                {...getInputProps({
                  onClick: () => openMenu(),
                  onBlur,
                  ref,
                })}
              />
              <InputRightElement
                pointerEvents="none"
                children={<ChevronDownIcon />}
              />
            </InputGroup>
          </PopoverAnchor>

          <PopoverContent
            w="full"
            overflow="hidden"
            boxShadow="dark-lg"
            maxHeight={48}
            overflowY="auto"
            {...getMenuProps()}
          >
            {isLoading ? (
              <Spinner />
            ) : (
              <>
                {items.map((item, key) => (
                  <Button
                    tabIndex={-1}
                    key={`${key}.${itemToID(item)}`}
                    fontWeight={
                      itemToID(selectedItem) === itemToID(item)
                        ? "bold"
                        : "normal"
                    }
                    bg={highlightedIndex === key ? "whiteAlpha.100" : "none"}
                    rounded="none"
                    justifyContent="flex-start"
                    p={2}
                    {...getItemProps({ item, index: key })}
                  >
                    {itemToString(item)}
                  </Button>
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
