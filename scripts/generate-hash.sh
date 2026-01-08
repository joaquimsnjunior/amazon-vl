#!/bin/bash
# Generate MD5 crypt hash for password
# Usage: ./generate-hash.sh <password>

set -e

if [ -z "$1" ]; then
    echo "Usage: $0 <password>"
    echo "Example: $0 mypassword"
    exit 1
fi

PASSWORD="$1"
SALT=$(openssl rand -base64 6 | tr -d '/+=' | head -c 8)

# Generate MD5 crypt hash
HASH=$(openssl passwd -1 -salt "$SALT" "$PASSWORD")

echo ""
echo "Password: $PASSWORD"
echo "Hash: $HASH"
echo ""
echo "To use this hash, set the environment variable:"
echo "  export AUTH_HASH='$HASH'"
echo ""
