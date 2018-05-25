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

package ipamd

import (
	"net"
	"testing"
	"time"

	"github.com/aws/amazon-vpc-cni-k8s/ipamd/datastore"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/awsutils"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/awsutils/mocks"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/k8sapi"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/k8sapi/mocks"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/networkutils/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	primaryENIid     = "eni-00000000"
	secENIid         = "eni-00000001"
	testAttachmentID = "eni-00000000-attach"
	eniID            = "eni-5731da78"
	primaryMAC       = "12:ef:2a:98:e5:5a"
	secMAC           = "12:ef:2a:98:e5:5b"
	primaryDevice    = 2
	secDevice        = 0
	primarySubnet    = "10.10.10.0/24"
	secSubnet        = "10.10.20.0/24"
	ipaddr01         = "10.10.10.11"
	ipaddr02         = "10.10.10.12"
	ipaddr03         = "10.10.10.13"
	ipaddr11         = "10.10.20.11"
	ipaddr12         = "10.10.20.12"
	ipaddr13         = "10.10.20.13"
	vpcCIDR          = "10.10.0.0/16"
)

func setup(t *testing.T) (*gomock.Controller,
	*mock_awsutils.MockAPIs,
	*mock_k8sapi.MockK8SAPIs,
	*mock_networkutils.MockNetworkAPIs) {
	ctrl := gomock.NewController(t)
	return ctrl,
		mock_awsutils.NewMockAPIs(ctrl),
		mock_k8sapi.NewMockK8SAPIs(ctrl),
		mock_networkutils.NewMockNetworkAPIs(ctrl)
}

func TestNodeInit(t *testing.T) {
	ctrl, mockAWS, mockK8S, mockNetwork := setup(t)
	defer ctrl.Finish()

	mockContext := &IPAMContext{
		awsClient:     mockAWS,
		k8sClient:     mockK8S,
		networkClient: mockNetwork}

	eni1 := awsutils.ENIMetadata{
		ENIID:          primaryENIid,
		MAC:            primaryMAC,
		DeviceNumber:   primaryDevice,
		SubnetIPv4CIDR: primarySubnet,
		LocalIPv4s:     []string{ipaddr01, ipaddr02},
	}

	eni2 := awsutils.ENIMetadata{
		ENIID:          secENIid,
		MAC:            secMAC,
		DeviceNumber:   secDevice,
		SubnetIPv4CIDR: secSubnet,
		LocalIPv4s:     []string{ipaddr11, ipaddr12},
	}
	mockAWS.EXPECT().GetENILimit().Return(4, nil)
	mockAWS.EXPECT().GetENIipLimit().Return(int64(56), nil)
	mockAWS.EXPECT().GetAttachedENIs().Return([]awsutils.ENIMetadata{eni1, eni2}, nil)
	mockAWS.EXPECT().GetVPCIPv4CIDR().Return(vpcCIDR)
	mockAWS.EXPECT().GetLocalIPv4().Return(ipaddr01)

	_, vpcCIDR, _ := net.ParseCIDR(vpcCIDR)
	primaryIP := net.ParseIP(ipaddr01)
	mockNetwork.EXPECT().SetupHostNetwork(vpcCIDR, &primaryIP).Return(nil)

	//primaryENIid
	mockAWS.EXPECT().GetPrimaryENI().Return(primaryENIid)
	mockAWS.EXPECT().AllocAllIPAddress(primaryENIid).Return(nil)
	attachmentID := testAttachmentID
	testAddr1 := ipaddr01
	testAddr2 := ipaddr02
	primary := true
	eniResp := []*ec2.NetworkInterfacePrivateIpAddress{
		&ec2.NetworkInterfacePrivateIpAddress{
			PrivateIpAddress: &testAddr1, Primary: &primary},
		&ec2.NetworkInterfacePrivateIpAddress{
			PrivateIpAddress: &testAddr2, Primary: &primary}}
	mockAWS.EXPECT().GetPrimaryENI().Return(primaryENIid)
	mockAWS.EXPECT().DescribeENI(primaryENIid).Return(eniResp, &attachmentID, nil)

	//secENIid
	mockAWS.EXPECT().GetPrimaryENI().Return(primaryENIid)
	mockAWS.EXPECT().AllocAllIPAddress(secENIid).Return(nil)
	attachmentID = testAttachmentID
	testAddr11 := ipaddr11
	testAddr12 := ipaddr12
	primary = false
	eniResp = []*ec2.NetworkInterfacePrivateIpAddress{
		&ec2.NetworkInterfacePrivateIpAddress{
			PrivateIpAddress: &testAddr11, Primary: &primary},
		&ec2.NetworkInterfacePrivateIpAddress{
			PrivateIpAddress: &testAddr12, Primary: &primary}}
	mockAWS.EXPECT().GetPrimaryENI().Return(primaryENIid)
	mockAWS.EXPECT().DescribeENI(secENIid).Return(eniResp, &attachmentID, nil)
	mockNetwork.EXPECT().SetupENINetwork(gomock.Any(), secMAC, secDevice, secSubnet)

	mockAWS.EXPECT().GetLocalIPv4().Return(ipaddr01)
	mockK8S.EXPECT().K8SGetLocalPodIPs(gomock.Any()).Return([]*k8sapi.K8SPodInfo{&k8sapi.K8SPodInfo{Name: "pod1",
		Namespace: "default"}}, nil)

	err := mockContext.nodeInit()
	assert.NoError(t, err)
}

