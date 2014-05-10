package core

import (
	"reflect"
	"testing"
)

func TestImageIDPartitions(t *testing.T) {
	ic := &ImageConfiguration{ID: "00ofrA"}
	partitions := ic.IDPartitions()
	expected := []string{"00", "of", "rA"}

	if !reflect.DeepEqual(expected, partitions) {
		t.Errorf("expected\n%v got\n%v", expected, partitions)
	}
}
