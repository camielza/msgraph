package slices

type IsEqual func(item1 interface{}, item2 interface{}) bool

func Union(slice1 []interface{}, slice2 []interface{}, equalityFunc IsEqual) []interface{} {
	var res = []interface{}{}
	for _, s1 := range slice1 {
		for _, s2 := range slice2 {
			if equalityFunc(s1, s2) {
				res = append(res, s1)
			}
		}
	}
	return res
}
