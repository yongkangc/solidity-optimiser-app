package binpack_test

// testing file for binpacking.go
import (
	"optimizer/optimizer/optimizer/binpack"
	"reflect"
	"testing"
)

func slotsToSizes(slots []binpack.Slot) [][]int {
	sizes := [][]int{}
	for _, slot := range slots {
		sizes = append(sizes, itemToSize(slot))
	}
	return sizes
}

func itemToSize(items []binpack.Item) []int {
	sizes := []int{}
	for _, item := range items {
		sizes = append(sizes, item.Size)
	}
	return sizes
}

func TestOptimalBinPacking(t *testing.T) {
	TestOptimalBinPackingHelper := func(t *testing.T, sizes []int, binCapacity int, expectedSizes [][]int) {
		items := []binpack.Item{}
		for i, size := range sizes {
			items = append(items, binpack.Item{i, size})
		}

		actualSlots := binpack.OptimalBinPacking(items, binCapacity)
		actualSizes := slotsToSizes(actualSlots)

		if !reflect.DeepEqual(expectedSizes, actualSizes) {
			t.Errorf("expected %v, got %v\n", expectedSizes, actualSizes)
		}
	}

	t.Run("Bin capacity 10", func(t *testing.T) {
		input := []int{9, 8, 2, 2, 5, 4}
		binCapacity := 10
		expected := [][]int{{9}, {8, 2}, {5, 4}, {2}}
		TestOptimalBinPackingHelper(t, input, binCapacity, expected)
	})
}
