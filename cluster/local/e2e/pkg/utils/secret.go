package utils

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetSecret from Kubernetes
func (f *Framework) GetSecret(nameSpace string, secretName string) (map[string][]byte, error) {
	secret, err := f.ClientSet.CoreV1().Secrets(nameSpace).Get(f.Context, secretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return secret.Data, nil
}
