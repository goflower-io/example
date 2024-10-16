CREATE TABLE `user` (
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id"',
    `name` varchar(100) NOT NULL COMMENT '名称|text|validate:"max=100,min=10"',
    `age` int(11) NOT NULL DEFAULT '0' COMMENT '年龄|number|validate:"max=140,min=18"',
    `sex` int(11) NOT NULL DEFAULT '2' COMMENT '性别|select|validate:"oneof=0 1 2"|0:女 1:男 2:无',
    `ctime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `mtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `ix_name` (`name`) USING BTREE,
    KEY `ix_mtime` (`mtime`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
/*
--id:'序号|text|validate:"oneof 1 2 3"'
--name:'名称|text|validate:"max=100,min=10"'
--age:'年龄|number|validate:"max=140,min=18"'
--sex:'性别|select|validate:"oneof=0 1 2"|0:女 1:男 2:无'
--ctime:'创建时间'
--mtime:'修改时间'
*/
