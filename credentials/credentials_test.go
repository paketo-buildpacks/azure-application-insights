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

package credentials_test

import (
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/microsoft-azure/credentials"
	"github.com/paketo-buildpacks/microsoft-azure/internal/common"
)

func testCredentials(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
	)

	context("launch", func() {
		var (
			l credentials.Launch
		)

		it("does not contribute if source is MetadataServer", func() {
			l.CredentialSource = common.MetadataServer

			Expect(l.Execute()).To(BeNil())
		})

		it("does not contribute if source is None", func() {
			l.CredentialSource = common.None

			Expect(l.Execute()).To(BeNil())
		})

		it("does not contribute if ApplicationCredentials does not exist", func() {
			l.Binding = libcnb.Binding{
				Path:   "/test/path/test-binding",
				Secret: map[string]string{},
			}
			l.CredentialSource = common.Binding

			Expect(l.Execute()).To(BeNil())
		})

		it("contributes credentials if ApplicationCredentials exists", func() {
			l.Binding = libcnb.Binding{
				Path:   "/test/path/test-binding",
				Secret: map[string]string{"InstrumentationKey": "test-value"},
			}
			l.CredentialSource = common.Binding

			Expect(l.Execute()).To(Equal(map[string]string{
				"APPINSIGHTS_INSTRUMENTATIONKEY": "test-value",
			}))
		})
	})
}
