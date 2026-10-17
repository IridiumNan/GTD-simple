#!/usr/bin/env bash

echo "replacing the mode of frontend javascript"

sed -i 's/dep/dev/g' ./front/js/capture/config/config.js

cat ./front/js/capture/config/config.js

rm ./internal/utils/front.tar.xz &>/dev/null
tar -cJf ./internal/utils/front.tar.xz ./front
echo "recovring the develop mode for frontend javascript"
