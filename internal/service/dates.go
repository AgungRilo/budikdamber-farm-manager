package service

import "time"

// Usaha ini di Indonesia, jadi "hari ini" dihitung dalam WIB,
// bukan zona waktu server (Railway memakai UTC).
var wib = func() *time.Location {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.FixedZone("WIB", 7*60*60) // cadangan kalau tzdata tidak ada
	}
	return loc
}()

func todayWIB() string {
	return time.Now().In(wib).Format("2006-01-02")
}

// isFutureDate: tanggal format YYYY-MM-DD bisa dibandingkan langsung sebagai string.
func isFutureDate(date string) bool {
	return date > todayWIB()
}
