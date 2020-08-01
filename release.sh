#!/bin/sh
make release
version=$(grep 'VERSION_TAG=' Makefile | cut -d'=' -f2)
cd release-"${version}" || exit
for binary in "$(pwd)"/*; do
  name=$(basename "${binary}")
  shasum -a 256 "${name}" > "${name}".'sha256sum'
  shasum -a 256 -c "${name}".'sha256sum'
  didctl sign dev --key 'code-sign' > "${name}".'sig.json' < "${name}".'sha256sum'
  didctl verify "${name}".'sig.json' < "${name}".'sha256sum'
done
