package utils

type UpdateTimeComparable interface {
	GetID() int
	GetUpdateTime() int64
}

func CheckUpsert[T UpdateTimeComparable](oldConfs, newConfs []T) (map[int]bool, map[int]T) {
	delM := make(map[int]bool, 0)
	for _, v1 := range oldConfs {
		id := v1.GetID()
		bFind := false
		for _, v2 := range newConfs {
			nID := v2.GetID()
			if nID == id {
				bFind = true
				break
			}
		}
		if bFind == false {
			delM[id] = true
		}
	}
	upsertM := make(map[int]T, 0)
	for _, v1 := range newConfs {
		id := v1.GetID()
		bFind := false
		for _, v2 := range oldConfs {
			nID := v2.GetID()
			if nID == id {
				if v2.GetUpdateTime() < v1.GetUpdateTime() {
					upsertM[id] = v1
				}
				bFind = true
				break
			}
		}
		if bFind == false {
			upsertM[id] = v1
		}
	}
	return delM, upsertM
}
