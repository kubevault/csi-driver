package kubernetes

import (
	"fmt"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetJWT(serviceAccountName, namespace string) (string, error) {
	kubeClient, err := getKubeClient()
	if err != nil {
		return "", err
	}

	serviceAccount, err := kubeClient.CoreV1().ServiceAccounts(namespace).Get(serviceAccountName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if len(serviceAccount.Secrets) == 0 {
		return "", errors.Errorf("No service account secret found")
	}
	secretName := serviceAccount.Secrets[0].Name
	fmt.Println(secretName)

	secret, err := kubeClient.CoreV1().Secrets(namespace).Get(secretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	tokenData, ok := secret.Data["token"]
	if !ok {
		return "", errors.New("No jwt token found")
	}

	return string(tokenData), nil
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
