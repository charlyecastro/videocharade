package indexes

import (
	"reflect"
	"testing"
)

//implement mutex -- store returned variable, unlock, return variable

func TestFind(t *testing.T) {
	//Return empty list because wor doesnt exist
	testTrie := NewTrie()
	list := testTrie.Find("nothin", 5)
	if len(list) != 0 {
		t.Errorf("incorrect list: expected empty list")
	}

	//simple add word and return a single value list
	testTrie.Add("Charlye", 1)
	size := testTrie.Len()
	if size != 1 {
		t.Errorf("incorrect size: expected %d but got %d",
			1, size)
	}

	//simple prefix
	list = testTrie.Find("charl", 1)
	expectedList := []int64{}
	expectedList = append(expectedList, 1)
	if !reflect.DeepEqual(list, expectedList) {
		t.Errorf("incorrect list")
	}

	//longer prefix and capital
	list = testTrie.Find("CHARLYE", 1)
	if list[0] != 1 {
		t.Errorf("incorrect size: expected %d but got %d",
			1, list[0])
	}

	testTrie.Add("Charlye", 3)
	list = testTrie.Find("charlYE", 2)
	expectedList = append(expectedList, 3)
	if !reflect.DeepEqual(list, expectedList) {
		t.Errorf("incorrect list")
	}

	//no duplicates
	testTrie.Add("OSCAR", 3)
	testTrie.Add("OSCOLOPE", 5)
	testTrie.Add("OSCARMEN", 4)
	testTrie.Add("MENDOZA", 3)
	list = testTrie.Find("OSC", 5)
	expectedList = []int64{3, 4, 5}
	if !reflect.DeepEqual(list, expectedList) {
		t.Errorf("incorrect list")
	}
}

func TestOrder(t *testing.T) {
	testTrie := NewTrie()
	testTrie.Add("Charlye", 1)
	testTrie.Add("charlis", 2)
	testTrie.Add("Charlie", 3)
	list := testTrie.Find("char", 5)

	expectedList := []int64{3, 2, 1}
	if !reflect.DeepEqual(list, expectedList) {
		t.Errorf("incorrect list")
	}
}

func TestRemove(t *testing.T) {
	testTrie := NewTrie()
	testTrie.Add("test", 22)
	list := testTrie.Find("te", 1)
	if list[0] != 22 {
		t.Errorf("incorrect id: expected %d but got %d",
			22, list[0])
	}

	testTrie.Remove("test", 22)
	list = testTrie.Find("te", 1)
	if len(list) != 0 {
		t.Errorf("incorrect list: expected empty list")
	}
}
