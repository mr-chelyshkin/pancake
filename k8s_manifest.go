package pancake

import (
	"fmt"
	"github.com/flosch/pongo2/v4"
	"gopkg.in/tomb.v1"
	"os"
	"path"
	"sync"
)

func GenerateManifest(templateUser K8STemplate, templateManifestsDir string) error {
	var goroutineTracker tomb.Tomb

	wg := &sync.WaitGroup{}
	wg.Add(len(templateUser.Applications))

	// -- >
	for _, app := range templateUser.Applications {
		go func(wg *sync.WaitGroup, app Application, manifestsDir string) {
			defer wg.Done()

			if app.Ingress != nil {
				err := __generateIngress__(
					manifestsDir,
					templateUser.Namespace,
					templateUser.Department,
					app.Name,
					app.Ingress,
				)
				if err != nil {
					goroutineTracker.Kill(err)
				}
			}
			//if app.Egress != nil {
			//	__generateEgress__(template.Department, app.Egress)
			//}

		}(wg, app, templateManifestsDir)
	}
	wg.Wait()
	// -- >

	select {
	case <-goroutineTracker.Dying():
		return goroutineTracker.Err()
	default:
		return nil
	}
}

//
func getTemplatePath(manifestsDir, department, filename string) (*string, error) {
	var specificTemplate = path.Join(manifestsDir, fmt.Sprintf("%s_%s", department, filename))
	var generalTemplate  = path.Join(manifestsDir, filename)

	if _, err := os.Stat(specificTemplate); err == nil {
		return &specificTemplate, nil
	}
	if _, err := os.Stat(generalTemplate); err == nil {
		return &generalTemplate, nil
	}

	return nil, fmt.Errorf(fmt.Sprintf("manifest template not found, in: %s", manifestsDir))
}

/*
	Internal functions for generating manifests blocks
*/

//
func __generateIngress__(manifestsDir, namespace, department, app string, blocks []Firewall) error {
	template, err := getTemplatePath(manifestsDir, department, "ingress.yaml.j2")
	if err != nil {
		return fmt.Errorf("ingress block, %s", err)
	}

	for _, block := range blocks {
		var tpl = pongo2.Must(pongo2.FromFile(*template))

		out, err := tpl.Execute(pongo2.Context{"item": block, "namespace": namespace, "app": app})
		if err != nil {
			return fmt.Errorf("ingress tpl execute, %s", err)
		}
		fmt.Println(out)
	}
	return nil
}

func __generateEgress__(department string, data []Firewall) {
	fmt.Println("##")
}
