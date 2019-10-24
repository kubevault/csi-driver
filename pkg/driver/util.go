package driver

import (
	"strings"

	"kubevault.dev/csi-driver/pkg/util"

	"github.com/pkg/errors"
)

func getPodInfo(attr map[string]string) (*util.PodInfo, error) {
	podInfo := &util.PodInfo{}
	var ok bool
	podInfo.Name, ok = attr[podName]
	if !ok {
		return nil, errors.Errorf("Pod name not found")
	}
	podInfo.Namespace, ok = attr[podNamespace]
	if !ok {
		return nil, errors.Errorf("Pod namespace not found")
	}
	podInfo.UID, ok = attr[podUID]
	if !ok {
		return nil, errors.Errorf("Pod UID not found")
	}
	podInfo.ServiceAccount, ok = attr[podServiceAccount]
	if !ok {
		return nil, errors.Errorf("Pod service account not found")
	}

	return podInfo, nil
}

func getAppBindingInfo(attr map[string]string) (string, string, error) {
	ref, ok := attr["ref"]
	if !ok {
		return "", "", errors.Errorf("App reference not found")
	}
	data := strings.Split(ref, "/") //namespace/name
	return data[0], data[1], nil
}
