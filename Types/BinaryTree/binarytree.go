package binarytree

func RemoveAtIndex[T any](array *[]T, index uint64) T {
	values := *array
	item := values[index]
	if uint64(len(values)) == index {
		values = values[:index]
		*array = values
		return item
	}
	values = append(values[:index], values[index+1:]...)
	*array = values
	return item
}

func InsertIntoSorted[T any](sorted *[]T, item T, compare func(T, T) int8) {
	sortedValue := *sorted
	length := len(sortedValue)

	if length == 0 {
		*sorted = []T{item}
		return
	}
	lowerBounds := uint64(0)
	upperBounds := uint64(length)

	for {
		index := uint64((lowerBounds + upperBounds) / 2)
		if upperBounds == lowerBounds {
			insert(sorted, index, item)
			return
		}
		comparison := compare(item, sortedValue[index])
		if comparison == 0 {
			insert(sorted, index, item)
			return
		} else if comparison > 0 {
			upperBounds = index
		} else if comparison < 0 {
			lowerBounds = index + 1
		}

	}

}
func insert[T any](array *[]T, index uint64, item T) bool {
	arrayValues := *array
	length := uint64(len(arrayValues))
	if length < index {
		return false
	}
	arrayValues = append(arrayValues[:index], append([]T{item}, arrayValues[index:]...)...)
	*array = arrayValues
	return true
}
