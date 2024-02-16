package salesforce

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetDescription - Returns a specific description.
func (c *Client) GetDescription(sfObject string) (*Description, error) {

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s/services/data/%s/sobjects/%s/describe",
			c.HostURL,
			c.ApiVersion,
			sfObject,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	description := &Description{}
	err = json.Unmarshal(body, description)
	if err != nil {
		return nil, err
	}

	return description, nil
}
