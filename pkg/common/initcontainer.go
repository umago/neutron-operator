package common

import (
	corev1 "k8s.io/api/core/v1"
)

// InitContainer information
type InitContainer struct {
	Privileged     bool
	ContainerImage string
	Database       string
	DatabaseHost   string
	NeutronSecret  string
	NovaSecret     string
	VolumeMounts   []corev1.VolumeMount
}

// GetInitContainer - init container for cinder services
func GetInitContainer(init InitContainer) []corev1.Container {
	runAsUser := int64(0)
	trueVar := true

	securityContext := &corev1.SecurityContext{
		RunAsUser: &runAsUser,
	}
	if init.Privileged {
		securityContext.Privileged = &trueVar
	}

	return []corev1.Container{
		{
			Name:            "init",
			Image:           init.ContainerImage,
			SecurityContext: securityContext,
			Command: []string{
				"/bin/bash", "-c", "/usr/local/bin/container-scripts/init.sh",
			},
			Env: []corev1.EnvVar{
				{
					Name: "TransportURL",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: init.NeutronSecret,
							},
							Key: "TransportUrl",
						},
					},
				},
				{
					Name:  "DatabaseHost",
					Value: init.DatabaseHost,
				},
				{
					Name:  "Database",
					Value: init.Database,
				},
				{
					Name: "DatabasePassword",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: init.NeutronSecret,
							},
							Key: "DatabasePassword",
						},
					},
				},
				{
					Name: "NeutronKeystoneAuthPassword",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: init.NeutronSecret,
							},
							Key: "NeutronKeystoneAuthPassword",
						},
					},
				},
				{
					Name: "NovaKeystoneAuthPassword",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: init.NovaSecret,
							},
							Key: "NovaKeystoneAuthPassword",
						},
					},
				},
			},
			VolumeMounts: init.VolumeMounts,
		},
	}
}
