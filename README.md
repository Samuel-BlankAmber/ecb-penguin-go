# ECB Penguin Go

A Golang implementation of the classic [ECB Penguin](https://words.filippo.io/the-ecb-penguin/).

## Overview

This project demonstrates the Electronic Codebook (ECB) mode of AES encryption using an image of the Tux penguin. ECB is a simple mode of operation for a block cipher, where each block of plaintext is encrypted separately. This results in a one-to-one mapping where if you encrypt the same block using the same key, you will always get the same output. This results in patterns being left behind in the data, potentially allowing information to be uncovered.

Original Tux:

![Tux](tux.png)

Encrypted Tux using ECB mode:

![ECB Tux](ecb_tux.png)

## How It Works

1. The program reads an image file (e.g., `tux.png`).
2. It converts the image into RGB data.
3. The RGB data is encrypted using ECB mode.
4. The encrypted data is used to create a new image (e.g., `ecb_tux.png`).

## Usage

To run the program, use the following command:

```sh
go run encrypt.go <image-file-path> [seed]
```

The seed is optional. If omitted, a seed will be randomly chosen.
