# MGCrypt

Dead simple AES file encryption.

**Disclaimer: I'm no security expert, use at your own risk.**

[Download the latest release here.](https://github.com/mgerb/mgcrypt/releases)

---

## Usage

**Arguments**

* **-d** Decrypt (encrypts by default)
* **-k** Specify your encryption key
* **-o** Specify an output file

---

### Encrypt

```
mgcrypt -k yourSecretKey inputFile
```

### Decrypt
```
mgcrypt -k yourSecretKey -d inputFile
```

#### Output to a file (this works for both encryption/decription)

```
mgcrypt -k yourSecretKey -o outputFile inputFile

or

mgcrypt -k yourSecretKey inputFile > outputFile

```

## How does it work?

* The Go standard crypto libraries are used to encrypt/decrypt files with [AES](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard)
* A salt is stored with each encrypted file (the first 128 bits of SHA256 checksum of the original file)
* [pbkdf2](https://en.wikipedia.org/wiki/PBKDF2) is used to generate an AES key with password/salt
