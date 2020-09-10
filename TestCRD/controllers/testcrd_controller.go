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
	v1 "k8s.io/api/core/v1"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	customv1 "k8s.io/kubernetes/api/v1"

	//testCrdV1 "github.com/hfeng101/operator/TestCRD/api/v1"
)

// TestCRDReconciler reconciles a TestCRD object
type TestCRDReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=custom.test.com,resources=testcrds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=custom.test.com,resources=testcrds/status,verbs=get;update;patch

func (r *TestCRDReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("testcrd", req.NamespacedName)

	// your logic here
	//get pods in ns:default

	crdObject := crd
	if err := r.Get(ctx, req.Namespace, )

	return ctrl.Result{}, nil
}

func (r *TestCRDReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&customv1.TestCRD{}).
		Complete(r)
}

//func (r *TestCRDReconciler) GetPodsByNamespace(namespace string, ctx context.Context) []*v1.Pod{
//	r.Client.Status()
//
//}
