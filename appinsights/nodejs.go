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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"

	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/effect"
	"github.com/paketo-buildpacks/libpak/sherpa"
	"github.com/paketo-buildpacks/microsoft-azure/internal/common"
	"github.com/paketo-buildpacks/microsoft-azure/internal/nodejs"
)

const (
	NodeModule = "applicationinsights"
	NodePath   = "BPI_AZURE_APPLICATION_INSIGHTS_NODE_PATH"
)

type NodeJSBuild struct {
	Executor         effect.Executor
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewNodeJSBuild(dependency libpak.BuildpackDependency, cache libpak.DependencyCache, plan *libcnb.BuildpackPlan) NodeJSBuild {
	return NodeJSBuild{
		Executor:         effect.NewExecutor(),
		LayerContributor: libpak.NewDependencyLayerContributor(dependency, cache, plan),
	}
}

func (n NodeJSBuild) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	n.LayerContributor.Logger = n.Logger

	return n.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		n.Logger.Bodyf("Installing to %s", layer.Path)

		if err := n.Executor.Execute(effect.Execution{
			Command: "npm",
			Args:    []string{"install", "--no-save", artifact.Name()},
			Dir:     layer.Path,
			Stdout:  n.Logger.InfoWriter(),
			Stderr:  n.Logger.InfoWriter(),
		}); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to run npm install\n%w", err)
		}

		layer.LaunchEnvironment.Default(NodePath, filepath.Join(layer.Path, "node_modules"))

		return layer, nil
	}, libpak.LaunchLayer)
}

func (n NodeJSBuild) Name() string {
	return n.LayerContributor.LayerName()
}

type NodeJSLaunch struct {
	ApplicationPath  string
	CredentialSource common.CredentialSource
	Logger           bard.Logger
}

// https://docs.microsoft.com/en-us/azure/azure-monitor/app/nodejs#get-started
func (n NodeJSLaunch) Execute() (map[string]string, error) {
	if n.CredentialSource == common.None {
		n.Logger.Info("Azure Application Insights disabled")
		return nil, nil
	}

	p, err := sherpa.GetEnvRequired(NodePath)
	if err != nil {
		return nil, err
	}

	n.Logger.Info("Azure Application Insights enabled")

	mod, err := sherpa.NodeJSMainModule(n.ApplicationPath)
	if err != nil {
		return nil, fmt.Errorf("unable to find main module in %s\n%w", n.ApplicationPath, err)
	}

	file := filepath.Join(n.ApplicationPath, mod)
	c, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read contents of %s\n%w", file, err)
	}

	e, err := nodejs.IsModuleRequired(NodeModule, c)
	if err != nil {
		return nil, fmt.Errorf("unable to determine if %s is already required\n%w", NodeModule, err)
	}

	if !e {
		n.Logger.Headerf("Requiring '%s' module", NodeModule)

		b := nodejs.RequireModule(NodeModule)

		if err := ioutil.WriteFile(file, append(b, c...), 0644); err != nil {
			return nil, fmt.Errorf("unable to write main module %s\n%w", file, err)
		}
	}

	return map[string]string{
		"NODE_PATH": sherpa.AppendToEnvVar("NODE_PATH", string(os.PathListSeparator), p),
	}, nil
}
