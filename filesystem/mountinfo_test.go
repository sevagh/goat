package filesystem

import "testing"

func TestIsSlashBootMounted(t *testing.T) {
	//i'm betting you have /boot

	isMounted, err := IsMounted("/boot")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !isMounted {
		t.Errorf("/boot isn't mounted?")
	}
}

func TestIsFakeThingNotMounted(t *testing.T) {
	isMounted, err := IsMounted("/boost")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if isMounted {
		t.Errorf("/boost shouldn't be mounted?")
	}
}
