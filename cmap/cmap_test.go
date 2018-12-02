package cmap

import "testing"

func TestCmapNew(t *testing.T) {
	var concurrency int
	var pairRedistributor PairRedistributor

	cm, err := NewConcurrentMap(concurrency, pairRedistributor)
	if err == nil {
		t.Logf("No error when new a concurrent map with concurrency %d, but should not be the case!",
			concurrency)
	} else {
		t.Logf("New cmap failed, %s", err)
	}

	concurrency = MAX_CONCURRENCY + 1
	cm, err = NewConcurrentMap(concurrency, pairRedistributor)
	if err == nil {
		t.Logf("No error when new a concurrent map with concurrency %d, but should not be the case!",
			concurrency)
	} else {
		t.Logf("New cmap failed, %s", err)
	}

	concurrency = 16
	cm, err = NewConcurrentMap(concurrency, pairRedistributor)
	if err != nil {
		t.Fatalf("An error occurs when new a concurrent map: %s (concurrency: %d, pairRedistributor: %#v)",
			err, concurrency, pairRedistributor)
	}

	t.Logf("The Cmap: %s", cm)
}

func TestCmapPut(t *testing.T) {
	number := 10
	testCases := genNoRepetitiveTestingPairs(number)
	t.Logf("testCases Len: %d, Data: %s", len(testCases), testCases)

	concurrency := 2
	var pairRedistributor PairRedistributor
	cm, _ := NewConcurrentMap(concurrency, pairRedistributor)

	for _, p := range testCases {
		key := p.Key()
		element := p.Element()
		ok, err := cm.Put(key, element)

		if err != nil {
			t.Fatalf("An error occurs when putting a key-element to the cmap: %s (key: %s, element: %#v)",
				err, key, element)
		}
		if !ok {
			t.Fatalf("Couldn't put key-element to the cmap! (key: %s, element: %#v)",
				key, element)
		}

		actualElement := cm.Get(key)
		if actualElement == nil {
			t.Fatalf("Inconsistent element: expected: %#v, actual: %#v",
				element, nil)
		}
	}

	if cm.Len() != uint64(number) {
		t.Fatalf("Inconsistent size: expected: %d, actual: %d",
			number, cm.Len())
	}

	t.Logf("The Cmap: %s", cm)
}
