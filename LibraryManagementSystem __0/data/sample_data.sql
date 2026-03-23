USE library;

-- 插入示例管理员
INSERT INTO Administrators (username, password_hash, name, role) VALUES
('admin', SHA2('admin123', 256), '系统管理员', '超级管理员');

-- 插入示例读者
INSERT INTO Readers (card_number, name, gender, phone, email, address, reader_type, max_borrow_limit, registration_date, expiry_date) VALUES
('R2023001', '张三', '男', '13800138000', 'zhangsan@example.com', '北京市海淀区', '学生', 5, '2023-01-01', '2024-12-31'),
('R2023002', '李四', '女', '13900139000', 'lisi@example.com', '上海市浦东区', '教师', 10, '2023-02-01', '2024-12-31');

-- 插入示例图书
INSERT INTO Books (isbn, title, author, publisher, publish_date, category, price, total_copies, available_copies, location) VALUES
('978-7-111-12345-6', 'C++ Primer Plus', 'Stephen Prata', '机械工业出版社', '2020-01-01', '编程', 128.00, 5, 5, 'A-01-01'),
('978-7-302-45678-9', '算法导论', 'Thomas H. Cormen', '清华大学出版社', '2018-05-01', '算法', 99.00, 3, 3, 'A-01-02'),
('978-7-121-34567-8', '深入理解计算机系统', 'Randal E. Bryant', '电子工业出版社', '2019-09-01', '系统', 139.00, 2, 2, 'A-02-01');

-- 插入示例借阅记录
INSERT INTO BorrowRecords (reader_id, book_id, borrow_date, due_date, status) VALUES
(1, 1, '2023-10-01 10:00:00', '2023-10-31', '借出'),
(2, 2, '2023-09-15 14:30:00', '2023-10-15', '已还');