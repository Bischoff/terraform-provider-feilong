/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package feilong

import (
)


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#create-token

func (c *Client) CreateToken(adminToken string) error {

	c.Token = &adminToken

	_, err := c.doRequest("POST", "/token", nil)

	return err
}
