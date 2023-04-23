package registry

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/werbenhu/chash"
)

// Http represents the http server object
type Http struct {
	addr     string       // the address that http server listens to
	listener net.Listener // the listener for the http server
}

// NewHttp returns a new Http object
func NewHttp() *Http {
	return &Http{}
}

// match assigns a service to a key using consistent hashing algorithm
func (h *Http) match(c *gin.Context) {
	name := c.Query("group")
	key := c.Query("key")

	// Get the group based on the provided name
	group, err := chash.GetGroup(name)
	if err != nil {
		// Return error response if group not found
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}

	// Match the key with a member in the group
	_, payload, err := group.Match(key)
	if err != nil {
		// Return error response if key not found
		c.JSON(http.StatusOK, gin.H{
			"code": 2,
			"msg":  err.Error(),
		})
		return
	}

	// Unmarshal the member payload
	m := &Member{}
	if err := m.Unmarshal(payload); err != nil {
		// Return error response if payload cannot be unmarshalled
		c.JSON(http.StatusOK, gin.H{
			"code": 3,
			"msg":  err.Error(),
		})
		return
	}

	// Return success response with the matched service
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"service": m.Service,
		},
	})
}

// members returns the list of services for a group
func (h *Http) members(c *gin.Context) {
	name := c.Query("group")
	// Get the group based on the provided name
	group, err := chash.GetGroup(name)
	if err != nil {
		// Return error response if group not found
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}

	// Get all the elements in the group and extract their services
	elements := group.GetElements()
	services := make([]Service, 0)
	for _, element := range elements {
		m := &Member{}
		if err := m.Unmarshal(element.Payload); err == nil {
			services = append(services, m.Service)
		}
	}

	// Return success response with the list of services
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"services": services,
		},
	})
}

// Start starts the http server
func (h *Http) Start(addr string) error {
	var err error
	h.addr = addr

	r := gin.Default()
	r.GET("/match", h.match)
	r.GET("/members", h.members)

	// Listen on the provided address and run the http server
	h.listener, err = net.Listen("tcp", h.addr)
	if err != nil {
		return err
	}
	return r.RunListener(h.listener)
}

// Stop stops the http server
func (h *Http) Stop() {
	if h.listener != nil {
		h.listener.Close()
	}
}
