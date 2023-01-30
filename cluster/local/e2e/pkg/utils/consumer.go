package utils

import (
	"context"
	"time"

	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/edgefarm/provider-nats/apis/consumer/v1alpha1"
)

func (f *Framework) GetConsumer(consumer string) (*v1alpha1.Consumer, error) {
	out, err := f.ConsumersClientset.Consumers().Get(context.Background(), consumer, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (f *Framework) ListConsumer() (*v1alpha1.ConsumerList, error) {
	consumers, err := f.ConsumersClientset.Consumers().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return consumers, nil
}

func (f *Framework) CreateConsumer(consumer *v1alpha1.Consumer) error {
	_, err := f.ConsumersClientset.Consumers().Create(context.Background(), consumer, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (f *Framework) UpdateConsumer(consumer *v1alpha1.Consumer) error {
	_, err := f.ConsumersClientset.Consumers().Update(context.Background(), consumer, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (f *Framework) DeleteConsumer(consumer string) error {
	err := f.ConsumersClientset.Consumers().Delete(context.Background(), consumer, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (f *Framework) WaitForConsumerDeleted(name string) error {
	return wait.PollImmediate(time.Second, PollTimeout, func() (bool, error) {
		_, err := f.ConsumersClientset.Consumers().Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			if kerrors.IsNotFound(err) {
				return true, nil
			}
			return false, err
		}
		return false, nil
	})
}

func (f *Framework) WaitForConsumerConditionsAvailable(name string, amount int) error {
	return wait.PollImmediate(time.Second, PollTimeout, func() (bool, error) {
		consumer, err := f.ConsumersClientset.Consumers().Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		status := consumer.Status.Conditions
		if len(status) < amount {
			return false, nil
		}
		return true, nil
	})
}

func (f *Framework) WaitForConsumerCondition(name string, condition *v1.Condition) error {
	return wait.PollImmediate(time.Second, PollTimeout, func() (bool, error) {
		consumer, err := f.ConsumersClientset.Consumers().Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		status := consumer.Status.Conditions
		for _, cond := range status {
			if cond.Type == condition.Type &&
				cond.Status == condition.Status &&
				cond.Reason == condition.Reason &&
				cond.Message == condition.Message {
				return true, nil
			}
		}

		return false, nil
	})
}

func (f *Framework) WaitForConsumerSyncAndReady(name string) error {
	return wait.PollImmediate(time.Second, PollTimeout, func() (bool, error) {
		consumer, err := f.ConsumersClientset.Consumers().Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		status := consumer.Status.Conditions
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

func (f *Framework) GetConsumerParameters(name string) (*v1alpha1.ConsumerParameters, error) {
	consumer, err := f.ConsumersClientset.Consumers().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &consumer.Spec.ForProvider, nil
}
