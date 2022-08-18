package scheduler

import (
	"fmt"
	"sync"
)

const (
	groupPrefix     = "Group_"
	groupFullValMax = 1024
	groupFullValMin = 900
)

var (
	// 边缘节点组(根据区域分组,每个组上行带宽为1GB)
	groupMap sync.Map // {key:groupID,val:*Group}

	geoGroupMap sync.Map // {key:geo,val:[]string{groupID}}
	// 上行宽带未满的节点组
	lessFullGroupMap sync.Map // {key:geo,val:map{key:groupID,val:bandwidth}}
	// 节点所在分组记录
	groupIDMap sync.Map // {key:deviceID,val:GroupID}

	groupCount int
)

// Group edge group
type Group struct {
	groupID        string
	edgeNodeMap    map[string]int
	toatlBandwidth int
	isFull         bool
}

// NewGroup new group
func NewGroup(groupID string) *Group {
	group := &Group{groupID: groupID, edgeNodeMap: make(map[string]int)}

	g, ok := groupMap.LoadOrStore(groupID, group)
	if ok && g != nil {
		return g.(*Group)
	}

	return group
}

func (g *Group) addEdge(dID string, bandwidth int) {
	if g.isFull {
		return
	}

	if g.toatlBandwidth+bandwidth > groupFullValMax {
		return
	}

	if _, ok := g.edgeNodeMap[dID]; ok {
		return
	}

	g.toatlBandwidth += bandwidth
	g.isFull = g.toatlBandwidth >= groupFullValMin

	g.edgeNodeMap[dID] = bandwidth
	// log.Infof("gID:%v,toatlBandwidth:%v,map:%v", g.GroupID, g.toatlBandwidth, g.edgeNodeMap)
}

func (g *Group) delEdge(dID string) {
	if bandwidth, ok := g.edgeNodeMap[dID]; ok {
		g.toatlBandwidth -= bandwidth
		g.isFull = g.toatlBandwidth >= groupFullValMin

		delete(g.edgeNodeMap, dID)
	}
}

func newGroupName() string {
	groupCount++

	return fmt.Sprintf("%s%d", groupPrefix, groupCount)
}

// 边缘节点分组
func edgeGrouping(node EdgeNode) string {
	deviceID := node.deviceInfo.DeviceId

	// 如果已经存在分组里 则不需要分组
	oldGroupID, ok := groupIDMap.Load(deviceID)
	if ok && oldGroupID != nil {
		g := oldGroupID.(string)
		return g
	}

	groupID := ""
	defer groupIDMap.Store(deviceID, groupID)

	geoKey := node.geoInfo.Geo
	bandwidth := node.bandwidth

	groups := loadGeoGroupMap(geoKey)
	if groups != nil {
		// 看看有没有未满的组可以加入
		lessFullMap := loadLessFullMap(geoKey)
		if lessFullMap != nil {
			findGroupID := ""
			bandwidthT := 0
			for groupID, bandwidth := range lessFullMap {
				bandwidthT = bandwidth + node.bandwidth

				if bandwidthT <= groupFullValMax {
					findGroupID = groupID
					break
				}
			}

			if findGroupID != "" {
				// 未满的组能加入
				groupID = addGroup(geoKey, deviceID, findGroupID, bandwidth, lessFullMap, groups)
			} else {
				// 未满的组不能加入
				groupID = addGroup(geoKey, deviceID, "", bandwidth, lessFullMap, groups)
			}
		} else {
			groupID = addGroup(geoKey, deviceID, "", bandwidth, nil, groups)
		}
	} else {
		groupID = addGroup(geoKey, deviceID, "", bandwidth, nil, nil)
	}

	return groupID
}

func addGroup(geoKey, deviceID, groupID string, bandwidth int, lessFullMap map[string]int, groups []string) string {
	if groups == nil {
		groups = make([]string, 0)
	}

	group := loadGroupMap(groupID)
	if groupID == "" {
		groupID = newGroupName()
		group = NewGroup(groupID)
		groups = append(groups, groupID)
	}

	group.addEdge(deviceID, bandwidth)

	storeGeoGroupMap(geoKey, groups)

	if lessFullMap == nil {
		lessFullMap = make(map[string]int)
	}

	if !group.isFull {
		// 如果组内带宽未满 则保存到未满map
		lessFullMap[groupID] = group.toatlBandwidth
		storeLessFullMap(geoKey, lessFullMap)
	} else {
		delete(lessFullMap, groupID)
		storeLessFullMap(geoKey, lessFullMap)
	}

	return groupID
}

func loadGroupMap(groupID string) *Group {
	group, ok := groupMap.Load(groupID)
	if ok && group != nil {
		return group.(*Group)
	}

	return nil
}

func loadGeoGroupMap(geoKey string) []string {
	groups, ok := geoGroupMap.Load(geoKey)
	if ok && groups != nil {
		return groups.([]string)
	}

	return nil
}

func storeGeoGroupMap(geoKey string, val []string) {
	geoGroupMap.Store(geoKey, val)
}

func loadLessFullMap(geoKey string) map[string]int {
	groups, ok := lessFullGroupMap.Load(geoKey)
	if ok && groups != nil {
		return groups.(map[string]int)
	}

	return nil
}

func storeLessFullMap(geoKey string, val map[string]int) {
	lessFullGroupMap.Store(geoKey, val)
}

// PrintlnMap Println
func testPrintlnEdgeGroupMap() {
	log.Info("geoGroupMap--------------------------------")
	geoGroupMap.Range(func(key, value interface{}) bool {
		g := key.(string)
		groups := value.([]string)
		log.Info("geo:", g)

		for _, gID := range groups {
			group := loadGroupMap(gID)
			if group == nil {
				continue
			}
			log.Info("gId:", gID, ",group:", group, ",bandwidth:", group.toatlBandwidth)
		}

		return true
	})

	log.Info("groupLessFullMap--------------------------------")
	lessFullGroupMap.Range(func(key, value interface{}) bool {
		g := key.(string)
		groups := value.(map[string]int)
		log.Info("geo:", g)

		for gID, bb := range groups {
			log.Info("gId:", gID, ",bandwidth:", bb)
		}

		return true
	})
}
