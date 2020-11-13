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

package nodejs_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/microsoft-azure/internal/nodejs"
)

func testModule(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
	)

	context("IsModuleRequired", func() {
		it("detects required module", func() {
			Expect(nodejs.IsModuleRequired("test-module", []byte(`require('test-module')`))).To(BeTrue())
			Expect(nodejs.IsModuleRequired("test-module", []byte(`require("test-module")`))).To(BeTrue())
			Expect(nodejs.IsModuleRequired("test-module", []byte(`test-data
require('test-module')
test-data`))).To(BeTrue())
		})

		it("detects not-required module", func() {
			Expect(nodejs.IsModuleRequired("test-module", []byte(`require('another-module')`))).To(BeFalse())
			Expect(nodejs.IsModuleRequired("test-module", []byte(`require("another-module")`))).To(BeFalse())
			Expect(nodejs.IsModuleRequired("test-module", []byte(`test-data
require('another-module')
test-data`))).To(BeFalse())
		})
	})

	context("RequireModule", func() {

		it("renders require module", func() {
			b := nodejs.RequireModule("test-module")

			Expect(string(b)).To(Equal(`require('test-module').start();
`))
		})
	})
}
