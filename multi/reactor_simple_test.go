package multi

import (
	"io"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/hsmlib"
	"github.com/mniak/hsmlib/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var _ Reactor = &_SimpleReactor{}

func TestSimpleReactor_HappyPath_SingleRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDManager := mocks.NewMockIDManager(ctrl)
	mockPacketStream := mocks.NewMockPacketStream(ctrl)
	mockCloser := mocks.NewMockCloser(ctrl)

	fakeRequest := hsmlib.Packet{
		Header:  []byte(gofakeit.Lexify("????")),
		Payload: []byte(gofakeit.BuzzWord()),
	}
	gofakeit.Struct(&fakeRequest)
	fakeResponse := hsmlib.Packet{
		Header:  []byte(gofakeit.Lexify("????")),
		Payload: []byte(gofakeit.BuzzWord()),
	}
	fakeResponseChan := make(chan []byte)

	mockIDManager.EXPECT().NewID().Return(fakeRequest.Header, fakeResponseChan)
	mockPacketStream.EXPECT().SendPacket(fakeRequest).Return(nil)
	mockPacketStream.EXPECT().ReceivePacket().Do(func() {
		time.Sleep(10 * time.Millisecond)
	}).Return(fakeResponse, nil)
	mockPacketStream.EXPECT().ReceivePacket().Return(hsmlib.Packet{}, io.EOF)
	mockIDManager.EXPECT().FindChannel(fakeResponse.Header).Return(fakeResponseChan, true)
	mockCloser.EXPECT().Close()

	reactor := NewSimpleReactor(nil)
	require.NotNil(t, reactor)
	require.NotEmpty(t, reactor)

	reactor.idManager = mockIDManager
	reactor.target = mockPacketStream
	reactor.connectionCloser = mockCloser

	err := reactor.Start()
	require.NoError(t, err)

	resp, err := reactor.Post(fakeRequest.Payload)
	require.NoError(t, err)
	assert.Equal(t, fakeResponse.Payload, resp)

	reactor.Wait()
}
