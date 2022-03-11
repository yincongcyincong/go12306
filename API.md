```bigquery
订单接口
curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "train_data":{
            "SecretStr":"vhrJYQ2HB9qi0fCAV1y6v8x01736PGt0xzLGwIlIaIAcfzNnaFT%2BaZ22LKUcgDYZIukNK59X0E5Q%0Anx5POI5maxGNQHzVviVUf%2BZxtOdZbdlo18eP7Xz0iN%2FBcRlx2g7iCCoUFgwplFYIBP7LC0y01aH1%0AbksEKVszDkmoegO5Y6b88pJHKX6vBCCAD4L74fR0UQufG3FTevjI0nuoYTCCeaXKU%2B8k0kUOhnRy%0AkoHv0hadUdztnsj3JV1xHWHhN0hi28OWXi4nKJ0Tmd2tI%2BEdJq6LYOhO%2BhWqvF6FcUzS7oMnleS1%0A4a5Nwxb6naE%3D",
            "TrainNo":"G45",
            "FromStationName":"北京南",
            "FromStation":"VNP",
            "ToStationName":"天津南",
            "ToStation":"TIP",
            "TrainLocation":"",
            "StationTrainCode":"",
            "LeftTicket":"",
            "StartTime":"06:54",
            "ArrivalTime":"07:25",
            "DistanceTime":"00:31",
            "Status":"预订",
            "SeatInfo":{
                "一等座":"无",
                "二等座":"有",
                "动卧":"",
                "商务座":"无",
                "无座":"",
                "特等座":"",
                "硬卧":"",
                "硬座":"",
                "软卧":"",
                "软座":""
            }
        },
  "search_param":{
    "train_date":"2022-02-17",
    "from_station":"BJP",
    "to_station":"TJP",
    "from_station_name":"北京",
    "to_station_name":"天津",
    "seat_type":"O"
  },
  "passenger_name":{
    "张三":true
  }
}' \
 'http://127.0.0.1:28178/order'
```

```bigquery
生成二维码接口
curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'' \
 'http://127.0.0.1:28178/create-image'
```

```bigquery
登陆接口
curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "image":"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",
  "result_code":"0",
  "result_message":"生成二维码成功",
  "uuid":"aeXAxBwl28bJPHzlezsXqYQ0f9OuIMD-xMd7PnhW-IIQBJPuO9OC0v4tDwtQq4QVEkKvyRpXj1HA511"
}' \
 'http://127.0.0.1:28178/login'
```
```bigquery
获取车辆信息接口
curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "train_date":   "2022-03-17",
  "from_station": "BJP",
  "to_station":   "TJP",
  "from_station_name": "北京",
  "to_station_name":   "天津",
  "seat_type": "O"
}' \
 'http://127.0.0.1:28178/search-train'
```

```bigquery
获取用户相关信息接口
curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'' \
 'http://127.0.0.1:28178/search-info'
```

```bigquery
退出登陆接口
curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'' \
 'http://127.0.0.1:28178/logout'
```

```bigquery
候补订单接口
curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "train_data":{
    "SecretStr":"oEN4xUyJL3Tdq8OJNdbLFhX8CxB%2BVE8DpXkIyuCZRYuIskPJA5nK81fUYorkgq1pSnlhZN3FsN8%2B%0AFYIWMNi%2B%2FM3eIAMHWNBrHsCjUimsPPX%2FJKblsouQmssQCJZBfXAcmPHsh2NnZ1cqfFIfx4tVi05f%0AwH7D95nxkryvnsUTUKSDXtOXzGc%2FDS1tVeWVMI6Wo4qIjSUy9ZJtPh2Rw35TAHypEkM%2BxKvwhdZ7%0AOQRGSjXIZ6MiQHJR%2Bn70W29LpuqS%2FQk4nnDo%2F2ZJ8P83VaAbCqtfVQ4DI%2F9QkqvT29aQh5lw06bf%0ANDjuCYnWQpE%3D",
    "TrainNo":"G45",
    "FromStationName":"北京南",
    "FromStation":"VNP",
    "ToStationName":"天津南",
    "ToStation":"TIP",
    "TrainLocation":"",
    "StationTrainCode":"",
    "LeftTicket":"",
    "StartTime":"06:54",
    "ArrivalTime":"07:25",
    "DistanceTime":"00:31",
    "Status":"预订",
    "SeatInfo":{
      "一等座":"无",
      "二等座":"有",
      "动卧":"",
      "商务座":"无",
      "无座":"",
      "特等座":"",
      "硬卧":"",
      "硬座":"",
      "软卧":"",
      "软座":""
    }
  },
  "search_param":{
    "train_date":"2022-02-17",
    "from_station":"BJP",
    "to_station":"TJP",
    "from_station_name":"北京",
    "to_station_name":"天津",
    "seat_type":"O"
  },
  "passenger_name":{
    "张三":true
  }
}' \
 'http://127.0.0.1:28178/hb'
```

### 使用方式
```bigquery
生成二维码——>拿到二维码接口的出参
用二维码接口的出参数放到登陆接口，这时接口会阻塞，需要扫描二维码
然后就登陆成功
后面获取用户信息
然后遍历列车信息即可
最后通过列车信息，列车搜索信息和选中的乘车人姓名进行
请求买票或者候补接口即可

```
