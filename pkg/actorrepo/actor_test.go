package actorrepo

import "testing"

func TestAppend(t *testing.T) {
	initialActors := Actors{}
	newActor := &Actor{FirstName: "John"}

	initialActors.Append(newActor)

	if len(initialActors) != 1 {
		t.Errorf("Expected length of 1, but got %d", len(initialActors))
	}

	if initialActors[0].FirstName != "John" {
		t.Errorf("Expected actor name to be 'John Doe', but got '%s'", initialActors[0].FirstName)
	}

	secondActor := &Actor{FirstName: "Jane"}
	initialActors.Append(secondActor)

	if len(initialActors) != 2 {
		t.Errorf("Expected length of 2, but got %d", len(initialActors))
	}

	if initialActors[1].FirstName != "Jane" {
		t.Errorf("Expected second actor name to be 'Jane Doe', but got '%s'", initialActors[1].FirstName)
	}
}
