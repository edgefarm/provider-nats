package utils

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func (f *Framework) DeleteNamespace(ns string, waitUntilDeleted bool) error {
	err := f.ClientSet.CoreV1().Namespaces().Delete(f.Context, ns, metav1.DeleteOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if !waitUntilDeleted {
		return nil
	}
	return wait.PollImmediate(time.Second, PollTimeout, func() (bool, error) {
		_, err = f.ClientSet.CoreV1().Namespaces().Get(f.Context, ns, metav1.GetOptions{})
		if kerrors.IsNotFound(err) {
			return true, nil
		}
		return false, err
	})
}

func (f *Framework) CreateNamespace(ns string) error {
	_, err := f.ClientSet.CoreV1().Namespaces().Create(f.Context, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
		Name: ns,
	}}, metav1.CreateOptions{})
	return err
}
