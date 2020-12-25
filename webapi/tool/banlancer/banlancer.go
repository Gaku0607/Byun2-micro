package banlancer

import (
	"hash/crc32"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Gaku0607/Byun2-micro/webapi/models"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func MainBanlancer(name string) func() selector.Strategy {
	switch name {
	case "RoundRobin":
		return roundRobin
	case "SmoothWeight":
		return SmoothWeight
	}
	return roundRobin
}

func getnodes(services []*registry.Service) []*registry.Node {
	nodes := make([]*registry.Node, 0)

	for _, service := range services {
		nodes = append(nodes, service.Nodes...)
	}
	return nodes
}

//輪詢算法
func roundRobin() selector.Strategy {

	var i int
	var mtx sync.Mutex

	return func(r []*registry.Service) selector.Next {

		nodes := getnodes(r)

		return func() (*registry.Node, error) {
			if len(nodes) == 0 {
				return nil, selector.ErrNoneAvailable
			}
			mtx.Lock()

			i = (i + 1) % len(nodes)

			node := nodes[i]

			mtx.Unlock()

			return node, nil
		}
	}
}

type srvnode struct {
	w  int       //weight 該node的真正權重
	cw int       //當前權重
	fw int       //當前降權值
	t  time.Time //該node

}

func configParms() (interval time.Duration, ratios int, err error) {
	if val := os.Getenv("weighted_ratios"); val == "" {
		return 0, 0, models.ErrWeightRatiosIsEmtpy
	} else {
		if ratios, err = strconv.Atoi(val); err != nil {
			return 0, 0, err
		}
	}

	if val := os.Getenv("banlancer_interval"); val == "" {
		return 0, 0, models.ErrBanlancerIntervalIsEmtpy
	} else {
		if intervalint, err := strconv.Atoi(val); err != nil {
			return 0, 0, err
		} else {
			interval = time.Duration(intervalint) * time.Second
		}
	}

	return
}

//平滑權重算法
func SmoothWeight() selector.Strategy {

	interval, ratios, err := configParms()
	if err != nil {
		// tool.ErrChan <- err
	}

	var (
		srvnodes map[string]*srvnode = map[string]*srvnode{}
		now      time.Time
		mux      sync.Mutex
	)

	go func() {
		for {
			<-time.After(interval)

			for key, node := range srvnodes {
				if node.t != now {
					mux.Lock()

					if now.Sub(node.t) <= time.Duration(10)*time.Second {
						node.fw = node.w
					} else {
						delete(srvnodes, key)
					}

					mux.Unlock()
				}

			}
		}
	}()

	return func(r []*registry.Service) selector.Next {

		nodes := getnodes(r)

		return func() (service *registry.Node, err error) {

			mux.Lock()
			defer mux.Unlock()

			now = time.Now()

			if len(nodes) == 0 {
				for _, node := range srvnodes {
					node.fw = node.w
				}
				return nil, selector.ErrNoneAvailable
			}

			var (
				exit   bool
				total  int
				sn     *srvnode
				data   string
				weight int
			)

			service = nodes[0]

			for _, node := range nodes {

				if sn, exit = srvnodes[node.Id]; !exit {
					//當該node沒有weights時默認為1
					if data, exit = node.Metadata["weights"]; exit {
						if weight, err = strconv.Atoi(data); err != nil {
							weight = 1
						}
					} else {
						weight = 1
					}
					srvnodes[node.Id] = &srvnode{w: weight}
					sn = srvnodes[node.Id]
				}

				sn.t = now

				total += sn.w

				sn.cw += (sn.w + sn.fw)

				if srvnodes[service.Id].cw < sn.cw {
					service = node
				}
			}
			//返回最大當全權重node前 扣除該node的當前權重
			srvnode := srvnodes[service.Id]

			srvnode.cw -= total
			//當返回的節點為可能的故障節點時等比例減去降權比例
			if srvnode.fw != 0 {
				fw := srvnode.w * (ratios / 10)
				if fw == 0 {
					fw = 1
				}
				if srvnode.fw -= fw; srvnode.fw < 0 {
					sn.fw = 0
				}
			}

			return service, nil
		}
	}
}

//ＩＰ哈希算法
func IPHash(ip string) selector.Strategy {
	return func(r []*registry.Service) selector.Next {

		nodes := getnodes(r)

		return func() (*registry.Node, error) {

			index := len(nodes) % int(crc32.ChecksumIEEE([]byte(ip)))

			if len(nodes) == 0 {
				return nil, selector.ErrNoneAvailable
			}
			return nodes[len(nodes)%index], nil
		}
	}
}
