# goMerge

HTML5在游戏方面有不少引擎,但对于新事物来说,引擎本身就在试错中,用到工作中不仅要付出学习成本,还要承担不知何时何处会爆发的坑,于是坚持自己造轮子(其实想通过此过程多学点东西).
时间长了,每个项目里会包含大致相同,而小地方有更新的代码,就把功能代码分割成一个个小文件,引入,通过文件名版本号区别,现在需要把小功能文件包含成一个文件,于是就去找包含合并功能的软件,发现大多数发布程序功能都很多,但需要建立发布规则什么的,我只是想把几个文件自动合并而已,没必要搞那么复杂.
年底前数个基础架构相同的项目要赶活,
于是...
用golang写了个小程序,识别代码里需要合并的文件

project/index.html包含以下内容

```html
<script src="../codeBase/scr/__-1.0.js" mergeTo="js/game.js"></script>
<script src="../codeBase/scr/Actor-2.0.js" mergeTo="js/game.js"></script>
<script src="../codeBase/scr/Camera-1.0.js" mergeTo="js/game.js"></script>
<script src="../codeBase/scr/Ui-1.01.js" mergeTo="js/game.js"></script>
<script src="js/self.js" mergeTo="js/game.js"></script>
<link type="text/css" rel="stylesheet" href="../codeBase/scr/css/__-1.0.css" mergeTo="res/css.css"/>
    <link type="text/css" rel="stylesheet" href="../codeBase/scr/css/textLine-1.0.css" mergeTo="res/css.css"/>
<link type="text/css" rel="stylesheet" href="res/self.css" mergeTo="res/css.css"/>
```

本地开发使用分离的小文件,项目在线使用合并后的文件.只要把index.html文件的本地文件地址录入程序即可.
然后加入了监控文件改动功能,这样任意一个引入的文件修改了,都会设置引入这个文件的所有项目为需要更新状态,选择这些项目更新,会重新合并生成目标文件,此时会在文件的顶部用注释写入更新时间,为了识别缓存.
工作中只需要修改代码,然后列出需要更新的项目,可能会有好几个项目,选择更新即可.不用过多操心那些项目与引入功能代码的问题.
可随时对功能代码复制一份后加版本号升级,不合用返回也只是改个版本号,重新生成的事情.
配置信息和数据保存为json.启动程序自动加载恢复.
做了一个简单的控制页面,使用websocket通讯,避免需要控制台录入指令.

起始目的其实很简单,就是为了工作中少点麻烦事儿,可以灵活写代码而已.
