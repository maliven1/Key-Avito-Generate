package logic

import (
	"avito/internal/logger"
	"avito/internal/storage/sqliteDB"
	"avito/internal/worker"
	"log/slog"
	"os"
	"strings"
	"unicode"

	"math/rand"
)

func Logic(storagePath string, pathtxt string) (string, string) {
	data := ReadFile(pathtxt)
	h, m, l := keyDistribution(data)
	HighKey, MiddleKey, LowKey := SplitKey(h, m, l)
	storage := worker.GetDbKey(storagePath, HighKey, MiddleKey, LowKey)
	return SortKey(storage), GetAllKeys(storage)
}

func ReadFile(filename string) string {
	const op = "internal.logic.ReadFile"
	log := logger.Log.With(
		slog.String("op", op),
	)
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Error("Read error", err)
	}
	return string(content)
}

func keyDistribution(data string) (string, string, string) {
	HighKey, cut, _ := strings.Cut(data, "ВЧ")
	MiddleKey, cut, _ := strings.Cut(cut, "СЧ")
	LowKey, _, _ := strings.Cut(cut, "НЧ")
	return HighKey, MiddleKey, LowKey
}

func SplitKey(HighKey string, MiddleKey string, LowKey string) ([]string, []string, []string) {
	HigtSplit := strings.Split(HighKey, ", ")
	MiddleSplit := strings.Split(MiddleKey, ", ")
	LowSplit := strings.Split(LowKey, ", ")
	return HigtSplit, MiddleSplit, LowSplit
}

func SortKey(db *sqliteDB.Storage) string {
	const op = "internal.logic.GetKey"
	log := logger.Log.With(
		slog.String("op", op),
	)
	first, err := db.GetHightKey()
	if err != nil {
		log.Error("GetHightKey error", err)
	}
	second, err := db.GetMiddleKey()
	if err != nil {
		log.Error("GetMiddleKey error", err)
	}
	thr, err := db.GetLowKey()
	if err != nil {
		log.Error("GetLowKey error", err)
	}
	firstWord := FirstLetter(first[rand.Intn(len(first))])
	secondWord := FirstLetter(second[rand.Intn(len(second))])
	thrWord := FirstLetter(thr[rand.Intn(len(thr))])
	result := firstWord + ", " + secondWord + ", " + thrWord
	return result
}

func FirstLetter(word string) string {
	rs := []rune(word)
	inWord := false
	for i, r := range rs {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			if !inWord {
				rs[i] = unicode.ToTitle(r)
			}
			inWord = true
		} else {
			inWord = false
		}
	}
	return string(rs)
}

func GetAllKeys(db *sqliteDB.Storage) string {
	const op = "internal.logic.GetAllKeys"
	log := logger.Log.With(
		slog.String("op", op),
	)

	highKeys, err := db.GetHightKey()
	if err != nil {
		log.Error("GetHightKey error", err)
		return ""
	}

	middleKeys, err := db.GetMiddleKey()
	if err != nil {
		log.Error("GetMiddleKey error", err)
		return ""
	}

	lowKeys, err := db.GetLowKey()
	if err != nil {
		log.Error("GetLowKey error", err)
		return ""
	}

	// Combine all keys into a single string
	allKeys := append(append(highKeys, middleKeys...), lowKeys...)
	return strings.Join(allKeys, ", ")
}
