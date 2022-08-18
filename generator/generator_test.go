package generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCodeGenerator_loadAsyncApiSpec(t *testing.T) {
	gen := ServerCodeGenerator{}

	spec := gen.loadAsyncApiSpec("./test-spec/test-spec.yaml")
	assert.Equal(t, "2.2.0", spec.AsyncApi)
	assert.Equal(t, "Example", spec.Info.Title)
	assert.Equal(t, "0.1.0", spec.Info.Version)
	assert.Equal(t, "broker.mycompany.com", spec.Servers["production"].Url)
	assert.Equal(t, "amqp", spec.Servers["production"].Protocol)
	assert.Equal(t, "This is \"My Company\" broker.", spec.Servers["production"].Description)
	assert.Equal(t, "An event describing that a user just signed up.", spec.Channels["USER_SIGNUP"].Subscribe.Message.Description)
	assert.Equal(t, "object", spec.Channels["USER_SIGNUP"].Subscribe.Message.Payload.Type)
	assert.Equal(t, false, spec.Channels["USER_SIGNUP"].Subscribe.Message.Payload.AdditionalProperties)
	assert.Equal(t, "string", spec.Channels["USER_SIGNUP"].Subscribe.Message.Payload.Properties["fullName"].Type)
	assert.Equal(t, "string", spec.Channels["USER_SIGNUP"].Subscribe.Message.Payload.Properties["email"].Type)
	assert.Equal(t, "email", spec.Channels["USER_SIGNUP"].Subscribe.Message.Payload.Properties["email"].Format)
	assert.Equal(t, "integer", spec.Channels["USER_SIGNUP"].Subscribe.Message.Payload.Properties["age"].Type)
	assert.Equal(t, 18, spec.Channels["USER_SIGNUP"].Subscribe.Message.Payload.Properties["age"].Minimum)
}
