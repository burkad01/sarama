package sarama

import (
	"fmt"
	"time"
)

type SyncGroupResponse struct {
	// Version defines the protocol version to use for encode and decode
	Version int16
	// ThrottleTime contains the duration in milliseconds for which the
	// request was throttled due to a quota violation, or zero if the request
	// did not violate any quota.
	ThrottleTime int32
	// Err contains the error code, or 0 if there was no error.
	Err KError
	// MemberAssignment contains the member assignment.
	MemberAssignment []byte
}

func (r *SyncGroupResponse) GetMemberAssignment() (*ConsumerGroupMemberAssignment, error) {
	assignment := new(ConsumerGroupMemberAssignment)
	err := decode(r.MemberAssignment, assignment, nil)
	return assignment, err
}

func (r *SyncGroupResponse) encode(pe packetEncoder) error {
	if r.Version >= 1 {
		pe.putInt32(r.ThrottleTime)
	}
	pe.putInt16(int16(r.Err))
	return pe.putBytes(r.MemberAssignment)
}

func (r *SyncGroupResponse) decode(pd packetDecoder, version int16) (err error) {
	r.Version = version
	if r.Version >= 1 {
		if r.ThrottleTime, err = pd.getInt32(); err != nil {
			fmt.Printf("[SyncGroupResponse] decode throttle getInt32() error: %v\n", err)
			return err
		}
	}
	kerr, err := pd.getInt16()
	if err != nil {
		fmt.Printf("[SyncGroupResponse] decode getInt16() error: %v\n", err)
		return err
	}

	r.Err = KError(kerr)
	if kerr == -1 {
		fmt.Println("sync_group_response: kErr=-1")
	}

	r.MemberAssignment, err = pd.getBytes()

	if kerr == -1 {
		fmt.Printf("[SyncGroupResponse] r.MemberAssignment=%v err=%v\n", r.MemberAssignment, err)
		fmt.Printf("[SyncGroupResponse] r.MemberAssignment as string=%s", string(r.MemberAssignment))
	}

	return
}

func (r *SyncGroupResponse) key() int16 {
	return 14
}

func (r *SyncGroupResponse) version() int16 {
	return r.Version
}

func (r *SyncGroupResponse) headerVersion() int16 {
	return 0
}

func (r *SyncGroupResponse) isValidVersion() bool {
	return r.Version >= 0 && r.Version <= 3
}

func (r *SyncGroupResponse) requiredVersion() KafkaVersion {
	switch r.Version {
	case 3:
		return V2_3_0_0
	case 2:
		return V2_0_0_0
	case 1:
		return V0_11_0_0
	case 0:
		return V0_9_0_0
	default:
		return V2_3_0_0
	}
}

func (r *SyncGroupResponse) throttleTime() time.Duration {
	return time.Duration(r.ThrottleTime) * time.Millisecond
}
