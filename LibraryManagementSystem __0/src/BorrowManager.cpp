#include "BorrowManager.hpp"
#include <ctime>
#include <iomanip>
#include <sstream>

bool BorrowRecord::saveToDatabase(Database&db){
    std::string sql="insert into BorrowRecord(reader_id,book_id,borrow_date,due_date,status) VALUES (" +
        std::to_string(readerId) + ", " + std::to_string(bookId) + ", '" +
        borrowDate + "', '" + dueDate + "', '借出')";
    return db.executeUpdate(sql) > 0;
}
bool BorrowRecord::updateInDatabase(Database&db){
    std::string sql = "UPDATE BorrowRecords SET actual_return_date='" + returnDate +
        "', renew_count=" + std::to_string(renewCount) +
        ", overdue_days=" + std::to_string(overdueDays) +
        ", fine_amount=" + std::to_string(fineAmount) +
        ", status='" + status + "' WHERE record_id=" + std::to_string(recordId);
    return db.executeUpdate(sql) > 0;
} 
std::unique_ptr<BorrowRecord> BorrowRecord::fetchFromDatabase(Database& db, int recordId) {
    std::string sql = "SELECT * FROM BorrowRecords WHERE record_id=" + std::to_string(recordId);
    auto res = db.executeQuery(sql);
    if (res && res->next()) {
        auto rec = std::make_unique<BorrowRecord>();
        rec->recordId = res->getInt("record_id");
        rec->readerId = res->getInt("reader_id");
        rec->bookId = res->getInt("book_id");
        rec->borrowDate = res->getString("borrow_date");
        rec->dueDate = res->getString("due_date");
        rec->returnDate = res->getString("actual_return_date");
        rec->renewCount = res->getInt("renew_count");
        rec->overdueDays = res->getInt("overdue_days");
        rec->fineAmount = res->getDouble("fine_amount");
        rec->status = res->getString("status");
        return rec;
    }
    return nullptr;
}
std::string currentDateTime() {
    time_t now = time(0);
    struct tm tstruct;
    char buf[80];
    tstruct = *localtime(&now);
    strftime(buf, sizeof(buf), "%Y-%m-%d %H:%M:%S", &tstruct);
    return buf;
}
std::string currentDate() {
    time_t now = time(0);
    struct tm tstruct;
    char buf[80];
    tstruct = *localtime(&now);
    strftime(buf, sizeof(buf), "%Y-%m-%d", &tstruct);
    return buf;
}
bool BorrowManager::borrowBook(int readerId, const std::string& isbn) {
    // 1. 获取读者信息
    auto reader = Reader::fetchFromDatabase(db, readerId);
    if (!reader) {
        std::cerr << "读者不存在！" << std::endl;
        return false;
    }
    if (!reader->canBorrowMore()) {
        std::cerr << "读者不能借更多书（已达上限或状态不正常）" << std::endl;
        return false;
    }
    if (reader->hasOverdueBooks(db)) {
        std::cerr << "读者有逾期图书未还！" << std::endl;
        return false;
    }
     std::string sql = "SELECT * FROM Books WHERE isbn='" + isbn + "'";
    auto res = db.executeQuery(sql);
    if (!res || !res->next()) {
        std::cerr << "图书不存在！" << std::endl;
        return false;
    }
    Book book;
    book.setBookId(res->getInt("book_id"));
    book.setTitle(res->getString("title"));
    book.setAvailableCopies(res->getInt("available_copies"));
    if (!book.isAvailable()) {
        std::cerr << "图书不可借！" << std::endl;
        return false;
    }

    // 3. 开始事务
    db.beginTransaction();
    try {
        // 更新图书库存
        book.decreaseAvailableCopies();
        std::string updateBook = "UPDATE Books SET available_copies=" +
            std::to_string(book.getAvailableCopies()) + " WHERE book_id=" +
            std::to_string(book.getBookId());
        if (db.executeUpdate(updateBook) <= 0) throw std::runtime_error("更新图书失败");

        // 更新读者借书数量
        reader->increaseBorrowedCount();
        std::string updateReader = "UPDATE Readers SET current_borrowed=" +
            std::to_string(reader->getCurrentBorrowed()) + " WHERE reader_id=" +
            std::to_string(readerId);
        if (db.executeUpdate(updateReader) <= 0) throw std::runtime_error("更新读者失败");

        // 创建借阅记录
        std::string borrowDate = currentDateTime();
        // 默认借期30天
        std::string dueDate = "DATE_ADD(CURDATE(), INTERVAL 30 DAY)"; // 需要在SQL中计算
        std::string insertRecord = "INSERT INTO BorrowRecords (reader_id, book_id, borrow_date, due_date, status) VALUES (" +
            std::to_string(readerId) + ", " + std::to_string(book.getBookId()) + ", '" +
            borrowDate + "', DATE_ADD(CURDATE(), INTERVAL 30 DAY), '借出')";
        if (db.executeUpdate(insertRecord) <= 0) throw std::runtime_error("插入借阅记录失败");

        db.commit();
        std::cout << "借书成功！" << std::endl;
        return true;
    } catch (const std::exception& e) {
        db.rollback();
        std::cerr << "借书失败，已回滚: " << e.what() << std::endl;
        return false;
    }
}

