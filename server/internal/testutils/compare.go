package testutils

import (
	"reflect"
	"sort"
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/stretchr/testify/require"
)

func IsEqualInAnyOrder[T any](tb testing.TB, sliceA []T, sliceB []T, comp func(a, b T) bool) {
	require.Equal(tb, len(sliceA), len(sliceB))

	sort.Slice(sliceA, func(i, j int) bool {
		return comp(sliceA[i], sliceA[j])
	})
	sort.Slice(sliceB, func(i, j int) bool {
		return comp(sliceB[i], sliceB[j])
	})

	for i := range sliceA {
		require.Truef(tb, reflect.DeepEqual(sliceA[i], sliceB[i]), "not equal: a: %v, b: %v", sliceA[i], sliceB[i])
	}
}

// Comparison functions

func CmpAccount(a *beans.Account, b *beans.Account) bool {
	return a.ID.String() < b.ID.String()
}

func CmpCategory(a *beans.Category, b *beans.Category) bool {
	return a.ID.String() < b.ID.String()
}

func CmpCategoryGroup(a *beans.CategoryGroup, b *beans.CategoryGroup) bool {
	return a.ID.String() < b.ID.String()
}

func CmpMonthCategory(a *beans.MonthCategory, b *beans.MonthCategory) bool {
	return a.ID.String() < b.ID.String()
}
