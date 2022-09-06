package youdao

import (
	"autoLearnEnglish/src/conf"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetOneWord(t *testing.T) {

	a := assert.New(t)

	conf.MockConf(map[string]string{
		//过期的cookie
		"cookie": "OUTFOX_SEARCH_USER_ID_NCOO=1677771653.2347355; OUTFOX_SEARCH_USER_ID=\"-1751478631@10.110.96.158\"; UM_distinctid=182fe46ac78618-0fe55b2d590bff-1b525635-13c680-182fe46ac79e4c; NTES_YD_SESS=OeXDxXgwUnrR8LolkZFeQOrnuMyt_VKiL2iW69oRgm46eypJekSXAhFF5MtspjaXcgXnUt4suVNiifZjS087hudsAYUaB1JpdqfTScUN0jOgqwCtMApZE7p9zH6F4BZOytsQBi_4U0xKDaF65f9xVOAqv85MwFkr519zw84QX4Ut99Gs1F2LUPnR48GnSJct.BbDN.xZelGZxLwPJA5J5J6NiBzgct.4Z; NTES_YD_PASSPORT=4EGZDdMpSK1F8AyEmg2xY7M8Glb2ckJwzmmiUaOA7QRxCIEACNjlM8xx4w2kEaGlqOl0L2fk.H1sn8BNXQNDnEdKbwcBziW2rkQNlPagsKwao2so7WhLL.t5JGu4eynj4Vokt.BFu9c62.vJvF6KBkvYFVM1OZ21eRVhLpFQXAEXuTlUjvEAQbZOlBrsv6irvOWAYUqoT4JUtJZRDBbtsksOG; S_INFO=1662123590|0|0&60##|17600460660; P_INFO=17600460660|1662123590|1|dict_logon|00&99|null&null&null#shd&370100#10#0|&0|null|17600460660; DICT_SESS=v2|r5bhPj3jARYE6Lqz0LYY0qLnMwF6MOE0eShHQ4Ofey0qLRfT4hfw40TF0fY5OLOERg4kMeFRH6B0wu6MkMhLlM0YfnHPLhMl5R; DICT_PERS=v2|urs-phone-web||DICT||web||604800000||1662123591205||2408:871a:3000:1::18c||urs-phoneyd.5f1f5487788e46329@163.com||QyPLkWhfQLRkY0MJuOfPK0PFP4PuOfU5RYEnM64RHp4RqukfY5nfzE0qFOLpZOMJuRUYhLTS0fqBRpBP4YA64zER; DICT_LOGIN=3||1662123591210",
		//这边需要填上sec和key，不然单元测试会不过
		"appSec": "",
		"appKey": "",
	})

	wb := NewWordBook()
	ab, err := wb.GetOneWord("abandon")
	a.Nil(err)
	a.NotNil(ab)
	a.Equal(ab.Query, "abandon")
	a.NotEqual(ab.TSpeakUrl, "")
	a.NotEqual(ab.SpeakUrl, "")
	fmt.Println(ab)
}
