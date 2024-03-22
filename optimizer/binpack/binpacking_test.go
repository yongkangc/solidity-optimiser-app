package binpack_test

// testing file for binpacking.go
import (
	"optimizer/optimizer/optimizer/binpack"
	"reflect"
	"testing"
)

func TestOptimalBinPackingHelper(t *testing.T, sizes []int, binCapacity int, expected [][]int) {
	pairs := []binpack.Item{}
	for i, size := range sizes {
		pairs = append(pairs, binpack.Item{i, size})
	}

	actual := binpack.OptimalBinPacking(pairs, binCapacity)
	actualSizes := [][]int{}
	for _, slot := range actual {
		sizes := []int{}
		for _, item := range slot {
			sizes = append(sizes, item.Size)
		}
	}

	if !reflect.DeepEqual(actualSizes, expected) {
		t.Errorf("expected %v, got %v\n", expected, actualSizes)
	}
}

func TestOptimalBinPacking(t *testing.T) {
	t.Run("bin cap 10", func(t *testing.T) {
		input := []int{9, 8, 2, 2, 5, 4}
		binCapacity := 10
		expected := [][]int{{9}, {8, 2}, {5, 4}}
		TestOptimalBinPackingHelper(t, input, binCapacity, expected)
	})
}
