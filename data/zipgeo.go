package data

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/tidwall/buntdb"
)

// ZipGeoService is a lookup service for zipcode->lat,long data
type ZipGeoService struct{}

// GetWeatherReport gets the weather report
func (s ZipGeoService) GetLatLong(ctx context.Context, zipcode string) (ZipGeo, error) {

	//	Start the service segment
	ctx, seg := xray.BeginSubsegment(ctx, "zipgeo-service")

	//	Our return value
	retval := ZipGeo{}

	//	Parse the zipcode:
	zip, err := strconv.Atoi(zipcode)
	if err != nil {
		seg.Close(err)
		return retval, fmt.Errorf("problem parsing the zipcode: %v", err)
	}
	retval.Zipcode = zip

	dbname := "zipgeo.db"
	sysdb, err := buntdb.Open(dbname)
	if err != nil {
		seg.Close(err)
		return retval, fmt.Errorf("problem opening the zipgeo.db: %v", err)
	}
	defer sysdb.Close()

	//	Create our indexes
	sysdb.CreateIndex("zip", "zip:*", buntdb.IndexString)

	//	Fetch our latlong
	latlong := ""
	err = sysdb.View(func(tx *buntdb.Tx) error {
		latlong, err = tx.Get(GetKey("zip", zipcode))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		seg.Close(err)
		return retval, fmt.Errorf("problem executing get in zipgeo.db: %v", err)
	}

	//	Parse the latlong into our return value
	latlongarray := strings.Split(latlong, ",")

	//	Sanity check that we parsed the string properly
	if len(latlongarray) != 2 {
		return retval, fmt.Errorf("couldn't find the lat long for zipcode: %v", zipcode)
	}

	//	Extract latitude
	latitude, err := strconv.ParseFloat(latlongarray[0], 64)
	if err != nil {
		seg.Close(err)
		return retval, fmt.Errorf("problem parsing latitude: %v", err)
	}
	retval.Latitude = latitude

	//	Extract longitude
	longitude, err := strconv.ParseFloat(latlongarray[1], 64)
	if err != nil {
		seg.Close(err)
		return retval, fmt.Errorf("problem parsing longitude: %v", err)
	}
	retval.Longitude = longitude

	//	Add the result to the request metadata
	xray.AddMetadata(ctx, "ZipGeoResult", retval)

	// Close the segment
	seg.Close(nil)

	//	Return the report
	return retval, nil
}

// GetKey returns a key to be used in the storage system
func GetKey(entityType string, keyPart ...string) string {
	allparts := []string{}
	allparts = append(allparts, entityType)
	allparts = append(allparts, keyPart...)
	return strings.Join(allparts, ":")
}
