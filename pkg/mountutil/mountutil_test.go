/*
   Copyright Farcloser.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package mountutil

import (
	"runtime"

	"go.farcloser.world/lepton/leptonic/api"
	"go.farcloser.world/lepton/pkg/mountutil/volumestore"
)

type MockVolumeStore struct {
	volumestore.VolumeService
}

func (mv *MockVolumeStore) CreateWithoutLock(name string, labels map[string]string) (*api.Volume, error) {
	if runtime.GOOS == "windows" {
		return &api.Volume{Name: "test_volume", Mountpoint: "C:\\test\\directory"}, nil
	}
	return &api.Volume{Name: "test_volume", Mountpoint: "/test/volume"}, nil
}

var mockVolumeStore = &MockVolumeStore{}
