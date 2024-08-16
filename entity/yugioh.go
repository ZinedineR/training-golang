package entity

type YugiohAPIResponse struct {
	Data []YugiohCard `json:"data"`
}

type YugiohCard struct {
	Id            int               `json:"id"`
	Name          string            `json:"name"`
	Type          string            `json:"type"`
	FrameType     string            `json:"frameType"`
	Desc          string            `json:"desc"`
	Atk           int               `json:"atk,omitempty"`
	Def           int               `json:"def,omitempty"`
	Level         int               `json:"level,omitempty"`
	Race          string            `json:"race"`
	Attribute     string            `json:"attribute,omitempty"`
	Archetype     string            `json:"archetype"`
	YgoprodeckUrl string            `json:"ygoprodeck_url"`
	CardSets      []YugiohCardSet   `json:"card_sets"`
	CardImages    []YugiohCardImage `json:"card_images"`
	CardPrices    []YugiohCardPrice `json:"card_prices"`
}

type YugiohCardSet struct {
	SetName       string `json:"set_name"`
	SetCode       string `json:"set_code"`
	SetRarity     string `json:"set_rarity"`
	SetRarityCode string `json:"set_rarity_code"`
	SetPrice      string `json:"set_price"`
}

type YugiohCardImage struct {
	Id              int    `json:"id"`
	ImageUrl        string `json:"image_url"`
	ImageUrlSmall   string `json:"image_url_small"`
	ImageUrlCropped string `json:"image_url_cropped"`
}

type YugiohCardPrice struct {
	CardmarketPrice   string `json:"cardmarket_price"`
	TcgplayerPrice    string `json:"tcgplayer_price"`
	EbayPrice         string `json:"ebay_price"`
	AmazonPrice       string `json:"amazon_price"`
	CoolstuffincPrice string `json:"coolstuffinc_price"`
}

type YugiohTemplate struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Desc          string `json:"desc"`
	Atk           int    `json:"atk,omitempty"`
	Def           int    `json:"def,omitempty"`
	Level         int    `json:"level,omitempty"`
	Race          string `json:"race"`
	Attribute     string `json:"attribute,omitempty"`
	Archetype     string `json:"archetype"`
	ImageUrl      string `json:"image_url"`
	YgoprodeckUrl string `json:"ygoprodeck_url"`
}
