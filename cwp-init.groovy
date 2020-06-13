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
        // backup current file, make sure it can be executed only one time
        File initFile = new File(instance.getRootDir(), "/war/WEB-INF/init.groovy.d/cwp-init.groovy");
        if(initFile.isFile()) {
            initFile.renameTo(new File(instance.getRootDir(), "/war/WEB-INF/init.groovy.d/cwp-init.groovy.bak"));
        }

        // remove bundled plugins
        File pluginsDir = new File(instance.getRootDir(), "/war/WEB-INF/plugins");
        if (pluginsDir.isDirectory()) {
            for(String plugin : pluginsDir.list()) {
                boolean result = new File(pluginsDir, plugin).delete();
                println "delete plugin " + plugin + " " + result
            }
        } else {
            println "plugins file is not a is directory"
        }
        println "Jenkins init ready..."
        break
    }
  }

  println "cert init thread done"
}