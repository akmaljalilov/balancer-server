package item

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/vault/sdk/helper/base62"
	"github.com/pingcap/errors"
	"velox/server/src/service/consistenHashing"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func Router(r fiber.Router) {
	r.Get("/", getItemById)
	r.Post("/", postItemById)
}

// @Summary GetInfoUser
// @Tags Admin
// @Security Bearer
// @Description item info

// @Accept json
// @Produce json
// @Router /item [get]
func getItemById(c *fiber.Ctx) error {
	id := c.Query("id")

	server := consistenHashing.GetServerByKey(id)
	if server == nil {
		return errors.New("Error i connect to Mysql")
	}

	res, err := server.Execute(fmt.Sprintf("select * from test.users where id='%s'", id))
	if err != nil {
		return err
	}
	return c.JSON(res["name"].String)
}

// @Summary GetInfoUser
// @Tags Admin
// @Security Bearer
// @Description item info
// @Param input body User true "user info"
// @Accept json
// @Produce json
// @Router /item [post]
func postItemById(c *fiber.Ctx) error {
	user := &User{}

	if err := c.BodyParser(user); err != nil {
		return err
	}

	id, err := base62.Random(16)
	if err != nil {
		return err
	}
	server := consistenHashing.GetServerByKey(id)
	if server == nil {
		return errors.New("Error i connect to Mysql")
	}
	_, err = server.Execute(fmt.Sprintf("INSERT INTO `test`.`users` (`id`,`name`) VALUES ('%s', '%s');", id, user.Name))
	if err != nil {
		return err
	}
	return c.JSON(id)
}
