package utils

import (
	"context"
	"fmt"
	"time"

	batchv1 "k8s.io/api/batch/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

type Pub struct {
	JWT     string
	Seed    string
	Address string
	Subject string
	Count   string
	Data    string
}

type StreamInfo struct {
	JWT      string
	Seed     string
	Address  string
	Count    string
	Stream   string
	Timeout  string
	Interval string
	Domain   string
}

func (f *Framework) NatsPublish(p *Pub) error {
	namespace := "nats"

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "generate",
			Namespace: namespace,
		},
		Data: map[string]string{
			"generate.sh": `#!/bin/sh
JWT=$JWT; SEED=$SEED
mkdir -p /creds
cat > /creds/creds <<EOF
-----BEGIN NATS USER JWT-----
$JWT
------END NATS USER JWT------
-----BEGIN USER NKEY SEED-----
$SEED
------END USER NKEY SEED------'
EOF`},
	}

	err := f.ApplyOrUpdate(cm)
	if err != nil {
		return err
	}

	// Define the Job object
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "publish-job",
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "publish",
							Image:   "natsio/nats-box:0.13.3",
							Command: []string{"/bin/sh", "-c", "cp /generate/generate.sh /generate.sh && chmod +x /generate.sh && /generate.sh && nats publish -s $ADDRESS --creds /creds/creds $SUBJECT --count $COUNT $DATA"},
							Env: []corev1.EnvVar{
								{
									Name:  "JWT",
									Value: p.JWT,
								},
								{
									Name:  "SEED",
									Value: p.Seed,
								},
								{
									Name:  "ADDRESS",
									Value: p.Address,
								},
								{
									Name:  "SUBJECT",
									Value: p.Subject,
								},
								{
									Name:  "COUNT",
									Value: p.Count,
								},
								{
									Name:  "DATA",
									Value: p.Data,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "generate",
									MountPath: "/generate/",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "generate",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "generate",
									},
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
			TTLSecondsAfterFinished: func() *int32 { i := int32(5); return &i }(),
			BackoffLimit:            func() *int32 { i := int32(4); return &i }(),
		},
	}

	err = f.ApplyOrUpdate(job)
	if err != nil {
		return err
	}
	err = f.WaitForJobFinished(job.Name, namespace)
	if err != nil {
		return err
	}
	err = f.DeleteObj(job)
	if err != nil {
		return err
	}
	err = f.DeleteObj(cm)
	if err != nil {
		return err
	}
	err = f.WairForDeleted(cm)
	if err != nil {
		return err
	}

	return nil
}

func (f *Framework) NatsCheckMessages(p *StreamInfo) error {
	namespace := "nats"

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "generate",
			Namespace: namespace,
		},
		Data: map[string]string{
			"generate.sh": `#!/bin/sh
JWT=$JWT; SEED=$SEED
mkdir -p /creds
cat > /creds/creds <<EOF
-----BEGIN NATS USER JWT-----
$JWT
------END NATS USER JWT------
-----BEGIN USER NKEY SEED-----
$SEED
------END USER NKEY SEED------'
EOF`,
			"check.sh": `#!/bin/sh
CMD=${CMD}
COUNT=${COUNT}
DOMAIN=${DOMAIN}
STREAM=${STREAM}
ADDRESS=${ADDRESS}
start_t=$(date +%s)
while true; do
	stop_t=$(date +%s)
	if [ $(($stop_t - $start_t)) -gt $TIMEOUT ]; then
		exit 1
	fi
	output=$(eval $CMD)
	if [ "$output" == "$COUNT" ]; then
		exit 0
	fi
	sleep $INTERVAL
done`},
	}

	err := f.ApplyOrUpdate(cm)
	if err != nil {
		return err
	}

	// Define the Pod object
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "check-messages",
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "check",
					Image:   "natsio/nats-box:0.13.3",
					Command: []string{"/bin/sh", "-c", "cp /generate/generate.sh /generate.sh && chmod +x /generate.sh && /generate.sh && cp /generate/check.sh /check.sh && chmod +x /check.sh && /check.sh"},
					Env: []corev1.EnvVar{
						{
							Name:  "CMD",
							Value: "nats stream info ${STREAM} -s ${ADDRESS} --domain \"${DOMAIN}\" --creds /creds/creds -j | jq -M '.state.messages'",
						},
						{
							Name:  "TIMEOUT",
							Value: p.Timeout,
						},
						{
							Name:  "INTERVAL",
							Value: p.Interval,
						},
						{
							Name:  "JWT",
							Value: p.JWT,
						},
						{
							Name:  "SEED",
							Value: p.Seed,
						},
						{
							Name:  "ADDRESS",
							Value: p.Address,
						},
						{
							Name:  "STREAM",
							Value: p.Stream,
						},
						{
							Name:  "COUNT",
							Value: p.Count,
						},
						{
							Name:  "DOMAIN",
							Value: p.Domain,
						}},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "generate",
							MountPath: "/generate/",
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "generate",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "generate",
							},
						},
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}

	err = f.ApplyOrUpdate(pod)
	if err != nil {
		return err
	}
	err = f.WaitForPodFinished(pod.Name, namespace)
	if err != nil {
		return err
	}
	err = f.DeleteObj(pod)
	if err != nil {
		return err
	}
	err = f.DeleteObj(cm)
	if err != nil {
		return err
	}
	err = f.WairForDeleted(cm)
	if err != nil {
		return err
	}

	return nil
}

func (f *Framework) WaitForJobFinished(name string, namespace string) error {
	return wait.PollImmediate(time.Second, PollTimeout, func() (bool, error) {
		obj, err := f.ClientSet.BatchV1().Jobs(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if obj.Status.CompletionTime != nil {
			return true, nil
		}
		return false, nil
	})
}

func (f *Framework) WaitForPodFinished(name string, namespace string) error {
	return wait.PollImmediate(time.Second, PollTimeout, func() (bool, error) {
		obj, err := f.ClientSet.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if obj.Status.Phase == corev1.PodSucceeded {
			return true, nil
		} else if obj.Status.Phase == corev1.PodFailed {
			return true, fmt.Errorf("%s", obj.Status.ContainerStatuses[0].State.Terminated.Reason)
		}
		return false, nil
	})
}
