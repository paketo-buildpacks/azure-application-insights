/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package appinsights_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/microsoft-azure/appinsights"
	"github.com/paketo-buildpacks/microsoft-azure/internal/common"
)

func testJava(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
	)

	context("Build", func() {
		var (
			ctx libcnb.BuildContext
		)

		it.Before(func() {
			var err error

			ctx.Buildpack.Path, err = ioutil.TempDir("", "appinsights-java-build-buildpack")
			Expect(err).NotTo(HaveOccurred())

			Expect(os.MkdirAll(filepath.Join(ctx.Buildpack.Path, "resources"), 0755)).To(Succeed())
			Expect(ioutil.WriteFile(filepath.Join(ctx.Buildpack.Path, "resources", "AI-Agent.xml"), []byte{}, 0644)).
				To(Succeed())

			ctx.Layers.Path, err = ioutil.TempDir("", "appinsights-java-build-layers")
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(os.RemoveAll(ctx.Buildpack.Path)).To(Succeed())
			Expect(os.RemoveAll(ctx.Layers.Path)).To(Succeed())
		})

		it("contributes Java agent", func() {
			dep := libpak.BuildpackDependency{
				URI:    "https://localhost/stub-azure-application-insights-agent.jar",
				SHA256: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			}
			dc := libpak.DependencyCache{CachePath: "testdata"}

			layer, err := ctx.Layers.Layer("test-layer")
			Expect(err).NotTo(HaveOccurred())

			layer, err = appinsights.NewJavaBuild(ctx.Buildpack.Path, dep, dc, &libcnb.BuildpackPlan{}).Contribute(layer)
			Expect(err).NotTo(HaveOccurred())

			Expect(layer.Launch).To(BeTrue())
			Expect(filepath.Join(layer.Path, "stub-azure-application-insights-agent.jar")).To(BeARegularFile())
			Expect(filepath.Join(layer.Path, "AI-Agent.xml")).To(BeARegularFile())
			Expect(layer.LaunchEnvironment[fmt.Sprintf("%s.default", appinsights.AgentPath)]).
				To(Equal(filepath.Join(layer.Path, "stub-azure-application-insights-agent.jar")))
		})

	})

	context("Launch", func() {
		var (
			l = appinsights.JavaLaunch{
				CredentialSource: common.MetadataServer,
			}
		)

		it.Before(func() {
			Expect(os.Setenv(appinsights.AgentPath, "test-path")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv(appinsights.AgentPath)).To(Succeed())
		})

		it("does not contribute if source is None", func() {
			l.CredentialSource = common.None

			Expect(l.Execute()).To(BeNil())
		})

		it("returns error if BPI_AZURE_APPLICATION_INSIGHTS_AGENT_PATH is not set", func() {
			Expect(os.Unsetenv(appinsights.AgentPath)).To(Succeed())

			_, err := l.Execute()
			Expect(err).To(MatchError("$BPI_AZURE_APPLICATION_INSIGHTS_AGENT_PATH must be set"))
		})

		it("contributes JAVA_TOOL_OPTIONS", func() {
			Expect(l.Execute()).To(Equal(map[string]string{"JAVA_TOOL_OPTIONS": "-javaagent:test-path"}))
		})

		context("existing $JAVA_TOOL_OPTIONS", func() {

			it.Before(func() {
				Expect(os.Setenv("JAVA_TOOL_OPTIONS", "test-java-tool-options")).To(Succeed())
			})

			it.After(func() {
				Expect(os.Unsetenv("JAVA_TOOL_OPTIONS")).To(Succeed())
			})

			it("contributes JAVA_TOOL_OPTIONS", func() {
				Expect(l.Execute()).To(Equal(map[string]string{
					"JAVA_TOOL_OPTIONS": "test-java-tool-options -javaagent:test-path",
				}))
			})
		})
	})
}
