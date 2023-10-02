package array

import (
	"bytes"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/mochammadshenna/arch-pba-template/internal/util/json"
)

func Contains[T comparable](data []T, v T) (int, bool) {
	for i, d := range data {
		if d == v {
			return i, true
		}
	}
	return -1, false
}

// ConvertStringToInt64 will parse comma separated numbers into []int64
// The non numerical value will be ignored
func ConvertStringToInt64(s string) []int64 {
	x := strings.Split(s, ",")
	res := []int64{}

	for _, v := range x {
		z, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			res = append(res, z)
		}
	}
	return res
}

// UniqueInt64 will return unique member with random sequence into new array
// this function shall be used before getting multiple class with multiple ID
func UniqueInt64(a []int64) []int64 {
	if len(a) == 1 {
		return a
	}

	dataMap := make(map[int64]bool)

	result := []int64{}
	for _, val := range a {
		if dataMap[val] == false {
			result = append(result, val)
			dataMap[val] = true
		}
	}
	return result
}

// UniqueString will return unique member with random sequence into new array
// this function shall be used before getting multiple class with multiple ID
func UniqueString(a []string) []string {
	if len(a) == 1 {
		return a
	}

	dataMap := make(map[string]bool)

	result := []string{}
	for _, val := range a {
		if dataMap[val] == false {
			result = append(result, val)
			dataMap[val] = true
		}
	}
	return result
}

// SafeExtract is the way to get the array by index without having fear of panics
func SafeExtract(a []int64, i int) int64 {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	return a[i]
}

// MergeDistinct is a func to merge distinct two slices
func MergeDistinct(source []int64, dest []int64) (res []int64) {
	resMap := make(map[int64]bool, 0)
	for _, s := range source {
		if _, ok := resMap[s]; ok {
			continue
		}
		resMap[s] = true
	}
	for _, d := range dest {
		if _, ok := resMap[d]; ok {
			continue
		}
		resMap[d] = true
	}
	for k := range resMap {
		res = append(res, k)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})

	return res
}

// ChunkSliceString is needed to make slice of string to smaller size
// smoothChunk is used for allocating chunk into smoother size
func ChunkSliceString(productsCache []string, divider int, smoothChunk bool) (divided [][]string) {
	lenSlice := len(productsCache)
	chunkSize := divider
	if smoothChunk {
		chunkSize = (lenSlice + divider - 1) / divider
	}
	for i := 0; i < lenSlice; i += chunkSize {
		end := i + chunkSize
		if end > lenSlice {
			end = lenSlice
		}
		divided = append(divided, productsCache[i:end])
	}
	return divided
}

// ChunkSliceInt64 is needed to make slice to smaller size
// smoothChunk is used for allocating chunk into smoother size
func ChunkSliceInt64(input []int64, divider int, smoothChunk bool) (divided [][]int64) {
	lenSlice := len(input)
	chunkSize := divider
	if smoothChunk {
		chunkSize = (lenSlice + divider - 1) / divider
	}
	for i := 0; i < lenSlice; i += chunkSize {
		end := i + chunkSize
		if end > lenSlice {
			end = lenSlice
		}
		divided = append(divided, input[i:end])
	}
	return divided
}

// RemoveNegativeValues will remove negative values from the array provided.
func RemoveNegativeValues(input []int64, removeZero bool) []int64 {
	result := []int64{}
	for _, val := range input {
		if val == 0 && !removeZero {
			result = append(result, val)
		}
		if val > 0 {
			result = append(result, val)
		}
	}
	return result
}

// SortStableInt64 is function to sort int64 ascending with stable sorting
func SortStableInt64(input []int64) {
	sort.SliceStable(input, func(i, j int) bool {
		return input[i] < input[j]
	})
}

// Remove is function to remove the first occurence of specified element
func Remove(input []int64, ele int64) []int64 {
	for i := 0; i < len(input); i++ {
		if input[i] == ele {
			return append(input[:i], input[i+1:]...)
		}
	}
	return input
}

// comparable element difference return elements in `a` but not in `b`
func ElementDifference[T comparable](a, b []T) []T {
	mb := make(map[T]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []T
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// StringSafeExtract is the way to get the array by index without having fear of panics
func StringSafeExtract(a []string, i int) string {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	return a[i]
}

func DistinctElement[T comparable](data []T) []T {
	m := make(map[T]struct{})
	for _, d := range data {
		m[d] = struct{}{}
	}

	res := make([]T, len(m))
	i := 0
	for k := range m {
		res[i] = k
		i++
	}
	return res
}

func Chunk[T comparable](data []T, chunkSize int) [][]T {
	size := len(data) / int(chunkSize)
	chunks := make([][]T, size)

	for i := range chunks {
		chunks[i] = make([]T, chunkSize)
	}

	mod := len(data) % chunkSize
	if mod > 0 {
		chunks = append(chunks, make([]T, mod))
	}

	for i, d := range data {
		x := i / chunkSize
		y := i % chunkSize
		chunks[x][y] = d
	}

	return chunks
}

func IsEqualIntSlice(x, y []int64) bool {
	if len(x) != len(y) {
		return false
	}
	// create a map of string -> int
	diff := make(map[int64]int, len(x))
	for _, _x := range x {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}
	for _, _y := range y {
		// If the string _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y] -= 1
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	return len(diff) == 0
}

func Construct[T any](v interface{}) []T {
	res := make([]T, 0)

	if v == nil {
		return res
	}

	switch reflect.TypeOf(v).Kind() {
	case reflect.Array, reflect.Slice:
		val := reflect.ValueOf(v)
		for i := 0; i < val.Len(); i++ {
			uVal := val.Index(i).Interface()
			switch reflect.TypeOf(uVal).Kind() {
			case reflect.Map:
				var element T
				var b bytes.Buffer
				err := json.NewEncoder(&b).Encode(uVal)
				if err != nil {
					continue
				}

				err = json.Unmarshal(b.Bytes(), &element)
				if err != nil {
					continue
				}

				res = append(res, element)
			default:
				val, ok := uVal.(T)
				if ok {
					res = append(res, val)
				}
			}
		}
	}

	return res
}

func Filter[T comparable](data []T, predicate func(T) bool) []T {
	var result []T
	for _, datum := range data {
		if predicate(datum) {
			result = append(result, datum)
		}
	}

	return result
}

func InArray[T comparable](val T, arr []T) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}

	return false
}

func CheckElements[T comparable](arrayA, arrayB []T) []T {
	elements := make([]T, 0)

	// Create a map of elements in arrayB for faster lookup
	lookup := make(map[T]struct{})
	for _, item := range arrayB {
		lookup[item] = struct{}{}
	}

	// Check each element in arrayA if it exists in arrayB
	for _, item := range arrayA {
		if _, ok := lookup[item]; ok {
			elements = append(elements, item)
		}
	}

	return elements
}
