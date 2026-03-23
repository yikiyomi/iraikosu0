class LibraryUI {
private:
    Database db;
    LibrarySystem system;
    int currentAdminId;
    
public:
    void mainMenu();
    void bookManagementMenu();
    void readerManagementMenu();
    void borrowReturnMenu();
    void reportMenu();
    void systemMenu();
    
    // 辅助函数
    void displayBooks(const std::vector<Book>& books);
    void displayReaders(const std::vector<Reader>& readers);
    void displayBorrowRecords(const std::vector<BorrowRecord>& records);
    
    // 输入验证
    int getIntInput(const std::string& prompt);
    std::string getStringInput(const std::string& prompt);
    Date getDateInput(const std::string& prompt);
};