package aws

import (
	//"time"

	//"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/hashicorp/errwrap"
	//"github.com/hashicorp/terraform/helper/resource"
)

type WAFToken struct {
	Connection *waf.WAF
	Region     string
	Token      string
}

func (t *WAFToken) Acquire() (*string, error) {
	awsMutexKV.Lock(t.Region)

	out, err := t.Connection.GetChangeToken(&waf.GetChangeTokenInput{})
	if err != nil {
		t.Release()
		return nil, errwrap.Wrapf("Failed to acquire change token: %s", err)
	}
	t.Token = *out.ChangeToken

	return out.ChangeToken, nil
}

func (t *WAFToken) Release() error {
	// out, err := t.Connection.GetChangeTokenStatus(&waf.GetChangeTokenStatusInput{
	// 	ChangeToken: aws.String(t.Token),
	// })
	// if err != nil {
	// 	return err
	// }
	// if *out.ChangeTokenStatus == "PROVISIONED" {
	// 	// Don't wait for token which wasn't used at all
	// 	awsMutexKV.Unlock(t.Region)
	// 	return nil
	// }

	// stateConf := resource.StateChangeConf{
	// 	Pending: []string{"PENDING"},
	// 	Target:  []string{"INSYNC"},
	// 	Timeout: 5 * time.Minute,
	// 	Refresh: func() (interface{}, string, error) {
	// 		out, err := t.Connection.GetChangeTokenStatus(&waf.GetChangeTokenStatusInput{
	// 			ChangeToken: aws.String(t.Token),
	// 		})
	// 		if err != nil {
	// 			return nil, "", err
	// 		}

	// 		return out, *out.ChangeTokenStatus, nil
	// 	},
	// }
	// _, err = stateConf.WaitForState()
	// if err != nil {
	// 	return err
	// }
	awsMutexKV.Unlock(t.Region)
	return nil
}

func newWAFToken(conn *waf.WAF, region string) *WAFToken {
	return &WAFToken{Connection: conn, Region: region}
}
