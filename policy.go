package cjsqldriver

import (
	"math/rand"
)

type Policy interface {
	ResolveRead([]*connectionItem) *connectionItem
	ResolveWrite([]*connectionItem) *connectionItem
}

type weightPolicy struct {
	R *rand.Rand
}

func (w weightPolicy) ResolveRead(connPools []*connectionItem) *connectionItem {

	wr := make([]*weightResolve, 0, len(connPools))
	for index := range connPools {
		wr = append(wr, &weightResolve{
			conn:   connPools[index],
			weight: connPools[index].readWeight,
		})
	}

	return w.resolve(wr)
}

func (w weightPolicy) ResolveWrite(connPools []*connectionItem) *connectionItem {

	wr := make([]*weightResolve, 0, len(connPools))
	for index := range connPools {
		wr = append(wr, &weightResolve{
			conn:   connPools[index],
			weight: connPools[index].writeWeight,
		})
	}

	return w.resolve(wr)
}

type weightResolve struct {
	conn   *connectionItem
	weight int
}

func (w weightPolicy) resolve(conns []*weightResolve) *connectionItem {
	type weightRange struct {
		begin int
		end   int
		conn  *connectionItem
	}

	step := 0
	ranges := make([]*weightRange, 0, len(conns))
	for index := range conns {
		temp := conns[index]
		if temp.weight <= 0 {
			continue
		}
		weight := temp.weight
		ranges = append(ranges, &weightRange{
			begin: step,
			end:   step + weight,
			conn:  temp.conn,
		})
		step = step + weight
	}

	randNum := w.R.Intn(step) + 1
	// 从头开始找起，落在哪个区间
	for _, temp := range ranges {
		if randNum > temp.begin && randNum <= temp.end {
			return temp.conn
		}
	}

	return conns[w.R.Intn(len(conns))].conn
}
