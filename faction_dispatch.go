package main

type dispatchRequestStruct struct {
	assignedOrderCode        orderCode
	requester                actor
	targetTileX, targetTileY int
	expirationTick           int
}

func (f *faction) addDispatchRequest(requester actor, ttx, tty int, orderCode orderCode, expirationTick int) {
	for i := f.dispatchRequests.Front(); i != nil; i = i.Next() {
		dr := i.Value.(*dispatchRequestStruct)
		if dr.assignedOrderCode == orderCode && dr.requester == requester && dr.targetTileX == ttx && dr.targetTileY == tty {
			return
		}
	}
	f.dispatchRequests.PushBack(&dispatchRequestStruct{
		assignedOrderCode: orderCode,
		requester:         requester,
		targetTileX:       ttx,
		targetTileY:       tty,
	})
}

func (f *faction) cleanExpiredFactionDispatchRequests(currentTick int) {
	for i := f.dispatchRequests.Front(); i != nil; i = i.Next() {
		if i.Value.(*dispatchRequestStruct).expirationTick < currentTick {
			removingElement := i
			if i.Prev() != nil {
				i = i.Prev()
			}
			f.dispatchRequests.Remove(removingElement)
		}
	}
}

func (f *faction) removeDispatchRequest(dr *dispatchRequestStruct) {
	for i := f.dispatchRequests.Front(); i != nil; i = i.Next() {
		if i.Value.(*dispatchRequestStruct) == dr {
			removingElement := i
			if i.Prev() != nil {
				i = i.Prev()
			}
			f.dispatchRequests.Remove(removingElement)
		}
	}
}
