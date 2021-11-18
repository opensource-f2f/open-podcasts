package broadcast

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/url"
	"strings"
	"time"
)

func Search(keyword string) (result *AlbumSearchResult, err error) {
	keyword = url.QueryEscape(keyword)
	apiurl := fmt.Sprintf("https://m.ximalaya.com/m-revision/page/search?kw=%s&core=album&page=1&rows=10", keyword)
	resp, err := HttpGet(apiurl, Android)
	if err != nil {
		return nil, fmt.Errorf("无法获取专辑信息: %v", err)
	}
	defer resp.Body.Close()

	result = &AlbumSearchResult{}
	err = jsoniter.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, fmt.Errorf("无法获取专辑信息: 无法解析Json: %v", err)
	}
	return
}

func GetAlbumInfo(albumID int) (ai *AlbumInfo, err error) {
	url := fmt.Sprintf("http://mobile.ximalaya.com/mobile-album/album/page/ts-%d?ac=WIFI&albumId=%d&device=android&pageId=1&pageSize=0",
		time.Now().Unix(), albumID)
	resp, err := HttpGet(url, Android)
	if err != nil {
		return nil, fmt.Errorf("无法获取专辑信息: %v", err)
	}
	defer resp.Body.Close()

	ai = &AlbumInfo{}
	err = jsoniter.NewDecoder(resp.Body).Decode(ai)
	if err != nil {
		return nil, fmt.Errorf("无法获取专辑信息: 无法解析Json: %v", err)
	}
	return ai, nil
}

func GetVipAudioInfo(trackId int, cookie string) (ai *AudioInfo, err error) {
	ts := time.Now().Unix()
	url := fmt.Sprintf(
		"https://mpay.ximalaya.com/mobile/track/pay/%d/%d?device=pc&isBackend=true&_=%d",
		trackId, ts, ts)

	resp, err := HttpGetByCookie(url, cookie, Android)
	if err != nil {
		return ai, fmt.Errorf("获取音频信息失败: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ai, err
	}

	var vipAi VipAudioInfo
	err = jsoniter.Unmarshal(body, &vipAi)
	if err != nil {
		return ai, fmt.Errorf("无法解析JSON: %v", err)
	}

	if vipAi.Ret != 0 {
		return ai, fmt.Errorf("无法获取VIP音频信息: %v", jsoniter.Get(body, "msg").ToString())
	}

	fileName := DecryptFileName(vipAi.Seed, vipAi.FileID)
	sign, _, token, timestamp := DecryptUrlParams(vipAi.Ep)

	args := fmt.Sprintf("?sign=%s&buy_key=%s&token=%d&timestamp=%d&duration=%d",
		sign, vipAi.BuyKey, token, timestamp, vipAi.Duration)

	ai = &AudioInfo{TrackID: trackId, Title: vipAi.Title}

	ai.PlayPathAacv164 = vipAi.Domain + "/download/" + vipAi.APIVersion + fileName + args
	return ai, nil
}

func GetAudioInfo(albumID, page, pageSize int) (audioList []AudioInfo, err error) {
	format := fmt.Sprintf("https://m.ximalaya.com/m-revision/common/album/queryAlbumTrackRecordsByPage?albumId=%d&page=%d&pageSize=%d&asc=true", albumID, page, pageSize)

	resp, err := client.Get(format)
	if err != nil {
		return nil, fmt.Errorf("http get %v fail:%v", format, err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	list := jsoniter.Get(body, "data").Get("trackDetailInfos")
	for i2 := 0; i2 < list.Size(); i2++ {
		v := list.Get(i2).Get("trackInfo")
		audioList = append(audioList, AudioInfo{TrackID: v.Get("id").ToInt(),
			PlayPathAacv164: v.Get("playPath").ToString(), Title: v.Get("title").ToString()})
	}

	return audioList, nil
}

func GetAllAudioInfo(albumID int) (list []*AudioInfo, err error) {
	firstPlayList, err := GetAudioInfoListByPageID(albumID, 0)
	if err != nil {
		return nil, fmt.Errorf("无法获取播放列表: 0, %s", err)
	}
	for _, v := range firstPlayList.List {
		list = append(list, v)
	}

	for i := 1; i < firstPlayList.MaxPageID; i++ {
		playList, err := GetAudioInfoListByPageID(albumID, i)
		if err != nil {
			return nil, fmt.Errorf("无法获取播放列表: %d, %s", i, err)
		}
		for _, v := range playList.List {
			list = append(list, v)
		}
	}
	return list, nil
}

func GetAudioInfoListByPageID(albumID, pageID int) (playlist *Playlist, err error) {
	url := fmt.Sprintf("http://mobwsa.ximalaya.com/mobile/playlist/album/page?albumId=%d&pageId=%d",
		albumID, pageID)
	resp, err := HttpGet(url, Android)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	playlist = &Playlist{}
	err = jsoniter.Unmarshal(data, playlist)
	if err != nil {
		return nil, err
	}

	return playlist, nil
}

func GetTrackList(albumID, pageID int, isAsc bool) (tracks *TrackList, err error) {
	url := fmt.Sprintf(
		"https://mobile.ximalaya.com/mobile/v1/album/track/ts-%d?ac=WIFI&albumId=%d&device=android&isAsc=%t&pageId=%d&pageSize=200",
		time.Now().Unix(), albumID, isAsc, pageID)
	resp, err := HttpGet(url, Android)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	tracks = &TrackList{}
	err = jsoniter.NewDecoder(resp.Body).Decode(tracks)
	if err != nil {
		return nil, err
	}
	return tracks, nil
}

func GetUserInfo(cookie string) (*UserInfo, error) {
	resp, err := HttpGetByCookie("https://www.ximalaya.com/revision/main/getCurrentUser", cookie, PC)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ui := &UserInfo{}
	err = jsoniter.Unmarshal(data, ui)
	if err != nil {
		return nil, err
	}

	return ui, nil
}

func GetQRCode() (qrCode *QRCode, err error) {
	resp, err := HttpGet("https://passport.ximalaya.com/web/qrCode/gen?level=L", PC)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	qrCode = &QRCode{}
	err = jsoniter.Unmarshal(data, qrCode)
	if err != nil {
		return nil, err
	}
	return qrCode, err
}

func CheckQRCodeStatus(qrID string) (status *QRCodeStatus, cookie string, err error) {
	url := fmt.Sprintf("https://passport.ximalaya.com/web/qrCode/check/%s/%d", qrID, time.Now().Unix())
	resp, err := HttpGet(url, PC)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	status = &QRCodeStatus{}
	err = jsoniter.Unmarshal(data, status)
	if err != nil {
		return nil, "", err
	}
	if status.Ret == 0 {
		token := resp.Header.Values("Set-Cookie")[1]
		index := strings.Index(token, ";")
		return status, token[0:index], nil
	}

	return status, "", nil
}
