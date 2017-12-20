# Telephone
A application like the Cloudgo


路由规则：
http://localhost:8181/   访问根文件

http://localhost:8181/add  向数据库（这里是data.json文件添加数据，以post的方式）

http://localhost:8181/search  列出数据库(这里是data.json里面所有的数据)

http://localhost:8181/search/brand={brandName}  在数据库（这里是data.json文件）里面搜索以brandName为brand的记录

http://localhost:8181/search/lowestPrice={low}&highestPrice={high}  在数据库（这里是data.json文件）里面搜索在这个价格范围内的所有记录

http://localhost:8181/search/color={colorArray}  在数据库（这里是data.json文件）里面搜索包含colorArray的颜色的记录


使用了gorilla/mux框架。理由是因为这个框架体积比较小，而且API使用方便。

根目录下的data.txt文件是我测试时所进行的步骤顺序。

curl测试的结果放在screenShot文件夹里面，和data.txt里面的顺序一一对应。

ab压力测试我只进行了一个简单的测试，测试结果（以ab开头的图片）也是放在screenShot文件夹中。

ab -n 1000 -c 1000 http://localhost:8181
-n 请求次数
-c 并发数
