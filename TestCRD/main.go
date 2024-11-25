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

package main

import (
	"context"
	"flag"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	customv1 "k8s.io/kubernetes/api/v1"
	"k8s.io/kubernetes/controllers"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(customv1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "e29887c8.test.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.TestCRDReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("TestCRD"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "TestCRD")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	// 1.构建context
	//ctx := context.CancelFunc(context.Background)
	ctx, cancel := context.WithCancel(context.Background())

	// 2.select等待选举结果
	c := leaderElection(ctx, mgr.GetClient())
	select {
	case ctx:
		setupLog.Info("exit ")
		os.Exit(0)
	case c:
		mgr.Start(ctrl.SetupSignalHandler())
	default:
		setupLog.Info("Waiting for leader election!")
	}

	// 3.销毁退出


	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}



}

// 多副本选举，选举可以通过configmap/lease/endpoint三种方式实现，k8s1.18之后，删除了configmap的方式
func leaderElection(ctx context.Context, k8sClient client.Client)<-chan{
	leaseLockName := ""
	leaseLockNamespace := ""

	resourcelock.LeaderElectionRecordToLeaseSpec()
	resourcelock.LeaseSpecToLeaderElectionRecord()

	leaseLock := &resourcelock.LeaseLock{
		LeaseMeta: v1.ObjectMeta{
			Name: leaseLockName,
			Namespace: leaseLockNamespace,
		},
		Client: k8sClient,
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: "TestCRD",
			EventRecorder: nil,
		},
	}

	//跑选举库做选举
	leaderelection.RunOrDie(ctx,
		leaderelection.LeaderElectionConfig{
			Lock: leaseLock,

			LeaseDuration: 60*time.Second,
			RenewDeadline: 15*time.Second,
			RetryPeriod: 5*time.Second,

			ReleaseOnCancel: true,

			Callbacks: leaderelection.LeaderCallbacks{
				OnNewLeader: func(identity string){
					setupLog.Info("TestCRD")
				},
				OnStartedLeading: func(ctx context.Context){
					setupLog.Info("TestCRD")
				},
				OnStoppedLeading: func(){
					setupLog.Info("TestCRD")
				},
			},
		})

	//获取lease，确认lease是否被占用

	//设置lease，确认当选leader

	//当选leader，加载配置及缓存，进入主程序

	//释放leader

	return c(<-chan)
}
