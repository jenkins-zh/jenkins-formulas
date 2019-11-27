import java.io.File;
import java.io.FileOutputStream;
import java.io.InputStream;
import java.net.URL;

URL url = new URL("https://raw.githubusercontent.com/jenkinsci/localization-zh-cn-plugin/master/src/main/resources/mirror-adapter.crt");
InputStream input = url.openStream();
FileOutputStream out = new FileOutputStream(System.getenv("JENKINS_HOME") +  "/war/WEB-INF/update-center-rootCAs/mirror-adapter.crt");
byte[] buf = new byte[1024];
int count = -1;

while((count = input.read(buf)) > 0) {
  out.write(buf, 0, count);
}
