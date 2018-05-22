// Copyright 2015-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/aws/amazon-ecs-agent/agent/async (interfaces: Cache)

package mock_async

import (
	async "github.com/aws/amazon-ecs-agent/agent/async"
	gomock "github.com/golang/mock/gomock"
)

// Mock of Cache interface
type MockCache struct {
	ctrl     *gomock.Controller
	recorder *_MockCacheRecorder
}

// Recorder for MockCache (not exported)
type _MockCacheRecorder struct {
	mock *MockCache
}

func NewMockCache(ctrl *gomock.Controller) *MockCache {
	mock := &MockCache{ctrl: ctrl}
	mock.recorder = &_MockCacheRecorder{mock}
	return mock
}

func (_m *MockCache) EXPECT() *_MockCacheRecorder {
	return _m.recorder
}

func (_m *MockCache) Get(_param0 string) (async.Value, bool) {
	ret := _m.ctrl.Call(_m, "Get", _param0)
	ret0, _ := ret[0].(async.Value)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

func (_mr *_MockCacheRecorder) Get(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Get", arg0)
}

func (_m *MockCache) Set(_param0 string, _param1 async.Value) {
	_m.ctrl.Call(_m, "Set", _param0, _param1)
}

func (_mr *_MockCacheRecorder) Set(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Set", arg0, arg1)
}