func TestIncreaseIPPool(t *testing.T) {
	ctrl, mockAWS, mockK8S, mockNetwork := setup(t)
	defer ctrl.Finish()

	mockContext := &IPAMContext{
		awsClient:     mockAWS,
		k8sClient:     mockK8S,
		networkClient: mockNetwork,
		primaryIP:     make(map[string]string),
	}

	mockContext.dataStore = datastore.NewDataStore()

	eni2 := secENIid

	mockAWS.EXPECT().GetENILimit().Return(4, nil)
	mockAWS.EXPECT().AllocENI().Return(eni2, nil)

	mockAWS.EXPECT().AllocAllIPAddress(eni2)

	mockAWS.EXPECT().GetAttachedENIs().Return([]awsutils.ENIMetadata{
		awsutils.ENIMetadata{
			ENIID:          primaryENIid,
			MAC:            primaryMAC,
			DeviceNumber:   primaryDevice,
			SubnetIPv4CIDR: primarySubnet,
			LocalIPv4s:     []string{ipaddr01, ipaddr02},
		},
		awsutils.ENIMetadata{
			ENIID:          secENIid,
			MAC:            secMAC,
			DeviceNumber:   secDevice,
			SubnetIPv4CIDR: secSubnet,
			LocalIPv4s:     []string{ipaddr11, ipaddr12}},
	}, nil)

	mockAWS.EXPECT().GetPrimaryENI().Return(primaryENIid)

	primary := false
	attachmentID := testAttachmentID
	testAddr11 := ipaddr11
	testAddr12 := ipaddr12

	mockAWS.EXPECT().DescribeENI(eni2).Return(
		[]*ec2.NetworkInterfacePrivateIpAddress{
			&ec2.NetworkInterfacePrivateIpAddress{
				PrivateIpAddress: &testAddr11, Primary: &primary},
			&ec2.NetworkInterfacePrivateIpAddress{
				PrivateIpAddress: &testAddr12, Primary: &primary}}, &attachmentID, nil)

	mockAWS.EXPECT().GetPrimaryENI().Return(primaryENIid)
	mockNetwork.EXPECT().SetupENINetwork(gomock.Any(), secMAC, secDevice, secSubnet)

	mockContext.increaseIPPool()
}

func TestDecreaseIPPool(t *testing.T) {
	ctrl, mockAWS, mockK8S, mockNetwork := setup(t)
	defer ctrl.Finish()

	mockContext := &IPAMContext{
		awsClient:     mockAWS,
		k8sClient:     mockK8S,
		networkClient: mockNetwork,
		primaryIP:     make(map[string]string),
	}

	ds := datastore.NewDataStore()
	ds.AddENI(secENIid, 1, false)
	// Have to sleep for a minute in order for it to qualify as a deletable ENI because this interface
	// is not exposed and not particularly testable...
	time.Sleep(1 * time.Minute)
	mockContext.dataStore = ds

	mockAWS.EXPECT().FreeENI(secENIid)

	mockContext.decreaseIPPool()
}

