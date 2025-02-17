// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package cloudservice

import (
	"os"
	"testing"

	"gotest.tools/assert"
)

func TestGetCloudServiceType(t *testing.T) {
	os.Setenv(ContainerAppNameEnvVar, "test-name")
	assert.Equal(t, GetCloudServiceType(), &ContainerApp{})

	os.Unsetenv(ContainerAppNameEnvVar)
	assert.Equal(t, GetCloudServiceType(), &CloudRun{})
}
