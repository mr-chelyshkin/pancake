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

	manifests := make(chan string, 1)
	defer close(manifests)

	wg := &sync.WaitGroup{}
	wg.Add(len(templateUser.Applications))

	// -- >
	for _, app := range templateUser.Applications {
		go func(wg *sync.WaitGroup, app Application, manifestsDir string) {
			var appManifests string
			defer wg.Done()

			if app.Ingress != nil {
				block, err := __generateIngress__(
					manifestsDir,
					templateUser.Namespace,
					templateUser.Department,
					app.Name,
					app.Ingress,
				)
				if err != nil {
					goroutineTracker.Kill(err)
				} else {
					appManifests += *block
				}
			}

			if app.Egress != nil {
				block, err := __generateEgress__(
					manifestsDir,
					templateUser.Namespace,
					templateUser.Department,
					app.Name,
					app.Egress,
				)
				if err != nil {
					goroutineTracker.Kill(err)
				} else {
					appManifests += *block
				}
			}

			manifests <-appManifests
		}(wg, app, templateManifestsDir)
	}

	go func() {
		wg.Wait()
	}()

	select {
	case f := <-manifests:
		fmt.Println(f)
	case <-goroutineTracker.Dying():
		fmt.Println(goroutineTracker.Err())
	}

	goroutineTracker.Done()
	return nil
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
func __generateIngress__(manifestsDir, namespace, department, app string, blocks []Firewall) (*string, error) {
	template, err := getTemplatePath(manifestsDir, department, "ingress.yaml.j2")
	if err != nil {
		return nil, fmt.Errorf("ingress block, %s", err)
	}
	var data string

	for _, block := range blocks {
		var tpl = pongo2.Must(pongo2.FromFile(*template))

		out, err := tpl.Execute(pongo2.Context{"item": block, "namespace": namespace, "app": app})
		if err != nil {
			return nil, fmt.Errorf("ingress tpl execute, %s", err)
		}
		data += out
	}
	return &data, nil
}

func __generateEgress__(manifestsDir, namespace, department, app string, blocks []Firewall) (*string, error) {
	template, err := getTemplatePath(manifestsDir, department, "egress.yaml.j2")
	if err != nil {
		return nil, fmt.Errorf("egress block, %s", err)
	}
	var data string

	for _, block := range blocks {
		var tpl = pongo2.Must(pongo2.FromFile(*template))

		out, err := tpl.Execute(pongo2.Context{"item": block, "namespace": namespace, "app": app})
		if err != nil {
			return nil, fmt.Errorf("egress tpl execute, %s", err)
		}
		data += out
	}
	return &data, nil
}
