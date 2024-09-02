package e2e

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mittwald/go-helm-client"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	clientset "k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	defaultKubeConfigPath = ".kube/config"
	defaultNS             = "kube-system"
	hashLength            = 6
	hashSeed              = "0123456789abcdefghijklmnopqrstuvwxyz"
)

var ctx *Context

// init initialize framework or die
func init() {
	if os.Getenv("SKIP_E2E_FRAMEWORK_INIT") != "" {
		fmt.Println("Framework initialization skipped")
		return
	}

	var err error
	ctx, err = NewFrameworkWithDefaultNamespace()
	if err != nil {
		panic(fmt.Sprintf("couldn't initialize: %v", err))
	}
}

// Context is container for components can be reused for each test
type Context struct {
	ClientConfig *rest.Config
	ClientSet    clientset.Interface
	Namespace    string
	HelmClient   helmclient.Client
	HelmChart    *helmclient.ChartSpec
}

func BackgroundContext() *Context {
	return ctx
}

// GenericBeforeSuiteFunc is a helper function to create a new framework with default namespace, this function should be called in BeforeSuite
func GenericBeforeSuiteFunc() {
	// list namespaces to ensure api-server accessibility
	list, err := ctx.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	Expect(err).To(BeNil())
	Expect(len(list.Items)).Should(BeNumerically(">=", 1))
}

// GenericRunTestSuiteFunc is a helper function to trigger ginkgo test suite, this function should be called within go test function with testing pkg.
func GenericRunTestSuiteFunc(t *testing.T, suitSpec string) {
	RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, suitSpec)
}

func NewFrameworkWithNamespace(namespace string) (*Context, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	kubeconfig := filepath.Join(homePath, defaultKubeConfigPath)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	helmChartClient, err := helmclient.NewClientFromRestConf(&helmclient.RestConfClientOptions{
		Options: &helmclient.Options{
			Namespace: namespace,
		},
		RestConfig: config,
	})
	if err != nil {
		return nil, err
	}

	return &Context{
		ClientConfig: config,
		ClientSet:    clientSet,
		HelmClient:   helmChartClient,
		Namespace:    namespace,
		HelmChart:    nil,
	}, nil
}

func NewFrameworkWithDefaultNamespace() (*Context, error) {
	return NewFrameworkWithNamespace(defaultNS)
}

func generateRandomAlphaNumeric6Digit() string {
	ret := make([]byte, hashLength)
	for i := 0; i < hashLength; i++ {
		ret[i] = hashSeed[rand.Intn(len(hashSeed))]
	}
	return string(ret)
}

func DeployHelmChart(relName, chartPath, values string) func() {
	return func() {
		helmChartSpec := &helmclient.ChartSpec{
			ReleaseName:     relName + generateRandomAlphaNumeric6Digit(),
			ChartName:       chartPath,
			Namespace:       ctx.Namespace,
			CreateNamespace: false,
			Wait:            false,
			Timeout:         5 * time.Minute,
			CleanupOnFail:   false,
			ValuesYaml:      values,
		}
		ctx.HelmChart = helmChartSpec

		_, err := ctx.HelmClient.InstallChart(context.TODO(), ctx.HelmChart, nil)
		Expect(err).To(BeNil())
	}
}

func DeleteHelmChart() func() {
	return func() {
		err := ctx.HelmClient.UninstallRelease(ctx.HelmChart)
		Expect(err).To(BeNil())
	}
}
