package conversations

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/RJPearson94/terraform-provider-twilio/twilio/common"
	"github.com/RJPearson94/terraform-provider-twilio/twilio/utils"
	"github.com/RJPearson94/twilio-sdk-go/service/conversations/v1/configuration/address"
	"github.com/RJPearson94/twilio-sdk-go/service/conversations/v1/configuration/addresses"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceConversationsAddressConfigurationWebhook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConversationsAddressConfigurationWebhookCreate,
		ReadContext:   resourceConversationsAddressConfigurationWebhookRead,
		UpdateContext: resourceConversationsAddressConfigurationWebhookUpdate,
		DeleteContext: resourceConversationsAddressConfigurationWebhookDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				format := "/Configuration/Addresses/(.*)"
				regex := regexp.MustCompile(format)
				match := regex.FindStringSubmatch(d.Id())

				if len(match) != 2 {
					return nil, fmt.Errorf("The imported ID (%s) does not match the format (%s)", d.Id(), format)
				}

				d.Set("sid", match[1])
				d.SetId(match[1])
				return []*schema.ResourceData{d}, nil
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"sid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_sid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_sid": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.ConversationServiceSidValidation(),
			},
			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"integration_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"sms",
					"whatsapp",
				}, false),
				ForceNew: true,
			},
			"webhook_filters": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"onMessageAdded",
						"onMessageUpdated",
						"onMessageRemoved",
						"onConversationUpdated",
						"onConversationStateUpdated",
						"onConversationRemoved",
						"onParticipantAdded",
						"onParticipantUpdated",
						"onParticipantRemoved",
						"onDeliveryUpdated",
					}, false),
				},
			},
			"webhook_method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "POST",
				ValidateFunc: validation.StringInSlice([]string{
					"GET",
					"POST",
				}, false),
			},
			"webhook_url": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceConversationsAddressConfigurationWebhookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).Conversations

	webhookType := "webhook"
	createInput := &addresses.CreateAddressInput{
		Address: d.Get("address").(string),
		AutoCreation: &addresses.CreateAutoCreationInput{
			ConversationServiceSid: utils.OptionalStringWithEmptyStringOnChange(d, "service_sid"),
			Enabled:                utils.OptionalBool(d, "enabled"),
			Type:                   &webhookType,
			WebhookFilters:         utils.OptionalStringSlice(d, "webhook_filters"),
			WebhookMethod:          utils.OptionalString(d, "webhook_method"),
			WebhookUrl:             utils.OptionalString(d, "webhook_url"),
		},
		FriendlyName: utils.OptionalStringWithEmptyStringOnChange(d, "friendly_name"),
		Type:         d.Get("type").(string),
	}

	createResult, err := client.Configuration().Addresses.CreateWithContext(ctx, createInput)
	if err != nil {
		return diag.Errorf("Failed to create address configuration webhook: %s", err.Error())
	}

	d.SetId(createResult.Sid)
	return resourceConversationsAddressConfigurationWebhookRead(ctx, d, meta)
}

func resourceConversationsAddressConfigurationWebhookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).Conversations

	getResponse, err := client.Configuration().Address(d.Id()).FetchWithContext(ctx)
	if err != nil {
		if utils.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Failed to read address configuration webhook: %s", err.Error())
	}

	d.Set("sid", getResponse.Sid)
	d.Set("account_sid", getResponse.AccountSid)
	d.Set("address", getResponse.Address)
	d.Set("service_sid", getResponse.AutoCreation.ConversationServiceSid)
	d.Set("friendly_name", getResponse.FriendlyName)
	d.Set("enabled", getResponse.AutoCreation.Enabled)
	d.Set("integration_type", getResponse.AutoCreation.Type)
	d.Set("type", getResponse.Type)
	d.Set("webhook_filters", getResponse.AutoCreation.WebhookFilters)
	d.Set("webhook_method", getResponse.AutoCreation.WebhookMethod)
	d.Set("webhook_url", getResponse.AutoCreation.WebhookUrl)
	d.Set("date_created", getResponse.DateCreated.Format(time.RFC3339))

	if getResponse.DateUpdated != nil {
		d.Set("date_updated", getResponse.DateUpdated.Format(time.RFC3339))
	}

	d.Set("url", getResponse.URL)

	return nil
}

func resourceConversationsAddressConfigurationWebhookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).Conversations

	webhookType := "webhook"
	updateInput := &address.UpdateAddressInput{
		AutoCreation: &address.UpdateAutoCreationInput{
			ConversationServiceSid: utils.OptionalStringWithEmptyStringOnChange(d, "service_sid"),
			Enabled:                utils.OptionalBool(d, "enabled"),
			Type:                   &webhookType,
			WebhookFilters:         utils.OptionalStringSlice(d, "webhook_filters"),
			WebhookMethod:          utils.OptionalString(d, "webhook_method"),
			WebhookUrl:             utils.OptionalString(d, "webhook_url"),
		},
		FriendlyName: utils.OptionalStringWithEmptyStringOnChange(d, "friendly_name"),
	}

	updateResp, err := client.Configuration().Address(d.Id()).UpdateWithContext(ctx, updateInput)
	if err != nil {
		return diag.Errorf("Failed to update address configuration webhook: %s", err.Error())
	}

	d.SetId(updateResp.Sid)
	return resourceConversationsAddressConfigurationWebhookRead(ctx, d, meta)
}

func resourceConversationsAddressConfigurationWebhookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).Conversations

	if err := client.Configuration().Address(d.Id()).DeleteWithContext(ctx); err != nil {
		return diag.Errorf("Failed to delete address configuration webhook: %s", err.Error())
	}
	d.SetId("")
	return nil
}
