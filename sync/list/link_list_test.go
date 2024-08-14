package _list

import (
	"testing"
)

func TestLinkList(t *testing.T) {
	// Test case 1: Test the basic operations of the linked list
	linkList := NewLinkList[int]()
	linkList.PushBack(1)
	linkList.PushBack(2)
	linkList.PushBack(3)
	if linkList.Len() != 3 {
		t.Errorf("Expected length 3, but got %d", linkList.Len())
	}
	if linkList.Front() != 1 {
		t.Errorf("Expected front element 1, but got %d", linkList.Front())
	}
	if linkList.Back() != 3 {
		t.Errorf("Expected back element 3, but got %d", linkList.Back())
	}

	// Test case 2: Test the PopBack operation [1, 2, 3]
	if linkList.PopBack() != 3 {
		t.Errorf("Expected popped element 3, but got %d", linkList.PopBack())
	}
	if linkList.Len() != 2 {
		t.Errorf("Expected length 2, but got %d", linkList.Len())
	}

	// Test case 3: Test the PopFront operation [1, 2]
	if linkList.PopFront() != 1 {
		t.Errorf("Expected popped element 1, but got %d", linkList.PopFront())
	}
	if linkList.Len() != 1 {
		t.Errorf("Expected length 1, but got %d", linkList.Len())
	}

	// Test case 4: Test the PushFront operation [2]
	linkList.PushFront(0)
	if linkList.Front() != 0 {
		t.Errorf("Expected front element 0, but got %d", linkList.Front())
	}
	if linkList.Len() != 2 {
		t.Errorf("Expected length 2, but got %d", linkList.Len())
	}

	// Test case 5: Test the Empty operation [0, 2]
	if linkList.Empty() {
		t.Errorf("Expected linked list to be empty, but it is not")
	}

	// Test case 6: Test the linked list with multiple types
	linkList2 := NewLinkList[string]()
	linkList2.PushBack("a")
	linkList2.PushBack("b")
	linkList2.PushBack("c")
	if linkList2.Len() != 3 {
		t.Errorf("Expected length 3, but got %d", linkList2.Len())
	}
	if linkList2.Front() != "a" {
		t.Errorf("Expected front element 'a', but got %s", linkList2.Front())
	}
	if linkList2.Back() != "c" {
		t.Errorf("Expected back element 'c', but got %s", linkList2.Back())
	}
	if linkList2.PopBack() != "c" {
		t.Errorf("Expected popped element 'c', but got %s", linkList2.PopBack())
	}
	if linkList2.Len() != 2 {
		t.Errorf("Expected length 2, but got %d", linkList2.Len())
	}
	if linkList2.PopFront() != "a" {
		t.Errorf("Expected popped element 'a', but got %s", linkList2.PopFront())
	}
	if linkList2.Len() != 1 {
		t.Errorf("Expected length 1, but got %d", linkList2.Len())
	}
	linkList2.PushFront("0")
	if linkList2.Front() != "0" {
		t.Errorf("Expected front element '0', but got %s", linkList2.Front())
	}
	if linkList2.Len() != 2 {
		t.Errorf("Expected length 2, but got %d", linkList2.Len())
	}
	if linkList2.Empty() {
		t.Errorf("Expected linked list to be empty, but it is not")
	}
}
