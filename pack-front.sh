#!/usr/bin/env bash

echo "replacing the mode of frontend javascript"

sed -i 's/ModeDevelop/ModeDeploy/g' ./front/js/capture/config/config.js

cat ./front/js/capture/config/config.js

tar -cJf ./internal/utils/front.tar.xz ./front

echo "recovring the develop mode for frontend javascript"

sed -i 's/ModeDeploy/ModeDevelop/g' ./front/js/capture/config/config.js

cat ./front/js/capture/config/config.js
