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

var _ Reactor = &SimpleReactor{}

func TestSimpleReactor_HappyPath_SingleRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDManager := mocks.NewMockIDManager(ctrl)
	mockPacketStream := mocks.NewMockPacketStream(ctrl)

	fakeRequest := hsmlib.Packet{
		Header:  []byte(gofakeit.Lexify("????")),
		Payload: []byte(gofakeit.BuzzWord()),
	}
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
	reactor := SimpleReactor{}

	reactor.IDManager = mockIDManager
	reactor.Target = mockPacketStream

	err := reactor.Start()
	require.NoError(t, err)

	resp, err := reactor.Post(fakeRequest.Payload)
	require.NoError(t, err)
	assert.Equal(t, fakeResponse.Payload, resp)

	reactor.Wait()
}

func TestSimpleReactor_HappyPath_TwoRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIDManager := mocks.NewMockIDManager(ctrl)
	mockPacketStream := mocks.NewMockPacketStream(ctrl)

	fakeRequest1 := hsmlib.Packet{
		Header:  []byte(gofakeit.Lexify("????")),
		Payload: []byte(gofakeit.BuzzWord()),
	}
	fakeResponse1 := hsmlib.Packet{
		Header:  []byte(gofakeit.Lexify("????")),
		Payload: []byte(gofakeit.BuzzWord()),
	}
	fakeResponseChan1 := make(chan []byte)
	fakeRequest2 := hsmlib.Packet{
		Header:  []byte(gofakeit.Lexify("????")),
		Payload: []byte(gofakeit.BuzzWord()),
	}
	fakeResponse2 := hsmlib.Packet{
		Header:  []byte(gofakeit.Lexify("????")),
		Payload: []byte(gofakeit.BuzzWord()),
	}
	fakeResponseChan2 := make(chan []byte)

	mockIDManager.EXPECT().NewID().Return(fakeRequest1.Header, fakeResponseChan1)
	mockIDManager.EXPECT().NewID().Return(fakeRequest2.Header, fakeResponseChan2)
	mockPacketStream.EXPECT().SendPacket(fakeRequest1).Return(nil)
	mockPacketStream.EXPECT().ReceivePacket().Do(func() {
		time.Sleep(10 * time.Millisecond)
	}).Return(fakeResponse1, nil)
	mockPacketStream.EXPECT().SendPacket(fakeRequest2).Return(nil)
	mockPacketStream.EXPECT().ReceivePacket().Do(func() {
		time.Sleep(10 * time.Millisecond)
	}).Return(fakeResponse2, nil)
	mockPacketStream.EXPECT().ReceivePacket().Return(hsmlib.Packet{}, io.EOF)
	mockIDManager.EXPECT().FindChannel(fakeResponse1.Header).Return(fakeResponseChan1, true)
	mockIDManager.EXPECT().FindChannel(fakeResponse2.Header).Return(fakeResponseChan2, true)

	reactor := SimpleReactor{
		IDManager: mockIDManager,
		Target:    mockPacketStream,
	}

	err := reactor.Start()
	require.NoError(t, err)

	resp, err := reactor.Post(fakeRequest1.Payload)
	require.NoError(t, err)
	assert.Equal(t, fakeResponse1.Payload, resp)

	resp, err = reactor.Post(fakeRequest2.Payload)
	require.NoError(t, err)
	assert.Equal(t, fakeResponse2.Payload, resp)

	reactor.Wait()
}
