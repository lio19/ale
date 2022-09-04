package youdao

type RespWordBook struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data RespWordBookData `json:"data"`
}

type RespWordBookData struct {
	Total    int                `json:"total"`
	ItemList []RespWordBookItem `json:"itemList"`
}

type RespWordBookItem struct {
	ItemId     string      `json:"itemId"`
	Word       string      `json:"word"`
	LanFrom    string      `json:"lanFrom"`
	LanTo      string      `json:"lanTo"`
	Trans      string      `json:"trans"`
	Usphone    interface{} `json:"usphone"`
	Ukphone    interface{} `json:"ukphone"`
	CreateTime string      `json:"createTime"`
}

type RespWord struct {
	ErrorCode   string   `json:"errorCode"`
	Query       string   `json:"query"`
	Translation []string `json:"translation"`
	Basic       struct {
		Phonetic   string   `json:"phonetic"`
		UkPhonetic string   `json:"uk-phonetic"`
		UsPhonetic string   `json:"us-phonetic"`
		UkSpeech   string   `json:"uk-speech"`
		UsSpeech   string   `json:"us-speech"`
		Explains   []string `json:"explains"`
	} `json:"basic"`
	Web []struct {
		Key   string   `json:"key"`
		Value []string `json:"value"`
	} `json:"web"`
	Dict struct {
		Url string `json:"url"`
	} `json:"dict"`
	Webdict struct {
		Url string `json:"url"`
	} `json:"webdict"`
	L         string `json:"l"`
	TSpeakUrl string `json:"tSpeakUrl"`
	SpeakUrl  string `json:"speakUrl"`
}
