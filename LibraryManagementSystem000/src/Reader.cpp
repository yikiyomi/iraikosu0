class Reader {
private:
    int readerId;
    std::string cardNumber;
    std::string name;
    std::string gender;
    std::string email;
    std::string phone;
    std::string address;
    std::string rederType;
    int maxBorrowLimit;
    int currentBorrowed;
    std::string registrationDate;
    std::string expiryDate;
    std::string status;
    double penalty;
    std::vector<Book*> borrowedBooks;
    
public:
    Reader();
    Reader(const std::string& cardNo, const std::string& name);
    
    // 注册和管理
    bool registerReader(Database& db);
    bool updateInfo(Database& db);
    
    // 借阅功能
    bool borrowBook(Database& db, Book& book);
    bool returnBook(Database& db, Book& book);
    bool renewBook(Database& db, Book& book);
    
    // 查询功能
    std::vector<Book> getBorrowedBooks(Database& db);
    std::vector<BorrowRecord> getBorrowHistory(Database& db);
    double calculatePenalty(Database& db);
    
    // 验证
    bool canBorrowMore() const;
    bool hasOverdueBooks(Database& db);
};