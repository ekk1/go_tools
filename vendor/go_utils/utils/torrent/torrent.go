package torrent

import (
	"errors"
	"go_utils/utils"
	"go_utils/utils/bencode"
	"os"
	"strings"
)

type Torrent struct {
	Announce     string
	AnnounceList []string
	Comment      string
	URL_List     []string
	CreatedBy    string
	CreateDate   int64
	Info         *TorrentInfo
}

type TorrentInfo struct {
	Length      int64
	PieceLength int64
	Private     int64
	Name        string
	Source      string
	Pieces      [][]byte
	Files       []*TorrentFile
}

type TorrentFile struct {
	Length int64
	Path   string
}

type BencodeTypes interface {
	int64 | []byte | []interface{} | map[string]interface{}
}

const failedMsg = "Failed to parse torrent file"

func decodeTorrentData[T BencodeTypes](data interface{}) T {
	dataValid := false
	switch data.(type) {
	case T:
		dataValid = true
	}
	if dataValid {
		return data.(T)
	}
	utils.LogPrintWarning("Failed to decode data")
	var emptyValue T
	return emptyValue
}

func ReadTorrentFile(filename string) (*Torrent, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.New("torrent: failed to read torrent file")
	}
	d := &bencode.BencodeDecoder{}
	ret, err := d.Decode(data)
	if err != nil {
		return nil, err
	}

	t := &Torrent{}

	retTrue, ok := ret.(map[string]interface{})
	if !ok {
		return nil, errors.New(failedMsg)
	}

	for topk, topv := range retTrue {
		utils.LogPrintDebug("Decoding: ", topk)
		switch topk {
		case "announce":
			t.Announce = string(decodeTorrentData[[]byte](topv))
		case "announce-list":
			t.AnnounceList = []string{}
			anListRaw := decodeTorrentData[[]interface{}](topv)
			for _, anRaw := range anListRaw {
				anBytes := decodeTorrentData[[]byte](anRaw)
				t.AnnounceList = append(t.AnnounceList, string(anBytes))
			}
		case "comment":
			t.Comment = string(decodeTorrentData[[]byte](topv))
		case "created by":
			t.CreatedBy = string(decodeTorrentData[[]byte](topv))
		case "creation date":
			t.CreateDate = decodeTorrentData[int64](topv)
		case "info":
			torrentInfo, err := decodeTorrentInfo(topv)
			if err != nil {
				return nil, errors.New("torrent: failed to parse torrent")
			}
			t.Info = torrentInfo
		default:
			utils.LogPrintWarning("Unknown field: ", topk)
		}
	}

	return t, nil
}

func decodeTorrentInfo(info interface{}) (*TorrentInfo, error) {
	torrentInfo := &TorrentInfo{}
	infoDict := decodeTorrentData[map[string]interface{}](info)
	for infoK, infoV := range infoDict {
		utils.LogPrintDebug("Decoding info-", infoK)
		switch infoK {
		case "files":
			fileListResult := []*TorrentFile{}
			fileList := decodeTorrentData[[]interface{}](infoV)
			for _, fileInfoRaw := range fileList {
				fileInfo, err := decodeTorrentFile(fileInfoRaw)
				if err != nil {
					return nil, errors.New("torrent: failed to parse torrent")
				}
				fileListResult = append(fileListResult, fileInfo)
			}
			torrentInfo.Files = fileListResult
		case "name":
			torrentInfo.Name = string(decodeTorrentData[[]byte](infoV))
		case "piece length":
			torrentInfo.PieceLength = decodeTorrentData[int64](infoV)
		case "pieces":
			torrentInfo.Pieces = [][]byte{}
			hashData := decodeTorrentData[[]byte](infoV)
			for i := 0; i < len(hashData)/20; i++ {
				torrentInfo.Pieces = append(
					torrentInfo.Pieces,
					hashData[i*20:i*20+20],
				)
			}
		case "private":
			torrentInfo.Private = decodeTorrentData[int64](infoV)
		case "source":
			torrentInfo.Source = string(decodeTorrentData[[]byte](infoV))
		default:
			utils.LogPrintWarning("Unknown info field: ", infoK)
		}
	}
	return torrentInfo, nil
}

func decodeTorrentFile(fileInfoRaw interface{}) (*TorrentFile, error) {
	fileInfo := &TorrentFile{}
	fileDict := decodeTorrentData[map[string]interface{}](fileInfoRaw)
	if fileDict == nil {
		return nil, errors.New("torrent: failed to parse torrent")
	}
	for k, v := range fileDict {
		switch k {
		case "path":
			pathResult := ""
			pathList := decodeTorrentData[[]interface{}](v)
			for _, path := range pathList {
				pathResult += string(decodeTorrentData[[]byte](path)) + "/"
			}
			pathResult = strings.TrimSuffix(pathResult, "/")
			fileInfo.Path = pathResult
		case "length":
			fileInfo.Length = decodeTorrentData[int64](v)
		}
	}
	return fileInfo, nil
}
