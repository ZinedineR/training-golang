package externalapi

import (
	"csv-xlsx-read/entity"
	"csv-xlsx-read/httpclient"
)

type YugiohExternalImpl struct {
	HttpClient httpclient.Client
}

func NewYugiohExternalImpl(
	HttpClient httpclient.Client,
) YugiohSvcExternal {
	return &YugiohExternalImpl{
		HttpClient: HttpClient,
	}
}
func (b *YugiohExternalImpl) Get(archetype string) (*entity.YugiohAPIResponse, int, error) {
	var response *entity.YugiohAPIResponse
	urlPath := "https://db.ygoprodeck.com/api/v7/cardinfo.php?archetype=" + archetype

	headers := map[string]string{
		"Content-Type": "application/json",
		"accept":       "application/json",
	}

	statusCode, err := b.HttpClient.Get(urlPath, headers, &response)
	if err != nil {
		return response, statusCode, err
	}
	return response, statusCode, nil
}
