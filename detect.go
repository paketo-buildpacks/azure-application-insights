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

package azure

import (
	"fmt"

	"github.com/buildpacks/libcnb"

	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/microsoft-azure/internal/common"
)

const (
	ApplicationInsightsEnabled = "BP_AZURE_APPLICATION_INSIGHTS_ENABLED"
)

type Detect struct{}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	cr, err := libpak.NewConfigurationResolver(context.Buildpack, nil)
	if err != nil {
		return libcnb.DetectResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	if !cr.ResolveBool(ApplicationInsightsEnabled) {
		return libcnb.DetectResult{Pass: false}, nil
	}

	return libcnb.DetectResult{
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
	}, nil
}
