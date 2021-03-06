package os

import (
	"context"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// GuestOsToCommonKey represents the element key in OS map of guest OS to common template mapping
	GuestOsToCommonKey = "guestos2common"

	// OsInfoToCommonKey represents the element key in OS map of OS type to common template mapping
	OsInfoToCommonKey = "osinfo2common"

	// OsConfigMapName represents the environment variable name that holds the OS config map name
	OsConfigMapName = "OS_CONFIGMAP_NAME"

	// OsConfigMapNamespace represents the environment variable name that holds the OS config map namespace
	OsConfigMapNamespace = "OS_CONFIGMAP_NAMESPACE"
)

// OSMaps is responsible for getting the operating systems maps which contain mapping of GuestOS to common templates
// and mapping of osinfo to common templates
type OSMaps struct {
	Client client.Client
}

// OSMapProvider is responsible for getting the operating systems maps
type OSMapProvider interface {
	GetOSMaps() (map[string]string, map[string]string, error)
}

// NewOSMapProvider creates new OSMapProvider
func NewOSMapProvider(client client.Client) *OSMaps {
	return &OSMaps{
		Client: client,
	}
}

// GetOSMaps retrieve the OS mapping config map
func (o *OSMaps) GetOSMaps() (map[string]string, map[string]string, error) {
	guestOsToCommon := initGuestOsToCommon()
	osInfoToCommon := initOsInfoToCommon()
	err := o.updateOsMapsByUserMaps(guestOsToCommon, osInfoToCommon)
	return guestOsToCommon, osInfoToCommon, err
}

func (o *OSMaps) updateOsMapsByUserMaps(guestOsToCommon map[string]string, osInfoToCommon map[string]string) error {
	osConfigMap := &corev1.ConfigMap{}
	osConfigMapName := os.Getenv(OsConfigMapName)
	osConfigMapNamespace := os.Getenv(OsConfigMapNamespace)
	if osConfigMapName == "" && osConfigMapNamespace == "" {
		return nil
	}

	err := o.Client.Get(
		context.TODO(),
		types.NamespacedName{Name: osConfigMapName, Namespace: osConfigMapNamespace},
		osConfigMap,
	)
	if err != nil {
		return fmt.Errorf(
			"Failed to read user OS config-map [%s/%s] due to: [%v]",
			osConfigMapNamespace,
			osConfigMapName,
			err,
		)
	}

	err = yaml.Unmarshal([]byte(osConfigMap.Data[GuestOsToCommonKey]), &guestOsToCommon)
	if err != nil {
		return fmt.Errorf(
			"Failed to parse user OS config-map [%s/%s] element %s due to: [%v]",
			osConfigMapNamespace,
			osConfigMapName,
			GuestOsToCommonKey,
			err,
		)
	}

	err = yaml.Unmarshal([]byte(osConfigMap.Data[OsInfoToCommonKey]), &osInfoToCommon)
	if err != nil {
		return fmt.Errorf(
			"Failed to parse user OS config-map [%s/%s] element %s due to: [%v]",
			osConfigMapNamespace,
			osConfigMapName,
			OsInfoToCommonKey,
			err,
		)
	}

	return nil
}

func initGuestOsToCommon() map[string]string {
	return map[string]string{
		"Red Hat Enterprise Linux Server": "rhel",
		"CentOS Linux":                    "centos",
		"Fedora":                          "fedora",
		"Ubuntu":                          "ubuntu",
		"openSUSE":                        "opensuse",
	}
}

func initOsInfoToCommon() map[string]string {
	return map[string]string{
		"rhel_6_9_plus_ppc64": "rhel6.9",
		"rhel_6_ppc64":        "rhel6.9",
		"rhel_6":              "rhel6.9",
		"rhel_6x64":           "rhel6.9",
		"rhel_7_ppc64":        "rhel7.7",
		"rhel_7_s390x":        "rhel7.7",
		"rhel_7x64":           "rhel7.7",
		"rhel_8x64":           "rhel8.1",
		"sles_11_ppc64":       "opensuse15.0",
		"sles_11":             "opensuse15.0",
		"sles_12_s390x":       "opensuse15.0",
		"ubuntu_12_04":        "ubuntu18.04",
		"ubuntu_12_10":        "ubuntu18.04",
		"ubuntu_13_04":        "ubuntu18.04",
		"ubuntu_13_10":        "ubuntu18.04",
		"ubuntu_14_04_ppc64":  "ubuntu18.04",
		"ubuntu_14_04":        "ubuntu18.04",
		"ubuntu_16_04_s390x":  "ubuntu18.04",
		"windows_10":          "win10",
		"windows_10x64":       "win10",
		"windows_2003":        "win10",
		"windows_2003x64":     "win10",
		"windows_2008R2x64":   "win2k8",
		"windows_2008":        "win2k8",
		"windows_2008x64":     "win2k8",
		"windows_2012R2x64":   "win2k12r2",
		"windows_2012x64":     "win2k12r2",
		"windows_2016x64":     "win2k16",
		"windows_2019x64":     "win2k19",
		"windows_7":           "win10",
		"windows_7x64":        "win10",
		"windows_8":           "win10",
		"windows_8x64":        "win10",
		"windows_xp":          "win10",
	}
}
