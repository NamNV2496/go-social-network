package service

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/namnv2496/user-service/internal/configs"
	"github.com/namnv2496/user-service/internal/entity"
)

type ILocation interface {
	GetLocation() []*entity.LocationInfo
}

type Location struct {
	location []*entity.LocationInfo
}

func NewLocation(
	conf *configs.Config,
) ILocation {
	records, err := readCsv(conf.Location.FilePath)
	if err != nil {
		return nil
	}
	return &Location{
		location: mapLocationInfo(records),
	}
}

func readCsv(filePath string) ([][]string, error) {
	f, err := os.Open("./data/" + filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func mapLocationInfo(records [][]string) []*entity.LocationInfo {
	//      0             1           2           3.         4            5            6.       7
	// Tên tỉnh cũ	| Mã TP cũ	| Tên QH cũ	| Tên PX cũ	| Mã QH cũ	| Mã PX cũ	| Tên PX mới	| Mã PX mới
	locations := make([]*entity.LocationInfo, 0)
	for i := 1; i < len(records); i++ {
		record := records[i]
		buildName := fmt.Sprintf("%s %s %s", record[3], record[2], record[0])
		buildName = strings.ReplaceAll(buildName, "xã", "")
		buildName = strings.ReplaceAll(buildName, "phường", "")
		buildName = strings.ReplaceAll(buildName, "quận", "")
		buildName = strings.ReplaceAll(buildName, "huyện", "")
		buildName = strings.ReplaceAll(buildName, "thành phố", "")
		locations = append(locations, &entity.LocationInfo{
			CityV1:       record[1],
			CityV1Name:   record[0],
			District:     record[4],
			DistrictName: record[2],
			WardV1:       record[5],
			WardV1Name:   record[3],
			WardV2:       record[7],
			WardV2Name:   record[6],
			BuilderName:  buildName,
		})
	}
	return locations
}

func (_self *Location) GetLocation() []*entity.LocationInfo {
	return _self.location
}
