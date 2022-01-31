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

func Test_NetworkSecurityGroupExtension_WhenSerializedToJson_DeserializesAsEqual(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	parameters.MaxSize = 10
	properties := gopter.NewProperties(parameters)
	properties.Property(
		"Round trip of NetworkSecurityGroupExtension via JSON returns original",
		prop.ForAll(RunJSONSerializationTestForNetworkSecurityGroupExtension, NetworkSecurityGroupExtensionGenerator()))
	properties.TestingRun(t, gopter.NewFormatedReporter(true, 240, os.Stdout))
}

// RunJSONSerializationTestForNetworkSecurityGroupExtension runs a test to see if a specific instance of NetworkSecurityGroupExtension round trips to JSON and back losslessly
func RunJSONSerializationTestForNetworkSecurityGroupExtension(subject NetworkSecurityGroupExtension) string {
	// Serialize to JSON
	bin, err := json.Marshal(subject)
	if err != nil {
		return err.Error()
	}

	// Deserialize back into memory
	var actual NetworkSecurityGroupExtension
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

// Generator of NetworkSecurityGroupExtension instances for property testing - lazily instantiated by
//NetworkSecurityGroupExtensionGenerator()
var networkSecurityGroupExtensionGenerator gopter.Gen

// NetworkSecurityGroupExtensionGenerator returns a generator of NetworkSecurityGroupExtension instances for property testing.
func NetworkSecurityGroupExtensionGenerator() gopter.Gen {
	if networkSecurityGroupExtensionGenerator != nil {
		return networkSecurityGroupExtensionGenerator
	}

	generators := make(map[string]gopter.Gen)
	networkSecurityGroupExtensionGenerator = gen.Struct(reflect.TypeOf(NetworkSecurityGroupExtension{}), generators)

	return networkSecurityGroupExtensionGenerator
}
