package icefile

type Crc32Alt struct {
	_table [0x100]uint
	//= 0xFFFFFFFF;
	DefaultSeed uint
}

func (this *Crc32Alt) Crc32Alt() {
	this.DefaultSeed = 0xFFFFFFFF

	for i := uint(0); i < 0x100; i++ {
		var crc uint = i
		for j := 0; j < 8; j++ {
			// crc = (crc >> 1) ^ ((crc & 1) > 0 ? 0xEDB88320 : 0);
			if (crc & 1) > 0 {
				crc = (crc >> 1) ^ 0xEDB88320
			} else {
				crc = (crc >> 1) ^ 0
			}
		}
		this._table[i] = crc
	}
}

func (this *Crc32Alt) GetCrc32(data []byte, crc uint) uint {
	var start int = 0
	var length int = len(data)
	crc ^= this.DefaultSeed

	if len(data) > 0 {
		for ; length > 0; length-- {
			crc = this._table[(uint(data[start])^crc)&0xFF] ^ (crc >> 8)
			start++
		}
	}

	return crc ^ this.DefaultSeed
}
