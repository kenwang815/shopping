# Shopping

### 環境需求
- [Golang](https://golang.org/doc/install) 
- [Docker](https://docs.docker.com/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### 使用套件
| 名稱                  | 描述 | 
| ------------------------------ | ---------|
| github.com/gin-gonic/gin       |  HTTP web framework  |
| github.com/gofrs/uuid          |  Support both the creation and parsing of UUIDs in different formats |
| github.com/jinzhu/gorm | The fantastic ORM library for Golang |
| github.com/mattn/go-sqlite3 | Sqlite3 driver conforming to the built-in database/sql interface |
| github.com/sirupsen/logrus | Structured, pluggable logging for Go |
| github.com/stretchr/testify | A toolkit with common assertions and mocks |

### 運行方式
- 建置後端服務
```
docker-compose build
```

- 啟用資料庫與後端服務
```
docker-compose up -d
```

- 關閉資料庫與後端服務
```
docker-compose down
```

### API 操作
**商品目錄**
- 新增商品目錄
```
curl --location --request POST '0.0.0.0:8080/v1/catalog' \
--form 'name="3C"' \
--form 'hide="false"'
```

- 取得商品目錄列表
```
curl --location --request GET '0.0.0.0:8080/v1/catalog?page=1&number=20'
```

- 修改商品目錄
```
curl --location --request PUT '0.0.0.0:8080/v1/catalog' \
--form 'id="1"' \
--form 'hide="true"'
```

- 刪除商品目錄
```
curl --location --request DELETE '0.0.0.0:8080/v1/catalog/1'
```

**商品**
- 新增商品
```
curl --location --request POST '0.0.0.0:8080/v1/commodity' \
--form 'catalog_id="1"' \
--form 'name="TV"' \
--form 'cost="5000"' \
--form 'price="7000"' \
--form 'description="43吋 LED顯示器 HF-43VA5"' \
--form 'sell="false"' \
--form 'start_time="2021-05-13 20:04:05"' \
--form 'end_time="2021-05-20 20:04:05"'
```

- 取得商品列表
```
curl --location --request GET '0.0.0.0:8080/v1/commodity?page=1&number=20'
```

- 修改商品
```
curl --location --request PUT '0.0.0.0:8080/v1/commodity' \
--form 'id="1"' \
--form 'end_time="2021-05-27 20:04:05"'
```

- 刪除商品
```
curl --location --request DELETE '0.0.0.0:8080/v1/commodity/1'
```

**進階功能**
- 取得商品目錄下的商品
```
curl --location --request GET '0.0.0.0:8080/v1/catalog/1?page=1&number=20'
```