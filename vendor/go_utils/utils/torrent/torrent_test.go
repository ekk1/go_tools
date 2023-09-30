package torrent

import (
	"fmt"
	"go_utils/utils"
	"testing"
)

func TestParseTorrent(t *testing.T) {
	utils.SetLogLevelDebug()

	tt, err := ReadTorrentFile("/path/to.torrent")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(tt.Announce)
	// fmt.Println(tt.AnnounceList)
	fmt.Println(tt.Comment)
	// fmt.Println(tt.URL_List)
	fmt.Println(tt.CreateDate)
	fmt.Println(tt.CreatedBy)
	// fmt.Println(tt.Info)
	// fmt.Println(tt.Info.Pieces)

}
