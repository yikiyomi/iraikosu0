#ifndef BORROWMANAGER_HPP
#define BORROWMANAGER_HPP

#include "Database.hpp"
#include "Book.hpp"
#include "Reader.hpp"
#include <string>
#include <memory>

struct BorrowRecord {
    int recordId;
    int readerId;
    int bookId;
    std::string borrowDate;
    std::string dueDate;
    std::string returnDate;
    int renewCount;
    int overdueDays;
    double fineAmount;
    std::string status;

    bool saveToDatabase(Database& db);
    static std::unique_ptr<BorrowRecord> fetchFromDatabase(Database& db, int recordId);
    bool updateInDatabase(Database& db);
};

struct ReturnResult {
    int overdueDays = 0;
    double fineAmount = 0.0;
};

class BorrowManager {
private:
    Database& db;

public:
    BorrowManager(Database& database) : db(database) {}

    // 借书
    bool borrowBook(int readerId, const std::string& isbn);
    // 还书
    ReturnResult returnBook(int recordId, const std::string& condition = "正常");
    // 续借
    bool renewBook(int recordId);

    // 计算罚款
    double calculateOverdueFine(int overdueDays, double bookPrice = 0);
    // 发送提醒（控制台模拟）
    void sendDueDateReminders();
    void sendOverdueNotices();
};

#endif // BORROWMANAGER_HPP