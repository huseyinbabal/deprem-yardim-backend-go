package feed

import (
	"fmt"
	"github.com/acikkaynak/backend-api-go/e2e"
	"github.com/acikkaynak/backend-api-go/feeds"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type FeedTestSuite struct {
	suite.Suite
	e2e.TestSuite
}

func (f *FeedTestSuite) SetupSuite() {
	err := f.SetupBaseSuite()
	f.NoError(err)
}

func (f *FeedTestSuite) TearDownSuite() {
	err := f.TearDownBaseSuite()
	f.NoError(err)
}

func TestAliasTestSuite(t *testing.T) {
	suite.Run(t, new(FeedTestSuite))
}

func (f *FeedTestSuite) Test_Should_List_Feeds() {
	feeds, err := f.getFeeds()
	f.NoError(err)
	f.NotNil(feeds)
}

func (f *FeedTestSuite) getFeeds() (feeds.Response, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/areas/feeds", f.Client.BaseURL), nil)
	if err != nil {
		return feeds.Response{}, fmt.Errorf("error listing feeds: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	var resp feeds.Response
	if err := f.Client.MakeRequest(req, &resp); err != nil {
		return feeds.Response{}, err
	}

	return resp, nil
}
