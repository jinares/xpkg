package xcrypt

import "hash/crc32"

//CRC32
func CRC32(data string) uint32 {
	return crc32.ChecksumIEEE([]byte(data))
}
