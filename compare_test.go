package testing_utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEqual(t *testing.T) {
	i := 123
	j := 321

	x1 := struct {
		SubStruct struct{
			Slice []int
			SliceOfStructs []*struct{
				StringVal string
			}
			Bool bool
		}
		SliceOfStructs []struct{
			IntVal *int
		}
		unexportedField bool
	}{
		SubStruct: struct{
			Slice []int
			SliceOfStructs []*struct{
				StringVal string
			}
			Bool bool
		}{
			Slice: []int{1,2,3},
			SliceOfStructs: []*struct{ StringVal string }{
				{
					"123",
				},{
					"321",
				},
			},
			Bool: true,
		},
		SliceOfStructs: []struct{ IntVal *int }{
			{IntVal: &i},
			{IntVal: &j},
		},
		unexportedField: true,
	}

	x2 := struct {
		SubStruct struct{
			Slice []int
			SliceOfStructs []*struct{
				StringVal string
			}
			Bool bool
		}
		SliceOfStructs []struct{
			IntVal *int
		}
	}{
		SubStruct: struct{
			Slice []int
			SliceOfStructs []*struct{
				StringVal string
			}
			Bool bool
		}{
			Slice: []int{1,3,2},
			SliceOfStructs: []*struct{ StringVal string }{
				{
					"321",
				},{
					"123",
				},
			},
			Bool: true,
		},
		SliceOfStructs: []struct{ IntVal *int }{
			{IntVal: &i},
			{IntVal: &j},
		},
	}

	x3 := struct {
		SubStruct struct{
			Slice []int
			SliceOfStructs []struct{
				StringVal string
			}
			Bool bool
		}
		SliceOfStructs []struct{
			IntVal int
		}
	}{
		SubStruct: struct{
			Slice []int
			SliceOfStructs []struct{
				StringVal string
			}
			Bool bool
		}{
			Slice: []int{1,3,2},
			SliceOfStructs: []struct{ StringVal string }{
				{
					"321",
				},{
					"123",
				},
			},
			Bool: true,
		},
		SliceOfStructs: []struct{ IntVal int }{
			{IntVal: i},
			{IntVal: j},
		},
	}

	ok, desc := Equal(x1, x2)
	assert.True(t, ok, desc)

	ok, desc = Equal(x2, x1)
	assert.True(t, ok, desc)

	ok, desc = Equal(x2, x3)
	assert.False(t, ok, desc)

	ok, desc = Equal(x1, x3)
	assert.False(t, ok, desc)
}
