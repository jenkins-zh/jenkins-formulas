import java.io.File;
import java.io.FileOutputStream;
import java.io.InputStream;
import java.net.URL;
import hudson.init.InitMilestone;
import jenkins.model.Jenkins;

Thread.start {
  while(true) {
    Jenkins instance = Jenkins.getInstance();
    InitMilestone initLevel = instance.getInitLevel();
    Thread.sleep(1500L);
    println "Jenkins not ready when handle init config..."
    if (initLevel >= InitMilestone.PLUGINS_STARTED) {
      InputStream input = new FileInputStream("/usr/share/jenkins/ref/mirror-adapter.crt")
      FileOutputStream out = new FileOutputStream(System.getenv("JENKINS_HOME") +  "/war/WEB-INF/update-center-rootCAs/mirror-adapter.crt");
      byte[] buf = new byte[1024];
      int count = -1;

      while((count = input.read(buf)) > 0) {
        out.write(buf, 0, count);
      }
      println "Jenkins init ready..."
      break
    }
  }
}