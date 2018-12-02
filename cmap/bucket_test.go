package cmap

// genTestingPairs 用于生成测试用的键-元素对的切片。
func genTestingPairs(number int) []Pair {
	testCases := make([]Pair, number)
	for i := 0; i < number; i++ {
		testCases[i], _ = newPair(randString(), randElement())
	}
	return testCases
}

// genNoRepetitiveTestingPairs 用于生成测试用的无重复的键-元素对的切片。
func genNoRepetitiveTestingPairs(number int) []Pair {
	testCases := make([]Pair, number)
	m := make(map[string]struct{})
	var p Pair
	for i := 0; i < number; i++ {
		for {
			p, _ = newPair(randString(), randElement())
			if _, ok := m[p.Key()]; !ok {
				testCases[i] = p
				m[p.Key()] = struct{}{}
				break
			}
		}
	}
	return testCases
}
