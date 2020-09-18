package xcrypt

import (
	"fmt"
	"github.com/jinares/xpkg/xtools"
	"testing"
)

func TestCRC322(t *testing.T) {
	id := xtools.GUID()
	fmt.Println(CRC32(MD5(id)))
}
