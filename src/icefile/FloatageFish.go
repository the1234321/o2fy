package icefile

func decrypt_block2(data_block []byte, length uint, key uint) []byte {
	/*
		byte xor_byte = (byte)((( key >> 16 ) ^ key) & 0xFF);
		byte[] to_return = new byte[length];

		for ( uint i = 0; i < length; ++i )
		{
			if (data_block[i] != 0 && data_block[i] != xor_byte)
				to_return[i] = (byte)(data_block[i] ^ xor_byte);
			else
				to_return[i] = data_block[i];
		}*/
	return decrypt_block(data_block, length, key, 16)
}

func decrypt_block(data_block []byte, length uint, key uint, shift int) []byte {
	var xor_byte byte = (byte)(((key >> shift) ^ key) & 0xFF)
	to_return := make([]byte, length)

	for i := uint(0); i < length; i++ {
		if data_block[i] != 0 && data_block[i] != xor_byte {
			to_return[i] = byte(data_block[i] ^ xor_byte)
		} else {
			to_return[i] = data_block[i]
		}
	}
	return to_return
}
