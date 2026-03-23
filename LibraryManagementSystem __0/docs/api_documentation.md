
# 图书馆管理系统 API 文档

## 数据库连接类 Database
- `Database(host, user, pass, db)` - 构造函数
- `connect()` - 连接数据库
- `disconnect()` - 断开连接
- `executeQuery(sql)` - 执行 SELECT 查询，返回 ResultSet
- `executeUpdate(sql)` - 执行 INSERT/UPDATE/DELETE，返回影响行数
- 事务: `beginTransaction()`, `commit()`, `rollback()`

## 图书类 Book
- 属性: bookId, isbn, title, author, publisher, publishDate, category, price, totalCopies, availableCopies, location, status
- `addToDatabase(db)` - 添加图书
- `updateInDatabase(db)` - 更新图书
- `deleteFromDatabase(db)` - 删除图书
- `fetchFromDatabase(db, bookId)` - 静态，根据ID获取图书
- `searchByTitle(db, keyword)` - 按书名搜索
- `searchByAuthor(db, author)` - 按作者搜索

## 读者类 Reader
- 属性: readerId, cardNumber, name, gender, phone, email, address, readerType, maxBorrowLimit, currentBorrowed, registrationDate, expiryDate, status, penalty
- `registerReader(db)` - 注册读者
- `updateInfo(db)` - 更新信息
- `deleteReader(db)` - 删除读者
- `fetchFromDatabase(db, readerId)` - 根据ID获取
- `fetchByCardNumber(db, cardNo)` - 根据借书证号获取
- `hasOverdueBooks(db)` - 检查是否有逾期图书

## 借阅管理 BorrowManager
- `borrowBook(readerId, isbn)` - 借书
- `returnBook(recordId, condition)` - 还书，condition 为 "正常"/"损坏"/"丢失"
- `renewBook(recordId)` - 续借
- `sendDueDateReminders()` - 发送到期提醒
- `sendOverdueNotices()` - 发送逾期通知

## 报表生成 ReportGenerator
- `generateBorrowStatistics(range)` - 借阅统计
- `generatePopularBooksReport(topN)` - 热门图书
- `generateReaderActivityReport()` - 读者活跃度
- `generateOverdueAnalysis()` - 逾期分析
- `generateInventoryReport()` - 馆藏统计
- `exportToCSV(filename, query)` - 导出查询结果为CSV