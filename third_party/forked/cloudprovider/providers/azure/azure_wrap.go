/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package azure

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-12-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/Azure/go-autorest/autorest"
	"k8s.io/apimachinery/pkg/types"
)

// checkExistsFromError inspects an error and returns a true if err is nil,
// false if error is an autorest.Error with StatusCode=404 and will return the
// error back if error is another status code or another type of error.
func checkResourceExistsFromError(err error) (bool, error) {
	if err == nil {
		return true, nil
	}
	v, ok := err.(autorest.DetailedError)
	if ok && v.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return false, v
}

func (az *Cloud) getVirtualMachine(ctx context.Context, nodeName types.NodeName) (vm compute.VirtualMachine, exists bool, err error) {
	var realErr error

	vmName := string(nodeName)
	vm, err = az.VirtualMachinesClient.Get(ctx, az.ResourceGroup, vmName, "")

	exists, realErr = checkResourceExistsFromError(err)
	if realErr != nil {
		return vm, false, realErr
	}

	if !exists {
		return vm, false, nil
	}

	return vm, exists, err
}

func (az *Cloud) getSecurityGroup(ctx context.Context) (sg network.SecurityGroup, exists bool, err error) {
	var realErr error

	sg, err = az.SecurityGroupsClient.Get(ctx, az.ResourceGroup, az.SecurityGroupName, "")

	exists, realErr = checkResourceExistsFromError(err)
	if realErr != nil {
		return sg, false, realErr
	}

	if !exists {
		return sg, false, nil
	}

	return sg, exists, err
}
