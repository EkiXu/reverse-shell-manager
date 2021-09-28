package util

import (
	"container/list"
	"errors"

	"sh.ieki.xyz/global"
	"sh.ieki.xyz/model"
)

/*
func wsBindsock(ws websocket.Conn, socks *net.TCPConn) {
	defer socks.Close()

	defer ws.Close()

	var wg sync.WaitGroup
	ioCopy := func(dst io.Writer, src io.Reader) {
		defer wg.Done()
		io.Copy(dst, src)
	}
	wg.Add(2)
	//go ioCopy(ws, socks)
	//go ioCopy(socks, ws)
	wg.Wait()
}*/

func GetListener(listnerList *list.List, targetListener model.Listener) (*model.Listener, error) {
	var listener *model.Listener
	// global.SERVER_LOG.Debugf("finding listener,target listener: %+v", targetListener)
	for element := listnerList.Front(); element != nil; element = element.Next() {
		listener = element.Value.(*model.Listener)

		// global.SERVER_LOG.Debugf("finding listener,now listener: %+v", *listener)
		if (*listener).Name == targetListener.Name {
			return listener, nil
		}
		if (*listener).Host == targetListener.Host && (*listener).Port == targetListener.Port {
			return listener, nil
		}
	}
	return nil, errors.New("listener not found")
}

func PointList2ValArray(list *list.List) ([]interface{}, error) {
	var res []interface{}
	global.SERVER_LOG.Debugf("now list %+v", list)
	for element := list.Front(); element != nil; element = element.Next() {
		p := element.Value.(*interface{})
		global.SERVER_LOG.Debugf("now element %+v", p)
		res = append(res, *p)
	}
	return res, nil
}

func GetBeaconList(becaonList *list.List) ([]model.Beacon, error) {
	var beaconListArray []model.Beacon
	var beacon *model.Beacon
	for element := becaonList.Front(); element != nil; element = element.Next() {
		beacon = element.Value.(*model.Beacon)

		beaconListArray = append(beaconListArray, *beacon)
	}
	return beaconListArray, nil
}

func GetListenerList(listnerList *list.List) ([]model.Listener, error) {
	var listenerListArray []model.Listener
	var listener *model.Listener
	for element := listnerList.Front(); element != nil; element = element.Next() {
		listener = element.Value.(*model.Listener)

		listenerListArray = append(listenerListArray, *listener)
	}
	return listenerListArray, nil
}

func GetBeacon(beaconList *list.List, targetBeacon model.Beacon) (*model.Beacon, error) {
	var beacon *model.Beacon
	//global.SERVER_LOG.Debugf("finding beacon,target beacon: %+v", targetBeacon)
	for element := beaconList.Front(); element != nil; element = element.Next() {
		beacon = element.Value.(*model.Beacon)

		//global.SERVER_LOG.Debugf("finding beacon,now beacon: %+v", *beacon)
		if (*beacon).Name == targetBeacon.Name || (*beacon).UUID == targetBeacon.UUID {
			return beacon, nil
		}
	}
	return nil, errors.New("beacon not found")
}

func DeleteBeacon(beaconList *list.List, targetBeacon *model.Beacon) error {
	var beacon *model.Beacon
	// global.SERVER_LOG.Debugf("finding beacon,target beacon: %+v", targetBeacon)
	for element := beaconList.Front(); element != nil; element = element.Next() {
		beacon = element.Value.(*model.Beacon)

		// global.SERVER_LOG.Debugf("finding beacon,now beacon: %+v", *beacon)
		if (*beacon).Name == targetBeacon.Name || (*beacon).UUID == targetBeacon.UUID {
			beaconList.Remove(element)
			return nil
		}
	}
	return errors.New("beacon not found")
}
