package collector

import (
	"net/http"

	"github.com/kaz/flos-hortus/common"
	"github.com/kaz/flos-hortus/database"
	"github.com/labstack/echo/v4"
)

func RegisterHandler(g *echo.Group) {
	g.GET("/instance", getInstances)
	g.PUT("/instance/:addr", putInstance)
	g.DELETE("/instance/:addr", deleteInstance)
}

func getInstances(c echo.Context) error {
	instances := []string{}

	mu.RLock()
	for addr, _ := range workers {
		instances = append(instances, addr)
	}
	mu.RUnlock()

	return c.JSON(http.StatusOK, instances)
}

func putInstance(c echo.Context) error {
	addr := c.Param("addr")
	bastion := c.QueryParam("bastion") != ""

	if _, err := database.DB().Exec("INSERT INTO instances VALUES (?, ?)", addr, bastion); err != nil {
		return err
	}

	if bastion {
		common.RegisterBastion(addr)
	} else {
		go runWorker(addr)
	}

	return c.NoContent(http.StatusOK)
}

func deleteInstance(c echo.Context) error {
	addr := c.Param("addr")

	result, err := database.DB().Exec("DELETE FROM instances WHERE host = ?", addr)
	if err != nil {
		return err
	}

	if affected, err := result.RowsAffected(); err != nil {
		return err
	} else if affected == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	mu.RLock()
	cancel, ok := workers[addr]
	mu.RUnlock()

	if ok {
		cancel()
		logger.Println("worker terminated:", addr)

		mu.Lock()
		delete(workers, addr)
		mu.Unlock()
	}

	return c.NoContent(http.StatusOK)
}
