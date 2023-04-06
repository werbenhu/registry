package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/werbenhu/chash"
	"github.com/werbenhu/srouter/discovery"
)

type Http struct {
	addr string
}

func New(addr string) *Http {
	return &Http{
		addr: addr,
	}
}

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

	agent := &discovery.Agent{}
	if err := agent.Unmarshal(payload); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 3,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "success",
		"data": gin.H{
			"service": agent.Service,
		},
	})
}

func (h *Http) elements(c *gin.Context) {
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
	services := make([]discovery.Service, 0)
	for _, element := range elements {
		agent := &discovery.Agent{}
		if err := agent.Unmarshal(element.Payload); err == nil {
			services = append(services, agent.Service)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "success",
		"data": gin.H{
			"services": services,
		},
	})
}

func (h *Http) Start() {
	r := gin.Default()
	r.GET("/match", h.match)
	r.GET("/elements", h.elements)
	go r.Run(h.addr)
}
