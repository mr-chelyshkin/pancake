package pancake

import (
	"fmt"
	"os/user"
	"strconv"
	"sync"
)

/*
k8s template data generator.

	GenerateTemplateObject():
		async collect K8STemplate struct fields data by internals funcs.
		(use named channels for each field of K8STemplate struct)

	Validate():
		async K8STemplate data validation by internals funcs.
*/

// pre-defined config values
const (
	confMaintainer = "<service_maintainer>"
	confDepartment = "<service_department>"
	confNamespace  = "<service_namespace>"

	confAppName           = "<app_name>"
	confAppType           = "<deploy/cronjob/ds/etc>"
	confAppVersioningBy   = "<tag/commit_hash>"
	confAppPostStart      = "<post-start_action>"
	confAppPreStop        = "<pre-stop_action>"
	confAppAffinity       = "<affinity>"
	confAppMaxSurge       = "<max_surge_percentage>"
	confAppMaxUnavailable = "<max_unavailable_percentage>"
	confAppReplicasNum    = "<replicas_num>"
	confAppInitContainers = "<init_list_actions>"
	confAppSideContainers = "<side_list_actions>"

	confLimitCpu = "<cpu_time_pod_limit>"
	confLimitGpu = "<gpu_time_pod_limit>"
	confLimitMem = "<mem_pod_limit>"

	confFirewallGroup         = "<group_name>"
	confFirewallService       = "<k8s_pod_service>"
	confFirewallMask          = "<ip_mask>"
	confFirewallPortsPort     = "<port>"
	confFirewallPortsProtocol = "<protocol>"
)

type K8STemplate struct {
	Maintainer   string `yaml:"maintainer"`
	Department   string `yaml:"department"`
	Namespace    string `yaml:"namespace"`

	Applications []Application `yaml:"applications"`
}

type Application struct {
	Name           string    `yaml:"name"`
	Type           string    `yaml:"type"`
	Affinity       string    `yaml:"affinity"`
	ReplicasNum    string    `yaml:"replicas_num"`
	VersioningBy   string    `yaml:"versioning_by"`
	PostStart      string    `yaml:"post_start,omitempty"`
	PreStop        string    `yaml:"pre_stop,omitempty"`
	Liveness       string    `yaml:"liveness,omitempty"`
	MaxSurge       string    `yaml:"max_surge,omitempty"`
	MaxUnavailable string    `yaml:"max_unavailable,omitempty"`
	InitContainers []string  `yaml:"init_containers,omitempty"`
	SideContainers []string  `yaml:"side_containers,omitempty"`

	Limit    Limit      `yaml:"limit"`
	Ingress  []Firewall `yaml:"ingress,omitempty"`
	Egress   []Firewall `yaml:"egress,omitempty"`
}

type Limit struct {
	Cpu string `yaml:"cpu"`
	Mem string `yaml:"mem"`
	Gpu string `yaml:"gpu,omitempty"`
}

type Firewall struct {
	Group   string  `yaml:"group"`
	Service string  `yaml:"service"`
	Mask    string  `yaml:"mask"`
	Ports   []Ports `yaml:"ports"`
}

type Ports struct {
	Protocol string `yaml:"protocol"`
	Port     string `yaml:"port"`
}

// -- >
func GenerateTemplateObject(appsCount int) K8STemplate {
	wait := make(chan struct{}, 1)
	defer close(wait)

	wait <-struct{}{}
	go __templateServiceIngress__()
	go __templateServiceEgress__()
	go __templateServiceLimits__()
	go __templateServiceAffinity__()
	go __templateServicePreStop__()
	go __templateServicePostStart__()
	go __templateServiceMaxUnavailable__()
	go __templateServiceMaxSurge__()
	go __templateServiceSideContainers__()
	go __templateServiceInitContainers__()
	go __templateServiceReplicas__()
	go __templateVersioningBy__()
	go __templateServiceType__()
	go __templateServiceName__()
	go __templateMaintainer__()
	go __templateDepartment__()
	go __templateNamespace__()
	<-wait

	app := Application{
		Name:           <-chTemplateServiceName,
		Type:           <-chTemplateServiceType,
		VersioningBy:   <-chTemplateVersioningBy,
		PostStart:      <-chTemplateServicePostStart,
		PreStop:        <-chTemplateServicePreStop,
		Affinity:       <-chTemplateServiceAffinity,
		MaxSurge:       <-chTemplateServiceMaxSurge,
		MaxUnavailable: <-chTemplateServiceMaxUnavailable,
		ReplicasNum:    <-chTemplateServiceReplicas,
		InitContainers: <-chTemplateServiceInitContainers,
		SideContainers: <-chTemplateServiceSideContainers,

		Limit:   <-chTemplateServiceLimits,
		Ingress: <-chTemplateServiceIngress,
		Egress:  <-chTemplateServiceEgress,
	}
	var apps []Application

	for i:=0;i<appsCount;i++ {
		apps = append(apps, app)
	}

	return K8STemplate{
		Maintainer: <-chTemplateMaintainer,
		Department: <-chTemplateDepartment,
		Namespace:  <-chTemplateNamespace,

		Applications: apps,
	}
}

