## Jenkins 中国定制版
目前定制版发行包括有：Docker 镜像、jenkins.war 文件。所有的 Jenkins 定制版本都包括如下的特性：

* 配置有部署在中国的[代理更新中心](https://github.com/jenkins-zh/mirror-proxy)（update center）
* [简体中文插件](https://github.com/jenkinsci/localization-zh-cn-plugin)

## 镜像
[![Docker Stars](https://img.shields.io/docker/stars/jenkinszh/jenkins-zh.svg)](https://hub.docker.com/r/jenkinszh/jenkins-zh/)
[![Docker Pulls](https://img.shields.io/docker/pulls/jenkinszh/jenkins-zh.svg)](https://hub.docker.com/r/jenkinszh/jenkins-zh/tags)

使用命令如下：

`docker run --rm -p 8080:8080 jenkinszh/jenkins-zh:lts`

下面的例子可以把 Jenkins 的数据目录挂载到本地：

`docker run -u root -v /var/jenkins/data:/var/jenkins_home -p 8080:8080 jenkinszh/jenkins-zh:latest`

[点击这里](https://github.com/jenkins-zh/docker-zh/packages/134536/versions)查看所有 `docker tag` 的版本。

## war
[![下载](https://api.bintray.com/packages/jenkins-zh/jenkins-cli/jenkins/images/download.svg) ](https://bintray.com/jenkins-zh/jenkins-cli/jenkins/_latestVersion)

这种发行版除了包含上述的公共特性外，还包括：

* [配置即代码插件](https://github.com/jenkinsci/configuration-as-code-plugin/)

[点击这里](https://dl.bintray.com/jenkins-zh/jenkins-cli/jenkins/)查看所有 `jenkins.war` 的版本。

## 参考
[Jenkins 官方 Docker Hub 地址](https://hub.docker.com/r/jenkins/jenkins/tags)

## 反馈

该项目还处于早起阶段，我们欢迎任何人以任何形式帮助完善、提出改进建议。
