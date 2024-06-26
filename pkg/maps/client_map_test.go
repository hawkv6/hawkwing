package maps

import (
	"fmt"
	"testing"

	"github.com/cilium/ebpf"
	"github.com/hawkv6/hawkwing/internal/test"
	"github.com/hawkv6/hawkwing/pkg/bpf"
	"github.com/hawkv6/hawkwing/pkg/bpf/client"
	"go.uber.org/mock/gomock"
)

func TestNewClientMap(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBpf := bpf.NewMockBpf(ctrl)
	mockClientBpfReader := client.NewMockClientBpfReader(ctrl)

	type args struct {
		bpf       bpf.Bpf
		bpfReader client.ClientBpfReader
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		mockReader func()
	}{
		{
			name: "client bpf reader returns error",
			args: args{
				bpf:       mockBpf,
				bpfReader: mockClientBpfReader,
			},
			wantErr: true,
			mockReader: func() {
				mockClientBpfReader.EXPECT().ReadClientBpfSpecs().Return(&ebpf.CollectionSpec{}, fmt.Errorf("error"))
			},
		},
		{
			name: "client bpf reader returns specs",
			args: args{
				bpf:       mockBpf,
				bpfReader: mockClientBpfReader,
			},
			wantErr: false,
			mockReader: func() {
				mockClientBpfReader.EXPECT().ReadClientBpfSpecs().Return(&test.MockClientCollectionSpec, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockReader()
			_, err := NewClientMap(tt.args.bpf, tt.args.bpfReader)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientMap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientMap_Create(t *testing.T) {
	test.SetupTestConfig(t)
	ctrl := gomock.NewController(t)
	mockBpf := bpf.NewMockBpf(ctrl)
	mockClientBpfReader := client.NewMockClientBpfReader(ctrl)
	mockClientBpfReader.EXPECT().ReadClientBpfSpecs().Return(&test.MockClientCollectionSpec, nil)
	clientDataMap, err := NewClientMap(mockBpf, mockClientBpfReader)
	if err != nil {
		t.Errorf("NewClientMap() should not return error: %s", err)
	}

	tests := []struct {
		name    string
		wantErr bool
		mockBpf func()
	}{
		{
			name:    "create lookup map returns error",
			wantErr: true,
			mockBpf: func() {
				mockBpf.EXPECT().CreateMap(gomock.Any(), gomock.Any()).Return(nil, nil).Times(5)
				mockBpf.EXPECT().CreateMap(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error")).Times(1)
			},
		},
		{
			name:    "create outer map returns error",
			wantErr: true,
			mockBpf: func() {
				mockBpf.EXPECT().CreateMap(gomock.Any(), gomock.Any()).Return(nil, nil).Times(4)
				mockBpf.EXPECT().CreateMap(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error")).Times(1)
			},
		},
		{
			name:    "create reverse map returns error",
			wantErr: true,
			mockBpf: func() {
				mockBpf.EXPECT().CreateMap(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error")).Times(1)
			},
		},
		{
			name:    "no error",
			wantErr: false,
			mockBpf: func() {
				mockBpf.EXPECT().CreateMap(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBpf()
			err := clientDataMap.Create()
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			ctrl.Finish()
		})
	}
}