func Validate(data K8STemplate) error {
	// !IMPORTANT: chBuf = count of concurrency validate functions
	// otherwise function will be wait channels or close early
	chBuf := 17
	chErrMsg := make(chan string, chBuf)

	// -- >
	// 17 concurrency validate functions
	go data.__validateNamespace__(chErrMsg)
	go data.__validateDepartment__(chErrMsg)
	go data.__validateMaintainer__(chErrMsg)
	go data.__validateServiceName__(chErrMsg)
	go data.__validateServiceType__(chErrMsg)
	go data.__validateVersioningBy__(chErrMsg)
	go data.__validateServiceReplicas(chErrMsg)
	go data.__validateServiceInitContainers__(chErrMsg)
	go data.__validateServiceSideContainers__(chErrMsg)
	go data.__validateMaxSurge__(chErrMsg)
	go data.__validateMaxUnavailable__(chErrMsg)
	go data.__validateServiceAffinity__(chErrMsg)
	go data.__validatePostStart__(chErrMsg)
	go data.__validatePreStop__(chErrMsg)
	go data.__validateServiceLimits__(chErrMsg)
	go data.__validateServiceEgress__(chErrMsg)
	go data.__validateServiceIngress__(chErrMsg)
	// -- >

	for {
		if len(chErrMsg) == chBuf {
			close(chErrMsg)

			var msg string
			for chOut := range chErrMsg {
				if chOut != "" {
					msg += chOut
				}
			}
			if msg != "" {
				return fmt.Errorf("config validation errors:\n%s", msg)
			}
			return nil
		}
	}
}

/*
	Internal async functions for getting "K8STemplate" struct data fields
*/

//
var chTemplateServiceIngress = make(chan []Firewall)
func __templateServiceIngress__() {
	chTemplateServiceIngress <- []Firewall{
		{
			Group:   confFirewallGroup,
			Service: confFirewallService,
			Mask:    confFirewallMask,
			Ports:   []Ports{
				{
					Protocol: confFirewallPortsProtocol,
					Port:     confFirewallPortsPort,
				},
			},
		},
	}
	defer close(chTemplateServiceIngress)
}

//
var chTemplateServiceEgress = make(chan []Firewall)
func __templateServiceEgress__() {
	chTemplateServiceEgress <- []Firewall{
		{
			Group:   confFirewallGroup,
			Service: confFirewallService,
			Mask:    confFirewallMask,
			Ports:   []Ports{
				{
					Protocol: confFirewallPortsProtocol,
					Port:     confFirewallPortsPort,
				},
			},
		},
	}
	defer close(chTemplateServiceEgress)
}

//
var chTemplateServiceLimits = make(chan Limit)
func __templateServiceLimits__() {
	chTemplateServiceLimits <-Limit{
		Cpu: confLimitCpu,
		Mem: confLimitMem,
		Gpu: confLimitGpu,
	}
	defer close(chTemplateServiceLimits)
}

//
var chTemplateServiceAffinity = make(chan string)
func __templateServiceAffinity__() {
	chTemplateServiceAffinity <-confAppAffinity
	defer close(chTemplateServiceAffinity)
}

//
var chTemplateServicePreStop = make(chan string)
func __templateServicePreStop__() {
	chTemplateServicePreStop <-confAppPreStop
	defer close(chTemplateServicePreStop)
}

