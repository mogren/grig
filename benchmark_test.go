package main

import (
	"testing"
)

func BenchmarkGenerateIdentities(b *testing.B) {
	// Initialize data once
	langFlag = "en_us"
	dict := loadData(langFlag)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getNext(dict)
	}
}

func BenchmarkGenerateMultipleIdentities(b *testing.B) {
	// Initialize data once
	langFlag = "en_us"
	dict := loadData(langFlag)

	const numIdentities = 100
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		identities := make([]Rig, 0, numIdentities)
		b.StartTimer()

		for j := 0; j < numIdentities; j++ {
			identities = append(identities, getNext(dict))
		}
	}
}

func BenchmarkAsJSON(b *testing.B) {
	// Initialize data once
	langFlag = "en_us"
	dict := loadData(langFlag)
	rig := getNext(dict)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rig.AsJSON(true)
	}
}

func BenchmarkAsXML(b *testing.B) {
	// Initialize data once
	langFlag = "en_us"
	dict := loadData(langFlag)
	rig := getNext(dict)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rig.AsXML()
	}
}

func BenchmarkLoadData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = loadData("en_us")
	}
}
