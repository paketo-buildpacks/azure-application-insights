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
	"github.com/paketo-buildpacks/libpak/bard"
)

type Build struct {
	Logger bard.Logger
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	b.Logger.Title(context.Buildpack)
	result := libcnb.NewBuildResult()

	pr := libpak.PlanEntryResolver{Plan: context.Plan}

	dr, err := libpak.NewDependencyResolver(context)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency resolver\n%w", err)
	}

	dc := libpak.NewDependencyCache(context.Buildpack)
	dc.Logger = b.Logger

	if _, ok, err := pr.Resolve("azure-application-insights-java"); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve azure-application-insights-java plan entry\n%w", err)
	} else if ok {
		dep, err := dr.Resolve("azure-application-insights-java", "")
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		ja := NewJavaAgent(context.Buildpack.Path, dep, dc, result.Plan)
		ja.Logger = b.Logger
		result.Layers = append(result.Layers, ja)
	}

	if _, ok, err := pr.Resolve("azure-application-insights-nodejs"); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve azure-application-insights-nodejs plan entry\n%w", err)
	} else if ok {
		dep, err := dr.Resolve("azure-application-insights-nodejs", "")
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		na := NewNodeJSAgent(context.Application.Path, dep, dc, result.Plan)
		na.Logger = b.Logger
		result.Layers = append(result.Layers, na)
	}

	p := NewProperties(context.Buildpack, result.Plan)
	p.Logger = b.Logger
	result.Layers = append(result.Layers, p)

	return result, nil
}
