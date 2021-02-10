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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	rcccmv0alpha0 "rc.ccm.dunescience.org/api/v0alpha0"
)

// DAQApplicationReconciler reconciles a DAQApplication object
type DAQApplicationReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	recorder record.EventRecorder
}

type appState struct {
	State         string `json:"state"`
	Transitioning bool   `json:"transitioning"`
}

// +kubebuilder:rbac:groups=rc.ccm.dunescience.org,resources=daqapplications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rc.ccm.dunescience.org,resources=daqapplications/status,verbs=get;update;patch

func (r *DAQApplicationReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("daqapplication", req.NamespacedName)

	// fetch DAQApplication resource that Reconcile was called for
	var daqapp rcccmv0alpha0.DAQApplication
	if err := r.Get(ctx, req.NamespacedName, &daqapp); err != nil {
		err = client.IgnoreNotFound(err) // make nil if NotFoundError
		if err != nil {
			log.Error(err, "unable to retrieve DAQ Application")
		}
		// if err != nil, this will trigger a re-queue with backoff
		// if err = nil, this was a NotFound error, we will be re-triggered when we actually exist...
		return ctrl.Result{}, err
	}

	// fetch Pod that this DAQApplication manages
	var pod corev1.Pod
	podName := client.ObjectKey{Name: daqapp.Spec.PodName, Namespace: req.Namespace}
	if err := r.Get(ctx, podName, &pod); err != nil {
		log.Error(err, "cannot find pod")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// get pod IP
	podIP := pod.Status.PodIP
	if podIP == "" {
		err := fmt.Errorf("pod has no IP")
		log.Error(err, "pod has no IP")
		return ctrl.Result{}, err
	}

	log.Info("Pod IP", "IP", podIP)

	httpclient := http.Client{Timeout: 5 * time.Second}
	httpurl := url.URL{Scheme: "http", Host: podIP, Path: "/api/v0/state"}
	resp, err := httpclient.Get(httpurl.String())
	if err != nil {
		log.Error(err, "cannot get DAQ State")
		return ctrl.Result{}, err
	}
	var body appState
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		log.Error(err, "cannot parse DAQ state")
		return ctrl.Result{}, err
	}

	if daqapp.Status.LastSeenState != body.State {
		r.recorder.Event(&daqapp, corev1.EventTypeNormal, "FSM Change", "DAQ App reached state '"+body.State+"'")
	}

	daqapp.Status.LastSeenState = body.State

	log.Info("retrieved DAQ state", "state", body.State)

	daqapp.Status.Status = "transitioning"
	if !body.Transitioning {
		nextCommand := ""
		switch body.State {
		case string(daqapp.Spec.DesiredState):
			daqapp.Status.Status = "happy"
		case "UNKNOWN":
			nextCommand = "init"
		case "INIT":
			nextCommand = "conf"
		case "CONFIGURED":
			nextCommand = "start"
		}

		if nextCommand == "" && daqapp.Status.Status != "happy" {
			daqapp.Status.Status = "stuck"
			err = fmt.Errorf("unknown FSM state")
			log.Error(err, "Unknown state. Don't know what to do", "state", body.State)
			return ctrl.Result{}, err
		} else {

			daqapp.Status.LastCommandSent = nextCommand
			resp, err = httpclient.Post(httpurl.String(), "application/json", bytes.NewBufferString("{\"command\": \""+nextCommand+"\"}"))
			if err != nil {
				log.Error(err, "cannot send DAQ Command")
				return ctrl.Result{}, err
			}
		}
	}

	log.Info("updating status")
	err = r.Status().Update(ctx, &daqapp)
	if err != nil {
		log.Error(err, "could not update status")
		return ctrl.Result{}, err
	}

	log.Info("DAQApplication reconciliation successfull")

	// DAQ Application state changes are not evented (yet?)
	// reconsile every 20s
	secondsToWait := 20 * time.Second
	if body.Transitioning {
		secondsToWait = 2 * time.Second
	}
	return ctrl.Result{RequeueAfter: secondsToWait}, nil
}

func (r *DAQApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.recorder = mgr.GetEventRecorderFor("DAQApplication")
	return ctrl.NewControllerManagedBy(mgr).
		For(&rcccmv0alpha0.DAQApplication{}).
		Complete(r)
}
