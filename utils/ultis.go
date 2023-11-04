package utils

func GetPagination(count int64) int64 {
	page := count / 10
	if count%10 != 0 {
		page++
	}
	return page
}
