package srvhttp

import (
	"context"
	"net"
	"testing"

	commonv1 "github.com/mcrgnt/proto/gen/go/common/v1"
	"github.com/stretchr/testify/assert"
)

func TestSome(t *testing.T) {
	var cfg = &Config{
		Port: commonv1.Port{Value: 8080},
	}
	var s, err = cfg.Build()
	assert.NoError(t, err)

	var testCaseList = []struct {
		name      string
		givenCtx  context.Context
		wantError error
	}{
		{
			name: "given context cancelled",
			givenCtx: func() context.Context {
				var ctx, cancel = context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			wantError: context.Canceled,
		},
		{
			name:      "given context background",
			givenCtx:  context.Background(),
			wantError: nil,
		},
	}

	for _, tc := range testCaseList {
		var actualErr = s.(*srv).Start(tc.givenCtx)
		assert.Equal(t, tc.wantError, actualErr)
	}
}

func TestListenerAddrAfterClose(t *testing.T) {
	// 1. Создаем листенер на любом свободном порту
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	// Получаем адрес рабочего листенера
	addrBefore := ln.Addr().String()
	t.Logf("Addr before close: %s", addrBefore)

	// 2. Закрываем листенер
	err = ln.Close()
	if err != nil {
		t.Fatalf("failed to close: %v", err)
	}

	// 3. Проверяем адрес ПОСЛЕ закрытия
	addrAfter := ln.Addr()

	if addrAfter == nil {
		t.Log("Addr is nil after close")
	} else {
		t.Logf("Addr after close: %s", addrAfter.String())
	}

	// Проверка: будет ли addrAfter по-прежнему содержать старый адрес?
	if addrAfter != nil && addrAfter.String() == addrBefore {
		t.Log("Theory confirmed: Addr() still returns the address after Close()")
	}
}
