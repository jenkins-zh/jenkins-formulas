#!/bin/sh

formulas=$(yq read config.yaml formulas.*.name)

for formula in $formulas;do
    echo "check formula $formula"
    oldMd5=$(yq read config.yaml 'formulas(name=='$formula').md5')
    echo $oldMd5

    newMd5=$(md5 -q formulas/${formula}.yaml)

    if [[ "$newMd5" != "oldMd5" ]]; then
    fi
done
