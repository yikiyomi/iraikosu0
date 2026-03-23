#include "ReportGenerator.hpp"
#include <iomanip>

void ReportGenerator::generateBorrowStatistics(const DateRange& range) {
    std::string sql = "SELECT DATE(borrow_date) as day, COUNT(*) as count FROM BorrowRecords "
                      "WHERE borrow_date BETWEEN '" + range.start + "' AND '" + range.end + "' "
                      "GROUP BY day ORDER BY day";
    auto res = db.executeQuery(sql);
    std::cout << "借阅统计 (" << range.start << " 至 " << range.end << "):" << std::endl;
    if (res) {
        while (res->next()) {
            std::cout << res->getString("day") << ": " << res->getInt("count") << " 次" << std::endl;
        }
    }
}

void ReportGenerator::generatePopularBooksReport(int topN) {
    std::string sql = "SELECT b.title, COUNT(*) as borrow_count FROM BorrowRecords br "
                      "JOIN Books b ON br.book_id = b.book_id "
                      "GROUP BY b.book_id ORDER BY borrow_count DESC LIMIT " + std::to_string(topN);
    auto res = db.executeQuery(sql);
    std::cout << "热门图书 Top " << topN << ":" << std::endl;
    if (res) {
        while (res->next()) {
            std::cout << res->getString("title") << " - 借阅次数: " << res->getInt("borrow_count") << std::endl;
        }
    }
}

void ReportGenerator::generateReaderActivityReport() {
    std::string sql = "SELECT r.name, COUNT(*) as borrow_count FROM BorrowRecords br "
                      "JOIN Readers r ON br.reader_id = r.reader_id "
                      "GROUP BY r.reader_id ORDER BY borrow_count DESC";
    auto res = db.executeQuery(sql);
    std::cout << "读者活跃度排行:" << std::endl;
    if (res) {
        while (res->next()) {
            std::cout << res->getString("name") << " - 借阅次数: " << res->getInt("borrow_count") << std::endl;
        }
    }
}

void ReportGenerator::generateOverdueAnalysis() {
    std::string sql = "SELECT r.name, b.title, br.due_date, DATEDIFF(CURDATE(), br.due_date) as overdue_days "
                      "FROM BorrowRecords br "
                      "JOIN Readers r ON br.reader_id = r.reader_id "
                      "JOIN Books b ON br.book_id = b.book_id "
                      "WHERE br.status='借出' AND br.due_date < CURDATE() "
                      "ORDER BY overdue_days DESC";
    auto res = db.executeQuery(sql);
    std::cout << "逾期分析:" << std::endl;
    if (res) {
        while (res->next()) {
            std::cout << "读者: " << res->getString("name") << "，图书: " << res->getString("title")
                      << "，逾期天数: " << res->getInt("overdue_days") << std::endl;
        }
    }
}

void ReportGenerator::generateInventoryReport() {
    std::string sql = "SELECT category, COUNT(*) as total, SUM(available_copies) as available FROM Books GROUP BY category";
    auto res = db.executeQuery(sql);
    std::cout << "馆藏统计:" << std::endl;
    if (res) {
        while (res->next()) {
            std::cout << "分类: " << res->getString("category")
                      << "，总册数: " << res->getInt("total")
                      << "，可借: " << res->getInt("available") << std::endl;
        }
    }
}

bool ReportGenerator::exportToCSV(const std::string& filename, const std::string& query) {
    std::ofstream file(filename);
    if (!file.is_open()) return false;
    auto res = db.executeQuery(query);
    if (!res) return false;

    // 写入列名（简化：获取第一个结果的元数据）
    sql::ResultSetMetaData* meta = res->getMetaData();
    int cols = meta->getColumnCount();
    for (int i = 1; i <= cols; ++i) {
        if (i > 1) file << ",";
        file << meta->getColumnName(i);
    }
    file << "\n";

    // 写入数据
    while (res->next()) {
        for (int i = 1; i <= cols; ++i) {
            if (i > 1) file << ",";
            file << res->getString(i);
        }
        file << "\n";
    }
    file.close();
    return true;
}