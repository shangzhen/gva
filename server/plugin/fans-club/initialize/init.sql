-- 粉丝团插件初始化SQL
-- 注意：GORM会自动创建表，这里的SQL主要用于添加索引和初始化数据

-- 为粉丝团表添加索引
CREATE INDEX IF NOT EXISTS idx_gva_fans_club_owner_id ON gva_fans_club(owner_id);
CREATE INDEX IF NOT EXISTS idx_gva_fans_club_status ON gva_fans_club(status);
CREATE INDEX IF NOT EXISTS idx_gva_fans_club_deleted_at ON gva_fans_club(deleted_at);

-- 为成员表添加索引
CREATE INDEX IF NOT EXISTS idx_gva_fans_club_member_club_id ON gva_fans_club_member(club_id);
CREATE INDEX IF NOT EXISTS idx_gva_fans_club_member_user_id ON gva_fans_club_member(user_id);
CREATE INDEX IF NOT EXISTS idx_gva_fans_club_member_role ON gva_fans_club_member(role);
CREATE INDEX IF NOT EXISTS idx_gva_fans_club_member_deleted_at ON gva_fans_club_member(deleted_at);
-- 联合索引用于查询用户在某个粉丝团的成员关系
CREATE INDEX IF NOT EXISTS idx_gva_fans_club_member_club_user ON gva_fans_club_member(club_id, user_id);

-- 为动态表添加索引
CREATE INDEX IF NOT EXISTS idx_gva_fans_club_post_club_id ON gva_fans_club_post(club_id);
CREATE INDEX IF NOT EXISTS idx_gva_fans_club_post_user_id ON gva_fans_club_post(user_id);
CREATE INDEX IF NOT EXISTS idx_gva_fans_club_post_deleted_at ON gva_fans_club_post(deleted_at);

-- 插入示例数据（可选）
-- INSERT INTO gva_fans_club (created_at, updated_at, name, description, avatar, owner_id, member_count, level, status)
-- VALUES (NOW(), NOW(), '示例粉丝团', '这是一个示例粉丝团', '', 1, 1, 1, 1);
