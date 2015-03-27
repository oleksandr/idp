BEGIN;

DROP TABLE IF EXISTS user_role;
DROP TABLE IF EXISTS domain_user;
DROP TABLE IF EXISTS role_permission;
DROP TABLE IF EXISTS user_session;
DROP TABLE IF EXISTS domain;
DROP TABLE IF EXISTS `permission`;
DROP TABLE IF EXISTS `role`;
DROP TABLE IF EXISTS `user`;

-- Create syntax for TABLE 'domain'
CREATE TABLE `domain` (
  `domain_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `object_id` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` text,
  `is_enabled` tinyint(1) DEFAULT NULL,
  `created_on` datetime NOT NULL,
  `updated_on` datetime NOT NULL,
  PRIMARY KEY (`domain_id`),
  UNIQUE KEY `object_id` (`object_id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Create syntax for TABLE 'user'
CREATE TABLE `user` (
  `user_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `object_id` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `passwd` varchar(255) NOT NULL,
  `is_enabled` tinyint(1) DEFAULT NULL,
  `created_on` datetime NOT NULL,
  `updated_on` datetime NOT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `object_id` (`object_id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Create syntax for TABLE 'permission'
CREATE TABLE `permission` (
  `permission_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` text,
  `evaluation_rule` text,
  `is_enabled` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`permission_id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Create syntax for TABLE 'role'
CREATE TABLE `role` (
  `role_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` text,
  `is_enabled` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`role_id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Create syntax for TABLE 'user_session'
CREATE TABLE `user_session` (
  `user_session_id` varchar(255) NOT NULL,
  `domain_id` int(11) unsigned NOT NULL,
  `user_id` int(11) unsigned NOT NULL,
  `user_agent` varchar(300) DEFAULT NULL,
  `remote_addr` varchar(20) DEFAULT NULL,
  `created_on` datetime NOT NULL,
  `updated_on` datetime NOT NULL,
  `expires_on` datetime NOT NULL,
  PRIMARY KEY (`user_session_id`),
  KEY `domain_id` (`domain_id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `user_session_ibfk_1` FOREIGN KEY (`domain_id`) REFERENCES `domain` (`domain_id`),
  CONSTRAINT `user_session_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Create syntax for TABLE 'domain_user'
CREATE TABLE `domain_user` (
  `domain_id` int(11) unsigned NOT NULL,
  `user_id` int(11) unsigned NOT NULL,
  PRIMARY KEY (`domain_id`,`user_id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `domain_user_ibfk_1` FOREIGN KEY (`domain_id`) REFERENCES `domain` (`domain_id`),
  CONSTRAINT `domain_user_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Create syntax for TABLE 'role_permission'
CREATE TABLE `role_permission` (
  `role_id` int(11) unsigned NOT NULL,
  `permission_id` int(11) unsigned NOT NULL,
  PRIMARY KEY (`role_id`,`permission_id`),
  KEY `permission_id` (`permission_id`),
  CONSTRAINT `role_permission_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `role` (`role_id`),
  CONSTRAINT `role_permission_ibfk_2` FOREIGN KEY (`permission_id`) REFERENCES `permission` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Create syntax for TABLE 'user_role'
CREATE TABLE `user_role` (
  `user_id` int(11) unsigned NOT NULL,
  `role_id` int(11) unsigned NOT NULL,
  PRIMARY KEY (`user_id`,`role_id`),
  KEY `role_id` (`role_id`),
  CONSTRAINT `user_role_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`),
  CONSTRAINT `user_role_ibfk_2` FOREIGN KEY (`role_id`) REFERENCES `role` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

COMMIT;