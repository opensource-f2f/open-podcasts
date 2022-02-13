> 以下想法源于[《开源面对面》](https://github.com/opensource-f2f/episode) 。既然播客节目的制作过程、内容是可以开源的，那么，
> 播客平台为什么不可以呢？

本项目，计划采用云原生的方式来开发后端程序，利用 [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) 作为
脚手架。

```shell
k3d cluster create -p 30000:30000 -p 30001:30001 -p 30002:30002 -p 30003:30003
```

## Create new API

```shell
kubebuilder create api --group osf2f --version v1alpha1 --kind Profile
```