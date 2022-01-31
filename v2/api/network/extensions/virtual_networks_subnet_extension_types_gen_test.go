// Code generated by azure-service-operator-codegen. DO NOT EDIT.
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package extensions

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/kr/pretty"
	"github.com/kylelemons/godebug/diff"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"os"
	"reflect"
	"testing"
)

func Test_VirtualNetworksSubnetExtension_WhenSerializedToJson_DeserializesAsEqual(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	parameters.MaxSize = 10
	properties := gopter.NewProperties(parameters)
	properties.Property(
		"Round trip of VirtualNetworksSubnetExtension via JSON returns original",
		prop.ForAll(RunJSONSerializationTestForVirtualNetworksSubnetExtension, VirtualNetworksSubnetExtensionGenerator()))
	properties.TestingRun(t, gopter.NewFormatedReporter(true, 240, os.Stdout))
}

// RunJSONSerializationTestForVirtualNetworksSubnetExtension runs a test to see if a specific instance of VirtualNetworksSubnetExtension round trips to JSON and back losslessly
func RunJSONSerializationTestForVirtualNetworksSubnetExtension(subject VirtualNetworksSubnetExtension) string {
	// Serialize to JSON
	bin, err := json.Marshal(subject)
	if err != nil {
		return err.Error()
	}

	// Deserialize back into memory
	var actual VirtualNetworksSubnetExtension
	err = json.Unmarshal(bin, &actual)
	if err != nil {
		return err.Error()
	}

	// Check for outcome
	match := cmp.Equal(subject, actual, cmpopts.EquateEmpty())
	if !match {
		actualFmt := pretty.Sprint(actual)
		subjectFmt := pretty.Sprint(subject)
		result := diff.Diff(subjectFmt, actualFmt)
		return result
	}

	return ""
}

// Generator of VirtualNetworksSubnetExtension instances for property testing - lazily instantiated by
//VirtualNetworksSubnetExtensionGenerator()
var virtualNetworksSubnetExtensionGenerator gopter.Gen

// VirtualNetworksSubnetExtensionGenerator returns a generator of VirtualNetworksSubnetExtension instances for property testing.
func VirtualNetworksSubnetExtensionGenerator() gopter.Gen {
	if virtualNetworksSubnetExtensionGenerator != nil {
		return virtualNetworksSubnetExtensionGenerator
	}

	generators := make(map[string]gopter.Gen)
	virtualNetworksSubnetExtensionGenerator = gen.Struct(reflect.TypeOf(VirtualNetworksSubnetExtension{}), generators)

	return virtualNetworksSubnetExtensionGenerator
}
