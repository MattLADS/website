func InitializeForumDB() {
	var err error
	forumDB, err = gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to forum database:", err)
	}

	// Enable foreign key constraint
	forumDB.Exec("PRAGMA foreign_keys = ON;")

	// move the User, Topic, Comment, and Assignment tables.
	forumDB.AutoMigrate(&User{}, &Topic{}, &Comment{}, &Assignment{})
}
