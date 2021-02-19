package pancake

import (
	"fmt"
	"os/user"
)

/*
k8s template data generator.

	GenerateTemplateObject():
		async collect K8STemplate struct fields data by internals funcs.
		(use named channels for each field of K8STemplate struct)
*/

const (
	confMaintainer = "<service_maintainer>"
	confDepartment = "<service_department>"
	confNamespace  = "<service_namespace>"

	confAppName           = "<app_name>"
	confAppType           = "{deploy/cronjob/ds/etc}"
	confAppVersioningBy   = "{tag/commit_hash}"
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
	VersioningBy   string    `yaml:"versioningBy"`
	PostStart      string    `yaml:"postStart,omitempty"`
	PreStop        string    `yaml:"preStop,omitempty"`
	Liveness       string    `yaml:"liveness,omitempty"`
	Affinity       string    `yaml:"affinity,omitempty"`
	MaxSurge       string    `yaml:"maxSurge,omitempty"`
	MaxUnavailable string    `yaml:"maxUnavailable,omitempty"`
	ReplicasNum    string    `yaml:"replicasNum,omitempty"`
	InitContainers []string  `yaml:"initContainers,omitempty"`
	SideContainers []string  `yaml:"sideContainers,omitempty"`

	Limit         Limit      `yaml:"limit"`
	Ingress       []Firewall `yaml:"ingress,omitempty"`
	Egress        []Firewall `yaml:"egress,omitempty"`
}

type Limit struct {
	Cpu string `yaml:"cpu"`
	Mem string `yaml:"mem"`
	Gpu string `yaml:"gpu"`
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
		Maintainer:   <-chTemplateMaintainer,
		Department:   <-chTemplateDepartment,
		Namespace:    <-chTemplateNamespace,

		Applications: apps,
	}
}

func Validate(data K8STemplate) {
	chErrMsg := make(chan string, 17)

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

	for {
		if len(chErrMsg) == 17 {
			close(chErrMsg)
			for i := range chErrMsg {
				fmt.Println(i)
			}
			fmt.Println("asd")

			return
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
	chTemplateServiceLimits <- Limit{
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

func (k K8STemplate) __validateMaintainer__(ch chan<- string) {
	if k.Maintainer == confMaintainer {
		ch <-"required value: namespace not filled"
		return
	}

	ch <-""
	return
}

func (k K8STemplate) __validateNamespace__(ch chan<- string) {
	if k.Namespace == confNamespace {
		ch <-"required value: namespace not filled"
		return
	}

	ch <-""
	return
}

func (k K8STemplate) __validateDepartment__(ch chan<- string) {
	if k.Department == confDepartment {
		ch <-"required value: department not filled"
		return
	}

	ch <-""
	return
}

func (k K8STemplate) __validateServiceName__(ch chan<- string) {
	var msg string

	for _, app := range k.Applications {
		if app.Name == confAppName {
			msg += "required value: app name not filled"
		}
	}

	ch <-msg
	return
}

func (k K8STemplate) __validateServiceType__(ch chan<- string) {
	var msg string

	for _, app := range k.Applications {
		if app.Type == confAppType {
			msg += "required value: app type not filled"
		}
	}

	ch <-msg
	return
}

func (k K8STemplate) __validateVersioningBy__(ch chan<- string) {
	var msg string

	for _, app := range k.Applications {
		if app.VersioningBy == confAppVersioningBy {
			msg += "required value: app versioningBy not filled"
		}
	}

	ch <-msg
	return
}

func (k K8STemplate) __validateServiceReplicas(ch chan<- string) {
	var msg string

	for _, app := range k.Applications {
		if app.ReplicasNum == confAppReplicasNum {
			msg += "required value: app replicas not filled"
		}
	}

	ch <-msg
	return
}

func (k K8STemplate) __validateServiceInitContainers__(ch chan<- string) {
	var msg string

	for _, app := range k.Applications {
		if len(app.InitContainers) > 0 {
			if app.InitContainers[0] == confAppInitContainers {
				msg += "required value: app initContainers not filled"
			}
		}
	}

	ch <-msg
	return
}

func (k K8STemplate) __validateServiceSideContainers__(ch chan<- string) {
	var msg string

	for _, app := range k.Applications {
		if len(app.SideContainers) > 0 {
			if app.SideContainers[0] == confAppSideContainers {
				msg += "required value: app sideContainers not filled"
			}
		}
	}

	ch <-msg
	return
}

func (k K8STemplate) __validateMaxSurge__(ch chan<- string) {
	var msg string

	for _, app := range k.Applications {
		if app.MaxSurge == confAppMaxSurge {
			msg += "required value: app maxSurge not filled"
		}
	}

	ch <-msg
	return
}

func (k K8STemplate) __validateMaxUnavailable__(ch chan<- string) {
	var msg string

	for _, app := range k.Applications {
		if app.MaxUnavailable == confAppMaxUnavailable {
			msg += "required value: app MaxUnavailable not filled"
		}
	}

	ch <-msg
	return
}

func (k K8STemplate) __validatePostStart__(ch chan<- string) {
	var msg string

	for _, app := range k.Applications {
		if app.PostStart == confAppPostStart {
			msg += "required value: app AppPostStart not filled"
		}
	}

	ch <-msg
	return
}

func (k K8STemplate) __validatePreStop__(ch chan<- string) {
	var msg string

	for _, app := range k.Applications {
		if app.PreStop == confAppPreStop {
			msg += "required value: app AppPreStop not filled"
		}
	}

	ch <-msg
	return
}

func (k K8STemplate) __validateServiceAffinity__(ch chan<- string) {
	var msg string

	for _, app := range k.Applications {
		if app.Affinity == confAppAffinity {
			msg += "required value: app confAppAffinity not filled"
		}
	}

	ch <-msg
	return
}

func (k K8STemplate) __validateServiceLimits__(ch chan<- string) {
	var msg string

	for _, app := range k.Applications {
		if app.Limit.Cpu == confLimitCpu {
			msg += "required value: app cpu not filled"
		}
		if app.Limit.Gpu == confLimitGpu {
			msg += "required value: app gpu not filled"
		}
		if app.Limit.Mem == confLimitMem {
			msg += "required value: app mem not filled"
		}
	}

	ch <-msg
	return
}


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