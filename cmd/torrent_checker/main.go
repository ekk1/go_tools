package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"go_utils/utils"
	"go_utils/utils/torrent"
	"io"
	"os"
)

func getBytesHash(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

func getReadPositionByTorrentOffset(tt *torrent.Torrent, offset int64) (string, int64, error) {
	var currentCursor int64 = 0

	for _, file := range tt.Info.Files {
		fileFullPath := tt.Info.Name + "/" + file.Path

		endCursor := currentCursor + file.Length

		if offset >= currentCursor && offset < endCursor {
			return fileFullPath, offset - currentCursor, nil
		}

		currentCursor += file.Length
	}

	return "", 0, errors.New("Failed to find file")
}

func main() {
	utils.SetLogLevelInfo()

	filename := os.Args[1]

	tt, err := torrent.ReadTorrentFile(filename)
	if err != nil {
		utils.ErrExit(err)
	}

	utils.LogPrintInfo("Performing check on: ", tt.Info.Name)
	utils.LogPrintInfo("Block size: ", tt.Info.PieceLength)

	var currentCursor int64 = 0

	totalBlocks := int64(len(tt.Info.Pieces))

	var successBlocks int64 = 0
	var failedBlocks int64 = 0
	var missingBlocks int64 = 0
	var presentBlocks int64 = 0

	var blockCanBeVerified map[int64]bool = map[int64]bool{}
	for _, file := range tt.Info.Files {
		fileFullPath := tt.Info.Name + "/" + file.Path

		currentCursorBlock := currentCursor / tt.Info.PieceLength
		endCursor := currentCursor + file.Length
		endCursorBlock := endCursor / tt.Info.PieceLength

		_, err := os.Stat(fileFullPath)
		if os.IsNotExist(err) {
			for i := currentCursorBlock; i <= endCursorBlock; i++ {
				if _, ok := blockCanBeVerified[i]; !ok {
					missingBlocks += 1
				}
				blockCanBeVerified[i] = false
			}
		} else {
			for i := currentCursorBlock; i <= endCursorBlock; i++ {
				if _, ok := blockCanBeVerified[i]; !ok {
					blockCanBeVerified[i] = true
				}
			}
		}
		currentCursor += file.Length
	}

	utils.LogPrintInfo("Torrent records", totalBlocks, "blocks")
	utils.LogPrintInfo("File spans", len(blockCanBeVerified), "blocks")

	for i := int64(0); i <= totalBlocks; i++ {
		if blockCanBeVerified[i] {
			// fmt.Println(i)
			presentBlocks += 1
		}
	}

	utils.LogPrintInfo("Present blocks:\t", presentBlocks)

	for i := int64(0); i <= totalBlocks; i++ {
		if blockCanBeVerified[i] {
			utils.LogPrintDebug(fmt.Sprintf("Block %v can be verified, start verify", i))

			currentSeek := i * tt.Info.PieceLength
			bytesToRead := tt.Info.PieceLength

			buf := bytes.NewBuffer(nil)

			for bytesToRead != 0 {
				filename, startPos, err := getReadPositionByTorrentOffset(tt, currentSeek)
				utils.ErrExit(err)
				f, err := os.Open(filename)
				utils.ErrExit(err)
				stat, _ := os.Stat(filename)
				_, err = f.Seek(startPos, 0)
				if err != nil {
					f.Close()
					utils.LogPrintError(err)
					return
				}
				fileRemainBytes := stat.Size() - startPos
				if fileRemainBytes >= bytesToRead {
					_, err := io.CopyN(buf, f, bytesToRead)
					if err != nil {
						f.Close()
						utils.LogPrintError(err)
						return
					}
					bytesToRead = 0
				} else {
					_, err := io.CopyN(buf, f, fileRemainBytes)
					if err != nil {
						f.Close()
						utils.LogPrintError(err)
						return
					}
					if i+1 == totalBlocks {
						bytesToRead = 0
					} else {
						bytesToRead -= fileRemainBytes
						currentSeek += fileRemainBytes
					}
				}
				f.Close()
			}

			hash := getBytesHash(buf.Bytes())

			if !bytes.Equal(hash, tt.Info.Pieces[i]) {
				utils.LogPrintWarning(fmt.Sprintf(
					"Hash check failed for block %v, got: %v, expect %v",
					i,
					hex.EncodeToString(hash),
					hex.EncodeToString(tt.Info.Pieces[i]),
				))
				failedBlocks += 1
			} else {
				successBlocks += 1
			}
			fmt.Printf("\rChecking: [%v/%v], failed: %v", successBlocks, presentBlocks, failedBlocks)
		}
	}

	fmt.Println()
	utils.LogPrintInfo("Success blocks:\t", successBlocks)
	utils.LogPrintInfo("Failure blocks:\t", failedBlocks)
}
