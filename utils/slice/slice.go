package slice

import (
	"sort"
	"strconv"
)

// StrInSlice 判断string是否在slice中
func StrInSlice(str string, arr []string) bool {
	if len(arr) == 0 {
		return false
	}
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// IntInSlice 判断int是否在slice中
func IntInSlice(num int, arr []int) bool {
	if len(arr) == 0 {
		return false
	}
	for _, v := range arr {
		if v == num {
			return true
		}
	}
	return false
}

// Interface2String interface的slice转string的slice
func Interface2String(arr []interface{}) []string {
	length := len(arr)
	strSlice := make([]string, length)
	if length == 0 {
		return strSlice
	}
	for k, v := range arr {
		strSlice[k] = v.(string)
	}
	return strSlice
}

// String2Int string的slice转int
func String2Int(arr []string) []int {
	length := len(arr)
	intSlice := make([]int, length)
	if length == 0 {
		return intSlice
	}
	for k, v := range arr {
		var ok error
		intSlice[k], ok = strconv.Atoi(v)
		if ok != nil {
			intSlice[k] = 0
		}
	}
	return intSlice
}

// Int642Uint32 int64转uint32
func Int642Uint32(arr []int64) []uint32 {
	uint32Slice := make([]uint32, len(arr))
	for _, v := range arr {
		uint32Slice = append(uint32Slice, uint32(v))
	}
	return uint32Slice
}

// String2Uint32 string的slice转uint32
func String2Uint32(arr []string) []uint32 {
	length := len(arr)
	intSlice := make([]uint32, length)
	if length == 0 {
		return intSlice
	}
	for k, v := range arr {
		tmp, _ := strconv.ParseInt(v, 10, 32)
		intSlice[k] = uint32(tmp)
	}
	return intSlice
}

func StrSliceRemove(slice []string, i int) []string {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func ReverseStringSlice(slice []string) []string {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

// 用于uid去重
type Uint32Slice []uint32

func (p Uint32Slice) Len() int           { return len(p) }
func (p Uint32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func RemoveRep(slc []uint32) []uint32 {
	if len(slc) <= 1 {
		return slc
	}
	sort.Sort(Uint32Slice(slc))

	var d int
	for i := 1; i < len(slc); i++ {
		if slc[d] != slc[i] {
			d++
			slc[d] = slc[i]
		}
	}
	return slc[:d+1]
}

type Uint64Slice []uint64

func (p Uint64Slice) Len() int           { return len(p) }
func (p Uint64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func RemoveRepUint64(slc []uint64) []uint64 {
	if len(slc) <= 1 {
		return slc
	}
	sort.Sort(Uint64Slice(slc))

	var d int
	for i := 1; i < len(slc); i++ {
		if slc[d] != slc[i] {
			d++
			slc[d] = slc[i]
		}
	}
	return slc[:d+1]
}

// func Uint32Diff(s1, s2 []uint32) []uint32 {
// 	l1 := len(s1)
// 	l2 := len(s2)
// 	var maxLen int
// 	if l1 >= l2 {
// 		maxLen = l1
// 	} else {
// 		maxLen = l2
// 	}
// 	returnSlice := make([]uint32, 0, maxLen)
// 	for _, v1 := range s1 {
// 		for _, v2 := range s2 {
// 			if v1 == v2 {
// 				continue
// 			}
// 			returnSlice = append(returnSlice, v2)
// 		}
// 	}
// }