// Testing the mutex on the context to make sure that when a delete is called, a lock is
// acquired on the context, so no add condition can cause a race
func TestDecreaseIPPoolRaceCondition(t *testing.T) {
	ctrl, mockAWS, mockK8S, mockNetwork := setup(t)
	defer ctrl.Finish()

	mockContext := &IPAMContext{
		awsClient:     mockAWS,
		k8sClient:     mockK8S,
		networkClient: mockNetwork,
		primaryIP:     make(map[string]string),
	}

	ds := datastore.NewDataStore()
	ds.AddENI(secENIid, 1, false)
	// Have to sleep for a minute in order for it to qualify as a deletable ENI because this interface
	// is not exposed and not particularly testable...
	time.Sleep(1 * time.Minute)
	mockContext.dataStore = ds

	// Ensures that the decrease and increase calls are called in the correct order
	mockAWS.EXPECT().FreeENI(secENIid)
	mockAWS.EXPECT().GetENILimit().Return(4, nil)
	mockAWS.EXPECT().AllocENI().Return(secENIid, nil)
	mockAWS.EXPECT().AllocAllIPAddress(secENIid)
	mockAWS.EXPECT().GetAttachedENIs().Return([]awsutils.ENIMetadata{
		awsutils.ENIMetadata{
			ENIID:          primaryENIid,
			MAC:            primaryMAC,
			DeviceNumber:   primaryDevice,
			SubnetIPv4CIDR: primarySubnet,
			LocalIPv4s:     []string{ipaddr01, ipaddr02},
		},
		awsutils.ENIMetadata{
			ENIID:          secENIid,
			MAC:            secMAC,
			DeviceNumber:   secDevice,
			SubnetIPv4CIDR: secSubnet,
			LocalIPv4s:     []string{ipaddr11, ipaddr12}},
	}, nil)
	mockAWS.EXPECT().GetPrimaryENI().Return(primaryENIid)
	mockAWS.EXPECT().DescribeENI(secENIid).Return(
		[]*ec2.NetworkInterfacePrivateIpAddress{
			&ec2.NetworkInterfacePrivateIpAddress{
				PrivateIpAddress: aws.String(ipaddr11), Primary: aws.Bool(false),
			},
			&ec2.NetworkInterfacePrivateIpAddress{
				PrivateIpAddress: aws.String(ipaddr12), Primary: aws.Bool(false),
			},
		},
		aws.String(testAttachmentID),
		nil,
	)
	mockAWS.EXPECT().GetPrimaryENI().Return(primaryENIid)
	mockNetwork.EXPECT().SetupENINetwork(gomock.Any(), secMAC, secDevice, secSubnet)

	mockContext.decreaseIPPool()
	go func() {
		time.Sleep(50 * time.Millisecond)
		mockContext.increaseIPPool()
	}()
}

func TestNodeIPPoolReconcile(t *testing.T) {
	ctrl, mockAWS, mockK8S, mockNetwork := setup(t)
	defer ctrl.Finish()

	mockContext := &IPAMContext{
		awsClient:     mockAWS,
		k8sClient:     mockK8S,
		networkClient: mockNetwork,
		primaryIP:     make(map[string]string),
	}

	mockContext.dataStore = datastore.NewDataStore()

	mockAWS.EXPECT().GetAttachedENIs().Return([]awsutils.ENIMetadata{
		awsutils.ENIMetadata{
			ENIID:          primaryENIid,
			MAC:            primaryMAC,
			DeviceNumber:   primaryDevice,
			SubnetIPv4CIDR: primarySubnet,
			LocalIPv4s:     []string{ipaddr01, ipaddr02},
		},
	}, nil)

	mockAWS.EXPECT().GetPrimaryENI().Return(primaryENIid)

	primary := true
	notPrimary := false
	attachmentID := testAttachmentID
	testAddr1 := ipaddr01
	testAddr2 := ipaddr02

	mockAWS.EXPECT().DescribeENI(primaryENIid).Return(
		[]*ec2.NetworkInterfacePrivateIpAddress{
			&ec2.NetworkInterfacePrivateIpAddress{
				PrivateIpAddress: &testAddr1, Primary: &primary},
			&ec2.NetworkInterfacePrivateIpAddress{
				PrivateIpAddress: &testAddr2, Primary: &notPrimary}}, &attachmentID, nil)
	mockAWS.EXPECT().GetPrimaryENI().Return(primaryENIid)

	mockContext.nodeIPPoolReconcile(0)

	curENIs := mockContext.dataStore.GetENIInfos()
	assert.Equal(t, len(curENIs.ENIIPPools), 1)
	assert.Equal(t, curENIs.TotalIPs, 1)

	// remove 1 IP
	mockAWS.EXPECT().GetAttachedENIs().Return([]awsutils.ENIMetadata{
		awsutils.ENIMetadata{
			ENIID:          primaryENIid,
			MAC:            primaryMAC,
			DeviceNumber:   primaryDevice,
			SubnetIPv4CIDR: primarySubnet,
			LocalIPv4s:     []string{ipaddr01},
		},
	}, nil)

	mockContext.nodeIPPoolReconcile(0)
	curENIs = mockContext.dataStore.GetENIInfos()
	assert.Equal(t, len(curENIs.ENIIPPools), 1)
	assert.Equal(t, curENIs.TotalIPs, 0)

	// remove eni
	mockAWS.EXPECT().GetAttachedENIs().Return(nil, nil)

	mockContext.nodeIPPoolReconcile(0)
	curENIs = mockContext.dataStore.GetENIInfos()
	assert.Equal(t, len(curENIs.ENIIPPools), 0)
	assert.Equal(t, curENIs.TotalIPs, 0)
}
