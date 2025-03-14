// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package cluster_test

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	"gitlab.eng.vmware.com/core-build/ako-operator/controllers/cluster"
	"gitlab.eng.vmware.com/core-build/ako-operator/pkg/aviclient"
	"gitlab.eng.vmware.com/core-build/ako-operator/pkg/test/builder"
	testutil "gitlab.eng.vmware.com/core-build/ako-operator/pkg/test/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	clustereaddonv1alpha3 "sigs.k8s.io/cluster-api/exp/addons/api/v1alpha3"
	ctrl "sigs.k8s.io/controller-runtime"

	ctrlmgr "sigs.k8s.io/controller-runtime/pkg/manager"
)

// suite is used for unit and integration testing this controller.
var suite = builder.NewTestSuiteForController(
	func(mgr ctrlmgr.Manager) error {

		builder.FakeAvi = aviclient.NewFakeAviClient()

		if err := (&cluster.ClusterReconciler{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("Cluster"),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr); err != nil {
			return err
		}
		return nil
	},
	func(scheme *runtime.Scheme) (err error) {
		err = corev1.AddToScheme(scheme)
		if err != nil {
			return err
		}
		err = clusterv1.AddToScheme(scheme)
		if err != nil {
			return err
		}
		err = clustereaddonv1alpha3.AddToScheme(scheme)
		if err != nil {
			return err
		}
		return nil
	},
	filepath.Join(testutil.FindModuleDir("sigs.k8s.io/cluster-api"), "config", "crd", "bases"),
)

func TestController(t *testing.T) {
	suite.Register(t, "AKO Operator Cluster Controller", intgTests, unitTests)
}

var _ = BeforeSuite(suite.BeforeSuite)

var _ = AfterSuite(suite.AfterSuite)

func intgTests() {
	Describe("ClusterController Test", intgTestEnsureClusterHAProvider)
}

func unitTests() {
}
