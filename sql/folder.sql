DROP DATABASE IF EXISTS `folder`;
CREATE DATABASE `folder`;
USE `folder`;

DROP TABLE IF EXISTS `folder_file`;
CREATE TABLE `folder_file` (
  `file_id` int(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '文件id',
  `file_name` varchar(255) NOT NULL DEFAULT '' COMMENT '文件名称',
  `file_type` varchar(20) NOT NULL DEFAULT '' COMMENT '文件类型',
  `file_status` varchar(20) NOT NULL DEFAULT 'valid' COMMENT '文件状态',
  `file_path` varchar(255) NOT NULL DEFAULT '' COMMENT '文件位置',
  `idnum_hash` char(32) NOT NULL COMMENT '文件MD5 hash',
  `add_time` bigint(30) NOT NULL DEFAULT 0 COMMENT '添加时间',	
  `update_time` bigint(30) NOT NULL DEFAULT 0 COMMENT '更新时间',	
  PRIMARY KEY (`file_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `folder_user`;
CREATE TABLE `folder_user` (
  `user_id` int(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户id',
  `user_name` varchar(255) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `user_passwd` varchar(255) NOT NULL DEFAULT '' COMMENT '用户passwd',
  `add_time` bigint(30) NOT NULL DEFAULT 0 COMMENT '添加时间',	
  `update_time` bigint(30) NOT NULL DEFAULT 0 COMMENT '更新时间',	
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
