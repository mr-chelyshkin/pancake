package pancake

import (
	"fmt"
	"os/user"
	"strconv"
	"strings"
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
	confMaintainer  = "<service_maintainer>"
	confDepartment  = "<service_department>"
	confNamespace   = "<service_namespace>"
	confResLimitCpu = "<cpu_time_pod_limit>"
	confResLimitGpu = "<gpu_time_pod_limit>"
	confResLimitMem = "<mem_pod_limit>"
	confResReqCpu   = "<cpu_time_pod_requested>"
	confResReqGpu   = "<gpu_time_pod_requested>"
	confResReqMem   = "<mem_pod_limit>"

	confAppName               = "<app_name>"
	confAppType               = "<deploy/cronjob/ds/etc>"
	confAppVersioningBy       = "<tag/commit_hash>"
	confAppPostStart          = "<post-start_action>"
	confAppPreStop            = "<pre-stop_action>"
	confAppAffinity           = "<affinity>"
	confAppMaxSurge           = "<max_surge_percentage>"
	confAppMaxUnavailable     = "<max_unavailable_percentage>"
	confAppReplicasNum        = "<replicas_num>"
	confAppInitContainers     = "<init_list_actions>"
	confAppSideContainers     = "<side_list_actions>"

	confFirewallGroup         = "<group_name>"
	confFirewallService       = "<k8s_pod_service>"
	confFirewallMask          = "<ip_mask>"
	confFirewallPortsPort     = "<port>"
	confFirewallPortsProtocol = "<protocol>"

	confPrometheusPath        = "<path>"
	confPrometheusScrape      = "<scrape>"
	confPrometheusPort        = "<port>"
)

type K8STemplate struct {
	Maintainer   string `yaml:"maintainer"`
	Department   string `yaml:"department"`
	Namespace    string `yaml:"namespace"`

	Applications []Application `yaml:"applications"`
}

type Application struct {
	Name               string    `yaml:"name"`
	Type               string    `yaml:"type"`
	Affinity           string    `yaml:"affinity"`
	ReplicasNum        string    `yaml:"replicas_num"`
	VersioningBy       string    `yaml:"versioning_by"`
	PostStart          string    `yaml:"post_start,omitempty"`
	PreStop            string    `yaml:"pre_stop,omitempty"`
	Liveness           string    `yaml:"liveness,omitempty"`
	MaxSurge           string    `yaml:"max_surge,omitempty"`
	MaxUnavailable     string    `yaml:"max_unavailable,omitempty"`
	InitContainers     []string  `yaml:"init_containers,omitempty"`
	SideContainers     []string  `yaml:"side_containers,omitempty"`

	Prometheus         Prometheus `yaml:"prometheus,omitempty"`
	ResourcesLimit     Resources  `yaml:"resources_limit,omitempty"`
	ResourcesRequested Resources  `yaml:"resources_requested,omitempty"`

	Ingress []Firewall `yaml:"ingress,omitempty"`
	Egress  []Firewall `yaml:"egress,omitempty"`
}

type Resources struct {
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

type Prometheus struct {
	Path   string `yaml:"path,omitempty"`
	Scrape string `yaml:"scrape,omitempty"`
	Port   string `yaml:"port,omitempty"`
}

// -- >
func GenerateTemplateObject(appsCount int) K8STemplate {
	wait := make(chan struct{}, 1)
	defer close(wait)

	wait <-struct{}{}
	go __templateServiceIngress__()
	go __templateServiceEgress__()
	go __templateServiceResourcesLimits__()
	go __templateServiceResourcesRequested__()
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
	go __templatePrometheus__()
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

		Prometheus: <-chTemplatePrometheus,
		Ingress: <-chTemplateServiceIngress,
		Egress:  <-chTemplateServiceEgress,

		ResourcesLimit:     <-chTemplateServiceResourcesLimits,
		ResourcesRequested: <-chTemplateServiceResourcesRequested,
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
	chBuf := 19
	chErrMsg := make(chan string, chBuf)

	// -- >
	// 19 concurrency validate functions
	go data.__validateNamespace__(chErrMsg)
	go data.__validateDepartment__(chErrMsg)
	go data.__validateMaintainer__(chErrMsg)
	go data.__validateServiceName__(chErrMsg)
	go data.__validateServiceType__(chErrMsg)
	go data.__validateVersioningBy__(chErrMsg)
	go data.__validateServiceReplicas__(chErrMsg)
	go data.__validateServiceInitContainers__(chErrMsg)
	go data.__validateServiceSideContainers__(chErrMsg)
	go data.__validateMaxSurge__(chErrMsg)
	go data.__validateMaxUnavailable__(chErrMsg)
	go data.__validateServiceAffinity__(chErrMsg)
	go data.__validatePostStart__(chErrMsg)
	go data.__validatePreStop__(chErrMsg)
	go data.__validateServiceResourcesRequested__(chErrMsg)
	go data.__validateServiceResourcesLimits__(chErrMsg)
	go data.__validateServiceEgress__(chErrMsg)
	go data.__validateServiceIngress__(chErrMsg)
	go data.__validateServicePrometheus__(chErrMsg)
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
				return fmt.Errorf(msg)
			}
			return nil
		}
	}
}

