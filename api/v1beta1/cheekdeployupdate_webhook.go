/*
Copyright 2022 aloys.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	errorApi "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var cheekdeployupdatelog = logf.Log.WithName("cheekdeployupdate-resource")

func (r *CheekDeployUpdate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-aloys-cheekdeployupdate-tech-v1beta1-cheekdeployupdate,mutating=true,failurePolicy=fail,sideEffects=None,groups=aloys.cheekdeployupdate.tech,resources=cheekdeployupdates,verbs=create;update,versions=v1beta1,name=mcheekdeployupdate.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &CheekDeployUpdate{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *CheekDeployUpdate) Default() {
	cheekdeployupdatelog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
	if r.Spec.DeploymentNamespace == "" {
		r.Spec.DeploymentNamespace = "default"
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:path=/validate-aloys-cheekdeployupdate-tech-v1beta1-cheekdeployupdate,mutating=false,failurePolicy=fail,sideEffects=None,groups=aloys.cheekdeployupdate.tech,resources=cheekdeployupdates,verbs=create;update,versions=v1beta1,name=vcheekdeployupdate.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &CheekDeployUpdate{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *CheekDeployUpdate) ValidateCreate() error {
	cheekdeployupdatelog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	if err := r.ValidateDeploymentName(); err != nil {
		return err
	}
	if err := r.ValidateDeploymentImage(); err != nil {
		return err
	}
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *CheekDeployUpdate) ValidateUpdate(old runtime.Object) error {
	cheekdeployupdatelog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	if err := r.ValidateDeploymentName(); err != nil {
		return err
	}
	if err := r.ValidateDeploymentImage(); err != nil {
		return err
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *CheekDeployUpdate) ValidateDelete() error {
	cheekdeployupdatelog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *CheekDeployUpdate) ValidateDeploymentName() error {
	if r.Spec.DeploymentName == "" {
		return errorApi.NewInvalid(GroupVersion.WithKind("CheekDeployUpdate").GroupKind(), r.Name,
			field.ErrorList{
				field.Invalid(field.NewPath("DeploymentName"),
					r.Spec.DeploymentName,
					"The Deployment Name cannot be empty"),
			},
		)
	}
	return nil
}

func (r *CheekDeployUpdate) ValidateDeploymentImage() error {
	if r.Spec.DeploymentImage == "" {
		return errorApi.NewInvalid(GroupVersion.WithKind("CheekDeployUpdate").GroupKind(), r.Name,
			field.ErrorList{
				field.Invalid(field.NewPath("DeploymentImage"),
					r.Spec.DeploymentName,
					"The Deployment image cannot be empty"),
			},
		)
	}
	return nil
}
