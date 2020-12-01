package util

func QuickSort(arr []string) []string {
	if len(arr) < 2 {
		return arr
	}

	i := len(arr) / 2
	left := QuickSort(arr[0:i])
	right := QuickSort(arr[i:])

	return merge(left, right)
}

func merge(l []string, r []string) []string {
	i, j := 0, 0
	m, n := len(l), len(r)
	var res []string
	for i < m && j < n {
		if l[i] < r[j] {
			res = append(res, l[i])
			i++
		} else {
			res = append(res, r[j])
			j++
		}
	}
	res = append(res, l[i:]...)
	res = append(res, r[j:]...)
	return res
}
