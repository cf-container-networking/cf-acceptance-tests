package apps

import (
	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"
	"github.com/cloudfoundry-incubator/cf-test-helpers/helpers"
	"github.com/cloudfoundry/cf-acceptance-tests/helpers/app_helpers"
	"github.com/cloudfoundry/cf-acceptance-tests/helpers/assets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("when an app has multiple instances", func() {
	var appName string
	BeforeEach(func() {
		appName = generator.PrefixedRandomName("CATS-APP-")
		Eventually(cf.Cf(
			"push", appName,
			"-p", assets.NewAssets().Dora,
			"--no-start",
			"-b", "ruby_buildpack",
			"-m", DEFAULT_MEMORY_LIMIT,
			"-d", config.AppsDomain,
			"-i", "2"),
			DEFAULT_TIMEOUT,
		).Should(Exit(0))

		app_helpers.SetBackend(appName)

		Eventually(cf.Cf("start", appName), CF_PUSH_TIMEOUT).Should(Exit(0))
	})

	AfterEach(func() {
		app_helpers.AppReport(appName, DEFAULT_TIMEOUT)
		Eventually(cf.Cf("delete", appName, "-f"), DEFAULT_TIMEOUT).Should(Exit(0))
	})

	Describe("when a specific instance is targeted", func() {
		It("returns correct instance index", func() {
			Eventually(func() string {
				return helpers.CurlApp(appName, "/env/INSTANCE_INDEX")
			}, DEFAULT_TIMEOUT).Should(Equal("1"))
		})
	})
})
