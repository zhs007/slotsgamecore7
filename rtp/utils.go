package sgc7rtp

func findResults(lst []*RTPReturnData, ret float64) *RTPReturnData {
	for _, v := range lst {
		if v.Return == ret {
			return v
		}
	}

	return nil
}

func addResults(lst []*RTPReturnData, ret float64, times int64) []*RTPReturnData {
	d := findResults(lst, ret)
	if d == nil {
		crd := &RTPReturnData{
			Return: ret,
			Times:  times,
		}

		lst = append(lst, crd)

		return lst
	}

	d.Times += times

	return lst
}

func addResults2(lst []*RTPReturnData, iret int64, times int64) []*RTPReturnData {
	fret := float64(iret) / 100
	d := findResults(lst, fret)
	if d == nil {
		crd := &RTPReturnData{
			Return: fret,
			Times:  times,
			Total:  fret,
		}

		lst = append(lst, crd)

		return lst
	}

	d.Times += times
	d.Total += fret

	return lst
}
