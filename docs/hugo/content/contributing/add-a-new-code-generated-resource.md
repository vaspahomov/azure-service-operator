# Adding a new code generated resource to ASO v2

This document discusses how to add a new resource to the ASO v2 code generation configuration. Check out [this PR](https://github.com/Azure/azure-service-operator/pull/1568) if you'd like to see what the end product looks like.

## What resources can be code generated?
Any ARM resource can be generated. There are a few ways to determine if a resource is an ARM resource:
1. If the resource is defined in the [ARM template JSON schema](https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json), or the [auto-generated ARM template JSON schema](https://schema.management.azure.com/schemas/common/autogeneratedResources.json) it is an ARM resource.
2. If the resource is defined in a `resource-manager` folder in the [Azure REST API specs](https://github.com/Azure/azure-rest-api-specs/tree/main/specification) repo, it is an ARM resource.

If the resource is not in either of the above places and cannot be deployed via an ARM template, then it's not an ARM resource and currently cannot be code generated.

## Determine a resource to add
There are three key pieces of information required before adding a resource to the code generation configuration file, and each of them can be found in the 
[ARM template JSON schema](https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json) or the [auto-generated ARM template JSON schema](https://schema.management.azure.com/schemas/common/autogeneratedResources.json). To get started, find the entry for the resource in one of the two templates mentioned. For example, the entry for Azure Network Security Groups looks like this:
`{ "$ref": "https://schema.management.azure.com/schemas/2020-11-01/Microsoft.Network.json#/resourceDefinitions/networkSecurityGroups" },`

**Note**: In many cases there will be multiple entries for the same resource, each with a different `api-version`. It is _strongly_ recommended that you use the latest available non-preview `api-version`.

1. The `name` of the resource.

   You usually know this going in. In our example above, the name of the resource is `networkSecurityGroups`.
2. The `group` the resource is in. 

   This is usually named after the Azure service, for example `resources` or `documentdb`. In our example entry from above, this is `network`.
3. The `api-version` of the resource.

   This is usually a date, sometimes with a `-preview` suffix. In our example entry from above, this is `2020-11-01`.

## Adding the resource to the code generation configuration file
The code generation configuration file is located [here](https://github.com/Azure/azure-service-operator/blob/main/v2/azure-arm.yaml). To add a new resource to this file, find the `exportFilters` section of the file and scroll down until you get to a block of `exportFilters` for individual resources. 

Add a new `exportFilter` of kind `include-transitive` at the end of that block, right _above_ this section:
```yaml
  # Exclude everything else as we are operating on an opt-in basis at the moment:
  - action: exclude
    because: We don't want to generate anything else, at the moment.
```

Your export filter should look like this:
```yaml
- action: include-transitive
  group: <group>
  version: v*api<api-version with dashes removed>
  name: <resource name, typically just remove the trailing "s">
  because: "including <resource>"
```

In the case of our example above, that ends up being:
```yaml
- action: include-transitive
  group: network
  version: v*api20201101
  name: NetworkSecurityGroup
  because: "including NSG"
```

## Run the code generator

Follow the steps in the [contributing guide](../contributing/) to set up your development environment.
Once you have a working development environment, run the `task` command to run the code generator.

## Fix any errors raised by the code generator

### \<Resource\> looks like a resource reference but was not labelled as one
Example:
>  Replace cross-resource references with genruntime.ResourceReference: 
> ["github.com/Azure/azure-service-operator/hack/generated/_apis/containerservice/v1alpha1api20210501/PrivateLinkResource.Id" looks like a resource reference but was not labelled as one. 
> It might need to be manually added to `newKnownReferencesMap`,

To fix this error, determine whether the property in question is an ARM ID or not, and then update the `newKnownReferencesMap` function 
in [add_cross_resource_references.go](https://github.com/Azure/azure-service-operator/blob/main/v2/tools/generator/internal/codegen/pipeline/add_cross_resource_references.go#:~:text=func-,newknownreferencesmap,-).

If the property is an ARM ID, update `newKnownReferencesMap` to flag that property as a reference:
```go
{
       typeName: astmodel.MakeTypeName(configuration.MakeLocalPackageReference("containerservice", "v1alpha1api20210501"), "PrivateLinkResource"),
       propName: "Id", 
}: true,
```

If the property is not an ARM ID, update `newKnownReferencesMap` to indicate that property is not a reference by providing the value **false** instead:
```go
{
       typeName: astmodel.MakeTypeName(configuration.MakeLocalPackageReference("containerservice", "v1alpha1api20210501"), "PrivateLinkResource"),
       propName: "Id", 
}: false,
``` 

TODO: expand on other common errors

## Examine the generated resource
After running the generator, the new resource you added should be in the [apis](https://github.com/Azure/azure-service-operator/blob/main/v2/api/) directory. 

Have a look through the files in the directory named after the `group` and `version` of the resource that was added.
In our `NetworkSecurityGroups` example, the best place to start is `/v2/api/network/v1alpha1api20201101/network_security_group_types_gen.go`
There may be other resources that already exist in that same directory - that's expected if ASO already supported some resources from that provider and API version.

Starting with the `network_security_group_types_gen.go` file, find the struct representing the resource you just added. It should be near the top and look something like this:
```go
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].reason"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].message"
//Generated from: https://schema.management.azure.com/schemas/2020-11-01/Microsoft.Network.json#/resourceDefinitions/networkSecurityGroups
type NetworkSecurityGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              NetworkSecurityGroups_Spec                                           `json:"spec,omitempty"`
	Status            NetworkSecurityGroup_Status_NetworkSecurityGroup_SubResourceEmbedded `json:"status,omitempty"`
}
```

Look over the `Spec` and `Status` types and their properties (and the properties of their properties and so-on).
The Azure REST API specs and ARM template JSON schemas which these types were derived from are not perfect. Sometimes they mark a `readonly` property as mutable or have another error or mistake. 

If you do identify properties which should be removed or changed, you can make customizations to the resource in the `typeTransformers` section of the code generation config. The most common issues have their own sections:

1. `# Deal with readonly properties that were not properly pruned in the JSON schema`
2. `# Deal with properties that should have been marked readOnly but weren't`

## Write a CRUD test for the resource
The best way to do this is to start from an [existing test](https://github.com/Azure/azure-service-operator/blob/main/v2/internal/controller/controllers/crd_cosmosdb_databaseaccount_test.go) and modify it to work for your resource. It can also be helpful to refer to examples in the [ARM templates GitHub repo](https://github.com/Azure/azure-quickstart-templates).

## Run the CRUD test for the resource and commit the recording
See [the code generator README](../contributing/#running-integration-tests) for how to run recording tests.

## Add a new sample
The samples are located in the [samples directory](https://github.com/Azure/azure-service-operator/blob/main/v2/config/samples). There should be at least one sample for each kind of supported resource. These currently need to be added manually. It's possible in the future we will automatically generate samples similar to how we automatically generate CRDs and types, but that doesn't happen today.

## Send a PR
You're all done!