/*
	Internal async functions for getting "K8STemplate" struct data fields
*/

//
var chTemplatePrometheus = make(chan Prometheus)
func __templatePrometheus__() {
	chTemplatePrometheus <-Prometheus{
		Path:   confPrometheusPath,
		Scrape: confPrometheusScrape,
		Port:   confPrometheusPort,
	}
	defer close(chTemplatePrometheus)
}

//
var chTemplateServiceIngress = make(chan []Firewall)
func __templateServiceIngress__() {
	chTemplateServiceIngress <-[]Firewall{
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
	chTemplateServiceEgress <-[]Firewall{
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
var chTemplateServiceResourcesLimits = make(chan Resources)
func __templateServiceResourcesLimits__() {
	chTemplateServiceResourcesLimits <-Resources{
		Cpu: confResLimitCpu,
		Mem: confResLimitGpu,
		Gpu: confResLimitMem,
	}
	defer close(chTemplateServiceResourcesLimits)
}

//
var chTemplateServiceResourcesRequested = make(chan Resources)
func __templateServiceResourcesRequested__() {
	chTemplateServiceResourcesRequested <-Resources{
		Cpu: confResReqCpu,
		Mem: confResReqGpu,
		Gpu: confResReqMem,
	}
	defer close(chTemplateServiceResourcesRequested)
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
func (k K8STemplate) __validateServiceReplicas__(ch chan<- string) {
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
			if !strings.HasSuffix(data.MaxSurge, "%") {
				internalCh <-"[ -app:max_surge ] should end with '%%'\n"
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
			if !strings.HasSuffix(data.MaxUnavailable, "%") {
				internalCh <-"[ -app:max_unavailable ] should end with '%%'\n"
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
func (k K8STemplate) __validateServiceResourcesLimits__(ch chan<- string) {
	internalCh := make(chan string, len(k.Applications))

	// -- >
	wg := &sync.WaitGroup{}
	wg.Add(len(k.Applications))
	for _, app := range k.Applications {
		go func(wg *sync.WaitGroup,data Application) {
			defer wg.Done()

			var msg string
			if data.ResourcesLimit.Gpu == confResLimitGpu {
				msg += "[ -app:ResourcesLimit:gpu ] is empty\n"
			}
			if data.ResourcesLimit.Cpu == confResLimitCpu || data.ResourcesLimit.Cpu == "" {
				msg += "[ -app:ResourcesLimit:cpu ] is empty\n"
			}
			if data.ResourcesLimit.Mem == confResLimitMem || data.ResourcesLimit.Mem == "" {
				msg += "[ -app:ResourcesLimit:mem ] is empty\n"
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
func (k K8STemplate) __validateServiceResourcesRequested__(ch chan<- string) {
	internalCh := make(chan string, len(k.Applications))

	// -- >
	wg := &sync.WaitGroup{}
	wg.Add(len(k.Applications))
	for _, app := range k.Applications {
		go func(wg *sync.WaitGroup,data Application) {
			defer wg.Done()

			var msg string
			if data.ResourcesRequested.Gpu == confResReqGpu {
				msg += "[ -app:ResourcesRequested:gpu ] is empty\n"
			}
			if data.ResourcesRequested.Cpu == confResReqCpu || data.ResourcesLimit.Cpu == "" {
				msg += "[ -app:ResourcesRequested:cpu ] is empty\n"
			}
			if data.ResourcesRequested.Mem == confResReqMem || data.ResourcesLimit.Mem == "" {
				msg += "[ -app:ResourcesRequested:mem ] is empty\n"
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
func (k K8STemplate) __validateServicePrometheus__(ch chan<- string) {
	internalCh := make(chan string, len(k.Applications))

	// -- >
	wg := &sync.WaitGroup{}
	wg.Add(len(k.Applications))
	for _, app := range k.Applications {
		go func(wg *sync.WaitGroup,data Application) {
			defer wg.Done()

			empty := Prometheus{}
			if data.Prometheus == empty {
				return
			}

			var msg string
			if data.Prometheus.Path == confPrometheusPath {
				msg += "[ -app:Prometheus:path ] is empty\n"
			}
			if data.Prometheus.Scrape == confPrometheusScrape {
				msg += "[ -app:Prometheus:scrape ] is empty\n"
			}
			if data.Prometheus.Port == confPrometheusPort {
				msg += "[ -app:Prometheus:port ] is empty\n"
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