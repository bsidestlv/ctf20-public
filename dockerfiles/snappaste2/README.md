# Snappaste (Part 2)

<img src="www/logo.png" alt="logo" width="32"> Here's an improved version of our pasting service. Some weird stuff was happening in the previous version, so we added an integrity check for the pastes. Also, due to misuse, we began limiting the size of the pasted snippets (please don't paste Kali Linux ISOs, thanks).

[BSidesTLV 2020](https://bsidestlv.com/) | [CTF](https://ctf20.bsidestlv.com/) | 750 points

# Compiling

## Visual Studio

Use the provided solution.

## GCC

Use the following commands:
```
gcc -c -std=c99 zlib/*.c
g++ -c -std=c++14 snappaste.cc 
g++ -o snappaste *.o -pthread
```
