package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	//mongodbConnect()
	//insert()
	//update()
	find()
}
//定义表结构
type Student struct {
	Mid int64
	Name string
	Age int64
}
//查找
func find() {
	var(
		clientOptions *options.ClientOptions
		client *mongo.Client
		err error
		collection *mongo.Collection
		filter bson.D
		student Student
	)
	//连接到mongodb
	clientOptions = options.Client().ApplyURI("mongodb://10.235.25.241:27017")
	if client, err = mongo.Connect(context.TODO(), clientOptions); err != nil {
		fmt.Println(err)
		return
	}
	collection = client.Database("my_db").Collection("my_collection")

	filter = bson.D{{"name", "id1"}}
	student = Student{}
	err = collection.FindOne(context.TODO(), filter).Decode(&student);
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(student)

}
//更新文档
func update() {
	var(
		clientOptions *options.ClientOptions
		client *mongo.Client
		err error
		collection *mongo.Collection
		filter bson.D
		update bson.D
		updateRes *mongo.UpdateResult
	)
	//连接到mongodb
	clientOptions = options.Client().ApplyURI("mongodb://10.235.25.241:27017")
	if client, err = mongo.Connect(context.TODO(), clientOptions); err != nil {
		fmt.Println(err)
		return
	}
	collection = client.Database("my_db").Collection("my_collection")

	filter = bson.D{{"name", "id1"}}
	update = bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}
	if updateRes, err = collection.UpdateMany(context.TODO(), filter, update); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("更改的行数：", updateRes.ModifiedCount)
}

//插入文档
func insert() {
	var(
		clientOptions *options.ClientOptions
		client *mongo.Client
		err error
		db *mongo.Database
		collection *mongo.Collection
		student Student
		insertRes *mongo.InsertOneResult
		insertId string
	)
	//连接到mongodb
	clientOptions = options.Client().ApplyURI("mongodb://10.235.25.241:27017")
	if client, err = mongo.Connect(context.TODO(), clientOptions); err != nil {
		fmt.Println(err)
		return
	}

	db = client.Database("my_db")
	collection = db.Collection("my_collection")
	student = Student{
		Mid:   1,
		Name: "id1",
		Age:  20,
	}
	if insertRes, err = collection.InsertOne(context.TODO(), student); err != nil {
		fmt.Println(err)
		return
	}
	insertId = insertRes.InsertedID.(primitive.ObjectID).Hex()
	fmt.Println("插入数据的_id为：", insertId)
}

//连接
func mongodbConnect() {
	var(
		clientOptions *options.ClientOptions
		client *mongo.Client
		err error
	)
	//连接到mongodb
	clientOptions = options.Client().ApplyURI("mongodb://10.235.25.241:27017")
	if client, err = mongo.Connect(context.TODO(), clientOptions); err != nil {
		fmt.Println(err)
		return
	}
	//检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Connected to MongoDb")
}
