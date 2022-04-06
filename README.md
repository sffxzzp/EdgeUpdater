EdgeUpdater-Go
======
由于[shuax](https://www.shuax.com)做的Edge增强版经常需要手动更新，于是提取了一些工具写了一个 Golang 版自用。

使用方法
------
下载后执行`edgeupdater.bat`即可。

带有自动检查更新功能，如果无需自动检查，在`msedge.ini`文件里注释掉`0=edgeupdater.bat`一行即可。

设置
------
修改`edgeupdater`目录下的`settings.json`即可。

`Branch`是分支，有`Stable`,`Beta`,`Dev`,`Canary`四个分支版本。

`Structure`是版本位数，分别为`x86`,`x64`。（仅`x64`可用，因为`Edge++`是64位版本。

`Version`是当前版本号，一般不推荐修改。

依赖/使用
------
[Edge++](https://shuax.com/project/edge/)

[7-Zip](https://www.7-zip.org/)
