package securechan

import "errors"

// LZ4 block format (https://github.com/lz4/lz4/blob/dev/doc/lz4_Block_format.md), compressor +
// decompressor. STOPGAP (CLAUDE.md "Core Deps Are PARENA-First"): written in Go to unblock the
// secure channel; the PARENA-native replacement is a tracked item (PARENA/stdlib/compress/lz4.prn
// is LZ4-*style* only, not wire-compatible, O(n^2), not streaming). Real LZ4 block output, so a
// PARENA implementation can be checked against this one byte for byte.

const (
	lz4MinMatch   = 4
	lz4HashLog    = 14
	lz4LastLits   = 5  // the last 5 bytes are always literals
	lz4MFLimit    = 12 // a match must start at least 12 bytes before the end
	lz4MaxOffset  = 65535
	errCorruptLZ4 = "securechan: corrupt lz4 block"
)

func lz4Hash(v uint32) uint32 { return (v * 2654435761) >> (32 - lz4HashLog) }

func le32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func lz4WriteLen(dst []byte, n int) []byte {
	for n >= 255 {
		dst = append(dst, 255)
		n -= 255
	}
	return append(dst, byte(n))
}

func lz4EmitSeq(dst, src []byte, litStart, litEnd, offset, matchLen int, last bool) []byte {
	lit := litEnd - litStart
	tok := byte(0)
	if lit >= 15 {
		tok = 15 << 4
	} else {
		tok = byte(lit) << 4
	}
	ml := 0
	if !last {
		ml = matchLen - lz4MinMatch
		if ml >= 15 {
			tok |= 15
		} else {
			tok |= byte(ml)
		}
	}
	dst = append(dst, tok)
	if lit >= 15 {
		dst = lz4WriteLen(dst, lit-15)
	}
	dst = append(dst, src[litStart:litEnd]...)
	if last {
		return dst
	}
	dst = append(dst, byte(offset), byte(offset>>8))
	if ml >= 15 {
		dst = lz4WriteLen(dst, ml-15)
	}
	return dst
}

// lz4Compress appends the LZ4 block encoding of src to dst. Greedy, single hash table.
func lz4Compress(dst, src []byte) []byte {
	n := len(src)
	if n < lz4MFLimit+1 {
		return lz4EmitSeq(dst, src, 0, n, 0, 0, true)
	}
	var table [1 << lz4HashLog]int32 // position+1; 0 = empty
	anchor, i := 0, 0
	limit := n - lz4MFLimit
	for i < limit {
		h := lz4Hash(le32(src[i:]))
		cand := int(table[h]) - 1
		table[h] = int32(i + 1)
		if cand < 0 || i-cand > lz4MaxOffset || le32(src[cand:]) != le32(src[i:]) {
			i++
			continue
		}
		m := lz4MinMatch
		for i+m < n-lz4LastLits && src[cand+m] == src[i+m] {
			m++
		}
		dst = lz4EmitSeq(dst, src, anchor, i, i-cand, m, false)
		i += m
		anchor = i
	}
	return lz4EmitSeq(dst, src, anchor, n, 0, 0, true)
}

// lz4Decompress decodes a block into exactly outLen bytes. outLen is attacker-influenced, so the
// caller must bound it before calling; every read and write here is bounds-checked.
func lz4Decompress(src []byte, outLen int) ([]byte, error) {
	dst := make([]byte, 0, outLen)
	bad := errors.New(errCorruptLZ4)
	i := 0
	for i < len(src) {
		tok := src[i]
		i++
		lit := int(tok >> 4)
		if lit == 15 {
			for {
				if i >= len(src) {
					return nil, bad
				}
				b := src[i]
				i++
				lit += int(b)
				if lit > outLen {
					return nil, bad
				}
				if b != 255 {
					break
				}
			}
		}
		if i+lit > len(src) || len(dst)+lit > outLen {
			return nil, bad
		}
		dst = append(dst, src[i:i+lit]...)
		i += lit
		if i == len(src) { // last sequence: literals only
			break
		}
		if i+2 > len(src) {
			return nil, bad
		}
		off := int(src[i]) | int(src[i+1])<<8
		i += 2
		if off == 0 || off > len(dst) {
			return nil, bad
		}
		ml := int(tok & 15)
		if ml == 15 {
			for {
				if i >= len(src) {
					return nil, bad
				}
				b := src[i]
				i++
				ml += int(b)
				if ml > outLen {
					return nil, bad
				}
				if b != 255 {
					break
				}
			}
		}
		ml += lz4MinMatch
		if len(dst)+ml > outLen {
			return nil, bad
		}
		for k := 0; k < ml; k++ { // byte-wise: overlapping matches (RLE) are legal
			dst = append(dst, dst[len(dst)-off])
		}
	}
	if len(dst) != outLen {
		return nil, bad
	}
	return dst, nil
}
