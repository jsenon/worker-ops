// Copyright © 2019 Delair <julien.senon@delair.aero>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:generate swagger generate spec -m -o ./swagger.yml

//Worker-ops generate can print static report, or serve an api with slack and metrics functionnality
package main

import "github.com/jsenon/worker-ops/cmd"

//Entry point of the Applicationm, will launch cobra command
func main() {
	cmd.Execute()
}
