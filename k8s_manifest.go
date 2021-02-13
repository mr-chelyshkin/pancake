package pancake

import (
	"fmt"
	"github.com/flosch/pongo2/v4"
	"sync"
)

func GenerateManifest(template K8STemplate)  {
	wg := &sync.WaitGroup{}

	// -
	wg.Add(len(template.Applications))
	for _, app := range template.Applications {
		go func(wg *sync.WaitGroup, app Application) {
			defer wg.Done()

			if app.Ingress != nil {
				__generateIngress__(template.Namespace, template.Department, app.Name, app.Ingress)
			}
			//if app.Egress != nil {
			//	__generateEgress__(template.Department, app.Egress)
			//}

		}(wg, app)
	}
	wg.Wait()
	// -


}

//
func __generateIngress__(namespace, department, app string, blocks []Firewall) {
	for _, block := range blocks {
		var tpl = pongo2.Must(pongo2.FromFile("/Users/i.chelyshkin/Desktop/templates/modules/ingress.yaml.j2"))

		out, err := tpl.Execute(pongo2.Context{"item": block, "namespace": namespace, "app": app})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(out)
	}
}

func __generateEgress__(department string, data []Firewall) {
	fmt.Println("##")
}
