package fpmanagement

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nextlinux/enterprise-client-go/pkg/external"
	"github.com/nextlinux/fpmanagement/internal/client"
	"io/ioutil"
)

func AddCorrections(apiClient *client.EnterpriseClient) {
	fmt.Println("parsing corrections JSON file")
	correctionData, err := ioutil.ReadFile("corrections.json")
	if err != nil {
		panic(err)
	}

	var corrections []external.Correction
	err = json.Unmarshal(correctionData, &corrections)
	if err != nil {
		panic(err)
	}
	fmt.Printf("found %d corrections to add\n", len(corrections))

	ctx := apiClient.NewRequestContext(context.Background())

	for _, correction := range corrections {
		AddCorrection(apiClient, correction, ctx)
	}
}

func AddCorrection(apiClient *client.EnterpriseClient, correction external.Correction, authedCtx context.Context) {
	fmt.Printf("adding correction for package: %s\n", getPackageName(correction))
	_, httpResponse, err := apiClient.Client.DefaultApi.AddCorrection(authedCtx, correction, &external.AddCorrectionOpts{})
	if httpResponse != nil {
		defer httpResponse.Body.Close()
	}
	err = client.HandleAPIError(httpResponse, err, "unable to add correction")
	// Looks like the swagger spec is a little bit off of the actual response here (this doesn't mean the request failed though)
	if err != nil && err.Error() != "json: cannot unmarshal object into Go value of type []external.Correction" {
		panic(err)
	}
}

func getPackageName(correction external.Correction) string {
	packageName := "unknown"
	for _, fieldMatch := range correction.Match.FieldMatches {
		if fieldMatch.FieldName == "package" {
			packageName = fieldMatch.FieldValue
		}
	}
	return packageName
}
