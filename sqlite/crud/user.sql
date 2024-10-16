CREATE TABLE "user" (
    "id" integer NOT NULL, 
    "name" text NOT NULL,
    "age" integer NOT NULL,
    "ctime" integer NOT NULL,
    "mtime" integer NOT NULL,
    PRIMARY KEY ("id")
);
--id:'序号|text|validate:"oneof 1 2 3"'
--name:'名称|text|validate:"max=100,min=10"'
--age:'年龄|number|validate:"max=140,min=18"'
--sex:'性别|select|validate:"oneof=0 1 2"|0:女 1:男 2:无'
--ctime:'创建时间'
--mtime:'修改时间'
