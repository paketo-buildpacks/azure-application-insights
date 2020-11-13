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

package appinsights

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"

	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
	"github.com/paketo-buildpacks/microsoft-azure/internal/common"
)

const AgentPath = "BPI_AZURE_APPLICATION_INSIGHTS_AGENT_PATH"

type JavaBuild struct {
	BuildpackPath    string
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewJavaBuild(buildpackPath string, dependency libpak.BuildpackDependency, cache libpak.DependencyCache,
	plan *libcnb.BuildpackPlan) JavaBuild {

	return JavaBuild{
		BuildpackPath:    buildpackPath,
		LayerContributor: libpak.NewDependencyLayerContributor(dependency, cache, plan),
	}
}

func (j JavaBuild) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	j.LayerContributor.Logger = j.Logger

	return j.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		j.Logger.Bodyf("Copying to %s", layer.Path)

		file := filepath.Join(layer.Path, filepath.Base(artifact.Name()))
		if err := sherpa.CopyFile(artifact, file); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to copy %s to %s\n%w", artifact.Name(), file, err)
		}

		layer.LaunchEnvironment.Default(AgentPath, file)

		file = filepath.Join(j.BuildpackPath, "resources", "AI-Agent.xml")
		in, err := os.Open(file)
		if err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to open %s\n%w", file, err)
		}
		defer in.Close()

		file = filepath.Join(layer.Path, "AI-Agent.xml")
		if err := sherpa.CopyFile(in, file); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to copy %s to %s\n%w", in.Name(), file, err)
		}

		return layer, nil
	}, libpak.LaunchLayer)
}

func (j JavaBuild) Name() string {
	return j.LayerContributor.LayerName()
}

type JavaLaunch struct {
	CredentialSource common.CredentialSource
	Logger           bard.Logger
}

// https://docs.microsoft.com/en-us/azure/azure-monitor/app/java-agent
func (j JavaLaunch) Execute() (map[string]string, error) {
	if j.CredentialSource == common.None {
		j.Logger.Info("Azure Application Insights disabled")
		return nil, nil
	}

	p, err := sherpa.GetEnvRequired(AgentPath)
	if err != nil {
		return nil, err
	}

	j.Logger.Info("Azure Application Insights enabled")

	return map[string]string{
		"JAVA_TOOL_OPTIONS": sherpa.AppendToEnvVar("JAVA_TOOL_OPTIONS", " ", fmt.Sprintf("-javaagent:%s", p)),
	}, nil
}
