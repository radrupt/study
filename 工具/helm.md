> Helm是Kubernetes的一个包管理工具，用来简化Kubernetes应用的部署和管理。可以把Helm比作CentOS的yum工具。 通过使用Helm可以管理Kubernetes manifest files、管理Helm安装包charts、基于chart的Kubernetes应用分发。

## Helm的基本概念
Chart: 代表着 Helm 包。它包含在 Kubernetes 集群内部运行应用程序，工具或服务所需的所有资源定义。你可以把它看作是 Homebrew formula，Apt dpkg，或 Yum RPM 在Kubernetes 中的等价物。

Repository:（仓库） 是用来存放和共享 charts 的地方。它就像 Perl 的 CPAN 档案库网络 或是 Fedora 的 软件包仓库，只不过它是供 Kubernetes 包所使用的。

Release: 是运行在 Kubernetes 集群中的 chart 的实例。一个 chart 通常可以在同一个集群中安装多次。每一次安装都会创建一个新的 release。以 MySQL chart为例，如果你想在你的集群中运行两个数据库，你可以安装该chart两次。每一个数据库都会拥有它自己的 release 和 release name。

Chart目录结构：
1、chart.yaml
Yaml文件，用来描述chart的摘要信息

2、readme.md
Markdown格式的readme文件，此文件为可选

3、LICENSE
文本文件，描述chart的许可信息，此文件为可选

4、requirements.yaml
Yaml文件，用来描述chart的的依赖关系，在安装过程中，依赖的chart也会被一起安装

5、value.yaml
Yaml文件，chart支持在安装的时候做对配置参数做定制化配置，value.yaml文件为配置参数的默认值

6、templates目录
各类k8s资源的配置模板目录

## 如何使用？
需求有哪些？
添加仓库：helm repo add

参考：
https://helm.sh/zh/docs/intro/using_helm/