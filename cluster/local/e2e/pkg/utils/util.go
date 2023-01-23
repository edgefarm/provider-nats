package utils

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"

	"k8s.io/utils/pointer"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func (f *Framework) GetDefaultSecurityContext() *corev1.SecurityContext {
	return &corev1.SecurityContext{
		RunAsUser: pointer.Int64(12345),
	}
}

// WairForDeleted waits for the given object to be deleted.
func (f *Framework) WairForDeleted(obj ctrlclient.Object) error {
	return wait.PollImmediate(time.Second, PollTimeout, func() (bool, error) {
		err := f.CtrlClient.Get(f.Context, ctrlclient.ObjectKeyFromObject(obj), obj)
		if err != nil {
			if kerrors.IsNotFound(err) {
				return true, nil
			}
			return false, err
		}
		return false, nil
	})
}

// ApplyOrUpdate applies or updates the given object.
func (f *Framework) ApplyOrUpdate(obj ctrlclient.Object) error {
	err := f.CtrlClient.Create(f.Context, obj)
	if err != nil {
		if !kerrors.IsAlreadyExists(err) {
			return err
		}
		err = f.CtrlClient.Update(f.Context, obj)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteObj deletes or updates the given object.
func (f *Framework) DeleteObj(obj ctrlclient.Object) error {
	err := f.CtrlClient.Delete(f.Context, obj, &ctrlclient.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
