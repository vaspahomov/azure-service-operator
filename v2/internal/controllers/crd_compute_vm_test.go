/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package controllers_test

import (
	"testing"

	"github.com/Azure/go-autorest/autorest/to"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	compute "github.com/Azure/azure-service-operator/v2/api/compute/v1alpha1api20201201"
	network "github.com/Azure/azure-service-operator/v2/api/network/v1alpha1api20201101"
	resources "github.com/Azure/azure-service-operator/v2/api/resources/v1alpha1api20200601"
	"github.com/Azure/azure-service-operator/v2/internal/testcommon"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

func newVM(
	tc *testcommon.KubePerTestContext,
	rg *resources.ResourceGroup,
	networkInterface *network.NetworkInterface) *compute.VirtualMachine {

	password := tc.Namer.GeneratePasswordOfLength(40)

	passwordKey := "password"
	secret := &v1.Secret{
		ObjectMeta: tc.MakeObjectMeta("vmsecret"),
		StringData: map[string]string{
			passwordKey: password,
		},
	}

	tc.CreateResource(secret)

	secretRef := genruntime.SecretReference{
		Name: secret.Name,
		Key:  passwordKey,
	}
	adminUsername := "bloom"
	size := compute.HardwareProfileVmSizeStandardA1V2

	return &compute.VirtualMachine{
		ObjectMeta: tc.MakeObjectMeta("vm"),
		Spec: compute.VirtualMachines_Spec{
			Location: tc.AzureRegion,
			Owner:    testcommon.AsOwner(rg),
			HardwareProfile: &compute.HardwareProfile{
				VmSize: &size,
			},
			OsProfile: &compute.OSProfile{
				AdminUsername: &adminUsername,
				// Specifying AdminPassword here rather than SSH Key to ensure that handling and injection
				// of secrets works.
				AdminPassword: &secretRef,
				ComputerName:  to.StringPtr("poppy"),
			},
			StorageProfile: &compute.StorageProfile{
				ImageReference: &compute.ImageReference{
					Offer:     to.StringPtr("UbuntuServer"),
					Publisher: to.StringPtr("Canonical"),
					Sku:       to.StringPtr("18.04-LTS"),
					Version:   to.StringPtr("latest"),
				},
			},
			NetworkProfile: &compute.VirtualMachines_Spec_Properties_NetworkProfile{
				NetworkInterfaces: []compute.VirtualMachines_Spec_Properties_NetworkProfile_NetworkInterfaces{{
					Reference: tc.MakeReferencePtrFromResource(networkInterface),
				}},
			},
		},
	}
}

func newVMNetworkInterface(tc *testcommon.KubePerTestContext, owner genruntime.KnownResourceReference, subnet *network.VirtualNetworksSubnet) *network.NetworkInterface {
	dynamic := network.NetworkInterfaceIPConfigurationPropertiesFormatPrivateIPAllocationMethodDynamic
	return &network.NetworkInterface{
		ObjectMeta: tc.MakeObjectMeta("nic"),
		Spec: network.NetworkInterfaces_Spec{
			Owner:    owner,
			Location: tc.AzureRegion,
			IpConfigurations: []network.NetworkInterfaces_Spec_Properties_IpConfigurations{{
				Name:                      "ipconfig1",
				PrivateIPAllocationMethod: &dynamic,
				Subnet: &network.SubResource{
					Reference: tc.MakeReferenceFromResource(subnet),
				},
			}},
		},
	}
}

func Test_Compute_VM_CRUD(t *testing.T) {
	t.Parallel()

	tc := globalTestContext.ForTest(t)
	rg := tc.CreateTestResourceGroupAndWait()

	vnet := newVMVirtualNetwork(tc, testcommon.AsOwner(rg))
	subnet := newVMSubnet(tc, testcommon.AsOwner(vnet))
	networkInterface := newVMNetworkInterface(tc, testcommon.AsOwner(rg), subnet)
	// Inefficient but avoids triggering the vnet/subnets problem.
	// https://github.com/Azure/azure-service-operator/issues/1944
	tc.CreateResourceAndWait(vnet)
	tc.CreateResourcesAndWait(subnet, networkInterface)
	vm := newVM(tc, rg, networkInterface)

	tc.CreateResourceAndWait(vm)
	tc.Expect(vm.Status.Id).ToNot(BeNil())
	armId := *vm.Status.Id

	// Perform a simple patch to turn on boot diagnostics
	old := vm.DeepCopy()
	vm.Spec.DiagnosticsProfile = &compute.DiagnosticsProfile{
		BootDiagnostics: &compute.BootDiagnostics{
			Enabled: to.BoolPtr(true),
		},
	}

	tc.Patch(old, vm)

	objectKey := client.ObjectKeyFromObject(vm)

	// Ensure state eventually gets updated in k8s from change in Azure.
	tc.Eventually(func() bool {
		var updated compute.VirtualMachine
		tc.GetResource(objectKey, &updated)

		diagProfile := updated.Status.DiagnosticsProfile
		if diagProfile == nil {
			return false
		}

		if diagProfile.BootDiagnostics == nil {
			return false
		}

		return *diagProfile.BootDiagnostics.Enabled
	}).Should(BeTrue())

	// Delete VM and resources.
	tc.DeleteResourcesAndWait(vm, networkInterface, subnet, vnet, rg)

	// Ensure that the resource was really deleted in Azure
	exists, retryAfter, err := tc.AzureClient.HeadByID(tc.Ctx, armId, string(compute.VirtualMachinesSpecAPIVersion20201201))
	tc.Expect(err).ToNot(HaveOccurred())
	tc.Expect(retryAfter).To(BeZero())
	tc.Expect(exists).To(BeFalse())
}
