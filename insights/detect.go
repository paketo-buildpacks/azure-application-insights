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

package insights

import (
	"fmt"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
)

type Detect struct{}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	br := libpak.BindingResolver{Bindings: context.Platform.Bindings}

	if _, ok, err := br.Resolve("ApplicationInsights", ""); err != nil {
		return libcnb.DetectResult{}, fmt.Errorf("unable to resolve binding ApplicationInsights\n%w", err)
	} else if !ok {
		return libcnb.DetectResult{Pass: false}, nil
	}

	return libcnb.DetectResult{
		Pass:  true,
		Plans: []libcnb.BuildPlan{
			{
				Provides: []libcnb.BuildPlanProvide{
					{Name: "azure-application-insights-java"},
				},
				Requires: []libcnb.BuildPlanRequire{
					{Name: "azure-application-insights-java"},
					{Name: "jvm-application"},
				},
			},
			{
				Provides: []libcnb.BuildPlanProvide{
					{Name: "azure-application-insights-nodejs"},
				},
				Requires: []libcnb.BuildPlanRequire{
					{Name: "azure-application-insights-nodejs"},
					{Name: "node", Metadata: map[string]interface{}{"build": true}},
				},
			},
		},
	}, nil
}
