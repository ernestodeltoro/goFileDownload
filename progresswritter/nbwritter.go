package progresswritter

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
)

// NonBlokingProgressWriter counts the number of bytes written to it. It implements to the io.Writer
// interface and we can pass this into io.TeeReader() which will report progress on each
// write cycle.
type NonBlokingProgressWriter struct {
	currentWritten uint64
	fullSize       uint64
	newWrite       chan bool
}

// NewNonBloking returns a Writer interface that allows to show the
// writing progress in percentage given the fullSize of the file is known
func NewNonBloking(fullSize uint64) *NonBlokingProgressWriter {
	pw := NonBlokingProgressWriter{
		currentWritten: 0,
		fullSize:       fullSize,
		newWrite:       make(chan bool),
	}

	go pw.serveUpdateChannel()

	return &pw
}

func (pw *NonBlokingProgressWriter) serveUpdateChannel() {

	defer close(pw.newWrite)

	for range pw.newWrite {
		pw.updateProgress()
	}
}

// DoneCh returns the done channel
/*func (pw *NonBlokingProgressWriter) DoneCh() chan bool {
	return pw.done
}*/

// Non Blocking write on a channel
func (pw *NonBlokingProgressWriter) execUpdateNonBlocking() {
	select {
	case pw.newWrite <- true:
		return
	default:
	}
}

func (pw *NonBlokingProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.currentWritten += uint64(n)
	pw.execUpdateNonBlocking()
	return n, nil
}

// updateProgress prints
func (pw *NonBlokingProgressWriter) updateProgress() {

	fmt.Printf("\r%s", strings.Repeat(" ", 40))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete of %s", humanize.Bytes(pw.currentWritten), humanize.Bytes(pw.fullSize))

}
