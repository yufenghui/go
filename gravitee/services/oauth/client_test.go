package oauth_test

import "github.com/yufenghui/go/gravitee/models"

func (suite *OauthTestSuite) TestFindClientByClientID() {
	var (
		client *models.OauthClient
		err    error
	)

	// When we try to find a client with a bogus client ID
	client, err = suite.service.FindClientByClientID("bogus")

	if err != nil {
		suite.T().Logf("error: %s", err)
	} else {
		suite.T().Logf("client: %p", client)
	}

}
