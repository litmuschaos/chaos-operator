/*
Copyright 2019 LitmusChaos Authors

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
	"fmt"
	"os"
	"runtime"
	"strings"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	"github.com/operator-framework/operator-sdk/pkg/leader"
	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	"github.com/operator-framework/operator-sdk/pkg/metrics"
	sdkVersion "github.com/operator-framework/operator-sdk/version"
	"github.com/spf13/pflag"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/litmuschaos/chaos-operator/pkg/analytics"
	"github.com/litmuschaos/chaos-operator/pkg/apis"
	"github.com/litmuschaos/chaos-operator/pkg/controller"
)

// Change below variables to serve metrics on different host or port.
var (
	metricsHost       = "0.0.0.0"
	metricsPort int32 = 8383
	log               = logf.Log.WithName("cmd")
)

func main() {
	// initializing the log configuration
	initializingLogConfiguration()

	// printing the operator and go configuration
	printVersion()

	// setting up initial configuration of chaos-operator
	mgr, err := initialConfiguration()
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Start the Chaos-Operator
	log.Info("Starting the Chaos-Operator...")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "Chaos-Operator exited non-zero")
		os.Exit(1)
	}
}

// initialize the log configuration
func initializingLogConfiguration() {
	// Add the zap logger flag set to the CLI. The flag set must
	// be added before calling pflag.Parse().
	pflag.CommandLine.AddFlagSet(zap.FlagSet())

	// Add flags registered by imported packages (e.g. glog and
	// controller-runtime)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	pflag.Parse()

	// Use a zap logr.Logger implementation. If none of the zap
	// flags are configured (or if the zap flag set is not being
	// used), this defaults to a production zap logger.
	//
	// The logger instantiated here can be changed to any logger
	// implementing the logr.Logger interface. This logger will
	// be propagated through the whole operator, generating
	// uniform and structured logs.
	logf.SetLogger(zap.Logger())
}

func printVersion() {
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
	log.Info(fmt.Sprintf("Version of operator-sdk: %v", sdkVersion.Version))
}

// initializing the configuration of chaos-operator
func initialConfiguration() (manager.Manager, error) {
	//setting up leader and analytics
	if err := initializingLeaderAndAnalytics(); err != nil {
		return nil, err
	}

	// creating metrics service
	cfg, namespace, err := initializeMetricsService()
	if err != nil {
		return nil, err
	}

	// registering the components of chaos-operator
	mgr, err := registerComponents(cfg, namespace)
	if err != nil {
		return nil, err
	}
	return mgr, nil
}

func initializingLeaderAndAnalytics() error {
	// Become the leader before proceeding
	if err := leader.Become(context.TODO(), "chaos-operator-lock"); err != nil {
		return err
	}

	// Trigger the Analytics if it's enabled
	if isAnalytics := strings.ToUpper(os.Getenv("ANALYTICS")); isAnalytics != "FALSE" {
		if err := analytics.TriggerAnalytics(); err != nil {
			log.Error(err, "")
		}
	}

	return nil
}

func initializeMetricsService() (*rest.Config, string, error) {
	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		return cfg, "", err
	}

	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get watch namespace %v", err)
	}

	// Create Service object to expose the metrics port(s).
	servicePorts := []v1.ServicePort{{Port: metricsPort, Name: metrics.OperatorPortName, Protocol: v1.ProtocolTCP, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: metricsPort}}}
	if _, err := metrics.CreateMetricsService(context.TODO(), cfg, servicePorts); err != nil {
		log.Info("Could not create metrics Service", "error", err.Error())
	}
	return cfg, namespace, nil
}

func registerComponents(cfg *rest.Config, namespace string) (manager.Manager, error) {

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{Namespace: namespace, MetricsBindAddress: fmt.Sprintf("%s:%d", metricsHost, metricsPort)})
	if err != nil {
		return mgr, err
	}

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		return nil, err
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		return nil, err
	}

	return mgr, nil
}
