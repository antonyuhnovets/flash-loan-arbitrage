package postgres

import (
	c "context"
	"reflect"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Storage struct {
	db *gorm.DB
}

func Connect(dsn string, i ...interface{}) (*Storage, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default,
		// DisableForeignKeyConstraintWhenMigrating: true,
		// FullSaveAssociations:                     true,
		// IgnoreRelationshipsWhenMigrating:         true,
		// AllowGlobalUpdate:                        true,
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(i...); err != nil {
		return nil, err
	}

	return &Storage{db}, nil
}

func (ps *Storage) Store(ctx c.Context, where string, item interface{}) (
	err error,
) {

	err = ps.db.Table(where).WithContext(ctx).Create(item).Error
	if err != nil {
		return err
	}

	return
}

func (ps *Storage) Read(ctx c.Context, where string, items interface{}) (
	err error,
) {
	err = ps.db.Table(where).WithContext(ctx).Find(items).Error

	return
}

func (ps *Storage) Remove(ctx c.Context, where string, item interface{}) (
	err error,
) {
	err = ps.db.Table(where).WithContext(ctx).Delete(item).Error

	return
}

func (ps *Storage) Clear(ctx c.Context, where string) (
	err error,
) {
	err = ps.db.Table(where).WithContext(ctx).Delete("*").Error

	return
}

func (ps *Storage) ClearAll(ctx c.Context) (
	err error,
) {
	err = ps.db.Table("*").Delete("*").Error

	return
}

type Serializers []*Serializer

func NewSerializers() Serializers {
	return make([]*Serializer, 0)
}

func (ss Serializers) Set(serializers ...*Serializer) {
	ss = append(ss, New("json", schema.JSONSerializer{}))
	ss = append(ss, New("gob", schema.GobSerializer{}))
	for _, s := range serializers {
		ss = append(ss, s)
	}
}

func (ss Serializers) RegisterAll() {
	for _, s := range ss {
		s.Register()
	}
}

type Serializer struct {
	name string
	si   schema.SerializerInterface
}

func New(name string, si schema.SerializerInterface) *Serializer {
	return &Serializer{
		name: name,
		si:   si,
	}
}

func (s *Serializer) Register() {
	schema.RegisterSerializer(s.name, s.si)
}

func (s *Serializer) GetActions() schema.SerializerInterface {
	return s.si
}

func (s *Serializer) Scan(
	ctx c.Context,
	field *schema.Field,
	dst reflect.Value,
	dbValue interface{},
) (
	err error,
) {
	err = s.GetActions().Scan(ctx, field, dst, dbValue)

	return
}

func (s *Serializer) Value(
	ctx c.Context,
	field *schema.Field,
	dst reflect.Value,
	fieldValue interface{},
) (
	out interface{},
	err error,
) {
	out, err = s.GetActions().Value(ctx, field, dst, fieldValue)

	return
}
