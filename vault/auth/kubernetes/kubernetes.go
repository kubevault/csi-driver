package kubernetes

import (
	"fmt"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	appcat_cs "kmodules.xyz/custom-resources/client/clientset/versioned/typed/appcatalog/v1alpha1"
)

func getServiceAccountSecret(kc *kubernetes.Clientset, svcName, svcNamespace string) (string, error) {
	serviceAccount, err := kc.CoreV1().ServiceAccounts(svcNamespace).Get(svcName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if len(serviceAccount.Secrets) == 0 {
		return "", errors.Errorf("No service account secret found")
	}
	secretName := serviceAccount.Secrets[0].Name
	fmt.Println(secretName)

	secret, err := kc.CoreV1().Secrets(svcNamespace).Get(secretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return secret.Name, nil
}

func getAppBinding(appName, appNamespace string) (*appcat.AppBinding, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	appClient, err := appcat_cs.NewForConfig(config)

	app, err := appClient.AppBindings(appNamespace).Get(appName, metav1.GetOptions{})
	return app, err
}

func getKubeClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return kubeClient, nil
}
