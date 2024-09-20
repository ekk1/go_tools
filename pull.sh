rm -rf vendor
mkdir -p vendor/go_utils/utils
PACKAGE_BASE="../go_utils/utils"

cp -r \
    $PACKAGE_BASE/bencode \
    $PACKAGE_BASE/torrent \
    $PACKAGE_BASE/screen \
    $PACKAGE_BASE/minikv \
    $PACKAGE_BASE/myhttp \
    $PACKAGE_BASE/webui \
    $PACKAGE_BASE/nbd \
    $PACKAGE_BASE/openai \
    $PACKAGE_BASE/redis \
    $PACKAGE_BASE/*.go \
    vendor/go_utils/utils/

mkdir -p vendor/official/md4
cp -r ../crypto/md4/md4* vendor/official/md4
