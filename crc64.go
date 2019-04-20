// Package provides an crc variants which is used by Redis.
// * Specification of this CRC64 variant follows:
// * Name: crc-64-jones
// * Width: 64 bites
// * Poly: 0xad93d23594c935a9
// * Reflected In: True
// * Xor_In: 0xffffffffffffffff
// * Reflected_Out: True
// * Xor_Out: 0x0
// * Check("123456789"): 0xe9c6d914c4b8d9ca
package crc64

import (
	"hash"
	"hash/crc64"
)

var table = makeTable()

// We can make any crc variants by https://github.com/tpircher/pycrc .
func makeTable() *[8]crc64.Table {
	poly := uint64(0xad93d23594c935a9)
	table := new([8]crc64.Table)
	for i := 0; i < 256; i++ {
		c := uint8(i)
		v := uint64(0)
		for j := uint8(0x01); j&0xff != 0; j <<= 1 {
			bit := v&0x8000000000000000 != 0
			if c&j != 0 {
				bit = !bit
			}
			v <<= 1
			if bit {
				v ^= poly
			}

		}
		vv := v & 0x01
		for j := 1; j < 64; j++ {
			v >>= 1
			vv = (vv << 1) | (v & 0x01)
		}
		table[0][i] = vv ^ 0
	}
	// slice-by-8
	for i := 0; i < 256; i++ {
		v := table[0][i]
		for j := 1; j < 8; j++ {
			v = table[0][v&0xff] ^ (v >> 8)
			table[j][i] = v
		}
	}
	return table
}

// Checksum returns the CRC-64 checksum of data.
func Checksum(crc uint64, p []byte) uint64 {
	for len(p) > 8 {
		crc ^= uint64(p[0]) | uint64(p[1])<<8 | uint64(p[2])<<16 | uint64(p[3])<<24 |
			uint64(p[4])<<32 | uint64(p[5])<<40 | uint64(p[6])<<48 | uint64(p[7])<<56
		crc = table[7][crc&0xff] ^
			table[6][(crc>>8)&0xff] ^
			table[5][(crc>>16)&0xff] ^
			table[4][(crc>>24)&0xff] ^
			table[3][(crc>>32)&0xff] ^
			table[2][(crc>>40)&0xff] ^
			table[1][(crc>>48)&0xff] ^
			table[0][crc>>56]
		p = p[8:]
	}
	for _, v := range p {
		crc = table[0][byte(crc)^v] ^ (crc >> 8)
	}
	return crc
}

type digest struct {
	crc uint64
}

// New returns a crc64 hasher.
func New() hash.Hash64 {
	return &digest{}
}

// Write implements the io.Writer of hash.Hash
func (d *digest) Write(p []byte) (int, error) {
	d.crc = Checksum(d.crc, p)
	return len(p), nil
}

// Sum encodes the sum64 in little endian
func (d *digest) Sum(in []byte) []byte {
	s := d.Sum64()
	return append(in, byte(s), byte(s>>8), byte(s>>16), byte(s>>24),
		byte(s>>32), byte(s>>40), byte(s>>48), byte(s>>56))
}

// Sum64 returns a uint64 checksum
func (d *digest) Sum64() uint64 { return d.crc }

// BlockSize returns the hash block size
func (d *digest) BlockSize() int { return 1 }

// Size returns the checksum size
func (d *digest) Size() int { return 8 }

// Reset resets this hash
func (d *digest) Reset() { d.crc = 0 }
