package pancake

import (
	"fmt"
	"github.com/flosch/pongo2/v4"
	"gopkg.in/tomb.v1"
	"os"
	"path"
	"sync"
)

/*
k8s manifest generator.

	GenerateManifests():
		async generate k8s manifests from income config and template patterns (jinja).
		for work with jinja use: "github.com/flosch/pongo2"

		return map as: [application_name]k8s_manifest.
*/

func GenerateManifests(templateUser K8STemplate, manifestsDir string) (*map[string]string, error) {
	serviceManifests := make(map[string]string)

	var goroutineTracker tomb.Tomb
	defer goroutineTracker.Done()

	manifests := make(chan [2]string, len(templateUser.Applications))
	defer close(manifests)

	wg := &sync.WaitGroup{}
	wg.Add(len(templateUser.Applications))

	// -- >
	for _, app := range templateUser.Applications {
		go func(wg *sync.WaitGroup, app Application, manifestsDir string) {
			defer wg.Done()

			template, err := generate(
				manifestsDir,
				templateUser.Namespace,
				templateUser.Department,
				templateUser.Maintainer,
				app,
			)

			if err != nil {
				goroutineTracker.Kill(err)
			} else {
				manifests <-[2]string{app.Name, *template}
			}
		}(wg, app, manifestsDir)
	}
	// -- >

	go func() {
		wg.Wait()
	}()

	for {
		select {
		// app: [2]string array: [app_name, k8s_manifest]
		case app := <-manifests:
			serviceManifests[app[0]] = app[1]

			// return if generate all manifests apps from config
			if len(templateUser.Applications) == len(serviceManifests) {
				return &serviceManifests, nil
			}
		case <-goroutineTracker.Dying():
			return nil, goroutineTracker.Err()
		}
	}
}

// -- >
func generate(manifestsDir, namespace, department, maintainer string, app Application) (*string, error) {
	template, err := __getTemplatePath__(manifestsDir, department, "base.yaml.j2")
	if err != nil {
		return nil, err
	}

	var tpl = pongo2.Must(pongo2.FromFile(*template))
	out, err := tpl.Execute(pongo2.Context{"namespace": namespace, "app": app, "maintainer": maintainer})
	if err != nil {
		return nil, err
	}

	return &out, nil
}

//
func __getTemplatePath__(manifestsDir, department, filename string) (*string, error) {
	var specificTemplate = path.Join(manifestsDir, fmt.Sprintf("%s_%s", department, filename))
	var generalTemplate  = path.Join(manifestsDir, filename)

	if _, err := os.Stat(specificTemplate); err == nil {
		return &specificTemplate, nil
	}
	if _, err := os.Stat(generalTemplate); err == nil {
		return &generalTemplate, nil
	}

	return nil, fmt.Errorf(fmt.Sprintf("template not found, in: %s", manifestsDir))
}