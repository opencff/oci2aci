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

package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	flagDebug = flag.Bool("debug", false, "Enables debug messages")
)

func usage() {
	fmt.Fprintf(os.Stderr, "NAME:\n")
        fmt.Fprintf(os.Stderr, "    oci2aci - Tool for conversion from oci to aci\n")

	fmt.Fprintf(os.Stderr, "USAGE:\n")
	fmt.Fprintf(os.Stderr, "    oci2aci [--debug] [arguments...]\n")

	fmt.Fprintf(os.Stderr, "VERSION:\n")
        fmt.Fprintf(os.Stderr, "    0.1.0\n")

	fmt.Fprintf(os.Stderr, "FLAGS:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		usage()
		return
	}

	if err := runOCI2ACI(args[0], *flagDebug); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}