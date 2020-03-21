该目录下，保存的是自定义 Jenkins 发行版的配置文件（模板）。

配方 YAML 文件的格式，遵循 [Custom Jenkins WAR packager for Jenkins](https://github.com/jenkinsci/custom-war-packager/) 中的规范，
这里在原有的基础上，增强了部分功能。例如：使得 YAML 还可以作为模板使用。

在您添加配方时，请考虑下面的几点：
* 已有的配方是否可以满足您的需求
* 配方要尽可能满足较多用户的需求，避免只针对特殊的场景定制配方

## 步骤
1. 拷贝一份已有的配方
2. 根据需要增加或者删除插件配置
3. 把新配方加到 [config.yaml](../config.yaml) 中的 `formulas.name`，`md5` 不需要填写
4. 执行命令构建包：`jcli cwp --config-path formulas/your.yaml --tmp-dir tmp --value-set version=2.223`
5. 启动构建好的 `jenkins.war` 后，检查是否有异常
6. 提交 Pull Request

## 变量
为了让同一份配方可以构建出多个不同版本的包，我们需要在设置 Jenkins 版本的地方，使用变量 `{{.version}}` 来替代。
