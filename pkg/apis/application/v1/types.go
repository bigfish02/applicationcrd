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
	// Cidr and Gateway are example custom spec fields
	//
	// this is where you would put your custom resource data
	Owner           string             `json:"owner"`
	Public          bool               `json:"public"`
	Template        string             `json:"template"`
	ImageRepository string             `json:"imageRepository"`
	ImageName       string             `json:"imageName"`
	Namespace       string             `json:"namespace"`
	GitAddr         string             `json:"gitAddr"`
	ChildApps       []ChildApplication `json:"childApps"`
}

type ChildApplication struct {
	Name       string `json:"name"`
	Controller string `json:"controller"`
	Type       string `json:"type"`
	Yaml       string `json:"yaml"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkList is a list of Network resources
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Application `json:"items"`
}
