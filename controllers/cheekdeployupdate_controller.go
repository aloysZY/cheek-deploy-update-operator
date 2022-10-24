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

package controllers

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	aloyscheekdeployupdatev1beta1 "cheek-deploy-update-operator/api/v1beta1"
)

// CheekDeployUpdateReconciler reconciles a CheekDeployUpdate object
type CheekDeployUpdateReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	logger   logr.Logger
	newImage string
	oldImage string
}

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update
// +kubebuilder:rbac:groups="",resources=pods,verbs=list;watch
// +kubebuilder:rbac:groups=aloys.cheekdeployupdate.tech,resources=cheekdeployupdates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aloys.cheekdeployupdate.tech,resources=cheekdeployupdates/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aloys.cheekdeployupdate.tech,resources=cheekdeployupdates/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CheekDeployUpdate object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *CheekDeployUpdateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.logger = log.FromContext(ctx)

	// TODO(user): your logic here
	// 检查自定义 cr 是否存在,获取 cr 设置的信息
	start := time.Now()
	cheekDeployUpdate := aloyscheekdeployupdatev1beta1.CheekDeployUpdate{}
	if err := r.Get(ctx, req.NamespacedName, &cheekDeployUpdate); err != nil {
		r.logger.Error(err, "unable to fetch cheekDeployUpdate")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var deployment appsv1.Deployment
	if err := r.Get(ctx, types.NamespacedName{
		Namespace: cheekDeployUpdate.Spec.DeploymentNamespace,
		Name:      cheekDeployUpdate.Spec.DeploymentName,
	}, &deployment); err != nil {
		r.logger.Error(err, "unable to find deploy", "name", cheekDeployUpdate.Spec.DeploymentName)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	r.logger.Info("Operation deployment for escalation", "deployment namespace", deployment.Namespace, "depl	oyment name", deployment.Name, "deployment namespace", deployment.Namespace)

	// 判断 deploy 的状态
	// 返回 true，更新成功，返回 false 更新失败
	isFinished := func(deployment *appsv1.Deployment) (bool, error) {
		// 这里要在去获取一下状态，不能直接判断
		// 这个判断应该是没有必要的，但是不排除有其他终端在操作
		for {
			// 这里会有一个获取的延迟，或者要判断几次获取都是可能也行
			// 执行更新后直接去请求，资源状态可能还没变更
			r.logger.Info("The upgrade is in progress", "deployment name", deployment.Name)
			// 1秒日志刷新的太快了，这里需要优化
			time.Sleep(time.Second * 15)
			if err := r.Get(ctx, types.NamespacedName{
				Namespace: deployment.Namespace,
				Name:      deployment.Name,
			}, deployment); err != nil {
				return false, err
			}
			for _, c := range deployment.Status.Conditions {
				if c.Type == appsv1.DeploymentProgressing {
					switch c.Reason {
					case "NewReplicaSetAvailable":
						return true, nil
					case "ProgressDeadlineExceeded":
						return false, nil
					}
				}
			}
		}
	}

	r.newImage = cheekDeployUpdate.Spec.DeploymentImage
	r.oldImage = deployment.Spec.Template.Spec.Containers[0].Image
	if r.newImage != r.oldImage {
		r.logger.Info("List of Pods that need to be upgraded")
		if err := r.getPodList(ctx, &deployment); err != nil {
			return ctrl.Result{}, err
		}
		r.logger.Info("deployment starts the upgrade", "deployment name", deployment.Name, "deployment namespace", deployment.Namespace)
		deployment.Spec.Template.Spec.Containers[0].Image = r.newImage
		// ProgressDeadlineSeconds 给出的是一个秒数值，Deployment 控制器在（通过 Deployment 状态） 标示 Deployment 进展停滞之前，需要等待所给的时长。默认 600 秒，感觉太久了
		progressDeadlineSeconds := int32(300)
		deployment.Spec.ProgressDeadlineSeconds = &progressDeadlineSeconds
		if err := r.Update(ctx, &deployment); err != nil {
			r.logger.Error(err, "Failed to update Deployment", "Deployment.Namespace", cheekDeployUpdate.Spec.DeploymentNamespace, "Deployment.Name", cheekDeployUpdate.Spec.DeploymentName)
			return reconcile.Result{}, err
		}

		finished, err := isFinished(&deployment)
		if err != nil {
			r.logger.Error(err, "unable to find deploy", "name", deployment.Name)
			return ctrl.Result{}, err
		}
		if !finished {
			// 状态错误就要回滚
			r.logger.Info("The deployment version is rolled back", "deployment name", deployment.Name, "deploy image", r.newImage)
			// 传入cr尝试看看，会不会带动 deploy 的回滚，其实就是重新执行了 CR，然后 deploy 再次更新了，感觉这个比直接更新 deoploy 好，因为 CR 也更新了回退 CR 有一个问题，就是多个 deloy 的时候就要一起回退了
			if err := r.rollingBackCR(ctx, &cheekDeployUpdate); err != nil {
				r.logger.Error(err, "Failed to roll back the version. Procedure", "deployment name", deployment.Name, "Wrong deployment image version", r.oldImage)
				return ctrl.Result{}, err
			}
			// 重置 CR 这里就不能这样输出了，不需要输出，直接返回，开始下一个循环
			// r.logger.Info("The version rollback succeeded. Procedure", "deployment name", deployment.Name, "deploy image", r.oldImage)
			return ctrl.Result{}, nil
		}
		r.logger.Info("deployment Current Version", "deployment name", deployment.Name, "deploy image", deployment.Spec.Template.Spec.Containers[0].Image)
		r.logger.Info("List of deployment Pods")
		if err := r.getPodList(ctx, &deployment); err != nil {
			return ctrl.Result{}, err
		}
	} else {
		r.logger.Info("The deployment image version is the same as the old version and does not require operation", "deployment name", deployment.Name, "deployment image", r.newImage)
		return ctrl.Result{}, nil
	}

	r.logger.Info("deployment as well as operation completion", "deployment name", deployment.Name, "deployment image", r.newImage)

	cost := time.Since(start)
	r.logger.Info("deployment upgrade time", "time", cost)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CheekDeployUpdateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aloyscheekdeployupdatev1beta1.CheekDeployUpdate{}).
		Complete(r)
}

