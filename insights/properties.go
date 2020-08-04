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
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"

	_ "github.com/paketo-buildpacks/azure-application-insights/insights/statik"
)

type Properties struct {
	LayerContributor libpak.HelperLayerContributor
	Logger           bard.Logger
}

func NewProperties(buildpack libcnb.Buildpack, plan *libcnb.BuildpackPlan) Properties {
	return Properties{
		LayerContributor: libpak.NewHelperLayerContributor(filepath.Join(buildpack.Path, "bin", "azure-application-insights-properties"),
			"Azure Application Insights Properties", buildpack.Info, plan),
	}
}

//go:generate statik -src . -include *.sh

func (p Properties) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	p.LayerContributor.Logger = p.Logger

	return p.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		p.Logger.Bodyf("Copying to %s", layer.Path)

		file := filepath.Join(layer.Path, "bin", "azure-application-insights-properties")
		if err := sherpa.CopyFile(artifact, file); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to copy %s to %s\n%w", artifact.Name(), file, err)
		}

		s, err := sherpa.StaticFile("/properties.sh")
		if err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to load properties.sh\n%w", err)
		}

		layer.Profile.Add("properties.sh", s)

		layer.Launch = true
		return layer, nil
	})
}

func (p Properties) Name() string {
	return p.LayerContributor.LayerName()
}
