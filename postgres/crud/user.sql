CREATE TABLE "public"."user" (
    "id" serial NOT NULL PRIMARY KEY,
    "name" varchar(255) NOT NULL, 
    "age" int4 NOT NULL, 
    "address" varchar(255)[] NOT NULL, 
    "ctime" timestamp(6) NOT NULL DEFAULT now(),
    "mtime" timestamp(6) NOT NULL DEFAULT now()  
);


--id:'序号|text|validate:"oneof 1 2 3"'
--name:'名称|text|validate:"max=100,min=10"'
--age:'年龄|number|validate:"max=140,min=18"'
--sex:'性别|select|validate:"oneof=0 1 2"|0:女 1:男 2:无'
--ctime:'创建时间'
--mtime:'修改时间'

COMMENT ON TABLE public.user IS '用户';
COMMENT ON TABLE public.user.id IS 'Id';
COMMENT ON TABLE public.user.name IS '姓名|text|validate:"max=20"';
COMMENT ON TABLE public.user.age IS '年龄|select|validate:"oneof 1 2 3"|1:x 2:y 3:z';
COMMENT ON TABLE public.user.address IS '地址';
COMMENT ON TABLE public.user.ctime IS '创建时间';
COMMENT ON TABLE public.user.mtime IS '修改时间';
