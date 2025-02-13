// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	SubscriptionIDVar   = "AZURE_SUBSCRIPTION_ID"
	targetNamespacesVar = "AZURE_TARGET_NAMESPACES"
	operatorModeVar     = "AZURE_OPERATOR_MODE"
	syncPeriodVar       = "AZURE_SYNC_PERIOD"
	podNamespaceVar     = "POD_NAMESPACE"
)

// Values stores configuration values that are set for the operator.
type Values struct {
	// SubscriptionID is the Azure subscription the operator will use
	// for ARM communication.
	SubscriptionID string

	// PodNamespace is the namespace the operator pods are running in.
	PodNamespace string

	// OperatorMode determines whether the operator should run
	// watchers, webhooks or both.
	OperatorMode OperatorMode

	// TargetNamespaces lists the namespaces the operator will watch
	// for Azure resources (if the mode includes running watchers). If
	// it's empty the operator will watch all namespaces.
	TargetNamespaces []string

	// SyncPeriod is the frequency at which resources are re-reconciled with Azure
	// when there have been no triggering changes in the Kubernetes resources. This sync
	// exists to detect and correct changes that happened in Azure that Kubernetes is not
	// aware about. BE VERY CAREFUL setting this value low - even a modest number of resources
	// can cause subscription level throttling if they are re-synced frequently.
	// If nil, no sync is performed. Durations are specified as "1h", "15m", or "60s". See
	// https://pkg.go.dev/time#ParseDuration for more details.
	//
	// This can be set to nil by specifying empty string for AZURE_SYNC_PERIOD explicitly in
	// the config.
	SyncPeriod *time.Duration
}

var _ fmt.Stringer = Values{}

// Returns the configuration as a string
func (v Values) String() string {
	return fmt.Sprintf(
		"SubscriptionID:%s/PodNamespace:%s/OperatorMode:%s/TargetNamespaces:%s/SyncPeriod:%s",
		v.SubscriptionID,
		v.PodNamespace,
		v.OperatorMode,
		strings.Join(v.TargetNamespaces, "|"),
		v.SyncPeriod)
}

// ReadFromEnvironment loads configuration values from the AZURE_*
// environment variables.
func ReadFromEnvironment() (Values, error) {
	var result Values
	modeValue := os.Getenv(operatorModeVar)
	if modeValue == "" {
		result.OperatorMode = OperatorModeBoth
	} else {
		mode, err := ParseOperatorMode(modeValue)
		if err != nil {
			return Values{}, err
		}
		result.OperatorMode = mode
	}

	var err error

	result.SubscriptionID = os.Getenv(SubscriptionIDVar)
	result.PodNamespace = os.Getenv(podNamespaceVar)
	result.TargetNamespaces = parseTargetNamespaces(os.Getenv(targetNamespacesVar))
	result.SyncPeriod, err = parseSyncPeriod()

	if err != nil {
		return result, errors.Wrapf(err, "parsing %q", syncPeriodVar)
	}

	// Not calling validate here to support using from tests where we
	// don't require consistent settings.
	return result, nil
}

// ReadAndValidate loads the configuration values and checks that
// they're consistent.
func ReadAndValidate() (Values, error) {
	result, err := ReadFromEnvironment()
	if err != nil {
		return Values{}, err
	}
	err = result.Validate()
	if err != nil {
		return Values{}, err
	}
	return result, nil
}

// Validate checks whether the configuration settings are consistent.
func (v Values) Validate() error {
	if v.SubscriptionID == "" {
		return errors.Errorf("missing value for %s", SubscriptionIDVar)
	}
	if v.PodNamespace == "" {
		return errors.Errorf("missing value for %s", podNamespaceVar)
	}
	if !v.OperatorMode.IncludesWatchers() && len(v.TargetNamespaces) > 0 {
		return errors.Errorf("%s must include watchers to specify target namespaces", targetNamespacesVar)
	}
	return nil
}

// parseTargetNamespaces splits a comma-separated string into a slice
// of strings with spaces trimmed.
func parseTargetNamespaces(fromEnv string) []string {
	if len(strings.TrimSpace(fromEnv)) == 0 {
		return nil
	}
	items := strings.Split(fromEnv, ",")
	// Remove any whitespace used to separate items.
	for i, item := range items {
		items[i] = strings.TrimSpace(item)
	}
	return items
}

// parseSyncPeriod parses the sync period from the environment
func parseSyncPeriod() (*time.Duration, error) {
	syncPeriodStr := envOrDefault(syncPeriodVar, "15m")
	if syncPeriodStr == "" {
		return nil, nil
	}

	syncPeriod, err := time.ParseDuration(syncPeriodStr)
	if err != nil {
		return nil, err
	}
	return &syncPeriod, nil
}

// envOrDefault returns the value of the specified env variable or the default value if
// the env variable was not set.
func envOrDefault(env string, def string) string {
	result, specified := os.LookupEnv(env)
	if !specified {
		return def
	}

	return result
}
