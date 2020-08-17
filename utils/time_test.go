package sgc7utils

import (
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
)

func Test_TimeI_Now(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockTimeI(ctrl)
	m.EXPECT().Now().Return(time.Unix(1597647832, 0))

	curtime := m.Now()
	if curtime.Unix() != 1597647832 {
		t.Fatalf("Test_buildLogFilename Now %d != %d",
			curtime.Unix(), 1597647832)
	}

	t.Logf("Test_buildLogFilename OK")
}

func Test_FormatNow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockTimeI(ctrl)
	m.EXPECT().Now().Return(time.Unix(1597647832, 0))

	strnow := FormatNow(m)
	if strnow != "2020-08-17_15:03:52" {
		t.Fatalf("Test_FormatNow Now %s != %s",
			strnow, "2020-08-17_15:03:52")
	}

	t.Logf("Test_FormatNow OK")
}
