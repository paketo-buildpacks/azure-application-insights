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

package credentials

import (
	"github.com/buildpacks/libcnb"

	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/microsoft-azure/internal/common"
)

type Launch struct {
	Binding          libcnb.Binding
	CredentialSource common.CredentialSource
	Logger           bard.Logger
}

func (l Launch) Execute() (map[string]string, error) {
	if l.CredentialSource == common.MetadataServer || l.CredentialSource == common.None {
		return nil, nil
	}

	if s, ok := l.Binding.Secret["InstrumentationKey"]; ok {
		l.Logger.Info("Configuring Azure Application Insights instrumentation key")
		return map[string]string{"APPINSIGHTS_INSTRUMENTATIONKEY": s}, nil
	}

	return nil, nil
}
