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
	"os"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	azure "github.com/paketo-buildpacks/microsoft-azure"
	"github.com/paketo-buildpacks/microsoft-azure/internal/common"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx    libcnb.DetectContext
		detect azure.Detect
	)

	it("fails without environment variable", func() {
		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{}))
	})

	context("$BP_AZURE_APPLICATION_INSIGHTS_ENABLED", func() {
		it.Before(func() {
			Expect(os.Setenv(azure.ApplicationInsightsEnabled, "true")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv(azure.ApplicationInsightsEnabled)).To(Succeed())
		})

		it("passes", func() {
			Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
				Pass: true,
				Plans: []libcnb.BuildPlan{
					{
						Provides: []libcnb.BuildPlanProvide{
							{Name: common.ApplicationInsightsJava},
							{Name: common.ApplicationInsightsNodeJS},
							{Name: common.Credentials},
						},
						Requires: []libcnb.BuildPlanRequire{
							{Name: common.ApplicationInsightsJava},
							{Name: "jvm-application"},
							{Name: common.ApplicationInsightsNodeJS},
							{Name: "node", Metadata: map[string]interface{}{"build": true}},
							{Name: common.Credentials},
						},
					},
					{
						Provides: []libcnb.BuildPlanProvide{
							{Name: common.ApplicationInsightsJava},
							{Name: common.Credentials},
						},
						Requires: []libcnb.BuildPlanRequire{
							{Name: common.ApplicationInsightsJava},
							{Name: "jvm-application"},
							{Name: common.Credentials},
						},
					},
					{
						Provides: []libcnb.BuildPlanProvide{
							{Name: common.ApplicationInsightsNodeJS},
							{Name: common.Credentials},
						},
						Requires: []libcnb.BuildPlanRequire{
							{Name: common.ApplicationInsightsNodeJS},
							{Name: "node", Metadata: map[string]interface{}{"build": true}},
							{Name: common.Credentials},
						},
					},
				},
			}))
		})
	})
}
