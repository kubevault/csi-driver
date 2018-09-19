package driver

import
(
	"testing"
	core "k8s.io/api/core/v1"
	"fmt"
)

func Test_PVC(t *testing.T) {
	pvc := core.PersistentVolumeClaim{

	}
	fmt.Println(pvc)
	pod := core.Pod{}
	fmt.Println(pod)
}

