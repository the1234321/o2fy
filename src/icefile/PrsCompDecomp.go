package icefile

type PrsCompDecomp struct {
	ctrlByteCounter int
	ctrlByte        byte //= 0;
	decompBuffer    []byte
	currDecompPos   int //= 0;

}

func (this *PrsCompDecomp) getCtrlBit() bool {
	this.ctrlByteCounter--
	if this.ctrlByteCounter == 0 {
		this.ctrlByte = this.decompBuffer[this.currDecompPos]
		this.currDecompPos++
		this.ctrlByteCounter = 8
	}
	var flag bool = (this.ctrlByte & 1) > 0
	this.ctrlByte >>= 1
	return flag
}

func Decompress(input []byte, outCount uint) []byte {
	Prs := new(PrsCompDecomp)
	return Prs.localDecompress(input, outCount)
}

func (this *PrsCompDecomp) localDecompress(input []byte, outCount uint) []byte {
	var outData []byte = make([]byte, outCount)
	this.decompBuffer = input
	this.ctrlByte = 0
	this.ctrlByteCounter = 1
	this.currDecompPos = 0
	var outIndex uint = 0

	for outIndex < outCount && this.currDecompPos < len(input) {
		for this.getCtrlBit() {
			outData[outIndex] = this.decompBuffer[this.currDecompPos]
			this.currDecompPos++
			outIndex++
		}
		var controlOffset int
		var controlSize int
		if this.getCtrlBit() {
			if this.currDecompPos < len(this.decompBuffer) {
				var data0 int = int(this.decompBuffer[this.currDecompPos])
				this.currDecompPos++
				var data1 int = int(this.decompBuffer[this.currDecompPos])
				this.currDecompPos++
				if data0 != 0 || data1 != 0 {
					controlOffset = (data1 << 5) + (data0 >> 3) - 8192
					var sizeTemp int = data0 & 7

					//controlSize = sizeTemp != 0 ? sizeTemp + 2 : decompBuffer[currDecompPos++] + 10;
					if sizeTemp != 0 {
						controlSize = sizeTemp + 2
					} else {
						controlSize = int(this.decompBuffer[this.currDecompPos]) + 10
						this.currDecompPos++
					}
				} else {
					break
				}
			} else {
				break
			}
		} else {
			controlSize = 2
			if this.getCtrlBit() {
				controlSize += 2
			}
			if this.getCtrlBit() {
				controlSize++
			}
			controlOffset = int(this.decompBuffer[this.currDecompPos]) - 256
			this.currDecompPos++
		}

		var loadIndex int = controlOffset + int(outIndex)
		for index := 0; index < controlSize && int(outIndex) < len(outData); index++ {
			outData[outIndex] = outData[loadIndex]
			outIndex++
			loadIndex++
		}
	}

	return outData
}

// public static byte[] compress(byte[] toCompress) => new PrsCompressor().compress(toCompress);
