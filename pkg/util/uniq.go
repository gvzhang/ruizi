package util

func Uniq(ar []string) []string {
	if len(ar) < 2 {
		return ar
	}
	ar = QuickSort(ar)
	var res []string
	res = append(res, ar[0])
	j := 1
	for i := 1; i < len(ar); i++ {
		if ar[i] != res[j-1] {
			res = append(res, ar[i])
			j++
		}
	}
	return res
}
