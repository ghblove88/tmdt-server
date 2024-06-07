/*
 Navicat Premium Data Transfer

 Source Server         : ecds
 Source Server Type    : MySQL
 Source Server Version : 80036 (8.0.36-0ubuntu0.22.04.1)
 Source Host           : 120.79.34.12:3306
 Source Schema         : tmdt

 Target Server Type    : MySQL
 Target Server Version : 80036 (8.0.36-0ubuntu0.22.04.1)
 File Encoding         : 65001

 Date: 07/06/2024 12:54:26
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for bed
-- ----------------------------
DROP TABLE IF EXISTS `bed`;
CREATE TABLE `bed` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `bed_code` varchar(191) COLLATE utf8mb4_general_ci NOT NULL,
  `bed_info` longtext COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_bed_bed_code` (`bed_code`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of bed
-- ----------------------------
BEGIN;
INSERT INTO `bed` (`id`, `bed_code`, `bed_info`) VALUES (1, 'A01', '检查室');
INSERT INTO `bed` (`id`, `bed_code`, `bed_info`) VALUES (2, 'A02', '检查室');
INSERT INTO `bed` (`id`, `bed_code`, `bed_info`) VALUES (3, 'I01', 'ICU');
INSERT INTO `bed` (`id`, `bed_code`, `bed_info`) VALUES (6, 'I02', 'ICU');
COMMIT;

-- ----------------------------
-- Table structure for device_info
-- ----------------------------
DROP TABLE IF EXISTS `device_info`;
CREATE TABLE `device_info` (
  `id` int NOT NULL AUTO_INCREMENT,
  `device_code` varchar(191) COLLATE utf8mb4_general_ci NOT NULL,
  `device_sequence` varchar(191) COLLATE utf8mb4_general_ci NOT NULL,
  `device_name` longtext COLLATE utf8mb4_general_ci NOT NULL,
  `device_info` longtext COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_device_info_device_code` (`device_code`),
  UNIQUE KEY `uni_device_info_device_sequence` (`device_sequence`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='CREATE TABLE device_info (\n    id INT AUTO_INCREMENT PRIMARY KEY,\n    device_code VARCHAR(255) NOT NULL UNIQUE,\n    device_sequence VARCHAR(255) NOT NULL UNIQUE,\n    device_name VARCHAR(255) NOT NULL,\n    device_info TEXT NOT NULL\n);\n\n	•	id：自增长的主键字段，用于唯一标识每条记录。\n	•	device_code：设备编号，设为 UNIQUE，确保值不能重复。\n	•	device_sequence：设备序号，设为 UNIQUE，确保值不能重复。\n	•	device_name：设备名称。\n	•	device_info：设备信息，存储较长的文本信息。\n';

-- ----------------------------
-- Records of device_info
-- ----------------------------
BEGIN;
INSERT INTO `device_info` (`id`, `device_code`, `device_sequence`, `device_name`, `device_info`) VALUES (1, '1001', '1', '体温检测1', '设备信息1');
INSERT INTO `device_info` (`id`, `device_code`, `device_sequence`, `device_name`, `device_info`) VALUES (2, '1002', '2', '体温检测2', '设备信息2');
INSERT INTO `device_info` (`id`, `device_code`, `device_sequence`, `device_name`, `device_info`) VALUES (3, '1003', '3', '体温监测3', '设备信息3');
INSERT INTO `device_info` (`id`, `device_code`, `device_sequence`, `device_name`, `device_info`) VALUES (4, '1004', '4', '体温检测4', '设备信息4');
INSERT INTO `device_info` (`id`, `device_code`, `device_sequence`, `device_name`, `device_info`) VALUES (5, '1005', '5', '体温检测5', '设备信息5');
INSERT INTO `device_info` (`id`, `device_code`, `device_sequence`, `device_name`, `device_info`) VALUES (9, '1006', '6', '体温监测6', '设备信息6');
COMMIT;

-- ----------------------------
-- Table structure for operator_info
-- ----------------------------
DROP TABLE IF EXISTS `operator_info`;
CREATE TABLE `operator_info` (
  `id` int NOT NULL AUTO_INCREMENT,
  `number` longtext,
  `name` longtext,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb3;

-- ----------------------------
-- Records of operator_info
-- ----------------------------
BEGIN;
INSERT INTO `operator_info` (`id`, `number`, `name`) VALUES (6, '0009666971', '赵雅芝');
INSERT INTO `operator_info` (`id`, `number`, `name`) VALUES (7, '0009668063', '翁美玲');
INSERT INTO `operator_info` (`id`, `number`, `name`) VALUES (8, '0009668425', '钟楚红');
INSERT INTO `operator_info` (`id`, `number`, `name`) VALUES (9, '0014393544', '关之琳');
INSERT INTO `operator_info` (`id`, `number`, `name`) VALUES (11, '0014381923', '张曼玉');
INSERT INTO `operator_info` (`id`, `number`, `name`) VALUES (12, '0014389991', '李赛凤');
INSERT INTO `operator_info` (`id`, `number`, `name`) VALUES (13, '0009929166', '曾华倩');
INSERT INTO `operator_info` (`id`, `number`, `name`) VALUES (15, '0005526995', '李丽珍');
INSERT INTO `operator_info` (`id`, `number`, `name`) VALUES (16, '0007388306', '周慧敏');
INSERT INTO `operator_info` (`id`, `number`, `name`) VALUES (18, '0009440932', '邱淑贞');
INSERT INTO `operator_info` (`id`, `number`, `name`) VALUES (20, '0010371733', '莫文蔚');
INSERT INTO `operator_info` (`id`, `number`, `name`) VALUES (21, '0008351013', '朱茵');
INSERT INTO `operator_info` (`id`, `number`, `name`) VALUES (22, '0008319245', '黎姿');
INSERT INTO `operator_info` (`id`, `number`, `name`) VALUES (23, '0009667336', '袁咏仪');
COMMIT;

-- ----------------------------
-- Table structure for temperature_record
-- ----------------------------
DROP TABLE IF EXISTS `temperature_record`;
CREATE TABLE `temperature_record` (
  `id` int NOT NULL AUTO_INCREMENT,
  `device_code` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `bed_code` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `operator` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `patient` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `temperature1` float NOT NULL,
  `temperature2` float DEFAULT NULL,
  `temperature3` float DEFAULT NULL,
  `record_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='CREATE TABLE temperature_records (\n    id INT AUTO_INCREMENT PRIMARY KEY,\n    device_code VARCHAR(255) NOT NULL,\n    bed_code VARCHAR(255) NOT NULL,\n    operator VARCHAR(255) NOT NULL,\n    patient VARCHAR(255) NOT NULL,\n    temperature1 FLOAT NOT NULL,\n    temperature2 FLOAT,\n    temperature3 FLOAT,\n    record_time DATETIME NOT NULL\n);\n此表包含以下字段：\n\n	•	id: 自增长的主键字段，用于唯一标识每条记录。\n	•	device_code: 设备编号。\n	•	bed_code: 病床号。\n	•	operator: 操作员。\n	•	patient: 病人。\n	•	temperature1: 体温1。\n	•	temperature2: 体温2。\n	•	temperature3: 体温3。\n	•	record_time: 记录时间。';

-- ----------------------------
-- Records of temperature_record
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `password` char(36) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb3;

-- ----------------------------
-- Records of user
-- ----------------------------
BEGIN;
INSERT INTO `user` (`id`, `username`, `password`) VALUES (1, 'admin', '0192023a7bbd73250516f069df18b500');
INSERT INTO `user` (`id`, `username`, `password`) VALUES (5, 'user', '0192023a7bbd73250516f069df18b500');
COMMIT;

-- ----------------------------
-- Procedure structure for NewProc
-- ----------------------------
DROP PROCEDURE IF EXISTS `NewProc`;
delimiter ;;
CREATE PROCEDURE `NewProc`(OUT `proc3` int)
BEGIN
	declare _number varchar(20);
declare state int default false; 
-- BETWEEN '2020-03-01 00:00:00' AND '2020-04-08 00:00:00'; 要更改记录的时间范围
declare ant_cur1 CURSOR for SELECT number FROM ant WHERE begin_time BETWEEN '2021-06-01 00:00:00' AND '2021-07-15 00:00:00';

open ant_cur1;

cur_loop:loop 
        fetch ant_cur1 into _number;
    --  id < 838212   838212 为ant_step 中 2020-03-01 这天第一个记录 对应洗消记录 的 ant_step 中正确纪律的 最小id
    DELETE FROM ant_step  WHERE number = _number and  id < 60305;

        if state then
            leave cur_loop; 
        end if; 
end loop; 
close ant_cur1;

END
;;
delimiter ;

-- ----------------------------
-- Procedure structure for p_get_endoscope_lastrecord
-- ----------------------------
DROP PROCEDURE IF EXISTS `p_get_endoscope_lastrecord`;
delimiter ;;
CREATE PROCEDURE `p_get_endoscope_lastrecord`(IN Endoscope_Number varchar(15))
BEGIN
	DECLARE time1 dateTime;
	DECLARE time2 dateTime;
	DECLARE time3 dateTime;
	
	DECLARE id1 INT;
	DECLARE id2 INT;
	DECLARE id3 INT;
	
	SELECT begin_time,id INTO time1,id1 FROM clear1.ant WHERE ant.endoscope_number = Endoscope_Number ORDER BY begin_time DESC LIMIT 0,1;
	SELECT begin_time,id INTO time2,id2 FROM clear2.ant WHERE ant.endoscope_number = Endoscope_Number ORDER BY begin_time DESC LIMIT 0,1;
	SELECT begin_time,id INTO time3,id3 FROM clear3.ant WHERE ant.endoscope_number = Endoscope_Number ORDER BY begin_time DESC LIMIT 0,1;
	
	IF ISNULL(time1) THEN
		SET time1 = MAKEDATE(2000,1);
	END IF;
	IF ISNULL(time2) THEN
		SET time2 = MAKEDATE(2000,1);
	END IF;
	IF ISNULL(time3) THEN
		SET time3 = MAKEDATE(2000,1);
	END IF;
	
	/*SELECT time1,time2,time3;
	SELECT id1,id2,id3;*/
	
	IF time1 > time2 and time1 > time3 THEN
		SELECT ant.id,ant.number,ant.endoscope_number,ant.operator,ant.begin_time,ant.end_time,ant.total_cost_time,ant.endoscope_type,ant.endoscope_info,ant_step.step,ant_step.cost_time 
		FROM clear1.ant LEFT JOIN clear1.ant_step on ant.number = ant_step.number WHERE ant.id = id1;
	ELSEIF  time2 > time1 and time2 > time3 THEN
		SELECT ant.id,ant.number,ant.endoscope_number,ant.operator,ant.begin_time,ant.end_time,ant.total_cost_time,ant.endoscope_type,ant.endoscope_info,ant_step.step,ant_step.cost_time 
		FROM clear2.ant LEFT JOIN clear2.ant_step on ant.number = ant_step.number WHERE ant.id = id2;
	ELSEIF time3 > time1 and time3 > time2 THEN
		SELECT ant.id,ant.number,ant.endoscope_number,ant.operator,ant.begin_time,ant.end_time,ant.total_cost_time,ant.endoscope_type,ant.endoscope_info,ant_step.step,ant_step.cost_time 
		FROM clear3.ant LEFT JOIN clear3.ant_step on ant.number = ant_step.number WHERE ant.id = id3;
	ELSE 
		SELECT NULL;
	END IF;
END
;;
delimiter ;

SET FOREIGN_KEY_CHECKS = 1;
