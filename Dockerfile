FROM jenkins/jenkins:2.190.1

ENV JENKINS_UC https://updates.jenkins-zh.cn

RUN mkdir -p $JENKINS_HOME/war/WEB-INF/update-center-rootCAs/ \
    && curl https://github.com/jenkinsci/localization-zh-cn-plugin/blob/master/src/main/resources/mirror-adapter.crt \
    -o $JENKINS_HOME/war/WEB-INF/update-center-rootCAs/mirror-adapter.crt

COPY active.txt active.txt
RUN plugins.sh active.txt
