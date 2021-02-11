package pancake

import (
	"os/user"
)

/*
k8s template data generator.

	GenerateTemplateObject():
		async collect K8STemplate struct fields data by internals funcs.
		(use named channels for each field of K8STemplate struct)
*/

type K8STemplate struct {
	Maintainer   string `yaml:"maintainer"`
	Department   string `yaml:"department"`

	Applications []Application `yaml:"applications"`
}

type Application struct {
	Name           string    `yaml:"name"`
	Type           string    `yaml:"type"`
	VersioningBy   string    `yaml:"versioningBy"`
	PostStart      string    `yaml:"postStart"`
	PreStop        string    `yaml:"preStop"`
	Liveness       string    `yaml:"liveness"`
	Affinity       string    `yaml:"affinity"`
	MaxSurge       string    `yaml:"maxSurge"`
	MaxUnavailable string    `yaml:"maxUnavailable"`
	ReplicasNum    string    `yaml:"replicasNum"`
	InitContainers []string  `yaml:"initContainers"`
	SideContainers []string  `yaml:"sideContainers"`

	Limit         Limit      `yaml:"limit"`
	Ingress       []Firewall `yaml:"ingress"`
	Egress        []Firewall `yaml:"egress"`
}

type Limit struct {
	Cpu string `yaml:"cpu"`
	Mem string `yaml:"mem"`
	Gpu string `yaml:"gpu"`
}

type Firewall struct {
	Group string   `yaml:"group"`
	Mask  string   `yaml:"mask"`
	Ports []string `yaml:"ports"`
}

// -- >
func GenerateTemplateObject(appsCount int) K8STemplate {
	wait := make(chan struct{}, 1)

	wait <- struct{}{}
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

		Applications: apps,
	}
}

/*
	Internal async funcs for getting "K8STemplate" struct data fields
*/

//
var chTemplateServiceIngress = make(chan []Firewall)
func __templateServiceIngress__() {
	chTemplateServiceIngress <- []Firewall{
		{
			Group: "[string] Set group name",
			Mask:  "[string] Set group ip mask",
			Ports: []string{"[slice] Set ports"},
		},
	}
	defer close(chTemplateServiceIngress)
}

//
var chTemplateServiceEgress = make(chan []Firewall)
func __templateServiceEgress__() {
	chTemplateServiceEgress <- []Firewall{
		{
			Group: "[string] Set group name",
			Mask:  "[string] Set group ip mask",
			Ports: []string{"[slice] Set ports"},
		},
	}
	defer close(chTemplateServiceEgress)
}

//
var chTemplateServiceLimits = make(chan Limit)
func __templateServiceLimits__() {
	chTemplateServiceLimits <- Limit{
		Cpu: "[string] Set processor time pod limit",
		Mem: "[string] Set memory pod limit",
		Gpu: "[string] Set gpu time pod limit",
	}
}

//
var chTemplateServiceAffinity = make(chan string)
func __templateServiceAffinity__() {
	chTemplateServiceAffinity <-"[choice] Set { much/less }"
	defer close(chTemplateServiceAffinity)
}

//
var chTemplateServicePreStop = make(chan string)
func __templateServicePreStop__() {
	chTemplateServicePreStop <-"[string] Set pre stop actions"
	defer close(chTemplateServicePreStop)
}

//
var chTemplateServicePostStart = make(chan string)
func __templateServicePostStart__() {
	chTemplateServicePostStart <-"[string] Set post start actions"
	defer close(chTemplateServicePostStart)
}

//
var chTemplateServiceMaxUnavailable = make(chan string)
func __templateServiceMaxUnavailable__() {
	chTemplateServiceMaxUnavailable <-"[string] Set percentage of max unavailable"
}

//
var chTemplateServiceMaxSurge = make(chan string)
func __templateServiceMaxSurge__() {
	chTemplateServiceMaxSurge <-"[string] Set percentage of max surge"
	defer close(chTemplateServiceMaxSurge)
}

//
var chTemplateServiceSideContainers = make(chan []string)
func __templateServiceSideContainers__() {
	chTemplateServiceSideContainers <-[]string{"[slice] Set side actions containers"}
	defer close(chTemplateServiceSideContainers)
}

//
var chTemplateServiceInitContainers = make(chan []string)
func __templateServiceInitContainers__() {
	chTemplateServiceInitContainers <-[]string{"[slice] Set init actions containers"}
	defer close(chTemplateServiceInitContainers)
}

//
var chTemplateServiceReplicas = make(chan string)
func __templateServiceReplicas__() {
	chTemplateServiceReplicas <-"[string] Set num of service replicas count"
	defer close(chTemplateServiceReplicas)
}

//
var chTemplateVersioningBy = make(chan string)
func __templateVersioningBy__() {
	chTemplateVersioningBy <-"[choice] Set { tag/commit_hash }"
	defer close(chTemplateVersioningBy)
}

//
var chTemplateServiceType = make(chan string)
func __templateServiceType__() {
	chTemplateServiceType <-"[choice] Set { deploy/cronjob/ds/etc }"
	defer close(chTemplateServiceType)
}

//
var chTemplateServiceName = make(chan string)
func __templateServiceName__() {
	chTemplateServiceName <-"[string] Set name of service"
	defer close(chTemplateServiceName)
}

//
var chTemplateMaintainer = make(chan string)
func __templateMaintainer__() {
	if currentUser, err := user.Current(); err != nil {
		chTemplateMaintainer <-"[string] Set user who maintain service"
	} else {
		chTemplateMaintainer <-currentUser.Username
	}
	defer close(chTemplateMaintainer)
}

//
var chTemplateDepartment = make(chan string)
func __templateDepartment__() {
	chTemplateDepartment <-"[string] Set department witch service belong"
	defer close(chTemplateDepartment)
}
