package externalapi

import "csv-xlsx-read/entity"

type YugiohSvcExternal interface {
	Get(archetype string) (*entity.YugiohAPIResponse, int, error)
}
