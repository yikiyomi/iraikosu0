class BorrowManager {
private:
    Database& db;
    
public:
    BorrowManager(Database& database) : db(database) {}
    
    // 借书流程
    BorrowRecord borrowBook(int readerId, int bookId);
    
    // 还书流程
    ReturnResult returnBook(int recordId, BookCondition condition);
    
    // 续借
    bool renewBook(int recordId);
    
    // 计算逾期罚款
    double calculateOverdueFine(int recordId);
    
    // 生成借书证
    void generateBorrowReceipt(int recordId);
    
    // 发送提醒
    void sendDueDateReminders();
    void sendOverdueNotices();
};