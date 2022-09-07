package cache

type WordListCache struct {
	AddTime int64
	Words   []string
}

type WordDetailCache struct {
	AddTime int64
	English string
	Chinese string
	EngMp3  string
	ChMp3   string
}

type InputWordCache struct {
	EngMp3Url string //english voice file http url
	ChMp3Url  string //chinese voice file http url
	WordEng   string //word english
	WordCh    string //word chinese
}

type innerWordFilePath struct {
	DirPath    string
	EngMp3Path string
	ChMp3Path  string
	InfoPath   string
}

type innerWordCache struct {
	WordEng string `json:"WordEng"`
	WordCh  string `json:"WordCh"`
	AddTime int64  `json:"AddTime"`
}
