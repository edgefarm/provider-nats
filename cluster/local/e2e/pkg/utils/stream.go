package utils

import (
	"context"
	"time"

	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/edgefarm/provider-nats/apis/stream/v1alpha1"
)

func (f *Framework) GetStream(stream string) (*v1alpha1.Stream, error) {
	out, err := f.StreamsClientset.Streams().Get(context.Background(), stream, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (f *Framework) ListStream() (*v1alpha1.StreamList, error) {
	streams, err := f.StreamsClientset.Streams().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return streams, nil
}

func (f *Framework) CreateStream(stream *v1alpha1.Stream) error {
	_, err := f.StreamsClientset.Streams().Create(context.Background(), stream, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (f *Framework) UpdateStream(stream *v1alpha1.Stream) error {
	_, err := f.StreamsClientset.Streams().Update(context.Background(), stream, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (f *Framework) DeleteStream(stream string) error {
	err := f.StreamsClientset.Streams().Delete(context.Background(), stream, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (f *Framework) WaitForStreamDeleted(name string) error {
	return wait.PollImmediate(time.Second, PollTimeout, func() (bool, error) {
		_, err := f.StreamsClientset.Streams().Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			if kerrors.IsNotFound(err) {
				return true, nil
			}
			return false, err
		}
		return false, nil
	})
}

func (f *Framework) WaitForStreamConditionsAvailable(name string, amount int) error {
	return wait.PollImmediate(time.Second, PollTimeout, func() (bool, error) {
		stream, err := f.StreamsClientset.Streams().Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			if kerrors.IsNotFound(err) {
				return false, nil
			}
			return false, err
		}
		status := stream.Status.Conditions
		if len(status) < amount {
			return false, nil
		}
		return true, nil
	})
}

func (f *Framework) WaitForStreamSyncAndReady(name string) error {
	return wait.PollImmediate(time.Second, PollTimeout, func() (bool, error) {
		stream, err := f.StreamsClientset.Streams().Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			if kerrors.IsNotFound(err) {
				return false, nil
			}
			return false, err
		}
		status := stream.Status.Conditions
		ready := false
		synced := false

		for _, condition := range status {
			if condition.Type == v1.TypeReady && condition.Status == corev1.ConditionTrue {
				ready = true
				continue
			}
			if condition.Type == v1.TypeSynced && condition.Status == corev1.ConditionTrue {
				synced = true
				continue
			}
		}

		if !ready || !synced {
			return false, nil
		}
		return true, nil
	})
}

func (f *Framework) WaitForStreamMessags(name string, messages uint64) error {
	return wait.PollImmediate(time.Second, PollTimeout, func() (bool, error) {
		stream, err := f.StreamsClientset.Streams().Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			if kerrors.IsNotFound(err) {
				return false, nil
			}
			return false, err
		}
		status := stream.Status.AtProvider.State.Messages
		if status < messages {
			return false, nil
		}
		return true, nil
	})
}

func (f *Framework) GetStreamParameters(name string) (*v1alpha1.StreamParameters, error) {
	stream, err := f.StreamsClientset.Streams().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &stream.Spec.ForProvider, nil
}

func (f *Framework) GetStreamUserPublicKey(name string) (string, error) {
	account, err := f.StreamsClientset.Streams().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return account.Status.AtProvider.Connection.PublicKey, nil
}

func (f *Framework) GetStreamDomain(name string) (string, error) {
	account, err := f.StreamsClientset.Streams().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return account.Status.AtProvider.Domain, nil
}
