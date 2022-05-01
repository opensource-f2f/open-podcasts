package handler

import (
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func Test_removeLocalObjectReference(t *testing.T) {
	type args struct {
		list      []v1.LocalObjectReference
		reference v1.LocalObjectReference
	}
	tests := []struct {
		name        string
		args        args
		wantResult  []v1.LocalObjectReference
		wantRemoved bool
	}{{
		name: "not exist",
		args: args{
			list:      []v1.LocalObjectReference{{Name: "good"}},
			reference: v1.LocalObjectReference{Name: "fake"},
		},
		wantResult:  []v1.LocalObjectReference{{Name: "good"}},
		wantRemoved: false,
	}, {
		name: "remove first",
		args: args{
			list:      []v1.LocalObjectReference{{Name: "fake"}, {Name: "good"}},
			reference: v1.LocalObjectReference{Name: "fake"},
		},
		wantResult:  []v1.LocalObjectReference{{Name: "good"}},
		wantRemoved: true,
	}, {
		name: "remove last",
		args: args{
			list:      []v1.LocalObjectReference{{Name: "good"}, {Name: "fake"}},
			reference: v1.LocalObjectReference{Name: "fake"},
		},
		wantResult:  []v1.LocalObjectReference{{Name: "good"}},
		wantRemoved: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotRemoved := removeLocalObjectReference(tt.args.list, tt.args.reference)
			assert.Equalf(t, tt.wantResult, gotResult, "removeLocalObjectReference(%v, %v)", tt.args.list, tt.args.reference)
			assert.Equalf(t, tt.wantRemoved, gotRemoved, "removeLocalObjectReference(%v, %v)", tt.args.list, tt.args.reference)
		})
	}
}
