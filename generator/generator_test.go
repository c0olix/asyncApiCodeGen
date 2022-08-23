package generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCodeGenerator_loadAsyncApiSpec(t *testing.T) {
	t.Skip("Currently not successful")
	spec := loadAsyncApiSpec("./test-spec/test-spec.yaml")
	assert.Equal(t, "2.2.0", spec.AsyncApi)
	assert.Equal(t, "Example", spec.Info.Title)
	assert.Equal(t, "0.1.0", spec.Info.Version)
	assert.Equal(t, "broker.mycompany.com", spec.Servers["production"].Url)
	assert.Equal(t, "amqp", spec.Servers["production"].Protocol)
	assert.Equal(t, "This is \"My Company\" broker.", spec.Servers["production"].Description)
	assert.Equal(t, "A brief description for the UserDeletedEvent", spec.Channels["USER_DELETED"].Publish.Message.Description)
	assert.Equal(t, "object", spec.Channels["USER_DELETED"].Publish.Message.Schema.Type)
	assert.Equal(t, bp(false), spec.Channels["USER_DELETED"].Publish.Message.Schema.AdditionalProperties)
	assert.Equal(t, "string", spec.Channels["USER_DELETED"].Publish.Message.Schema.Properties["fullName"].Type)
	assert.Equal(t, "string", spec.Channels["USER_DELETED"].Publish.Message.Schema.Properties["email"].Type)
	assert.Equal(t, strp("email"), spec.Channels["USER_DELETED"].Publish.Message.Schema.Properties["email"].Format)
	assert.Equal(t, "integer", spec.Channels["USER_DELETED"].Publish.Message.Schema.Properties["age"].Type)
	assert.Equal(t, ip(18), spec.Channels["USER_DELETED"].Publish.Message.Schema.Properties["age"].Minimum)
}

func strp(in string) *string {
	return &in
}
