// Code generated by azure-service-operator-codegen. DO NOT EDIT.
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package extensions

import (
	compute "github.com/Azure/azure-service-operator/v2/api/compute/v1alpha1api20201201"
	"github.com/Azure/azure-service-operator/v2/api/compute/v1alpha1api20201201storage"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

type VirtualMachineExtension struct {
}

// GetExtendedResources Returns the KubernetesResource slice for Resource versions
func (extension *VirtualMachineExtension) GetExtendedResources() []genruntime.KubernetesResource {
	return []genruntime.KubernetesResource{
		&compute.VirtualMachine{},
		&v1alpha1api20201201storage.VirtualMachine{}}
}
