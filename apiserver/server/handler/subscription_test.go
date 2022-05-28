package handler

import (
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func Test_uniqueAppend(t *testing.T) {
	type args struct {
		list      []v1.LocalObjectReference
		reference v1.LocalObjectReference
	}
	tests := []struct {
		name string
		args args
		want []v1.LocalObjectReference
	}{{
		name: "no duplicated item",
		args: args{
			list: []v1.LocalObjectReference{},
			reference: v1.LocalObjectReference{
				Name: "fake",
			},
		},
		want: []v1.LocalObjectReference{{Name: "fake"}},
	}, {
		name: "have duplicated item",
		args: args{
			list: []v1.LocalObjectReference{{Name: "fake"}},
			reference: v1.LocalObjectReference{
				Name: "fake",
			},
		},
		want: []v1.LocalObjectReference{{Name: "fake"}},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, uniqueAppend(tt.args.list, tt.args.reference), "uniqueAppend(%v, %v)", tt.args.list, tt.args.reference)
		})
	}
}
