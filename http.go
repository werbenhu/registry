package srouter

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/werbenhu/chash"
)

type Http struct {
	port     string
	listener net.Listener
}

func NewHttp() *Http {
	return &Http{}
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

	agent := &Agent{}
	if err := agent.Unmarshal(payload); err != nil {
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
			"service": agent.Service,
		},
	})
}

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
		agent := &Agent{}
		if err := agent.Unmarshal(element.Payload); err == nil {
			services = append(services, agent.Service)
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

func (h *Http) Start(port string) error {
	var err error
	h.port = port

	r := gin.Default()
	r.GET("/match", h.match)
	r.GET("/members", h.members)

	h.listener, err = net.Listen("tcp", ":"+h.port)
	if err != nil {
		return err
	}
	return r.RunListener(h.listener)
}

func (h *Http) Stop() {
	h.listener.Close()
}
