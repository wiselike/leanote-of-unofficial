# Leanote of non-official version
* for secondary development only
* update with official patches
* push to mainline as needed
* with experimental features
* with full ChangeLogs

  

# ChangeLogs
1. searched from https://github.com/leanote/leanote  
		with git(d58fd64)[gen tmp tool without revel] on 15 Aug 2021
2. patched https://github.com/ctaoist/leanote/commit/2cee584f793e21c7469e8701874d1548bee1be17
		which comes from https://github.com/leanote/leanote/compare/c4bb20fd129e63edd14bc7ecd229bbad3b13bcb7..450deb09bdf1ebc47ea31b0ed209b8d85492f7fa
		and https://github.com/leanote/leanote/pull/933/commits/92db56f4f141e477dbd1fa01232ea2c6536fe027	
3. patched https://github.com/ctaoist/leanote/commit/c5c19e32e0cb892fe35178a14dfe927049f5b3a9
4. patched https://github.com/ctaoist/leanote/commit/c2c4a5536301132a78594c2311d1dbd0d957b304
5. 自研的优化
6. patched "markdown编辑器增加字数统计功能" https://github.com/ctaoist/leanote/commit/297ca0c3ef15db680a7fe395b0283497dd768b2d and https://github.com/ctaoist/leanote/commit/7060829c7ab015431d05a529c4f2d31822992f15
7. 自研：修改配置文件，默认为中文的语言
8. 自研：添加自定义的git忽略文件
9. 自研：整理node图片，按标题来存放，以便于到服务器上检索维护
10. 自研：修复Site's URL设置后，不同步配置文件，导致重启生效的问题
11. 自研：添加在配置文件中自定义note的图片、附件存放路径
12. 自研：修改默认note历史数为5，并且添加app.conf配置文件可配。优化历史记录新增删除算法。修改note历史顺序，与官方原生不兼容，如使用，会自动删除之前的旧历史，无其他副影响
13. 自研：将所有配置参数，调整为从系统全局变量中读取，而不是每次都从文件中读。优化了读取速度和效率
14. 自研：使用gofmt格式化所有go代码，不对源码做任何手动改动
15. 自研：禁用github.io，改为使用本地css文件
16. 自研：禁用demo账号，自己用的话demo没有必要存在啊，直接用admin不就行啦
17. 自研：修复无法退出登录的故障
18. 自研：修正保存note历史记录的算法，调整note自动保存到历史记录的功能，用起来更顺畅
19. 上传原始package.json文件里定义的项目GPLv2 license
20. 自研：前端实现博客置顶设置
21. 优化note的字数统计功能
22. 自研：修复移动端界面的博客图标显示异常
23. 自研：改进验证码登录流程，降低爆破的可能性
24. 自研：添加图片备份文件夹，防止图片丢失
25. 自研：屏蔽首页的广告页，改为直接跳转为登录或者note页
26. 自研：清理数据库中冗余数据，将chirpy主题(非自研)合入为默认主题之一
27. 自研：修复发送邮件的中文标题乱码故障

