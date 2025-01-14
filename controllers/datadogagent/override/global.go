// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package override

import (
	"github.com/DataDog/datadog-operator/apis/datadoghq/v2alpha1"
	"github.com/DataDog/datadog-operator/controllers/datadogagent/feature"
	corev1 "k8s.io/api/core/v1"
)

// ApplyGlobalSettings use to apply global setting to a PodTemplateSpec
func ApplyGlobalSettings(manager feature.PodTemplateManagers, config *v2alpha1.GlobalConfig) *corev1.PodTemplateSpec {
	// TODO(operator-ga): implement ApplyGlobalSettings

	// set image registry

	return manager.PodTemplateSpec()
}
