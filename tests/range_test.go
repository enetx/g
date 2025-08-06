package g_test

import (
	"reflect"
	"testing"

	. "github.com/enetx/g"
)

func TestRange(t *testing.T) {
	tests := []struct {
		name        string
		start, stop int
		step        []int
		want        Slice[int]
	}{
		{
			name:  "default step",
			start: 0, stop: 3,
			step: nil,
			want: Slice[int]{0, 1, 2},
		},
		{
			name:  "positive step 1",
			start: 0, stop: 5,
			step: []int{1},
			want: Slice[int]{0, 1, 2, 3, 4},
		},
		{
			name:  "custom positive step",
			start: 2, stop: 10,
			step: []int{2},
			want: Slice[int]{2, 4, 6, 8},
		},
		{
			name:  "negative step",
			start: 5, stop: 0,
			step: []int{-1},
			want: Slice[int]{5, 4, 3, 2, 1},
		},
		{
			name:  "zero step yields nothing",
			start: 0, stop: 5,
			step: []int{0},
			want: Slice[int]{},
		},
		{
			name:  "empty range when start == stop",
			start: 3, stop: 3,
			step: nil,
			want: Slice[int]{},
		},
		{
			name:  "step never approaches stop",
			start: 0, stop: 5,
			step: []int{-1},
			want: Slice[int]{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Range(tc.start, tc.stop, tc.step...).Collect()
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Range(%d, %d, %v) = %v; want %v",
					tc.start, tc.stop, tc.step, got, tc.want)
			}
		})
	}
}

func TestRangeInclusive(t *testing.T) {
	tests := []struct {
		name        string
		start, stop int
		step        []int
		want        Slice[int]
	}{
		{
			name:  "default step",
			start: 0, stop: 3,
			step: nil,
			want: Slice[int]{0, 1, 2, 3},
		},
		{
			name:  "positive step 1",
			start: 0, stop: 5,
			step: []int{1},
			want: Slice[int]{0, 1, 2, 3, 4, 5},
		},
		{
			name:  "custom positive step",
			start: 2, stop: 10,
			step: []int{2},
			want: Slice[int]{2, 4, 6, 8, 10},
		},
		{
			name:  "negative step",
			start: 5, stop: 0,
			step: []int{-1},
			want: Slice[int]{5, 4, 3, 2, 1, 0},
		},
		{
			name:  "zero step yields nothing",
			start: 0, stop: 5,
			step: []int{0},
			want: Slice[int]{},
		},
		{
			name:  "inclusive when start == stop",
			start: 3, stop: 3,
			step: nil,
			want: Slice[int]{3},
		},
		{
			name:  "step never approaches stop",
			start: 0, stop: 5,
			step: []int{-1},
			want: Slice[int]{},
		},
		{
			name:  "step overshoots stop",
			start: 0, stop: 5,
			step: []int{6},
			want: Slice[int]{0},
		},
		{
			name:  "step exactly reaches stop",
			start: 0, stop: 6,
			step: []int{3},
			want: Slice[int]{0, 3, 6},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := RangeInclusive(tc.start, tc.stop, tc.step...).Collect()
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("RangeInclusive(%d, %d, %v) = %v; want %v",
					tc.start, tc.stop, tc.step, got, tc.want)
			}
		})
	}
}
