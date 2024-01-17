#!/bin/bash
set -e

sudo apt-get update && \
sudo apt-get install -y gnupg2 pass

mkdir -p ~/.gnupg && \
chown -R "$(whoami)" ~/.gnupg && \
find ~/.gnupg -type f -exec chmod 600 {} \; && \
find ~/.gnupg -type d -exec chmod 700 {} \;

gpg --command-fd 0 --no-tty --batch --gen-key <<- EOF
	%echo Generating a basic OpenPGP key
	Key-Type: RSA
	Key-Length: 2048
	Subkey-Type: RSA
	Subkey-Length: 2048
	Name-Real: TestUser
	Name-Email: user@example.com
	Expire-Date: 0
	%no-ask-passphrase
	%no-protection
	%pubring pubring.gpg
	%secring secring.gpg
	%commit
	%echo done
EOF

pass init "$(gpg --list-keys --with-colons | head -n 5 | awk -F: '/^pub:/ { print $5 }')"
