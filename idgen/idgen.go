package idgen

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	epoch          int64 = 1288834974657
	timestampBits  uint  = 41
	machineIdBits  uint  = 10
	sequenceBits   uint  = 12
	maxMachineId         = -1 ^ (-1 << machineIdBits)
	maxSequence          = -1 ^ (-1 << sequenceBits)
	timestampShift       = sequenceBits + machineIdBits
	machineIdShift       = sequenceBits
)

type IdGen struct {
	mu        sync.Mutex
	lastTime  int64
	machineId int64
	sequence  int64
	randMaker *RandomMaker
}

func NewIdGen(machineId int64) (*IdGen, error) {
	if machineId < 0 || machineId > maxMachineId {
		return nil, errors.New("machine ID out of range")
	}
	return &IdGen{
		machineId: machineId,
		randMaker: NewRandomMaker(),
		lastTime:  time.Now().UnixNano(),
	}, nil
}

func (s *IdGen) NextId() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixNano() / 1e6
	if now == s.lastTime {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			for now <= s.lastTime {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}
	s.lastTime = now

	id := ((now - epoch) << timestampShift) |
		(s.machineId << machineIdShift) |
		s.sequence

	return id
}

func (s *IdGen) NextString() string {
	buf := make([]byte, 12)
	s.randMaker.Read(buf)
	str := base64.URLEncoding.EncodeToString(buf)
	// 只保留字母和数字，去掉填充字符
	str = strings.TrimRight(str, "=")
	return str
}

type RandomMaker struct {
	randSrc rand.Source
}

func NewRandomMaker() *RandomMaker {
	return &RandomMaker{
		randSrc: rand.NewSource(time.Now().UnixNano()),
	}
}

// Read satisfies io.Reader
func (s *RandomMaker) Read(p []byte) (n int, err error) {
	todo := len(p)
	offset := 0
	for {
		val := s.randSrc.Int63()
		for i := 0; i < 8; i++ {
			p[offset] = byte(val)
			todo--
			if todo == 0 {
				return len(p), nil
			}
			offset++
			val >>= 8
		}
	}
}
