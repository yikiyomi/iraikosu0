class Books{
    private:
    int bookId;
    std::string title;
    std::string isbn;
    std::string author;
    std::string publisher;
    std::string publisher_data;
    std::string category;
    double price;
    int total_copies;
    int available_copies;
    std::string location;
    std::string status;
    std::string created_at;

    public:
    Book();
    Book(const std::string& isbn, const std::string& title, 
         const std::string& author);
    
    // CRUD 操作
    bool addToDatabase(Database& db);
    bool updateInDatabase(Database& db);
    bool deleteFromDatabase(Database& db);
    static Book* fetchFromDatabase(Database& db, int bookId);
    
    // 查询方法
    static std::vector<Book> searchByTitle(Database& db, 
                                          const std::string& keyword);
    static std::vector<Book> searchByAuthor(Database& db, 
                                           const std::string& author);
    
    // 借阅相关
    bool isAvailable() const;
    int getAvailableCopies(Database& db);
    
    // Getter/Setter
    std::string getTitle() const { return title; }
    std::string getAuthor() const { return author; }
    // ...
};