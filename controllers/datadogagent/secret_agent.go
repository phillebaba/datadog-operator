// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package datadogagent

import (
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	apicommon "github.com/DataDog/datadog-operator/apis/datadoghq/common"
	datadoghqv1alpha1 "github.com/DataDog/datadog-operator/apis/datadoghq/v1alpha1"
	"github.com/DataDog/datadog-operator/controllers/datadogagent/object"
	"github.com/DataDog/datadog-operator/pkg/config"
	"github.com/DataDog/datadog-operator/pkg/controller/utils"
)

func (r *Reconciler) manageAgentSecret(logger logr.Logger, dda *datadoghqv1alpha1.DatadogAgent) (reconcile.Result, error) {
	return r.manageSecret(logger, managedSecret{name: utils.GetDefaultCredentialsSecretName(dda), requireFunc: needAgentSecret, createFunc: newAgentSecret}, dda)
}

func newAgentSecret(name string, dda *datadoghqv1alpha1.DatadogAgent) *corev1.Secret {
	labels := object.GetDefaultLabels(dda, apicommon.DefaultClusterAgentResourceSuffix, getClusterAgentVersion(dda))
	annotations := object.GetDefaultAnnotations(dda)

	creds := dda.Spec.Credentials
	data := datadoghqv1alpha1.GetKeysFromCredentials(&creds.DatadogCredentials)

	if creds.Token != "" {
		data[apicommon.DefaultTokenKey] = []byte(creds.Token)
	} else if isClusterAgentEnabled(dda.Spec.ClusterAgent) {
		defaultedToken := datadoghqv1alpha1.DefaultedClusterAgentToken(&dda.Status)
		if defaultedToken != "" {
			data[apicommon.DefaultTokenKey] = []byte(defaultedToken)
		}
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   dda.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Type: corev1.SecretTypeOpaque,
		Data: data,
	}
	return secret
}

// needAgentSecret checks if a secret should be used or created.
func needAgentSecret(dda *datadoghqv1alpha1.DatadogAgent) bool {
	// If credentials is nil, there is nothing to create a secret from.
	if dda.Spec.Credentials == nil {
		return false
	}

	// If API key, app key _and_ token don't need a new secret, then don't create one.
	if datadoghqv1alpha1.CheckAPIKeySufficiency(&dda.Spec.Credentials.DatadogCredentials, config.DDAPIKeyEnvVar) &&
		datadoghqv1alpha1.CheckAppKeySufficiency(&dda.Spec.Credentials.DatadogCredentials, config.DDAppKeyEnvVar) &&
		!isClusterAgentEnabled(dda.Spec.ClusterAgent) {
		return false
	}

	return true
}
