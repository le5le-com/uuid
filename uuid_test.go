package uuid

import (
	"testing"
	"time"
)

func TestNewV7(t *testing.T) {
	var arr [1000000]UUID
	for i := 0; i < 1000000; i++ {
		uuid, err := NewV7()
		if err != nil {
			t.Fatal(err)
		}

		if uuid[6]>>4 != 7 {
			t.Errorf("Not UUIDv7 format: %s", uuid)
		}

		arr[i] = uuid

		if i < 16 {
			t.Logf("uuid v7: %s %v", uuid, uuid[:])
		}
	}

	if hasDuplicate(arr[:]) {
		t.Errorf("Has duplicate!")
	}
}

func TestParse(t *testing.T) {
	s := "0189dd43-c284-7f4f-806e-e7d238e9babb"
	uuid, err := Parse(s)
	if err != nil {
		t.Fatal(err)
	}

	if uuid[6]>>4 != 7 || uuid.String() != s {
		t.Errorf("Parse Error: %s", uuid)
	}

	t.Logf("uuid v7: %s %v", uuid, uuid[:])
}

func TestUUIDV7FromObjectID(t *testing.T) {
	objectId := "63ede45a8d0137fc1b631091"
	uuid, err := UUIDV7FromObjectID(objectId)
	if err != nil {
		t.Fatal(err)
	}

	if uuid[6]>>4 != 7 {
		t.Errorf("Not UUIDv7 format: %s", uuid)
	}

	t.Logf("uuid v7: %s %v", uuid, uuid[:])
}

func TestString(t *testing.T) {
	s := "0189dd43-c284-7f4f-806e-e7d238e9babb"
	uuid, err := Parse(s)
	if err != nil {
		t.Fatal(err)
	}

	if uuid.String() != s {
		t.Errorf("Parse Error: %s, %s", s, uuid)
	}

	t.Logf("uuid v7: %s, %s", s, uuid)
}

func TestShortString(t *testing.T) {
	s := "0189dd43c2847f4f806ee7d238e9babb"
	uuid, err := Parse(s)
	if err != nil {
		t.Fatal(err)
	}

	if uuid.ShortString() != s {
		t.Errorf("Parse Error: %s, %s", s, uuid)
	}

	t.Logf("uuid v7: %s, %s", s, uuid)
}

func TestTimeFromV7(t *testing.T) {
	uuid, err := NewV7()
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	ms := uuid.TimeFromV7().Sub(now).Milliseconds()

	if ms > 10 {
		t.Errorf("TimeFromV7 excessive Error: %v, uuid.time=%v, time.Now=%v", ms, uuid.TimeFromV7(), now)
	}

	t.Logf("%v, uuid.time=%v, time.Now=%v ", ms, uuid.TimeFromV7(), now)
}

func TestObjectIDHex(t *testing.T) {
	objectId := "63ede45a8d0137fc1b631091"
	uuid, err := UUIDV7FromObjectID(objectId)
	if err != nil {
		t.Fatal(err)
	}

	if uuid.ObjectIDHex() != objectId {
		t.Errorf("Convert Object Error: uuid.ObjectIDHex=%s, objectId=%s", uuid.ObjectIDHex(), objectId)
	}
}

func hasDuplicate(arr []UUID) bool {
	seen := make(map[UUID]bool)
	for _, val := range arr {
		if seen[val] {
			return true
		}
		seen[val] = true
	}
	return false
}
