package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/Permify/permify/pkg/database"
	base "github.com/Permify/permify/pkg/pb/base/v1"
	"github.com/Permify/permify/pkg/token"
)

// RelationshipReader is an autogenerated mock type for the RelationshipReader type
type RelationshipReader struct {
	mock.Mock
}

// QueryRelationships -
func (_m *RelationshipReader) QueryRelationships(ctx context.Context, filter *base.TupleFilter, snap token.SnapToken) (database.ITupleCollection, error) {
	ret := _m.Called(filter, snap)

	var r0 *database.TupleCollection
	if rf, ok := ret.Get(0).(func(context.Context, *base.TupleFilter, token.SnapToken) *database.TupleCollection); ok {
		r0 = rf(ctx, filter, snap)
	} else {
		r0 = ret.Get(0).(*database.TupleCollection)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *base.TupleFilter, token.SnapToken) error); ok {
		r1 = rf(ctx, filter, snap)
	} else {
		if e, ok := ret.Get(1).(error); ok {
			r1 = e
		} else {
			r1 = nil
		}
	}

	return r0, r1
}