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
    $PACKAGE_BASE/*.go \
    vendor/go_utils/utils/
