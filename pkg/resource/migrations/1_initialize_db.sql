-- +migrate Up
CREATE TABLE `person` (
                        `id`            VARCHAR(32) NOT NULL,
                        `name`         VARCHAR(32) NOT NULL COMMENT     '姓名',,
                        `create_time` DATETIME NOT NULL,
                        `update_time` DATETIME NOT NULL,
                        PRIMARY KEY (`id`)
)
  COLLATE='utf8mb4_general_ci'
  ENGINE=InnoDB COMMENT '';


-- +migrate Down
DROP TABLE `lessor_unit`;
