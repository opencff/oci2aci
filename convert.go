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
	"fmt"
	"os"
	"os/exec"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"github.com/opencontainers/specs"
	"github.com/appc/spec/schema/types"
	"github.com/appc/spec/schema"
)

func runOCI2ACI(path string, flagDebug bool) error {
	if flagDebug {
		InitDebug()
	}
	if bValidate := validateOCIProc(path); bValidate != true {
                fmt.Println("Conversion stop.")
                return nil
        }

	dirWork := createWorkDir()

	convertProc(path, dirWork)

	buildACI(dirWork)

	return nil
}

func createWorkDir() string {
	idir, err := ioutil.TempDir("", "oci2aci")
        if err != nil {
                return ""
        }
        rootfs := filepath.Join(idir, "rootfs")
        os.MkdirAll(rootfs, 0755)

	data := []byte{}
	if err := ioutil.WriteFile(filepath.Join(idir, "manifest"), data, 0644); err != nil {
		return ""
	}
	return idir
}

func genManifest(path string) *schema.ImageManifest {
	//runtimePath := path + "/runtime.json"
	configPath := path + "/config.json"

	/*runtime, err := ioutil.ReadFile(runtimePath)
	if err != nil {
		return nil
	}*/
	
	config, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil
	}
	
	var spec specs.Spec
	err = json.Unmarshal(config, &spec)
	if err != nil {
		return nil
	}

	m := new(schema.ImageManifest)

	// Assemble "acKind", "acVersion" and "name" fields
	m.ACVersion = schema.AppContainerVersion
	m.ACKind = "ImageManifest"
	m.Name = "example"
	
	// Assemble "mountPoints" field
	/*var i int
	for index := range spec.Mounts {
		i = index
	}

	app := new(types.App)
	app.MountPoints = make([]types.MountPoint, i+1)

	for index := range spec.Mounts {
		app.MountPoints[index].Name = types.ACName(spec.Mounts[index].Name)
		app.MountPoints[index].Path = spec.Mounts[index].Path
	}
	
	m.App = app*/
	
	// Assemble "labels" field
	label := new(types.Label)
        label.Name = types.ACIdentifier("version")
        label.Value = spec.Version
        m.Labels = append(m.Labels, *label)

	label = new(types.Label)
	label.Name = types.ACIdentifier("os")
	label.Value = spec.Platform.OS
	m.Labels = append(m.Labels, *label)

	label = new(types.Label)
	label.Name = types.ACIdentifier("arch")
	label.Value = spec.Platform.Arch
        m.Labels = append(m.Labels, *label)

	fmt.Println(m)
	return m
}

func convertProc(srcPath, dstPath string) error {
	src := srcPath + "/rootfs"
	if err := run(exec.Command("cp", "-rf", src, dstPath)); err != nil {
		return err
	}

	m := genManifest(srcPath)
	
	bytes, err := json.Marshal(m)
	if err != nil {
		fmt.Println("ERROR")
		return err
	}
	
	manifestPath := dstPath + "/manifest"

	ioutil.WriteFile(manifestPath, bytes, 0644)
	return nil
}
