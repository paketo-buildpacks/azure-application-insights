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

package properties_test

import (
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/azure-application-insights/properties"
	"github.com/sclevine/spec"
)

func testProperties(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		p properties.Properties
	)

	it("does not contribute properties if no binding exists", func() {
		Expect(p.Execute()).To(BeNil())
	})

	it("contributes properties is binding exists", func() {
		p.Bindings = map[string]libcnb.Binding{
			"test-binding": {
				Metadata: map[string]string{"kind": "ApplicationInsights"},
				Secret:   map[string]string{"instrumentationkey": "test-value"},
			},
		}

		Expect(p.Execute()).To(Equal([]string{`export APPINSIGHTS_INSTRUMENTATIONKEY="test-value"`}))
	})
}
