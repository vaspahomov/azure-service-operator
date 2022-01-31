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

func Test_SqlDatabaseContainerStoredProcedureExtension_WhenSerializedToJson_DeserializesAsEqual(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	parameters.MaxSize = 10
	properties := gopter.NewProperties(parameters)
	properties.Property(
		"Round trip of SqlDatabaseContainerStoredProcedureExtension via JSON returns original",
		prop.ForAll(RunJSONSerializationTestForSqlDatabaseContainerStoredProcedureExtension, SqlDatabaseContainerStoredProcedureExtensionGenerator()))
	properties.TestingRun(t, gopter.NewFormatedReporter(true, 240, os.Stdout))
}

// RunJSONSerializationTestForSqlDatabaseContainerStoredProcedureExtension runs a test to see if a specific instance of SqlDatabaseContainerStoredProcedureExtension round trips to JSON and back losslessly
func RunJSONSerializationTestForSqlDatabaseContainerStoredProcedureExtension(subject SqlDatabaseContainerStoredProcedureExtension) string {
	// Serialize to JSON
	bin, err := json.Marshal(subject)
	if err != nil {
		return err.Error()
	}

	// Deserialize back into memory
	var actual SqlDatabaseContainerStoredProcedureExtension
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

// Generator of SqlDatabaseContainerStoredProcedureExtension instances for property testing - lazily instantiated by
//SqlDatabaseContainerStoredProcedureExtensionGenerator()
var sqlDatabaseContainerStoredProcedureExtensionGenerator gopter.Gen

// SqlDatabaseContainerStoredProcedureExtensionGenerator returns a generator of SqlDatabaseContainerStoredProcedureExtension instances for property testing.
func SqlDatabaseContainerStoredProcedureExtensionGenerator() gopter.Gen {
	if sqlDatabaseContainerStoredProcedureExtensionGenerator != nil {
		return sqlDatabaseContainerStoredProcedureExtensionGenerator
	}

	generators := make(map[string]gopter.Gen)
	sqlDatabaseContainerStoredProcedureExtensionGenerator = gen.Struct(reflect.TypeOf(SqlDatabaseContainerStoredProcedureExtension{}), generators)

	return sqlDatabaseContainerStoredProcedureExtensionGenerator
}
