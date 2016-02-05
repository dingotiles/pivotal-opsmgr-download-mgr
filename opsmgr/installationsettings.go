package opsmgr

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// InstallationSettingsCommon allows detection of which OpsMgr API version is available
type InstallationSettingsCommon struct {
	GUID                      string `json:"guid"`
	InstallationSchemaVersion string `json:"installation_schema_version"`
}

// InstallationSettingsVersion16 documents all the installation settings for tiles on a v1.6 API OpsMgr
type InstallationSettingsVersion16 struct {
	GUID                      string `json:"guid"`
	InstallationSchemaVersion string `json:"installation_schema_version"`
	Infrastructure            struct {
		Type                  string `json:"type"`
		VMPasswordType        string `json:"vm_password_type"`
		DirectorConfiguration struct {
			ResurrectorEnabled bool     `json:"resurrector_enabled"`
			NtpServers         []string `json:"ntp_servers"`
			BlobstoreType      string   `json:"blobstore_type"`
			DatabaseType       string   `json:"database_type"`
		} `json:"director_configuration"`
		Networks []struct {
			GUID                  string `json:"guid"`
			Name                  string `json:"name"`
			IaasNetworkIdentifier string `json:"iaas_network_identifier"`
			Subnet                string `json:"subnet"`
			DNS                   string `json:"dns"`
			Gateway               string `json:"gateway"`
			ReservedIPRanges      string `json:"reserved_ip_ranges"`
		} `json:"networks"`
		AvailabilityZones []struct {
			GUID    string `json:"guid"`
			Name    string `json:"name"`
			Cluster string `json:"cluster"`
		} `json:"availability_zones"`
		IaasConfiguration struct {
			Datacenter         string   `json:"datacenter"`
			VcenterIP          string   `json:"vcenter_ip"`
			VcenterUsername    string   `json:"vcenter_username"`
			VcenterPassword    string   `json:"vcenter_password"`
			Datastores         []string `json:"datastores"`
			BoshVMFolder       string   `json:"bosh_vm_folder"`
			BoshTemplateFolder string   `json:"bosh_template_folder"`
			BoshDiskPath       string   `json:"bosh_disk_path"`
		} `json:"iaas_configuration"`
	} `json:"infrastructure"`
	Products []struct {
		GUID                               string              `json:"guid"`
		InstallationName                   string              `json:"installation_name"`
		ProductVersion                     string              `json:"product_version"`
		SingletonAvailabilityZoneReference string              `json:"singleton_availability_zone_reference"`
		Ips                                map[string][]string `json:"ips"`
		Prepared                           bool                `json:"prepared"`
		Stemcell                           struct {
			Infrastructure string `json:"infrastructure"`
			Hypervisor     string `json:"hypervisor"`
			Os             string `json:"os"`
			Version        string `json:"version"`
			File           string `json:"file"`
			Name           string `json:"name"`
		} `json:"stemcell"`
		Jobs []struct {
			GUID             string `json:"guid"`
			InstallationName string `json:"installation_name"`
			Properties       []struct {
				Value      interface{} `json:"value"`
				Identifier string      `json:"identifier"`
			} `json:"properties"`
			Instances []struct {
				Value      int    `json:"value"`
				Identifier string `json:"identifier"`
			} `json:"instances"`
			Resources []struct {
				Value      int    `json:"value"`
				Identifier string `json:"identifier"`
			} `json:"resources"`
			Partitions []struct {
				JobReference              string `json:"job_reference"`
				InstallationName          string `json:"installation_name"`
				InstanceCount             int    `json:"instance_count"`
				AvailabilityZoneReference string `json:"availability_zone_reference"`
			} `json:"partitions"`
			Identifier string `json:"identifier"`
		} `json:"jobs"`
		InfrastructureNetworkReference string `json:"infrastructure_network_reference"`
		DeploymentNetworkReference     string `json:"deployment_network_reference"`
		Identifier                     string `json:"identifier"`
	} `json:"products"`
}

// SecurityProperty represents a job property for a security/credential item
type SecurityProperty struct {
	Identity string `json:"identity"`
	Salt     string `json:"salt"`
	Password string `json:"password"`
}

// GetInstallationSettings gets the installation settings from target OpsMgr
// Currently limited to v1.6 API
func (opsmgr *OpsMgr) GetInstallationSettings() (err error) {
	fmt.Println("Fetching OpsMgr installation settings...")
	req, err := http.NewRequest("GET", opsmgr.apiURL("/api/installation_settings"), nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(opsmgr.Username, opsmgr.Password)

	resp, err := opsmgr.httpClient().Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	opsmgr.InstallationSettings = &InstallationSettingsVersion16{}
	err = json.NewDecoder(resp.Body).Decode(opsmgr.InstallationSettings)
	return
}
