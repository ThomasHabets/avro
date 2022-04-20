package avro_test

import (
	"bytes"
	"testing"

	"github.com/ThomasHabets/avro"
	"github.com/stretchr/testify/assert"
)

func TestEncoder_MapInvalidType(t *testing.T) {
	defer ConfigTeardown()

	schema := `{"type":"map", "values": "string"}`
	buf := bytes.NewBuffer([]byte{})
	enc, err := avro.NewEncoder(schema, buf)
	assert.NoError(t, err)

	err = enc.Encode("test")

	assert.Error(t, err)
}

func TestEncoder_Map(t *testing.T) {
	defer ConfigTeardown()

	schema := `{"type":"map", "values": "string"}`
	buf := bytes.NewBuffer([]byte{})
	enc, err := avro.NewEncoder(schema, buf)
	assert.NoError(t, err)

	err = enc.Encode(map[string]string{})

	assert.NoError(t, err)
	assert.Equal(t, []byte{0x00}, buf.Bytes())
}

func TestEncoder_MapEmpty(t *testing.T) {
	defer ConfigTeardown()

	schema := `{"type":"map", "values": "string"}`
	buf := bytes.NewBuffer([]byte{})
	enc, err := avro.NewEncoder(schema, buf)
	assert.NoError(t, err)

	err = enc.Encode(map[string]string{"foo": "foo"})

	assert.NoError(t, err)
	assert.Equal(t, []byte{0x01, 0x10, 0x06, 0x66, 0x6F, 0x6F, 0x06, 0x66, 0x6F, 0x6F, 0x00}, buf.Bytes())
}

func TestEncoder_MapOfStruct(t *testing.T) {
	defer ConfigTeardown()

	schema := `{"type":"map", "values": {"type": "record", "name": "test", "fields" : [{"name": "a", "type": "long"}, {"name": "b", "type": "string"}]}}`
	buf := bytes.NewBuffer([]byte{})
	enc, err := avro.NewEncoder(schema, buf)
	assert.NoError(t, err)

	err = enc.Encode(map[string]TestRecord{"foo": {A: 27, B: "foo"}})

	assert.NoError(t, err)
	assert.Equal(t, []byte{0x01, 0x12, 0x06, 0x66, 0x6F, 0x6F, 0x36, 0x06, 0x66, 0x6f, 0x6f, 0x0}, buf.Bytes())
}

func TestEncoder_MapInvalidKeyType(t *testing.T) {
	defer ConfigTeardown()

	schema := `{"type":"map", "values": "string"}`
	buf := bytes.NewBuffer([]byte{})
	enc, err := avro.NewEncoder(schema, buf)
	assert.NoError(t, err)

	err = enc.Encode(map[int]string{1: "foo"})

	assert.Error(t, err)
}

func TestEncoder_MapError(t *testing.T) {
	defer ConfigTeardown()

	schema := `{"type":"map", "values": "string"}`
	buf := bytes.NewBuffer([]byte{})
	enc, err := avro.NewEncoder(schema, buf)
	assert.NoError(t, err)

	err = enc.Encode(map[string]int{"foo": 1})

	assert.Error(t, err)
}

func TestEncoder_MapWithMoreThanBlockLengthKeys(t *testing.T) {
	avro.DefaultConfig = avro.Config{
		TagKey:               "avro",
		BlockLength:          1,
		UnionResolutionError: true,
	}.Freeze()

	schema := `{"type":"map", "values": "int"}`
	buf := bytes.NewBuffer([]byte{})
	enc, err := avro.NewEncoder(schema, buf)
	assert.NoError(t, err)

	err = enc.Encode(map[string]int{"foo": 1, "bar": 2})

	assert.NoError(t, err)
	assert.Condition(t, func() bool {
		// {"foo": 1, "bar": 2}
		foobar := bytes.Equal([]byte{0x01, 0x0a, 0x06, 0x66, 0x6F, 0x6F, 0x02, 0x01, 0x0a, 0x06, 0x62, 0x61, 0x72, 0x04, 0x0}, buf.Bytes())
		// {"bar": 2, "foo": 1}
		barfoo := bytes.Equal([]byte{0x01, 0x0a, 0x06, 0x62, 0x61, 0x72, 0x04, 0x01, 0x0a, 0x06, 0x66, 0x6F, 0x6F, 0x02, 0x0}, buf.Bytes())
		return (foobar || barfoo)
	})
}
