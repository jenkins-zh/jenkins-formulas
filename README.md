[![Docker Pulls](https://img.shields.io/docker/pulls/jenkinszh/jenkins-zh.svg)](https://hub.docker.com/r/jenkinszh/jenkins-zh/tags)

# docker-zh

该 Jenkins 镜像，自带了部署在中国的代理更新中心。

`docker run --rm -p 8080:8080 jenkinszh/jenkins-zh:latest`

下面的例子可以把 Jenkins 的数据目录挂载到本地：

`docker run -u root -v /var/jenkins/data:/var/jenkins_home -p 18080:8080 jenkinszh/jenkins-zh:latest`
