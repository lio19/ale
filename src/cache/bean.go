package cache

type WordList struct {
	AddTime int64
	Words   []string
}

type Word struct {
	AddTime int64
	English string
	Chinese string
	EngMp3  string
	ChMp3   string
}

type WordFilePath struct {
	DirPath    string
	EngMp3Path string
	ChMp3Path  string
	InfoPath   string
}

type WordCache struct {
	EngMp3Url string //english voice file http url
	ChMp3Url  string //chinese voice file http url
	WordEng   string //word english
	WordCh    string //word chinese
}

type WordInfoCache struct {
	WordEng string `json:"WordEng"`
	WordCh  string `json:"WordCh"`
	AddTime int64  `json:"AddTime"`
}
