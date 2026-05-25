package entity

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

var (
	globalRand *rand.Rand
	randMutex  sync.Mutex
)

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	globalRand = rand.New(source)
}

var syllables = []string{
	"a", "ar", "ax", "ba", "br", "ca", "cr", "da", "dr", "el",
	"en", "fa", "fl", "ga", "gr", "ha", "il", "in", "ir", "ka",
	"la", "li", "lo", "ma", "mi", "na", "ne", "ni", "or", "pa",
	"qu", "ra", "ri", "ro", "sa", "se", "si", "ta", "te", "ti",
	"ul", "va", "ve", "vi", "xa", "ze", "zo", "zu", "th", "sh", "ch",
}

var firstNames = []string{
	"Aria", "Bex", "Cael", "Dax", "Evo", "Flux", "Grit", "Halo",
	"Iris", "Jinx", "Kai", "Luna", "Mox", "Nyx", "Orion", "Pax",
	"Quin", "Rex", "Sage", "Tess", "Umbra", "Vex", "Wren", "Zed",
}

var lastNames = []string{
	"Storm", "Shadow", "Steel", "Night", "Fire", "Ice", "Wolf",
	"Raven", "Blade", "Phoenix", "Viper", "Dragon", "Ghost",
}

func GenerateName() string {
	randMutex.Lock()
	defer randMutex.Unlock()

	r := globalRand.Float32()
	if r < 0.5 {
		return genericName()
	}
	return phoneticName()
}

func phoneticName() string {
	numSyllables := globalRand.Intn(2) + 2
	parts := make([]string, numSyllables)
	for i := range parts {
		parts[i] = syllables[globalRand.Intn(len(syllables))]
	}
	name := strings.Join(parts, "")
	if len(name) > 0 {
		name = strings.ToUpper(string(name[0])) + name[1:]
	}
	if globalRand.Float32() < 0.1 {
		suffixes := []string{"on", "ix", "us", "a", "os", "en"}
		name += suffixes[globalRand.Intn(len(suffixes))]
	}
	return name
}

func genericName() string {
	first := firstNames[globalRand.Intn(len(firstNames))]
	last := lastNames[globalRand.Intn(len(lastNames))]
	return first + " " + last
}
