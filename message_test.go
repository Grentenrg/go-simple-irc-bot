package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePacket(t *testing.T) {
	type testCase struct {
		name      string
		input     []byte
		want      []string
		wantRest  string
		wantError error
	}

	cases := []testCase{
		{
			name:      "empty",
			input:     []byte{},
			want:      nil,
			wantRest:  "",
			wantError: ErrEmptyPacket,
		},
		{
			name:      "one packet",
			input:     []byte("PING :tmi.twitch.tv\r\n"),
			want:      []string{"PING :tmi.twitch.tv"},
			wantRest:  "",
			wantError: nil,
		},
		{
			name:      "two packets",
			input:     []byte("PING :tmi.twitch.tv\r\nPING :tmi.twitch.tv\r\n"),
			want:      []string{"PING :tmi.twitch.tv", "PING :tmi.twitch.tv"},
			wantRest:  "",
			wantError: nil,
		},
		{
			name:      "three partial packets",
			input:     []byte(".tv\r\nPING :tmi.twitch.tv\r\nPIN"),
			want:      []string{".tv", "PING :tmi.twitch.tv"},
			wantRest:  "PIN",
			wantError: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, rest, err := parsePacket(c.input)
			if err != nil {
				assert.EqualError(t, err, c.wantError.Error())
			}

			assert.Equal(t, c.want, res)
			assert.Equal(t, c.wantRest, rest)

			res2, rest2, err2 := parsePacket(c.input)
			if err2 != nil {
				assert.EqualError(t, err2, c.wantError.Error())
			}
			assert.Equal(t, c.want, res2)
			assert.Equal(t, c.wantRest, rest2)
		})
	}
}

func BenchmarkParsePacket(b *testing.B) {
	benchmarks := []struct {
		name string
		data []byte
	}{
		{"SingleSmallPacket", []byte("small\r\n")},
		{
			"SingleLargePacket",
			[]byte("large_packet_with_more_content_here_to_test_larger_sizes\r\n"),
		},
		{"MultiplePackets", []byte("packet1\r\npacket2\r\npacket3\r\n")},
		{"PacketsWithoutTermination", []byte("packet1\r\npacket2\r\npacket3")},
		{"EmptyPackets", []byte("\r\npacket1\r\n\r\npacket2\r\n")},
		{"LargeMultiPacket", generateLargePackets(100)}, // 100 packets
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs() // Report memory allocations
			b.ResetTimer()   // Reset timer before the loop

			for i := 0; i < b.N; i++ {
				packets, rest, err := parsePacket(bm.data)
				if err != nil {
					b.Fatal(err)
				}
				// Prevent compiler optimization
				_ = packets
				_ = rest
			}
		})
	}
}

// Helper function to generate large test data
func generateLargePackets(n int) []byte {
	// Pre-allocate builder with estimated size
	var result []byte
	packet := []byte("this_is_a_test_packet_with_reasonable_length\r\n")

	// Allocate slice with exact size needed
	result = make([]byte, 0, len(packet)*n)

	// Generate n packets
	for i := 0; i < n; i++ {
		result = append(result, packet...)
	}

	return result
}

// Optional: Benchmark comparison with original function
func BenchmarkParsePacketOld(b *testing.B) {
	// Use same test cases as above
	benchmarks := []struct {
		name string
		data []byte
	}{
		{"SingleSmallPacket", []byte("small\r\n")},
		{
			"SingleLargePacket",
			[]byte("large_packet_with_more_content_here_to_test_larger_sizes\r\n"),
		},
		{"MultiplePackets", []byte("packet1\r\npacket2\r\npacket3\r\n")},
		{"PacketsWithoutTermination", []byte("packet1\r\npacket2\r\npacket3")},
		{"EmptyPackets", []byte("\r\npacket1\r\n\r\npacket2\r\n")},
		{"LargeMultiPacket", generateLargePackets(100)}, // 100 packets
	}

	for _, bm := range benchmarks {
		b.Run("Old_"+bm.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				packets, err := parsePacketOld(bm.data)
				if err != nil {
					b.Fatal(err)
				}
				_ = packets
			}
		})
	}
}
