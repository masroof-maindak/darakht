#!/usr/bin/env bash

#
# Utils
#
file_exists() {
	if [ ! -f "$1" ]; then
		echo "File not found: $1"
		return 1
	fi
}

md2exe="md2sum"
md2src="${md2exe}.c"

if [ ! -f "$md2src" ]; then
	cat <<EOF >"$md2src"
#include <md2.h>
#include <stdio.h>
int main(int argc, char **argv) {
	if (argc != 2)
		return fprintf(stderr, "Usage: %s <file>\n", argv[0]), 1;
	char buf[MD2_DIGEST_STRING_LENGTH];
	char *tmp = MD2File(argv[1], buf);
	if (tmp == NULL)
		return fprintf(stderr, "Usage: %s <file>\n", argv[0]), 1;
	puts(buf);
	return 0;
}
EOF
fi
if [ ! -f "$md2exe" ]; then gcc "$md2src" -o "$md2exe" -lmd; fi
if [ ! -x "$md2exe" ]; then chmod +x "$md2exe"; fi

#
# Q1: Hashing
#
hf() {
	if [ -z "$1" ]; then echo "Usage: hf <file>"; return 1; fi

	file_exists "$1" || return 1
	
	echo -n "md5: " && rhash --md5 "$1" | awk '{print $1}'
	echo -n "md4: " && rhash --md4 "$1" | awk '{print $1}'
	echo -n "md2: " && ./md2sum "$1"
	echo -n "sha256: " && rhash --sha256 "$1" | awk '{print $1}'
	echo -n "sha3-224: " && rhash --sha3-224 "$1" | awk '{print $1}'
}

#
# Q2: Encrypt and Decrypt
#
aes_encrypt() {
	if [ -z "$1" ]; then echo "Usage: aes_encrypt <file> [<password>]"; return 1; fi

	local file=$1
	local password=$2

	file_exists "$file" || return 1

	openssl enc -aes128 -in "$file" -out "${file}.aes128" -pass pass:"$password"
}

aes_decrypt() {
	if [ -z "$1" ]; then echo "Usage: aes_decrypt <file>.aes128 [<password>]"; return 1; fi

	local file=$1
	local password=$2

	file_exists "$file" || return 1

	openssl enc -aes128 -d -in "$file" -out "${file%.aes128}" -pass pass:"$password"
}

#
# Q3: Digital Signature
#
generate_keys() {
	# Usage: generate_keys [<private_key_file>] [<public_key_file>]
	local privKey=${1:-"key.pem"}
	local pubKey=${2:-"key.pub"}

	openssl genrsa -out "$privKey" 4096
	openssl rsa -in "$privKey" -pubout >"$pubKey"
}

sign_file() {
	if [ -z "$1" ]; then echo "Usage: sign_file <file> [<private_key_file>] [<signature_file_name>]"; return 1; fi

	local file=$1
	local privKey=${2:-"key.pem"}
	local sigFile=${3:-"${file}.sig"}

	file_exists "$file" || return 1
	file_exists "$privKey" || return 1

	openssl dgst -sha256 -sign "$privKey" -out "$sigFile" "$file"
}

verify_signature() {
	if [ -z "$1" ]; then echo "Usage: verify_signature <file.sig> [<public_key_file>] [<original_file>]"; return 1; fi

	local sigFile=$1
	local pubKey=${2:-"key.pub"}
	local file=${3:-"${sigFile%.sig}"}

	file_exists "$sigFile" || return 1
	file_exists "$file" || return 1
	file_exists "$pubKey" || return 1

	openssl dgst -sha256 -verify "$pubKey" -signature "$sigFile" "$file"
}
