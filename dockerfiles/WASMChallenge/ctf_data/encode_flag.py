import os 

key = ord('s')
b=bytearray([key, *os.environ["FLAG"].encode()])
print(",".join(hex(b[i] ^ b[i-1]) for i in range(1,len(b))))