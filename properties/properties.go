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

package properties

import (
	"fmt"
	"strings"

	"github.com/buildpacks/libcnb"

	"github.com/paketo-buildpacks/libpak"
)

type Properties struct {
	Bindings libcnb.Bindings
}

func (p Properties) Execute() (map[string]string, error) {
	br := libpak.BindingResolver{Bindings: p.Bindings}

	b, ok, err := br.Resolve("ApplicationInsights")
	if err != nil {
		return nil, fmt.Errorf("unable to resolve binding ApplicationInsights\n%w", err)
	} else if !ok {
		return nil, nil
	}

	fmt.Println("Configuring Azure Application Insight properties")

	e := make(map[string]string, len(b.Secret))
	for k, v := range b.Secret {
		s := strings.ToUpper(k)
		s = strings.ReplaceAll(s, "-", "_")
		s = strings.ReplaceAll(s, ".", "_")

		e[fmt.Sprintf("APPINSIGHTS_%s", s)] = v
	}

	return e, nil
}