ReturnResult BorrowManager::returnBook(int recordId, const std::string& condition) {
    ReturnResult result;
    auto record = BorrowRecord::fetchFromDatabase(db, recordId);
    if (!record) {
        std::cerr << "借阅记录不存在！" << std::endl;
        return result;
    }

    // 获取图书信息
    auto book = Book::fetchFromDatabase(db, record->bookId);
    if (!book) {
        std::cerr << "图书信息不存在！" << std::endl;
        return result;
    }

    // 获取读者信息
    auto reader = Reader::fetchFromDatabase(db, record->readerId);
    if (!reader) {
        std::cerr << "读者信息不存在！" << std::endl;
        return result;
    }

    db.beginTransaction();
    try {
        std::string returnDate = currentDateTime();
        // 计算逾期天数（简化：假设due_date格式为YYYY-MM-DD）
        // 实际应该解析字符串比较，这里用SQL计算
        std::string calcOverdue = "SELECT DATEDIFF('" + returnDate.substr(0,10) + "', due_date) AS days FROM BorrowRecords WHERE record_id=" + std::to_string(recordId);
        auto res = db.executeQuery(calcOverdue);
        int overdueDays = 0;
        if (res && res->next()) {
            overdueDays = res->getInt("days");
            if (overdueDays < 0) overdueDays = 0;
        }
        result.overdueDays = overdueDays;

        // 计算罚款：假设每天0.1元，上限书价
        double fine = 0.0;
        if (overdueDays > 0) {
            fine = overdueDays * 0.1;
            if (fine > book->getPrice()) fine = book->getPrice();
        }
        // 根据图书状况增加罚款
        if (condition == "损坏") {
            fine += book->getPrice() * 0.2; // 20%损坏费
        } else if (condition == "丢失") {
            fine += book->getPrice(); // 全价赔偿
        }
        result.fineAmount = fine;

        // 更新借阅记录
        std::string updateRecord = "UPDATE BorrowRecords SET actual_return_date='" + returnDate +
            "', overdue_days=" + std::to_string(overdueDays) +
            ", fine_amount=" + std::to_string(fine) +
            ", status='已还' WHERE record_id=" + std::to_string(recordId);
        if (db.executeUpdate(updateRecord) <= 0) throw std::runtime_error("更新借阅记录失败");

        // 更新图书库存
        if (condition == "丢失") {
            // 丢失：总册数减少
            book->setTotalCopies(book->getTotalCopies() - 1);
            // 可借册数不变（因为已经借出）
        } else {
            book->increaseAvailableCopies();
        }
        book->updateInDatabase(db);

        // 更新读者借书数量
        reader->decreaseBorrowedCount();
        // 加上罚款
        if (fine > 0) {
            reader->addPenalty(fine);
        }
        reader->updateInfo(db);

        db.commit();
        std::cout << "还书成功！罚款：" << fine << " 元" << std::endl;
    } catch (const std::exception& e) {
        db.rollback();
        std::cerr << "还书失败，已回滚: " << e.what() << std::endl;
    }
    return result;
}

bool BorrowManager::renewBook(int recordId) {
    auto record = BorrowRecord::fetchFromDatabase(db, recordId);
    if (!record || record->status != "借出") {
        std::cerr << "记录无效或已归还！" << std::endl;
        return false;
    }
    // 检查续借次数限制（最多2次）
    if (record->renewCount >= 2) {
        std::cerr << "已达到最大续借次数！" << std::endl;
        return false;
    }
    // 检查是否逾期（逾期不能续借）
    std::string checkOverdue = "SELECT due_date < CURDATE() AS overdue FROM BorrowRecords WHERE record_id=" + std::to_string(recordId);
    auto res = db.executeQuery(checkOverdue);
    if (res && res->next() && res->getBoolean("overdue")) {
        std::cerr << "图书已逾期，不能续借！" << std::endl;
        return false;
    }

    // 续借：延长应还日期30天
    std::string update = "UPDATE BorrowRecords SET due_date = DATE_ADD(due_date, INTERVAL 30 DAY), renew_count = renew_count + 1 WHERE record_id=" + std::to_string(recordId);
    return db.executeUpdate(update) > 0;
}

double BorrowManager::calculateOverdueFine(int overdueDays, double bookPrice) {
    double fine = overdueDays * 0.1;
    if (fine > bookPrice) fine = bookPrice;
    return fine;
}

void BorrowManager::sendDueDateReminders() {
    // 查询3天内到期的借阅
    std::string sql = "SELECT r.name, r.phone, b.title, br.due_date FROM BorrowRecords br "
                      "JOIN Readers r ON br.reader_id = r.reader_id "
                      "JOIN Books b ON br.book_id = b.book_id "
                      "WHERE br.status='借出' AND br.due_date BETWEEN CURDATE() AND DATE_ADD(CURDATE(), INTERVAL 3 DAY)";
    auto res = db.executeQuery(sql);
    if (res) {
        while (res->next()) {
            std::cout << "提醒：" << res->getString("name") << "，您借阅的《"
                      << res->getString("title") << "》将于 "
                      << res->getString("due_date") << " 到期，请及时归还。" << std::endl;
        }
    }
}

void BorrowManager::sendOverdueNotices() {
    std::string sql = "SELECT r.name, r.phone, b.title, br.due_date, DATEDIFF(CURDATE(), br.due_date) as overdue FROM BorrowRecords br "
                      "JOIN Readers r ON br.reader_id = r.reader_id "
                      "JOIN Books b ON br.book_id = b.book_id "
                      "WHERE br.status='借出' AND br.due_date < CURDATE()";
    auto res = db.executeQuery(sql);
    if (res) {
        while (res->next()) {
            std::cout << "逾期通知：" << res->getString("name") << "，您借阅的《"
                      << res->getString("title") << "》已逾期 "
                      << res->getInt("overdue") << " 天，请尽快归还并缴纳罚款。" << std::endl;
        }
    }
}