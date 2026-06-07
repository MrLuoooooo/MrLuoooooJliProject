package repository

import (
	"community-server/internal/db/mysql"

	"gorm.io/gorm"
)

// ============================================
// MySQL 数据访问层
// ============================================

// --- UserRepo ---
type userRepo struct{ db *gorm.DB }

func NewUserRepo(db *gorm.DB) UserRepository               { return &userRepo{db} }
func (r *userRepo) Create(u *mysql.User) error              { return r.db.Create(u).Error }
func (r *userRepo) FindByUsername(username string) (*mysql.User, error) {
	var u mysql.User; err := r.db.Where("username = ?", username).First(&u).Error; return &u, err
}
func (r *userRepo) FindByEmail(email string) (*mysql.User, error) {
	var u mysql.User; err := r.db.Where("email = ?", email).First(&u).Error; return &u, err
}
func (r *userRepo) FindByID(id uint) (*mysql.User, error) {
	var u mysql.User; err := r.db.First(&u, id).Error; return &u, err
}
func (r *userRepo) Update(userID uint, updates map[string]interface{}) error {
	return r.db.Model(&mysql.User{}).Where("id = ?", userID).Updates(updates).Error
}
func (r *userRepo) Delete(userID uint) error {
	return r.db.Delete(&mysql.User{}, userID).Error
}
func (r *userRepo) FindByIDs(ids []uint) ([]mysql.User, error) {
	var users []mysql.User
	err := r.db.Where("id IN ?", ids).Find(&users).Error
	return users, err
}
func (r *userRepo) Search(keyword string, page, pageSize int) ([]mysql.User, int64, error) {
	var users []mysql.User; var total int64
	q := r.db.Model(&mysql.User{}).Where("status = 1 AND (username LIKE ? OR nickname LIKE ? OR bio LIKE ?)", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	q.Count(&total)
	err := q.Offset(offset(page, pageSize)).Limit(pageSize).Order("created_at DESC").Find(&users).Error
	return users, total, err
}
func (r *userRepo) List(page, pageSize int) ([]mysql.User, int64, error) {
	var users []mysql.User; var total int64
	r.db.Model(&mysql.User{}).Count(&total)
	err := r.db.Offset(offset(page, pageSize)).Limit(pageSize).Order("created_at DESC").Find(&users).Error
	return users, total, err
}
func (r *userRepo) Count() (int64, error) {
	var c int64; err := r.db.Model(&mysql.User{}).Count(&c).Error; return c, err
}
func (r *userRepo) CountToday() (int64, error) {
	var c int64; err := r.db.Model(&mysql.User{}).Where("DATE(created_at) = CURDATE()").Count(&c).Error; return c, err
}

// --- PostRepo ---
type postRepo struct{ db *gorm.DB }

func NewPostRepo(db *gorm.DB) PostRepository               { return &postRepo{db} }
func (r *postRepo) Create(p *mysql.Post) error              { return r.db.Create(p).Error }
func (r *postRepo) FindByID(id uint) (*mysql.Post, error) {
	var p mysql.Post; err := r.db.Where("id = ?", id).First(&p).Error; return &p, err
}
func (r *postRepo) Update(postID uint, updates map[string]interface{}) error {
	return r.db.Model(&mysql.Post{}).Where("id = ?", postID).Updates(updates).Error
}
func (r *postRepo) UpdateColumn(postID uint, column string, value interface{}) error {
	return r.db.Model(&mysql.Post{}).Where("id = ?", postID).UpdateColumn(column, value).Error
}
func (r *postRepo) List(req PostListQuery) ([]mysql.Post, int64, error) {
	var posts []mysql.Post; var total int64
	q := r.db.Model(&mysql.Post{})
	if req.StatusExclude > 0 {
		q = q.Where("status != ?", req.StatusExclude)
	}
	if req.CategoryID > 0 {
		q = q.Where("category_id = ?", req.CategoryID)
	}
	if req.Keyword != "" {
		q = q.Where("title LIKE ? OR content LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.UserID > 0 {
		q = q.Where("user_id = ?", req.UserID)
	}
	if len(req.UserIDs) > 0 {
		q = q.Where("user_id IN ?", req.UserIDs)
	}
	q.Count(&total)
	orderBy := "is_top DESC, created_at DESC"
	switch req.Sort {
	case "hot":
		orderBy = "is_top DESC, comment_count DESC, view_count DESC"
	case "essence":
		orderBy = "is_essence DESC, is_top DESC, created_at DESC"
	}
	err := q.Order(orderBy).Offset(offset(req.Page, req.PageSize)).Limit(req.PageSize).Find(&posts).Error
	return posts, total, err
}
func (r *postRepo) Search(keyword string, page, pageSize int) ([]mysql.Post, int64, error) {
	var posts []mysql.Post; var total int64
	q := r.db.Model(&mysql.Post{}).Where("status != 2 AND (title LIKE ? OR content LIKE ? OR summary LIKE ?)", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	q.Count(&total)
	err := q.Order("is_top DESC, created_at DESC").Offset(offset(page, pageSize)).Limit(pageSize).Find(&posts).Error
	return posts, total, err
}
func (r *postRepo) Delete(postID uint) error {
	return r.db.Delete(&mysql.Post{}, postID).Error
}
func (r *postRepo) SoftDelete(postID uint) error {
	return r.db.Model(&mysql.Post{}).Where("id = ?", postID).Update("status", 2).Error
}
func (r *postRepo) Count() (int64, error) {
	var c int64; err := r.db.Model(&mysql.Post{}).Count(&c).Error; return c, err
}
func (r *postRepo) CountToday() (int64, error) {
	var c int64; err := r.db.Model(&mysql.Post{}).Where("DATE(created_at) = CURDATE()").Count(&c).Error; return c, err
}

// --- PostLikeRepo ---
type postLikeRepo struct{ db *gorm.DB }

func NewPostLikeRepo(db *gorm.DB) PostLikeRepository       { return &postLikeRepo{db} }
func (r *postLikeRepo) Create(l *mysql.PostLike) error      { return r.db.Create(l).Error }
func (r *postLikeRepo) Delete(userID, postID uint) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&mysql.PostLike{}).Error
}
func (r *postLikeRepo) Exists(userID, postID uint) (bool, error) {
	var count int64
	err := r.db.Model(&mysql.PostLike{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error
	return count > 0, err
}
func (r *postLikeRepo) FindByUserID(userID uint, page, pageSize int) ([]mysql.PostLike, int64, error) {
	var likes []mysql.PostLike; var total int64
	r.db.Model(&mysql.PostLike{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Offset(offset(page, pageSize)).Limit(pageSize).Find(&likes).Error
	return likes, total, err
}

// --- PostFavoriteRepo ---
type postFavRepo struct{ db *gorm.DB }

func NewPostFavoriteRepo(db *gorm.DB) PostFavoriteRepository { return &postFavRepo{db} }
func (r *postFavRepo) Create(f *mysql.PostFavorite) error   { return r.db.Create(f).Error }
func (r *postFavRepo) Delete(userID, postID uint) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&mysql.PostFavorite{}).Error
}
func (r *postFavRepo) Exists(userID, postID uint) (bool, error) {
	var count int64
	err := r.db.Model(&mysql.PostFavorite{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error
	return count > 0, err
}
func (r *postFavRepo) FindByUserID(userID uint, page, pageSize int) ([]mysql.PostFavorite, int64, error) {
	var favs []mysql.PostFavorite; var total int64
	r.db.Model(&mysql.PostFavorite{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Offset(offset(page, pageSize)).Limit(pageSize).Find(&favs).Error
	return favs, total, err
}

// --- CommentRepo ---
type commentRepo struct{ db *gorm.DB }

func NewCommentRepo(db *gorm.DB) CommentRepository         { return &commentRepo{db} }
func (r *commentRepo) Create(c *mysql.Comment) error        { return r.db.Create(c).Error }
func (r *commentRepo) FindByID(id uint) (*mysql.Comment, error) {
	var c mysql.Comment; err := r.db.First(&c, id).Error; return &c, err
}
func (r *commentRepo) FindRootByPostID(postID uint, page, pageSize int) ([]mysql.Comment, int64, error) {
	var comments []mysql.Comment; var total int64
	r.db.Model(&mysql.Comment{}).Where("post_id = ? AND parent_id = 0 AND status != 2", postID).Count(&total)
	err := r.db.Where("post_id = ? AND parent_id = 0 AND status != 2", postID).Order("created_at DESC").Offset(offset(page, pageSize)).Limit(pageSize).Find(&comments).Error
	return comments, total, err
}
func (r *commentRepo) FindRepliesByPostAndParents(postID uint, parentIDs []uint) ([]mysql.Comment, error) {
	var replies []mysql.Comment
	err := r.db.Where("post_id = ? AND parent_id IN ? AND status != 2", postID, parentIDs).Order("created_at ASC").Find(&replies).Error
	return replies, err
}
func (r *commentRepo) Update(commentID uint, updates map[string]interface{}) error {
	return r.db.Model(&mysql.Comment{}).Where("id = ?", commentID).Updates(updates).Error
}
func (r *commentRepo) UpdateColumn(commentID uint, column string, value interface{}) error {
	return r.db.Model(&mysql.Comment{}).Where("id = ?", commentID).UpdateColumn(column, value).Error
}
func (r *commentRepo) SoftDelete(commentID uint) error {
	return r.db.Model(&mysql.Comment{}).Where("id = ?", commentID).Update("status", 2).Error
}
func (r *commentRepo) SoftDeleteByPostID(postID uint) error {
	return r.db.Model(&mysql.Comment{}).Where("post_id = ?", postID).Update("status", 2).Error
}
func (r *commentRepo) Count() (int64, error) {
	var c int64; err := r.db.Model(&mysql.Comment{}).Count(&c).Error; return c, err
}

// --- CommentLikeRepo ---
type commentLikeRepo struct{ db *gorm.DB }

func NewCommentLikeRepo(db *gorm.DB) CommentLikeRepository  { return &commentLikeRepo{db} }
func (r *commentLikeRepo) Create(l *mysql.CommentLike) error { return r.db.Create(l).Error }
func (r *commentLikeRepo) Delete(userID, commentID uint) error {
	return r.db.Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&mysql.CommentLike{}).Error
}
func (r *commentLikeRepo) Exists(userID, commentID uint) (bool, error) {
	var count int64
	err := r.db.Model(&mysql.CommentLike{}).Where("user_id = ? AND comment_id = ?", userID, commentID).Count(&count).Error
	return count > 0, err
}

// --- TagRepo ---
type tagRepo struct{ db *gorm.DB }

func NewTagRepo(db *gorm.DB) TagRepository                 { return &tagRepo{db} }
func (r *tagRepo) Create(t *mysql.Tag) error                { return r.db.Create(t).Error }
func (r *tagRepo) FindByID(id uint) (*mysql.Tag, error) {
	var t mysql.Tag; err := r.db.First(&t, id).Error; return &t, err
}
func (r *tagRepo) FindByName(name string) (*mysql.Tag, error) {
	var t mysql.Tag; err := r.db.Where("name = ?", name).First(&t).Error; return &t, err
}
func (r *tagRepo) FindByIDs(ids []uint) ([]mysql.Tag, error) {
	var tags []mysql.Tag; err := r.db.Where("id IN ?", ids).Find(&tags).Error; return tags, err
}
func (r *tagRepo) Update(tagID uint, updates map[string]interface{}) error {
	return r.db.Model(&mysql.Tag{}).Where("id = ?", tagID).Updates(updates).Error
}
func (r *tagRepo) Delete(tagID uint) error {
	return r.db.Delete(&mysql.Tag{}, tagID).Error
}
func (r *tagRepo) List(page, pageSize int) ([]mysql.Tag, int64, error) {
	var tags []mysql.Tag; var total int64
	r.db.Model(&mysql.Tag{}).Where("status = 1").Count(&total)
	err := r.db.Where("status = 1").Offset(offset(page, pageSize)).Limit(pageSize).Order("post_count DESC").Find(&tags).Error
	return tags, total, err
}
func (r *tagRepo) IncrementPostCount(tagID uint, delta int) error {
	return r.db.Model(&mysql.Tag{}).Where("id = ?", tagID).UpdateColumn("post_count", gorm.Expr("post_count + ?", delta)).Error
}

// --- PostTagRepo ---
type postTagRepo struct{ db *gorm.DB }

func NewPostTagRepo(db *gorm.DB) PostTagRepository          { return &postTagRepo{db} }
func (r *postTagRepo) Create(pt *mysql.PostTag) error       { return r.db.Create(pt).Error }
func (r *postTagRepo) Delete(postID, tagID uint) error {
	return r.db.Where("post_id = ? AND tag_id = ?", postID, tagID).Delete(&mysql.PostTag{}).Error
}
func (r *postTagRepo) FindByPostID(postID uint) ([]mysql.PostTag, error) {
	var pts []mysql.PostTag; err := r.db.Where("post_id = ?", postID).Find(&pts).Error; return pts, err
}
func (r *postTagRepo) FindByTagID(tagID uint, page, pageSize int) ([]mysql.PostTag, int64, error) {
	var pts []mysql.PostTag; var total int64
	r.db.Model(&mysql.PostTag{}).Where("tag_id = ?", tagID).Count(&total)
	err := r.db.Where("tag_id = ?", tagID).Offset(offset(page, pageSize)).Limit(pageSize).Find(&pts).Error
	return pts, total, err
}
func (r *postTagRepo) DeleteByPostID(postID uint) error {
	return r.db.Where("post_id = ?", postID).Delete(&mysql.PostTag{}).Error
}

// --- CategoryRepo ---
type categoryRepo struct{ db *gorm.DB }

func NewCategoryRepo(db *gorm.DB) CategoryRepository        { return &categoryRepo{db} }
func (r *categoryRepo) Create(c *mysql.Category) error      { return r.db.Create(c).Error }
func (r *categoryRepo) FindByID(id uint) (*mysql.Category, error) {
	var c mysql.Category; err := r.db.First(&c, id).Error; return &c, err
}
func (r *categoryRepo) Update(catID uint, updates map[string]interface{}) error {
	return r.db.Model(&mysql.Category{}).Where("id = ?", catID).Updates(updates).Error
}
func (r *categoryRepo) Delete(catID uint) error {
	return r.db.Delete(&mysql.Category{}, catID).Error
}
func (r *categoryRepo) List() ([]mysql.Category, error) {
	var cats []mysql.Category; err := r.db.Where("status = 1").Order("sort_order ASC").Find(&cats).Error; return cats, err
}
func (r *categoryRepo) IncrementPostCount(catID uint, delta int) error {
	return r.db.Model(&mysql.Category{}).Where("id = ?", catID).UpdateColumn("post_count", gorm.Expr("post_count + ?", delta)).Error
}

// --- FollowRepo ---
type followRepo struct{ db *gorm.DB }

func NewFollowRepo(db *gorm.DB) FollowRepository            { return &followRepo{db} }
func (r *followRepo) Create(f *mysql.UserFollow) error      { return r.db.Create(f).Error }
func (r *followRepo) Delete(userID, followID uint) error {
	return r.db.Where("user_id = ? AND follow_id = ?", userID, followID).Delete(&mysql.UserFollow{}).Error
}
func (r *followRepo) Exists(userID, followID uint) (bool, error) {
	var count int64
	err := r.db.Model(&mysql.UserFollow{}).Where("user_id = ? AND follow_id = ?", userID, followID).Count(&count).Error
	return count > 0, err
}
func (r *followRepo) FindFollowers(userID uint, page, pageSize int) ([]mysql.UserFollow, int64, error) {
	var follows []mysql.UserFollow; var total int64
	r.db.Model(&mysql.UserFollow{}).Where("follow_id = ?", userID).Count(&total)
	err := r.db.Where("follow_id = ?", userID).Order("created_at DESC").Offset(offset(page, pageSize)).Limit(pageSize).Find(&follows).Error
	return follows, total, err
}
func (r *followRepo) FindFollowing(userID uint, page, pageSize int) ([]mysql.UserFollow, int64, error) {
	var follows []mysql.UserFollow; var total int64
	r.db.Model(&mysql.UserFollow{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Offset(offset(page, pageSize)).Limit(pageSize).Find(&follows).Error
	return follows, total, err
}
func (r *followRepo) CountFollowers(userID uint) int64 {
	var count int64; r.db.Model(&mysql.UserFollow{}).Where("follow_id = ?", userID).Count(&count); return count
}
func (r *followRepo) CountFollowing(userID uint) int64 {
	var count int64; r.db.Model(&mysql.UserFollow{}).Where("user_id = ?", userID).Count(&count); return count
}
func (r *followRepo) FindAllFollowing(userID uint) ([]mysql.UserFollow, error) {
	var follows []mysql.UserFollow; err := r.db.Where("user_id = ?", userID).Find(&follows).Error; return follows, err
}

// --- MessageRepo ---
type messageRepo struct{ db *gorm.DB }

func NewMessageRepo(db *gorm.DB) MessageRepository          { return &messageRepo{db} }
func (r *messageRepo) Create(m *mysql.Message) error        { return r.db.Create(m).Error }
func (r *messageRepo) FindByConversation(senderID, receiverID uint, page, pageSize int) ([]mysql.Message, int64, error) {
	var msgs []mysql.Message; var total int64
	q := r.db.Model(&mysql.Message{}).Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", senderID, receiverID, receiverID, senderID)
	q.Count(&total)
	err := q.Order("created_at DESC").Offset(offset(page, pageSize)).Limit(pageSize).Find(&msgs).Error
	return msgs, total, err
}
func (r *messageRepo) FindReceived(userID uint, page, pageSize int) ([]mysql.Message, int64, error) {
	var msgs []mysql.Message; var total int64
	r.db.Model(&mysql.Message{}).Where("receiver_id = ?", userID).Count(&total)
	err := r.db.Where("receiver_id = ?", userID).Order("created_at DESC").Offset(offset(page, pageSize)).Limit(pageSize).Find(&msgs).Error
	return msgs, total, err
}
func (r *messageRepo) CountUnread(userID uint) (int64, error) {
	var count int64; err := r.db.Model(&mysql.Message{}).Where("receiver_id = ? AND is_read = ?", userID, false).Count(&count).Error; return count, err
}
func (r *messageRepo) MarkAsRead(messageID, userID uint) error {
	return r.db.Model(&mysql.Message{}).Where("id = ? AND receiver_id = ?", messageID, userID).Update("is_read", true).Error
}
func (r *messageRepo) FindConversationMessages(userID uint) ([]mysql.Message, error) {
	var msgs []mysql.Message; err := r.db.Where("sender_id = ? OR receiver_id = ?", userID, userID).Order("created_at DESC").Find(&msgs).Error; return msgs, err
}
func (r *messageRepo) CountUnreadBySenders(receiverID uint, senderIDs []uint) (map[uint]int64, error) {
	type row struct { SenderID uint; Count int64 }
	var rows []row
	err := r.db.Model(&mysql.Message{}).Select("sender_id, count(*) as count").Where("receiver_id = ? AND is_read = ? AND sender_id IN ?", receiverID, false, senderIDs).Group("sender_id").Find(&rows).Error
	result := make(map[uint]int64)
	for _, r := range rows { result[r.SenderID] = r.Count }
	return result, err
}

// --- NotificationRepo ---
type notifRepo struct{ db *gorm.DB }

func NewNotificationRepo(db *gorm.DB) NotificationRepository { return &notifRepo{db} }
func (r *notifRepo) Create(n *mysql.Notification) error      { return r.db.Create(n).Error }
func (r *notifRepo) FindByUserID(userID uint, page, pageSize int) ([]mysql.Notification, int64, error) {
	var notifs []mysql.Notification; var total int64
	r.db.Model(&mysql.Notification{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Offset(offset(page, pageSize)).Limit(pageSize).Find(&notifs).Error
	return notifs, total, err
}
func (r *notifRepo) CountUnread(userID uint) (int64, error) {
	var count int64; err := r.db.Model(&mysql.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error; return count, err
}
func (r *notifRepo) MarkRead(notifID, userID uint) error {
	return r.db.Model(&mysql.Notification{}).Where("id = ? AND user_id = ?", notifID, userID).Update("is_read", true).Error
}
func (r *notifRepo) MarkAllRead(userID uint) error {
	return r.db.Model(&mysql.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}

// --- helpers ---
func offset(page, pageSize int) int {
	o := (page - 1) * pageSize
	if o < 0 { return 0 }
	return o
}

// ============================================
// 密码重置
// ============================================
type passwordResetRepo struct{ db *gorm.DB }

func NewPasswordResetRepo(db *gorm.DB) PasswordResetRepository {
	_ = db.AutoMigrate(&mysql.PasswordReset{})
	return &passwordResetRepo{db: db}
}

func (r *passwordResetRepo) Create(reset *mysql.PasswordReset) error {
	return r.db.Create(reset).Error
}

func (r *passwordResetRepo) FindByToken(token string) (*mysql.PasswordReset, error) {
	var reset mysql.PasswordReset
	err := r.db.Where("token = ?", token).First(&reset).Error
	return &reset, err
}

func (r *passwordResetRepo) MarkUsed(id uint) error {
	return r.db.Model(&mysql.PasswordReset{}).Where("id = ?", id).Update("used", true).Error
}
