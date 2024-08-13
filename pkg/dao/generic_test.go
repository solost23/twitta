package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
	"twitta/pkg/mongoos"
)

// go test -parallel=1 串行执行单元测试

func initDB() (db *mongo.Database) {
	conn, err := mongoos.NewMongoConnect(context.Background(), &mongoos.Config{
		Hosts:      []string{"localhost:27017"},
		AuthSource: "",
		Username:   "",
		Password:   "",
		Timeout:    30,
	})
	if err != nil {
		panic(err)
	}
	return conn.Database("twitta")
}

//func TestGInsertOne(t *testing.T) {
//	db := initDB()
//
//	user := User{
//		CreatedAt: time.Now(),
//		UpdatedAt: time.Now(),
//		Username:  "alex",
//	}
//	_, err := GInsertOne[*User](context.TODO(), db, &user)
//	if err != nil {
//		panic(err)
//	}
//}

func TestGPaginatorOrder(t *testing.T) {
	db := initDB()

	// 分页查询数据
	users, total, pages, err := GPaginatorOrder[*User](context.TODO(), db, &ListPageInput{
		Page: 1,
		Size: 10,
	}, bson.M{"_id": 1}, bson.M{"username": "alex"})
	if err != nil {
		panic(err)
	}
	fmt.Println(users)
	fmt.Println(total, pages)
}

func TestGInsertMany(t *testing.T) {
	db := initDB()

	users := []*User{
		{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Username:  "alex1",
		},
		{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Username:  "alex2",
		},
	}
	err := GInsertMany[*User](context.TODO(), db, users)
	if err != nil {
		panic(err)
	}
}

func TestGWhereFirst(t *testing.T) {
	db := initDB()

	result, err := GWhereFirst[*User](context.TODO(), db, bson.M{"username": "alex1"})
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func TestGWhereFind(t *testing.T) {
	db := initDB()

	result, err := GWhereFind[*User](context.TODO(), db, bson.M{"username": "alex1"})
	if err != nil {
		panic(err)
	}
	for i := 0; i != len(result); i++ {
		fmt.Println(result[i])
	}
}

func TestGWhereUpdate(t *testing.T) {
	db := initDB()

	result, err := GWhereUpdate[*User](context.TODO(), db, bson.M{"username": "alex1"}, bson.M{"$set": bson.M{"username": "egon"}})
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func TestGWhereDelete(t *testing.T) {
	db := initDB()

	result, err := GWhereDelete[*User](context.TODO(), db, bson.M{"username": "egon"})
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func TestGWhereCount(t *testing.T) {
	db := initDB()

	result, err := GWhereCount[*User](context.TODO(), db, bson.M{"username": "alex2"})
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
