func Connect() error {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PWD")
	dbname := os.Getenv("MYSQL_DATABASE")
	instanceConnName := os.Getenv("INSTANCE_CONNECTION_NAME")

	if user == "" || pass == "" || dbname == "" || instanceConnName == "" {
		return fmt.Errorf("環境変数が不足しています")
	}

	dsn := fmt.Sprintf(
		"%s:%s@unix(/cloudsql/%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FTokyo",
		user, pass, instanceConnName, dbname)

	log.Println("🔗 DSN:", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("GORM DB接続失敗: %w", err)
	}

	log.Println("✅ GORM: DB接続成功")
	DB = db
	return nil
}
