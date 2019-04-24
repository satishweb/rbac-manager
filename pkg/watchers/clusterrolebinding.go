/*
Copyright 2019 ReactiveOps.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package watchers

import (
	rbacdef "github.com/reactiveops/rbac-manager/pkg/controller/rbacdefinition"
	kube "github.com/reactiveops/rbac-manager/pkg/kube"
	"github.com/sirupsen/logrus"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

func watchClusterRoleBindings(clientset *kubernetes.Clientset) error {
	watcher, err := clientset.RbacV1().ClusterRoleBindings().Watch(kube.ListOptions)

	if err != nil {
		logrus.Error(err, "unable to watch Cluster Role Bindings")
		return err
	}

	ch := watcher.ResultChan()

	for event := range ch {
		crb, ok := event.Object.(*rbacv1.ClusterRoleBinding)
		if !ok {
			logrus.Error(err, "Could not parse Cluster Role Binding")
		}

		if event.Type == watch.Modified || event.Type == watch.Deleted {
			rbacdef.ReconcileOwners(crb.OwnerReferences, "ClusterRoleBinding")
		}
	}
	return nil
}