func (r *CheekDeployUpdateReconciler) getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		for _, v := range pod.Status.ContainerStatuses {
			if v.Ready {
				podNames = append(podNames, pod.Name)
			}
		}
	}
	return podNames
	// // 不需要上面的额外判断了，只要上面限制了访问的时间
	// for _, pod := range pods {
	// 	podNames = append(podNames, pod.Name)
	// }
	// return podNames
}

func (r *CheekDeployUpdateReconciler) getPodList(ctx context.Context, deployment *appsv1.Deployment) error {
	labels := deployment.GetLabels()
	listOpts := []client.ListOption{
		client.InNamespace(deployment.Namespace),
		client.MatchingLabels(labels),
	}
	var podList corev1.PodList
	// time.Sleep(time.Second * 15)
	if err := r.List(ctx, &podList, listOpts...); err != nil {
		r.logger.Error(err, "unable to list pods", "pod list", podList)
		return err
	}
	podNames := r.getPodNames(podList.Items)
	r.logger.Info("find the pod list", "pod name ", podNames)
	return nil
}

func (r *CheekDeployUpdateReconciler) rollingBackCR(ctx context.Context, acdu *aloyscheekdeployupdatev1beta1.CheekDeployUpdate) error {
	acdu.Spec.DeploymentImage = r.oldImage
	if err := r.Update(ctx, acdu); err != nil {
		return err
	}
	return nil
}

// 暂时不用了，被rollingBackCR替换了，直接回退 CR
// func (r *CheekDeployUpdateReconciler) rollingBackDeployment(ctx context.Context, deployment *appsv1.Deployment) error {
// 	deployment.Spec.Template.Spec.Containers[0].Image = r.oldImage
// 	deployment.ResourceVersion = ""
// 	if err := r.Update(ctx, deployment); err != nil {
// 		return err
// 	}
// 	return nil
// }
