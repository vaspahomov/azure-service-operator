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

func Test_MongodbDatabaseCollectionExtension_WhenSerializedToJson_DeserializesAsEqual(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	parameters.MaxSize = 10
	properties := gopter.NewProperties(parameters)
	properties.Property(
		"Round trip of MongodbDatabaseCollectionExtension via JSON returns original",
		prop.ForAll(RunJSONSerializationTestForMongodbDatabaseCollectionExtension, MongodbDatabaseCollectionExtensionGenerator()))
	properties.TestingRun(t, gopter.NewFormatedReporter(true, 240, os.Stdout))
}

// RunJSONSerializationTestForMongodbDatabaseCollectionExtension runs a test to see if a specific instance of MongodbDatabaseCollectionExtension round trips to JSON and back losslessly
func RunJSONSerializationTestForMongodbDatabaseCollectionExtension(subject MongodbDatabaseCollectionExtension) string {
	// Serialize to JSON
	bin, err := json.Marshal(subject)
	if err != nil {
		return err.Error()
	}

	// Deserialize back into memory
	var actual MongodbDatabaseCollectionExtension
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

// Generator of MongodbDatabaseCollectionExtension instances for property testing - lazily instantiated by
//MongodbDatabaseCollectionExtensionGenerator()
var mongodbDatabaseCollectionExtensionGenerator gopter.Gen

// MongodbDatabaseCollectionExtensionGenerator returns a generator of MongodbDatabaseCollectionExtension instances for property testing.
func MongodbDatabaseCollectionExtensionGenerator() gopter.Gen {
	if mongodbDatabaseCollectionExtensionGenerator != nil {
		return mongodbDatabaseCollectionExtensionGenerator
	}

	generators := make(map[string]gopter.Gen)
	mongodbDatabaseCollectionExtensionGenerator = gen.Struct(reflect.TypeOf(MongodbDatabaseCollectionExtension{}), generators)

	return mongodbDatabaseCollectionExtensionGenerator
}
