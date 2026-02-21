package benchmarks

import (
	"testing"

	"github.com/your-org/p2p-network/pkg/protocol"
	"github.com/your-org/p2p-network/pkg/utils"
)

func BenchmarkMessageEncode(b *testing.B) {
	msg := protocol.NewRequest("benchmark-id", make([]byte, 1024))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := msg.Encode()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMessageDecode(b *testing.B) {
	msg := protocol.NewRequest("benchmark-id", make([]byte, 1024))
	encoded, _ := msg.Encode()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := protocol.DecodeMessage(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMessageEncodeDecode(b *testing.B) {
	msg := protocol.NewRequest("benchmark-id", make([]byte, 1024))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoded, err := msg.Encode()
		if err != nil {
			b.Fatal(err)
		}

		_, err = protocol.DecodeMessage(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCryptoSign(b *testing.B) {
	key, err := utils.GenerateKey()
	if err != nil {
		b.Fatal(err)
	}

	data := make([]byte, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := utils.SignMessage(key, data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCryptoVerify(b *testing.B) {
	key, err := utils.GenerateKey()
	if err != nil {
		b.Fatal(err)
	}

	data := make([]byte, 1024)
	signature, _ := utils.SignMessage(key, data)

	publicKey := &key.PublicKey

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ok := utils.VerifySignature(publicKey, data, signature)
		if !ok {
			b.Fatal("verification failed")
		}
	}
}

func BenchmarkCryptoHash(b *testing.B) {
	data := make([]byte, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.ComputeHash(data)
	}
}

func BenchmarkLoggerInfo(b *testing.B) {
	logger, err := utils.NewLogger("benchmark", utils.LogLevelError)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark log message", "iteration", i)
	}
}

func BenchmarkLoggerDebug(b *testing.B) {
	logger, err := utils.NewLogger("benchmark", utils.LogLevelDebug)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("Benchmark log message", "iteration", i)
	}
}

func BenchmarkSmallMessageEncode(b *testing.B) {
	msg := protocol.NewRequest("id", []byte("hello"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = msg.Encode()
	}
}

func BenchmarkLargeMessageEncode(b *testing.B) {
	payload := make([]byte, 64*1024)
	msg := protocol.NewRequest("id", payload)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = msg.Encode()
	}
}
