// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/artback/mvp/pkg/products (interfaces: Repository)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	products "github.com/artback/mvp/pkg/products"
	gomock "github.com/golang/mock/gomock"
)

// ProductRepository is a mock of Repository interface.
type ProductRepository struct {
	ctrl     *gomock.Controller
	recorder *ProductRepositoryMockRecorder
}

// ProductRepositoryMockRecorder is the mock recorder for ProductRepository.
type ProductRepositoryMockRecorder struct {
	mock *ProductRepository
}

// NewProductRepository creates a new mock instance.
func NewProductRepository(ctrl *gomock.Controller) *ProductRepository {
	mock := &ProductRepository{ctrl: ctrl}
	mock.recorder = &ProductRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *ProductRepository) EXPECT() *ProductRepositoryMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *ProductRepository) Delete(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *ProductRepositoryMockRecorder) Delete(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*ProductRepository)(nil).Delete), arg0, arg1, arg2)
}

// Get mocks base method.
func (m *ProductRepository) Get(arg0 context.Context, arg1 string) (*products.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*products.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *ProductRepositoryMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*ProductRepository)(nil).Get), arg0, arg1)
}

// Insert mocks base method.
func (m *ProductRepository) Insert(arg0 context.Context, arg1 products.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *ProductRepositoryMockRecorder) Insert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*ProductRepository)(nil).Insert), arg0, arg1)
}

// Update mocks base method.
func (m *ProductRepository) Update(arg0 context.Context, arg1 products.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *ProductRepositoryMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*ProductRepository)(nil).Update), arg0, arg1)
}
