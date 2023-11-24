package syncer

import (
	"fmt"
	"testing"
	"time"

	"github.com/hawkv6/hawkwing/internal/test"
	"github.com/hawkv6/hawkwing/pkg/bpf"
	"github.com/hawkv6/hawkwing/pkg/bpf/client"
	"github.com/hawkv6/hawkwing/pkg/entities"
	"github.com/hawkv6/hawkwing/pkg/maps"
	"github.com/hawkv6/hawkwing/pkg/messaging"
	"go.uber.org/mock/gomock"
)

func TestSyncer_NewSyncer(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBpf := bpf.NewMockBpf(ctrl)
	mockBpfReader := client.NewMockClientBpfReader(ctrl)
	mockBpfReader.EXPECT().ReadClientBpfSpecs().Return(&test.MockClientCollectionSpec, nil)
	adapterChannels := messaging.NewAdapterChannels()
	clientMap, err := maps.NewClientMap(mockBpf, mockBpfReader)
	if err != nil {
		t.Fatalf("could not create client map: %v", err)
	}

	syncer := NewSyncer(mockBpf, adapterChannels, clientMap)
	if syncer.bpf != mockBpf {
		t.Errorf("NewSyncer() bpf = %v, want %v", syncer.bpf, mockBpf)
	}
	if syncer.adapterChannels != adapterChannels {
		t.Errorf("NewSyncer() adapterChannels = %v, want %v", syncer.adapterChannels, adapterChannels)
	}
	if syncer.cm != clientMap {
		t.Errorf("NewSyncer() cm = %v, want %v", syncer.cm, clientMap)
	}
	if syncer.reqChan == nil {
		t.Errorf("NewSyncer() reqChan = %v, want %v", syncer.reqChan, make(chan *entities.PathRequest))
	}
	if syncer.resolver == nil {
		t.Errorf("NewSyncer() resolver = %v, want %v", syncer.resolver, &ResolverService{})
	}
	if syncer.ErrCh == nil {
		t.Errorf("NewSyncer() ErrCh = %v, want %v", syncer.ErrCh, make(chan error))
	}
}

func TestSyncer_handleIntentMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBpf := bpf.NewMockBpf(ctrl)
	mockBpfReader := client.NewMockClientBpfReader(ctrl)
	mockBpfReader.EXPECT().ReadClientBpfSpecs().Return(&test.MockClientCollectionSpec, nil)
	adapterChannels := messaging.NewAdapterChannels()
	clientMap, err := maps.NewClientMap(mockBpf, mockBpfReader)
	if err != nil {
		t.Fatalf("could not create client map: %v", err)
	}

	syncer := NewSyncer(mockBpf, adapterChannels, clientMap)

	go syncer.handleIntentMessages()

	syncer.reqChan <- &entities.PathRequest{
		Ipv6DestinationAddress: "2001:db8::1",
		Intents:                []entities.Intent{},
	}

	select {
	case receivedRequest := <-adapterChannels.ChAdapterIntentRequest:
		if receivedRequest.Ipv6DestinationAddress != "2001:db8::1" {
			t.Errorf("receivedRequest = %v, want %v", receivedRequest, "2001:db8::1")
		}
	case <-time.After(time.Second * 1):
		t.Error("Timeout: did not receive message on ChAdapterIntentRequest")
	}

	mockBpf.EXPECT().LoadPinnedMap(gomock.Any()).Return(nil, fmt.Errorf("error")).AnyTimes()
	adapterChannels.ChAdapterIntentResponse <- &entities.PathResult{
		Ipv6DestinationAddress: "2001:db8::1",
		Ipv6SidAddresses:       []string{"2001:db8::1"},
	}

}

func TestSyncer_FetchAll(t *testing.T) {
	test.SetupTestConfig(t)
	ctrl := gomock.NewController(t)
	mockBpf := bpf.NewMockBpf(ctrl)
	mockBpfReader := client.NewMockClientBpfReader(ctrl)
	mockBpfReader.EXPECT().ReadClientBpfSpecs().Return(&test.MockClientCollectionSpec, nil)
	adapterChannels := messaging.NewAdapterChannels()
	clientMap, err := maps.NewClientMap(mockBpf, mockBpfReader)
	if err != nil {
		t.Fatalf("could not create client map: %v", err)
	}

	syncer := NewSyncer(mockBpf, adapterChannels, clientMap)

	go syncer.FetchAll()

	var pathRequests []*entities.PathRequest
	select {
	case receivedRequest := <-syncer.reqChan:
		pathRequests = append(pathRequests, receivedRequest)
	case <-time.After(time.Second * 1):
		t.Error("Timeout: did not receive message on reqChan")
	}

	if len(pathRequests) != 1 {
		t.Errorf("len(pathRequests) = %v, want %v", len(pathRequests), 3)
	}
}
