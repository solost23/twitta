echo -e "****Twitta 开始自动部署****"

echo -e "---step1: 合并代码---"
git pull origin master
echo -e "合并代码成功"

echo -e "---step2: 编译代码---"
go mod tidy && go build -o ./build/start cmd/twitta/main.go
echo -e "编译代码成功"

echo -e "---step3: 更改权限---"
chmod -R 777 ./build/start
echo -e "更改权限完成"

echo -e "---step4: 启动项目---"
i1=`ps -ef|grep -E "./build/start"|grep -v grep|awk '{print $2}'`
if (( $i1 > 0 ))
then
  echo -e "杀掉进程$i1"
  kill 9 $i1
fi
nohup ./build/start >/dev/null 2>&1 &
i2=`ps -ef|grep -E "./build/start"|grep -v grep|awk '{print $2}'`
echo -e "****部署成功,部署的进程ID为:$i2****"