//
var chTemplateServicePostStart = make(chan string)
func __templateServicePostStart__() {
	chTemplateServicePostStart <-confAppPostStart
	defer close(chTemplateServicePostStart)
}

//
var chTemplateServiceMaxUnavailable = make(chan string)
func __templateServiceMaxUnavailable__() {
	chTemplateServiceMaxUnavailable <-confAppMaxUnavailable
	defer close(chTemplateServiceMaxUnavailable)
}

//
var chTemplateServiceMaxSurge = make(chan string)
func __templateServiceMaxSurge__() {
	chTemplateServiceMaxSurge <-confAppMaxSurge
	defer close(chTemplateServiceMaxSurge)
}

//
var chTemplateServiceSideContainers = make(chan []string)
func __templateServiceSideContainers__() {
	chTemplateServiceSideContainers <-[]string{confAppSideContainers}
	defer close(chTemplateServiceSideContainers)
}

//
var chTemplateServiceInitContainers = make(chan []string)
func __templateServiceInitContainers__() {
	chTemplateServiceInitContainers <-[]string{confAppInitContainers}
	defer close(chTemplateServiceInitContainers)
}

//
var chTemplateServiceReplicas = make(chan string)
func __templateServiceReplicas__() {
	chTemplateServiceReplicas <-confAppReplicasNum
	defer close(chTemplateServiceReplicas)
}

//
var chTemplateVersioningBy = make(chan string)
func __templateVersioningBy__() {
	chTemplateVersioningBy <-confAppVersioningBy
	defer close(chTemplateVersioningBy)
}

//
var chTemplateServiceType = make(chan string)
func __templateServiceType__() {
	chTemplateServiceType <-confAppType
	defer close(chTemplateServiceType)
}

//
var chTemplateServiceName = make(chan string)
func __templateServiceName__() {
	chTemplateServiceName <-confAppName
	defer close(chTemplateServiceName)
}

//
var chTemplateMaintainer = make(chan string)
func __templateMaintainer__() {
	if currentUser, err := user.Current(); err != nil {
		chTemplateMaintainer <-confMaintainer
	} else {
		chTemplateMaintainer <-currentUser.Username
	}
	defer close(chTemplateMaintainer)
}

//
var chTemplateDepartment = make(chan string)
func __templateDepartment__() {
	chTemplateDepartment <-confDepartment
	defer close(chTemplateDepartment)
}

//
var chTemplateNamespace = make(chan string)
func __templateNamespace__() {
	chTemplateNamespace <-confNamespace
	defer close(chTemplateNamespace)
}

/*
	Internal async functions for validate template values
*/

//
func (k K8STemplate) __validateMaintainer__(ch chan<- string) {
	if k.Maintainer == confMaintainer || k.Maintainer == "" {
		ch <-"[ maintainer ] is empty"
		return
	}
	ch <-""
	return
}

//
func (k K8STemplate) __validateNamespace__(ch chan<- string) {
	if k.Namespace == confNamespace || k.Namespace == "" {
		ch <-"[ namespace ] is empty\n"
		return
	}
	ch <-""
	return
}

//
func (k K8STemplate) __validateDepartment__(ch chan<- string) {
	if k.Department == confDepartment || k.Department == "" {
		ch <-"[ department ] is empty\n"
		return
	}
	ch <-""
	return
}

//
func (k K8STemplate) __validateServiceName__(ch chan<- string) {
	internalCh := make(chan string, len(k.Applications))

	// -- >
	wg := &sync.WaitGroup{}
	wg.Add(len(k.Applications))
	for _, app := range k.Applications {
		go func(wg *sync.WaitGroup,data Application) {
			defer wg.Done()

			if data.Name == confAppName || data.Name == "" {
				internalCh <-"[ -app:name ] is empty\n"
				return
			}

		}(wg, app)
	}
	wg.Wait()
	close(internalCh)
	// -- >

	var msg string
	for validateMsg := range internalCh {
		msg += validateMsg
	}
	ch <-msg
	return
}

//
func (k K8STemplate) __validateServiceType__(ch chan<- string) {
	internalCh := make(chan string, len(k.Applications))

	wg := &sync.WaitGroup{}
	wg.Add(len(k.Applications))

	// -- >
	for _, app := range k.Applications {
		go func(wg *sync.WaitGroup,data Application) {
			defer wg.Done()

			if data.Type == confAppType || data.Type == "" {
				internalCh <-"[ -app:type ] is empty\n"
				return
			}

		}(wg, app)
	}
	wg.Wait()
	close(internalCh)
	// -- >

	var msg string
	for validateMsg := range internalCh {
		msg += validateMsg
	}
	ch <-msg
	return
}

