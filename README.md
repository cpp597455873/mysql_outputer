# 一个方便的sql导出器  

为什么有这个工具？因为我们经常需要导出数据，但是我们的普通的机器可能没有权限导出数据。我们就需要登录到目标机器上去导出数据，但是我们登录的数据库用户又没有导出权限,这个时候这个工具的就体现出来了。原本写的python脚本，但是python脚本有他的劣势就是机器上要安装python以及相关的依赖，这个时候我们很有可能被限制不能安装这些依赖或者我们不想动到机器上的依赖，这个工具就是一个很适合的工具。

### 如何使用
将data_outputer拷贝到目标机器上去  
1、执行chmod +x data_outputer, 赋予程序执行权限  
2、在data_outputer同级目录填写一个qry.sql的文件，并填入相关的查询sql  
3、然后就可以使用具体的命令来导出了，具体下文的使用示例  

### 参数解释
* **-s**  指定你需要执行的sql文件，默认为qry.sql
* **-f**  指定导出的格式 支持csv,json,sql。默认为csv
* **-t**  csv文件分隔符 默认逗号  可选值comma,tab
* **-o**  指需要指定导出为特定的文件名的时候请传入
* **-n**  当指定为sql导出的时候的表名，默认为xxxtable
* **-r**  当你想更换数据库连接信息的时候可以传入，当值为1的时候就重置数据库连接信息，当然你可以直接删除outputer.conf文件就好了  

### 如何指定数据库连接？
直接运行程序就会提示，输入user:password@tcp(ip:port)/dbname?charset=utf8类似的链接就可以连上数据库了。后续导出器将会加密保存你的连接信息到outputer.conf里面。

## 使用示例
### 1.导出为csv
```
./data_outputer 

./data_outputer -t=tab -s=xxx.sql -r=1
说明:上述命令导出了一个csv,使用的sql文件时xxx.sql,使用的csv分隔符为TAB，并且此时使用了新的数据库连接
```

### 2.导出为json
```
./data_outputer -f=json
```

### 3.导出为sql
```
./data_outputer -f=sql -n=your_table_name
```

### 4.批量sql导出
```
qry.sql

#file=student.csv
#format=csv
select * from student;

#file=teacher.json
#format=json
select * from teacher;

#file=school.sql
#format=sql
#table=school
select * from school;



./data_outputer 
说明:上述命令导出了两个SQL，并使用student.csv，teacher.json 作为文件名，写法#file=  #format=  #table= 作为其文件名 导出格式
注意: 1、sql之前请用分号隔开 2、 #file=  #format=  #table= 需要单独的一行
```




