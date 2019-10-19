package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	"github.com/litmuschaos/chaos-operator/pkg/apis"
	"github.com/litmuschaos/chaos-operator/pkg/controller"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	"github.com/operator-framework/operator-sdk/pkg/leader"
	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	"github.com/operator-framework/operator-sdk/pkg/metrics"
	sdkVersion "github.com/operator-framework/operator-sdk/version"
	"github.com/spf13/pflag"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
	"k8s.io/client-go/rest"
	v1 "k8s.io/api/core/v1"
)

// Change below variables to serve metrics on different host or port.
const (
	host =  "0.0.0.0"
	port = 8383
	lockName = "chaos-operator-lock"
)
var (
	metricsHost       = host
	metricsPort int32 = port
)
var log = logf.Log.WithName("cmd")

func printVersion() {
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
	log.Info(fmt.Sprintf("Version of operator-sdk: %v", sdkVersion.Version))
}

// unit testing is done to these
func getK8Namespace()(string, error){
	return k8sutil.GetWatchNamespace()
}

func getK8RestConfig()(*rest.Config, error){
	return config.GetConfig()
}

func becomeLeader(ctx context.Context)error{
	return leader.Become(ctx, lockName)
}

//
func createNewManager(cfg *rest.Config, namespace string)(manager.Manager, error){
	return manager.New(cfg, manager.Options{
		Namespace:          namespace,
		MetricsBindAddress: fmt.Sprintf("%s:%d", metricsHost, metricsPort),
	})
}

// Decoupled these statements to create so that unit testing is possible
func addToApiSchema(mgr manager.Manager)error{
	return apis.AddToScheme(mgr.GetScheme())
}

// Decoupled so that unit testing is possible
// These packages are gonna have its own kind of unit testing
func addToControllerSchema(mgr manager.Manager)error{
	return controller.AddToManager(mgr)
}

// The function exposes port
func addToMetricsPort(ctx context.Context, metricsPort int32)(* v1.Service, error){
	return metrics.ExposeMetricsPort(ctx, metricsPort)
}

func startCmd(mgr manager.Manager)error{
	return mgr.Start(signals.SetupSignalHandler())
}

func main() {

	log.Info("-----------NEW RUN--------")
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

	printVersion()

	namespace, err := getK8Namespace()
	if err != nil {
		log.Error(err, "Failed to get watch namespace")
		os.Exit(1)
	}

	// Get a config to talk to the apiserver
	cfg, err := getK8RestConfig()
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	ctx := context.TODO()

	// Become the leader before proceeding
	err = becomeLeader(ctx)
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := createNewManager(cfg, namespace)
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	log.Info("Registering Components.")

	// Setup Scheme for all resources
	if err := addToApiSchema(mgr); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Setup all Controllers
	if err := addToControllerSchema(mgr); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Create Service object to expose the metrics port.
	_, err = addToMetricsPort(ctx, metricsPort)
	if err != nil {
		log.Info(err.Error())
	}

	log.Info("Starting the Cmd.")

	// Start the Cmd
	if err := startCmd(mgr); err != nil {
		log.Error(err, "Manager exited non-zero")
		os.Exit(1)
	}
}
