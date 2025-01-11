jwt zai go-zero 中有点不熟悉, 来不及调试了
task 监控来不及写了

test.csv 放在  etc 目录下

curl --location 'http://0.0.0.0:8888/tasks' \
--form 'file=@"/Users/zhengwei/Downloads/test.csv"'
curl --location 'http://0.0.0.0:8888/tasks/3'
curl --location --request POST 'http://0.0.0.0:8888/tasks/3/translate'
curl --location 'http://0.0.0.0:8888/tasks/3/download'


curl --location 'http://0.0.0.0:8888/auth/users' \
--header 'Content-Type: application/json' \
--data '{
    "user_name": "likun",
    "user_pwd": "123"
}'

curl --location 'http://0.0.0.0:8888/auth/login' \
--header 'Content-Type: application/json' \
--data '{
    "user_name": "liun",
    "user_pwd": "12"
}'
