[1mdiff --git a/internal/repository/interfaces.go b/internal/repository/interfaces.go[m
[1mindex 5f4f51e..06f96dc 100644[m
[1m--- a/internal/repository/interfaces.go[m
[1m+++ b/internal/repository/interfaces.go[m
[36m@@ -14,7 +14,7 @@[m [mtype ArticleRepository interface {[m
 }[m
 [m
 type UserRepository interface {[m
[31m-	GetByID(id string) (models.User, error)[m
[32m+[m	[32mGetByUserName(username string) (models.User, error)[m
 	Create(data dto.CreateUserDTO) (string, error)[m
 	Update(data dto.UpdateUserDTO) error[m
 	Delete(id string) error[m
[1mdiff --git a/internal/repository/sqlite/user_sqlite.go b/internal/repository/sqlite/user_sqlite.go[m
[1mindex 26dda61..d951def 100644[m
[1m--- a/internal/repository/sqlite/user_sqlite.go[m
[1m+++ b/internal/repository/sqlite/user_sqlite.go[m
[36m@@ -52,8 +52,35 @@[m [mfunc NewUserSqliteRepo(db *sql.DB, clock timeutil.Clock) (UserSqliteRepo, *sql.D[m
 	return r, db[m
 }[m
 [m
[31m-func (r *UserSqliteRepo) GetByID(id string) (models.User, error) {[m
[31m-	return models.User{}, errors.New("Not Implemented")[m
[32m+[m[32mfunc (r *UserSqliteRepo) GetByUserName(userName string) (models.User, error) {[m
[32m+[m	[32msql := sq.Select([m
[32m+[m		[32m"id",[m
[32m+[m		[32m"user_name",[m
[32m+[m		[32m"first_name",[m
[32m+[m		[32m"last_name",[m
[32m+[m		[32m"email",[m
[32m+[m		[32m"hashed_password",[m
[32m+[m		[32m"is_active",[m
[32m+[m		[32m"created_at",[m
[32m+[m		[32m"updated_at",[m
[32m+[m	[32m).From( "users").Where(sq.Eq{"user_name": userName})[m
[32m+[m	[32mrow := sql.RunWith(r.db).QueryRow()[m
[32m+[m	[32muser := models.User{}[m
[32m+[m	[32merr := row.Scan([m
[32m+[m		[32m&user.ID,[m[41m [m
[32m+[m		[32m&user.UserName,[m[41m [m
[32m+[m		[32m&user.FirstName,[m[41m [m
[32m+[m		[32m&user.LastName,[m[41m [m
[32m+[m		[32m&user.Email,[m
[32m+[m		[32m&user.HashedPassword,[m
[32m+[m		[32m&user.IsActive,[m
[32m+[m		[32m&user.CreatedAt,[m
[32m+[m		[32m&user.UpdatedAt,[m
[32m+[m	[32m)[m
[32m+[m	[32mif err != nil {[m
[32m+[m		[32mreturn user, err[m
[32m+[m	[32m}[m
[32m+[m	[32mreturn user, nil[m
 }[m
 [m
 [m
[1mdiff --git a/internal/repository/sqlite/user_sqlite_test.go b/internal/repository/sqlite/user_sqlite_test.go[m
[1mindex d96cdf3..08cf119 100644[m
[1m--- a/internal/repository/sqlite/user_sqlite_test.go[m
[1m+++ b/internal/repository/sqlite/user_sqlite_test.go[m
[36m@@ -18,6 +18,50 @@[m [mfunc NewUserTestDB(t *testing.T) *sql.DB {[m
 	return db[m
 }[m
 [m
[32m+[m[32mfunc TestGetUserByUserName(t *testing.T){[m
[32m+[m	[32mur, db := repo.NewUserSqliteRepo(repo.NewTestDB(), testutil.Fc0)[m
[32m+[m	[32mq := `INSERT[m
[32m+[m	[32mINTO users ([m
[32m+[m		[32mid,[m
[32m+[m		[32muser_name,[m
[32m+[m		[32mfirst_name,[m
[32m+[m		[32mlast_name,[m
[32m+[m		[32memail,[m
[32m+[m		[32mhashed_password,[m
[32m+[m		[32mis_active,[m
[32m+[m		[32mcreated_at,[m
[32m+[m		[32mupdated_at[m
[32m+[m	[32m)[m
[32m+[m	[32mVALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)[m
[32m+[m	[32m`[m
[32m+[m	[32mid := "id"[m
[32m+[m	[32muserName := "tyler.durden"[m
[32m+[m	[32mfirstName := "first name"[m
[32m+[m	[32mlastName := "last name"[m
[32m+[m	[32memail := "email@test.com"[m
[32m+[m	[32mhashedPassword := "random-hashed-value"[m
[32m+[m	[32misActive := true[m
[32m+[m	[32mcreatedAt:= testutil.Fc0.FixedTime.UTC().Unix()[m
[32m+[m	[32mupdatedAt := testutil.Fc0.FixedTime.UTC().Unix()[m
[32m+[m	[32m_, err := db.Exec([m
[32m+[m		[32mq, id, userName, firstName, lastName,[m
[32m+[m		[32memail, hashedPassword, isActive,[m
[32m+[m		[32mcreatedAt, updatedAt,[m
[32m+[m	[32m)[m
[32m+[m	[32mif err != nil{[m
[32m+[m		[32mt.Fatalf("could not insert user for test due to error: %v", err)[m
[32m+[m	[32m}[m
[32m+[m
[32m+[m	[32muser, err := ur.GetByUserName(userName)[m
[32m+[m	[32mif err != nil {[m
[32m+[m		[32mt.Fatalf("could not get user by user name due to error: %v", err)[m
[32m+[m	[32m}[m
[32m+[m
[32m+[m	[32mif user.UserName != userName {[m
[32m+[m		[32mt.Errorf("want %s got %s", userName, user.UserName)[m
[32m+[m	[32m}[m
[32m+[m[32m}[m
[32m+[m
 func TestCreateUser(t *testing.T) {[m
 	var id string[m
 	var err error[m
