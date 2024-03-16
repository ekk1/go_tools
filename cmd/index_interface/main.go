package main

import (
	"flag"
	"go_utils/utils"
	"go_utils/utils/nbd/backend"
	"go_utils/utils/nbd/server"
	"net"
)

func main() {
	size := flag.Int64("size", 1073741824, "NBD Size")
	laddr := flag.String("laddr", ":10809", "Listen address")
	name := flag.String("name", "default", "Export name")
	description := flag.String("description", "The default export", "Export description")
	readOnly := flag.Bool("read-only", false, "Whether the export should be read-only")
	minimumBlockSize := flag.Uint("minimum-block-size", 1, "Minimum block size")
	preferredBlockSize := flag.Uint("preferred-block-size", 4096, "Preferred block size")
	maximumBlockSize := flag.Uint("maximum-block-size", 0xffffffff, "Maximum blocksize")
	multiConn := flag.Bool("multi-conn", false, "Whether to advertise support for multiple simultaneous connections")

	flag.Parse()

	l, err := net.Listen("tcp", *laddr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	utils.LogPrintInfo("Listening on", l.Addr())

	b := backend.NewMemoryBackend(make([]byte, *size))

	clients := 0

	for {
		conn, err := l.Accept()
		if err != nil {
			utils.LogPrintError("Could not accept connection, continuing:", err)
			continue
		}
		clients++
		utils.LogPrintInfo(clients, "clients connected")

		go func() {
			defer func() {
				_ = conn.Close()
				clients--
				if err := recover(); err != nil {
					utils.LogPrintError("Client disconnected with error:", err)
				}
				utils.LogPrintInfo(clients, "clients connected")
			}()
			if err := server.Handle(
				conn,
				[]*server.Export{
					{
						Name:        *name,
						Description: *description,
						Backend:     b,
					},
				},
				&server.Options{
					ReadOnly:           *readOnly,
					MinimumBlockSize:   uint32(*minimumBlockSize),
					PreferredBlockSize: uint32(*preferredBlockSize),
					MaximumBlockSize:   uint32(*maximumBlockSize),
					SupportsMultiConn:  *multiConn,
				}); err != nil {
				panic(err)
			}
		}()
	}

}
