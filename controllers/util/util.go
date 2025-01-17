//
// Copyright 2022 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package util

import (
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"k8s.io/client-go/discovery"
)

// GetOperatorNamespace returns the Namespace of the operator
func GetOperatorNamespace() string {
	ns, found := os.LookupEnv("OPERATOR_NAMESPACE")
	if !found {
		return ""
	}
	return ns
}

// GetWatchNamespace returns the Namespace of the operator
func GetWatchNamespace() string {
	ns, found := os.LookupEnv("WATCH_NAMESPACE")
	if !found {
		return GetOperatorNamespace()
	}
	return ns
}

// GetInstallScope returns the scope of the installation
func GetInstallScope() string {
	ns, found := os.LookupEnv("INSTALL_SCOPE")
	if !found {
		return "cluster"
	}
	return ns
}

func GetIsolatedMode() bool {
	isEnable, found := os.LookupEnv("ISOLATED_MODE")
	if !found || isEnable != "true" {
		return false
	}
	return true
}

func GetoperatorCheckerMode() bool {
	isEnable, found := os.LookupEnv("OPERATORCHECKER_MODE")
	if found && isEnable == "false" {
		return true
	}
	return false
}

// ResourceExists returns true if the given resource kind exists
// in the given api groupversion
func ResourceExists(dc discovery.DiscoveryInterface, apiGroupVersion, kind string) (bool, error) {
	_, apiLists, err := dc.ServerGroupsAndResources()
	if err != nil {
		return false, err
	}
	for _, apiList := range apiLists {
		if apiList.GroupVersion == apiGroupVersion {
			for _, r := range apiList.APIResources {
				if r.Kind == kind {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

//StringSliceContentEqual checks if the contant from two string slice are the same
func StringSliceContentEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// WaitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

// ResourceNamespaced returns true if the given resource is namespaced
func ResourceNamespaced(dc discovery.DiscoveryInterface, apiGroupVersion, kind string) (bool, error) {
	_, apiLists, err := dc.ServerGroupsAndResources()
	if err != nil {
		return false, err
	}
	for _, apiList := range apiLists {
		if apiList.GroupVersion == apiGroupVersion {
			for _, r := range apiList.APIResources {
				if r.Kind == kind {
					return r.Namespaced, nil
				}
			}
		}
	}
	return false, nil
}

func CompareChannelVersion(v1, v2 string) (v1IsLarger bool, err error) {
	_, v1Cut, isExist := strings.Cut(v1, "v")
	if !isExist {
		v1Cut = "0.0"
	}
	v1Slice := strings.Split(v1Cut, ".")
	if len(v1Slice) == 1 {
		v1Cut = v1Cut + ".0"
	}

	_, v2Cut, isExist := strings.Cut(v2, "v")
	if !isExist {
		v1Cut = "0.0"
	}
	v2Slice := strings.Split(v2Cut, ".")
	if len(v2Slice) == 1 {
		v2Cut = v2Cut + ".0"
	}

	v1Slice = strings.Split(v1Cut, ".")
	v2Slice = strings.Split(v2Cut, ".")
	for index := range v1Slice {
		v1SplitInt, e1 := strconv.Atoi(v1Slice[index])
		if e1 != nil {
			return false, e1
		}
		v2SplitInt, e2 := strconv.Atoi(v2Slice[index])
		if e2 != nil {
			return false, e2
		}

		if v1SplitInt > v2SplitInt {
			return true, nil
		} else if v1SplitInt == v2SplitInt {
			continue
		} else {
			return false, nil
		}
	}
	return false, nil
}

func Contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
