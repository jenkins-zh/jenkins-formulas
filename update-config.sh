#!/bin/sh

ltsVersion=$(curl http://updates.jenkins.io/stable/latestCore.txt -L -s)
weeklyVersion=$(curl http://updates.jenkins.io/latestCore.txt -L -s)

newLts=$(yq read config.yaml lts | grep $ltsVersion)
if [[ "$newLts" == "" ]]; then
    echo "found a new version of lts $ltsVersion"

    yq write config.yaml 'lts.+' $ltsVersion -i
else
    echo "no new version of lts"
fi

newWeekly=$(yq read config.yaml weekly | grep $weeklyVersion)
if [[ "$newWeekly" == "" ]]; then
    echo "found a new version of lts $weeklyVersion"

    yq write config.yaml 'weekly.+' $weeklyVersion -i
else
    echo "no new version of weekly"
fi