//
func (k K8STemplate) __validateVersioningBy__(ch chan<- string) {
	internalCh := make(chan string, len(k.Applications))

	// -- >
	wg := &sync.WaitGroup{}
	wg.Add(len(k.Applications))
	for _, app := range k.Applications {
		go func(wg *sync.WaitGroup,data Application) {
			defer wg.Done()

			if data.VersioningBy == confAppVersioningBy {
				internalCh <-"[ -app:versioning_by ] is empty\n"
				return
			}
			if data.VersioningBy != "tag" && data.VersioningBy != "hash" {
				internalCh <-"[ -app:versioning_by ] should be 'tag' or 'hash'\n"
				return
			}
		}(wg, app)
	}
	wg.Wait()
	close(internalCh)
	// -- >

	var msg string
	for validateMsg := range internalCh {
		msg += validateMsg
	}
	ch <-msg
	return
}

//
func (k K8STemplate) __validateServiceReplicas(ch chan<- string) {
	internalCh := make(chan string, len(k.Applications))

	// -- >
	wg := &sync.WaitGroup{}
	wg.Add(len(k.Applications))
	for _, app := range k.Applications {
		go func(wg *sync.WaitGroup,data Application) {
			defer wg.Done()

			if data.ReplicasNum == confAppReplicasNum {
				internalCh <-"[ -app:replicas ] is empty\n"
				return
			}
			if _, err := strconv.Atoi(data.ReplicasNum); err != nil {
				internalCh <-"[ -app:replicas ] has invalid chars (should be int)\n"
				return
			}

		}(wg, app)
	}
	wg.Wait()
	close(internalCh)
	// -- >

	var msg string
	for validateMsg := range internalCh {
		msg += validateMsg
	}
	ch <-msg
	return
}

//
func (k K8STemplate) __validateServiceInitContainers__(ch chan<- string) {
	internalCh := make(chan string, len(k.Applications))

	// -- >
	wg := &sync.WaitGroup{}
	wg.Add(len(k.Applications))
	for _, app := range k.Applications {
		go func(wg *sync.WaitGroup,data Application) {
			defer wg.Done()

			if len(data.InitContainers) > 0 {
				if data.InitContainers[0] == confAppInitContainers {
					internalCh <-"[ -app:init_containers ] is empty\n"
					return
				}
			}
		}(wg, app)
	}
	wg.Wait()
	close(internalCh)
	// -- >

	var msg string
	for validateMsg := range internalCh {
		msg += validateMsg
	}
	ch <-msg
	return
}

//
func (k K8STemplate) __validateServiceSideContainers__(ch chan<- string) {
	internalCh := make(chan string, len(k.Applications))

	// -- >
	wg := &sync.WaitGroup{}
	wg.Add(len(k.Applications))
	for _, app := range k.Applications {
		go func(wg *sync.WaitGroup,data Application) {
			defer wg.Done()

			if len(data.SideContainers) > 0 {
				if data.SideContainers[0] == confAppSideContainers {
					internalCh <-"[ -app:side_containers ] is empty\n"
					return
				}
			}

		}(wg, app)
	}
	wg.Wait()
	close(internalCh)
	// -- >

	var msg string
	for validateMsg := range internalCh {
		msg += validateMsg
	}
	ch <-msg
	return
}

//
func (k K8STemplate) __validateMaxSurge__(ch chan<- string) {
	internalCh := make(chan string, len(k.Applications))

	// -- >
	wg := &sync.WaitGroup{}
	wg.Add(len(k.Applications))
	for _, app := range k.Applications {
		go func(wg *sync.WaitGroup,data Application) {
			defer wg.Done()

			if data.MaxSurge == confAppMaxSurge {
				internalCh <-"[ -app:max_surge ] is empty\n"
				return
			}

		}(wg, app)
	}
	wg.Wait()
	close(internalCh)
	// -- >

	var msg string
	for validateMsg := range internalCh {
		msg += validateMsg
	}
	ch <-msg
	return
}

