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

package azure_test

import (
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"

	"github.com/paketo-buildpacks/libpak"
	azure "github.com/paketo-buildpacks/microsoft-azure"
	"github.com/paketo-buildpacks/microsoft-azure/internal/common"

	"github.com/sclevine/spec"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx libcnb.BuildContext
	)

	it("contributes credentials", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: common.Credentials})

		result, err := azure.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(1))
		Expect(result.Layers[0].Name()).To(Equal("helper"))
		Expect(result.Layers[0].(libpak.HelperLayerContributor).Names).To(Equal([]string{common.Credentials}))
	})

	it("contributes appinsights Java", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: common.ApplicationInsightsJava})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      common.ApplicationInsightsJava,
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
				},
			},
		}
		ctx.StackID = "test-stack-id"

		result, err := azure.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(2))
		Expect(result.Layers[0].Name()).To(Equal(common.ApplicationInsightsJava))
		Expect(result.Layers[1].Name()).To(Equal("helper"))
		Expect(result.Layers[1].(libpak.HelperLayerContributor).Names).To(Equal([]string{common.ApplicationInsightsJava}))
	})

	it("contributes appinsights NodeJS", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: common.ApplicationInsightsNodeJS})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      common.ApplicationInsightsNodeJS,
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
				},
			},
		}
		ctx.StackID = "test-stack-id"

		result, err := azure.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(2))
		Expect(result.Layers[0].Name()).To(Equal(common.ApplicationInsightsNodeJS))
		Expect(result.Layers[1].Name()).To(Equal("helper"))
		Expect(result.Layers[1].(libpak.HelperLayerContributor).Names).To(Equal([]string{common.ApplicationInsightsNodeJS}))
	})
}
