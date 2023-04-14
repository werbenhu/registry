package registry

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/werbenhu/chash"
)

// http server object
type Http struct {
	addr     string
	listener net.Listener
}

// NetHttp new a http server object
func NewHttp() *Http {
	return &Http{}
}

// match assign a service to a key with consistent hashing algorithm
func (h *Http) match(c *gin.Context) {
	name := c.Query("group")
	key := c.Query("key")

	group, err := chash.GetGroup(name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}

	_, payload, err := group.Match(key)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2,
			"msg":  err.Error(),
		})
		return
	}

	m := &Member{}
	if err := m.Unmarshal(payload); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 3,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"service": m.Service,
		},
	})
}

// members() get services list of a group
// groupName:
//
//	the group name of the services
func (h *Http) members(c *gin.Context) {
	name := c.Query("group")
	group, err := chash.GetGroup(name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}

	elements := group.GetElements()
	services := make([]Service, 0)
	for _, element := range elements {
		m := &Member{}
		if err := m.Unmarshal(element.Payload); err == nil {
			services = append(services, m.Service)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"services": services,
		},
	})
}

// Start() start http server
// addr:
//
//	the addr that http server listen to
func (h *Http) Start(addr string) error {
	var err error
	h.addr = addr

	r := gin.Default()
	r.GET("/match", h.match)
	r.GET("/members", h.members)

	h.listener, err = net.Listen("tcp", h.addr)
	if err != nil {
		return err
	}
	return r.RunListener(h.listener)
}

// Stop() stop http server
func (h *Http) Stop() {
	if h.listener != nil {
		h.listener.Close()
	}
}
