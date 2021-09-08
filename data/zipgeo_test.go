package data_test

import (
	"context"
	"testing"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/danesparza/zipgeo/data"
)

func TestZipGeo_GetLatLong_ReturnsValidData(t *testing.T) {

	//	Arrange
	service := data.ZipGeoService{}
	zipcode := "30019"
	ctx := context.Background()
	ctx, seg := xray.BeginSegment(ctx, "unit-test")
	defer seg.Close(nil)

	//	Act
	response, err := service.GetLatLong(ctx, zipcode)

	//	Assert
	if err != nil {
		t.Errorf("Error calling GetLatLong: %v", err)
	}

	t.Logf("Returned object: %+v", response)

}
