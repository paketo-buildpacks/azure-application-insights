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

package insights_test

import (
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/azure-application-insights/v5/insights"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx libcnb.BuildContext
	)

	it("contributes Java agent for API <=0.6", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "azure-application-insights-java"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "azure-application-insights-java",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
				},
			},
		}
		ctx.StackID = "test-stack-id"
		ctx.Buildpack.API = "0.6"

		result, err := insights.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(2))
		Expect(result.Layers[0].Name()).To(Equal("azure-application-insights-java"))
		Expect(result.Layers[1].Name()).To(Equal("helper"))
		Expect(result.Layers[1].(libpak.HelperLayerContributor).Names).To(Equal([]string{"properties"}))

		Expect(result.BOM.Entries).To(HaveLen(2))
		Expect(result.BOM.Entries[0].Name).To(Equal("azure-application-insights-java"))
		Expect(result.BOM.Entries[1].Name).To(Equal("helper"))
	})

	it("contributes Java agent for API 0.7+", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "azure-application-insights-java"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "azure-application-insights-java",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
					"cpes":    []string{"cpe:2.3:a:microsoft:azure-application-insights-java:3.2.3:*:*:*:*:*:*:*"},
					"purl":    "pkg:generic/azure-application-insights-java@3.2.3",
				},
			},
		}
		ctx.StackID = "test-stack-id"
		ctx.Buildpack.API = "0.7"

		result, err := insights.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(2))
		Expect(result.Layers[0].Name()).To(Equal("azure-application-insights-java"))
		Expect(result.Layers[1].Name()).To(Equal("helper"))
		Expect(result.Layers[1].(libpak.HelperLayerContributor).Names).To(Equal([]string{"properties"}))

		Expect(result.BOM.Entries).To(HaveLen(2))
		Expect(result.BOM.Entries[0].Name).To(Equal("azure-application-insights-java"))
		Expect(result.BOM.Entries[1].Name).To(Equal("helper"))
	})

	it("contributes NodeJS agent for API <=0.6", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "azure-application-insights-nodejs"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "azure-application-insights-nodejs",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
				},
			},
		}
		ctx.StackID = "test-stack-id"
		ctx.Buildpack.API = "0.6"

		result, err := insights.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(2))
		Expect(result.Layers[0].Name()).To(Equal("azure-application-insights-nodejs"))
		Expect(result.Layers[1].Name()).To(Equal("helper"))
		Expect(result.Layers[1].(libpak.HelperLayerContributor).Names).To(Equal([]string{"properties"}))

		Expect(result.BOM.Entries).To(HaveLen(2))
		Expect(result.BOM.Entries[0].Name).To(Equal("azure-application-insights-nodejs"))
		Expect(result.BOM.Entries[1].Name).To(Equal("helper"))
	})
	it("contributes NodeJS agent for API 0.7+", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "azure-application-insights-nodejs"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "azure-application-insights-nodejs",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
					"cpes":    []string{"cpe:2.3:a:microsoft:azure-application-insights-nodejs:3.2.3:*:*:*:*:*:*:*"},
					"purl":    "pkg:generic/azure-application-insights-nodejs@3.2.3",
				},
			},
		}
		ctx.StackID = "test-stack-id"
		ctx.Buildpack.API = "0.7"

		result, err := insights.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(2))
		Expect(result.Layers[0].Name()).To(Equal("azure-application-insights-nodejs"))
		Expect(result.Layers[1].Name()).To(Equal("helper"))
		Expect(result.Layers[1].(libpak.HelperLayerContributor).Names).To(Equal([]string{"properties"}))

		Expect(result.BOM.Entries).To(HaveLen(2))
		Expect(result.BOM.Entries[0].Name).To(Equal("azure-application-insights-nodejs"))
		Expect(result.BOM.Entries[1].Name).To(Equal("helper"))
	})
}
