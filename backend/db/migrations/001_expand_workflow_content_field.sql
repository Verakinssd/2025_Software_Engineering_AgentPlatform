-- 扩大工作流内容字段的存储容量
-- 将workflows表的content字段从TEXT类型改为LONGTEXT类型
-- 这样可以存储更大的工作流数据（从65KB增加到4GB）

ALTER TABLE workflows MODIFY COLUMN content LONGTEXT NOT NULL; 