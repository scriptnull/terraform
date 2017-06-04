package github

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/github"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGithubUserSSHKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceGithubUserSSHKeyCreateOrUpdate,
		Read:   resourceGithubUserSSHKeyRead,
		Update: resourceGithubUserSSHKeyCreateOrUpdate,
		Delete: resourceGithubUserSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceGithubUserSSHKeyCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Organization).client
	title := d.Get("title").(string)
	key := d.Get("key").(string)
	// keys, err := client.Users.ListKeys(context.TODO(), "", nil)
	// if err

	respKey, resp, err := client.Users.CreateKey(context.TODO(), &github.Key{
		Title: &title,
		Key:   &key,
	})
	if err != nil {
		return err
	}
	fmt.Printf("%v \n %v \n", respKey, resp)
	log.Printf("[DEBUG] %v %v", respKey, resp)
	// d.SetId(strconv.Itoa(key.GetID()))

	return nil
}

func resourceGithubUserSSHKeyRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceGithubUserSSHKeyDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
