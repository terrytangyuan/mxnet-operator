// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testutil

import (
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mxv1 "github.com/kubeflow/mxnet-operator/pkg/apis/mxnet/v1"
)

func NewMXJobWithCleanPolicy(scheduler, worker, server int, policy mxv1.CleanPodPolicy) *mxv1.MXJob {

	var mxJob *mxv1.MXJob

	if scheduler > 0 {
		mxJob = NewMXJobWithScheduler(worker, server)
	} else {
		mxJob = NewMXJob(worker, server)
	}

	mxJob.Spec.CleanPodPolicy = &policy
	return mxJob
}

func NewMXJobWithCleanupJobDelay(scheduler, worker, server int, ttl *int32) *mxv1.MXJob {

	var mxJob *mxv1.MXJob

	if scheduler > 0 {
		mxJob = NewMXJobWithScheduler(worker, server)
	} else {
		mxJob = NewMXJob(worker, server)
	}

	mxJob.Spec.TTLSecondsAfterFinished = ttl
	policy := mxv1.CleanPodPolicyNone
	mxJob.Spec.CleanPodPolicy = &policy
	return mxJob
}

func NewMXJobWithActiveDeadlineSeconds(scheduler, worker, ps int, ads *int64) *mxv1.MXJob {
	if scheduler == 1 {
		mxJob := NewMXJobWithScheduler(worker, ps)
		mxJob.Spec.ActiveDeadlineSeconds = ads
		policy := mxv1.CleanPodPolicyAll
		mxJob.Spec.CleanPodPolicy = &policy
		return mxJob
	}
	mxJob := NewMXJob(worker, ps)
	mxJob.Spec.ActiveDeadlineSeconds = ads
	policy := mxv1.CleanPodPolicyAll
	mxJob.Spec.CleanPodPolicy = &policy
	return mxJob
}

func NewMXJobWithBackoffLimit(scheduler, worker, ps int, backoffLimit *int32) *mxv1.MXJob {
	if scheduler == 1 {
		mxJob := NewMXJobWithScheduler(worker, ps)
		mxJob.Spec.BackoffLimit = backoffLimit
		mxJob.Spec.MXReplicaSpecs["Worker"].RestartPolicy = "OnFailure"
		policy := mxv1.CleanPodPolicyAll
		mxJob.Spec.CleanPodPolicy = &policy
		return mxJob
	}
	mxJob := NewMXJob(worker, ps)
	mxJob.Spec.BackoffLimit = backoffLimit
	mxJob.Spec.MXReplicaSpecs["Worker"].RestartPolicy = "OnFailure"
	policy := mxv1.CleanPodPolicyAll
	mxJob.Spec.CleanPodPolicy = &policy
	return mxJob
}

func NewMXJobWithScheduler(worker, server int) *mxv1.MXJob {
	mxJob := NewMXJob(worker, server)
	mxJob.Spec.MXReplicaSpecs[mxv1.MXReplicaTypeScheduler] = &mxv1.MXReplicaSpec{
		Template: NewMXReplicaSpecTemplate(),
	}
	return mxJob
}

func NewMXJob(worker, server int) *mxv1.MXJob {
	mxJob := &mxv1.MXJob{
		TypeMeta: metav1.TypeMeta{
			Kind: mxv1.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      TestMXJobName,
			Namespace: metav1.NamespaceDefault,
		},
		Spec: mxv1.MXJobSpec{
			MXReplicaSpecs: make(map[mxv1.MXReplicaType]*mxv1.MXReplicaSpec),
		},
	}

	if worker > 0 {
		worker := int32(worker)
		workerReplicaSpec := &mxv1.MXReplicaSpec{
			Replicas: &worker,
			Template: NewMXReplicaSpecTemplate(),
		}
		mxJob.Spec.MXReplicaSpecs[mxv1.MXReplicaTypeWorker] = workerReplicaSpec
	}

	if server > 0 {
		server := int32(server)
		serverReplicaSpec := &mxv1.MXReplicaSpec{
			Replicas: &server,
			Template: NewMXReplicaSpecTemplate(),
		}
		mxJob.Spec.MXReplicaSpecs[mxv1.MXReplicaTypeServer] = serverReplicaSpec
	}
	return mxJob
}

func NewMXReplicaSpecTemplate() v1.PodTemplateSpec {
	return v1.PodTemplateSpec{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  mxv1.DefaultContainerName,
					Image: TestImageName,
					Args:  []string{"Fake", "Fake"},
					Ports: []v1.ContainerPort{
						{
							Name:          mxv1.DefaultPortName,
							ContainerPort: mxv1.DefaultPort,
						},
					},
				},
			},
		},
	}
}

func SetMXJobCompletionTime(mxJob *mxv1.MXJob) {
	now := metav1.Time{Time: time.Now()}
	mxJob.Status.CompletionTime = &now
}
