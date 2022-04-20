package avro_test

import "github.com/ThomasHabets/avro"

func ConfigTeardown() {
	// Reset the caches
	avro.DefaultConfig = avro.Config{}.Freeze()
}
