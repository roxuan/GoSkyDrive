CREATE TABLE `tbl_file`(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `file_sha1` char(40) NOT NULL DEFAULT '' COMMENT '文件hash',
    `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
    `file_size` bigint(20) DEFAULT '0' COMMENT '文件大小',
    `file_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
    `create_at` datetime default NOW() COMMENT '创建日期',
    `update_at` datetime default NOW() on update current_timestamp() COMMENT '更新日期',
    `stuatus` bigint(20) DEFAULT '0' COMMENT '文件大小',

)