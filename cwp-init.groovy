import hudson.Plugin;
import java.io.File;
import java.io.FileOutputStream;
import java.io.InputStream;
import hudson.init.InitMilestone;
import jenkins.model.Jenkins;

Thread.start {
  while(true) {
    Jenkins instance = Jenkins.getInstance();
    InitMilestone initLevel = instance.getInitLevel();
    Thread.sleep(1500L);
    println "Jenkins not ready when handle init config..."
    if (initLevel >= InitMilestone.PLUGINS_STARTED) {
      Plugin zhPlugin = instance.getPlugin("localization-zh-cn");
      if (zhPlugin != null) {
        InputStream input = zhPlugin.getWrapper().classLoader.getResourceAsStream("mirror-adapter.crt");
        if (input == null) {
          System.err.println("cannot found mirror-adapter.crt from localization-zh-cn, would not copy cert file");
          break;
        }

        File certPath = new File(instance.getRootDir(), "/war/WEB-INF/update-center-rootCAs/mirror-adapter.crt");
        FileOutputStream out = new FileOutputStream(certPath);

        byte[] buf = new byte[1024];
        int count = -1;

        while((count = input.read(buf)) > 0) {
          out.write(buf, 0, count);
        }

        println "Jenkins init ready..."
      } else {
        System.err.println("cannot found localization-zh-cn, would not copy cert file");
      }
      break
    }
  }
}