/*
Copyright 2016 The Kubernetes Authors.

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

package gcp

import (
	"fmt"

	"k8s.io/kubernetes/test/e2e/cloud/gcp/common"
	"k8s.io/kubernetes/test/e2e/framework"
	e2eskipper "k8s.io/kubernetes/test/e2e/framework/skipper"
	"k8s.io/kubernetes/test/e2e/upgrades"
	"k8s.io/kubernetes/test/e2e/upgrades/apps"
	"k8s.io/kubernetes/test/e2e/upgrades/auth"
	"k8s.io/kubernetes/test/e2e/upgrades/autoscaling"
	"k8s.io/kubernetes/test/e2e/upgrades/network"
	"k8s.io/kubernetes/test/e2e/upgrades/node"
	"k8s.io/kubernetes/test/e2e/upgrades/storage"
	"k8s.io/kubernetes/test/utils/junit"

	"github.com/onsi/ginkgo"
)

var upgradeTests = []upgrades.Test{
	&network.ServiceUpgradeTest{},
	&node.SecretUpgradeTest{},
	&apps.ReplicaSetUpgradeTest{},
	&apps.StatefulSetUpgradeTest{},
	&apps.DeploymentUpgradeTest{},
	&apps.JobUpgradeTest{},
	&node.ConfigMapUpgradeTest{},
	&autoscaling.HPAUpgradeTest{},
	&storage.PersistentVolumeUpgradeTest{},
	&apps.DaemonSetUpgradeTest{},
	&node.AppArmorUpgradeTest{},
	&storage.VolumeModeDowngradeTest{},
}

var statefulsetUpgradeTests = []upgrades.Test{
	&apps.MySQLUpgradeTest{},
	&apps.EtcdUpgradeTest{},
	&apps.CassandraUpgradeTest{},
}

var kubeProxyUpgradeTests = []upgrades.Test{
	&network.KubeProxyUpgradeTest{},
	&network.ServiceUpgradeTest{},
}

var kubeProxyDowngradeTests = []upgrades.Test{
	&network.KubeProxyDowngradeTest{},
	&network.ServiceUpgradeTest{},
}

var serviceaccountAdmissionControllerMigrationTests = []upgrades.Test{
	&auth.ServiceAccountAdmissionControllerMigrationTest{},
}

func kubeProxyDaemonSetExtraEnvs(enableKubeProxyDaemonSet bool) []string {
	return []string{fmt.Sprintf("KUBE_PROXY_DAEMONSET=%v", enableKubeProxyDaemonSet)}
}

// TODO(#98326): Split the test by SIGs, move to appropriate directories and use SIGDescribe.
var _ = ginkgo.Describe("Upgrade [Feature:Upgrade]", func() {
	f := framework.NewDefaultFramework("cluster-upgrade")

	// Create the frameworks here because we can only create them
	// in a "Describe".
	ginkgo.Describe("master upgrade", func() {
		ginkgo.It("should maintain a functioning cluster [Feature:MasterUpgrade]", func() {
			upgCtx, err := common.GetUpgradeContext(f.ClientSet.Discovery())
			framework.ExpectNoError(err)

			testSuite := &junit.TestSuite{Name: "Master upgrade"}
			masterUpgradeTest := &junit.TestCase{
				Name:      "[sig-cloud-provider-gcp] master-upgrade",
				Classname: "upgrade_tests",
			}
			testSuite.TestCases = append(testSuite.TestCases, masterUpgradeTest, nil)

			upgradeFunc := common.ControlPlaneUpgradeFunc(f, upgCtx, masterUpgradeTest, nil)
			upgrades.RunUpgradeSuite(upgCtx, upgradeTests, testSuite, upgrades.MasterUpgrade, upgradeFunc)
		})
	})

	ginkgo.Describe("cluster upgrade", func() {
		ginkgo.It("should maintain a functioning cluster [Feature:ClusterUpgrade]", func() {
			upgCtx, err := common.GetUpgradeContext(f.ClientSet.Discovery())
			framework.ExpectNoError(err)

			testSuite := &junit.TestSuite{Name: "Cluster upgrade"}
			clusterUpgradeTest := &junit.TestCase{Name: "[sig-cloud-provider-gcp] cluster-upgrade", Classname: "upgrade_tests"}
			testSuite.TestCases = append(testSuite.TestCases, clusterUpgradeTest)

			upgradeFunc := common.ClusterUpgradeFunc(f, upgCtx, clusterUpgradeTest, nil, nil)
			upgrades.RunUpgradeSuite(upgCtx, upgradeTests, testSuite, upgrades.ClusterUpgrade, upgradeFunc)
		})
	})
})

// TODO(#98326): Split the test by SIGs, move to appropriate directories and use SIGDescribe.
var _ = ginkgo.Describe("Downgrade [Feature:Downgrade]", func() {
	f := framework.NewDefaultFramework("cluster-downgrade")

	ginkgo.Describe("cluster downgrade", func() {
		ginkgo.It("should maintain a functioning cluster [Feature:ClusterDowngrade]", func() {
			upgCtx, err := common.GetUpgradeContext(f.ClientSet.Discovery())
			framework.ExpectNoError(err)

			testSuite := &junit.TestSuite{Name: "Cluster downgrade"}
			clusterDowngradeTest := &junit.TestCase{Name: "[sig-cloud-provider-gcp] cluster-downgrade", Classname: "upgrade_tests"}
			testSuite.TestCases = append(testSuite.TestCases, clusterDowngradeTest)

			upgradeFunc := common.ClusterDowngradeFunc(f, upgCtx, clusterDowngradeTest, nil, nil)
			upgrades.RunUpgradeSuite(upgCtx, upgradeTests, testSuite, upgrades.ClusterUpgrade, upgradeFunc)
		})
	})
})

// TODO(#98326): Move the test to sig-aps dir and use SIGDescribe.
var _ = ginkgo.Describe("[sig-apps] stateful Upgrade [Feature:StatefulUpgrade]", func() {
	f := framework.NewDefaultFramework("stateful-upgrade")

	ginkgo.Describe("stateful upgrade", func() {
		ginkgo.It("should maintain a functioning cluster", func() {
			upgCtx, err := common.GetUpgradeContext(f.ClientSet.Discovery())
			framework.ExpectNoError(err)

			testSuite := &junit.TestSuite{Name: "Stateful upgrade"}
			statefulUpgradeTest := &junit.TestCase{Name: "[sig-apps] stateful-upgrade", Classname: "upgrade_tests"}
			testSuite.TestCases = append(testSuite.TestCases, statefulUpgradeTest)

			upgradeFunc := common.ClusterUpgradeFunc(f, upgCtx, statefulUpgradeTest, nil, nil)
			upgrades.RunUpgradeSuite(upgCtx, statefulsetUpgradeTests, testSuite, upgrades.ClusterUpgrade, upgradeFunc)
		})
	})
})

// TODO(#98326): Move the test to sig-network dir and use SIGDescribe.
var _ = ginkgo.Describe("kube-proxy migration [Feature:KubeProxyDaemonSetMigration]", func() {
	f := framework.NewDefaultFramework("kube-proxy-ds-migration")

	ginkgo.BeforeEach(func() {
		e2eskipper.SkipUnlessProviderIs("gce")
	})

	ginkgo.Describe("Upgrade kube-proxy from static pods to a DaemonSet", func() {
		ginkgo.It("should maintain a functioning cluster [Feature:KubeProxyDaemonSetUpgrade]", func() {
			upgCtx, err := common.GetUpgradeContext(f.ClientSet.Discovery())
			framework.ExpectNoError(err)

			testSuite := &junit.TestSuite{Name: "kube-proxy upgrade"}
			kubeProxyUpgradeTest := &junit.TestCase{
				Name:      "kube-proxy-ds-upgrade",
				Classname: "upgrade_tests",
			}
			testSuite.TestCases = append(testSuite.TestCases, kubeProxyUpgradeTest)

			extraEnvs := kubeProxyDaemonSetExtraEnvs(true)
			upgradeFunc := common.ClusterUpgradeFunc(f, upgCtx, kubeProxyUpgradeTest, extraEnvs, extraEnvs)
			upgrades.RunUpgradeSuite(upgCtx, kubeProxyUpgradeTests, testSuite, upgrades.ClusterUpgrade, upgradeFunc)
		})
	})

	ginkgo.Describe("Downgrade kube-proxy from a DaemonSet to static pods", func() {
		ginkgo.It("should maintain a functioning cluster [Feature:KubeProxyDaemonSetDowngrade]", func() {
			upgCtx, err := common.GetUpgradeContext(f.ClientSet.Discovery())
			framework.ExpectNoError(err)

			testSuite := &junit.TestSuite{Name: "kube-proxy downgrade"}
			kubeProxyDowngradeTest := &junit.TestCase{
				Name:      "kube-proxy-ds-downgrade",
				Classname: "upgrade_tests",
			}
			testSuite.TestCases = append(testSuite.TestCases, kubeProxyDowngradeTest)

			extraEnvs := kubeProxyDaemonSetExtraEnvs(false)
			upgradeFunc := common.ClusterDowngradeFunc(f, upgCtx, kubeProxyDowngradeTest, extraEnvs, extraEnvs)
			upgrades.RunUpgradeSuite(upgCtx, kubeProxyDowngradeTests, testSuite, upgrades.ClusterUpgrade, upgradeFunc)
		})
	})
})

// TODO(#98326): Move the test to sig-auth dir and use SIGDescribe.
var _ = ginkgo.Describe("[sig-auth] ServiceAccount admission controller migration [Feature:BoundServiceAccountTokenVolume]", func() {
	f := framework.NewDefaultFramework("serviceaccount-admission-controller-migration")

	ginkgo.Describe("master upgrade", func() {
		ginkgo.It("should maintain a functioning cluster", func() {
			upgCtx, err := common.GetUpgradeContext(f.ClientSet.Discovery())
			framework.ExpectNoError(err)

			testSuite := &junit.TestSuite{Name: "ServiceAccount admission controller migration"}
			serviceaccountAdmissionControllerMigrationTest := &junit.TestCase{
				Name:      "[sig-auth] serviceaccount-admission-controller-migration",
				Classname: "upgrade_tests",
			}
			testSuite.TestCases = append(testSuite.TestCases, serviceaccountAdmissionControllerMigrationTest)

			extraEnvs := []string{"KUBE_FEATURE_GATES=BoundServiceAccountTokenVolume=true"}
			upgradeFunc := common.ControlPlaneUpgradeFunc(f, upgCtx, serviceaccountAdmissionControllerMigrationTest, extraEnvs)
			upgrades.RunUpgradeSuite(upgCtx, serviceaccountAdmissionControllerMigrationTests, testSuite, upgrades.MasterUpgrade, upgradeFunc)
		})
	})
})
