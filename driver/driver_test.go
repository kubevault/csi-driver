package driver

import
(
	"testing"
	//core "k8s.io/api/core/v1"
	"fmt"
	vaultapi "github.com/hashicorp/vault/api"
	//"k8s.io/kubernetes/pkg/apis/core"
)

func Test_PVC(t *testing.T) {
	//pvc := core.PersistentVolumeClaim{

	//}
	//fmt.Println(pvc)
	//pod := core.Pod{}
	//fmt.Println(pod)
}

func Test_Vault(t *testing.T) {
	vc, err:= NewVaultClient("", "root",  &vaultapi.TLSConfig{Insecure: false})
	fmt.Println(err, vc.vc.Token())

	//vc.vc.Auth().Token().Lookup()

	sec, err := vc.getTokenForPolicy([]string{"nginx"}, "localtest")
	fmt.Println(sec, err)
}
