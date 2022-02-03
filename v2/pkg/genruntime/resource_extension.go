/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

/*
This file contains manual implementations to reduce code bloat in generated code.
*/

// ResourceExtension is defines extended functionality of a resource used by the reconciler
type ResourceExtension interface {
	// GetExtendedResources Returns the KubernetesResource slice for Resource versions
	GetExtendedResources() []KubernetesResource
}
