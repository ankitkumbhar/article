package models

type article struct {
	app *Application
}

// ArticleStore holds all method
type ArticleStore interface {
	Store(article *Article) (int64, error)
	GetByID(articleID int) (*Article, error)
	GetAll() ([]*Article, error)
}

// Article holds article fields
type Article struct {
	ID      int    `db:"id"`
	Title   string `db:"title"`
	Content string `db:"content"`
	Author  string `db:"author"`
}

// Store used to store article in database
func (a *article) Store(article *Article) (lastInsertedID int64, err error) {
	// prepare query to insert record
	query := `INSERT INTO article (title, content, author) 
		VALUES(?, ?, ?)`

	// execute query
	res, err := a.app.db.Exec(query, article.Title, article.Content, article.Author)
	if err != nil {
		return lastInsertedID, err
	}

	// get last inserted record ID
	lastInsertedID, err = res.LastInsertId()

	return lastInsertedID, err
}

// GetByID fetches article by articleID
func (a *article) GetByID(articleID int) (*Article, error) {
	query := `SELECT id, title, content, author FROM article  
		WHERE id=?`

	row, err := a.app.db.Query(query, articleID)
	if err != nil {
		return nil, err
	}

	var article Article

	for row.Next() {
		err = row.Scan(&article.ID, &article.Title, &article.Content, &article.Author)
		if err != nil {
			return nil, err
		}
	}

	return &article, nil
}

// GetAll fetches all articles
func (a *article) GetAll() ([]*Article, error) {
	query := `SELECT id, title, content, author FROM article`

	row, err := a.app.db.Query(query)
	if err != nil {
		return nil, err
	}

	var articles []*Article

	for row.Next() {
		var article Article

		err = row.Scan(&article.ID, &article.Title, &article.Content, &article.Author)
		if err != nil {
			return nil, err
		}

		articles = append(articles, &article)
	}

	return articles, nil
}
