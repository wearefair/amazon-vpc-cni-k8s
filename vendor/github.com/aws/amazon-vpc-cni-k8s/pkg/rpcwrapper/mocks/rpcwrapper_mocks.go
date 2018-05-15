// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/aws/amazon-vpc-cni-k8s/pkg/rpcwrapper (interfaces: RPC)

// Package mock_rpcwrapper is a generated GoMock package.
package mock_rpcwrapper

import (
	reflect "reflect"

	rpc "github.com/aws/amazon-vpc-cni-k8s/rpc"
	gomock "github.com/golang/mock/gomock"
	"google.golang.org/grpc"
)

// MockRPC is a mock of RPC interface
type MockRPC struct {
	ctrl     *gomock.Controller
	recorder *MockRPCMockRecorder
}

// MockRPCMockRecorder is the mock recorder for MockRPC
type MockRPCMockRecorder struct {
	mock *MockRPC
}

// NewMockRPC creates a new mock instance
func NewMockRPC(ctrl *gomock.Controller) *MockRPC {
	mock := &MockRPC{ctrl: ctrl}
	mock.recorder = &MockRPCMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRPC) EXPECT() *MockRPCMockRecorder {
	return m.recorder
}

// NewCNIBackendClient mocks base method
func (m *MockRPC) NewCNIBackendClient(arg0 *grpc.ClientConn) rpc.CNIBackendClient {
	ret := m.ctrl.Call(m, "NewCNIBackendClient", arg0)
	ret0, _ := ret[0].(rpc.CNIBackendClient)
	return ret0
}

// NewCNIBackendClient indicates an expected call of NewCNIBackendClient
func (mr *MockRPCMockRecorder) NewCNIBackendClient(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewCNIBackendClient", reflect.TypeOf((*MockRPC)(nil).NewCNIBackendClient), arg0)
}
