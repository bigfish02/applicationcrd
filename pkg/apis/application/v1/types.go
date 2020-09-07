package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Network describes a Network resource
type Application struct {
	// TypeMeta is the metadata for the resource, like kind and apiversion
	metav1.TypeMeta `json:",inline"`
	// ObjectMeta contains the metadata for the particular object, including
	// things like...
	//  - name
	//  - namespace
	//  - self link
	//  - labels
	//  - ... etc ...
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the custom resource spec
	Spec ApplicationSpec `json:"spec"`
}

// NetworkSpec is the spec for a Network resource
type ApplicationSpec struct {
	Owner           string             `json:"owner"`
	Public          bool               `json:"public"`
	Template        string             `json:"template"`
	ImageRepository string             `json:"imageRepository"`
	ImageName       string             `json:"imageName"`
	Namespace       string             `json:"namespace"`
	ServiceType     string             `json:"serviceType"`
	DefaultPort     int                `json:"defaultPort"`
	GitAddr         string             `json:"gitAddr"`
	ChildApps       []ChildApplication `json:"childApps"`
}

type ChildApplication struct {
	Name            string     `json:"name"`
	Controller      string     `json:"controller"`
	Type            string     `json:"type"`
	Port            int        `json:"port"`
	Command         string     `json:"command"`
	Debug           bool       `json:"debug"`
	ImageName       string     `json:"imageName"`
	ImageRepository string     `json:"imageRepository"`
	TriggerTag      string     `json:"triggerTag"`
	TriggerEnable   bool       `json:"triggerEnable"`
	Pipelines       []Pipeline `json:"pipelines"`
	Clusters        []string   `json:"clusters"`
}

type Pipeline struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApplicationList is a list of Application resources
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Application `json:"items"`
}
