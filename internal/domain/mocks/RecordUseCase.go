// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	dns "github.com/miekg/dns"

	domain "github.com/cewuandy/go-restful-dns/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// RecordUseCase is an autogenerated mock type for the RecordUseCase type
type RecordUseCase struct {
	mock.Mock
}

// CreateRecord provides a mock function with given fields: ctx, rr
func (_m *RecordUseCase) CreateRecord(ctx context.Context, rr dns.RR) error {
	ret := _m.Called(ctx, rr)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, dns.RR) error); ok {
		r0 = rf(ctx, rr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteRecord provides a mock function with given fields: ctx, question
func (_m *RecordUseCase) DeleteRecord(ctx context.Context, question domain.Question) error {
	ret := _m.Called(ctx, question)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Question) error); ok {
		r0 = rf(ctx, question)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetRecord provides a mock function with given fields: ctx, question
func (_m *RecordUseCase) GetRecord(ctx context.Context, question domain.Question) (dns.RR, error) {
	ret := _m.Called(ctx, question)

	var r0 dns.RR
	if rf, ok := ret.Get(0).(func(context.Context, domain.Question) dns.RR); ok {
		r0 = rf(ctx, question)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(dns.RR)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.Question) error); ok {
		r1 = rf(ctx, question)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListRecords provides a mock function with given fields: ctx
func (_m *RecordUseCase) ListRecords(ctx context.Context) ([]dns.RR, error) {
	ret := _m.Called(ctx)

	var r0 []dns.RR
	if rf, ok := ret.Get(0).(func(context.Context) []dns.RR); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dns.RR)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateRecord provides a mock function with given fields: ctx, rr
func (_m *RecordUseCase) UpdateRecord(ctx context.Context, rr dns.RR) error {
	ret := _m.Called(ctx, rr)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, dns.RR) error); ok {
		r0 = rf(ctx, rr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}