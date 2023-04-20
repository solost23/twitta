package models

const (
	LogPrivateLatterTypeAcceptOrReject = iota
	LogPrivateLatterTypeAccept
	LogPrivateLatterTypeReject
	LogPrivateLatterTypePrivateLatter
)

type LogPrivateLatter struct {
	BaseModel
	ID       string `json:"id" bson:"_id"`
	UserId   string `json:"userId" bson:"user_id"`
	TargetId string `json:"targetId" bson:"target_id"`
	Content  string `json:"content" bson:"content"`
	Type     uint   `json:"type" bson:"type" comment:"是否通过好友 0 接受/拒绝 1 接受 2 拒绝 3 朋友私信内容"`
}

func NewLogPrivateLatter() *LogPrivateLatter {
	return &LogPrivateLatter{}
}

func (*LogPrivateLatter) TableName() string {
	return "log_private_latter"
}

//func (t *LogPrivateLatter) InsertOne(ctx context.Context, db *mongo.Client, dbName string, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
//	result, err := db.Database(dbName).Collection(t.TableName()).InsertOne(ctx, document, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return result, nil
//}
//
//func (t *LogPrivateLatter) InsertMany(ctx context.Context, db *mongo.Client, dbName string, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
//	result, err := db.Database(dbName).Collection(t.TableName()).InsertMany(ctx, documents, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return result, nil
//}
//
//func (t *LogPrivateLatter) FindOne(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, result interface{}, opts ...*options.FindOneOptions) error {
//	err := db.Database(dbName).Collection(t.TableName()).FindOne(ctx, filter, opts...).Decode(result)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (t *LogPrivateLatter) Find(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, results interface{}, opts ...*options.FindOptions) error {
//	cur, err := db.Database(dbName).Collection(t.TableName()).Find(ctx, filter, opts...)
//	if err != nil {
//		return err
//	}
//	if err = cur.All(ctx, results); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (t *LogPrivateLatter) Update(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
//	updateResult, err := db.Database(dbName).Collection(t.TableName()).UpdateMany(ctx, filter, update, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return updateResult, nil
//}
//
//func (t *LogPrivateLatter) Delete(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
//	deleteResult, err := db.Database(dbName).Collection(t.TableName()).DeleteMany(ctx, filter, opts...)
//	if err != nil {
//		return nil, err
//	}
//	return deleteResult, nil
//}
//
//func (t *LogPrivateLatter) Count(ctx context.Context, db *mongo.Client, dbName string, filter interface{}, opts ...*options.CountOptions) (int64, error) {
//	count, err := db.Database(dbName).Collection(t.TableName()).CountDocuments(ctx, filter, opts...)
//	if err != nil {
//		return 0, err
//	}
//	return count, nil
//}
