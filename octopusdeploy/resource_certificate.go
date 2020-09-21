package octopusdeploy

import (
	"fmt"
	"log"
	"strconv"

	"github.com/asaskevich/govalidator"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCertificateCreate,
		Read:   resourceCertificateRead,
		Update: resourceCertificateUpdate,
		Delete: resourceCertificateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificate_data": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"environment_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tenanted_deployment_participation": getTenantedDeploymentSchema(),
			"tenant_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tenant_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceCertificateRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	if apiClient == nil {
		log.Println("Client is empty. go-octopusdeploy SDK may be facing issues.")
	}

	if d == nil {
		return createInvalidParameterError("esourceCertificateRead", "d")
	}

	if m == nil {
		return createInvalidParameterError("esourceCertificateRead", "m")
	}

	certificateID := d.Id()
	certificate, err := apiClient.Certificates.Get(certificateID)

	if err == client.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading certificate %s: %s", certificateID, err.Error())
	}

	d.Set("name", certificate.Name)
	d.Set("notes", certificate.Notes)
	d.Set("environment_ids", certificate.EnvironmentIDs)
	d.Set("tenanted_deployment_participation", certificate.TenantedDeploymentParticipation)
	d.Set("tenant_ids", certificate.TenantIds)
	d.Set("tenant_tags", certificate.TenantTags)

	return nil
}

func buildCertificateResource(d *schema.ResourceData) *model.Certificate {
	if d == nil {
		log.Println("The schema for certificate resource is nil")
	}

	certificateName := d.Get("name").(string)

	if govalidator.IsNull("name") {
		fmt.Println("Please confirm the certificate name is a string and is not null")
	}

	str, intErr := strconv.Atoi("name")
	if intErr != nil {
		log.Println(str)
	} else {
		fmt.Println("Please ensure that the name is of type: string")
	}

	var notes string
	var certificateData string
	var password string
	var environmentIds []string
	var tenantedDeploymentParticipation string
	var tenantIds []string
	var tenantTags []string

	notesInterface, ok := d.GetOk("notes")
	if ok {
		notes = notesInterface.(string)
	}

	certificateDataInterface, ok := d.GetOk("certificate_data")
	if ok {
		certificateData = certificateDataInterface.(string)
	}

	passwordInterface, ok := d.GetOk("password")
	if ok {
		password = passwordInterface.(string)
	}

	environmentIdsInterface, ok := d.GetOk("environment_ids")
	if ok {
		environmentIds = getSliceFromTerraformTypeList(environmentIdsInterface)
	}

	if environmentIds == nil {
		environmentIds = []string{}
	}

	tenantedDeploymentParticipationInterface, ok := d.GetOk("tenanted_deployment_participation")
	if ok {
		tenantedDeploymentParticipation = tenantedDeploymentParticipationInterface.(string)
	}

	tenantIdsInterface, ok := d.GetOk("tenant_ids")
	if ok {
		tenantIds = getSliceFromTerraformTypeList(tenantIdsInterface)
	}

	if tenantIds == nil {
		tenantIds = []string{}
	}

	tenantTagsInterface, ok := d.GetOk("tenant_tags")
	if ok {
		tenantTags = getSliceFromTerraformTypeList(tenantTagsInterface)
	}

	if tenantTags == nil {
		tenantTags = []string{}
	}

	var certificate, err = model.NewCertificate(certificateName, model.SensitiveValue{NewValue: &certificateData}, model.SensitiveValue{NewValue: &password})
	certificate.Notes = notes
	certificate.EnvironmentIDs = environmentIds
	certificate.TenantedDeploymentParticipation, _ = enum.ParseTenantedDeploymentMode(tenantedDeploymentParticipation)
	certificate.TenantIds = tenantIds
	certificate.TenantTags = tenantTags

	if err != nil {
		log.Println(err)
	}

	return certificate
}

func resourceCertificateCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	newCertificate := buildCertificateResource(d)
	certificate, err := apiClient.Certificates.Add(newCertificate)

	if err != nil {
		return fmt.Errorf("error creating certificate %s: %s", newCertificate.Name, err.Error())
	}

	d.SetId(certificate.ID)

	return nil
}

func resourceCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	certificate := buildCertificateResource(d)
	certificate.ID = d.Id()

	apiClient := m.(*client.Client)

	updatedCertificate, err := apiClient.Certificates.Update(*certificate)

	if err != nil {
		return fmt.Errorf("error updating certificate id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedCertificate.ID)
	return nil
}

func resourceCertificateDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	certificateID := d.Id()

	err := apiClient.Certificates.Delete(certificateID)

	if err != nil {
		return fmt.Errorf("error deleting certificate id %s: %s", certificateID, err.Error())
	}

	d.SetId("")
	return nil
}