//
func (k K8STemplate) __validateMaxUnavailable__(ch chan<- string) {
	internalCh := make(chan string, len(k.Applications))

	// -- >
	wg := &sync.WaitGroup{}
	wg.Add(len(k.Applications))
	for _, app := range k.Applications {
		go func(wg *sync.WaitGroup,data Application) {
			defer wg.Done()

			if data.MaxUnavailable == confAppMaxUnavailable {
				internalCh <-"[ -app:max_unavailable ] is empty\n"
				return
			}

		}(wg, app)
	}
	wg.Wait()
	close(internalCh)
	// -- >

	var msg string
	for validateMsg := range internalCh {
		msg += validateMsg
	}
	ch <-msg
	return
}

//
func (k K8STemplate) __validatePostStart__(ch chan<- string) {
	internalCh := make(chan string, len(k.Applications))

	// -- >
	wg := &sync.WaitGroup{}
	wg.Add(len(k.Applications))
	for _, app := range k.Applications {
		go func(wg *sync.WaitGroup,data Application) {
			defer wg.Done()

			if data.PostStart == confAppPostStart {
				internalCh <-"[ -app:post_start ] is empty\n"
				return
			}

		}(wg, app)
	}
	wg.Wait()
	close(internalCh)
	// -- >

	var msg string
	for validateMsg := range internalCh {
		msg += validateMsg
	}
	ch <-msg
	return
}

//
func (k K8STemplate) __validatePreStop__(ch chan<- string) {
	internalCh := make(chan string, len(k.Applications))

	// -- >
	wg := &sync.WaitGroup{}
	wg.Add(len(k.Applications))
	for _, app := range k.Applications {
		go func(wg *sync.WaitGroup,data Application) {
			defer wg.Done()

			if data.PreStop == confAppPreStop {
				internalCh <-"[ -app:pre_stop ] is empty\n"
				return
			}

		}(wg, app)
	}
	wg.Wait()
	close(internalCh)
	// -- >

	var msg string
	for validateMsg := range internalCh {
		msg += validateMsg
	}
	ch <-msg
	return
}

//
func (k K8STemplate) __validateServiceAffinity__(ch chan<- string) {
	internalCh := make(chan string, len(k.Applications))

	// -- >
	wg := &sync.WaitGroup{}
	wg.Add(len(k.Applications))
	for _, app := range k.Applications {
		go func(wg *sync.WaitGroup,data Application) {
			defer wg.Done()

			if data.Affinity == confAppAffinity || data.Affinity == "" {
				internalCh <-"[ -app:affinity ] is empty\n"
				return
			}

		}(wg, app)
	}
	wg.Wait()
	close(internalCh)
	// -- >

	var msg string
	for validateMsg := range internalCh {
		msg += validateMsg
	}
	ch <-msg
	return
}

//
func (k K8STemplate) __validateServiceLimits__(ch chan<- string) {
	internalCh := make(chan string, len(k.Applications))

	// -- >
	wg := &sync.WaitGroup{}
	wg.Add(len(k.Applications))
	for _, app := range k.Applications {
		go func(wg *sync.WaitGroup,data Application) {
			defer wg.Done()

			var msg string
			if data.Limit.Gpu == confLimitGpu {
				msg += "[ -app:limit:gpu ] is empty\n"
			}
			if data.Limit.Cpu == confLimitCpu || data.Limit.Cpu == "" {
				msg += "[ -app:limit:cpu ] is empty\n"
			}
			if data.Limit.Mem == confLimitMem || data.Limit.Mem == "" {
				msg += "[ -app:limit:mem ] is empty\n"
			}
			internalCh <-msg

		}(wg, app)
	}
	wg.Wait()
	close(internalCh)
	// -- >

	var msg string
	for validateMsg := range internalCh {
		msg += validateMsg
	}
	ch <-msg
	return
}

//
func (k K8STemplate) __validateServiceEgress__(ch chan<- string) {
	var msg string

	for _, app := range k.Applications {
		for _, _ = range app.Egress {
			///
		}
	}
	ch <-msg
	return
}

//
func (k K8STemplate) __validateServiceIngress__(ch chan<- string) {
	var msg string

	for _, app := range k.Applications {
		for _, _ = range app.Ingress {
			///
		}
	}
	ch <-msg
	return
}