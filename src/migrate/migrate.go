package migrate

import "velox/server/src/service"

func Start() {
	master, _, master2, _, err := service.ConnectAndGetToServer()
	if err != nil {
		return
	}
	_, _ = master.Execute("CREATE SCHEMA `test` ;")
	_, _ = master.Execute("CREATE TABLE `test`.`users` (`id` VARCHAR(45) NOT NULL, `name` VARCHAR(45) NULL,  PRIMARY KEY (`id`));")

	_, _ = master2.Execute("CREATE SCHEMA `test` ;")
	_, _ = master2.Execute("CREATE TABLE `test`.`users` (`id` VARCHAR(45) NOT NULL, `name` VARCHAR(45) NULL,  PRIMARY KEY (`id`));")
}
