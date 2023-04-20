/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MachineSpec defines the desired state of Machine
type MachineSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Machine. Edit machine_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// NetworkInterfaceState is the infrastructure attachment state a NetworkInterface can be in.
type NetworkInterfaceState string

const (
	// NetworkInterfaceStatePending indicates that the attachment of a network interface is pending.
	NetworkInterfaceStatePending NetworkInterfaceState = "Pending"
	// NetworkInterfaceStateAttached indicates that a network interface has been successfully attached.
	NetworkInterfaceStateAttached NetworkInterfaceState = "Attached"
	// NetworkInterfaceStateDetached indicates that a network interface has been successfully detached.
	NetworkInterfaceStateDetached NetworkInterfaceState = "Detached"
)

type MachineState string

const (
	// MachineStatePending means the Machine has been accepted by the system, but not yet completely started.
	// This includes time before being bound to a MachinePool, as well as time spent setting up the Machine on that
	// MachinePool.
	MachineStatePending MachineState = "Pending"
	// MachineStateRunning means the machine is running on a MachinePool.
	MachineStateRunning MachineState = "Running"
	// MachineStateShutdown means the machine is shut down.
	MachineStateShutdown MachineState = "Shutdown"
	// MachineStateTerminated means the machine has been permanently stopped and cannot be started.
	MachineStateTerminated MachineState = "Terminated"
)

// MachineStatus defines the observed state of Machine
type MachineStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	State         MachineState          `json:"state,omitempty"`
	NetworkStatus NetworkInterfaceState `json:"networkstate,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Machine is the Schema for the machines API
type Machine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachineSpec   `json:"spec,omitempty"`
	Status MachineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MachineList contains a list of Machine
type MachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Machine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Machine{}, &MachineList{})
}
