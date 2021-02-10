/*


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
	"strings"
	"time"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	rcccmv0alpha0 "rc.ccm.dunescience.org/api/v0alpha0"
)

// PartitionReconciler reconciles a Partition object
type PartitionReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

type PhysicsDetectorType string

const (
	TPCType PhysicsDetectorType = "TPC"
	PDSType PhysicsDetectorType = "PDS"
)

// +kubebuilder:rbac:groups=rc.ccm.dunescience.org,resources=partitions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rc.ccm.dunescience.org,resources=partitions/status,verbs=get;update;patch

func (r *PartitionReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("partition", req.NamespacedName)

	// fetch Partition resource that Reconcile was called for
	var partition rcccmv0alpha0.Partition
	if err := r.Get(ctx, req.NamespacedName, &partition); err != nil {
		err = client.IgnoreNotFound(err) // make nil if NotFoundError
		if err != nil {
			log.Error(err, "unable to retrieve partition")
		}
		// if err != nil, this will trigger a re-queue with backoff
		// if err = nil, this was a NotFound error, we will be re-triggered when we actually exist...
		return ctrl.Result{}, err
	}

	applyOpts := []client.PatchOption{client.ForceOwnership, client.FieldOwner("partition-controller")}

	deploys := []appsv1.Deployment{}
	daqApps := []*rcccmv0alpha0.DAQApplication{}
	for _, moduleSpec := range partition.Spec.Resources {
		for _, tpcApa := range moduleSpec.TPC.APAs {
			deployment, err := r.deploymentFor(partition, TPCType, moduleSpec.Module, "APA", tpcApa, partition.Spec.ConfigName)
			if err != nil {
				log.Error(err, "cannot create deployment")
				return ctrl.Result{}, err
			}
			daqApp, err := r.daqAppFor(partition, TPCType, moduleSpec.Module, "APA", tpcApa, partition.Spec.ConfigName)
			if err != nil {
				log.Error(err, "cannot create DAQ App")
				return ctrl.Result{}, err
			}
			deploys = append(deploys, deployment)
			if daqApp != nil {
				daqApps = append(daqApps, daqApp)
			}
		}
		for _, pdsApa := range moduleSpec.PDS.APAs {
			deployment, err := r.deploymentFor(partition, PDSType, moduleSpec.Module, "APA", pdsApa, partition.Spec.ConfigName)
			if err != nil {
				log.Error(err, "cannot create deployment")
				return ctrl.Result{}, err
			}
			daqApp, err := r.daqAppFor(partition, PDSType, moduleSpec.Module, "APA", pdsApa, partition.Spec.ConfigName)
			if err != nil {
				log.Error(err, "cannot create DAQ App")
				return ctrl.Result{}, err
			}
			deploys = append(deploys, deployment)
			if daqApp != nil {
				daqApps = append(daqApps, daqApp)
			}
		}
	}

	for _, deployment := range deploys {
		err := r.Patch(ctx, &deployment, client.Apply, applyOpts...)
		if err != nil {
			log.Error(err, "cannot patch deployment")
			return ctrl.Result{}, err
		}
	}

	partition.Status.Status = "happy"
	for _, daqApp := range daqApps {
		err := r.Patch(ctx, daqApp, client.Apply, applyOpts...)
		if err != nil {
			log.Error(err, "cannot patch DAQApplication")
			return ctrl.Result{}, err
		}
		var realDaqApp rcccmv0alpha0.DAQApplication
		appName := client.ObjectKey{Name: daqApp.ObjectMeta.Name, Namespace: req.Namespace}
		if err := r.Get(ctx, appName, &realDaqApp); err == nil {
			log.Info("real state of daq app", "state", realDaqApp.Status.Status)
			if realDaqApp.Status.Status != "happy" {
				partition.Status.Status = "sad"
			}
		}
	}

	log.Info("updating status")
	err := r.Status().Update(ctx, &partition)
	if err != nil {
		log.Error(err, "could not update status")
		return ctrl.Result{}, err
	}

	log.Info("Partition reconciliation successfull")

	return ctrl.Result{RequeueAfter: 20 * time.Second}, nil
}

func (r *PartitionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&rcccmv0alpha0.Partition{}).
		Complete(r)
}

func (r *PartitionReconciler) deploymentFor(partition rcccmv0alpha0.Partition, physicsType PhysicsDetectorType, moduleID string, hardwareType string, hardwareID string, configName string) (appsv1.Deployment, error) {
	labelMap := map[string]string{
		"app":       "daq",
		"partition": partition.Name,
		"module":    moduleID,
		"physics":   string(physicsType),
		"hwtype":    hardwareType,
		"hwid":      hardwareID}

	name := strings.ToLower(strings.Join([]string{"ru", hardwareID, hardwareType, string(physicsType), moduleID}, "-"))
	replicas := int32(1)

	podResources := corev1.ResourceList{
		corev1.ResourceName("rc.ccm/" + name): resource.MustParse("1"),
	}

	deployment := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: appsv1.SchemeGroupVersion.String(),
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: partition.Namespace,
			Labels:    labelMap,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: labelMap,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labelMap,
				},
				Spec: corev1.PodSpec{
					Tolerations: []corev1.Toleration{{
						Key:      "dedicated",
						Operator: corev1.TolerationOpEqual,
						Value:    hardwareType,
						Effect:   corev1.TaintEffectNoSchedule,
					}},
					Containers: []corev1.Container{{
						Name:            name,
						Image:           "gitlab-registry.cern.ch/gdirkx/dune-daq-app-mmvp/app:latest",
						ImagePullPolicy: corev1.PullAlways,
						Ports: []corev1.ContainerPort{
							{ContainerPort: 80, Name: "http", Protocol: "TCP"},
						},
						VolumeMounts: []corev1.VolumeMount{{
							MountPath: "/cvmfs",
							Name:      "cvmfs",
						}, {
							MountPath: "/mnt",
							Name:      "config",
						}},
						Resources: corev1.ResourceRequirements{
							Requests: podResources,
							Limits:   podResources,
						},
					}},
					RestartPolicy: corev1.RestartPolicyAlways,
					Volumes: []corev1.Volume{{
						Name: "cvmfs",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/DUNE-RC-RC/cvmfs",
							},
						},
					}, {
						Name: "config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: configName,
								},
							},
						},
					}},
				},
			},
		},
	}

	err := ctrl.SetControllerReference(&partition, &deployment, r.Scheme)
	return deployment, err
}

func (r *PartitionReconciler) daqAppFor(partition rcccmv0alpha0.Partition, physicsType PhysicsDetectorType, moduleID string, hardwareType string, hardwareID string, configName string) (*rcccmv0alpha0.DAQApplication,
	error) {
	labelMap := map[string]string{
		"app":       "daq",
		"partition": partition.Name,
		"module":    moduleID,
		"physics":   string(physicsType),
		"hwtype":    hardwareType,
		"hwid":      hardwareID}

	name := strings.ToLower(strings.Join([]string{"ru", hardwareID, hardwareType, string(physicsType), moduleID}, "-"))

	listOptions := []client.ListOption{
		client.MatchingLabels(labelMap),
		// in the right namespace
		client.InNamespace(partition.GetNamespace()),
	}
	var list corev1.PodList
	if err := r.List(context.Background(), &list, listOptions...); err != nil {
		return nil, nil
	}
	if len(list.Items) == 0 {
		r.Log.Info("no pods match our labels")
		return nil, nil
	}
	if len(list.Items) != 1 {
		r.Log.Info("more than one pod matching", "amount", len(list.Items))
	}
	podName := list.Items[0].ObjectMeta.Name

	daqApp := rcccmv0alpha0.DAQApplication{
		TypeMeta: metav1.TypeMeta{
			APIVersion: rcccmv0alpha0.GroupVersion.String(),
			Kind:       "DAQApplication",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: partition.Namespace,
			Labels:    labelMap,
		},
		Spec: rcccmv0alpha0.DAQApplicationSpec{
			PodName:      podName,
			DesiredState: "STARTED",
		},
	}

	err := ctrl.SetControllerReference(&partition, &daqApp, r.Scheme)
	return &daqApp, err
}
